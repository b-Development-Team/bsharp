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
	err := MatchTypes(n.Pos(), args, ext.Params)
	if err != nil {
		return nil, err
	}

	return &ExtensionCall{
		Name: n.Name,
		Args: args,
		typ:  ext.RetType,
		pos:  n.Pos(),
	}, nil
}