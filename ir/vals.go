package ir

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type CastNode struct {
	Value Node

	typ types.Type
}

func (c *CastNode) Type() types.Type { return c.typ }
func (c *CastNode) Pos() *tokens.Pos { return c.Value.Pos() }

func NewCastNode(val Node, typ types.Type) *CastNode {
	return &CastNode{
		Value: val,
		typ:   typ,
	}
}

type MathOperation int

const (
	MathOperationAdd MathOperation = iota
	MathOperationSub
	MathOperationMul
	MathOperationDiv
	MathOperationPow
)

var mathOps = map[string]MathOperation{
	"+": MathOperationAdd,
	"-": MathOperationSub,
	"*": MathOperationMul,
	"/": MathOperationDiv,
	"^": MathOperationPow,
}

type MathNode struct {
	Op  MathOperation
	Lhs Node
	Rhs Node
	typ types.Type
}

func (m *MathNode) Type() types.Type { return m.typ }

func init() {
	nodeBuilders["MATH"] = nodeBuilder{
		ArgTypes: []types.Type{types.NewMulType(types.INT, types.FLOAT), types.IDENT, types.NewMulType(types.INT, types.FLOAT)},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			op, exists := mathOps[args[1].(*Const).Value.(string)]
			if !exists {
				return nil, fmt.Errorf("unknown math operation: %s", args[1].(*Const).Value.(string))
			}

			// Get common type
			var outTyp types.Type
			if !args[0].Type().Equal(args[2].Type()) { // one is int, one is float, so cast to float
				outTyp = types.FLOAT
				if !outTyp.Equal(args[0].Type()) {
					args[0] = NewCastNode(args[0], outTyp)
				}
				if !outTyp.Equal(args[2].Type()) {
					args[2] = NewCastNode(args[2], outTyp)
				}
			} else {
				outTyp = args[0].Type() // Both are same
			}

			// Build node
			return &MathNode{
				Op:  op,
				Lhs: args[0],
				Rhs: args[2],
				typ: outTyp,
			}, nil
		},
	}
}
