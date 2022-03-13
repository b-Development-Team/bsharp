package optimize

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
)

func (o *Optimizer) optimizeFunc(name string) {
	fn := o.ir.Funcs[name]
	body := make([]ir.Node, 0, len(fn.Body))
	o.scope.Push()
	for _, n := range fn.Body {
		r := o.OptimizeNode(n)
		if r.Stmt != nil && r.NotDead {
			body = append(body, r.Stmt)
		}
	}
	o.scope.Pop()

	fn.Body = body
}

func (o *Optimizer) optimizeRet(n *ir.ReturnNode, pos *tokens.Pos) *Result {
	v := o.OptimizeNode(n.Value)
	return &Result{
		Stmt: ir.NewCallNode(
			&ir.ReturnNode{Value: v.Stmt},
			pos,
		),
		NotDead: true,
	}
}

func (o *Optimizer) optimizeCall(n *ir.FnCallNode, pos *tokens.Pos) *Result {
	fn := o.OptimizeNode(n.Fn)
	pars := make([]ir.Node, len(n.Args))
	for i, arg := range n.Args {
		pars[i] = o.OptimizeNode(arg).Stmt
	}

	return &Result{
		Stmt:    ir.NewFnCallNode(fn.Stmt, pars, n.Type(), pos),
		NotDead: true,
	}
}

func (o *Optimizer) optimizeFn(n *ir.FnNode, pos *tokens.Pos) *Result {
	return &Result{
		Stmt:    ir.NewCallNode(n, pos),
		IsConst: true,
	}
}
