package ir

import (
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
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

		// Check if exists
		_, exists := b.Funcs[name]
		if exists {
			return call.Pos().Error("function %s already exists!", name)
		}

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
				Pos:  call.Pos(),
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
				body = body[1:]
			}
		}

		// Add name to end of body
		body = append(body, &parser.IdentNode{
			Value: name,
		})

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

type FnCallNode struct {
	Name string
	Args []Node

	typ types.Type
	pos *tokens.Pos
}

func (c *FnCallNode) Type() types.Type { return c.typ }
func (c *FnCallNode) Pos() *tokens.Pos { return c.pos }

func (b *Builder) buildFnCall(n *parser.CallNode) (Node, error) {
	fn := b.Funcs[n.Name] // We know this exists because this function won't be called if the function doesn't exist

	args := make([]Node, len(n.Args))
	for i, arg := range n.Args {
		node, err := b.buildNode(arg)
		if err != nil {
			return nil, err
		}
		args[i] = node
	}

	expected := make([]types.Type, len(fn.Params))
	for i, par := range fn.Params {
		expected[i] = par.Type
	}
	err := MatchTypes(n.Pos(), args, expected)
	if err != nil {
		return nil, err
	}

	return &FnCallNode{
		Name: n.Name,
		Args: args,
		typ:  fn.RetType,
		pos:  n.Pos(),
	}, nil
}

func (b *Builder) buildFnDef(n *parser.CallNode) error {
	if b.Scope.CurrType() != ScopeTypeGlobal {
		return n.Pos().Error("functions can only be defined in global scope")
	}

	// Get name
	name := n.Args[len(n.Args)-1].(*parser.IdentNode).Value
	n.Args = n.Args[:len(n.Args)-1]
	fn := b.Funcs[name]

	// Add params
	b.Scope.Push(ScopeTypeFunction)
	b.currFn = fn.Name
	for _, par := range fn.Params {
		par.ID = b.Scope.AddVariable(par.Name, par.Type, par.Pos)
	}

	// Build body
	body := make([]Node, len(n.Args))
	for i, arg := range n.Args {
		node, err := b.buildNode(arg)
		if err != nil {
			return err
		}
		body[i] = node
	}
	fn.Body = body

	// Check if return
	if !types.NULL.Equal(fn.RetType) {
		if len(body) == 0 {
			return n.Pos().Error("expected return statement")
		}
		call, ok := body[len(body)-1].(*CallNode)
		if !ok {
			return n.Pos().Error("expected return statement")
		}
		_, ok = call.Call.(*ReturnNode)
		if !ok {
			return n.Pos().Error("expected return statement")
		}
	}

	// Cleanup
	b.Scope.Pop()
	return nil
}

type ReturnNode struct {
	NullCall
	Value Node
}

func init() {
	nodeBuilders["RETURN"] = nodeBuilder{
		ArgTypes: []types.Type{types.ANY},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			if !b.Scope.HasType(ScopeTypeFunction) {
				return nil, pos.Error("return statement outside of function")
			}
			retTyp := b.Funcs[b.currFn].RetType
			if !args[0].Type().Equal(retTyp) {
				return nil, args[0].Pos().Error("expected return type %s, got %s", retTyp.String(), args[0].Type().String())
			}
			return &ReturnNode{
				Value: args[0],
			}, nil
		},
	}
}