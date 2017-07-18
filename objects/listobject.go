
package objects

import (
    "bytes"
    "github.com/Magnus9/blue/errpkg"
)
const LIST_MAX = 0x00ffffff

type BlListObject struct {
    header blHeader
    list   []BlObject
    lsize  int
}
func (blo *BlListObject) BlType() *BlTypeObject {
    return blo.header.typeobj
}
func (blo *BlListObject) Append(obj BlObject) {
    blo.list = append(blo.list, obj)
    blo.lsize++
}
var blListSequence = BlSequenceMethods{
    SeqItem      : blListSeqItem,
    SeqAssItem   : blListSeqAssItem,
    SeqRepeat    : blListSeqRepeat,
}
var BlListType BlTypeObject

func NewBlList(lsize int) *BlListObject {
    return &BlListObject{
        header: blHeader{&BlListType},
        list  : make([]BlObject, lsize),
        lsize : lsize,
    }
}

func blListSeqItem(obj BlObject, num int) BlObject {
    lobj := obj.(*BlListObject)
    if num > lobj.lsize || num < 0 {
        errpkg.SetErrmsg("subscript position out of bounds")
        return nil
    }
    return lobj.list[num]
}

func blListSeqAssItem(obj, value BlObject, num int) int {
    lobj := obj.(*BlListObject)
    if num > lobj.lsize || num < 0 {
        errpkg.SetErrmsg("subscript position out of bounds")
        return -1
    }
    lobj.list[num] = value
    return 0
}

func blListSeqRepeat(a, b BlObject) BlObject {
    iobj, ok := b.(*BlIntObject)
    if !ok {
        errpkg.SetErrmsg("cant multiply sequence with" +
                         " non-integer")
        return nil
    }
    lobj := a.(*BlListObject)
    if iobj.Value == 1 {
        return a
    }
    size := lobj.lsize * int(iobj.Value)
    if size > LIST_MAX {
        errpkg.SetErrmsg("repeated list became too large")
        return nil
    }
    ret := NewBlList(size)
    for i := 0; i < size; i += lobj.lsize {
        copy(ret.list[i:], lobj.list)
    }
    return ret
}

func blListRepr(obj BlObject) BlObject {
    lobj := obj.(*BlListObject)
    
    var buf bytes.Buffer
    buf.WriteByte('[')
    for i, elem := range lobj.list {
        if i > 0 {
            buf.WriteString(", ")
        }
        sobj := elem.BlType().Repr(elem).(*BlStringObject)
        buf.WriteString(sobj.Value)
    }
    buf.WriteByte(']')
    return NewBlString(buf.String())
}

func blListEvalCond(obj BlObject) bool {
    lobj := obj.(*BlListObject)
    if lobj.lsize > 0 {
        return true
    }
    return false
}

func blListCompare(a, b BlObject) int {
    aLobj := a.(*BlListObject)
    bLobj := b.(*BlListObject)
    for i := 0; i < aLobj.lsize && i < bLobj.lsize; i++ {
        ret := BlCompare(aLobj.list[i], bLobj.list[i])
        if ret != 0 {
            return ret
        }
    }
    /*
     * If i becomes exhausted it means the lists are equal
     * up to one point, but their length might still be
     * different.
     */
    switch {
    case aLobj.lsize < bLobj.lsize:
        return -1
    case aLobj.lsize > bLobj.lsize:
        return 1
    default:
        return 0
    }
}

func blInitList() {
    BlListType = BlTypeObject{
        header  : blHeader{&BlTypeType},
        Name    : "list",
        Repr    : blListRepr,
        EvalCond: blListEvalCond,
        Compare : blListCompare,
        Sequence: &blListSequence,
    }
}