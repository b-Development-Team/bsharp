package bsp

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

type BSP struct {
	ir     *ir.IR
	Config *BSPConfig
}

type BSPConfig struct {
	Seperator string
	Tab       string
	Newline   string
}

func NewBSP(i *ir.IR) *BSP {
	return &BSP{
		ir: i,
		Config: &BSPConfig{
			Seperator: " ",
			Tab:       "\t",
			Newline:   "\n",
		},
	}
}

func (b *BSP) Build() (string, error) {
	out := &strings.Builder{}

	// Move global vars to the top
	toRemove := make([]int, 0)
	pre := make([]ir.Node, 0)
	for i, v := range b.ir.Body {
		c, ok := v.(*ir.CallNode)
		if !ok {
			continue
		}
		d, ok := c.Call.(*ir.DefineNode)
		if !ok {
			continue
		}
		if b.ir.Variables[d.Var].NeedsGlobal {
			toRemove = append(toRemove, i)
			pre = append(pre, v)
		}
	}

	// Remove those definitions
	newBody := make([]ir.Node, 0, len(b.ir.Body)-len(toRemove))
	for i, v := range b.ir.Body {
		if len(toRemove) > 0 && i == toRemove[0] {
			toRemove = toRemove[1:]
			continue
		}
		newBody = append(newBody, v)
	}
	b.ir.Body = newBody

	// Add pre
	for _, n := range pre {
		v, err := b.buildNode(n)
		if err != nil {
			return "", err
		}
		out.WriteString(v)
		out.WriteString(b.Config.Newline)
	}

	// Add functions
	for _, fn := range b.ir.Funcs {
		fmt.Fprintf(out, "[FUNC%s%s", b.Config.Seperator, fn.Name)
		for _, p := range fn.Params {
			fmt.Fprintf(out, "%s[PARAM%s%s%s%s]", b.Config.Seperator, b.Config.Seperator, p.Name, b.Config.Seperator, p.Type.String())
		}
		if !types.NULL.Equal(fn.RetType) {
			fmt.Fprintf(out, "%s[RETURNS%s%s]", b.Config.Seperator, b.Config.Seperator, fn.RetType.String())
		}
		out.WriteString(b.Config.Newline)
		for _, n := range fn.Body {
			v, err := b.buildNode(n)
			if err != nil {
				return "", err
			}
			out.WriteString(b.Tab(v))
			out.WriteString(b.Config.Newline)
		}
		out.WriteString("]" + b.Config.Newline)
	}

	// Add body
	for _, n := range b.ir.Body {
		v, err := b.buildNode(n)
		if err != nil {
			return "", err
		}
		out.WriteString(v)
		out.WriteString(b.Config.Newline)
	}
	return out.String(), nil
}
