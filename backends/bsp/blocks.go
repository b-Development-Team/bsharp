package bsp

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
)

func (b *BSP) Tab(val string) string {
	lines := strings.Split(val, "\n")
	out := &strings.Builder{}
	for _, line := range lines {
		out.WriteString(b.Config.Tab)
		out.WriteString(line)
	}
	return out.String()
}

func (b *BSP) buildIf(n *ir.IfNode) (string, error) {
	out := &strings.Builder{}
	cond, err := b.buildNode(n.Condition)
	if err != nil {
		return "", err
	}
	fmt.Fprintf(out, "[IF%s%s%s", b.Config.Seperator, cond, b.Config.Newline)
	for _, v := range n.Body {
		s, err := b.buildNode(v)
		if err != nil {
			return "", err
		}
		out.WriteString(b.Tab(s))
		out.WriteString(b.Config.Newline)
	}
	if n.Else != nil {
		fmt.Fprintf(out, "ELSE%s", b.Config.Newline)
		for _, v := range n.Else {
			s, err := b.buildNode(v)
			if err != nil {
				return "", err
			}
			out.WriteString(b.Tab(s))
			out.WriteString(b.Config.Newline)
		}
	}
	fmt.Fprintf(out, "]")
	return out.String(), nil
}

func (b *BSP) buildWhile(n *ir.WhileNode) (string, error) {
	out := &strings.Builder{}
	cond, err := b.buildNode(n.Condition)
	if err != nil {
		return "", err
	}
	fmt.Fprintf(out, "[WHILE%s%s%s", b.Config.Seperator, cond, b.Config.Newline)
	for _, v := range n.Body {
		s, err := b.buildNode(v)
		if err != nil {
			return "", err
		}
		out.WriteString(b.Tab(s))
		out.WriteString(b.Config.Newline)
	}
	fmt.Fprintf(out, "]")
	return out.String(), nil
}

func (b *BSP) buildCase(n *ir.Case) (string, error) {
	out := &strings.Builder{}
	cond, err := b.buildNode(n.Value)
	if err != nil {
		return "", err
	}
	fmt.Fprintf(out, "[CASE%s%s%s", b.Config.Seperator, cond, b.Config.Newline)
	for _, v := range n.Body {
		s, err := b.buildNode(v)
		if err != nil {
			return "", err
		}
		out.WriteString(b.Tab(s))
		out.WriteString(b.Config.Newline)
	}
	fmt.Fprintf(out, "]")
	return out.String(), nil
}

func (b *BSP) buildSwitch(n *ir.SwitchNode) (string, error) {
	out := &strings.Builder{}
	cond, err := b.buildNode(n.Value)
	if err != nil {
		return "", err
	}
	fmt.Fprintf(out, "[SWITCH%s%s%s", b.Config.Seperator, cond, b.Config.Newline)
	for _, v := range n.Cases {
		s, err := b.buildCase(v.Block.(*ir.Case))
		if err != nil {
			return "", err
		}
		out.WriteString(b.Tab(s))
		out.WriteString(b.Config.Newline)
	}
	if n.Default != nil {
		fmt.Fprintf(out, "[DEFAULT%s", b.Config.Newline)
		for _, v := range n.Default.Block.(*ir.Default).Body {
			s, err := b.buildNode(v)
			if err != nil {
				return "", err
			}
			out.WriteString(b.Tab(s))
			out.WriteString(b.Config.Newline)
		}
		fmt.Fprintf(out, "]%s", b.Config.Newline)
	}
	fmt.Fprintf(out, "]")
	return out.String(), nil
}
