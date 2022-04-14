package bspgen

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/ssa"
)

func (b *BSPGen) Add(id ssa.ID) ir.Node {
	info := b.Instructions[id]
	instr := b.Blocks[info.Block].Instructions[id]
	switch i := instr.(type) {
	case *ssa.LiveIRValue:
		switch i.Kind {
		case ssa.IRNodePrint:
			return ir.NewCallNode(&ir.PrintNode{Arg: b.Add(i.Params[0])}, info.Pos)

		default:
			panic(fmt.Errorf("unknown live ir kind: %s", i.Kind.String()))
		}

	case *ssa.Const:
		return ir.NewConst(i.Typ, info.Pos, i.Value)

	default:
		panic(fmt.Errorf("unknown instruction type: %T", i))
	}
}
