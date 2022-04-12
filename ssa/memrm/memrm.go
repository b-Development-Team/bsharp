package memrm

import "github.com/Nv7-Github/bsharp/ssa"

// Implements memrm pass, where variables are turned into PHI nodes using https://pp.info.uni-karlsruhe.de/uploads/publikationen/braun13cc.pdf

type BlockData struct {
	Filled    bool
	Variables map[int]ssa.ID
}

func NewBlockData() *BlockData {
	return &BlockData{
		Filled:    false,
		Variables: make(map[int]ssa.ID),
	}
}

type MemRM struct {
	ssa *ssa.SSA

	blockData map[string]*BlockData
}

func NewMemRM(ssa *ssa.SSA) *MemRM {
	m := &MemRM{
		ssa:       ssa,
		blockData: make(map[string]*BlockData),
	}
	for k := range m.ssa.Blocks {
		m.blockData[k] = NewBlockData()
	}
	return m
}

func (m *MemRM) Eval() {
	todo := []string{m.ssa.EntryBlock}
	for len(todo) > 0 {
		t := todo[0]
		blk := m.ssa.Blocks[t]
		todo = todo[1:]
		done := m.evalBlock(t)

		if !done {
			switch blk.End.Type() {
			case ssa.EndInstructionTypeJmp:
				todo = append(todo, blk.End.(*ssa.EndInstructionJmp).Label)

			case ssa.EndInstructionTypeCondJmp:
				j := blk.End.(*ssa.EndInstructionCondJmp)
				todo = append(todo, j.IfTrue)
				todo = append(todo, j.IfFalse)
			}
		}
	}
}
