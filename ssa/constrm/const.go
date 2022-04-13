package constrm

import (
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/ssa"
	"github.com/Nv7-Github/bsharp/types"
)

func evalConcat(blk *ssa.Block, i *ssa.Concat) *ssa.Const {
	out := &strings.Builder{}
	for _, par := range i.Values {
		out.WriteString(cnst(blk, par).(string))
	}
	return &ssa.Const{Typ: types.STRING, Value: out.String()}
}

func evalLogicalOp(blk *ssa.Block, i *ssa.LogicalOp) *ssa.Const {
	switch i.Op {
	case ir.LogicalOpAnd:
		return &ssa.Const{Typ: types.BOOL, Value: cnst(blk, i.Lhs).(bool) && cnst(blk, *i.Rhs).(bool)}

	case ir.LogicalOpOr:
		return &ssa.Const{Typ: types.BOOL, Value: cnst(blk, i.Lhs).(bool) || cnst(blk, *i.Rhs).(bool)}

	case ir.LogicalOpNot:
		return &ssa.Const{Typ: types.BOOL, Value: !cnst(blk, i.Lhs).(bool)}

	default:
		return nil
	}
}
