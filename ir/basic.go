package ir

import (
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type PrintNode struct {
	NullCall
	Arg Node
}

func init() {
	nodeBuilders["PRINT"] = nodeBuilder{
		ArgTypes: []types.Type{types.STRING},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return &PrintNode{
				Arg: args[0],
			}, nil
		},
	}
}
