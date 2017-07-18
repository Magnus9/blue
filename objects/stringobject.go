
package objects

import (
    "bytes"
    "github.com/Magnus9/blue/errpkg"
)
const STRING_MAX = 0x00ffffff

type BlStringObject struct {
    header blHeader
    Value  string
    vsize  int
}
func (bso *BlStringObject) BlType() *BlTypeObject {
    return bso.header.typeobj
}
var blStringSequence = BlSequenceMethods{
    SeqItem  : blStringSeqItem,
    SeqConcat: blStringSeqConcat,
    SeqRepeat: blStringSeqRepeat,
}
var blStringMethods = []BlGFunctionObject {
    NewBlGFunction("concat", blStringConcat, 1),
    NewBlGFunction("toupper", blStringToUpper, 0),
    NewBlGFunction("tolower", blStringToLower, 0),
}
var BlStringType BlTypeObject

func NewBlString(value string) *BlStringObject {
    return &BlStringObject{
        header  : blHeader{&BlStringType},
        Value   : value,
        vsize   : len(value),
    }
}

func blStringSeqItem(obj BlObject, num int) BlObject {
    sobj := obj.(*BlStringObject)
    if num >= sobj.vsize || num < 0 {
        errpkg.SetErrmsg("subscript position out of bounds")
        return nil
    }
    return NewBlString(string(sobj.Value[num]))
}

func blStringSeqConcat(a, b BlObject) BlObject {
    sobj := a.(*BlStringObject)
    t, ok := b.(*BlStringObject)
    if !ok {
        errpkg.SetErrmsg("cannot add '%s' to string",
                         b.BlType().Name)
        return nil
    }
    return NewBlString(sobj.Value + t.Value)
}

func blStringSeqRepeat(a, b BlObject) BlObject {
    iobj, ok := b.(*BlIntObject)
    if !ok {
        errpkg.SetErrmsg("cant multiply sequence with" + 
                         " non-integer")
        return nil
    }
    sobj := a.(*BlStringObject)
    if iobj.Value == 1 {
        return a
    }
    size := sobj.vsize * int(iobj.Value)
    /*
     * Maximum string size for now is 24bits. Funny thing
     * is to replace the buffer writing to string
     * concentation using the '+=' operator.
     */
    if size > STRING_MAX {
        errpkg.SetErrmsg("repeated string became too large")
        return nil
    }
    var buf bytes.Buffer
    for i := 0; i < size; i += sobj.vsize {
        buf.WriteString(sobj.Value)
    }
    return NewBlString(buf.String())
}

func blStringRepr(obj BlObject) BlObject {
    return obj
}

func blStringGetMember(obj BlObject,
                       name string) BlObject {
    return genericGetMember(obj.BlType(), name, obj)
}

func blStringEvalCond(obj BlObject) bool {
    sobj := obj.(*BlStringObject)
    if sobj.vsize > 0 {
        return true
    }
    return false
}

func blStringCompare(a, b BlObject) int {
    aSobj := a.(*BlStringObject)
    bSobj := b.(*BlStringObject)
    for i := 0; i < aSobj.vsize && i < bSobj.vsize; i++ {
        v := int(aSobj.Value[i]) - int(bSobj.Value[i])
        if v != 0 {
            return v
        }
    }
    switch {
    case aSobj.vsize < bSobj.vsize:
        return -1
    case aSobj.vsize > bSobj.vsize:
        return 1
    default:
        return 0
    }
}

func blStringInit(obj *BlTypeObject,
                  args ...BlObject) BlObject {
    var str string
    if blParseArguments("s", args, &str) == -1 {
        return nil
    }
    return NewBlString(str)
}

/*
 * The beginning of string methods.
 */
func blStringConcat(obj BlObject,
                    args ...BlObject) BlObject {
    sobj, ok := obj.(*BlStringObject)
    if !ok {
        errpkg.InternError("expected string object")
    }
    var str string
    if blParseArguments("s", args, &str) == -1 {
        return nil
    }
    return NewBlString(sobj.Value + str)
}

func blStringToUpper(obj BlObject,
                     args ...BlObject) BlObject {
    sobj, ok := obj.(*BlStringObject)
    if !ok {
        errpkg.InternError("expected string object")
    }
    var buf bytes.Buffer
    for _, ch := range sobj.Value {
        num := ch
        if ch >= 'a' && ch <= 'z' {
            num = ch - rune(32)
        }
        buf.WriteByte(byte(num))
    }
    return NewBlString(buf.String())
}

func blStringToLower(obj BlObject,
                     args ...BlObject) BlObject {
    sobj, ok := obj.(*BlStringObject)
    if !ok {
        errpkg.InternError("expected string object")
    }
    var buf bytes.Buffer
    for _, ch := range sobj.Value {
        num := ch
        if ch >= 'A' && ch <= 'Z' {
            num = ch + rune(32)
        }
        buf.WriteByte(byte(num))
    }
    return NewBlString(buf.String())
}

func blInitString() {
    BlStringType = BlTypeObject{
        header   : blHeader{&BlTypeType},
        Name     : "string",
        Repr     : blStringRepr,
        GetMember: blStringGetMember,
        EvalCond : blStringEvalCond,
        Compare  : blStringCompare,
        Init     : blStringInit,
        Sequence : &blStringSequence,
        methods  : blStringMethods,
    }
    blTypeFinish(&BlStringType)
}