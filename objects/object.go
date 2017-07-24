
package objects

import (
    "strings"
    "github.com/Magnus9/blue/errpkg"
)
type BlObject interface {
    BlType() *BlTypeObject
}
type blHeader struct {
    typeobj *BlTypeObject
}
const (
    T_STRING = iota
    T_INT
    T_FLOAT
    T_BOOL
    T_NIL
)
var (
    BlTrue  = NewBlBool(true)
    BlFalse = NewBlBool(false)
    BlNil   = NewBlNil()
)

type reprfunc func(BlObject) *BlStringObject
type fillfunc func(BlObject, map[string]struct{}, int) int
type getterfunc func(BlObject, string) BlObject
type setterfunc func(BlObject, string, BlObject) int
type evalcondfunc func(BlObject) bool
type initfunc func(*BlTypeObject, ...BlObject) BlObject
type compfunc func(BlObject, BlObject) int

type unaryfunc func(BlObject) BlObject
type binaryfunc func(BlObject, BlObject) BlObject

type BlFields struct {
    name      string
    blType    int
    value     interface{}
}

type BlNumberMethods struct {
    NumNeg    unaryfunc
    NumCompl  unaryfunc
    NumOr     binaryfunc
    NumAnd    binaryfunc
    NumXor    binaryfunc
    NumLshift binaryfunc
    NumRshift binaryfunc
    NumAdd    binaryfunc
    NumSub    binaryfunc
    NumMul    binaryfunc
    NumDiv    binaryfunc
    NumMod    binaryfunc
    NumCoerce func(*BlObject, *BlObject) int
}
type BlSequenceMethods struct {
    SeqItem       func(BlObject, int) BlObject
    SeqAssItem    func(BlObject, BlObject, int) int
    SeqConcat     func(BlObject, BlObject) BlObject
    SeqRepeat     func(BlObject, BlObject) BlObject
    SeqSize       func(BlObject) int
}
/*
 * The implementation object of a datatype.
 * Every Blue object is attached to one of these.
 */
type BlTypeObject struct {
    header      blHeader
    Name        string
    Repr        reprfunc
    Fill        fillfunc
    GetMember   getterfunc
    SetMember   setterfunc
    EvalCond    evalcondfunc
    Compare     compfunc
    Init        initfunc
    Numbers     *BlNumberMethods
    Sequence    *BlSequenceMethods
    methods     []BlGFunctionObject
    fields      []BlFields
    members     map[string]BlObject
    base        *BlTypeObject
}
func (bio *BlTypeObject) BlType() *BlTypeObject {
    return bio.header.typeobj
}

func genericGetMember(typeobj *BlTypeObject,
                      name string, self BlObject) BlObject {
    ret := locate(typeobj, name)
    if ret == nil {
        errpkg.SetErrmsg("'%s' object has no member '%s'",
                         typeobj.Name, name)
        return nil
    }
    f, ok := ret.(*BlGFunctionObject)
    if ok {
        return newBlGMethod(typeobj, self, f)
    }
    return ret
}

func locate(typeobj *BlTypeObject, name string) BlObject {
    ret, ok := typeobj.members[name]
    if !ok && typeobj.base != nil {
        ret = locate(typeobj.base, name)
    }
    return ret
}

/*
 * Used to parse arguments for builtin functions.
 */
func blParseArguments(fmts string, args []BlObject,
                      values ...interface{}) int {
    var reqLen, maxLen int
    var fmtLen = len(fmts)
    pos := strings.IndexByte(fmts, '|')
    if pos >= 0 {
        reqLen = pos
        maxLen = fmtLen - 1
    } else {
        reqLen = fmtLen
        maxLen = reqLen
    }
    arglen := len(args)
    /*
     * This algorithm might be pretty shabby. But it yielded
     * the shortest amount of code i found. The last case
     * will also run if pos >= 0 and arglen >= maxLen or
     * arglen <= reqLen1.
     */
    switch {
        case arglen > maxLen && pos >= 0:
            errpkg.SetErrmsg("expected at most (%d) arguments" +
                             ", got (%d)", maxLen, arglen)
            return -1
        case arglen < reqLen && pos >= 0:
            errpkg.SetErrmsg("expected at least (%d) arguments" +
                             ", got (%d)", reqLen, arglen)
            return -1
        case arglen > maxLen || arglen < reqLen:
            errpkg.SetErrmsg("expected exactly (%d) arguments" +
                             ", got (%d)", reqLen, arglen)
            return -1
    }
    /*
     * Make sure that maxLen is equal to the length of the
     * values passed. If the wrong value type is passed,
     * just continue the for loop below, there is just
     * too many errors to report..
     */
    valueLen := len(values)
    if maxLen != valueLen {
        errpkg.InternError("expected exactly (%d) values" +
                           ", got (%d)", maxLen, valueLen)
    }
    var argpos int = -1
    for i := 0; i < fmtLen; i++ {
        argpos++
        if argpos >= arglen {
            break
        }
        switch ch := fmts[i]; ch {
        case 's':
            sobj, ok := args[argpos].(*BlStringObject)
            if !ok {
                errpkg.SetErrmsg("expected string")
                return -1
            }
            sval, ok := values[argpos].(*string)
            if !ok {
                continue
            }
            *sval = sobj.Value
        case 'i':
            iobj, ok := args[argpos].(*BlIntObject)
            if !ok {
                errpkg.SetErrmsg("expected integer")
                return -1
            }
            ival, ok := values[argpos].(*int64)
            if !ok {
                continue
            }
            *ival = iobj.Value
        case 'f':
            fobj, ok := args[argpos].(*BlFloatObject)
            if !ok {
                errpkg.SetErrmsg("expected float")
                return -1
            }
            fval, ok := values[argpos].(*float64)
            if !ok {
                continue
            }
            *fval = fobj.value
        case 'b':
            bobj, ok := args[argpos].(*BlBoolObject)
            if !ok {
                errpkg.SetErrmsg("expected boolean")
                return -1
            }
            bval, ok := values[argpos].(*bool)
            if !ok {
                continue
            }
            *bval = bobj.value
        case 'o':
            oval, ok := values[argpos].(*BlObject)
            if !ok {
                continue
            }
            *oval = args[argpos]
        case '|':
            argpos--
        }
    }
    return 0
}

