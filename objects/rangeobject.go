
package objects

import (
    "fmt"
)
type BlRangeObject struct {
    header blHeader
    S      int
    E      int
}
func (bro *BlRangeObject) BlType() *BlTypeObject {
    return bro.header.typeobj
}
var BlRangeType BlTypeObject

func NewBlRange(s, e int) *BlRangeObject {
    return &BlRangeObject{
        header: blHeader{&BlRangeType},
        S     : s,
        E     : e,
    }
}

func blRangeRepr(obj BlObject) *BlStringObject {
    robj := obj.(*BlRangeObject)

    return NewBlString(fmt.Sprintf("%d..%d", robj.S,
                                   robj.E))
}

func blInitRange() {
    BlRangeType = BlTypeObject{
        header: blHeader{&BlTypeType},
        Name  : "range",
        Repr  : blRangeRepr,
    }
}