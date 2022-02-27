package interpreter

import "github.com/Nv7-Github/bsharp/ir"

func (i *Interpreter) evalPrint(c *ir.PrintNode) error {
	v, err := i.evalNode(c.Arg)
	if err != nil {
		return err
	}
	_, err = i.stdout.Write([]byte(v.Value.(string) + "\n"))
	return err
}
