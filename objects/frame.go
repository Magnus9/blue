
package objects

import "github.com/Magnus9/blue/interm"

type BlFrame struct {
    Prev     *BlFrame
    Globals  map[string]BlObject
    Locals   map[string]BlObject
    Pathname string
    Node     *interm.Node
}

func NewBlFrame(prev *BlFrame,
                globals, locals map[string]BlObject,
                pathname string) *BlFrame {
    return &BlFrame{
        Prev    : prev,
        Globals : globals,
        Locals  : locals,
        Pathname: pathname,
    }
}
func (bf *BlFrame) SetNode(node *interm.Node) {
    bf.Node = node
}