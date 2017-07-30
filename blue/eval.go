
package blue

import (
    "os"
    "fmt"
    "github.com/Magnus9/blue/token"
    "github.com/Magnus9/blue/interm"
    "github.com/Magnus9/blue/errpkg"
    "github.com/Magnus9/blue/objects"
)
const (
    DIVEOUT_RETURN   = 0
    DIVEOUT_CONTINUE = 1
    DIVEOUT_BREAK    = 2
    DIVEOUT_NONE     = 3
)
// The module map. Holds all loaded modules.
var modules = make(map[string]*objects.BlModuleObject,
                   0)
type Eval struct {
    pathname     string
    root         *interm.Node
    globals      map[string]objects.BlObject
    builtins     map[string]objects.BlObject
    frame        *objects.BlFrame
    tracefunc    func(frame *objects.BlFrame)
    diveout      objects.BlDiveout
    cobj         objects.BlObject
    inFunction   bool
    loopCount    int
}
type tracefunction func(frame *objects.BlFrame)

func Init(argv []string) {
    // Initialize all the type objects.
    objects.BlInitTypes()
    // Initialize the builtins module.
    blInitBuiltins()
    // Initialize the system module(core).
    blInitSystem(argv)
}

func GetModuleMap() map[string]*objects.BlModuleObject {
    return modules
}

func New(pathname string, root *interm.Node,
         globals map[string]objects.BlObject) *Eval {
    eval := &Eval{
        pathname : pathname,
        root     : root,
        globals  : globals,
        builtins : builtins,
        tracefunc: genericTraceFunc,
    }
    return eval
}

/*
 * Run an interpretation on this evaluation context.
 * An evaluation context equals a compiled file.
 */
func (e *Eval) Run() {
    e.evalCode(e.root, nil)
}

func (e *Eval) evalCode(
node *interm.Node,
locals map[string]objects.BlObject) {
    e.frame = objects.NewBlFrame(e.frame, locals,
                                 e.pathname)
    e.exec(node)
    e.frame = e.frame.Prev
}

/*
 * Kind of messy to read over, it takes care of
 * allocating a map filled with varnames => values.
 * It might become the function dispatcher in the
 * future.

 * The important variables:
 * argpos => Points to the stared parameter if it
             exists, otherwise argument length.
 * i      => The parameter iterator.
 * j      => The argument iterator.
 */
func (e *Eval) buildLocals(
f *objects.BlFunctionObject,
args *interm.Node,
self *objects.BlInstanceObject,
) map[string]objects.BlObject {
    /*
     * argpos points to the position in the parameter
     * list where a stared parameter occurs.
     */
    argpos := args.Nchildren
    if self != nil {
        if f.StarParam && f.ParamLen == 0 {
        } else {
            argpos++
        }
    }
    if f.StarParam && argpos > f.ParamLen {
        argpos = argpos - (argpos - f.ParamLen)
    }
    if !e.verifyParamArgCount(f.ParamLen, argpos) {
        return nil
    }
    locals := make(map[string]objects.BlObject, 0)
    list := objects.NewBlList(0)
    var i int
    if self != nil {
        // Place the receiver.
        if f.StarParam && f.ParamLen == 0 {
            list.Append(self)
        } else {
            locals[f.Params[0]] = self
            i++
        }
    }
    var j int
    for ; i < argpos; i++ {
        locals[f.Params[i]] = e.exec(args.Children[j])
        j++
    }
    if f.StarParam {
        for ; j < args.Nchildren; j++ {
            list.Append(e.exec(args.Children[j]))
        }
        locals[f.Params[argpos]] = list
    }
    return locals
}

func (e *Eval) block(node *interm.Node) {
    for _, n := range node.Children {
        e.exec(n)
    }
}

func (e *Eval) verifyParamArgCount(params,
                                   args int) bool {
    if params != args {
        errpkg.SetErrmsg("argument mismatch. Expected" +
                         " (%d), got (%d)", params, args)
        return false
    }
    return true
}

