package ir

import (
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type VarNode struct {
	ID  int
	typ types.Type

	name string // Variable name for use in Code()
}

func (v *VarNode) Type() types.Type { return v.typ }
func (v *VarNode) Args() []Node     { return []Node{} }

type DefineNode struct {
	NullCall
	Var     int
	Value   Node
	InScope bool // Check if accessing var outside of function (helps with recursion)

	name string // Variable name for use in Code()
}

func (d *DefineNode) Args() []Node { return []Node{d.Value} }

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
				b.Error(ErrorLevelError, args[0].Pos(), "undefined variable: %s", name)
				return NewTypedValue(types.INVALID), nil
			}
			va := b.Scope.Variable(v)
			if !va.NeedsGlobal && va.ScopeType == ScopeTypeGlobal && b.Scope.HasType(ScopeTypeFunction) {
				va.NeedsGlobal = true
			}

			return &VarNode{
				ID:   v,
				typ:  va.Type,
				name: name,
			}, nil
		},
	}

	nodeBuilders["DEFINE"] = nodeBuilder{
		ArgTypes: []types.Type{types.IDENT, types.ALL},
		Build: func(b *Builder, pos *tokens.Pos, args []Node) (Call, error) {
			// Get var
			name := args[0].(*Const).Value.(string)
			id, exists := b.Scope.GetVar(name)
			if !exists {
				id = b.Scope.AddVariable(name, args[1].Type(), pos)
			} else {
				v := b.Scope.Variable(id)
				if !v.Type.Equal(args[1].Type()) {
					// Check if in current scope, if so redefine
					_, exists = b.Scope.CurrScopeGetVar(name)
					if !exists {
						id = b.Scope.AddVariable(name, args[1].Type(), pos)
					} else {
						return nil, pos.Error("cannot redefine variable %s to type %s", name, args[1].Type())
					}
				}
			}
			_, exists = b.Scope.CurrScopeGetVar(name)

			// check if needs global
			v := b.Scope.Variable(id)
			if !v.NeedsGlobal && v.ScopeType == ScopeTypeGlobal && b.Scope.HasType(ScopeTypeFunction) {
				v.NeedsGlobal = true
			}

			return &DefineNode{
				Var:     id,
				Value:   args[1],
				name:    name,
				InScope: exists,
			}, nil
		},
	}
}
