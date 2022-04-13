package ssagen

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/ssa"
	"github.com/Nv7-Github/bsharp/tokens"
)

func (s *SSAGen) addDefine(pos *tokens.Pos, n *ir.DefineNode) ssa.ID {
	v := s.Add(n.Value)
	va := s.ir.Variables[n.Var]
	if va.NeedsGlobal {
		return s.blk.AddInstruction(&ssa.GlobalSetVariable{
			Variable: n.Var,
			Value:    v,
		}, pos)
	}
	return s.blk.AddInstruction(&ssa.SetVariable{
		Variable: n.Var,
		Value:    v,
	}, pos)
}

func (s *SSAGen) addVar(pos *tokens.Pos, n *ir.VarNode) ssa.ID {
	v := s.ir.Variables[n.ID]
	if s.fn != nil {
		isParam := true
		for _, par := range s.fn.Params {
			if par.ID != n.ID {
				isParam = false
				break
			}
		}
		if isParam {
			return s.blk.AddInstruction(&ssa.GetParam{
				Variable: n.ID,
			}, pos)
		}
	}
	if v.NeedsGlobal {
		return s.blk.AddInstruction(&ssa.GlobalGetVariable{
			Variable: n.ID,
		}, pos)
	}
	return s.blk.AddInstruction(&ssa.GetVariable{
		Variable: n.ID,
	}, pos)
}
