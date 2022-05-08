package ir

import (
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type BodyBlock interface {
	Block

	Block() []Node
	ScopeInfo() *ScopeInfo
}

type IfNode struct {
	Condition Node
	Body      []Node
	Scope     *ScopeInfo

	Else      []Node // nil if no else
	ElseScope *ScopeInfo
}

type WhileNode struct {
	Condition Node
	Body      []Node
	Scope     *ScopeInfo
}

func (w *WhileNode) Block() []Node {
	return w.Body
}

func (w *WhileNode) ScopeInfo() *ScopeInfo {
	return w.Scope
}

type Case struct {
	Value *Const
	Body  []Node
	Scope *ScopeInfo
}

func (c *Case) Block() []Node {
	return c.Body
}

func (c *Case) ScopeInfo() *ScopeInfo {
	return c.Scope
}

type Default struct {
	Body  []Node
	Scope *ScopeInfo
}

func (d *Default) Block() []Node {
	return d.Body
}

func (d *Default) ScopeInfo() *ScopeInfo {
	return d.Scope
}

type SwitchNode struct {
	Value   Node
	Cases   []*BlockNode // *Case is Block
	Default *BlockNode   // if nil, no default
}

func init() {
	blockBuilders["IF"] = blockBuilder{
		Build: func(b *Builder, pos *tokens.Pos, args []parser.Node) (Block, error) {
			if len(args) < 2 {
				b.Error(ErrorLevelError, pos, "IF requires at least 2 arguments")
				return &IfNode{Condition: NewTypedNode(types.INVALID, pos)}, nil
			}

			cond, err := b.buildNode(args[0])
			if err != nil {
				return nil, err
			}
			if !types.BOOL.Equal(cond.Type()) {
				b.Error(ErrorLevelError, cond.Pos(), "expected boolean for condition in IF statement")
				return &IfNode{Condition: NewTypedNode(types.INVALID, pos)}, nil
			}

			body := make([]Node, 0, len(args)-1)
			var els []parser.Node
			b.Scope.Push(ScopeTypeIf)
			for i, arg := range args[1:] {
				// ELSE?
				_, ok := arg.(*parser.IdentNode)
				if ok && arg.(*parser.IdentNode).Value == "ELSE" {
					els = args[i+2:]
					break
				}

				node, err := b.buildNode(arg)
				if err != nil {
					return nil, err
				}
				body = append(body, node)
			}
			scope := b.Scope.CurrScopeInfo()
			b.Scope.Pop()

			if els != nil {
				elsBody := make([]Node, 0, len(els))
				b.Scope.Push(ScopeTypeIf)
				for _, v := range els {
					node, err := b.buildNode(v)
					if err != nil {
						return nil, err
					}
					elsBody = append(elsBody, node)
				}
				elsscope := b.Scope.CurrScopeInfo()
				b.Scope.Pop()
				return &IfNode{
					Condition: cond,
					Body:      body,
					Scope:     scope,
					Else:      elsBody,
					ElseScope: elsscope,
				}, nil
			}

			return &IfNode{
				Condition: cond,
				Body:      body,
				Scope:     scope,
			}, nil
		},
	}

	blockBuilders["WHILE"] = blockBuilder{
		Build: func(b *Builder, pos *tokens.Pos, args []parser.Node) (Block, error) {
			if len(args) < 2 {
				b.Error(ErrorLevelError, pos, "WHILE requires at least 2 arguments")
				return &WhileNode{Condition: NewTypedNode(types.INVALID, pos)}, nil
			}
			cond, err := b.buildNode(args[0])
			if err != nil {
				return nil, err
			}

			body := make([]Node, len(args)-1)
			b.Scope.Push(ScopeTypeWhile)
			for i, v := range args[1:] {
				body[i], err = b.buildNode(v)
				if err != nil {
					return nil, err
				}
			}
			b.Scope.Pop()

			return &WhileNode{
				Condition: cond,
				Body:      body,
			}, nil
		},
	}

	blockBuilders["CASE"] = blockBuilder{
		Build: func(b *Builder, pos *tokens.Pos, args []parser.Node) (Block, error) {
			if b.Scope.CurrType() != ScopeTypeSwitch {
				b.Error(ErrorLevelError, pos, "CASE can only be used inside SWITCH")
				return NewTypedValue(types.INVALID), nil
			}
			if len(args) < 2 {
				b.Error(ErrorLevelError, pos, "CASE requires at least 2 arguments")
				return &Case{}, nil
			}

			val, err := b.buildNode(args[0])
			if err != nil {
				return nil, err
			}
			cnst, ok := val.(*Const)
			if !ok {
				b.Error(ErrorLevelError, val.Pos(), "expected constant for CASE value")
				return &Case{}, nil
			}
			if !hashable.Equal(cnst.Type()) {
				b.Error(ErrorLevelError, val.Pos(), "expected hashable type for CASE value")
				return &Case{}, nil
			}

			body := make([]Node, len(args)-1)
			b.Scope.Push(ScopeTypeCase)
			for i, v := range args[1:] {
				body[i], err = b.buildNode(v)
				if err != nil {
					return nil, err
				}
			}
			b.Scope.Pop()
			return &Case{
				Value: cnst,
				Body:  body,
			}, nil
		},
	}

	blockBuilders["DEFAULT"] = blockBuilder{
		Build: func(b *Builder, pos *tokens.Pos, args []parser.Node) (Block, error) {
			if b.Scope.CurrType() != ScopeTypeSwitch {
				b.Error(ErrorLevelError, pos, "DEFAULT can only be used inside SWITCH")
				return NewTypedValue(types.INVALID), nil
			}
			if len(args) < 1 {
				b.Error(ErrorLevelError, pos, "DEFAULT requires at least 1 argument")
				return &Default{}, nil
			}

			body := make([]Node, len(args))
			b.Scope.Push(ScopeTypeCase)
			var err error
			for i, v := range args {
				body[i], err = b.buildNode(v)
				if err != nil {
					return nil, err
				}
			}
			b.Scope.Pop()
			return &Default{
				Body: body,
			}, nil
		},
	}

	blockBuilders["SWITCH"] = blockBuilder{
		Build: func(b *Builder, pos *tokens.Pos, args []parser.Node) (Block, error) {
			if len(args) < 2 {
				b.Error(ErrorLevelError, pos, "SWITCH requires at least 2 arguments")
				return &SwitchNode{Value: NewTypedNode(types.INVALID, pos)}, nil
			}

			val, err := b.buildNode(args[0])
			if err != nil {
				return nil, err
			}
			if !hashable.Equal(val.Type()) {
				b.Error(ErrorLevelError, val.Pos(), "expected hashable type for SWITCH value")
			}

			// Get cases
			cases := make([]*BlockNode, 0, len(args)-1)
			var def *BlockNode
			b.Scope.Push(ScopeTypeSwitch)
			for _, v := range args[1:] {
				node, err := b.buildNode(v)
				if err != nil {
					return nil, err
				}
				blk, ok := node.(*BlockNode)
				if !ok {
					b.Error(ErrorLevelError, v.Pos(), "expected case")
					continue
				}
				cs, ok := blk.Block.(*Case)
				if cs.Value == nil {
					continue
				}
				if !ok {
					// Default
					_, ok := blk.Block.(*Default)
					if ok {
						if def != nil {
							b.Error(ErrorLevelError, node.Pos(), "only one default case allowed")
						}
						def = blk
						continue
					}
					b.Error(ErrorLevelError, v.Pos(), "expected case")
					continue
				}
				if !cs.Value.Type().Equal(val.Type()) {
					b.Error(ErrorLevelError, v.Pos(), "expected case with type %s", val.Type())
					continue
				}
				cases = append(cases, blk)
			}
			b.Scope.Pop()
			return &SwitchNode{
				Value:   val,
				Cases:   cases,
				Default: def,
			}, nil
		},
	}
}
