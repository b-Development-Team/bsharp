package ssagen

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/ssa"
)

func (s *SSAGen) addDefine(n *ir.DefineNode) ssa.ID {
	v := s.Add(n.Value)
	return s.blk.AddInstruction(&ssa.SetVariable{
		Variable: n.Var,
		Value:    v,
	})
}

func (s *SSAGen) addVar(n *ir.VarNode) ssa.ID {
	return s.blk.AddInstruction(&ssa.GetVariable{
		Variable: n.ID,
	})
}