func (e *Eval) exec(node *interm.Node) objects.BlObject {
    e.frame.SetNode(node)
    switch node.NodeType {
        case token.INTERACTIVE:
            for _, n := range node.Children {
                ret := e.exec(n)
                if ret != nil {
                    blPrint(ret)
                }
            }
        case token.FILE_INPUT:
            fallthrough
        case token.CLASSBLOCK:
            fallthrough
        case token.BLOCK:
            e.block(node)
        case token.IMPORT:
            for _, n := range node.Children {
                name, ret := blImportModule(n)
                if ret == nil {
                    goto err
                }
                e.globals[name] = ret
            }
        case token.MAKE_CLASS:
            name := node.Children[0].Str
            obj := e.makeClass(name, node.Children[1],
                               node.Children[2])
            if obj == nil {
                goto err
            }
            e.cobj = nil
            e.set(name, obj)
        case token.MAKE_FUNC:
            name := node.Children[0].Str
            obj := e.makeFunc(name, node.Children[1],
                              node.Children[2])
            if obj == nil {
                goto err
            }
            e.set(name, obj)
        case token.IF:
            e.ifStmt(node)
        case token.WHILE:
            e.whileStmt(node)
        case token.RETURN:
            if !e.inFunction {
                errpkg.SetErrmsg("return outside function")
                goto err
            }
            e.diveout.Type = DIVEOUT_RETURN
            if node.Nchildren > 0 {
                e.diveout.Value = e.exec(node.Children[0])
            }
            panic(e.diveout)
        case token.BREAK:
            if e.loopCount == 0 {
                errpkg.SetErrmsg("break outside loop")
                goto err
            }
            e.diveout.Type = DIVEOUT_BREAK
        case token.CONTINUE:
            if e.loopCount == 0 {
                errpkg.SetErrmsg("continue outside loop")
                goto err
            }
            e.diveout.Type = DIVEOUT_CONTINUE
        case token.ASSIGN:
            ret := e.assign(node)
            if ret == -1 {
                goto err
            }
        case token.AUGASSIGN:
            ret := e.augassign(node.Children[0])
            if ret == -1 {
                goto err
            }
        case token.LOGICAL_AND:
            a := e.exec(node.Children[0])
            if a.BlType().EvalCond(a) {
                b := e.exec(node.Children[1])
                if b.BlType().EvalCond(b) {
                    return objects.BlTrue
                }
            }
            return objects.BlFalse
        case token.LOGICAL_OR:
            a := e.exec(node.Children[0])
            if a.BlType().EvalCond(a) {
                return objects.BlTrue
            }
            b := e.exec(node.Children[1])
            if b.BlType().EvalCond(b) {
                return objects.BlTrue
            }
            return objects.BlFalse
        case token.NOT:
            /*
             * We take advantage that every object
             * contains an EvalCond subroutine here.
             * We simply just reverse the boolean value
             * returned. Maybe move it into BlNumberMethods
             * as a wrapper since its an unary operator?.
             */
            obj := e.exec(node.Children[0])
            ret := obj.BlType().EvalCond(obj)
            if ret {
                return objects.BlFalse
            }
            return objects.BlTrue
        case token.NEGATE:
            obj := e.exec(node.Children[0])
            ret := blNumNegate(obj)
            if ret == nil {
                goto err
            }
            return ret
        case token.COMPL:
            obj := e.exec(node.Children[0])
            ret := blNumCompl(obj)
            if ret == nil {
                goto err
            }
            return ret
        case token.BITWISE_OR:
            a := e.exec(node.Children[0])
            b := e.exec(node.Children[1])
            ret := blNumOr(a, b, "|")
            if ret == nil {
                goto err
            }
            return ret
        case token.BITWISE_AND:
            a := e.exec(node.Children[0])
            b := e.exec(node.Children[1])
            ret := blNumAnd(a, b, "&")
            if ret == nil {
                goto err
            }
            return ret
        case token.XOR:
            a := e.exec(node.Children[0])
            b := e.exec(node.Children[1])
            ret := blNumXor(a, b, "^")
            if ret == nil {
                goto err
            }
            return ret
        case token.LEFTSHIFT:
            a := e.exec(node.Children[0])
            b := e.exec(node.Children[1])
            ret := blNumLshift(a, b, "<<")
            if ret == nil {
                goto err
            }
            return ret
        case token.RIGHTSHIFT:
            a := e.exec(node.Children[0])
            b := e.exec(node.Children[1])
            ret := blNumRshift(a, b, ">>")
            if ret == nil {
                goto err
            }
            return ret
        case token.ADD:
            a := e.exec(node.Children[0])
            b := e.exec(node.Children[1])
            ret := blNumAddition(a, b, "+")
            if ret == nil {
                goto err
            }
            return ret
        case token.SUB:
            a := e.exec(node.Children[0])
            b := e.exec(node.Children[1])
            ret := blNumSubtract(a, b, "-")
            if ret == nil {
                goto err
            }
            return ret
        case token.MUL:
            a := e.exec(node.Children[0])
            b := e.exec(node.Children[1])
            ret := blNumMultiply(a, b, "*")
            if ret == nil {
                goto err
            }
            return ret
        case token.DIV:
            a := e.exec(node.Children[0])
            b := e.exec(node.Children[1])
            ret := blNumDivide(a, b, "/")
            if ret == nil {
                goto err
            }
            return ret
        case token.MODULO:
            a := e.exec(node.Children[0])
            b := e.exec(node.Children[1])
            ret := blNumModulo(a, b, "%")
            if ret == nil {
                goto err
            }
            return ret
        /*
         * Begin comparison operators. All of them have
         * a non-terminal with type token.COMP_OP.
         */
        case token.COMP_OP:
            o := node.Children[0]
            a := e.exec(o.Children[0])
            b := e.exec(o.Children[1])
            res := blCmp(a, b, o.NodeType)
            if res == nil {
                goto err
            }
            return res
        case token.NAME:
            ret := e.get(node.Str)
            if ret == nil {
                goto err
            }
            return ret
        case token.MEMBER:
            obj := e.exec(node.Children[0])
            field := blGetMember(obj, node.Children[1].Str)
            if field == nil {
                goto err
            }
            return field
        case token.MAKE_INSTANCE:
            obj := e.exec(node.Children[0])
            switch t := obj.(type) {
            case *objects.BlClassObject:
                ret := e.newInstance(t, node.Children[1])
                if ret == nil {
                    goto err
                }
                return ret
            case *objects.BlTypeObject:
                ret := e.callType(t, node.Children[1])
                if ret == nil {
                    goto err
                }
                return ret
            }
        case token.CALL:
            var ret objects.BlObject
            obj := e.exec(node.Children[0])
            switch t := obj.(type) {
            case *objects.BlGFunctionObject:
                ret = e.callBuiltin(t, node.Children[1],
                                    nil, false)
                if ret == nil {
                    goto err
                }
            case *objects.BlGMethodObject:
                args := node.Children[1]
                /*
                 * Methods without a receiver expects the first
                 * arg to be the object of the class the method
                 * belongs to. This means that if args.Nchildren
                 * is < 1 or the first child is not the expected
                 * type, the call fails.
                 */
                rcv := t.Self
                if rcv == nil {
                    if args.Nchildren > 0 {
                        rcv = e.exec(args.Children[0])
                    }
                    if rcv == nil || rcv.BlType() != t.Class {
                        tobj := t.Class.(*objects.BlTypeObject)
                        errpkg.SetErrmsg("method '%s' requires a '%s'" +
                                         " object as receiver", t.F.Name,
                                         tobj.Name)
                        goto err
                    }
                    args.Children = args.Children[1:]
                    args.Nchildren--
                }
                ret = e.callBuiltin(t.F, args, rcv, true)
                if ret == nil {
                    goto err
                }
            case *objects.BlFunctionObject:
                locals := e.buildLocals(t, node.Children[1], nil)
                if locals == nil {
                    goto err
                }
                ret = e.callFunction(t, locals)
            case *objects.BlMethodObject:
                locals := e.buildLocals(t.F, node.Children[1], t.Self)
                if locals == nil {
                    goto err
                }
                ret = e.callFunction(t.F, locals)
            default:
                errpkg.SetErrmsg("'%s' object is not callable",
                                 obj.BlType().Name)
                goto err
            }
            if ret == nil {
                goto err
            }
            return ret
        case token.PRINT:
            obj := e.exec(node.Children[0])
            ret := blPrint(obj)
            if ret == -1 {
                goto err
            }
        case token.SUBSCRIPT:
            obj := e.exec(node.Children[0])
            key := e.exec(node.Children[1])
            ret := blGetSeqItem(obj, key)
            if ret == nil {
                goto err
            }
            return ret
        case token.LIST:
            list := objects.NewBlList(0)
            for _, elem := range node.Children {
                list.Append(e.exec(elem))
            }
            return list
        case token.RANGE:
            a := e.exec(node.Children[0])
            b := e.exec(node.Children[1])
            var aIobj, bIobj *objects.BlIntObject
            aIobj, ok := a.(*objects.BlIntObject)
            if !ok {
                goto fail
            }
            bIobj, ok = b.(*objects.BlIntObject)
            if !ok {
                goto fail
            }
            return objects.NewBlRange(aIobj.Value,
                                      bIobj.Value)
fail:
            errpkg.SetErrmsg("types of the range construct" +
                             " must be integers")
            goto err
        case token.STRING:
            value := parseString(node.Str)
            if value == nil {
                goto err
            }
            return objects.NewBlString(*value)
        case token.INTEGER:
            value := parseInt(node.Str)
            if value == -1 {
                goto err
            }
            return objects.NewBlInt(value)
        case token.FLOAT:
            value := parseFloat(node.Str)
            if value == -1.0 {
                goto err
            }
            return objects.NewBlFloat(value)
        case token.TRUE:
            return objects.BlTrue
        case token.FALSE:
            return objects.BlFalse
        case token.NIL:
            return objects.BlNil
        default:
            fmt.Fprintf(os.Stderr, "unrecognized node type" +
                        " (%d)\n", node.NodeType)
            os.Exit(1)
    }
    return nil
err:
    e.tracefunc(e.frame)
    return nil
}

