package ssagen

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/ssa"
	"github.com/Nv7-Github/bsharp/types"
)

type SSAGen struct {
	ir  *ir.IR
	ssa *ssa.SSA
	blk *ssa.Block
}

func NewSSAGen(i *ir.IR) *SSAGen {
	s := ssa.NewSSA()
	b := s.NewBlock("entry")
	return &SSAGen{
		ir:  i,
		ssa: s,
		blk: b,
	}
}

func (s *SSAGen) Build() {
	s.ssa.VariableTypes = make([]types.Type, len(s.ir.Variables))
	for _, v := range s.ir.Variables {
		s.ssa.VariableTypes[v.ID] = v.Type
	}
	for _, node := range s.ir.Body {
		s.Add(node)
	}
	s.blk.EndInstructionExit()
	s.blk = nil
}

func (s *SSAGen) SSA() *ssa.SSA {
	return s.ssa
}
