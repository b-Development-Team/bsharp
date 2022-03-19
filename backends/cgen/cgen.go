package cgen

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
)

type Code struct {
	Pre   string
	Value string
}

type CodeConfig struct {
	Tab string
}

func DefaultCodeConfig() CodeConfig {
	return CodeConfig{Tab: "\t"}
}

func JoinCode(vals ...string) string {
	out := &strings.Builder{}
	first := true
	for _, val := range vals {
		if val == "" {
			continue
		}
		if !first {
			out.WriteString("\n")
		}
		out.WriteString(val)
		first = false
	}
	return out.String()
}

func Indent(code string, cnf CodeConfig) string {
	out := &strings.Builder{}
	for _, line := range strings.Split(code, "\n") {
		out.WriteString(cnf.Tab + line + "\n")
	}
	return out.String()
}

type CGen struct {
	Config CodeConfig
	stack  *stack
	ir     *ir.IR
	tmps   map[string]int
}

func NewCGen(i *ir.IR) *CGen {
	return &CGen{
		Config: DefaultCodeConfig(),
		stack: &stack{
			vals: make([]scope, 0),
		},
		tmps: make(map[string]int),
		ir:   i,
	}
}

//go:embed std/std.c
var std string

func (c *CGen) Build() (string, error) {
	out := &strings.Builder{}
	out.WriteString(std)

	// Add fn types
	for _, fn := range c.ir.Funcs {
		fmt.Fprintf(out, "%s %s(", c.CType(fn.RetType), Namespace+fn.Name)
		for i, arg := range fn.Params {
			fmt.Fprintf(out, "%s %s", c.CType(arg.Type), arg.Name)
			if i != len(fn.Params)-1 {
				out.WriteString(", ")
			}
		}
		out.WriteString(");\n")
	}
	out.WriteString("\n")

	// Add fns
	for _, fn := range c.ir.Funcs {
		fmt.Fprintf(out, "%s %s(", c.CType(fn.RetType), Namespace+fn.Name)
		for i, arg := range fn.Params {
			fmt.Fprintf(out, "%s %s", c.CType(arg.Type), arg.Name)
			if i != len(fn.Params)-1 {
				out.WriteString(", ")
			}
		}
		out.WriteString(") {\n")
		c.stack.Push()
		for _, stmt := range fn.Body {
			code, err := c.AddNode(stmt)
			if err != nil {
				return "", err
			}
			if code.Pre != "" {
				out.WriteString(Indent(code.Pre, c.Config))
			}
			if code.Value != "" {
				out.WriteString("\n" + Indent(code.Value+"\n", c.Config))
			}
		}
		free := c.stack.FreeCode()
		c.stack.Pop()
		if free != "" {
			out.WriteString(Indent(free, c.Config))
		}
		out.WriteString("}\n\n")
	}

	// Add main
	out.WriteString("int main() {\n")
	c.stack.Push()
	for _, stmt := range c.ir.Body {
		code, err := c.AddNode(stmt)
		if err != nil {
			return "", err
		}
		if code.Pre != "" {
			out.WriteString(Indent(code.Pre+"\n", c.Config))
		}
		if code.Value != "" {
			out.WriteString(Indent(code.Value+"\n", c.Config))
		}
	}
	free := c.stack.FreeCode()
	c.stack.Pop()
	if free != "" {
		out.WriteString(Indent(free, c.Config))
	}
	out.WriteString(c.Config.Tab + "return 0;\n}\n")

	return out.String(), nil
}

func (c *CGen) GetTmp(name string) string {
	cnt := c.tmps[name]
	c.tmps[name]++
	return fmt.Sprintf("%s%s_%d", Namespace, name, cnt)
}
