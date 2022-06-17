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
	Body    []Node
	NoPrint bool
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

func blockNode(body ...Node) Node {
	return &BlockNode{Body: body}
}

func NewBStar(i *ir.IR) *BStar {
	return &BStar{i}
}

func (b *BStar) Build() ([]Node, error) {
	out := make([]Node, 0)
	for _, n := range b.ir.Body {
		node, err := b.buildNode(n)
		if err != nil {
			return nil, err
		}
		out = append(out, node)
	}
	return out, nil
}
