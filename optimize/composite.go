package optimize

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

func (o *Optimizer) optimizeIndex(n *ir.IndexNode, pos *tokens.Pos) *Result {
	v := o.OptimizeNode(n.Value)
	ind := o.OptimizeNode(n.Index)
	if v.IsConst && ind.IsConst { // Arrays cant be const
		return &Result{
			Stmt:    ir.NewConst(n.Type(), pos, v.Stmt.(*ir.Const).Value.(string)[ind.Stmt.(*ir.Const).Value.(int)]),
			IsConst: true,
		}
	}
	return &Result{
		Stmt:    ir.NewCallNode(ir.NewIndexNode(v.Stmt, ind.Stmt), pos),
		IsConst: false,
	}
}

func (o *Optimizer) optimizeLength(n *ir.LengthNode, pos *tokens.Pos) *Result {
	v := o.OptimizeNode(n.Value)
	if v.IsConst {
		if types.STRING.Equal(v.Stmt.Type()) {
			return &Result{
				Stmt:    ir.NewConst(n.Type(), pos, len(v.Stmt.(*ir.Const).Value.(string))),
				IsConst: true,
			}
		}
		return &Result{
			Stmt:    ir.NewConst(n.Type(), pos, len(*v.Stmt.(*ir.Const).Value.(*[]interface{}))),
			IsConst: true,
		}
	}
	return &Result{
		Stmt: ir.NewCallNode(&ir.LengthNode{
			Value: v.Stmt,
		}, pos),
		IsConst: false,
	}
}

func (o *Optimizer) optimizeMake(n *ir.MakeNode, pos *tokens.Pos) *Result {
	return &Result{
		Stmt:    ir.NewCallNode(n, pos),
		IsConst: false,
	}
}

func (o *Optimizer) optimizeSet(n *ir.SetNode, pos *tokens.Pos) *Result {
	m := o.OptimizeNode(n.Map)
	k := o.OptimizeNode(n.Key)
	v := o.OptimizeNode(n.Value)

	// If there were map literals, then check if const and if so, assign
	return &Result{
		Stmt: ir.NewCallNode(&ir.SetNode{
			Map:   m.Stmt,
			Key:   k.Stmt,
			Value: v.Stmt,
		}, pos),
		IsConst: false,
		NotDead: true,
	}
}

func (o *Optimizer) optimizeGet(n *ir.GetNode, pos *tokens.Pos) *Result {
	m := o.OptimizeNode(n.Map)
	k := o.OptimizeNode(n.Key)

	// If there were map literals, then check if const and if so, assign
	return &Result{
		Stmt:    ir.NewCallNode(ir.NewGetNode(m.Stmt, k.Stmt), pos),
		IsConst: false,
	}
}

func (o *Optimizer) optimizeArray(n *ir.ArrayNode, pos *tokens.Pos) *Result {
	results := make([]ir.Node, len(n.Values))
	for i, v := range n.Values {
		r := o.OptimizeNode(v)
		results[i] = r.Stmt
	}
	return &Result{
		Stmt:    ir.NewCallNode(ir.NewArrayNode(results, n.Type()), pos),
		IsConst: false,
	}
}

func (o *Optimizer) optimizeAppend(n *ir.AppendNode, pos *tokens.Pos) *Result {
	// Don't optimize array since its a pointer
	v := o.OptimizeNode(n.Value)
	return &Result{
		Stmt: ir.NewCallNode(&ir.AppendNode{
			Array: n.Array,
			Value: v.Stmt,
		}, pos),
		NotDead: true,
	}
}