func blTypeFinish(typeobj *BlTypeObject) {
    m := make(map[string]BlObject)
    if typeobj.methods != nil {
        blInsertFunctions(typeobj.methods, m)
    }
    if typeobj.fields  != nil {
        blInsertFields(typeobj.fields, m)
    }
    typeobj.members = m
}

/*
 * Insert a function into a map that either belongs
 * to a module object or a type object.
 */
func blInsertFunctions(funcs []BlGFunctionObject,
                       m map[string]BlObject) {
    for i := 0; i < len(funcs); i++ {
        f := &funcs[i]; m[f.Name] = f       
    }
}

/*
 * Insert a field into a map that either belongs to
 * a module object or a type object. It is done this
 * way to ensure rules on what can be added as a field.
 * i.e, adding a method as a field is not allowed, since
 * methods are represented with typeobj.methods.
 */
func blInsertFields(fields []BlFields, m map[string]BlObject) {
    for i := 0; i < len(fields); i++ {
        f := &fields[i]
        switch f.blType {
            case T_STRING:
                m[f.name] = NewBlString(f.value.(string))
            case T_INT:
                m[f.name] = NewBlInt(f.value.(int64))
            case T_FLOAT:
                m[f.name] = NewBlFloat(f.value.(float64))
            case T_BOOL:
                m[f.name] = NewBlBool(f.value.(bool))
            case T_NIL:
                m[f.name] = NewBlNil()
            default:
                errpkg.InternError("trying to add field with" +
                                   " unknown type (%d)", f.blType)
        }
    }
}

/*
 * BlNumCoerce and BlCompare is defined here because
 * pkg blue requires the object pkg, and visa versa
 * because of calling into BlCompare on sequence objects
 * for ordering.
 */
func BlNumCoerce(a, b *BlObject) int {
    aTobj := (*a).BlType()
    bTobj := (*b).BlType()
    if aTobj == bTobj {
        return 0
    }
    if aTobj.Numbers.NumCoerce != nil {
        ret := aTobj.Numbers.NumCoerce(a, b)
        if ret == 0 {
            return ret
        }
    }
    if bTobj.Numbers != nil && bTobj.Numbers.NumCoerce != nil {
        ret := bTobj.Numbers.NumCoerce(b, a)
        if ret == 0 {
            return ret
        }
    }
    return -1
}

/*
 * Small and simple comparison function that returns
 * < 0 for LT, > 0 for GT and 0 for EQ. It will grow
 * as the language grows.. For now it does its job.
 */
func BlCompare(a, b BlObject) int {
    if a == b {
        return 0
    }
    aTobj := a.BlType()
    bTobj := b.BlType()
    if aTobj != bTobj {
        if aTobj.Numbers != nil {
            if BlNumCoerce(&a, &b) == 0 {
                aTobj = a.BlType()
                if fn := aTobj.Compare; fn != nil {
                    return fn(a, b)
                }
            }
        }
        goto out
    }
    if fn := aTobj.Compare; fn != nil {
        return fn(a, b)
    }
out:
    errpkg.SetErrmsg("types cannot be ordered, '%s' and '%s'",
                     aTobj.Name, bTobj.Name)
    return -2
}

func BlInitTypes() {
    // Initialize the string type.
    blInitString()
    // Initialize the int type.
    blInitInt()
    // Initialize the float type.
    blInitFloat()
    // Initialize the list type.
    blInitList()
    // Initialize the range type.
    blInitRange()
    // Initialize the file type.
    blInitFile()
    // Initialize the socket type.
    blInitSocket()
    // Initialize the bool type.
    blInitBool()
    // Initialize the nil type.
    blInitNil()
    // Initialize the class type.
    blInitClass()
    // Initialize the instance type.
    blInitInstance()
    // Initialize the method type.
    blInitMethod()
    // Initialize the function type.
    blInitFunction()
    // Initialize the gmethod type.
    blInitGMethod()
    // Initialize the gfunction type.
    blInitGFunction()
    // Initialize the `type' type.
    blInitType()
}