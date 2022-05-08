package ir

import (
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

var hashable = types.NewMulType(types.INT, types.BYTE, types.STRING, types.FLOAT) // Only types as key to map and switch

type ArrayNode struct {
	Values []Node
	typ    types.Type
}

func NewArrayNode(vals []Node, typ types.Type) *ArrayNode {
	return &ArrayNode{
		Values: vals,
		typ:    typ,
	}
}

func (a *ArrayNode) Type() types.Type { return a.typ }
func (a *ArrayNode) Args() []Node     { return a.Values }

type AppendNode struct {
	NullCall

	Array Node
	Value Node
}

func (a *AppendNode) Args() []Node { return []Node{a.Array, a.Value} }

type IndexNode struct {
	Value Node
	Index Node
	typ   types.Type
}

func (i *IndexNode) Type() types.Type { return i.typ }
func (i *IndexNode) Args() []Node     { return []Node{i.Value, i.Index} }

type SetIndexNode struct {
	NullCall

	Array Node
	Index Node
	Value Node
}

func (s *SetIndexNode) Args() []Node {
	return []Node{s.Array, s.Index, s.Value}
}

type SetStructNode struct {
	NullCall

	Struct Node
	Field  int
	Value  Node
	pos    *tokens.Pos
}

func (s *SetStructNode) Args() []Node {
	return []Node{s.Struct, s.Value, NewConst(types.INT, s.pos, s.Field)}
}

func NewIndexNode(val, index Node) *IndexNode {
	outTyp := types.Type(types.STRING)
	if types.ARRAY.Equal(val.Type()) {
		outTyp = val.Type().(*types.ArrayType).ElemType
	}
	return &IndexNode{
		Value: val,
		Index: index,
		typ:   outTyp,
	}
}

type LengthNode struct {
	Value Node
}

func (l *LengthNode) Type() types.Type { return types.INT }
func (l *LengthNode) Args() []Node     { return []Node{l.Value} }

type MakeNode struct {
	typ types.Type
}

func (m *MakeNode) Type() types.Type { return m.typ }
func (m *MakeNode) Args() []Node     { return []Node{} }

type SetNode struct {
	NullCall

	Map   Node
	Key   Node
	Value Node
}

func (s *SetNode) Args() []Node { return []Node{s.Map, s.Key, s.Value} }

type GetNode struct {
	Map Node
	Key Node

	typ types.Type
}

func (g *GetNode) Type() types.Type { return g.typ }
func (g *GetNode) Args() []Node     { return []Node{g.Map, g.Key} }

func NewGetNode(m, k Node) *GetNode {
	return &GetNode{
		Map: m,
		Key: k,
		typ: m.Type().(*types.MapType).ValType,
	}
}

type GetStructNode struct {
	Struct Node
	Field  int
	pos    *tokens.Pos

	typ types.Type
}

func (g *GetStructNode) Type() types.Type { return g.typ }
func (g *GetStructNode) Args() []Node     { return []Node{g.Struct, NewConst(types.INT, g.pos, g.Field)} }

type ExistsNode struct {
	Map Node
	Key Node
}

func (g *ExistsNode) Type() types.Type { return types.BOOL }
func (g *ExistsNode) Args() []Node     { return []Node{g.Map, g.Key} }

type KeysNode struct {
	Map Node
	typ types.Type
}

func (k *KeysNode) Type() types.Type { return k.typ }
func (k *KeysNode) Args() []Node     { return []Node{k.Map} }

type SliceNode struct {
	Value Node
	Start Node
	End   Node
}

func (s *SliceNode) Args() []Node { return []Node{s.Value, s.Start, s.End} }
func (s *SliceNode) Type() types.Type {
	if types.STRING.Equal(s.Value.Type()) {
		return types.STRING // Strings are immutable
	}
	return types.NULL
}

func init() {
	nodeBuilders["ARRAY"] = nodeBuilder{
		ArgTypes: []types.Type{types.ALL, types.VARIADIC},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			// Check types
			typ := args[0].Type()
			if len(args) > 1 {
				for _, arg := range args[1:] {
					if !typ.Equal(arg.Type()) {
						b.Error(ErrorLevelError, arg.Pos(), "expected type %s, got %s", typ.String(), arg.Type().String())
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
		ArgTypes: []types.Type{types.ARRAY, types.ALL},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			arrTyp, ok := args[0].Type().(*types.ArrayType)
			if !ok { // invalid
				return NewTypedValue(types.INVALID), nil
			}
			if !arrTyp.ElemType.Equal(args[1].Type()) {
				b.Error(ErrorLevelError, args[1].Pos(), "cannot append value of type %s to array with element type %s", args[1].Type().String(), arrTyp.ElemType.String())
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
				outTyp = types.BYTE
			}

			return &IndexNode{
				Value: args[0],
				Index: args[1],
				typ:   outTyp,
			}, nil
		},
	}

	nodeBuilders["LENGTH"] = nodeBuilder{
		ArgTypes: []types.Type{types.NewMulType(types.STRING, types.ARRAY, types.MAP)},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return &LengthNode{
				Value: args[0],
			}, nil
		},
	}

	nodeBuilders["MAKE"] = nodeBuilder{
		ArgTypes: []types.Type{types.IDENT},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			typV := args[0].(*Const).Value.(string)
			typ, err := types.ParseType(typV, b.typeNames)
			if err != nil {
				b.Error(ErrorLevelError, args[0].Pos(), "%s", err.Error())
				return NewTypedValue(types.INVALID), nil
			}
			if !types.ARRAY.Equal(typ) && !types.MAP.Equal(typ) && !types.STRUCT.Equal(typ) {
				b.Error(ErrorLevelError, args[0].Pos(), "expected map, array, or struct type, got %s", typ.String())
				return NewTypedValue(types.INVALID), nil
			}
			if types.MAP.Equal(typ) { // Check key type
				if !hashable.Equal(typ.(*types.MapType).KeyType) {
					b.Error(ErrorLevelError, args[0].Pos(), "unhashable key type: %s", typ.(*types.MapType).KeyType.String())
				}
			}
			return &MakeNode{
				typ: typ,
			}, nil
		},
	}

	nodeBuilders["SET"] = nodeBuilder{
		ArgTypes: []types.Type{types.NewMulType(types.MAP, types.ARRAY, types.STRUCT), types.NewMulType(hashable, types.IDENT), types.ALL},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			// If map type
			if types.MAP.Equal(args[0].Type()) {
				// Check types
				mapTyp := args[0].Type().(*types.MapType)
				if !mapTyp.KeyType.Equal(args[1].Type()) {
					b.Error(ErrorLevelError, args[1].Pos(), "expected type %s for map key, got %s", mapTyp.KeyType.String(), args[1].Type().String())
				}
				if !mapTyp.ValType.Equal(args[2].Type()) {
					b.Error(ErrorLevelError, args[2].Pos(), "expected type %s for map value, got %s", mapTyp.ValType.String(), args[2].Type().String())
				}

				return &SetNode{
					Map:   args[0],
					Key:   args[1],
					Value: args[2],
				}, nil
			}

			// If struct type
			if types.STRUCT.Equal(args[0].Type()) {
				// Check types
				structTyp := args[0].Type().(*types.StructType)
				if !types.IDENT.Equal(args[1].Type()) {
					b.Error(ErrorLevelError, args[1].Pos(), "expected type %s for struct field name, got %s", types.IDENT.String(), args[1].Type().String())
					return NewTypedValue(types.INVALID), nil
				}
				name := args[1].(*Const).Value.(string)
				var field *types.StructField
				id := -1
				for i, f := range structTyp.Fields {
					if f.Name == name {
						field = &f
						id = i
						break
					}
				}
				if field == nil {
					b.Error(ErrorLevelError, args[1].Pos(), "unknown struct field: %s", name)
				} else if !field.Type.Equal(args[2].Type()) {
					b.Error(ErrorLevelError, args[2].Pos(), "expected type %s for struct field value, got %s", field.Type.String(), args[2].Type().String())
				}

				return &SetStructNode{
					Struct: args[0],
					Field:  id,
					Value:  args[2],
					pos:    args[1].Pos(),
				}, nil
			}

			// Array type
			arrTyp := args[0].Type().(*types.ArrayType)
			if !types.INT.Equal(args[1].Type()) {
				b.Error(ErrorLevelError, args[1].Pos(), "expected type %s for array index, got %s", types.INT.String(), args[1].Type().String())
			}
			if !arrTyp.ElemType.Equal(args[2].Type()) {
				b.Error(ErrorLevelError, args[2].Pos(), "expected type %s for array value, got %s", arrTyp.ElemType.String(), args[2].Type().String())
			}
			return &SetIndexNode{
				Array: args[0],
				Index: args[1],
				Value: args[2],
			}, nil
		},
	}

	nodeBuilders["GET"] = nodeBuilder{
		ArgTypes: []types.Type{types.NewMulType(types.MAP, types.STRUCT), types.NewMulType(hashable, types.IDENT)},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			// If map type
			if types.MAP.Equal(args[0].Type()) {
				// Check types
				mapTyp := args[0].Type().(*types.MapType)
				if !mapTyp.KeyType.Equal(args[1].Type()) {
					b.Error(ErrorLevelError, args[1].Pos(), "expected type %s for map key, got %s", mapTyp.KeyType.String(), args[1].Type().String())
				}

				return &GetNode{
					Map: args[0],
					Key: args[1],
					typ: mapTyp.ValType,
				}, nil
			}

			// Struct type
			structTyp, ok := args[0].Type().(*types.StructType)
			if !ok { // invalid
				return NewTypedValue(types.INVALID), nil
			}
			if !types.IDENT.Equal(args[1].Type()) {
				b.Error(ErrorLevelError, args[1].Pos(), "expected type %s for struct field name, got %s", types.IDENT.String(), args[1].Type().String())
				return NewTypedValue(types.INVALID), nil
			}
			name := args[1].(*Const).Value.(string)
			var field *types.StructField
			id := -1
			for i, f := range structTyp.Fields {
				if f.Name == name {
					field = &f
					id = i
					break
				}
			}
			if field == nil {
				b.Error(ErrorLevelError, args[1].Pos(), "unknown struct field: %s", name)
				return NewTypedValue(types.INVALID), nil
			}
			return &GetStructNode{
				Struct: args[0],
				Field:  id,
				typ:    field.Type,
				pos:    args[1].Pos(),
			}, nil
		},
	}

	nodeBuilders["EXISTS"] = nodeBuilder{
		ArgTypes: []types.Type{types.MAP, hashable},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			// Check types
			mapTyp := args[0].Type().(*types.MapType)
			if !mapTyp.KeyType.Equal(args[1].Type()) {
				b.Error(ErrorLevelError, args[1].Pos(), "expected type %s for map key, got %s", mapTyp.KeyType.String(), args[1].Type().String())
			}

			return &ExistsNode{
				Map: args[0],
				Key: args[1],
			}, nil
		},
	}

	nodeBuilders["KEYS"] = nodeBuilder{
		ArgTypes: []types.Type{types.MAP},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return &KeysNode{
				Map: args[0],
				typ: types.NewArrayType(args[0].Type().(*types.MapType).KeyType),
			}, nil
		},
	}

	nodeBuilders["SLICE"] = nodeBuilder{
		ArgTypes: []types.Type{types.NewMulType(types.STRING, types.ARRAY), types.INT, types.INT},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return &SliceNode{
				Value: args[0],
				Start: args[1],
				End:   args[2],
			}, nil
		},
	}
}
