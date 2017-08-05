
package objects

import (
    "fmt"
    "github.com/Magnus9/blue/errpkg"
)
type BlRangeObject struct {
    header blHeader
    S      int
    E      int
}
func (bro *BlRangeObject) BlType() *BlTypeObject {
    return bro.header.typeobj
}
var blRangeSequence = BlSequenceMethods{
    SqItem: blRangeItem,
    SqSize: blRangeSize,
}
var BlRangeType BlTypeObject

func NewBlRange(s, e int) *BlRangeObject {
    return &BlRangeObject{
        header: blHeader{&BlRangeType},
        S     : s,
        E     : e,
    }
}

func blRangeItem(obj BlObject, num int) BlObject {
    robj := obj.(*BlRangeObject)
    if num < 0 || num >= (robj.E - robj.S) {
        errpkg.SetErrmsg("subscript position out of bounds")
        return nil
    }
    return NewBlInt(int64(robj.S + num))
}

func blRangeSize(obj BlObject) int {
    robj := obj.(*BlRangeObject)
    return robj.E - robj.S
}

func blRangeRepr(obj BlObject) *BlStringObject {
    robj := obj.(*BlRangeObject)
    return NewBlString(fmt.Sprintf("%d..%d", robj.S,
                                   robj.E))
}

func blInitRange() {
    BlRangeType = BlTypeObject{
        header  : blHeader{&BlTypeType},
        Name    : "range",
        Repr    : blRangeRepr,
        Sequence: &blRangeSequence,
    }
}