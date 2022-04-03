package main

import (
	"log"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
)

func GetScope(i *ir.IR, pos *tokens.Pos) *ir.ScopeInfo {
	for _, fn := range i.Funcs {
		log.Println(fn.Name, fn.Pos(), pos, fn.Pos().EndLine, fn.Pos().EndChar, fn.Pos().Contains(pos))
		if fn.Pos().Contains(pos) {
			for _, node := range fn.Body {
				_, ok := node.(*ir.BlockNode)
				if ok && node.Pos().Contains(pos) {
					s := scopeRecurse(node.(*ir.BlockNode), pos)
					if s != nil {
						return s
					}
				}
			}
			return fn.Scope
		}
	}

	for _, node := range i.Body {
		_, ok := node.(*ir.BlockNode)
		if ok && node.Pos().Contains(pos) {
			s := scopeRecurse(node.(*ir.BlockNode), pos)
			if s != nil {
				return s
			}
		}
	}
	return i.GlobalScope
}

func scopeRecurse(bl *ir.BlockNode, pos *tokens.Pos) *ir.ScopeInfo {
	switch b := bl.Block.(type) {
	case ir.BodyBlock:
		for _, node := range b.Block() {
			_, ok := node.(*ir.BlockNode)
			if ok && node.Pos().Contains(pos) {
				s := scopeRecurse(node.(*ir.BlockNode), pos)
				if s != nil {
					return s
				}
			}
		}
		return b.ScopeInfo()

	case *ir.IfNode:
		if b.Else[0].Pos().Extend(b.Else[len(b.Else)-1].Pos()).Contains(pos) { // In else
			for _, node := range b.Else {
				_, ok := node.(*ir.BlockNode)
				if ok && node.Pos().Contains(pos) {
					s := scopeRecurse(node.(*ir.BlockNode), pos)
					if s != nil {
						return s
					}
				}
			}
			return b.ElseScope
		}

		for _, node := range b.Body {
			_, ok := node.(*ir.BlockNode)
			if ok && node.Pos().Contains(pos) {
				s := scopeRecurse(node.(*ir.BlockNode), pos)
				if s != nil {
					return s
				}
			}
		}
		return b.Scope

	case *ir.SwitchNode:
		for _, cs := range b.Cases {
			if cs.Pos().Contains(pos) {
				return scopeRecurse(cs, pos)
			}
		}

		if b.Default != nil && b.Default.Pos().Contains(pos) {
			return scopeRecurse(b.Default, pos)
		}
	}

	return nil
}
