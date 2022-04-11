package ssa

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/types"
)

type Print struct {
	Value ID
}

func (p *Print) Type() types.Type { return types.INT }
func (p *Print) String() string   { return fmt.Sprintf("Print (%s)", p.Value.String()) }
