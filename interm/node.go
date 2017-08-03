
package interm

import (
    "bytes"
)

const (
    FLAG_STARPARAM = 1 << 0
    FLAG_RANGELHS  = 1 << 0
    FLAG_RANGERHS  = 1 << 1
)
type Node struct {
    Str       string
    Line      string
    NodeType  int
    LineNum   int
    Nchildren int
    Flags     int
    Children  []*Node
}

func New(str, line string, nodeType, lineNum int) *Node {
    return &Node{
        Str     : str,
        Line    : line,
        NodeType: nodeType,
        LineNum : lineNum,
        Children: make([]*Node, 0, 2),  
    }
}

func (n *Node) Add(node *Node) {
    n.Children = append(n.Children, node)
    n.Nchildren++
}

func (n *Node) GiveRootTo(node *Node) *Node {
    node.Add(n)

    return node
}

func (n *Node) String() string {
    return n.Str
}

func (n *Node) ListTree() string {
    if n.Nchildren == 0 {
        return n.String()
    }
    var buf bytes.Buffer
    buf.WriteByte('(')
    buf.WriteString(n.String())
    buf.WriteByte(' ')
    for i := 0; i < n.Nchildren; i++ {
        if i > 0 {
            buf.WriteByte(' ')
        }
        buf.WriteString(n.Children[i].ListTree())
    }
    buf.WriteByte(')')

    return buf.String()
}