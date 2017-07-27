
package objects

import (
    "fmt"
    "bytes"
    "strings"
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
    SeqItem  : blStringItem,
    SeqConcat: blStringConcat,
    SeqRepeat: blStringRepeat,
    SeqSize  : blStringSize,
}
var blStringMethods = []BlGFunctionObject {
    NewBlGFunction("index",      stringIndex,      GFUNC_VARARGS),
    NewBlGFunction("split",      stringSplit,      GFUNC_VARARGS),
    NewBlGFunction("concat",     stringConcat,     GFUNC_VARARGS),
    NewBlGFunction("toupper",    stringToUpper,    GFUNC_NOARGS ),
    NewBlGFunction("tolower",    stringToLower,    GFUNC_NOARGS ),
    NewBlGFunction("startswith", stringStartsWith, GFUNC_VARARGS),
    NewBlGFunction("endswith",   stringEndsWith,   GFUNC_VARARGS),
}
var BlStringType BlTypeObject

func NewBlString(value string) *BlStringObject {
    return &BlStringObject{
        header  : blHeader{&BlStringType},
        Value   : value,
        vsize   : len(value),
    }
}

func blStringItem(obj BlObject, num int) BlObject {
    sobj := obj.(*BlStringObject)
    if num >= sobj.vsize || num < 0 {
        errpkg.SetErrmsg("subscript position out of bounds")
        return nil
    }
    return NewBlString(string(sobj.Value[num]))
}

func blStringConcat(a, b BlObject) BlObject {
    sobj := a.(*BlStringObject)
    t, ok := b.(*BlStringObject)
    if !ok {
        errpkg.SetErrmsg("cannot add '%s' to string",
                         b.BlType().Name)
        return nil
    }
    return NewBlString(sobj.Value + t.Value)
}

func blStringRepeat(a, b BlObject) BlObject {
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

func blStringSize(obj BlObject) int {
    return obj.(*BlStringObject).vsize
}

func blStringRepr(obj BlObject) *BlStringObject {
    sobj := obj.(*BlStringObject)
    return NewBlString(fmt.Sprintf("\"%s\"",
                       sobj.Value))
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
    var arg BlObject
    if blParseArguments("o", args, &arg) == -1 {
        return nil
    }
    typeobj := arg.BlType()
    if fn := typeobj.Repr; fn != nil {
        return fn(arg)
    }
    /*
     * This is more of an internal error. Every object
     * should have a Repr function attached to it.
     */
    errpkg.SetErrmsg("'%s' object has no representation",
                     typeobj.Name)
    return nil
}

/*
 * Splits a string into a list using var 'sep' as
 * the separator. This function does not split
 * if the left/right side is empty.
 */
func stringSplit(obj BlObject, args ...BlObject) BlObject {
    var sep string = " "
    var max int64  = -1
    if blParseArguments("|si", args, &sep, &max) == -1 {
        return nil
    }
    self := obj.(*BlStringObject)
    lobj := NewBlList(0)
    if max == 0 {
        lobj.Append(NewBlString(self.Value))
    } else {
        siz := len(sep)
        var j uint; var mark int
        for i := 0; i < self.vsize; i++ {
            if int(j) == int(max) {
                break
            }
            tot := i + siz
            if tot >= self.vsize {
                break
            }
            if self.Value[i:tot] == sep {
                if i != 0 {
                    lobj.Append(NewBlString(self.Value[mark:i]))
                }
                mark = tot
                j++
            }
        }
        if mark < self.vsize {
            lobj.Append(NewBlString(self.Value[mark:]))
        }
    }
    return lobj
}

func stringIndex(obj BlObject, args ...BlObject) BlObject {
    var str string
    if blParseArguments("s", args, &str) == -1 {
        return nil
    }
    self := obj.(*BlStringObject)
    return NewBlInt(int64(strings.Index(self.Value, str)))
}

func stringConcat(obj BlObject,
                  args ...BlObject) BlObject {
    var str string
    if blParseArguments("s", args, &str) == -1 {
        return nil
    }
    sobj := obj.(*BlStringObject)
    return NewBlString(sobj.Value + str)
}

func stringToUpper(obj BlObject,
                   args ...BlObject) BlObject {
    sobj := obj.(*BlStringObject)
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

func stringToLower(obj BlObject,
                   args ...BlObject) BlObject {
    sobj := obj.(*BlStringObject)
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

func stringStartsWith(obj BlObject,
                      args ...BlObject) BlObject {
    var str string
    var s int64
    if blParseArguments("s|i", args, &str, &s) == -1 {
        return nil
    }
    self := obj.(*BlStringObject)
    if s < int64(0) || s > int64(self.vsize) {
        return BlFalse
    }
    if strings.HasPrefix(self.Value[s:], str) {
        return BlTrue
    }
    return BlFalse
}

func stringEndsWith(obj BlObject,
                    args ...BlObject) BlObject {
    self := obj.(*BlStringObject)
    var str string
    var s int64 = int64(self.vsize - 1)
    if blParseArguments("s|i", args, &str, &s) == -1 {
        return nil
    }
    s++
    if s < int64(0) || s > int64(self.vsize) {
        return BlFalse
    }
    if strings.HasSuffix(self.Value[:s], str) {
        return BlTrue
    }
    return BlFalse
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