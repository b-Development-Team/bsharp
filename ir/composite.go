package ir

import (
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

var hashable = types.NewMulType(types.INT, types.STRING, types.FLOAT) // Only types as key to map and switch

type ArrayNode struct {
	Values []Node
	typ    types.Type
}

func (a *ArrayNode) Type() types.Type { return a.typ }

type AppendNode struct {
	NullCall

	Array Node
	Value Node
}

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

	nodeBuilders["APPEND"] = nodeBuilder{
		ArgTypes: []types.Type{types.ARRAY, types.ANY},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			arrTyp := args[0].Type().(*types.ArrayType)
			if !arrTyp.ElemType.Equal(args[1].Type()) {
				return nil, args[1].Pos().Error("cannot append value of type %s to array with element type %s", args[1].Type().String(), arrTyp.ElemType.String())
			}

			return &AppendNode{
				Array: args[0],
				Value: args[1],
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
			if !hashable.Equal(typ.(*types.MapType).KeyType) {
				return nil, args[0].Pos().Error("unhashable key type: %s", typ.(*types.MapType).KeyType.String())
			}
			return &MakeNode{
				typ: typ,
			}, nil
		},
	}

	nodeBuilders["SET"] = nodeBuilder{
		ArgTypes: []types.Type{types.MAP, hashable, types.ANY},
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
		ArgTypes: []types.Type{types.MAP, hashable},
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