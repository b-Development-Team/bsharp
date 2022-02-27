package ir

import (
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type PrintNode struct {
	NullCall
	Arg Node
}

type ConcatNode struct {
	Values []Node
}

func (c *ConcatNode) Type() types.Type { return types.STRING }

type RandintNode struct {
	Lower Node
	Upper Node
}

func (r *RandintNode) Type() types.Type { return types.INT }

type RandomNode struct {
	Lower Node
	Upper Node
}

func (r *RandomNode) Type() types.Type { return types.FLOAT }

func init() {
	nodeBuilders["PRINT"] = nodeBuilder{
		ArgTypes: []types.Type{types.STRING},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return &PrintNode{
				Arg: args[0],
			}, nil
		},
	}

	nodeBuilders["CONCAT"] = nodeBuilder{
		ArgTypes: []types.Type{types.STRING, types.VARIADIC},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return &ConcatNode{
				Values: args,
			}, nil
		},
	}

	nodeBuilders["RANDINT"] = nodeBuilder{
		ArgTypes: []types.Type{types.INT, types.INT},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return &RandintNode{
				Lower: args[0],
				Upper: args[1],
			}, nil
		},
	}

	nodeBuilders["RANDOM"] = nodeBuilder{
		ArgTypes: []types.Type{types.FLOAT, types.FLOAT},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return &RandomNode{
				Lower: args[0],
				Upper: args[1],
			}, nil
		},
	}
}
