package bstar

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
)

type BStar struct {
	ir *ir.IR
}

type Node interface {
	Code(opts *BStarConfig) string
}

type BlockNode struct {
	Body      []Node
	DoesPrint bool
}

type BStarConfig struct {
	Seperator string
}

func (n *BlockNode) Code(opts *BStarConfig) string {
	out := &strings.Builder{}
	out.WriteByte('[')
	for i, n := range n.Body {
		if i > 0 {
			out.WriteString(opts.Seperator)
		}
		out.WriteString(n.Code(opts))
	}
	out.WriteByte(']')
	return out.String()
}

type ConstNode struct{ any }

func (c ConstNode) Code(opts *BStarConfig) string {
	return fmt.Sprintf("%v", c.any)
}

func constNode(v any) Node {
	return ConstNode{v}
}

func blockNode(doesPrint bool, body ...Node) Node {
	return &BlockNode{Body: body, DoesPrint: doesPrint}
}

func NewBStar(i *ir.IR) *BStar {
	return &BStar{i}
}

func (b *BStar) Build() (Node, error) {
	out, err := b.buildNodes(b.ir.Body)
	if err != nil {
		return nil, err
	}
	return blockNode(false, append([]Node{constNode("BLOCK")}, append(out, b.noPrintNode())...)...), nil
}

func (b *BStar) noPrint(node Node) Node {
	return blockNode(false, constNode("DEFINE"), constNode("noprint"), node)
}

func (b *BStar) noPrintNode() Node {
	return blockNode(false, constNode("DEFINE"), constNode("noprint"), constNode("NOPRINT"))
}
