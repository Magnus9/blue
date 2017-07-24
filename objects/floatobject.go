
/*
 * Pretty primitive float object. For now it
 * can only present integral and fraction bits.
 */
package objects

import (
    "fmt"
    "math"
    "github.com/Magnus9/blue/errpkg"
)
type BlFloatObject struct {
    header blHeader
    value  float64
}
func (bfo *BlFloatObject) BlType() *BlTypeObject {
    return bfo.header.typeobj
}
var blFloatNumbers = BlNumberMethods{
    NumNeg   : blFloatNeg,
    NumAdd   : blFloatAdd,
    NumSub   : blFloatSub,
    NumMul   : blFloatMul,
    NumDiv   : blFloatDiv,
    NumMod   : blFloatMod,
    NumCoerce: blFloatCoerce,
}
var BlFloatType BlTypeObject

func NewBlFloat(value float64) *BlFloatObject {
    return &BlFloatObject{
        header: blHeader{&BlFloatType},
        value : value,
    }
}

func blFloatNeg(obj BlObject) BlObject {
    fobj := obj.(*BlFloatObject)
    return NewBlFloat(-fobj.value)
}

func blFloatAdd(a, b BlObject) BlObject {
    aFobj := a.(*BlFloatObject)
    bFobj := b.(*BlFloatObject)
    return NewBlFloat(aFobj.value + bFobj.value)
}

func blFloatSub(a, b BlObject) BlObject {
    aFobj := a.(*BlFloatObject)
    bFobj := b.(*BlFloatObject)
    return NewBlFloat(aFobj.value - bFobj.value)
}

func blFloatMul(a, b BlObject) BlObject {
    aFobj := a.(*BlFloatObject)
    bFobj := b.(*BlFloatObject)
    return NewBlFloat(aFobj.value * bFobj.value)
}

func blFloatDiv(a, b BlObject) BlObject {
    aFobj := a.(*BlFloatObject)
    bFobj := b.(*BlFloatObject)
    if bFobj.value == 0.0 {
        errpkg.SetErrmsg("float division by zero")
        return nil
    }
    return NewBlFloat(aFobj.value / bFobj.value)
}

func blFloatMod(a, b BlObject) BlObject {
    aFobj := a.(*BlFloatObject)
    bFobj := b.(*BlFloatObject)
    if bFobj.value == 0.0 {
        errpkg.SetErrmsg("float modulo by zero")
        return nil
    }
    return NewBlFloat(math.Mod(aFobj.value, bFobj.value))
}

func blFloatCoerce(a, b *BlObject) int {
    t, ok := (*b).(*BlIntObject)
    if !ok {
        return -1
    }
    obj := NewBlFloat(float64(t.Value))
    *b = obj

    return 0
}

func blFloatRepr(obj BlObject) *BlStringObject {
    fobj := obj.(*BlFloatObject)
    return NewBlString(fmt.Sprintf("%f", fobj.value))
}

func blFloatEvalCond(obj BlObject) bool {
    fobj := obj.(*BlFloatObject)
    if fobj.value > 0.0 {
        return true
    }
    return false
}

func blFloatCompare(a, b BlObject) int {
    aFobj := a.(*BlFloatObject)
    bFobj := b.(*BlFloatObject)
    switch {
    case aFobj.value < bFobj.value:
        return -1
    case aFobj.value > bFobj.value:
        return 1
    default:
        return 0
    }
}

func blInitFloat() {
    BlFloatType = BlTypeObject{
        header  : blHeader{&BlTypeType},
        Name    : "float",
        Repr    : blFloatRepr,
        EvalCond: blFloatEvalCond,
        Compare : blFloatCompare,
        Numbers : &blFloatNumbers,
    }
}