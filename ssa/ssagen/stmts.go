package ssagen

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/ssa"
)

func (s *SSAGen) Add(node ir.Node) ssa.ID {
	switch n := node.(type) {
	default:
		panic(fmt.Sprintf("unknown node type: %T", n))
	}
}
