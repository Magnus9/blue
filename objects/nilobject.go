
package objects

import (

)
type BlNilObject struct {
    header blHeader
}
func (bno *BlNilObject) BlType() *BlTypeObject {
    return bno.header.typeobj
}
var BlNilType BlTypeObject

func NewBlNil() BlObject {
    return &BlNilObject{
        header: blHeader{&BlNilType},
    }
}

func blNilRepr(obj BlObject) *BlStringObject {
    return NewBlString("nil")
}

func blNilEvalCond(obj BlObject) bool {
    return false
}

func blInitNil() {
    BlNilType = BlTypeObject{
        header  : blHeader{&BlTypeType},
        Name    : "nil",
        Repr    : blNilRepr,
        EvalCond: blNilEvalCond,
    }
}