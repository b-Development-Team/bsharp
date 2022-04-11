package ssa

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Nv7-Github/bsharp/types"
)

type Instruction interface {
	fmt.Stringer
	Type() types.Type
}

type ID string

func NullID() ID {
	return ""
}

func (i ID) String() string {
	return string(i)
}

type Block struct {
	Parent       *SSA
	Label        string
	Instructions map[ID]Instruction
	Order        []ID // The list of instructions in the block
	End          EndInstruction

	Before []string
}

type SSA struct {
	EntryBlock   string
	Blocks       map[string]*Block
	Instructions map[ID]string // map[id]block

	// Before the memrm pass, variable data needs to be stored
	VariableTypes []types.Type

	// For IDs
	cnt int
}

func NewSSA() *SSA {
	return &SSA{
		EntryBlock:   "",
		Blocks:       make(map[string]*Block),
		Instructions: make(map[ID]string),
		cnt:          0,
	}
}

func (s *SSA) BlockName(name string) string {
	i := int64(0)
	for {
		n := name + strconv.FormatInt(i, 16)
		_, exists := s.Blocks[n]
		if !exists {
			return n
		}
		i++
	}
}

func (s *SSA) NewBlock(label string) *Block {
	if s.EntryBlock == "" {
		s.EntryBlock = label
	}
	b := &Block{
		Parent:       s,
		Label:        label,
		Instructions: make(map[ID]Instruction),
		Order:        make([]ID, 0),
	}
	s.Blocks[b.Label] = b
	return b
}

func (s *SSA) genID() ID {
	n := s.cnt
	s.cnt++
	return ID(strconv.FormatInt(int64(n), 16))
}

func (s *SSA) String() string {
	out := &strings.Builder{}
	todo := []string{s.EntryBlock}
	done := make(map[string]struct{})
	for len(todo) > 0 {
		_, exists := done[todo[0]]
		if exists {
			todo = todo[1:]
			continue
		}

		blk := s.Blocks[todo[0]]
		done[todo[0]] = struct{}{}
		todo = todo[1:]

		out.WriteString(blk.Label + ":\n")
		for _, instr := range blk.Instructions {
			out.WriteString("\t" + instr.String() + "\n")
		}
		out.WriteString("\t" + blk.End.String() + "\n")
		out.WriteString("\n")

		switch blk.End.Type() {
		case EndInstructionTypeJmp:
			todo = append(todo, blk.End.(*EndInstructionJmp).Label)

		case EndInstructionTypeCondJmp:
			j := blk.End.(*EndInstructionCondJmp)
			todo = append(todo, j.IfTrue)
			todo = append(todo, j.IfFalse)
		}
	}
	return out.String()
}

func (b *Block) AddInstruction(i Instruction) ID {
	id := b.Parent.genID()
	b.Instructions[id] = i
	b.Order = append(b.Order, id)
	b.Parent.Instructions[id] = b.Label
	return id
}

type EndInstructionType int

const (
	EndInstructionTypeJmp EndInstructionType = iota
	EndInstructionTypeCondJmp
	EndInstructionTypeExit
)

func (e EndInstructionType) Type() EndInstructionType { return e }
func (e EndInstructionType) String() string {
	return "Exit" // only time used as end instruction
}

type EndInstruction interface {
	fmt.Stringer

	Type() EndInstructionType
}

type EndInstructionJmp struct {
	Label string
}

func (e *EndInstructionJmp) Type() EndInstructionType {
	return EndInstructionTypeJmp
}

func (e *EndInstructionJmp) String() string {
	return fmt.Sprintf("Jmp (%s)", e.Label)
}

type EndInstructionCondJmp struct {
	Cond    ID
	IfTrue  string
	IfFalse string
}

func (c *EndInstructionCondJmp) Type() EndInstructionType {
	return EndInstructionTypeCondJmp
}

func (c *EndInstructionCondJmp) String() string {
	return fmt.Sprintf("CondJmp [%s](%s, %s)", c.Cond.String(), c.IfTrue, c.IfFalse)
}

func (b *Block) EndInstrutionJmp(a *Block) {
	e := &EndInstructionJmp{
		Label: a.Label,
	}
	a.Before = append(a.Before, b.Label)
	b.End = e
}

func (b *Block) EndInstrutionCondJmp(cond ID, iftrue *Block, iffalse *Block) {
	e := &EndInstructionCondJmp{
		Cond:    cond,
		IfTrue:  iftrue.Label,
		IfFalse: iffalse.Label,
	}
	iftrue.Before = append(iftrue.Before, b.Label)
	iffalse.Before = append(iffalse.Before, b.Label)
	b.End = e
}

func (b *Block) EndInstructionExit() {
	b.End = EndInstructionTypeExit
}
