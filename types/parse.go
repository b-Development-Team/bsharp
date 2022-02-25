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

	// Map or nothing
	if strings.HasPrefix(typ, "map{") {
		// Get key and value types
		typ = typ[4:]

		// Get key by matching brackets
		openBrackets := 1
		key := ""
		for len(typ) > 0 {
			c := typ[0]
			if c == '{' {
				openBrackets++
			} else if c == '}' {
				openBrackets--
			}

			if openBrackets == 0 {
				break
			} else {
				key += string(c)
				typ = typ[1:]
			}
		}

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
