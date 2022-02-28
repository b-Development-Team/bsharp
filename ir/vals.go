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
	MathOperationMod
)

var mathOps = map[string]MathOperation{
	"+": MathOperationAdd,
	"-": MathOperationSub,
	"*": MathOperationMul,
	"/": MathOperationDiv,
	"^": MathOperationPow,
	"%": MathOperationMod,
}

type MathNode struct {
	Op  MathOperation
	Lhs Node
	Rhs Node
	typ types.Type
}

func (m *MathNode) Type() types.Type { return m.typ }

type CompareOperation int

const (
	CompareOperationEqual CompareOperation = iota
	CompareOperationNotEqual
	CompareOperationLess
	CompareOperationGreater
	CompareOperationLessEqual
	CompareOperationGreaterEqual
)

var compareOps = map[string]CompareOperation{
	"==": CompareOperationEqual,
	"!=": CompareOperationNotEqual,
	"<":  CompareOperationLess,
	">":  CompareOperationGreater,
	"<=": CompareOperationLessEqual,
	">=": CompareOperationGreaterEqual,
}

type CompareNode struct {
	Op  CompareOperation
	Lhs Node
	Rhs Node
}

func (c *CompareNode) Type() types.Type { return types.BOOL }

type MathFunction int

const (
	MathFunctionCeil MathFunction = iota
	MathFunctionFloor
	MathFunctionRound
)

type MathFunctionNode struct {
	Func MathFunction
	Arg  Node
	typ  types.Type
}

func (m *MathFunctionNode) Type() types.Type { return m.typ }

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

	nodeBuilders["COMPARE"] = nodeBuilder{
		ArgTypes: []types.Type{types.NewMulType(types.STRING, types.INT, types.FLOAT), types.IDENT, types.NewMulType(types.STRING, types.INT, types.FLOAT)},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			// Get op
			op, exists := compareOps[args[1].(*Const).Value.(string)]
			if !exists {
				return nil, fmt.Errorf("unknown math operation: %s", args[1].(*Const).Value.(string))
			}

			// Check if can compare
			if !args[0].Type().Equal(args[2].Type()) {
				return nil, fmt.Errorf("cannot compare type %s to type %s", args[0].Type(), args[2].Type())
			}
			// Return node
			return &CompareNode{
				Op:  op,
				Lhs: args[0],
				Rhs: args[2],
			}, nil
		},
	}

	nodeBuilders["FLOAT"] = nodeBuilder{
		ArgTypes: []types.Type{types.NewMulType(types.INT, types.STRING)},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return NewCastNode(args[0], types.FLOAT), nil
		},
	}

	nodeBuilders["STRING"] = nodeBuilder{
		ArgTypes: []types.Type{types.NewMulType(types.INT, types.FLOAT, types.BOOL)},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return NewCastNode(args[0], types.STRING), nil
		},
	}

	nodeBuilders["INT"] = nodeBuilder{
		ArgTypes: []types.Type{types.NewMulType(types.FLOAT, types.STRING)},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return NewCastNode(args[0], types.INT), nil
		},
	}

	nodeBuilders["FLOOR"] = nodeBuilder{
		ArgTypes: []types.Type{types.FLOAT},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return &MathFunctionNode{
				Func: MathFunctionFloor,
				Arg:  args[0],
				typ:  types.INT,
			}, nil
		},
	}

	nodeBuilders["CEIL"] = nodeBuilder{
		ArgTypes: []types.Type{types.FLOAT},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return &MathFunctionNode{
				Func: MathFunctionCeil,
				Arg:  args[0],
				typ:  types.INT,
			}, nil
		},
	}

	nodeBuilders["ROUND"] = nodeBuilder{
		ArgTypes: []types.Type{types.FLOAT},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return &MathFunctionNode{
				Func: MathFunctionRound,
				Arg:  args[0],
				typ:  types.INT,
			}, nil
		},
	}
}