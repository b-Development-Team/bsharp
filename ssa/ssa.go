package ssa

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type Instruction interface {
	fmt.Stringer
	Type() types.Type
	Args() []ID
	SetArgs([]ID)
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
	Pos          *tokens.Pos

	Before []string
}

func (b *Block) Remove(id ID) {
	delete(b.Instructions, id)
	ind := -1
	for i, val := range b.Order {
		if val == id {
			ind = i
			break
		}
	}
	copy(b.Order[ind:], b.Order[ind+1:])
	b.Order = b.Order[:len(b.Order)-1]

	delete(b.Parent.Instructions, id)
}

// Following returns a list of blocks that can directly be gone to from this block
func (b *Block) After() []string {
	switch b.End.Type() {
	case EndInstructionTypeJmp:
		return []string{b.End.(*EndInstructionJmp).Label}

	case EndInstructionTypeCondJmp:
		j := b.End.(*EndInstructionCondJmp)
		return []string{j.IfTrue, j.IfFalse}

	case EndInstructionTypeSwitch:
		j := b.End.(*EndInstructionSwitch)
		out := make([]string, 0, len(j.Cases))
		for _, cs := range j.Cases {
			out = append(out, cs.Label)
		}
		out = append(out, j.Default)
		return out
	}

	return []string{}
}

type InstructionInfo struct {
	Block string
	Pos   *tokens.Pos
}

type SSA struct {
	EntryBlock   string
	Blocks       map[string]*Block
	Instructions map[ID]InstructionInfo // map[id]block
	Funcs        map[string]*SSA

	// Before the memrm pass, variable data needs to be stored
	VariableTypes []types.Type
	ParamTypes    []types.Type // Param types if accepts params

	// For IDs
	cnt int
}

func NewSSA() *SSA {
	return &SSA{
		EntryBlock:   "",
		Blocks:       make(map[string]*Block),
		Instructions: make(map[ID]InstructionInfo),
		Funcs:        make(map[string]*SSA),
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

func (s *SSA) NewBlock(label string, pos *tokens.Pos) *Block {
	if s.EntryBlock == "" {
		s.EntryBlock = label
	}
	b := &Block{
		Parent:       s,
		Label:        label,
		Instructions: make(map[ID]Instruction),
		Order:        make([]ID, 0),
		Pos:          pos,
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
	// Add functions
	for name, fn := range s.Funcs {
		fmt.Fprintf(out, "fn %s(", name)
		for i, typ := range fn.ParamTypes {
			if i != 0 {
				out.WriteString(", ")
			}
			out.WriteString(typ.String())
		}
		out.WriteString(") {\n")
		out.WriteString(fn.String())
		out.WriteString("}\n\n")
	}

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
		for _, instr := range blk.Order {
			i := blk.Instructions[instr]
			fmt.Fprintf(out, "\t%s = %s\n", instr.String(), i.String())
		}
		out.WriteString("\t" + blk.End.String() + "\n")
		out.WriteString("\n")

		todo = append(todo, blk.After()...)
	}
	return out.String()
}

func (b *Block) AddInstruction(i Instruction, pos *tokens.Pos) ID {
	id := b.Parent.genID()
	b.Instructions[id] = i
	b.Order = append(b.Order, id)
	b.Parent.Instructions[id] = InstructionInfo{
		Block: b.Label,
		Pos:   pos,
	}
	return id
}

type EndInstructionType int

const (
	EndInstructionTypeJmp EndInstructionType = iota
	EndInstructionTypeCondJmp
	EndInstructionTypeExit
	EndInstructionTypeReturn
	EndInstructionTypeSwitch
)

func (e EndInstructionType) Type() EndInstructionType { return e }
func (e EndInstructionType) String() string {
	return [...]string{"Jmp", "CondJmp", "Exit", "Return", "Switch"}[e]
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

type EndInstructionReturn struct {
	Value ID
}

func (c *EndInstructionReturn) Type() EndInstructionType {
	return EndInstructionTypeReturn
}

func (c *EndInstructionReturn) String() string {
	return fmt.Sprintf("Return (%s)", c.Value.String())
}

type EndInstructionCase struct {
	Cond  *Const
	Label string
}

type EndInstructionSwitch struct {
	Cond    ID
	Cases   []EndInstructionCase
	Default string // nil if no default
}

func (e *EndInstructionSwitch) Type() EndInstructionType {
	return EndInstructionTypeSwitch
}

func (c *EndInstructionSwitch) String() string {
	out := &strings.Builder{}
	fmt.Fprintf(out, "Switch [%s](", c.Cond.String())
	for i, cs := range c.Cases {
		if i != 0 {
			out.WriteString(", ")
		}
		fmt.Fprintf(out, "{%s, %s}", cs.Cond, cs.Label)
	}
	out.WriteRune(')')
	return out.String()
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

func (b *Block) EndInstructionReturn(value ID) {
	b.End = &EndInstructionReturn{
		Value: value,
	}
}

func (b *Block) EndInstructionSwitch(cond ID, def string, cases []EndInstructionCase) {
	b.End = &EndInstructionSwitch{
		Cond:    cond,
		Cases:   cases,
		Default: def,
	}
	for _, cs := range cases {
		blk := b.Parent.Blocks[cs.Label]
		blk.Before = append(blk.Before, b.Label)
	}
	d := b.Parent.Blocks[def]
	d.Before = append(d.Before, b.Label)
}
