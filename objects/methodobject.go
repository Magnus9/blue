
package objects

import (
    "fmt"
)
type BlMethodObject struct {
    header blHeader
    class  *BlClassObject
    Self   *BlInstanceObject
    F      *BlFunctionObject
}
func (bmo *BlMethodObject) BlType() *BlTypeObject {
    return bmo.header.typeobj
}
var BlMethodType BlTypeObject

func NewBlMethod(class *BlClassObject,
self *BlInstanceObject, f *BlFunctionObject) BlObject {
    return &BlMethodObject{
        header: blHeader{&BlMethodType},
        class : class,
        Self  : self,
        F     : f,
    }
}

func blMethodRepr(obj BlObject) *BlStringObject {
    mobj := obj.(*BlMethodObject)
    str := fmt.Sprintf("<method '%s.%s', params=%d>",
                       mobj.class.name, mobj.F.name,
                       mobj.F.ParamLen)
    return NewBlString(str)
}

func blMethodEvalCond(obj BlObject) bool {
    mobj := obj.(*BlMethodObject)
    return mobj.F.BlType().EvalCond(mobj.F)
}

func blInitMethod() {
    BlMethodType = BlTypeObject{
        header: blHeader{&BlTypeType},
        Name    : "method",
        Repr    : blMethodRepr,
        EvalCond: blMethodEvalCond,
    }
}