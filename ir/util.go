package ir

import (
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type NullCall struct{}

func (n NullCall) Type() types.Type { return types.NULL }

func (b *Builder) MatchTypes(pos *tokens.Pos, args []Node, typs []types.Type) bool {
	if len(typs) > 0 && typs[len(typs)-1] == types.VARIADIC { // last type is variadic, check up to that
		if len(args) < len(typs)-2 { // Ignore last 2 since variadic means 0 or more
			b.Error(ErrorLevelError, pos, "wrong number of arguments: expected at least %d, got %d", len(typs)-2, len(args))
			return true
		}

		hasErr := false
		for i, v := range args {
			if i < len(typs)-2 && !typs[i].Equal(v.Type()) && !types.INVALID.Equal(v.Type()) {
				b.Error(ErrorLevelError, pos, "wrong argument type: expected %s, got %s", typs[i], args[i].Type())
				hasErr = true
			}

			if i >= len(typs)-2 && !typs[len(typs)-2].Equal(v.Type()) && !types.INVALID.Equal(v.Type()) {
				b.Error(ErrorLevelError, pos, "wrong variadic argument type: expected %s, got %s", typs[len(typs)-2], args[i].Type())
				hasErr = true
			}
		}

		return hasErr
	}

	if len(args) != len(typs) {
		b.Error(ErrorLevelError, pos, "wrong number of arguments: expected %d, got %d", len(typs), len(args))
		return true
	}

	hasErr := false
	for i, arg := range args {
		if !typs[i].Equal(arg.Type()) && !types.INVALID.Equal(arg.Type()) {
			b.Error(ErrorLevelError, pos, "wrong argument type: expected %s, got %s", typs[i], arg.Type())
			hasErr = true
		}
	}
	return hasErr
}

func (b *Builder) FixTypes(args *[]Node, typs []types.Type, pos *tokens.Pos) {
	if len(typs) > 0 && typs[len(typs)-1] == types.VARIADIC {
		if len(*args) < len(typs)-1 {
			start := len(*args) - 1
			*args = append(*args, make([]Node, len(typs)-len(*args)-1)...)
			for i := start; i < len(*args); i++ {
				(*args)[i] = NewTypedNode(typs[i], pos)
			}
		}
	}

	if len(*args) < len(typs) {
		start := len(*args) - 1
		*args = append(*args, make([]Node, len(typs)-len(*args))...)
		for i := start; i < len(*args); i++ {
			(*args)[i] = NewTypedNode(typs[i], pos)
		}
	}
}
