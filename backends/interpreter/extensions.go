package interpreter

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

type Extension struct {
	Name     string
	Fn       func([]interface{}) (interface{}, error) // reflect.FuncOf
	ParTypes []types.Type
	RetType  types.Type
}

func (e *Extension) IRExtension() *ir.Extension {
	return &ir.Extension{
		Name:    e.Name,
		Params:  e.ParTypes,
		RetType: e.RetType,
	}
}

func NewExtension(name string, fn func([]interface{}) (interface{}, error), parTypes []types.Type, retTyp types.Type) *Extension {
	return &Extension{
		Fn:       fn,
		Name:     name,
		ParTypes: parTypes,
		RetType:  retTyp,
	}
}
func (i *Interpreter) evalExtensionCall(n *ir.ExtensionCall) (*Value, error) {
	ext := i.extensions[n.Name]
	// Build args
	args := make([]interface{}, len(n.Args))
	for ind, arg := range n.Args {
		val, err := i.evalNode(arg)
		if err != nil {
			return nil, err
		}
		args[ind] = val.Value
	}
	// Call
	res, err := ext.Fn(args)
	if err != nil {
		return nil, err
	}

	return NewValue(ext.RetType, res), nil
}
