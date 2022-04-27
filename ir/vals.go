package ir

import (
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

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

func (m MathOperation) String() string {
	return mathOpsNames[m]
}

type MathNode struct {
	Op  MathOperation
	Lhs Node
	Rhs Node
	typ types.Type
}

func (m *MathNode) Type() types.Type { return m.typ }
func (m *MathNode) Args() []Node     { return []Node{m.Lhs, m.Rhs} }

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

func (c CompareOperation) String() string {
	return compareOpsNames[c]
}

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
func (c *CompareNode) Args() []Node     { return []Node{c.Lhs, c.Rhs} }

type LogicalOp int

const (
	LogicalOpAnd LogicalOp = iota
	LogicalOpOr
	LogicalOpNot
)

func (l LogicalOp) String() string {
	return [...]string{"And", "Or", "Not"}[l]
}

type LogicalOpNode struct {
	Op  LogicalOp
	Val Node
	Rhs Node // nil in NOT
}

func (l *LogicalOpNode) Type() types.Type { return types.BOOL }
func (l *LogicalOpNode) Args() []Node     { return []Node{l.Val, l.Rhs} }

func NewMathNode(op MathOperation, lhs, rhs Node, typ types.Type) *MathNode {
	return &MathNode{
		Op:  op,
		Lhs: lhs,
		Rhs: rhs,
		typ: typ,
	}
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
		ArgTypes: []types.Type{types.NewMulType(types.STRING, types.INT, types.BYTE, types.FLOAT), types.IDENT, types.NewMulType(types.STRING, types.INT, types.BYTE, types.FLOAT)},
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

}
