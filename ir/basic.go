package ir

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type PrintNode struct {
	NullCall
	Arg Node
}

func (p *PrintNode) Code(cnf CodeConfig) string { return fmt.Sprintf("[PRINT %s]", p.Arg.Code(cnf)) }

type ConcatNode struct {
	Values []Node
}

func (c *ConcatNode) Type() types.Type { return types.STRING }
func (c *ConcatNode) Code(cnf CodeConfig) string {
	args := &strings.Builder{}
	for i, v := range c.Values {
		args.WriteString(v.Code(cnf))
		if i != len(c.Values)-1 {
			args.WriteString(" ")
		}
	}
	return fmt.Sprintf("[CONCAT %s]", args.String())
}

type RandintNode struct {
	Lower Node
	Upper Node
}

func (r *RandintNode) Type() types.Type { return types.INT }

func (r *RandintNode) Code(cnf CodeConfig) string {
	return fmt.Sprintf("[RANDINT %s %s]", r.Lower.Code(cnf), r.Upper.Code(cnf))
}

type RandomNode struct {
	Lower Node
	Upper Node
}

func (r *RandomNode) Type() types.Type { return types.FLOAT }

func (r *RandomNode) Code(cnf CodeConfig) string {
	return fmt.Sprintf("[RANDOM %s %s]", r.Lower.Code(cnf), r.Upper.Code(cnf))
}

type TimeNode struct{}

func (t *TimeNode) Code(cnf CodeConfig) string { return "[TIME]" }
func (t *TimeNode) Type() types.Type           { return types.INT }

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

	nodeBuilders["TIME"] = nodeBuilder{
		ArgTypes: []types.Type{},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return &TimeNode{}, nil
		},
	}
}
