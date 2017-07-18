
package objects

import (
    "fmt"
)
type BlGMethodObject struct {
    header   blHeader
    class    BlObject
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
        class   : class,
        Self    : self,
        F       : f,
    }
}

func blGMethodRepr(obj BlObject) BlObject {
    mobj := obj.(*BlGMethodObject)
    
    params := mobj.F.Params
    if mobj.Self == nil {
        params++   
    }
    typeobj := mobj.class.(*BlTypeObject)
    str := fmt.Sprintf("<builtin-method '%s.%s', params=%d>",
                       typeobj.Name, mobj.F.name,
                       params)
    return NewBlString(str)
}

func blGMethodEvalCond(obj BlObject) bool {
    mobj := obj.(*BlGMethodObject)

    params := mobj.F.Params
    if mobj.Self != nil {
        params++
    }
    if params > 0 {
        return true
    }
    return false
}

func blInitGMethod() {
    BlGMethodType = BlTypeObject{
        header  : blHeader{&BlTypeType},
        Name    : "builtin-method",
        Repr    : blGMethodRepr,
        EvalCond: blGMethodEvalCond,
    }
}