package interpreter

import (
	"fmt"
	"reflect"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

type Extension struct {
	Name     string
	Fn       reflect.Value // reflect.FuncOf
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

func (e *Extension) Call(vals []interface{}) (interface{}, error) {
	parVals := make([]reflect.Value, len(vals))
	for i, val := range vals {
		parVals[i] = reflect.ValueOf(val)
	}

	retVals := e.Fn.Call(parVals)
	if len(retVals) == 0 {
		return nil, nil
	}

	if !types.NULL.Equal(e.RetType) {
		return retVals[0].Interface(), nil
	}

	return nil, nil
}

func NewExtension(fnObj interface{}, name string) (*Extension, error) {
	fn := reflect.ValueOf(fnObj)

	// Get par types
	fnTyp := fn.Type()
	parCnt := fnTyp.NumIn()
	parTypes := make([]types.Type, parCnt)
	var err error

	for i := 0; i < parCnt; i++ {
		parTypes[i], err = getType(fnTyp.In(i))
		if err != nil {
			return nil, err
		}
	}

	// Get ret type
	numOut := fnTyp.NumOut()
	hasRet := numOut > 0
	var retType types.Type = types.NULL

	if hasRet && numOut > 1 {
		return nil, fmt.Errorf("extension cannot return multiple values")
	} else if hasRet {
		retType, err = getType(fnTyp.Out(0))
		if err != nil {
			return nil, err
		}
	}

	return &Extension{
		Fn:       fn,
		Name:     name,
		ParTypes: parTypes,
		RetType:  retType,
	}, nil
}

var typMap = map[reflect.Kind]types.BasicType{
	reflect.Int:     types.INT,
	reflect.Float64: types.FLOAT,
	reflect.String:  types.STRING,
}

func getType(typ reflect.Type) (types.Type, error) {
	basicTyp, exists := typMap[typ.Kind()]
	if exists {
		return basicTyp, nil
	}

	if typ.Kind() == reflect.Slice {
		valTyp, err := getType(typ.Elem())
		if err != nil {
			return nil, err
		}
		return types.NewArrayType(valTyp), nil
	}

	return nil, fmt.Errorf("unknown type \"%s\"", typ.String())
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
	res, err := ext.Call(args)
	if err != nil {
		return nil, err
	}

	return NewValue(ext.RetType, res), nil
}
