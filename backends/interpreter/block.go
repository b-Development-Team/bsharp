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

func (i *Interpreter) evalSwitchNode(n *ir.SwitchNode) error {
	v, err := i.evalNode(n.Value)
	if err != nil {
		return err
	}

	// Check cases (O(n), not O(1) like expected from switch)
	for _, cs := range n.Cases {
		if cs.Value.Value == v.Value { // This works for int, float, string, all the hashable types
			i.scope.Push()
			for _, node := range cs.Body {
				if _, err := i.evalNode(node); err != nil {
					return err
				}
			}
			i.pop()
			return nil
		}
	}

	// Default case
	if n.Default != nil {
		i.scope.Push()
		for _, node := range n.Default {
			if _, err := i.evalNode(node); err != nil {
				return err
			}
		}
		i.pop()
	}
	return nil
}
