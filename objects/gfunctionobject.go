
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
type gfunction func(obj BlObject,
                    args ...BlObject) BlObject
type BlGFunctionObject struct {
    header   blHeader
    name     string
    Function gfunction
    Params   int
}
func (bgfo *BlGFunctionObject) BlType() *BlTypeObject {
    return bgfo.header.typeobj
}
var BlGFunctionType BlTypeObject

func NewBlGFunction(name string, function gfunction,
                    params int) BlGFunctionObject {
    return BlGFunctionObject{
        header  : blHeader{&BlGFunctionType},
        name    : name,
        Function: function,
        Params  : params,
    }
}

func blGFunctionRepr(obj BlObject) BlObject {
    fobj := obj.(*BlGFunctionObject)
    str := fmt.Sprintf("<builtin-function '%s', params=%d>\n",
                       fobj.name, fobj.Params)
    return NewBlString(str)
}

func blGFunctionEvalCond(obj BlObject) bool {
    fobj := obj.(*BlGFunctionObject)
    if fobj.Params > 0 {
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