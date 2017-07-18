
package objects

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

func blBoolRepr(obj BlObject) BlObject {
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

func blInitBool() {
    BlBoolType = BlTypeObject{
        header  : blHeader{&BlTypeType},
        Name    : "bool",
        Repr    : blBoolRepr,
        EvalCond: blBoolEvalCond,
        Compare : blBoolCompare,
    }
    blTypeFinish(&BlBoolType)
}