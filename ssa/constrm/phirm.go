package constrm

import "github.com/Nv7-Github/bsharp/ssa"

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