func (e *Eval) get(name string) objects.BlObject {
    var obj objects.BlObject

    if e.cobj != nil {
        obj = blGetMember(e.cobj, name)
    } else {
        if e.frame.Locals != nil {
            obj = e.frame.Locals[name]
        }
    }
    if obj == nil {
        obj = e.globals[name]
        if obj == nil {
            obj = e.builtins[name]
        }
    }    
    if obj == nil {
        errpkg.SetErrmsg("failed to resolve variable" +
                         " '%s'", name)
    }
    return obj
}

func (e *Eval) set(name string, v objects.BlObject) int {
    if e.cobj != nil {
        return blSetMember(e.cobj, v, name)
    }
    if e.frame.Locals != nil {
        e.frame.Locals[name] = v
    } else {
        e.globals[name] = v
    }
    return 0
}

func (e *Eval) makeClass(
name string,
extends, classblock *interm.Node) objects.BlObject {
    var base objects.BlObject
    if extends.Nchildren > 0 {
        base = e.get(extends.Children[0].Str)
        if base == nil {
            return nil
        }
    }
    e.cobj = objects.NewBlClass(name, base)
    e.exec(classblock)

    return e.cobj
}

func (e *Eval) makeFunc(
name string,
paramsNode, block *interm.Node) objects.BlObject {
    var params []string
    for _, i := range paramsNode.Children {
        params = append(params, i.Str)
    }
    starParam := false
    if (paramsNode.Flags & interm.FLAG_STARPARAM) != 0 {
        starParam = true
    } 
    return objects.NewBlFunction(name, params,
                                 paramsNode.Nchildren,
                                 block, starParam) 
}

