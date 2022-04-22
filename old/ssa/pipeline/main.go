package pipeline

import (
	"github.com/Nv7-Github/bsharp/old/ssa"
	"github.com/Nv7-Github/bsharp/old/ssa/constrm"
	"github.com/Nv7-Github/bsharp/old/ssa/dce"
	"github.com/Nv7-Github/bsharp/old/ssa/memrm"
)

type Pipeline struct {
	constrm bool
	dce     bool
}

func New() *Pipeline {
	return &Pipeline{}
}

func (p *Pipeline) ConstantPropagation() {
	p.constrm = true
}

func (p *Pipeline) DeadCodeElimination() {
	p.constrm = true // required for DCE
	p.dce = true
}

func (p *Pipeline) Run(s *ssa.SSA) {
	// Apply to every function
	for _, fn := range s.Funcs {
		p.Run(fn)
	}

	// Phi removal
	memrm := memrm.NewMemRM(s)
	memrm.Eval()

	// Constant folding and propagation
	if p.constrm {
		constrm.Constrm(s) // Constants
		constrm.Phirm(s)   // Constant phi nodes
	}

	if p.dce {
		dce := dce.NewDCE(s)
		dce.Remove()
	}
}
