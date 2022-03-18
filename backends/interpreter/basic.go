package interpreter

import (
	"strings"
	"time"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

func (i *Interpreter) evalPrint(c *ir.PrintNode) error {
	v, err := i.evalNode(c.Arg)
	if err != nil {
		return err
	}
	_, err = i.stdout.Write([]byte(v.Value.(string) + "\n"))
	return err
}

func (i *Interpreter) evalConcat(c *ir.ConcatNode) (*Value, error) {
	out := &strings.Builder{}
	for _, arg := range c.Values {
		v, err := i.evalNode(arg)
		if err != nil {
			return nil, err
		}
		out.WriteString(v.Value.(string))
	}
	return NewValue(types.STRING, out.String()), nil
}

func (i *Interpreter) evalTime(c *ir.TimeNode) *Value {
	now := time.Now()
	var out int64
	switch c.Mode {
	case ir.TimeModeSeconds:
		out = now.Unix()

	case ir.TimeModeMicro:
		out = now.UnixMicro()

	case ir.TimeModeMilli:
		out = now.UnixMilli()

	case ir.TimeModeNano:
		out = now.UnixNano()
	}
	return NewValue(types.INT, int(out))
}
