package ir

import (
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type IfNode struct {
	NullCall
	Condition Node
	Body      Node
}

type IfElseNode struct {
	NullCall
	Condition Node
	Body      Node
	Else      Node
}

type WhileNode struct {
	NullCall
	Condition Node
	Body      Node
}

func init() {
	nodeBuilders["IF"] = nodeBuilder{
		ArgTypes: []types.Type{types.BOOL, types.ANY, types.VARIADIC},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			if len(args) == 2 { // If node
				return &IfNode{
					Condition: args[0],
					Body:      args[1],
				}, nil
			}

			if len(args) == 4 { // If node
				if !types.IDENT.Equal(args[2].Type()) || args[2].(*Const).Value.(string) != "ELSE" {
					return nil, args[2].Pos().Error("expected ELSE")
				}
				return &IfElseNode{
					Condition: args[0],
					Body:      args[1],
					Else:      args[3],
				}, nil
			}

			return nil, pos.Error("expected IF or IF ELSE")
		},
	}

	nodeBuilders["WHILE"] = nodeBuilder{
		ArgTypes: []types.Type{types.BOOL, types.ANY},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return &WhileNode{
				Condition: args[0],
				Body:      args[1],
			}, nil
		},
	}
}

type BlockNode struct {
	NullCall
	pos  *tokens.Pos
	Body []Node
}

func (b *BlockNode) Pos() *tokens.Pos { return b.pos }

func (b *Builder) addBlockNode(n *parser.CallNode) (Node, error) {
	b.Scope.Push(ScopeTypeBlock)

	body := make([]Node, len(n.Args))
	var err error
	for i, arg := range n.Args {
		body[i], err = b.buildNode(arg)
		if err != nil {
			return nil, err
		}
	}

	b.Scope.Pop()

	return &BlockNode{
		Body: body,
		pos:  n.Pos(),
	}, nil
}
