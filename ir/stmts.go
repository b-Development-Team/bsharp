package ir

import (
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type nodeBuilder struct {
	ArgTypes []types.Type
	Build    func(b *Builder, pos *tokens.Pos, args []Node) (Call, error)
}

var nodeBuilders = make(map[string]nodeBuilder)

type blockBuilder struct {
	Build func(b *Builder, pos *tokens.Pos, args []parser.Node) (Block, error)
}

var blockBuilders = make(map[string]blockBuilder)

type CallNode struct {
	pos  *tokens.Pos
	Call Call
}

func (c *CallNode) Pos() *tokens.Pos { return c.pos }
func (c *CallNode) Type() types.Type { return c.Call.Type() }

type BlockNode struct {
	pos   *tokens.Pos
	Block Block
}

func (b *BlockNode) Pos() *tokens.Pos { return b.pos }
func (b *BlockNode) Type() types.Type { return types.NULL }

func NewCallNode(call Call, pos *tokens.Pos) *CallNode {
	return &CallNode{
		pos:  pos,
		Call: call,
	}
}

func NewBlockNode(blk Block, pos *tokens.Pos) *BlockNode {
	return &BlockNode{
		pos:   pos,
		Block: blk,
	}
}

func (b *Builder) buildNode(node parser.Node) (Node, error) {
	switch n := node.(type) {
	case *parser.CallNode:
		builder, exists := nodeBuilders[n.Name]
		if !exists {
			// Block?
			blkBuilder, exists := blockBuilders[n.Name]
			if exists {
				blk, err := blkBuilder.Build(b, n.Pos(), n.Args)
				if err != nil {
					return nil, err
				}
				return &BlockNode{n.Pos(), blk}, nil
			}

			// Special case?
			switch n.Name {
			case "FUNC":
				return nil, b.buildFnDef(n)

			case "TYPEDEF":
				return nil, b.checkTypeDef(n)

			case "CONSTDEF":
				return nil, b.checkConst(n)

			case "CONST":
				arg, err := b.buildNode(n.Args[0])
				if err != nil {
					return nil, err
				}
				if !types.IDENT.Equal(arg.Type()) {
					return nil, arg.Pos().Error("expected identifier as argument to CONST")
				}
				name := arg.(*Const).Value.(string)
				v, exists := b.consts[name]
				if !exists {
					return nil, arg.Pos().Error("undefined constant: %s", name)
				}
				return &Const{
					typ: v.typ,
					pos: n.Pos(),

					Value: v.Value,
				}, nil

			case "IMPORT":
				if b.Scope.CurrType() != ScopeTypeGlobal {
					return nil, n.Pos().Error("import must be at the top level")
				}
				return nil, nil
			}

			// Is function?
			_, exists = b.Funcs[n.Name]
			if exists {
				return b.buildFnCall(n)
			}

			// Is extension?
			_, exists = b.extensions[n.Name]
			if exists {
				return b.buildExtensionCall(n)
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
		hasErr := b.MatchTypes(n.Pos(), args, builder.ArgTypes)
		if hasErr {
			b.FixTypes(&args, builder.ArgTypes, n.Pos())
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

	case *parser.ByteNode:
		return b.buildByte(n), nil

	case *parser.NumberNode:
		return b.buildNumber(n)

	case *parser.BoolNode:
		return b.buildBool(n), nil

	case *parser.NullNode:
		return NewConst(types.NULL, n.Pos(), nil), nil

	default:
		return nil, n.Pos().Error("unknown node type: %T", n)
	}
}
