package bstar

import (
	"fmt"

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
		bod, err := b.buildNodes(bl.Body)
		if err != nil {
			return nil, err
		}
		bod = append(bod, blockNode(false, constNode("BLOCK"))) // make sure it doesnt return anything
		if bl.Else == nil {
			return blockNode(false, constNode("IF"), blockNode(false, append([]Node{constNode("BLOCK")}, bod...)...), blockNode(false, constNode("BLOCK"))), nil
		}
		els, err := b.buildNodes(bl.Else)
		if err != nil {
			return nil, err
		}
		els = append(els, blockNode(false, constNode("BLOCK"))) // make sure it doesnt return anything
		return blockNode(false, constNode("IF"), blockNode(false, append([]Node{constNode("BLOCK")}, bod...)...), blockNode(false, append([]Node{constNode("BLOCK")}, els...)...)), nil

	default:
		return nil, fmt.Errorf("unknown block: %T", b)
	}
}
