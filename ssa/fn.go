package ssa

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/types"
)

type FnCallNode struct {
	Fn     ID
	Params []ID
	Typ    types.Type
}

func (f *FnCallNode) Type() types.Type { return f.Typ }
func (f *FnCallNode) String() string {
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
func (f *FnCallNode) Args() []ID { return append([]ID{f.Fn}, f.Params...) }
func (f *FnCallNode) SetArgs(args []ID) {
	f.Fn = args[0]
	f.Params = args[1:]
}
