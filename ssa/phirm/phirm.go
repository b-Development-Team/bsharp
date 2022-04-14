// Package phirm replaces phi nodes with variable defines
package phirm

import (
	"github.com/Nv7-Github/bsharp/ssa"
)

type PhiRM struct {
	*ssa.SSA
	*stack
}

func NewPhiRM(s *ssa.SSA) *PhiRM {
	return &PhiRM{
		SSA:   s,
		stack: &stack{SSA: s, pastFrames: make([]*frame, 0), frames: make([]*frame, 0)},
	}
}

func (p *PhiRM) Remove() {
	blk := p.EntryBlock
	popOn := make(map[string][]string)
	p.Push(p.Blocks[blk])

loop:
	for {
		bl := p.Blocks[blk]
		v, exists := popOn[blk]
		if exists {
			if len(v) == 0 {
				p.Pop()
			} else {
				popOn[blk] = v[1:]
				blk = v[0]
				continue
			}
		}

		p.processBlk(bl.Label)
		switch bl.End.Type() {
		case ssa.EndInstructionTypeExit:
			break loop

		case ssa.EndInstructionTypeJmp:
			blk = bl.End.(*ssa.EndInstructionJmp).Label

		case ssa.EndInstructionTypeCondJmp:
			j := bl.End.(*ssa.EndInstructionCondJmp)
			isLoop := p.checkLoop(bl.Label)
			if isLoop {
				// If it is a loop, then the IfTrue is the body and IfFalse is the end
				blk = j.IfTrue
				popOn[j.IfFalse] = []string{}
				p.Push(bl)
			} else {
				// Its an if statement, calculate the end by following the true chain
				end := p.getIfEnd(j.IfTrue)
				popOn[end] = []string{j.IfFalse}
				p.Push(bl)

				blk = j.IfTrue
			}

		case ssa.EndInstructionTypeSwitch:
			j := bl.End.(*ssa.EndInstructionSwitch)
			end := p.getIfEnd(j.Cases[0].Label)
			cases := make([]string, 0, len(j.Cases))
			if len(j.Cases) > 1 {
				for _, cs := range j.Cases[1:] {
					cases = append(cases, cs.Label)
				}
			}
			if j.Default != end {
				cases = append(cases, j.Default)
			}
			popOn[end] = cases
			p.Push(bl)
			blk = j.Cases[0].Label
		}
	}
	p.Pop()
	p.Apply()
}

func (p *PhiRM) processBlk(name string) {
	bl := p.Blocks[name]
	for _, id := range bl.Order {
		in := bl.Instructions[id]
		phi, ok := in.(*ssa.Phi)
		if ok {
			// Insert
			p.Add(phi.Variable, phi.Values)
			// Replace with getvariable
			bl.Instructions[id] = &ssa.GetVariable{Typ: phi.Typ, Variable: phi.Variable}
		}
	}
}

func (p *PhiRM) checkLoop(blk string) bool {
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
		bl := p.Blocks[bn]
		todo = append(todo, bl.After()...)
	}
	return false
}

func (p *PhiRM) getIfEnd(blk string) string {
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
		bl := p.Blocks[bn]

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
