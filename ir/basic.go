package ir

import (
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type PrintNode struct {
	NullCall
	Arg Node
}

func (p *PrintNode) Args() []Node { return []Node{p.Arg} }

type ConcatNode struct {
	Values []Node
}

func (c *ConcatNode) Type() types.Type { return types.STRING }
func (c *ConcatNode) Args() []Node {
	return c.Values
}

type TimeMode int

const (
	TimeModeSeconds TimeMode = iota
	TimeModeMicro
	TimeModeMilli
	TimeModeNano
)

var timeModeNames = map[string]TimeMode{
	"SECONDS": TimeModeSeconds,
	"MICRO":   TimeModeMicro,
	"MILLI":   TimeModeMilli,
	"NANO":    TimeModeNano,
}

type TimeNode struct {
	Mode TimeMode
	pos  *tokens.Pos
}

func (t *TimeNode) Type() types.Type { return types.INT }
func (t *TimeNode) Args() []Node     { return []Node{NewConst(types.INT, t.pos, t.Mode)} }

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

	nodeBuilders["TIME"] = nodeBuilder{
		ArgTypes: []types.Type{types.IDENT, types.VARIADIC},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			if len(args) > 1 {
				return nil, pos.Error("TIME takes 0 or 1 arguments")
			}

			mode := TimeModeSeconds
			if len(args) == 1 {
				i, ok := args[0].(*Const)
				if !ok {
					return nil, pos.Error("TIME mode must be an identifier")
				}

				m, ok := timeModeNames[i.Value.(string)]
				if !ok {
					return nil, pos.Error("TIME mode must be SECONDS, MICRO, MILLI, or NANO")
				}
				mode = m
			}

			return &TimeNode{
				Mode: mode,
				pos:  pos,
			}, nil
		},
	}
}
