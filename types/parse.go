package types

import (
	"fmt"
)

func parse(tokens []token) (Type, []token, error) {
	if len(tokens) < 1 {
		return nil, nil, fmt.Errorf("expected type")
	}
	switch tokens[0].typ {
	case tokenTypeConst:
		switch *tokens[0].value {
		case "INT":
			return INT, tokens[1:], nil

		case "FLOAT":
			return FLOAT, tokens[1:], nil

		case "BOOL":
			return BOOL, tokens[1:], nil

		case "NULL":
			return NULL, tokens[1:], nil

		case "STRING":
			return STRING, tokens[1:], nil

		case "ARRAY":
			// Eat "ARRAY"
			tokens = tokens[1:]
			if tokens[0].typ != tokenTypeLBrack {
				return nil, nil, fmt.Errorf("expected '{' after 'ARRAY'")
			}

			// Eat "{"
			tokens = tokens[1:]

			// Get element type
			elemType, tokens, err := parse(tokens)
			if err != nil {
				return nil, nil, err
			}

			// Eat "}"
			if tokens[0].typ != tokenTypeRBrack {
				return nil, nil, fmt.Errorf("expected '}' after ARRAY element type")
			}
			tokens = tokens[1:]

			return NewArrayType(elemType), tokens, nil

		case "MAP":
			// Eat "MAP"
			tokens = tokens[1:]
			if tokens[0].typ != tokenTypeLBrack {
				return nil, nil, fmt.Errorf("expected '{' after 'MAP'")
			}

			// Eat "{"
			tokens = tokens[1:]

			// Get key type
			keyType, tokens, err := parse(tokens)
			if err != nil {
				return nil, nil, err
			}

			// Eat ","
			if tokens[0].typ != tokenTypeComma {
				return nil, nil, fmt.Errorf("expected ',' after MAP key type")
			}
			tokens = tokens[1:]

			// Get value type
			valueType, tokens, err := parse(tokens)
			if err != nil {
				return nil, nil, err
			}

			// Eat "}"
			if tokens[0].typ != tokenTypeRBrack {
				return nil, nil, fmt.Errorf("expected '}' after MAP value type")
			}
			tokens = tokens[1:]

			return NewMapType(keyType, valueType), tokens, nil

		case "FUNC":
			// Eat "FUNC"
			tokens = tokens[1:]
			if tokens[0].typ != tokenTypeLBrack {
				return nil, nil, fmt.Errorf("expected '{' after 'FUNC'")
			}

			// Eat "{"
			tokens = tokens[1:]

			// Get argument types
			argTypes := make([]Type, 0)
			for {
				var argType Type
				var err error
				argType, tokens, err = parse(tokens)
				if err != nil {
					return nil, nil, err
				}

				argTypes = append(argTypes, argType)

				if tokens[0].typ == tokenTypeRBrack {
					break
				}

				if tokens[0].typ != tokenTypeComma {
					return nil, nil, fmt.Errorf("expected ',' or '}' after FUNC argument type")
				}
				tokens = tokens[1:]
			}

			// Eat "}"
			if tokens[0].typ != tokenTypeRBrack {
				return nil, nil, fmt.Errorf("expected '}' after FUNC argument types")
			}
			tokens = tokens[1:]

			// Get return type
			returnType, tokens, err := parse(tokens)
			if err != nil {
				return nil, nil, err
			}

			return NewFuncType(argTypes, returnType), tokens, nil

		case "STRUCT":
			// Eat "STRUCT"
			tokens = tokens[1:]
			if tokens[0].typ != tokenTypeLBrack {
				return nil, nil, fmt.Errorf("expected '{' after 'STRUCT'")
			}

			// Eat "{"
			tokens = tokens[1:]

			// Get fields
			fields := make([]StructField, 0)
			for {
				// Check if done
				if tokens[0].typ == tokenTypeRBrack {
					break
				}

				// Get name
				if tokens[0].typ != tokenTypeIdent {
					return nil, nil, fmt.Errorf("expected field name")
				}
				name := *tokens[0].value
				tokens = tokens[1:]

				// Eat ":"
				if tokens[0].typ != tokenTypeColon {
					return nil, nil, fmt.Errorf("expected ':' after field name")
				}
				tokens = tokens[1:]

				// Get type
				var fieldType Type
				var err error
				fieldType, tokens, err = parse(tokens)
				if err != nil {
					return nil, nil, err
				}

				fields = append(fields, NewStructField(name, fieldType))
			}

			// Eat RBrack
			tokens = tokens[1:]
			if tokens[0].typ != tokenTypeRBrack {
				return nil, nil, fmt.Errorf("expected '}' after STRUCT fields")
			}

			return NewStruct(fields...), tokens, nil

		default:
			return nil, nil, fmt.Errorf("unexpected const: %s", *tokens[0].value)
		}
	}

	return nil, nil, fmt.Errorf("unexpected token: %s", tokens[0].String())
}