func (e *Eval) ifStmt(node *interm.Node) {
    var i int
    for i = 0; i < node.Nchildren; i += 3 {
        cond := e.exec(node.Children[i])
        if cond.BlType().EvalCond(cond) {
            e.exec(node.Children[i + 1])
            return
        }
    }
    if i == node.Nchildren {
        e.exec(node.Children[i - 1])
    }
}

func (e *Eval) whileStmt(node *interm.Node) {
    block := node.Children[1]
    e.loopCount++
outer:
    for true {
        cond := e.exec(node.Children[0])
        if !cond.BlType().EvalCond(cond) {
            break
        }
inner:
        for _, stmt := range block.Children {
            e.exec(stmt)
            switch {
                case e.diveout.Type == DIVEOUT_BREAK:
                    e.diveout.Type = DIVEOUT_NONE
                    break outer
                case e.diveout.Type == DIVEOUT_CONTINUE:
                    e.diveout.Type = DIVEOUT_NONE
                    break inner
            }
        }
    }
    e.loopCount--
}

func (e *Eval) assign(node *interm.Node) int {
    left := node.Children[0]
    switch left.NodeType {
        case token.NAME:
            ret := e.set(left.Str, e.exec(node.Children[1]))
            return ret
        case token.MEMBER:
            obj := e.exec(left.Children[0])
            return blSetMember(obj, e.exec(node.Children[1]),
                               left.Children[1].Str)
        case token.SUBSCRIPT:
            obj := e.exec(left.Children[0])
            key := e.exec(left.Children[1])
            return blSetSeqItem(obj, e.exec(node.Children[1]),
                                key)
    }
    return 0
}
// ('+=' EXPR EXPR)
// We can use our exec to our advantage here on the zeroeth child.
// Getting the value through the dispatcher before we case the OP
// in question. Then we can run the operator in question and produce
// the exec RETVAL OP 1nth child.
// In the end we still need to case the zeroeth child, we dont know
// how we are going to save the value (the most important part).
// NAME      = e.set(...)
// DOT       = blSetMember(...)
// SUBSCRIPT = blSetSeqItem(...) 
func (e *Eval) augassign(node *interm.Node) int {
    // Start with dispatching the zeroeth child.
    a := e.exec(node.Children[0])
    // Dispatch the RHS value.
    b := e.exec(node.Children[1])
    // Now case the OP to handle the correct oper.
    var ret objects.BlObject
    switch node.NodeType {
        case token.ASS_BITWISE_OR:
            ret = blNumOr(a, b, "|=")
        case token.ASS_BITWISE_AND:
            ret = blNumAnd(a, b, "&=")
        case token.ASS_XOR:
            ret = blNumXor(a, b, "^=")
        case token.ASS_LEFTSHIFT:
            ret = blNumLshift(a, b, "<<=")
        case token.ASS_RIGHTSHIFT:
            ret = blNumRshift(a, b, ">>=")
        case token.ASS_ADD:
            ret = blNumAddition(a, b, "+=")
        case token.ASS_SUB:
            ret = blNumSubtract(a, b, "-=")
        case token.ASS_MUL:
            ret = blNumMultiply(a, b, "*=")
        case token.ASS_DIV:
            ret = blNumDivide(a, b, "/=")
        case token.ASS_MODULO:
            ret = blNumModulo(a, b, "%=")
    }
    if ret == nil {
        return -1
    }
    left := node.Children[0]
    switch left.NodeType {
        case token.NAME:
            return e.set(left.Str, ret)
        case token.MEMBER:
            obj := e.exec(left.Children[0])
            return blSetMember(obj, ret,
                               left.Children[1].Str)
        case token.SUBSCRIPT:
            obj := e.exec(left.Children[0])
            key := e.exec(left.Children[1])
            return blSetSeqItem(obj, ret, key)
    }
    return 0
}

