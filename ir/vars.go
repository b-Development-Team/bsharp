package ir

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type VarNode struct {
	ID  int
	typ types.Type
}

func (v *VarNode) Type() types.Type { return v.typ }

func init() {
	nodeBuilders["VAR"] = nodeBuilder{
		ArgTypes: []types.Type{types.IDENT},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			name := args[0].(*Const).Value.(string)
			v, exists := b.Scope.GetVar(name)
			if !exists {
				return nil, fmt.Errorf("unknown variable: %s", name)
			}
			va := b.Scope.Variable(v)

			return &VarNode{
				ID:  v,
				typ: va.Type,
			}, nil
		},
	}
}
