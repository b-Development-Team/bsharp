package ssagen

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/ssa"
)

func (s *SSAGen) addFnCall(n *ir.FnCallNode) ssa.ID {
	args := make([]ssa.ID, len(n.Params))
	for i, arg := range n.Params {
		args[i] = s.Add(arg)
	}
	return s.blk.AddInstruction(&ssa.FnCallNode{
		Fn:     s.Add(n.Fn),
		Params: args,
		Typ:    n.Type(),
	}, n.Pos())
}

func (s *SSAGen) addReturn(n *ir.ReturnNode) ssa.ID {
	s.blk.EndInstructionReturn(s.Add(n.Value))
	s.blk = nil
	return ssa.NullID()
}
