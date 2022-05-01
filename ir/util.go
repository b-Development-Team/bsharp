package ir

import (
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type NullCall struct{}

func (n NullCall) Type() types.Type { return types.NULL }

func MatchTypes(pos *tokens.Pos, args []Node, typs []types.Type) error {
	if len(typs) > 0 && typs[len(typs)-1] == types.VARIADIC { // last type is variadic, check up to that
		if len(args) < len(typs)-2 { // Ignore last 2 since variadic means 0 or more
			return pos.Error("wrong number of arguments: expected at least %d, got %d", len(typs)-2, len(args))
		}
		for i, v := range args {
			if i < len(typs)-2 && !typs[i].Equal(v.Type()) {
				return args[i].Pos().Error("wrong argument type: expected %s, got %s", typs[i], args[i].Type())
			}

			if i >= len(typs)-2 && !typs[len(typs)-2].Equal(v.Type()) {
				return args[i].Pos().Error("wrong variadic argument type: expected %s, got %s", typs[len(typs)-2], args[i].Type())
			}
		}

		return nil
	}
	if len(args) != len(typs) {
		return pos.Error("wrong number of arguments: expected %d, got %d", len(typs), len(args))
	}
	for i, arg := range args {
		if !typs[i].Equal(arg.Type()) {
			return arg.Pos().Error("wrong argument type: expected %s, got %s", typs[i], arg.Type())
		}
	}
	return nil
}
