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
