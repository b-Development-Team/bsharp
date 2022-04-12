package memrm

import (
	"github.com/Nv7-Github/bsharp/ssa"
)

func (m *MemRM) evalBlock(label string) bool {
	d := m.blockData[label]
	if d.Filled {
		return true
	}
	b := m.ssa.Blocks[label]

	// Remove all assigns, store in block data, map all reads to that
	i := 0
	for i < len(b.Order) {
		id := b.Order[i]
		instr := b.Instructions[id]
		assign, ok := instr.(*ssa.SetVariable)
		if ok {
			d.Variables[assign.Variable] = assign.Value
			b.Remove(id)
			continue
		}

		read, ok := instr.(*ssa.GetVariable)
		if ok {
			val, exists := d.Variables[read.Variable]
			var replaceVal *ssa.ID = nil
			if !exists {
				// Step 1: Insert empty phi node
				b.Instructions[id] = &ssa.Phi{Values: make([]ssa.ID, 0)}

				// Step 2: Get vals using recurse functions
				vals := make([]ssa.ID, 0, len(b.Before))
				for _, v := range b.Before {
					vs := m.getVals(v, read.Variable)
					vals = append(vals, vs...)
				}

				// If just one val, you can replace
				if len(vals) == 1 {
					replaceVal = &vals[0]
				} else {
					// Step 3: Update using PHI node
					b.Instructions[id] = &ssa.Phi{Values: vals}
					val = id

					// Also, save the value
					d.Variables[read.Variable] = id
				}
			} else {
				replaceVal = &val
			}

			// Replace if can
			if replaceVal != nil {
				val = *replaceVal

				// Remove getvar
				b.Remove(id)

				// Replace all future occurences
				for j := i; j < len(b.Order); j++ {
					instr := b.Instructions[b.Order[j]]
					args := instr.Args()
					changed := false
					for k, arg := range args {
						if arg == id {
							args[k] = val
							changed = true
						}
					}
					if changed {
						instr.SetArgs(args)
					}
				}
			}

			continue
		}

		i++
	}

	d.Filled = true
	return false
}

func (m *MemRM) getVals(block string, variable int) []ssa.ID {
	b := m.ssa.Blocks[block]
	d := m.blockData[block]
	if !d.Filled {
		// Eval it
		m.evalBlock(block)
	}

	v, exists := d.Variables[variable]
	if exists {
		return []ssa.ID{v}
	}

	out := make([]ssa.ID, 0, len(b.Before))
	for _, v := range b.Before {
		vals := m.getVals(v, variable)
		out = append(out, vals...)
	}
	return out
}
