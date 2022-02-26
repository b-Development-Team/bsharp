package ir

import (
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type ArrayNode struct {
	Values []Node
	typ    types.Type
}

func (a *ArrayNode) Type() types.Type { return a.typ }

type IndexNode struct {
	Value Node
	Index Node
	typ   types.Type
}

func (i *IndexNode) Type() types.Type { return i.typ }

type LengthNode struct {
	Value Node
}

func (l *LengthNode) Type() types.Type { return types.INT }

type MakeNode struct {
	typ types.Type
}

func (m *MakeNode) Type() types.Type { return m.typ }

type SetNode struct {
	NullCall

	Map   Node
	Key   Node
	Value Node
}

type GetNode struct {
	Map Node
	Key Node

	typ types.Type
}

func (g *GetNode) Type() types.Type { return g.typ }

func init() {
	nodeBuilders["ARRAY"] = nodeBuilder{
		ArgTypes: []types.Type{types.ANY, types.VARIADIC},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			// Check types
			typ := args[0].Type()
			if len(args) > 1 {
				for _, arg := range args[1:] {
					if !typ.Equal(arg.Type()) {
						return nil, arg.Pos().Error("expected type %s, got %s", typ.String(), arg.Type().String())
					}
				}
			}

			return &ArrayNode{
				Values: args,
				typ:    types.NewArrayType(typ),
			}, nil
		},
	}

	nodeBuilders["INDEX"] = nodeBuilder{
		ArgTypes: []types.Type{types.NewMulType(types.STRING, types.ARRAY), types.INT},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			var outTyp types.Type
			if types.ARRAY.Equal(args[0].Type()) {
				outTyp = args[0].Type().(*types.ArrayType).ElemType
			} else {
				outTyp = types.STRING
			}

			return &IndexNode{
				Value: args[0],
				Index: args[1],
				typ:   outTyp,
			}, nil
		},
	}

	nodeBuilders["LENGTH"] = nodeBuilder{
		ArgTypes: []types.Type{types.NewMulType(types.STRING, types.ARRAY)},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return &LengthNode{
				Value: args[0],
			}, nil
		},
	}

	nodeBuilders["MAKE"] = nodeBuilder{
		ArgTypes: []types.Type{types.NewMulType(types.IDENT)},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			typV := args[0].(*Const).Value.(string)
			typ, err := types.ParseType(typV)
			if err != nil {
				return nil, args[0].Pos().Error("%s", err.Error())
			}
			if !types.MAP.Equal(typ) {
				return nil, args[0].Pos().Error("expected map type, got %s", typ.String())
			}
			return &MakeNode{
				typ: typ,
			}, nil
		},
	}

	nodeBuilders["SET"] = nodeBuilder{
		ArgTypes: []types.Type{types.MAP, types.ANY, types.ANY},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			// Check types
			mapTyp := args[0].Type().(*types.MapType)
			if !mapTyp.KeyType.Equal(args[1].Type()) {
				return nil, args[1].Pos().Error("expected type %s for map key, got %s", mapTyp.KeyType.String(), args[1].Type().String())
			}
			if !mapTyp.ValType.Equal(args[2].Type()) {
				return nil, args[2].Pos().Error("expected type %s for map value, got %s", mapTyp.ValType.String(), args[2].Type().String())
			}

			return &SetNode{
				Map:   args[0],
				Key:   args[1],
				Value: args[2],
			}, nil
		},
	}

	nodeBuilders["GET"] = nodeBuilder{
		ArgTypes: []types.Type{types.MAP, types.ANY},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			// Check types
			mapTyp := args[0].Type().(*types.MapType)
			if !mapTyp.KeyType.Equal(args[1].Type()) {
				return nil, args[1].Pos().Error("expected type %s for map key, got %s", mapTyp.KeyType.String(), args[1].Type().String())
			}

			return &GetNode{
				Map: args[0],
				Key: args[1],
				typ: mapTyp.ValType,
			}, nil
		},
	}
}
