package ssa

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/types"
)

// Everything in this file is not real SSA, will be removed in a pass

type SetVariable struct {
	Variable int
	Value    ID
}

func (s *SetVariable) Type() types.Type { return types.NULL }
func (s *SetVariable) String() string {
	return fmt.Sprintf("SetVariable (%s) -> [%d]", s.Value.String(), s.Variable)
}

type GetVariable struct {
	Variable int
	Typ      types.Type
}

func (g *GetVariable) Type() types.Type { return g.Typ }
func (g *GetVariable) String() string {
	return fmt.Sprintf("GetVariable [%d]", g.Variable)
}
