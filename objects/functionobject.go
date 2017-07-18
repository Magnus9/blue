
package objects

import (
    "fmt"
    "github.com/Magnus9/blue/interm"
)
type BlFunctionObject struct {
    header    blHeader
    name      string
    Params    []string
    ParamLen  int
    StarParam bool
    Block     *interm.Node
}
func (bfo *BlFunctionObject) BlType() *BlTypeObject {
    return bfo.header.typeobj
}
var BlFunctionType BlTypeObject

func NewBlFunction(name string, params []string,
                   paramLen int, block *interm.Node,
                   starParam bool) BlObject {
    bfo := &BlFunctionObject{
        header   : blHeader{&BlFunctionType},
        name     : name,
        Params   : params,
        ParamLen : paramLen,
        StarParam: starParam,
        Block    : block,
    }
    /*
     * If starParam == true, we reduce bfo.ParamLen with one,
     * since it gets much easier to evaluate a func call
     */
    if starParam {
        bfo.ParamLen--
    }
    return bfo
}

func blFunctionRepr(obj BlObject) BlObject {
    fobj := obj.(*BlFunctionObject)
    str := fmt.Sprintf("<function '%s', params=%d>\n",
                       fobj.name, fobj.ParamLen)
    return NewBlString(str)
}

func blFunctionEvalCond(obj BlObject) bool {
    sobj := obj.(*BlFunctionObject)
    if sobj.ParamLen > 0 {
        return true
    }
    return false
}

func blInitFunction() {
    BlFunctionType = BlTypeObject{
        header  : blHeader{&BlTypeType},
        Name    : "function",
        Repr    : blFunctionRepr,
        EvalCond: blFunctionEvalCond,
    }
}