package optimize

import "github.com/Nv7-Github/bsharp/ir"

type Optimizer struct {
	ir    *ir.IR
	scope Scope
}

type Result struct {
	Stmt    ir.Node
	IsConst bool
	NotDead bool
}

func NewOptimizer(i *ir.IR) *Optimizer {
	s := Scope{Scopes: make([]scope, 0)}
	s.Push()
	return &Optimizer{
		ir:    i,
		scope: s,
	}
}

func (o *Optimizer) Optimize() *ir.IR {
	for k := range o.ir.Funcs {
		o.optimizeFunc(k)
	}

	out := make([]ir.Node, 0, len(o.ir.Body))
	for _, node := range o.ir.Body {
		res := o.OptimizeNode(node)
		if res.Stmt != nil && res.NotDead {
			out = append(out, res.Stmt)
		}
	}

	o.ir.Body = out
	return o.ir
}
