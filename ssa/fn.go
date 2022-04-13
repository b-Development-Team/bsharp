package ssa

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/types"
)

type FnCall struct {
	Fn     ID
	Params []ID
	Typ    types.Type
}

func (f *FnCall) Type() types.Type { return f.Typ }
func (f *FnCall) String() string {
	out := &strings.Builder{}
	fmt.Fprintf(out, "FnCall [%s](", f.Fn.String())
	for i, val := range f.Params {
		out.WriteString(val.String())
		if i != len(f.Params)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteRune(')')
	return out.String()
}
func (f *FnCall) Args() []ID { return append([]ID{f.Fn}, f.Params...) }
func (f *FnCall) SetArgs(args []ID) {
	f.Fn = args[0]
	f.Params = args[1:]
}

type ExtensionCall struct {
	Fn     string
	Params []ID
	Typ    types.Type
}

func (e *ExtensionCall) Type() types.Type { return e.Typ }
func (e *ExtensionCall) String() string {
	out := &strings.Builder{}
	fmt.Fprintf(out, "FnCall [%s](", e.Fn)
	for i, val := range e.Params {
		out.WriteString(val.String())
		if i != len(e.Params)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteRune(')')
	return out.String()
}
func (e *ExtensionCall) Args() []ID        { return e.Params }
func (e *ExtensionCall) SetArgs(args []ID) { e.Params = args }
