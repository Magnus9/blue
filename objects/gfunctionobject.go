
/*
 * Object represents a builtin function. These
 * can be attached to a module, or wrapped inside
 * of a BlGMethodObject. Instead of evaluating the
 * BlGMethodObject as a builtin-function/builtin method
 * they are split up instead to make them easier to
 * call.
 */
package objects

import (
    "fmt"
)
const (
    GFUNC_NOARGS  = 1
    GFUNC_VARARGS = 2
)
type gfunction func(BlObject, ...BlObject) BlObject
type BlGFunctionObject struct {
    header   blHeader
    Name     string
    Function gfunction
    Flags    int
    Params   int
}
func (bgfo *BlGFunctionObject) BlType() *BlTypeObject {
    return bgfo.header.typeobj
}
var BlGFunctionType BlTypeObject

func NewBlGFunction(name string, function gfunction,
                    flags int) BlGFunctionObject {
    return BlGFunctionObject{
        header  : blHeader{&BlGFunctionType},
        Name    : name,
        Function: function,
        Flags   : flags,
    }
}

func blGFunctionRepr(obj BlObject) *BlStringObject {
    fobj := obj.(*BlGFunctionObject)
    str := fmt.Sprintf("<builtin-function '%s'>",
                       fobj.Name)
    return NewBlString(str)
}

func blGFunctionEvalCond(obj BlObject) bool {
    fobj := obj.(*BlGFunctionObject)
    if (fobj.Flags & GFUNC_VARARGS) != 0 {
        return true
    }
    return false
}

func blInitGFunction() {
    BlGFunctionType = BlTypeObject{
        header  : blHeader{&BlTypeType},
        Name    : "builtin-function",
        Repr    : blGFunctionRepr,
        EvalCond: blGFunctionEvalCond,
    }
}