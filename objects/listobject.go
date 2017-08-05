
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
func (blo *BlListObject) GetList() []BlObject {
    return blo.list
}
var blListSequence = BlSequenceMethods{
    SqItem      : blListItem,
    SqAssItem   : blListAssItem,
    SqRepeat    : blListRepeat,
    SqSlice     : blListSlice,
    SqAssSlice  : blListAssSlice,
    SqSize      : blListSize,
}
var blListMethods = []BlGFunctionObject{
    NewBlGFunction("append",  listAppend,  GFUNC_VARARGS),
    NewBlGFunction("prepend", listPrepend, GFUNC_VARARGS),
    NewBlGFunction("insert",  listInsert,  GFUNC_VARARGS),
    NewBlGFunction("trunc",   listTrunc,   GFUNC_NOARGS ),
    NewBlGFunction("reverse", listReverse, GFUNC_NOARGS ),
}
var BlListType BlTypeObject

func NewBlList(lsize int) *BlListObject {
    return &BlListObject{
        header: blHeader{&BlListType},
        list  : make([]BlObject, lsize),
        lsize : lsize,
    }
}

func blListItem(obj BlObject, num int) BlObject {
    lobj := obj.(*BlListObject)
    if num >= lobj.lsize || num < 0 {
        errpkg.SetErrmsg("subscript position out of bounds")
        return nil
    }
    return lobj.list[num]
}

func blListAssItem(obj, value BlObject, num int) int {
    lobj := obj.(*BlListObject)
    if num > lobj.lsize || num < 0 {
        errpkg.SetErrmsg("subscript position out of bounds")
        return -1
    }
    lobj.list[num] = value
    return 0
}

func blListRepeat(a, b BlObject) BlObject {
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

func blListSlice(obj BlObject, s, e int) BlObject {
    lobj := obj.(*BlListObject)
    if s < 0 {
        s = 0
    } else if s > lobj.lsize {
        s = lobj.lsize
    }
    if e < s {
        e = s
    } else if e > lobj.lsize {
        e = lobj.lsize
    }
    list := NewBlList(0)
    for i := s; i < e; i++ {
        list.Append(lobj.list[i])
    }
    return list
}

func blListAssSlice(obj, value BlObject, s, e int) int {
    lobj2, ok := value.(*BlListObject)
    if !ok {
        errpkg.SetErrmsg("expected list for slice assigment" +
                         " found '%s'", value.BlType().Name)
        return -1
    }
    lobj := obj.(*BlListObject)
    if s < 0 {
        s = 0
    } else if s > lobj.lsize {
        s = lobj.lsize
    }
    if e < s {
        e = s
    } else if e > lobj.lsize {
        e = lobj.lsize
    }
    /*
     * Algorithm below might work better with for loops;
     * removing items, changing item locations etc..
     */
    siz := lobj2.lsize - (e - s)
    if siz < 0 {
        lobj.list = append(lobj.list[s:e + siz], lobj.list[e:]...)
        for i := s; i < e + siz; i++ {
            lobj.list[i] = lobj2.list[i - s]
        }
    } else {
        lobj.list = append(lobj.list[:s], append(lobj2.list,
                           lobj.list[e:]...)...)
    }
    lobj.lsize += siz
    return 0
}

func blListSize(obj BlObject) int {
    return obj.(*BlListObject).lsize
}

func blListRepr(obj BlObject) *BlStringObject {
    lobj := obj.(*BlListObject)
    
    var buf bytes.Buffer
    buf.WriteByte('[')
    for i, elem := range lobj.list {
        if i > 0 {
            buf.WriteString(", ")
        }
        sobj := elem.BlType().Repr(elem)
        buf.WriteString(sobj.Value)
    }
    buf.WriteByte(']')
    return NewBlString(buf.String())
}

func blListGetMember(obj BlObject, name string) BlObject {
    return genericGetMember(obj.BlType(), name, obj)
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

/*
 * The list's constructor takes an object that is
 * iterable (it contains SeqItem and SeqSize).
 * Later it will instead dispatch into a subroutine
 * that takes care of running the 'for' construct.
 */
func blListInit(obj *BlTypeObject, args ...BlObject) BlObject {
    var arg BlObject
    if blParseArguments("|o", args, &arg) == -1 {
        return nil
    }
    lobj := NewBlList(0)
    if arg == nil {
        return lobj
    }
    typeobj := arg.BlType()
    if seq := typeobj.Sequence; seq != nil {
        if seq.SqItem == nil || seq.SqSize == nil {
            goto err
        }
        for i := 0; i < seq.SqSize(arg); i++ {
            lobj.Append(seq.SqItem(arg, i))
        }
        return lobj
    }
err:
    errpkg.SetErrmsg("'%s' object is not iterable",
                     typeobj.Name)
    return nil
}

/*
 * The beginning of list methods..
 */
func listAppend(self BlObject, args ...BlObject) BlObject {
    var obj BlObject
    if blParseArguments("o", args, &obj) == -1 {
        return nil
    }
    lobj := self.(*BlListObject)
    lobj.list = append(lobj.list, obj)
    lobj.lsize++

    return BlNil
}

func listPrepend(self BlObject, args ...BlObject) BlObject {
    var obj BlObject
    if blParseArguments("o", args, &obj) == -1 {
        return nil
    }
    lobj := self.(*BlListObject)
    lobj.list = append([]BlObject{obj}, lobj.list...)
    lobj.lsize++

    return BlNil
}

func listInsert(self BlObject, args ...BlObject) BlObject {
    var obj BlObject
    var pos int64
    if blParseArguments("oi", args, &obj, &pos) == -1 {
        return nil
    }
    lobj := self.(*BlListObject)
    if pos < 0 || pos >= int64(lobj.lsize) {
        errpkg.SetErrmsg("position out of bounds")
        return nil
    }
    lobj.list[pos] = obj

    return BlNil
}

func listTrunc(self BlObject, args ...BlObject) BlObject {
    lobj := self.(*BlListObject)
    lobj.list  = make([]BlObject, 0)
    lobj.lsize = 0
    return BlNil
}

func listReverse(self BlObject, args ...BlObject) BlObject {
    lobj := self.(*BlListObject)
    list := make([]BlObject, lobj.lsize)
    var j int
    for i := lobj.lsize - 1; i >= 0; i-- {
        list[j] = lobj.list[i]
        j++
    }
    lobj.list = list
    return BlNil
}

func blInitList() {
    BlListType = BlTypeObject{
        header   : blHeader{&BlTypeType},
        Name     : "list",
        Repr     : blListRepr,
        GetMember: blListGetMember,
        EvalCond : blListEvalCond,
        Compare  : blListCompare,
        Init     : blListInit,
        Sequence : &blListSequence,
        methods  : blListMethods,
    }
    blTypeFinish(&BlListType)
}