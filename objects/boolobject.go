
package objects

import "github.com/Magnus9/blue/errpkg"

type BlBoolObject struct {
    header blHeader
    value  bool
}
func (bbo *BlBoolObject) BlType() *BlTypeObject {
    return bbo.header.typeobj
}
var BlBoolType BlTypeObject

func NewBlBool(value bool) BlObject {
    return &BlBoolObject{
        header: blHeader{&BlBoolType},
        value : value,
    }
}

func blBoolRepr(obj BlObject) *BlStringObject {
    bobj := obj.(*BlBoolObject)
    if bobj.value == true {
        return NewBlString("true")
    }
    return NewBlString("false")
}

func blBoolEvalCond(obj BlObject) bool {
    bobj := obj.(*BlBoolObject)
    if bobj.value == true {
        return true
    }
    return false
}

func blBoolCompare(a, b BlObject) int {
    aBobj := a.(*BlBoolObject)
    bBobj := b.(*BlBoolObject)
    switch {
    case aBobj.value == bBobj.value:
        return 0
    case aBobj.value == false && bBobj.value == true:
        return -1
    default:
        return 1
    }
}

/*
 * Patch in the use of EvalCond in later stages.
 * For now there is enough conditions as it is,
 * so keep it simple.
 */
func blBoolInit(obj *BlTypeObject,
                args ...BlObject) BlObject {
    var arg BlObject
    if blParseArguments("|o", args, &arg) == -1 {
        return nil
    }
    if arg == nil {
        return NewBlBool(false)
    }
    switch t := arg.(type) {
        case *BlBoolObject:
            return NewBlBool(t.value)
        case *BlIntObject:
            if t.Value != 0 {
                return NewBlBool(true)
            }
            return NewBlBool(false)
        default:
            errpkg.SetErrmsg("expected bool or int")
    }
    return nil
}

func blInitBool() {
    BlBoolType = BlTypeObject{
        header  : blHeader{&BlTypeType},
        Name    : "bool",
        Repr    : blBoolRepr,
        EvalCond: blBoolEvalCond,
        Compare : blBoolCompare,
        Init    : blBoolInit,
    }
    blTypeFinish(&BlBoolType)
}