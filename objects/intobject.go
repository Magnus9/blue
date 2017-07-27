
package objects

import (
    "fmt"
    "math"
    "github.com/Magnus9/blue/errpkg"
)
type BlIntObject struct {
    header blHeader
    Value  int64
}
func (bio *BlIntObject) BlType() *BlTypeObject {
    return bio.header.typeobj
}

var blIntNumbers = BlNumberMethods{
    NumNeg   : blIntNeg,
    NumCompl : blIntCompl,
    NumOr    : blIntOr,
    NumAnd   : blIntAnd,
    NumXor   : blIntXor,
    NumLshift: blIntLshift,
    NumRshift: blIntRshift,
    NumAdd   : blIntAdd,
    NumSub   : blIntSub,
    NumMul   : blIntMul,
    NumDiv   : blIntDiv,
    NumMod   : blIntMod,
}
var blIntSequence = BlSequenceMethods{
    SeqItem      : blIntItem,
    SeqAssItem   : blIntAssItem,
}
var BlIntType BlTypeObject

func NewBlInt(value int64) *BlIntObject {
    return &BlIntObject{
        header  : blHeader{&BlIntType},
        Value   : value,
    }
}

func blIntNeg(obj BlObject) BlObject {
    iobj := obj.(*BlIntObject)
    return NewBlInt(-iobj.Value)
}

func blIntCompl(obj BlObject) BlObject {
    iobj := obj.(*BlIntObject)
    return NewBlInt(^iobj.Value)
}

func blIntOr(a, b BlObject) BlObject {
    aIobj := a.(*BlIntObject)
    bIobj := b.(*BlIntObject)
    return NewBlInt(aIobj.Value | bIobj.Value)
}

func blIntAnd(a, b BlObject) BlObject {
    aIobj := a.(*BlIntObject)
    bIobj := b.(*BlIntObject)
    return NewBlInt(aIobj.Value & bIobj.Value)
}

func blIntXor(a, b BlObject) BlObject {
    aIobj := a.(*BlIntObject)
    bIobj := b.(*BlIntObject)
    return NewBlInt(aIobj.Value ^ bIobj.Value)
}

func blIntLshift(a, b BlObject) BlObject {
    aIobj := a.(*BlIntObject)
    bIobj := b.(*BlIntObject)
    return NewBlInt(aIobj.Value << uint(bIobj.Value))
}

func blIntRshift(a, b BlObject) BlObject {
    aIobj := a.(*BlIntObject)
    bIobj := b.(*BlIntObject)
    return NewBlInt(aIobj.Value >> uint(bIobj.Value))
}

func blIntAdd(a, b BlObject) BlObject {
    aIobj := a.(*BlIntObject)
    bIobj := b.(*BlIntObject)
    return NewBlInt(aIobj.Value + bIobj.Value)
}

func blIntSub(a, b BlObject) BlObject {
    aIobj := a.(*BlIntObject)
    bIobj := b.(*BlIntObject)
    return NewBlInt(aIobj.Value - bIobj.Value)
}

func blIntMul(a, b BlObject) BlObject {
    aIobj := a.(*BlIntObject)
    bIobj := b.(*BlIntObject)
    return NewBlInt(aIobj.Value * bIobj.Value)
}

func blIntDiv(a, b BlObject) BlObject {
    aIobj := a.(*BlIntObject)
    bIobj := b.(*BlIntObject)
    if bIobj.Value == 0 {
        errpkg.SetErrmsg("int division by zero")
        return nil
    }
    return NewBlInt(aIobj.Value / bIobj.Value)
}

func blIntMod(a, b BlObject) BlObject {
    aIobj := a.(*BlIntObject)
    bIobj := b.(*BlIntObject)
    if bIobj.Value == 0 {
        errpkg.SetErrmsg("int modulo by zero")
        return nil
    }
    return NewBlInt(aIobj.Value % bIobj.Value)
}

func blIntItem(obj BlObject, num int) BlObject {
    iobj := obj.(*BlIntObject)
    if num > 63  || num < 0 {
        errpkg.SetErrmsg("subscript position out of bounds")
        return nil
    }
    return NewBlInt((iobj.Value >> uint(num)) & 0x01)
}

func blIntAssItem(obj BlObject, value BlObject,
                     num int) int {
    iobj := obj.(*BlIntObject)
    if num > 63 || num < 0 {
        errpkg.SetErrmsg("subscript position out of bounds")
        return -1
    }
    t, ok := value.(*BlIntObject)
    if !ok {
        errpkg.SetErrmsg("value must be an integer")
        return -1
    }
    if t.Value != 1 && t.Value != 0 {
        errpkg.SetErrmsg("value must be either 1 or 0")
        return -1
    }
    bitValue := int64(math.Pow(2, float64(num)))
    if t.Value == 1 {
        iobj.Value |= bitValue
    } else {
        iobj.Value &= ^bitValue
    }
    return 0
}

func blIntRepr(obj BlObject) *BlStringObject {
    iobj := obj.(*BlIntObject)
    str := fmt.Sprintf("%d", iobj.Value)

    return NewBlString(str)
}

func blIntEvalCond(obj BlObject) bool {
    iobj := obj.(*BlIntObject)
    if iobj.Value > 0 {
        return true
    }
    return false
}

func blIntCompare(a, b BlObject) int {
    aIobj := a.(*BlIntObject)
    bIobj := b.(*BlIntObject)
    switch {
    case aIobj.Value < bIobj.Value:
        return -1
    case aIobj.Value > bIobj.Value:
        return 1
    default:
        return 0
    }
}

func blIntInit(obj *BlTypeObject,
               args ...BlObject) BlObject {
    var arg BlObject
    if blParseArguments("o", args, &arg) == -1 {
        return nil
    }
    switch t := arg.(type) {
        case *BlIntObject:
            return t
        case *BlFloatObject:
            return NewBlInt(int64(t.value))
    }
    errpkg.SetErrmsg("expected number")
    return nil
}

func blInitInt() {
    BlIntType = BlTypeObject{
        header  : blHeader{&BlTypeType},
        Name    : "int",
        Repr    : blIntRepr,
        EvalCond: blIntEvalCond,
        Compare : blIntCompare,
        Init    : blIntInit,
        Numbers : &blIntNumbers,
        Sequence: &blIntSequence,
    }
    blTypeFinish(&BlIntType)
}