func (e *Eval) callBuiltin(
f *objects.BlGFunctionObject, args *interm.Node,
rcv objects.BlObject, meth bool) objects.BlObject {
    if (f.Flags & objects.GFUNC_NOARGS) != 0 &&
        args.Nchildren > 0 {
        errpkg.SetErrmsg("%s() takes no arguments",
                         f.Name)
        return nil
    }
    arglist := make([]objects.BlObject, args.Nchildren)
    for i, arg := range args.Children {
        arglist[i] = e.exec(arg)
    }
    if !meth {
        return f.Function(nil, arglist...)
    }
    return f.Function(rcv, arglist...)
}

func (e *Eval) callFunction(f *objects.BlFunctionObject,
locals map[string]objects.BlObject) (obj objects.BlObject) {
    obj = objects.BlNil
    defer func() {
        err := recover()
        if err != nil {
            switch err.(type) {
                case string:
                    panic(err)
                case objects.BlDiveout:
                    if e.diveout.Value != nil {
                        obj = e.diveout.Value
                    }
            }
            e.frame = e.frame.Prev
            e.inFunction = false
            e.loopCount = 0
        }
    }()
    e.inFunction = true
    e.evalCode(f.Block, locals)
    e.inFunction = false
    e.loopCount = 0
    return obj
}

func (e *Eval) newInstance(class *objects.BlClassObject,
                           args *interm.Node) objects.BlObject {
    iobj := objects.NewBlInstance(class)
    mobj := blGetMember(class, "__init__")
    if mobj != nil {
        mobj := mobj.(*objects.BlMethodObject)
        locals := e.buildLocals(mobj.F, args, iobj)
        if locals == nil {
            return nil
        }
        ret := e.callFunction(mobj.F, locals)
        if ret == nil {
            return nil
        }
    }
    return iobj
}

func (e *Eval) callType(obj *objects.BlTypeObject,
                        args *interm.Node) objects.BlObject {
    arglist := make([]objects.BlObject, args.Nchildren)
    for i, arg := range args.Children {
        arglist[i] = e.exec(arg)
    }
    return obj.BlType().Init(obj, arglist...)
}