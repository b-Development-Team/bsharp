package ssagen

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/ssa"
)

func (s *SSAGen) addIf(n *ir.IfNode) ssa.ID {
	cond := s.Add(n.Condition)
	body := s.ssa.NewBlock(s.ssa.BlockName("ifbody"))
	end := s.ssa.NewBlock(s.ssa.BlockName("ifend"))
	if n.Else != nil {
		els := s.ssa.NewBlock(s.ssa.BlockName("ifelse"))
		s.blk.EndInstrutionCondJmp(cond, body, els)

		s.blk = body
		for _, node := range n.Body {
			s.Add(node)
		}
		body.EndInstrutionJmp(end)

		s.blk = els
		for _, node := range n.Else {
			s.Add(node)
		}
		els.EndInstrutionJmp(end)

		s.blk = end
	} else {
		s.blk.EndInstrutionCondJmp(cond, body, end)
		s.blk = body
		for _, node := range n.Body {
			s.Add(node)
		}
		body.EndInstrutionJmp(end)
		s.blk = end
	}
	return ssa.NullID()
}

func (s *SSAGen) addWhile(n *ir.WhileNode) ssa.ID {
	pre := s.ssa.NewBlock(s.ssa.BlockName("whilepre"))
	s.blk.EndInstrutionJmp(pre)
	s.blk = pre
	cond := s.Add(n.Condition)

	body := s.ssa.NewBlock(s.ssa.BlockName("whilebody"))
	end := s.ssa.NewBlock(s.ssa.BlockName("whileend"))
	pre.EndInstrutionCondJmp(cond, body, end)

	s.blk = body
	for _, node := range n.Body {
		s.Add(node)
	}
	body.EndInstrutionJmp(pre)

	s.blk = end
	return ssa.NullID()
}

func (s *SSAGen) addSwitch(n *ir.SwitchNode) ssa.ID {
	cond := s.Add(n.Value)
	blk := s.blk
	cases := make([]ssa.EndInstructionCase, len(n.Cases))

	end := s.ssa.NewBlock(s.ssa.BlockName("switchend"))
	def := end.Label
	for i, cs := range n.Cases {
		s.blk = s.ssa.NewBlock(s.ssa.BlockName(fmt.Sprintf("switchcase%d_", i)))
		c := cs.Block.(*ir.Case)
		for _, node := range c.Body {
			s.Add(node)
		}
		s.blk.EndInstrutionJmp(end)

		// Add case
		cases[i] = ssa.EndInstructionCase{
			Cond:  &ssa.Const{Typ: c.Value.Type(), Value: c.Value.Value},
			Label: s.blk.Label,
		}
	}

	// Add default
	if n.Default != nil {
		s.blk = s.ssa.NewBlock(s.ssa.BlockName("switchdefault"))
		for _, node := range n.Default.Block.(*ir.Default).Body {
			s.Add(node)
		}
		s.blk.EndInstrutionJmp(end)
		def = s.blk.Label
	}

	// Go to end
	s.blk = end
	blk.EndInstructionSwitch(cond, def, cases)
	return ssa.NullID()
}
