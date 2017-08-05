
package objects

import (
    "fmt"
    "github.com/Magnus9/blue/interm"
)
type BlFunctionObject struct {
    header    blHeader
    Path      string
    Name      string
    Globals   map[string]BlObject
    Params    []string
    ParamLen  int
    StarParam bool
    Block     *interm.Node
}
func (bfo *BlFunctionObject) BlType() *BlTypeObject {
    return bfo.header.typeobj
}
var BlFunctionType BlTypeObject

func NewBlFunction(path, name string, globals map[string]BlObject,
                   params []string, paramLen int,
                   block *interm.Node, starParam bool) BlObject {
    bfo := &BlFunctionObject{
        header   : blHeader{&BlFunctionType},
        Path     : path,
        Name     : name,
        Globals  : globals,
        Params   : params,
        ParamLen : paramLen,
        StarParam: starParam,
        Block    : block,
    }
    /*
     * If starParam == true, we reduce bfo.ParamLen with one,
     * since it gets much easier to evaluate a func call.
     */
    if starParam {
        bfo.ParamLen--
    }
    return bfo
}

func blFunctionRepr(obj BlObject) *BlStringObject {
    fobj := obj.(*BlFunctionObject)
    mesg := fmt.Sprintf("<function '%s', params=%d>",
                        fobj.Name, fobj.ParamLen)
    return NewBlString(mesg)
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