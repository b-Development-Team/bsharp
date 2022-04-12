package ssagen

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/ssa"
	"github.com/Nv7-Github/bsharp/tokens"
)

func (s *SSAGen) addDefine(pos *tokens.Pos, n *ir.DefineNode) ssa.ID {
	v := s.Add(n.Value)
	return s.blk.AddInstruction(&ssa.SetVariable{
		Variable: n.Var,
		Value:    v,
	}, pos)
}

func (s *SSAGen) addVar(pos *tokens.Pos, n *ir.VarNode) ssa.ID {
	return s.blk.AddInstruction(&ssa.GetVariable{
		Variable: n.ID,
	}, pos)
}
