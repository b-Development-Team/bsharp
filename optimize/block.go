package optimize

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
)

func (o *Optimizer) optimizeWhile(n *ir.WhileNode, pos *tokens.Pos) *Result {
	canUnroll := true
	notDead := false
	body := make([]ir.Node, 0)
	o.scope.Push()
	for {
		cond := o.OptimizeNode(n.Condition)
		if !cond.IsConst {
			canUnroll = false
			break
		}

		if len(body) > 1000 { // Too much!
			canUnroll = false
			break
		}

		// Unroll
		if cond.Stmt.(*ir.Const).Value.(bool) {
			for _, node := range n.Body {
				res := o.OptimizeNode(node)
				if res.Stmt != nil && res.NotDead {
					body = append(body, res.Stmt)
					notDead = true
				}
			}
		} else {
			break
		}
	}
	o.scope.Pop()
	if canUnroll {
		return &Result{
			Stmt: ir.NewCallNode(&ir.BlockNode{
				Body: body,
			}, pos),
			NotDead: notDead,
		}
	}

	cond := o.OptimizeNode(n.Condition)

	// Can't unroll, instead analyze variable scopes
	o.scope.Push()
	body = make([]ir.Node, 0, len(n.Body))
	for _, node := range n.Body {
		res := o.OptimizeNode(node)
		if res.Stmt != nil && res.NotDead {
			body = append(body, res.Stmt)
		}
	}
	o.scope.Pop()

	return &Result{
		Stmt: ir.NewCallNode(&ir.WhileNode{
			Condition: cond.Stmt,
			Body:      body,
		}, pos),
		NotDead: true,
	}
}

func (o *Optimizer) optimizeIf(i *ir.IfNode, pos *tokens.Pos) *Result {
	cond := o.OptimizeNode(i.Condition)
	if cond.IsConst {
		var out []ir.Node
		if cond.Stmt.(*ir.Const).Value.(bool) {
			out = make([]ir.Node, 0, len(i.Body))
			o.scope.Push()
			for _, node := range i.Body {
				v := o.OptimizeNode(node)
				if v.Stmt != nil && v.NotDead {
					out = append(out, v.Stmt)
				}
			}
			o.scope.Pop()
		} else {
			out = make([]ir.Node, 0, len(i.Else))
			o.scope.Push()
			for _, node := range i.Else {
				v := o.OptimizeNode(node)
				if v.Stmt != nil && v.NotDead {
					out = append(out, v.Stmt)
				}
			}
			o.scope.Pop()
		}

		return &Result{
			Stmt: ir.NewCallNode(&ir.BlockNode{
				Body: out,
			}, pos),
			NotDead: len(out) > 0,
		}
	}

	// Build
	body := make([]ir.Node, 0, len(i.Body))
	o.scope.Push()
	for _, node := range i.Body {
		r := o.OptimizeNode(node)
		if r.Stmt != nil && r.NotDead {
			body = append(body, r.Stmt)
		}
	}
	o.scope.Pop()

	var els []ir.Node
	if i.Else != nil {
		els = make([]ir.Node, 0, len(i.Else))
		o.scope.Push()
		for _, node := range i.Else {
			r := o.OptimizeNode(node)
			if r.Stmt != nil && r.NotDead {
				els = append(els, r.Stmt)
			}
		}
		o.scope.Pop()
	}

	return &Result{
		Stmt: ir.NewCallNode(&ir.IfNode{
			Condition: cond.Stmt,
			Body:      body,
			Else:      els,
		}, pos),
		NotDead: true,
	}
}

func (o *Optimizer) optimizeCase(i *ir.Case) []ir.Node {
	out := make([]ir.Node, 0, len(i.Body))
	o.scope.Push()
	for _, node := range i.Body {
		r := o.OptimizeNode(node)
		if r.Stmt != nil && r.NotDead {
			out = append(out, r.Stmt)
		}
	}
	o.scope.Pop()
	return out
}

func (o *Optimizer) optimizeSwitch(i *ir.SwitchNode, pos *tokens.Pos) *Result {
	v := o.OptimizeNode(i.Value)

	// Cases
	cases := make([]*ir.Case, len(i.Cases))
	for i, cse := range i.Cases {
		cases[i] = &ir.Case{
			Value: cse.Value,
			Body:  o.optimizeCase(cse),
		}
	}

	// Default
	var def []ir.Node
	if i.Default != nil {
		def = make([]ir.Node, 0, len(i.Default))
		o.scope.Push()
		for _, node := range i.Default {
			r := o.OptimizeNode(node)
			if r.Stmt != nil && r.NotDead {
				def = append(def, r.Stmt)
			}
		}
		o.scope.Pop()
	}

	return &Result{
		Stmt: ir.NewCallNode(&ir.SwitchNode{
			Value:   v.Stmt,
			Cases:   cases,
			Default: def,
		}, pos),
		NotDead: true,
	}
}
