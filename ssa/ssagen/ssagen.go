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
	fn  *ir.Function
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

	// Build functions
	if s.fn == nil { // Don't build if within function
		for _, fn := range s.ir.Funcs {
			i := &ir.IR{
				Funcs:     s.ir.Funcs,
				Variables: s.ir.Variables,
				Body:      fn.Body,
			}
			gen := NewSSAGen(i)
			gen.fn = fn
			gen.Build()
			ssa := gen.SSA()
			s.ssa.Funcs[fn.Name] = ssa
		}
	} else {
		// Build param types
		s.ssa.ParamTypes = make([]types.Type, len(s.fn.Params))
		for i, par := range s.fn.Params {
			s.ssa.ParamTypes[i] = par.Type
		}
	}

	// Build body
	for _, node := range s.ir.Body {
		s.Add(node)
	}
	s.blk.EndInstructionExit()
	s.blk = nil
}

func (s *SSAGen) SSA() *ssa.SSA {
	return s.ssa
}
