package ir

import (
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/types"
)

func (b *Builder) functionPass(p *parser.Parser) error {
	for _, node := range p.Nodes {
		call, ok := node.(*parser.CallNode)
		if !ok {
			continue
		}
		if call.Name != "FUNC" {
			continue
		}

		// Its a function!
		body := call.Args
		if len(body) < 1 {
			return call.Pos().Error("expected body")
		}

		// Get name
		nm := body[0]
		nameV, ok := nm.(*parser.IdentNode)
		if !ok {
			return nm.Pos().Error("expected function name")
		}
		name := nameV.Value
		body = body[1:]

		// Get params
		params := make([]*Param, 0)
		for {
			if len(body) == 0 {
				break
			}

			node := body[0]
			call, ok := node.(*parser.CallNode)
			if !ok {
				break
			}
			if call.Name != "PARAM" {
				break
			}

			// Its a param!
			if len(call.Args) != 2 {
				return call.Pos().Error("expected param")
			}

			// Get name
			nm := call.Args[0]
			nameV, ok := nm.(*parser.IdentNode)
			if !ok {
				return nm.Pos().Error("expected param name")
			}
			name := nameV.Value

			// Get type
			typVal := call.Args[1]
			typV, ok := typVal.(*parser.IdentNode)
			if !ok {
				return typVal.Pos().Error("expected param type")
			}
			typ, err := types.ParseType(typV.Value)
			if err != nil {
				return typV.Pos().Error("%s", err.Error())
			}

			// Add param
			params = append(params, &Param{
				ID:   -1,
				Name: name,
				Type: typ,
			})
			body = body[1:]
		}

		// Return type
		retType := types.Type(types.NULL)
		if len(body) > 0 {
			retTypFn := body[0]
			retTypVal, ok := retTypFn.(*parser.CallNode)
			if ok && retTypVal.Name == "RETURNS" {
				if len(retTypVal.Args) != 1 {
					return retTypVal.Pos().Error("expected return type")
				}
				retTypV, ok := retTypVal.Args[0].(*parser.IdentNode)
				if !ok {
					return retTypV.Pos().Error("expected return type")
				}
				var err error
				retType, err = types.ParseType(retTypV.Value)
				if err != nil {
					return retTypV.Pos().Error("%s", err.Error())
				}
			}
		}

		// Update call
		call.Args = body

		// Save function
		fn := &Function{
			Name:    name,
			Params:  params,
			RetType: retType,
		}
		b.Funcs[name] = fn
	}

	return nil
}
