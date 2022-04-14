package phirm

import (
	"github.com/Nv7-Github/bsharp/ssa"
	"github.com/Nv7-Github/bsharp/types"
)

type frame struct {
	blk       *ssa.Block
	variables map[int]struct{}
	blocks    map[string]struct{}
}

func (f *frame) Remove(v int) {
	delete(f.variables, v)
}

func (f *frame) Add(v int) {
	f.variables[v] = struct{}{}
}

func (f *frame) AddBlock(blk string) {
	f.blocks[blk] = struct{}{}
}

type stack struct {
	*ssa.SSA

	pastFrames []*frame
	frames     []*frame
}

func (s *stack) Push(blk *ssa.Block) {
	s.frames = append(s.frames, &frame{
		blk:       blk,
		variables: make(map[int]struct{}),
		blocks:    map[string]struct{}{blk.Label: {}},
	})
}

func (s *stack) Pop() {
	v := s.frames[len(s.frames)-1]
	s.pastFrames = append(s.pastFrames, v)
	s.frames = s.frames[:len(s.frames)-1]
}

// Add removes it from all nested frames and adds it to this frame
func (s *stack) Add(v int, uses []ssa.ID) {
	for _, f := range s.pastFrames {
		f.Remove(v)
	}

	// Insert SetMemory
	for _, use := range uses {
		// Insert SetMemory
		pos := -1
		bl := s.Blocks[s.Instructions[use].Block]
		for i, v := range bl.Order {
			if v == use {
				pos = i
			}
		}
		v := bl.AddInstruction(&ssa.SetVariable{Variable: v, Value: use}, s.Instructions[use].Pos)
		bl.Order = bl.Order[:len(bl.Order)-1] // Remove last and instead insert in correct spot
		bl.Order = insert(bl.Order, pos+1, v)
	}

	// If one of the uses is in this frame, then it doesn't need to be redefined
	f := s.frames[len(s.frames)-1]
	contains := false
	for _, use := range uses {
		_, exists := f.blocks[s.Instructions[use].Block]
		if exists {
			contains = true
			break
		}
	}
	if !contains {
		f.Add(v)
	}
}

func zValue(typ types.Type) interface{} {
	switch typ.BasicType() {
	case types.INT:
		return int(0)

	case types.FLOAT:
		return float64(0)

	case types.BOOL:
		return false

	case types.STRING:
		return ""

	default:
		return nil
	}
}

// Apply adds variable declarations
func (s *stack) Apply() {
	for _, f := range s.pastFrames {
		prev := make([]ssa.ID, 0, len(f.variables)*2)
		for v := range f.variables {
			cnst := f.blk.AddInstruction(&ssa.Const{
				Value: zValue(s.VariableTypes[v]),
				Typ:   s.VariableTypes[v],
			}, f.blk.Pos)
			val := f.blk.AddInstruction(&ssa.SetVariable{
				Variable: v,
				Value:    cnst,
			}, f.blk.Pos)

			// Add to prev
			f.blk.Order = f.blk.Order[:len(f.blk.Order)-2] // Remove last 2
			prev = append(prev, cnst)
			prev = append(prev, val)
		}
		f.blk.Order = append(prev, f.blk.Order...)
	}
}

func insert[T any](a []T, index int, value T) []T {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}
