package ssagen

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/old/ssa"
)

func (s *SSAGen) newIRNode(kind ssa.IRNode, n *ir.CallNode) ssa.ID {
	a := n.Call.Args()
	pars := make([]ssa.ID, len(a))
	for i, par := range a {
		pars[i] = s.Add(par)
	}
	return s.blk.AddInstruction(&ssa.IRValue{
		Kind:   kind,
		Params: pars,
		Typ:    n.Type(),
	}, n.Pos())
}

func (s *SSAGen) newLiveIRNode(kind ssa.LiveIRNode, n *ir.CallNode) ssa.ID {
	a := n.Call.Args()
	pars := make([]ssa.ID, len(a))
	for i, par := range a {
		pars[i] = s.Add(par)
	}
	return s.blk.AddInstruction(&ssa.LiveIRValue{
		Kind:   kind,
		Params: pars,
		Typ:    n.Type(),
	}, n.Pos())
}
