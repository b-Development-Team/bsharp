package ir

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

var hashable = types.NewMulType(types.INT, types.STRING, types.FLOAT) // Only types as key to map and switch

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
func (a *ArrayNode) Code(cnf CodeConfig) string {
	args := &strings.Builder{}
	for i, v := range a.Values {
		args.WriteString(v.Code(cnf))
		if i != len(a.Values)-1 {
			args.WriteString(" ")
		}
	}
	return fmt.Sprintf("[ARRAY %s]", args.String())
}

type AppendNode struct {
	NullCall

	Array Node
	Value Node
}

func (a *AppendNode) Code(cnf CodeConfig) string {
	return fmt.Sprintf("[APPEND %s %s]", a.Array.Code(cnf), a.Value.Code(cnf))
}

type IndexNode struct {
	Value Node
	Index Node
	typ   types.Type
}

func (i *IndexNode) Type() types.Type { return i.typ }

func (i *IndexNode) Code(cnf CodeConfig) string {
	return fmt.Sprintf("[INDEX %s %s]", i.Value.Code(cnf), i.Index.Code(cnf))
}

type SetIndexNode struct {
	NullCall

	Array Node
	Index Node
	Value Node
}

func (s *SetIndexNode) Code(cnf CodeConfig) string {
	return fmt.Sprintf("[SETINDEX %s %s %s]", s.Array.Code(cnf), s.Index.Code(cnf), s.Value.Code(cnf))
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
func (l *LengthNode) Code(cnf CodeConfig) string {
	return fmt.Sprintf("[LENGTH %s]", l.Value.Code(cnf))
}

type MakeNode struct {
	typ types.Type
}

func (m *MakeNode) Type() types.Type { return m.typ }
func (m *MakeNode) Code(cnf CodeConfig) string {
	return fmt.Sprintf("[MAKE %s]", m.Type().String())
}

type SetNode struct {
	NullCall

	Map   Node
	Key   Node
	Value Node
}

func (s *SetNode) Code(cnf CodeConfig) string {
	return fmt.Sprintf("[SET %s %s %s]", s.Map.Code(cnf), s.Key.Code(cnf), s.Value.Code(cnf))
}

type GetNode struct {
	Map Node
	Key Node

	typ types.Type
}

func (g *GetNode) Type() types.Type { return g.typ }

func (g *GetNode) Code(cnf CodeConfig) string {
	return fmt.Sprintf("[GET %s %s]", g.Map.Code(cnf), g.Key.Code(cnf))
}

func NewGetNode(m, k Node) *GetNode {
	return &GetNode{
		Map: m,
		Key: k,
		typ: m.Type().(*types.MapType).ValType,
	}
}

type ExistsNode struct {
	Map Node
	Key Node
}

func (g *ExistsNode) Type() types.Type { return types.BOOL }

func (e *ExistsNode) Code(cnf CodeConfig) string {
	return fmt.Sprintf("[EXISTS %s %s]", e.Map.Code(cnf), e.Key.Code(cnf))
}

type KeysNode struct {
	Map Node
	typ types.Type
}

func (k *KeysNode) Type() types.Type { return k.typ }
func (k *KeysNode) Code(cnf CodeConfig) string {
	return fmt.Sprintf("[KEYS %s]", k.Map.Code(cnf))
}

type SliceNode struct {
	Value Node
	Start Node
	End   Node
}

func (s *SliceNode) Type() types.Type {
	if types.STRING.Equal(s.Value.Type()) {
		return types.STRING // Strings are immutable
	}
	return types.NULL
}
func (s *SliceNode) Code(cnf CodeConfig) string {
	return fmt.Sprintf("[SLICE %s %s %s]", s.Value.Code(cnf), s.Start.Code(cnf), s.End.Code(cnf))
}

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
		ArgTypes: []types.Type{types.NewMulType(types.STRING, types.ARRAY, types.MAP)},
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
			if !types.ARRAY.Equal(typ) && !types.MAP.Equal(typ) {
				return nil, args[0].Pos().Error("expected map or array type, got %s", typ.String())
			}
			if types.MAP.Equal(typ) { // Check key type
				if !hashable.Equal(typ.(*types.MapType).KeyType) {
					return nil, args[0].Pos().Error("unhashable key type: %s", typ.(*types.MapType).KeyType.String())
				}
			}
			return &MakeNode{
				typ: typ,
			}, nil
		},
	}

	nodeBuilders["SETINDEX"] = nodeBuilder{
		ArgTypes: []types.Type{types.ARRAY, types.INT, types.ANY},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			// Check types
			arrTyp := args[0].Type().(*types.ArrayType)
			if !arrTyp.ElemType.Equal(args[2].Type()) {
				return nil, args[2].Pos().Error("expected type %s for array value, got %s", arrTyp.ElemType.String(), args[2].Type().String())
			}

			return &SetIndexNode{
				Array: args[0],
				Index: args[1],
				Value: args[2],
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

	nodeBuilders["EXISTS"] = nodeBuilder{
		ArgTypes: []types.Type{types.MAP, hashable},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			// Check types
			mapTyp := args[0].Type().(*types.MapType)
			if !mapTyp.KeyType.Equal(args[1].Type()) {
				return nil, args[1].Pos().Error("expected type %s for map key, got %s", mapTyp.KeyType.String(), args[1].Type().String())
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
