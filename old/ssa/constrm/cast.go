package constrm

import (
	"strconv"

	"github.com/Nv7-Github/bsharp/old/ssa"
	"github.com/Nv7-Github/bsharp/types"
)

func evalCast(blk *ssa.Block, i *ssa.Cast) *ssa.Const {
	if types.ANY.Equal(i.From) || types.ANY.Equal(i.To) {
		return nil
	}

	switch i.From.BasicType() {
	case types.INT:
		switch i.To.BasicType() {
		case types.FLOAT:
			return &ssa.Const{Typ: types.FLOAT, Value: float64(cnst(blk, i.Value).(int))}

		case types.STRING:
			return &ssa.Const{Typ: types.STRING, Value: strconv.Itoa(cnst(blk, i.Value).(int))}
		}
		fallthrough

	case types.FLOAT:
		switch i.To.BasicType() {
		case types.INT:
			return &ssa.Const{Typ: types.INT, Value: int(cnst(blk, i.Value).(float64))}

		case types.STRING:
			return &ssa.Const{Typ: types.STRING, Value: strconv.FormatFloat(cnst(blk, i.Value).(float64), 'f', -1, 64)}
		}
		fallthrough

	case types.STRING:
		switch i.To.BasicType() {
		case types.INT:
			val, err := strconv.Atoi(cnst(blk, i.Value).(string))
			if err != nil {
				return nil
			}
			return &ssa.Const{Typ: types.INT, Value: val}

		case types.FLOAT:
			val, err := strconv.ParseFloat(cnst(blk, i.Value).(string), 64)
			if err != nil {
				return nil
			}
			return &ssa.Const{Typ: types.FLOAT, Value: val}
		}
		fallthrough

	case types.BOOL:
		switch i.To.BasicType() {
		case types.STRING:
			val := "false"
			if cnst(blk, i.Value).(bool) {
				val = "true"
			}
			return &ssa.Const{Typ: types.STRING, Value: val}
		}
		fallthrough

	default:
		return nil
	}
}
