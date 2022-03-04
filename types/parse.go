package types

import (
	"fmt"
	"strings"
)

// int, string, float, int{}, map{string}int, etc.
func ParseType(typ string) (Type, error) {
	for k, v := range basicTypeNames {
		// Either array or basic type
		if strings.HasPrefix(typ, v) {
			if typ == v { // Basic type
				return k, nil
			}
		}
	}

	// Array
	if strings.HasSuffix(typ, "{}") {
		typ, err := ParseType(typ[:len(typ)-2])
		if err != nil {
			return nil, err
		}
		return NewArrayType(typ), nil
	}

	// Function
	if strings.HasPrefix(typ, "FUNC{") {
		typ = strings.TrimPrefix(typ, "FUNC{")
		var v string
		argTyps := make([]Type, 0)
		for len(typ) > 0 && typ[0] != '}' { // Match until closing brace
			v, typ = matchBrackets(typ)
			argTyp, err := ParseType(v)
			if err != nil {
				return nil, err
			}
			argTyps = append(argTyps, argTyp)
		}
		typ = typ[1:] // Get rid of closing bracket
		retTyp := Type(NULL)
		// Check if return type
		if len(typ) > 0 {
			var err error
			retTyp, err = ParseType(typ)
			if err != nil {
				return nil, err
			}
		}
		return NewFuncType(argTyps, retTyp), nil
	}

	// Map or nothing
	if strings.HasPrefix(typ, "MAP{") {
		// Get key and value types
		typ = typ[4:]

		// Get key by matching brackets
		key, typ := matchBrackets(typ)

		// Parse types
		keyType, err := ParseType(key)
		if err != nil {
			return nil, err
		}
		valTyp, err := ParseType(typ[1:])
		if err != nil {
			return nil, err
		}
		return NewMapType(keyType, valTyp), nil
	}
	return nil, fmt.Errorf("unknown type: %s", typ)
}

func matchBrackets(val string) (string, string) {
	openBrackets := 1
	key := ""
	for len(val) > 0 {
		c := val[0]
		if c == '{' {
			openBrackets++
		} else if c == '}' {
			openBrackets--
		}

		if openBrackets == 0 {
			break
		} else {
			key += string(c)
			val = val[1:]
		}
	}
	return key, val
}
