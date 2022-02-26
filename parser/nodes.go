package parser

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/tokens"
)

type Node interface {
	fmt.Stringer
	Pos() *tokens.Pos
}

type CallNode struct {
	pos  *tokens.Pos
	Name string
	Args []Node
}

func (n *CallNode) String() string {
	args := &strings.Builder{}
	for i, arg := range n.Args {
		args.WriteString(arg.String())
		if i < len(n.Args)-1 {
			args.WriteString(" ")
		}
	}
	return fmt.Sprintf("(%s)[%s %s]", n.Pos().String(), n.Name, args)
}

func (n *CallNode) Pos() *tokens.Pos {
	return n.pos
}

type IdentNode struct {
	pos   *tokens.Pos
	Value string
}

func (i *IdentNode) Pos() *tokens.Pos {
	return i.pos
}

func (i *IdentNode) String() string {
	return fmt.Sprintf("Ident(%s, %s)", i.Value, i.Pos().String())
}

type StringNode struct {
	pos   *tokens.Pos
	Value string
}

func (s *StringNode) Pos() *tokens.Pos {
	return s.pos
}

func (s *StringNode) String() string {
	return fmt.Sprintf("String(\"%s\", %s)", s.Value, s.Pos().String())
}

type NumberNode struct {
	pos   *tokens.Pos
	Value string
}

func (n *NumberNode) Pos() *tokens.Pos {
	return n.pos
}

func (n *NumberNode) String() string {
	return fmt.Sprintf("Number(%s, %s)", n.Value, n.Pos().String())
}
