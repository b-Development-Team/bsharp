package ir

import (
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type Extension struct {
	Name    string
	Params  []types.Type
	RetType types.Type
}

type ExtensionCall struct {
	Name string
	Args []Node
	typ  types.Type
	pos  *tokens.Pos
}

func (e *ExtensionCall) Type() types.Type { return e.typ }
func (e *ExtensionCall) Pos() *tokens.Pos { return e.pos }

func (b *Builder) buildExtensionCall(n *parser.CallNode) (Node, error) {
	ext := b.extensions[n.Name]

	// Add params
	args := make([]Node, len(n.Args))
	for i, arg := range n.Args {
		node, err := b.buildNode(arg)
		if err != nil {
			return nil, err
		}
		args[i] = node
	}
	err := b.MatchTypes(n.Pos(), args, ext.Params)
	if err {
		return NewTypedNode(ext.RetType, n.Pos()), nil
	}

	return &ExtensionCall{
		Name: n.Name,
		Args: args,
		typ:  ext.RetType,
		pos:  n.Pos(),
	}, nil
}

type BuiltinFn struct {
	Name   string
	Params []types.Type // nil if block builder
}

func BuiltinFns() []*BuiltinFn {
	out := make([]*BuiltinFn, len(nodeBuilders)+len(blockBuilders))
	i := 0
	for name, builder := range nodeBuilders {
		out[i] = &BuiltinFn{
			Name:   name,
			Params: builder.ArgTypes,
		}
		i++
	}
	for name := range blockBuilders {
		out[i] = &BuiltinFn{
			Name: name,
		}
		i++
	}
	return out
}
