
package objects

import "github.com/Magnus9/blue/interm"

type BlFrame struct {
    Prev     *BlFrame
    Locals   map[string]BlObject
    Pathname string
    Node     *interm.Node
}

func NewBlFrame(prev *BlFrame, locals map[string]BlObject,
                pathname string) *BlFrame {
    return &BlFrame{
        Prev    : prev,
        Locals  : locals,
        Pathname: pathname,
    }
}
func (bf *BlFrame) SetNode(node *interm.Node) {
    bf.Node = node
}