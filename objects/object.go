
package objects

import (
    "github.com/Magnus9/blue/errpkg"
)
type BlObject interface {
    BlType() *BlTypeObject
}
type blHeader struct {
    typeobj *BlTypeObject
}
var (
    BlTrue  = NewBlBool(true)
    BlFalse = NewBlBool(false)
    BlNil   = NewBlNil()
)

type reprfunc func(BlObject) BlObject
type fillfunc func(BlObject, map[string]struct{}, int) int
type getterfunc func(BlObject, string) BlObject
type setterfunc func(BlObject, string, BlObject) int
type evalcondfunc func(BlObject) bool
type initfunc func(*BlTypeObject, ...BlObject) BlObject
type compfunc func(BlObject, BlObject) int

type unaryfunc func(BlObject) BlObject
type binaryfunc func(BlObject, BlObject) BlObject

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

func blParseArguments(format string, args []BlObject,
                      values ...interface{}) int {
    fmtLen := len(format)
    if fmtLen != len(args) {
        errpkg.SetErrmsg("argument mismatch. Expected (%d)" +
                         " , got (%d)", fmtLen, len(args))
        return -1
    }
    for i := 0; i < fmtLen; i++ {
        switch ch := format[i]; ch {
            case 'i':
                iObj, ok := args[i].(*BlIntObject)
                if !ok {
                    errpkg.SetErrmsg("expected int object")
                    return -1
                }
                iptr := values[i].(*int64)
                *iptr = iObj.Value
            case 's':
                sObj, ok := args[i].(*BlStringObject)
                if !ok {
                    errpkg.SetErrmsg("expected string object")
                    return -1
                }
                sptr := values[i].(*string)
                *sptr = sObj.Value
        }
    }
    return 0
}

func blTypeFinish(typeobj *BlTypeObject) {
    typeobj.members = make(map[string]BlObject)
    if typeobj.methods != nil {
        for i := 0; i < len(typeobj.methods); i++ {
            m := &typeobj.methods[i]
            typeobj.members[m.name] = m
        }
    }
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