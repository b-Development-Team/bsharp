package ir

import (
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

var mathOpsNames = make(map[MathOperation]string)

func init() {
	for k, v := range mathOps {
		mathOpsNames[v] = k
	}
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

var compareOpsNames = make(map[CompareOperation]string)

func init() {
	for k, v := range compareOps {
		compareOpsNames[v] = k
	}
}

type CompareNode struct {
	Op  CompareOperation
	Lhs Node
	Rhs Node
}

func (c *CompareNode) Type() types.Type { return types.BOOL }

type LogicalOp int

const (
	LogicalOpAnd LogicalOp = iota
	LogicalOpOr
	LogicalOpNot
)

type LogicalOpNode struct {
	Op  LogicalOp
	Val Node
	Rhs Node // nil in NOT
}

func (l *LogicalOpNode) Type() types.Type { return types.BOOL }

func NewMathNode(op MathOperation, lhs, rhs Node, typ types.Type) *MathNode {
	return &MathNode{
		Op:  op,
		Lhs: lhs,
		Rhs: rhs,
		typ: typ,
	}
}

type CanCastNode struct {
	Value Node
	Typ   types.Type
}

func (c *CanCastNode) Type() types.Type {
	return types.BOOL
}

func init() {
	nodeBuilders["MATH"] = nodeBuilder{
		ArgTypes: []types.Type{types.NewMulType(types.INT, types.FLOAT), types.IDENT, types.NewMulType(types.INT, types.FLOAT)},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			op, exists := mathOps[args[1].(*Const).Value.(string)]
			if !exists {
				return nil, args[1].Pos().Error("unknown math operation: %s", args[1].(*Const).Value.(string))
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
				return nil, args[1].Pos().Error("unknown compare operation: %s", args[1].(*Const).Value.(string))
			}

			// Check if can compare
			if !args[0].Type().Equal(args[2].Type()) {
				return nil, pos.Error("cannot compare type %s to type %s", args[0].Type(), args[2].Type())
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

	nodeBuilders["AND"] = nodeBuilder{
		ArgTypes: []types.Type{types.BOOL, types.BOOL},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return &LogicalOpNode{
				Op:  LogicalOpAnd,
				Val: args[0],
				Rhs: args[1],
			}, nil
		},
	}

	nodeBuilders["OR"] = nodeBuilder{
		ArgTypes: []types.Type{types.BOOL, types.BOOL},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return &LogicalOpNode{
				Op:  LogicalOpOr,
				Val: args[0],
				Rhs: args[1],
			}, nil
		},
	}

	nodeBuilders["NOT"] = nodeBuilder{
		ArgTypes: []types.Type{types.BOOL},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return &LogicalOpNode{
				Op:  LogicalOpNot,
				Val: args[0],
			}, nil
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
			typ, err := types.ParseType(args[1].(*Const).Value.(string))
			if err != nil {
				return nil, err
			}
			return &CanCastNode{
				Value: args[0],
				Typ:   typ,
			}, nil
		},
	}

	nodeBuilders["CAST"] = nodeBuilder{
		ArgTypes: []types.Type{types.ALL, types.IDENT},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			typ, err := types.ParseType(args[1].(*Const).Value.(string))
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
