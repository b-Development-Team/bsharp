package ir

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type VarNode struct {
	ID  int
	typ types.Type

	name string // Variable name for use in Code()
}

func (v *VarNode) Type() types.Type { return v.typ }
func (v *VarNode) Code(cnf CodeConfig) string {
	return fmt.Sprintf("[VAR %s]", v.name)
}

type DefineNode struct {
	NullCall
	Var   int
	Value Node

	name string // Variable name for use in Code()
}

func (d *DefineNode) Code(cnf CodeConfig) string {
	return fmt.Sprintf("[DEFINE %s %s]", d.name, d.Value.Code(cnf))
}

func NewDefineNode(id int, val Node, i *IR) *DefineNode {
	return &DefineNode{
		Var:   id,
		Value: val,
		name:  i.Variables[id].Name,
	}
}

func init() {
	nodeBuilders["VAR"] = nodeBuilder{
		ArgTypes: []types.Type{types.IDENT},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			name := args[0].(*Const).Value.(string)
			v, exists := b.Scope.GetVar(name)
			if !exists {
				return nil, args[0].Pos().Error("unknown variable: %s", name)
			}
			va := b.Scope.Variable(v)

			return &VarNode{
				ID:   v,
				typ:  va.Type,
				name: name,
			}, nil
		},
	}

	nodeBuilders["DEFINE"] = nodeBuilder{
		ArgTypes: []types.Type{types.IDENT, types.ANY},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			// Get var
			name := args[0].(*Const).Value.(string)
			id, exists := b.Scope.GetVar(name)
			if !exists {
				id = b.Scope.AddVariable(name, args[1].Type(), pos)
			} else {
				v := b.Scope.Variable(id)
				if v.Type != args[1].Type() {
					// Check if in current scope, if so redefine
					_, exists = b.Scope.CurrScopeGetVar(name)
					if !exists {
						id = b.Scope.AddVariable(name, args[1].Type(), pos)
					} else {
						return nil, pos.Error("cannot redefine variable %s to type %s", name, args[1].Type())
					}
				}
			}

			return &DefineNode{
				Var:   id,
				Value: args[1],
				name:  name,
			}, nil
		},
	}
}
