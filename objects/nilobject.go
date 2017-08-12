
package objects

import (
    "fmt"
    "strconv"
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

func blNilHash(obj BlObject) int64 {
    /*
     * Hash the address of the nil object. Yes
     * this can actually crash with other hashes..
     */
    hashstr := fmt.Sprintf("%p", obj.(*BlNilObject))
    hash, _ := strconv.ParseInt(hashstr, 10, 64)
    return hash
}

func blInitNil() {
    BlNilType = BlTypeObject{
        header  : blHeader{&BlTypeType},
        Name    : "nil",
        Repr    : blNilRepr,
        EvalCond: blNilEvalCond,
        hash    : blNilHash,
    }
}