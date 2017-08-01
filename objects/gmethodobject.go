
package objects

import (
    "fmt"
)
type BlGMethodObject struct {
    header   blHeader
    Class    BlObject
    Self     BlObject
    F        *BlGFunctionObject
}
func (bgmo *BlGMethodObject) BlType() *BlTypeObject {
    return bgmo.header.typeobj
}
var BlGMethodType BlTypeObject

func newBlGMethod(class BlObject, self BlObject,
                  f *BlGFunctionObject) BlObject {
    return &BlGMethodObject{
        header  : blHeader{&BlGMethodType},
        Class   : class,
        Self    : self,
        F       : f,
    }
}

func blGMethodRepr(obj BlObject) *BlStringObject {
    mobj := obj.(*BlGMethodObject)
    typeobj := mobj.Class.(*BlTypeObject)
    str := fmt.Sprintf("<builtin-method '%s.%s'>",
                       typeobj.Name, mobj.F.Name)
    return NewBlString(str)
}

func blGMethodEvalCond(obj BlObject) bool {
    mobj := obj.(*BlGMethodObject)
    return mobj.F.BlType().EvalCond(mobj.F)
}

func blInitGMethod() {
    BlGMethodType = BlTypeObject{
        header  : blHeader{&BlTypeType},
        Name    : "builtin-method",
        Repr    : blGMethodRepr,
        EvalCond: blGMethodEvalCond,
    }
}