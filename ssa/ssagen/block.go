package ssagen

import (
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
