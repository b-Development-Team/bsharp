// Package constrm applies constant folding optimizations
package constrm

import (
	"github.com/Nv7-Github/bsharp/ssa"
)

func checkInstrConst(instr ssa.Instruction, block *ssa.Block) bool {
	for _, arg := range instr.Args() {
		_, ok := block.Instructions[arg].(*ssa.Const)
		if !ok {
			return false
		}
	}
	return true
}

func cnst(blk *ssa.Block, id ssa.ID) interface{} {
	return blk.Instructions[id].(*ssa.Const).Value
}

func globalcnst(s *ssa.SSA, id ssa.ID) *ssa.Const {
	blk := s.Blocks[s.Instructions[id].Block]
	return blk.Instructions[id].(*ssa.Const)
}

// Constrm initially converts all non-phi consts into constant values
func Constrm(s *ssa.SSA) {
	todo := []string{s.EntryBlock}
	done := make(map[string]struct{})
	for len(todo) > 0 {
		t := todo[0]
		blk := s.Blocks[t]
		todo = todo[1:]
		_, exists := done[blk.Label]
		if !exists {
			todo = append(todo, blk.After()...)
			done[blk.Label] = struct{}{}
		}

		// Constant folding
		for _, id := range blk.Order {
			instr := blk.Instructions[id]

			// TODO: The rest of the constant instructions
			switch i := instr.(type) {
			case *ssa.Math:
				if checkInstrConst(instr, blk) {
					v := evalMath(blk, i)
					if v != nil {
						blk.Instructions[id] = v
					}
				}
			}
		}
	}
}

func Phirm(s *ssa.SSA) {
	todo := []string{s.EntryBlock}
	done := make(map[string]struct{})
	for len(todo) > 0 {
		t := todo[0]
		blk := s.Blocks[t]
		todo = todo[1:]
		_, exists := done[blk.Label]
		if !exists {
			todo = append(todo, blk.After()...)
			done[blk.Label] = struct{}{}
		}

		// Constant phi removal
		for _, id := range blk.Order {
			instr := blk.Instructions[id]
			p, ok := instr.(*ssa.Phi)
			if ok {
				// Check if const
				isConst := true
				for _, val := range p.Values {
					b := s.Blocks[s.Instructions[val].Block]
					_, ok := b.Instructions[val].(*ssa.Const)
					if !ok {
						isConst = false
						break
					}
				}

				if isConst {
					isSame := true
					first := globalcnst(s, p.Values[0])
					for _, v := range p.Values[1:] {
						if globalcnst(s, v).Value != first.Value {
							isSame = false
							break
						}
					}

					if isSame {
						blk.Instructions[id] = first
					}
				}
			}
		}
	}
}
