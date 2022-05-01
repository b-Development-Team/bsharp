package ir

import (
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

// Typedef pass
func (b *Builder) typePass(p *parser.Parser) error {
	for _, node := range p.Nodes {
		call, ok := node.(*parser.CallNode)
		if !ok {
			continue
		}
		if call.Name != "TYPEDEF" {
			continue
		}

		// Check if valid signature
		if len(call.Args) != 2 {
			return call.Pos().Error("expected two arguments to TYPEDEF")
		}

		// Check name
		name, ok := call.Args[0].(*parser.IdentNode)
		if !ok {
			return call.Args[0].Pos().Error("expected identifier as first argument to TYPEDEF")
		}

		// Get type
		typV, ok := call.Args[1].(*parser.IdentNode)
		if !ok {
			return call.Args[1].Pos().Error("expected identifier as second argument to TYPEDEF")
		}
		typ, err := types.ParseType(typV.Value, b.typeNames)
		if err != nil {
			return call.Args[1].Pos().Error("%s", err.Error())
		}

		// Check if type already exists
		if _, ok := b.typeNames[name.Value]; ok {
			return name.Pos().Error("type already exists")
		}

		// Add type
		b.typeNames[name.Value] = typ
	}

	return nil
}

func (b *Builder) checkTypeDef(n *parser.CallNode) error {
	if b.Scope.CurrType() != ScopeTypeGlobal {
		return n.Pos().Error("TYPEDEF can only be used in global scope")
	}
	return nil
}

type CastNode struct {
	Value Node

	typ types.Type
}

func (c *CastNode) Type() types.Type { return c.typ }
func (c *CastNode) Pos() *tokens.Pos { return c.Value.Pos() }
func (c *CastNode) Args() []Node     { return []Node{c.Value} }

func NewCastNode(val Node, typ types.Type) *CastNode {
	return &CastNode{
		Value: val,
		typ:   typ,
	}
}

type CanCastNode struct {
	Value Node
	Typ   types.Type
	pos   *tokens.Pos
}

func (c *CanCastNode) Type() types.Type { return types.BOOL }
func (c *CanCastNode) Args() []Node {
	return []Node{c.Value, NewConst(types.STRING, c.pos, c.Typ.String())}
}

func init() {
	nodeBuilders["FLOAT"] = nodeBuilder{
		ArgTypes: []types.Type{types.NewMulType(types.INT, types.STRING)},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return NewCastNode(args[0], types.FLOAT), nil
		},
	}

	nodeBuilders["STRING"] = nodeBuilder{
		ArgTypes: []types.Type{types.NewMulType(types.INT, types.FLOAT, types.BOOL, types.BYTE, types.NewArrayType(types.BYTE))},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return NewCastNode(args[0], types.STRING), nil
		},
	}

	nodeBuilders["INT"] = nodeBuilder{
		ArgTypes: []types.Type{types.NewMulType(types.FLOAT, types.STRING, types.BYTE)},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return NewCastNode(args[0], types.INT), nil
		},
	}

	nodeBuilders["BYTE"] = nodeBuilder{
		ArgTypes: []types.Type{types.NewMulType(types.INT)},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return NewCastNode(args[0], types.BYTE), nil
		},
	}

	nodeBuilders["ANY"] = nodeBuilder{
		ArgTypes: []types.Type{types.ALL},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return &CastNode{
				Value: args[0],
				typ:   types.ANY,
			}, nil
		},
	}

	nodeBuilders["CANCAST"] = nodeBuilder{
		ArgTypes: []types.Type{types.ALL, types.IDENT},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			typ, err := types.ParseType(args[1].(*Const).Value.(string), b.typeNames)
			if err != nil {
				return nil, args[1].Pos().Error("%s", err.Error())
			}
			return &CanCastNode{
				Value: args[0],
				Typ:   typ,
				pos:   args[1].Pos(),
			}, nil
		},
	}

	nodeBuilders["CAST"] = nodeBuilder{
		ArgTypes: []types.Type{types.ALL, types.IDENT},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			typ, err := types.ParseType(args[1].(*Const).Value.(string), b.typeNames)
			if err != nil {
				return nil, err
			}
			return &CastNode{
				Value: args[0],
				typ:   typ,
			}, nil
		},
	}
}
