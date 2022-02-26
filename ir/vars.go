package ir

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type VarNode struct {
	ID  int
	typ types.Type
}

func (v *VarNode) Type() types.Type { return v.typ }

type DefineNode struct {
	NullCall
	Var   int
	Value Node
}

func init() {
	nodeBuilders["VAR"] = nodeBuilder{
		ArgTypes: []types.Type{types.IDENT},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			name := args[0].(*Const).Value.(string)
			v, exists := b.Scope.GetVar(name)
			if !exists {
				return nil, fmt.Errorf("unknown variable: %s", name)
			}
			va := b.Scope.Variable(v)

			return &VarNode{
				ID:  v,
				typ: va.Type,
			}, nil
		},
	}

	nodeBuilders["DEFINE"] = nodeBuilder{
		ArgTypes: []types.Type{types.IDENT, types.ANY},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			// Get var
			name := args[0].(*Const).Value.(string)
			id, exists := b.Scope.CurrScopeGetVar(name) // Only check in current scope
			if !exists {
				id = b.Scope.AddVariable(name, args[1].Type(), pos)
			} else {
				v := b.Scope.Variable(id)
				if v.Type != args[1].Type() {
					return nil, fmt.Errorf("cannot redefine variable %s to type %s", name, args[1].Type())
				}
			}

			return &DefineNode{
				Var:   id,
				Value: args[1],
			}, nil
		},
	}
}
