package ir

import (
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type Call interface {
	Type() types.Type
}

type nodeBuilder struct {
	ArgTypes []types.Type
	Build    func(b *Builder, pos *tokens.Pos, args []Node) (Call, error)
}

var nodeBuilders = make(map[string]nodeBuilder)

type CallNode struct {
	pos  *tokens.Pos
	Call Call
}

func (c *CallNode) Pos() *tokens.Pos { return c.pos }
func (c *CallNode) Type() types.Type { return c.Call.Type() }

func (b *Builder) buildNode(node parser.Node) (Node, error) {
	switch n := node.(type) {
	case *parser.CallNode:
		builder, exists := nodeBuilders[n.Name]
		if !exists {
			// Function def
			if n.Name == "FUNC" {
				return nil, b.buildFnDef(n)
			}

			// Import
			if n.Name == "IMPORT" {
				if b.Scope.CurrType() != ScopeTypeGlobal {
					return nil, n.Pos().Error("import must be at the top level")
				}
				return nil, nil
			}

			// Is function?
			_, exists := b.Funcs[n.Name]
			if exists {
				return b.buildFnCall(n)
			}
			return nil, n.Pos().Error("unknown function: " + n.Name)
		}
		args := make([]Node, len(n.Args))
		for i, arg := range n.Args {
			node, err := b.buildNode(arg)
			if err != nil {
				return nil, err
			}
			args[i] = node
		}
		call, err := builder.Build(b, n.Pos(), args)
		if err != nil {
			return nil, err
		}
		return &CallNode{
			pos:  n.Pos(),
			Call: call,
		}, nil

	case *parser.IdentNode:
		return b.buildIdent(n), nil

	case *parser.StringNode:
		return b.buildString(n), nil

	case *parser.NumberNode:
		return b.buildNumber(n), nil

	default:
		return nil, n.Pos().Error("unknown node type: %T", n)
	}
}
