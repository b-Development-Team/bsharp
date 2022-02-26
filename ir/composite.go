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
}
