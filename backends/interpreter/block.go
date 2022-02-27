package interpreter

import "github.com/Nv7-Github/bsharp/ir"

func (i *Interpreter) evalIfNode(n *ir.IfNode) error {
	cond, err := i.evalNode(n.Condition)
	if err != nil {
		return err
	}
	if cond.Value.(bool) {
		i.scope.Push()
		for _, node := range n.Body {
			if _, err := i.evalNode(node); err != nil {
				return err
			}
		}
		i.pop()
	} else if n.Else != nil {
		i.scope.Push()
		for _, node := range n.Else {
			if _, err := i.evalNode(node); err != nil {
				return err
			}
		}
		i.pop()
	}
	return nil
}

func (i *Interpreter) evalWhileNode(n *ir.WhileNode) error {
	for {
		if i.retVal != nil {
			return nil
		}
		cond, err := i.evalNode(n.Condition)
		if err != nil {
			return err
		}
		if !cond.Value.(bool) {
			break
		}

		i.scope.Push()
		for _, node := range n.Body {
			if _, err := i.evalNode(node); err != nil {
				return err
			}
		}
		i.pop()
	}
	return nil
}
