package constrm

import (
	"github.com/Nv7-Github/bsharp/ssa"
	"github.com/Nv7-Github/bsharp/types"
)

func evalMath(blk *ssa.Block, i *ssa.Math) *ssa.Const {
	switch i.Type().BasicType() {
	case types.INT:
		return &ssa.Const{
			Value: cnst(blk, i.Lhs).(int) + cnst(blk, i.Rhs).(int),
			Typ:   types.INT,
		}

	case types.FLOAT:
		return &ssa.Const{
			Value: cnst(blk, i.Lhs).(float64) + cnst(blk, i.Rhs).(float64),
			Typ:   types.INT,
		}
	}

	return nil
}
