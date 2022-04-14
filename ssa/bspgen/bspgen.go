package bspgen

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/ssa"
)

type BSPGen struct {
	*ssa.SSA
	Scope *ir.Scope
	ir    *ir.IR // Original IR data
}

func NewBSPGen(s *ssa.SSA, i *ir.IR) *BSPGen {
	return &BSPGen{
		SSA:   s,
		Scope: ir.NewScope(),
	}
}

func (b *BSPGen) Build() *ir.IR {
	// Traverse the SSA and try to reconstruct control flow
	blk := b.EntryBlock
	b.Scope.Push(ir.ScopeTypeGlobal)
	body := make([]ir.Node, 0)

loop:
	for {
		bl := b.Blocks[blk]
		switch bl.End.Type() {
		case ssa.EndInstructionTypeExit:
			for _, id := range bl.Order {
				n := b.Add(id)
				body = append(body, n)
			}
			break loop

		case ssa.EndInstructionTypeJmp:
			// Calc whether loop or not
		}
	}
	i := b.Scope.CurrScopeInfo()
	b.Scope.Pop()

	return &ir.IR{
		GlobalScope: i,
		Body:        body,
	}
}

func (b *BSPGen) checkLoop(blk string) bool {
	// See if you ever end up back at blk by traversing SSA
	done := make(map[string]struct{})
	todo := []string{blk}
	doneFirst := false
	for len(todo) > 0 {
		bn := todo[0]
		if bn == blk && doneFirst {
			return true
		}
		if !doneFirst {
			doneFirst = true
		}
		_, exists := done[bn]
		if !exists {
			done[bn] = struct{}{}
		}
		bl := b.Blocks[bn]
		todo = append(todo, bl.After()...)
	}
	return false
}
