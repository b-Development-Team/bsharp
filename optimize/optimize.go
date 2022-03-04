package optimize

import "github.com/Nv7-Github/bsharp/ir"

type VariableInfo struct {
	CurrValue *Result // may be nil
}

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
	s := Scope{scopes: make([]scope, 0)}
	s.Push()
	return &Optimizer{
		ir:    i,
		scope: s,
	}
}

func (o *Optimizer) Optimize() *ir.IR {
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
