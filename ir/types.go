package ir

import (
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

// Typedef & Constdef pass
func (b *Builder) defPass(p *parser.Parser) error {
	for _, node := range p.Nodes {
		call, ok := node.(*parser.CallNode)
		if !ok {
			continue
		}

		switch call.Name {
		case "TYPEDEF":
			// Check if valid signature
			if len(call.Args) != 2 {
				b.Error(ErrorLevelError, call.Pos(), "expected two arguments to TYPEDEF")
				continue
			}

			// Check name
			name, ok := call.Args[0].(*parser.IdentNode)
			if !ok {
				b.Error(ErrorLevelError, call.Args[0].Pos(), "expected identifier as first argument to TYPEDEF")
				continue
			}

			// Get type
			typV, ok := call.Args[1].(*parser.IdentNode)
			var typ types.Type
			if !ok {
				// check if its null
				_, ok := call.Args[1].(*parser.NullNode)
				if ok {
					typ = types.NULL
				} else {
					b.Error(ErrorLevelError, call.Args[1].Pos(), "expected identifier as second argument to TYPEDEF")
					continue
				}
			}
			if typ == nil {
				var err error
				typ, err = types.ParseType(typV.Value, b.typeNames)
				if err != nil {
					b.Error(ErrorLevelError, call.Args[1].Pos(), "%s", err.Error())
					continue
				}
			}

			// Check if type already exists
			if _, ok := b.typeNames[name.Value]; ok {
				b.Error(ErrorLevelError, name.Pos(), "cannot redefine type %s", name.Value)
				continue
			}

			// Add type
			b.typeNames[name.Value] = typ

		case "CONSTDEF":
			// Check if valid signature
			if len(call.Args) != 2 {
				b.Error(ErrorLevelError, call.Pos(), "expected two arguments to CONSTDEF")
				continue
			}

			// Check name
			name, ok := call.Args[0].(*parser.IdentNode)
			if !ok {
				b.Error(ErrorLevelError, call.Args[0].Pos(), "expected identifier as first argument to CONSTDEF")
				continue
			}

			// Build val
			_, ok = call.Args[1].(*parser.CallNode)
			if ok {
				b.Error(ErrorLevelError, call.Args[1].Pos(), "expected constant value as second argument to CONSTDEF")
				continue
			}
			val, err := b.buildNode(call.Args[1])
			if err != nil {
				return err
			}
			v, ok := val.(*Const)
			if !ok {
				b.Error(ErrorLevelError, call.Args[1].Pos(), "expected constant value as second argument to CONSTDEF")
				continue
			}

			// Add consts
			b.consts[name.Value] = v
		}

	}

	return nil
}

func (b *Builder) checkTypeDef(n *parser.CallNode) error {
	if b.Scope.CurrType() != ScopeTypeGlobal {
		b.Error(ErrorLevelError, n.Pos(), "TYPEDEF can only be used in global scope")
	}
	return nil
}

func (b *Builder) checkConst(n *parser.CallNode) error {
	if b.Scope.CurrType() != ScopeTypeGlobal {
		b.Error(ErrorLevelError, n.Pos(), "CONSTDEF can only be used in global scope")
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
				b.Error(ErrorLevelError, args[1].Pos(), "%s", err.Error())
				typ = types.INVALID
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
