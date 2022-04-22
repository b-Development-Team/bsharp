package ssagen

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/old/ssa"
	"github.com/Nv7-Github/bsharp/tokens"
)

func (s *SSAGen) addConcat(pos *tokens.Pos, n *ir.ConcatNode) ssa.ID {
	vals := make([]ssa.ID, len(n.Values))
	for i, val := range n.Values {
		vals[i] = s.Add(val)
	}
	return s.blk.AddInstruction(&ssa.Concat{
		Values: vals,
	}, pos)
}

func (s *SSAGen) addLogicalOp(pos *tokens.Pos, n *ir.LogicalOpNode) ssa.ID {
	lhs := s.Add(n.Val)
	if n.Rhs == nil {
		return s.blk.AddInstruction(&ssa.LogicalOp{
			Op:  n.Op,
			Lhs: lhs,
		}, pos)
	}
	rhs := s.Add(n.Rhs)
	return s.blk.AddInstruction(&ssa.LogicalOp{
		Op:  n.Op,
		Lhs: lhs,
		Rhs: &rhs,
	}, pos)
}
