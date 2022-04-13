package dce

import (
	"github.com/Nv7-Github/bsharp/ssa"
)

type DCE struct {
	*ssa.SSA

	notDead map[ssa.ID]struct{}
}

func NewDCE(s *ssa.SSA) *DCE {
	return &DCE{
		SSA:     s,
		notDead: make(map[ssa.ID]struct{}),
	}
}

func (d *DCE) Remove() {
	todo := []string{d.EntryBlock}
	done := make(map[string]struct{})
	for len(todo) > 0 {
		t := todo[0]
		blk := d.Blocks[t]
		todo = todo[1:]
		_, exists := done[blk.Label]
		if !exists {
			// Calc whats after
			switch blk.End.Type() {
			case ssa.EndInstructionTypeJmp:
				todo = append(todo, blk.End.(*ssa.EndInstructionJmp).Label)

			case ssa.EndInstructionTypeCondJmp:
				j := blk.End.(*ssa.EndInstructionCondJmp)
				cond := blk.Instructions[j.Cond]
				c, ok := cond.(*ssa.Const)
				if ok {
					if c.Value.(bool) {
						todo = append(todo, j.IfTrue)
						d.removeBlockBefore(j.IfFalse, blk.Label)
						blk.EndInstrutionJmp(d.Blocks[j.IfTrue])
					} else {
						todo = append(todo, j.IfFalse)
						d.removeBlockBefore(j.IfTrue, blk.Label)
						blk.EndInstrutionJmp(d.Blocks[j.IfFalse])
					}
				} else {
					todo = append(todo, j.IfTrue, j.IfFalse)
				}

			case ssa.EndInstructionTypeSwitch: // TODO: Actually check all constants and resolve all block references
				todo = append(todo, blk.After()...)
			}

			done[blk.Label] = struct{}{}
		}

		// Mark not dead
		for _, id := range blk.Order {
			instr := blk.Instructions[id]

			switch instr.(type) {
			case *ssa.LiveIRValue, *ssa.GlobalSetVariable, *ssa.FnCall:
				d.markNotDead(id)
			}
		}

		switch blk.End.Type() {
		case ssa.EndInstructionTypeCondJmp:
			d.markNotDead(blk.End.(*ssa.EndInstructionCondJmp).Cond)

		case ssa.EndInstructionTypeSwitch:
			d.markNotDead(blk.End.(*ssa.EndInstructionSwitch).Cond)

		case ssa.EndInstructionTypeReturn:
			d.markNotDead(blk.End.(*ssa.EndInstructionReturn).Value)
		}
	}

	// Calculate dead blocks
	deadBlocks := make(map[string]struct{})
	for _, blk := range d.Blocks {
		deadBlocks[blk.Label] = struct{}{}
	}
	for instr := range d.notDead {
		delete(deadBlocks, d.Instructions[instr].Block)
	}

	// Remove dead blocks
	todo = []string{d.EntryBlock}
	done = make(map[string]struct{})
	for len(todo) > 0 {
		blk := todo[0]
		todo = todo[1:]
		_, exists := deadBlocks[blk]
		if !exists {
			continue
		}

		b, exists := d.Blocks[blk]
		if !exists {
			continue
		}
		_, exists = done[b.Label]
		if !exists {
			todo = append(todo, b.After()...)
			done[b.Label] = struct{}{}
		}

		// Replace jumps and entry
		next := b.End.(*ssa.EndInstructionJmp).Label
		for _, blk := range b.Before {
			bl, exists := d.Blocks[blk]
			if !exists {
				continue
			}

			// Remove from after
			switch bl.End.Type() {
			case ssa.EndInstructionTypeJmp:
				bl.EndInstrutionJmp(d.Blocks[next])

			case ssa.EndInstructionTypeCondJmp:
				j := bl.End.(*ssa.EndInstructionCondJmp)
				if j.IfTrue == b.Label {
					j.IfTrue = next
				} else {
					j.IfFalse = next
				}

			case ssa.EndInstructionTypeSwitch:
				j := bl.End.(*ssa.EndInstructionSwitch)
				for _, cs := range j.Cases {
					if cs.Label == b.Label {
						cs.Label = next
					}
				}
				if j.Default == b.Label {
					j.Default = next
				}
			}
		}
		if d.EntryBlock == b.Label {
			d.EntryBlock = next
		}

		for _, blk := range b.After() {
			bl, exists := d.Blocks[blk]
			if !exists {
				continue
			}
			d.removeBlockBefore(bl.Label, b.Label)
		}

		// Remove block and pointers to instructions within the block
		delete(d.Blocks, blk)
		for id, instr := range d.Instructions {
			if instr.Block == blk {
				delete(d.Instructions, id)
			}
		}
	}

	// Remove dead instructions
	for _, bl := range d.Blocks {
		newOrder := make([]ssa.ID, 0, len(bl.Order))
		for _, id := range bl.Order {
			_, exists := d.notDead[id]
			if !exists {
				delete(bl.Instructions, id)
				delete(d.Instructions, id)
			} else {
				newOrder = append(newOrder, id)
			}
		}
		bl.Order = newOrder
	}
}

func (d *DCE) markNotDead(id ssa.ID) {
	_, exists := d.notDead[id]
	if exists {
		return
	}

	b := d.Blocks[d.Instructions[id].Block]
	i := b.Instructions[id]
	d.notDead[id] = struct{}{}
	for _, arg := range i.Args() {
		d.markNotDead(arg)
	}
}

func (d *DCE) removeBlockBefore(block string, val string) {
	bl := d.Blocks[block]
	before := make(map[string]struct{}, len(bl.Before))
	for _, val := range bl.Before {
		before[val] = struct{}{}
	}
	delete(before, val)
	bef := make([]string, len(before))
	i := 0
	for k := range before {
		bef[i] = k
		i++
	}
	bl.Before = bef
}
