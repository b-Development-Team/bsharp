package ssagen

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/ssa"
	"github.com/Nv7-Github/bsharp/tokens"
)

func (s *SSAGen) addConst(n *ir.Const) ssa.ID {
	return s.blk.AddInstruction(&ssa.Const{
		Value: n.Value,
		Typ:   n.Type(),
	}, n.Pos())
}

func (s *SSAGen) addCompare(pos *tokens.Pos, n *ir.CompareNode) ssa.ID {
	return s.blk.AddInstruction(&ssa.Compare{
		Op:  n.Op,
		Typ: n.Type(),
		Lhs: s.Add(n.Lhs),
		Rhs: s.Add(n.Rhs),
	}, pos)
}

func (s *SSAGen) addMath(pos *tokens.Pos, n *ir.MathNode) ssa.ID {
	return s.blk.AddInstruction(&ssa.Math{
		Op:  n.Op,
		Typ: n.Type(),
		Lhs: s.Add(n.Lhs),
		Rhs: s.Add(n.Rhs),
	}, pos)
}
