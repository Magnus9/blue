
package objects

import "github.com/Magnus9/blue/interm"

type BlFrame struct {
    Prev     *BlFrame
    Globals  map[string]BlObject
    Locals   map[string]BlObject
    Pathname string
    Name     string
    Node     *interm.Node
}

func NewBlFrame(prev *BlFrame,
                globals, locals map[string]BlObject,
                pathname, name string) *BlFrame {
    return &BlFrame{
        Prev    : prev,
        Globals : globals,
        Locals  : locals,
        Pathname: pathname,
        Name    : name,
    }
}
func (bf *BlFrame) SetNode(node *interm.Node) {
    bf.Node = node
}