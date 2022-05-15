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
			b.Error(ErrorLevelError, call.Pos(), "invalid function definition")
			continue
		}

		// Get name
		nm := body[0]
		nameV, ok := nm.(*parser.IdentNode)
		if !ok {
			b.Error(ErrorLevelError, nm.Pos(), "expected function name")
			continue
		}
		name := nameV.Value
		body = body[1:]

		// Check if exists
		_, exists := b.Funcs[name]
		if exists {
			b.Error(ErrorLevelError, nm.Pos(), "cannot redefine function %s", name)
			return nil
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
				b.Error(ErrorLevelError, call.Pos(), "expected param")
				continue
			}

			// Get name
			hasErr := false
			nm := call.Args[0]
			nameV, ok := nm.(*parser.IdentNode)
			if !ok {
				b.Error(ErrorLevelError, nm.Pos(), "expected param name")
				hasErr = true
			}
			name := nameV.Value

			// Get type
			typVal := call.Args[1]
			typV, ok := typVal.(*parser.IdentNode)
			if !ok {
				b.Error(ErrorLevelError, typVal.Pos(), "expected param type")
				continue
			}
			typ, err := types.ParseType(typV.Value, b.typeNames)
			if err != nil {
				b.Error(ErrorLevelError, typV.Pos(), "%s", err.Error())
			}

			// Add param
			if hasErr {
				continue
			}
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
					b.Error(ErrorLevelError, retTypVal.Pos(), "expected return type")
					continue
				}
				retTypV, ok := retTypVal.Args[0].(*parser.IdentNode)
				if !ok {
					b.Error(ErrorLevelError, retTypVal.Args[0].Pos(), "expected return type")
					continue
				}
				var err error
				retType, err = types.ParseType(retTypV.Value, b.typeNames)
				if err != nil {
					b.Error(ErrorLevelError, retTypV.Pos(), "%s", err.Error())
					continue
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
			pos:     call.Pos(),
		}
		b.Funcs[name] = fn
	}

	return nil
}

type FnCallNode struct {
	Fn     Node
	Params []Node

	typ types.Type
	pos *tokens.Pos
}

func (c *FnCallNode) Type() types.Type { return c.typ }
func (c *FnCallNode) Pos() *tokens.Pos { return c.pos }
func (c *FnCallNode) Args() []Node     { return c.Params }

func NewFnCallNode(fn Node, args []Node, typ types.Type, pos *tokens.Pos) *FnCallNode {
	return &FnCallNode{
		Fn:     fn,
		Params: args,
		typ:    typ,
		pos:    pos,
	}
}

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
	hasE := b.MatchTypes(n.Pos(), args, expected)
	if hasE {
		return NewTypedNode(fn.RetType, n.Pos()), nil
	}

	return &FnCallNode{
		Fn: &CallNode{
			Call: &FnNode{
				Name: n.Name,
				pos:  n.Pos(),
				typ:  types.NewFuncType(expected, fn.RetType),
			},
			pos: n.Pos(),
		},
		Params: args,
		typ:    fn.RetType,
		pos:    n.Pos(),
	}, nil
}

func (b *Builder) buildFnDef(n *parser.CallNode) error {
	if b.Scope.CurrType() != ScopeTypeGlobal {
		b.Error(ErrorLevelError, n.Pos(), "functions can only be defined in global scope")
		return nil
	}

	// Get name
	if len(n.Args) == 0 {
		b.Error(ErrorLevelError, n.Pos(), "invalid function definition")
		return nil
	}
	_, ok := n.Args[len(n.Args)-1].(*parser.IdentNode)
	if !ok {
		b.Error(ErrorLevelError, n.Pos(), "invalid function definition")
		return nil
	}
	name := n.Args[len(n.Args)-1].(*parser.IdentNode).Value
	n.Args = n.Args[:len(n.Args)-1]
	fn, exists := b.Funcs[name]
	if !exists {
		b.Error(ErrorLevelError, n.Pos(), "invalid function definition")
		return nil
	}

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
			b.Error(ErrorLevelError, n.Pos(), "expected return statement")
			return nil
		}
		call, ok := body[len(body)-1].(*CallNode)
		if !ok {
			b.Error(ErrorLevelError, n.Pos(), "expected return statement")
			return nil
		}
		_, ok = call.Call.(*ReturnNode)
		if !ok {
			b.Error(ErrorLevelError, n.Pos(), "expected return statement")
			return nil
		}
	}

	// Cleanup
	s := b.Scope.CurrScopeInfo()
	fn.Scope = s
	b.Scope.Pop()
	return nil
}

type ReturnNode struct {
	NullCall
	Value Node
}

func (r *ReturnNode) Args() []Node { return []Node{r.Value} }

type FnNode struct {
	Name string
	typ  types.Type
	pos  *tokens.Pos
}

func (f *FnNode) Type() types.Type { return f.typ }
func (f *FnNode) Args() []Node     { return []Node{NewConst(types.IDENT, f.pos, f.Name)} }

func init() {
	nodeBuilders["RETURN"] = nodeBuilder{
		ArgTypes: []types.Type{types.ALL},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			if !b.Scope.HasType(ScopeTypeFunction) {
				b.Error(ErrorLevelError, pos, "RETURN statement outside of function")
				return &ReturnNode{Value: args[0]}, nil
			}

			retTyp := b.Funcs[b.currFn].RetType
			if !args[0].Type().Equal(retTyp) {
				b.Error(ErrorLevelError, args[0].Pos(), "expected return type %s, got %s", retTyp.String(), args[0].Type().String())
			}

			return &ReturnNode{
				Value: args[0],
			}, nil
		},
	}

	nodeBuilders["FN"] = nodeBuilder{
		ArgTypes: []types.Type{types.IDENT},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			name := args[0].(*Const).Value.(string)
			fn, exists := b.Funcs[name]
			if !exists {
				b.Error(ErrorLevelError, pos, "undefined function: %s", name)
				return NewTypedValue(types.INVALID), nil
			}
			pars := make([]types.Type, len(fn.Params))
			for i, par := range fn.Params {
				pars[i] = par.Type
			}

			return &FnNode{
				Name: name,
				typ:  types.NewFuncType(pars, fn.RetType),
				pos:  args[0].Pos(),
			}, nil
		},
	}

	nodeBuilders["CALL"] = nodeBuilder{
		ArgTypes: []types.Type{types.FUNCTION, types.ALL, types.VARIADIC},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			// Match types
			typ, ok := args[0].Type().(*types.FuncType)
			if !ok {
				return NewTypedValue(types.INVALID), nil
			}
			if len(typ.ParTypes) > 0 {
				err := b.MatchTypes(pos, args[1:], typ.ParTypes)
				if err {
					return NewTypedValue(types.INVALID), nil
				}
			} else {
				args = args[:1]
			}

			return &FnCallNode{
				Fn:     args[0],
				Params: args[1:],
				typ:    typ.RetType,
				pos:    pos,
			}, nil
		},
	}
}
