package optimize

import (
	"strconv"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

func (o *Optimizer) optimizeCast(c *ir.CastNode, pos *tokens.Pos) *Result {
	v := o.OptimizeNode(c.Value)
	if !v.IsConst {
		return &Result{
			Stmt:    ir.NewCastNode(c, c.Type()),
			IsConst: false,
		}
	}

	var out interface{}
	val := v.Stmt.(*ir.Const).Value
	switch v.Stmt.Type().BasicType() {
	case types.INT:
		switch c.Type().BasicType() {
		case types.FLOAT:
			out = float64(val.(int))

		case types.STRING:
			out = strconv.Itoa(val.(int))
		}

	case types.FLOAT:
		switch c.Type().BasicType() {
		case types.INT:
			out = int(val.(float64))

		case types.STRING:
			out = strconv.FormatFloat(val.(float64), 'f', -1, 64)
		}

	case types.STRING:
		switch c.Type().BasicType() {
		case types.INT:
			val, err := strconv.Atoi(out.(string))
			if err != nil {
				return &Result{
					Stmt:    ir.NewCastNode(c, c.Type()),
					IsConst: false,
				}
			}
			out = val

		case types.FLOAT:
			val, err := strconv.ParseFloat(val.(string), 64)
			if err != nil {
				return &Result{
					Stmt:    ir.NewCastNode(c, c.Type()),
					IsConst: false,
				}
			}
			out = val
		}

	case types.BOOL:
		switch c.Type().BasicType() {
		case types.STRING:
			out = "false"
			if val.(bool) {
				out = "true"
			}
		}
	}

	return &Result{
		Stmt:    ir.NewConst(c.Type(), pos, out),
		IsConst: true,
	}
}
