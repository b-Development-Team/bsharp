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
	b.Scope.Push(ir.ScopeTypeGlobal)
	body := b.BuildUntil(b.EntryBlock, "")
	i := b.Scope.CurrScopeInfo()
	b.Scope.Pop()

	return &ir.IR{
		GlobalScope: i,
		Body:        body,
	}
}

func (b *BSPGen) BuildUntil(start string, endBlk string) []ir.Node {
	// Traverse the SSA and try to reconstruct control flow
	blk := start
	body := make([]ir.Node, 0)

loop:
	for {
		bl := b.Blocks[blk]
		if blk == endBlk {
			break
		}

		switch bl.End.Type() {
		case ssa.EndInstructionTypeExit:
			for _, id := range bl.Order {
				n := b.Add(id)
				body = append(body, n)
			}
			break loop

		case ssa.EndInstructionTypeJmp:
			blk = bl.End.(*ssa.EndInstructionJmp).Label

		case ssa.EndInstructionTypeCondJmp:
			j := bl.End.(*ssa.EndInstructionCondJmp)
			isLoop := b.checkLoop(bl.Label)
			if isLoop {
				// If it is a loop, then the IfTrue is the body and IfFalse is the end
				cond := b.Add(j.Cond)
				b.Scope.Push(ir.ScopeTypeWhile)
				wBody := b.BuildUntil(j.IfTrue, j.IfFalse)
				i := b.Scope.CurrScopeInfo()
				b.Scope.Pop()
				body = append(body, ir.NewBlockNode(&ir.WhileNode{Condition: cond, Body: wBody, Scope: i}, b.Blocks[j.IfTrue].Pos))
				blk = j.IfFalse
			} else {
				// Its an if statement, calculate the end by following the true chain
				// end := b.getIfEnd(j.IfTrue)
				// TODO
			}
		}
	}

	return body
}

func (b *BSPGen) checkLoop(blk string) bool {
	// See if you ever end up back at blk by traversing SSA
	done := make(map[string]struct{})
	todo := []string{blk}
	doneFirst := false
	for len(todo) > 0 {
		bn := todo[0]
		todo = todo[1:]
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

func (b *BSPGen) getIfEnd(blk string) string {
	// See if you ever end up back at blk by traversing SSA
	done := make(map[string]struct{})
	todo := []string{blk}
	depth := 1 // In body
	for len(todo) > 0 {
		bn := todo[0]
		todo = todo[1:]
		if depth == 0 {
			return bn
		}

		_, exists := done[bn]
		if !exists {
			done[bn] = struct{}{}
		}
		bl := b.Blocks[bn]

		switch bl.End.Type() {
		case ssa.EndInstructionTypeSwitch, ssa.EndInstructionTypeCondJmp:
			depth++

		case ssa.EndInstructionTypeJmp:
			depth--
		}

		todo = append(todo, bl.After()...)
	}
	return ""
}
