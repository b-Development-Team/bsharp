package ssagen

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/ssa"
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
	for _, node := range s.ir.Body {
		s.Add(node)
	}
}

func (s *SSAGen) SSA() *ssa.SSA {
	return s.ssa
}
