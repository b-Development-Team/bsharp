package bstar

import (
	"github.com/Nv7-Github/bsharp/ir"
)

func (b *BStar) buildNodes(n []ir.Node) ([]Node, error) {
	out := make([]Node, 0)
	for _, n := range n {
		node, err := b.buildNode(n)
		if err != nil {
			return nil, err
		}
		_, ok := node.(*BlockNode)
		if ok && node.(*BlockNode).DoesPrint {
			node = b.noPrint(node)
		}
		out = append(out, node)
	}
	return out, nil
}

func (b *BStar) buildBlock(n *ir.BlockNode) (Node, error) {
	switch bl := n.Block.(type) {
	case *ir.IfNode:
		cond, err := b.buildNode(bl.Condition)
		if err != nil {
			return nil, err
		}
		bod, err := b.buildNodes(bl.Body)
		if err != nil {
			return nil, err
		}
		bod = append(bod, b.noPrintNode()) // make sure it doesnt return anything
		if bl.Else == nil {
			return blockNode(false, constNode("IF"), cond, blockNode(false, append([]Node{constNode("BLOCK")}, bod...)...), b.noPrintNode()), nil
		}
		els, err := b.buildNodes(bl.Else)
		if err != nil {
			return nil, err
		}
		els = append(els, b.noPrintNode()) // make sure it doesnt return anything
		return blockNode(false, constNode("IF"), cond, blockNode(false, append([]Node{constNode("BLOCK")}, bod...)...), blockNode(false, append([]Node{constNode("BLOCK")}, els...)...)), nil

	case *ir.WhileNode:
		cond, err := b.buildNode(bl.Condition)
		if err != nil {
			return nil, err
		}
		bod, err := b.buildNodes(bl.Body)
		if err != nil {
			return nil, err
		}
		bod = append([]Node{constNode("WHILE"), cond}, blockNode(false, append([]Node{constNode("BLOCK")}, append(bod, b.noPrintNode())...)...))
		return blockNode(true, bod...), nil

	default:
		return nil, n.Pos().Error("unknown block: %T", bl)
	}
}
