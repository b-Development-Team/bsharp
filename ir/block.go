package ir

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

func bodyCode(cnf CodeConfig, body []Node) string {
	tab := strings.Repeat(" ", cnf.Indent)
	out := &strings.Builder{}
	for _, node := range body {
		code := node.Code(cnf)
		lines := strings.Split(code, "\n")
		for _, line := range lines {
			out.WriteString(tab)
			out.WriteString(line)
			out.WriteString("\n")
		}
	}
	return out.String()
}

type IfNode struct {
	NullCall
	Condition Node
	Body      []Node
	Else      []Node // nil if no else
}

func (i *IfNode) Code(cnf CodeConfig) string {
	if i.Else != nil {
		return fmt.Sprintf("[IF %s\n%sELSE\n%s]", i.Condition.Code(cnf), bodyCode(cnf, i.Body), bodyCode(cnf, i.Else))
	}
	return fmt.Sprintf("[IF %s\n%s]", i.Condition.Code(cnf), bodyCode(cnf, i.Body))
}

type WhileNode struct {
	NullCall
	Condition Node
	Body      []Node
}

func (w *WhileNode) Code(cnf CodeConfig) string {
	return fmt.Sprintf("[WHILE %s\n%s]", w.Condition.Code(cnf), bodyCode(cnf, w.Body))
}

type Case struct {
	NullCall
	Value *Const
	Body  []Node
}

func (c *Case) Code(cnf CodeConfig) string {
	return fmt.Sprintf("[CASE %s\n%s]", c.Value.Code(cnf), bodyCode(cnf, c.Body))
}

type Default struct {
	NullCall
	Body []Node
}

func (d *Default) Code(cnf CodeConfig) string {
	return fmt.Sprintf("[DEFAULT\n%s]", bodyCode(cnf, d.Body))
}

type SwitchNode struct {
	NullCall
	Value   Node
	Cases   []*Case
	Default []Node // if nil, no default
}

type BlockNode struct {
	NullCall
	Body []Node
}

func (b *BlockNode) Code(cnf CodeConfig) string {
	return fmt.Sprintf("[BLOCK\n%s]", bodyCode(cnf, b.Body))
}

func (s *SwitchNode) Code(cnf CodeConfig) string {
	args := &strings.Builder{}
	tab := strings.Repeat(" ", cnf.Indent)
	for i, c := range s.Cases {
		lines := strings.Split(c.Code(cnf), "\n")
		for _, line := range lines {
			args.WriteString(tab)
			args.WriteString(line)
			args.WriteString("\n")
		}
		if i != len(s.Cases)-1 {
			args.WriteString("\n")
		}
	}
	if s.Default != nil {
		lines := strings.Split(bodyCode(cnf, s.Default), "\n")
		for _, line := range lines {
			args.WriteString(tab)
			args.WriteString(line)
			args.WriteString("\n")
		}
	}
	return fmt.Sprintf("[SWITCH\n%s]", args.String())
}

func init() {
	blockBuilders["IF"] = blockBuilder{
		Build: func(b *Builder, pos *tokens.Pos, args []parser.Node) (Call, error) {
			if len(args) < 2 {
				return nil, pos.Error("IF requires at least 2 arguments")
			}

			cond, err := b.buildNode(args[0])
			if err != nil {
				return nil, err
			}
			if !types.BOOL.Equal(cond.Type()) {
				return nil, cond.Pos().Error("expected boolean for condition in IF statement")
			}

			body := make([]Node, 0, len(args)-1)
			hasElse := false
			var els []Node
			b.Scope.Push(ScopeTypeIf)
			for _, arg := range args[1:] {
				// ELSE?
				_, ok := arg.(*parser.IdentNode)
				if ok && arg.(*parser.IdentNode).Value == "ELSE" {
					if hasElse {
						return nil, arg.Pos().Error("ELSE can only be used once in IF statement")
					}
					hasElse = true
					els = make([]Node, 0, len(args)-1)
					continue
				}

				node, err := b.buildNode(arg)
				if err != nil {
					return nil, err
				}
				if !hasElse {
					body = append(body, node)
				} else {
					els = append(els, node)
				}
			}
			b.Scope.Pop()

			return &IfNode{
				Condition: cond,
				Body:      body,
				Else:      els,
			}, nil
		},
	}

	blockBuilders["WHILE"] = blockBuilder{
		Build: func(b *Builder, pos *tokens.Pos, args []parser.Node) (Call, error) {
			if len(args) < 2 {
				return nil, pos.Error("WHILE requires at least 2 arguments")
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
		Build: func(b *Builder, pos *tokens.Pos, args []parser.Node) (Call, error) {
			if b.Scope.CurrType() != ScopeTypeSwitch {
				return nil, pos.Error("CASE can only be used inside SWITCH")
			}
			if len(args) < 2 {
				return nil, pos.Error("CASE requires at least 2 arguments")
			}

			val, err := b.buildNode(args[0])
			if err != nil {
				return nil, err
			}
			cnst, ok := val.(*Const)
			if !ok {
				return nil, val.Pos().Error("expected constant for CASE value")
			}
			if !hashable.Equal(cnst.Type()) {
				return nil, val.Pos().Error("expected hashable type for CASE value")
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
		Build: func(b *Builder, pos *tokens.Pos, args []parser.Node) (Call, error) {
			if b.Scope.CurrType() != ScopeTypeSwitch {
				return nil, pos.Error("DEFAULT can only be used inside SWITCH")
			}
			if len(args) < 1 {
				return nil, pos.Error("DEFAULT requires at least 1 argument")
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
		Build: func(b *Builder, pos *tokens.Pos, args []parser.Node) (Call, error) {
			if len(args) < 2 {
				return nil, pos.Error("SWITCH requires at least 2 arguments")
			}

			val, err := b.buildNode(args[0])
			if err != nil {
				return nil, err
			}
			if !hashable.Equal(val.Type()) {
				return nil, val.Pos().Error("expected hashable type for SWITCH value")
			}

			// Get cases
			cases := make([]*Case, 0, len(args)-1)
			var def []Node
			b.Scope.Push(ScopeTypeSwitch)
			for _, v := range args[1:] {
				node, err := b.buildNode(v)
				if err != nil {
					return nil, err
				}
				call, ok := node.(*CallNode)
				if !ok {
					return nil, v.Pos().Error("expected case")
				}
				cs, ok := call.Call.(*Case)
				if !ok {
					// Default
					defaul, ok := call.Call.(*Default)
					if ok {
						if def != nil {
							return nil, node.Pos().Error("only one default case allowed")
						}
						def = defaul.Body
						continue
					}
					return nil, v.Pos().Error("expected case")
				}
				if !cs.Value.Type().Equal(val.Type()) {
					return nil, v.Pos().Error("expected case with type %s", val.Type())
				}
				cases = append(cases, cs)
			}
			b.Scope.Pop()
			return &SwitchNode{
				Value:   val,
				Cases:   cases,
				Default: def,
			}, nil
		},
	}

	blockBuilders["BLOCK"] = blockBuilder{
		Build: func(b *Builder, pos *tokens.Pos, args []parser.Node) (Call, error) {
			if len(args) < 1 {
				return nil, pos.Error("BLOCK requires at least 1 argument")
			}
			body := make([]Node, len(args)-1)
			var err error
			b.Scope.Push(ScopeTypeBlock)
			for i, v := range args[1:] {
				body[i], err = b.buildNode(v)
				if err != nil {
					return nil, err
				}
			}
			b.Scope.Pop()

			return &BlockNode{
				Body: body,
			}, nil
		},
	}
}
