package interpreter

import (
	"strings"

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
