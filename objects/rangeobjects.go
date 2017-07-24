
package objects

import (
    "fmt"
)
type BlRangeObject struct {
    header blHeader
    s      int64
    e      int64
}
func (bro *BlRangeObject) BlType() *BlTypeObject {
    return bro.header.typeobj
}
var BlRangeType BlTypeObject

func NewBlRange(s, e int64) *BlRangeObject {
    return &BlRangeObject{
        header: blHeader{&BlRangeType},
        s     : s,
        e     : e,
    }
}

func blRangeRepr(obj BlObject) *BlStringObject {
    robj := obj.(*BlRangeObject)

    return NewBlString(fmt.Sprintf("%d..%d", robj.s,
                                   robj.e))
}

func blInitRange() {
    BlRangeType = BlTypeObject{
        header: blHeader{&BlTypeType},
        Name  : "range",
        Repr  : blRangeRepr,
    }
}