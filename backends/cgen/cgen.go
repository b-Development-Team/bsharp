package cgen

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

type Code struct {
	Pre   string
	Value string
}

type CodeConfig struct {
	Tab         string
	BoundsCheck bool
}

func DefaultCodeConfig() CodeConfig {
	return CodeConfig{Tab: "\t", BoundsCheck: true}
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
	Config       CodeConfig
	stack        *stack
	ir           *ir.IR
	tmps         map[string]int
	declaredVars []bool
	isReturn     bool
	addedFns     map[string]struct{}

	typIDs map[string]int
	idTyps []types.Type

	globals    *strings.Builder
	globalfns  *strings.Builder
	globaltyps *strings.Builder
}

func NewCGen(i *ir.IR) *CGen {
	return &CGen{
		Config: DefaultCodeConfig(),
		stack: &stack{
			vals: make([]scope, 0),
		},
		tmps:         make(map[string]int),
		ir:           i,
		declaredVars: make([]bool, len(i.Variables)),
		globals:      &strings.Builder{},
		addedFns:     make(map[string]struct{}),
		globalfns:    &strings.Builder{},
		typIDs:       make(map[string]int),
		idTyps:       make([]types.Type, 0),
	}
}

//go:embed std/std.c
var std string

func (c *CGen) addCode(bld *strings.Builder, code *Code) {
	if code.Pre != "" {
		bld.WriteString(Indent(code.Pre, c.Config))
	}
	if code.Value != "" {
		if code.Pre != "" {
			bld.WriteString(Indent("\n"+code.Value+";\n", c.Config))
		} else {
			bld.WriteString(Indent(code.Value+";", c.Config))
		}
	}
}

func (c *CGen) addFree(bld *strings.Builder) {
	free := c.stack.FreeCode()
	if free != "" {
		bld.WriteString(Indent(free, c.Config))
	}
}

func (c *CGen) Build() (string, error) {
	top := &strings.Builder{}
	c.globaltyps = &strings.Builder{}

	out := &strings.Builder{}
	// Add fn types
	for _, fn := range c.ir.Funcs {
		fmt.Fprintf(top, "%s %s(", c.CType(fn.RetType), Namespace+fn.Name)
		for i, arg := range fn.Params {
			fmt.Fprintf(top, "%s", c.CType(arg.Type))
			if i != len(fn.Params)-1 {
				top.WriteString(", ")
			}
		}
		top.WriteString(");\n")
	}

	// Add fns
	for _, fn := range c.ir.Funcs {
		fmt.Fprintf(out, "%s %s(", c.CType(fn.RetType), Namespace+fn.Name)
		for i, arg := range fn.Params {
			fmt.Fprintf(out, "%s %s", c.CType(arg.Type), Namespace+arg.Name+strconv.Itoa(arg.ID))
			if i != len(fn.Params)-1 {
				out.WriteString(", ")
			}
			c.declaredVars[arg.ID] = true
		}
		out.WriteString(") {\n")
		c.stack.Push()
		for _, stmt := range fn.Body {
			code, err := c.AddNode(stmt)
			if err != nil {
				return "", err
			}
			c.addCode(out, code)
		}
		c.addFree(out)
		c.stack.Pop()
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
		c.addCode(out, code)
	}
	c.addFree(out)
	c.stack.Pop()
	out.WriteString(c.Config.Tab + "return 0;\n}\n")

	// Combine
	topCode := strings.Replace(std, "const char* const anytyps[];", c.typCode(), 1)
	code := &strings.Builder{}
	code.WriteString(topCode)
	code.WriteString(c.globaltyps.String())
	code.WriteString(top.String())
	code.WriteRune('\n')
	code.WriteString(c.globals.String())
	code.WriteRune('\n')
	code.WriteString(c.globalfns.String())
	code.WriteString(out.String())

	return code.String(), nil
}

func (c *CGen) GetTmp(name string) string {
	cnt := c.tmps[name]
	c.tmps[name]++
	return fmt.Sprintf("%s%s_%d", Namespace, name, cnt)
}
