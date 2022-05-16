package bsp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

func (b *BSP) addFnCall(fn *ir.FnCallNode) (string, error) {
	name, err := b.buildNode(fn.Fn)
	if err != nil {
		return "", err
	}

	out := &strings.Builder{}
	fmt.Fprintf(out, "[CALL%s%s", b.Config.Seperator, name)
	for _, arg := range fn.Params {
		v, err := b.buildNode(arg)
		if err != nil {
			return "", err
		}
		out.WriteString(b.Config.Seperator)
		out.WriteString(v)
	}
	out.WriteRune(']')
	return out.String(), nil
}

func (b *BSP) addCall(name string, val ir.Call) (string, error) {
	out := &strings.Builder{}
	fmt.Fprintf(out, "[%s", name)
	for _, arg := range val.Args() {
		v, err := b.buildNode(arg)
		if err != nil {
			return "", err
		}
		out.WriteString(b.Config.Seperator)
		out.WriteString(v)
	}
	out.WriteRune(']')
	return out.String(), nil
}

func (b *BSP) addConst(v *ir.Const) (string, error) {
	switch v.Type() {
	case types.INT:
		return strconv.Itoa(v.Value.(int)), nil

	case types.BYTE:
		switch v.Value.(byte) {
		case '\n':
			return `'\n'`, nil

		case '\'':
			return `'\''`, nil

		case '\\':
			return `'\\'`, nil

		case '\t':
			return `'\t'`, nil

		default:
			return "'" + string(v.Value.(byte)) + "'", nil
		}

	case types.FLOAT:
		v := strconv.FormatFloat(v.Value.(float64), 'f', -1, 64)
		if !strings.Contains(v, ".") {
			v += ".0"
		}
		return v, nil

	case types.STRING:
		return fmt.Sprintf("%q", v.Value.(string)), nil

	case types.IDENT:
		return fmt.Sprintf("%s", v.Value.(string)), nil

	case types.BOOL:
		if v.Value.(bool) {
			return "TRUE", nil
		}
		return "FALSE", nil

	case types.NULL:
		return "NULL", nil

	default:
		return "", v.Pos().Error("unknown const type: %s", v.Type())
	}
}

func (b *BSP) addCast(n *ir.CastNode) (string, error) {
	if types.ANY.Equal(n.Value.Type()) {
		return b.addCall("CAST", n)
	}

	v, err := b.buildNode(n.Value)
	if err != nil {
		return "", err
	}

	switch n.Type().BasicType() {
	case types.ANY:
		return fmt.Sprintf("[ANY%s%s]", b.Config.Seperator, v), nil

	case types.STRING:
		return fmt.Sprintf("[STRING%s%s]", b.Config.Seperator, v), nil

	case types.FLOAT:
		return fmt.Sprintf("[FLOAT%s%s]", b.Config.Seperator, v), nil

	case types.INT:
		return fmt.Sprintf("[INT%s%s]", b.Config.Seperator, v), nil

	case types.BYTE:
		return fmt.Sprintf("[BYTE%s%s]", b.Config.Seperator, v), nil
	}

	return "", n.Pos().Error("cannot cast from %s to %s", n.Value.Type(), n.Type())
}
