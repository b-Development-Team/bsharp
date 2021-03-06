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

type PanicNode struct {
	NullCall
	Arg Node
}

func (p *PanicNode) Args() []Node { return []Node{p.Arg} }

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

func (t TimeMode) String() string {
	return [...]string{"SECONDS", "MICRO", "MILLI", "NANO"}[t]
}

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
func (t *TimeNode) Args() []Node     { return []Node{NewConst(types.IDENT, t.pos, t.Mode.String())} }

func init() {
	nodeBuilders["PRINT"] = nodeBuilder{
		ArgTypes: []types.Type{types.STRING},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return &PrintNode{
				Arg: args[0],
			}, nil
		},
	}

	nodeBuilders["PANIC"] = nodeBuilder{
		ArgTypes: []types.Type{types.STRING},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			return &PanicNode{
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
				b.Error(ErrorLevelError, pos, "TIME takes 0 or 1 arguments")
				return &TimeNode{
					Mode: TimeModeSeconds,
					pos:  pos,
				}, nil
			}

			mode := TimeModeSeconds
			if len(args) == 1 {
				i, ok := args[0].(*Const)
				if !ok {
					b.Error(ErrorLevelError, pos, "TIME mode must be a constant")
					return &TimeNode{
						Mode: TimeModeSeconds,
						pos:  pos,
					}, nil
				}

				_, ok = i.Value.(string)
				if !ok {
					b.Error(ErrorLevelError, pos, "TIME mode must be a string")
					return &TimeNode{
						Mode: TimeModeSeconds,
						pos:  pos,
					}, nil
				}

				m, ok := timeModeNames[i.Value.(string)]
				if !ok {
					b.Error(ErrorLevelError, pos, "TIME mode must be SECONDS, MICRO, MILLI, or NANO")
					return &TimeNode{
						Mode: TimeModeSeconds,
						pos:  pos,
					}, nil
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
