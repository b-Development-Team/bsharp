package interpreter

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

func (i *Interpreter) evalCallNode(n *ir.FnCallNode) (*Value, error) {
	fn := i.ir.Funcs[n.Name]
	// Build args
	args, err := i.evalNodes(n.Args)
	if err != nil {
		return nil, err
	}
	for ind, par := range fn.Params {
		i.Variables[par.ID] = args[ind]
	}
	// Run
	for _, v := range fn.Body {
		_, err = i.evalNode(v)
		if err != nil {
			return nil, err
		}
	}
	retVal := i.retVal
	if types.NULL.Equal(fn.RetType) {
		retVal = NewValue(types.NULL, nil)
	}
	i.retVal = nil // Un-return
	return retVal, nil
}

func (i *Interpreter) evalReturnNode(n *ir.ReturnNode) (*Value, error) {
	v, err := i.evalNode(n.Value)
	if err != nil {
		return nil, err
	}
	i.retVal = v
	return nil, nil
}
