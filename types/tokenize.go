package types

import (
	"fmt"
	"strings"
)

// Tokenize types
type tokenType int

const (
	tokenTypeConst tokenType = iota // ARRAY, MAP, FUNC, STRING, INT, FLOAT, BOOL, NULL
	tokenTypeLBrack
	tokenTypeRBrack
	tokenTypeComma
	tokenTypeIdent
	tokenTypeColon
)

// ARRAY{FUNC{INT,INT}INT}
type token struct {
	typ   tokenType
	value *string
}

func (t token) String() string {
	switch t.typ {
	case tokenTypeConst:
		return *t.value
	case tokenTypeLBrack:
		return "{"
	case tokenTypeRBrack:
		return "}"
	case tokenTypeComma:
		return ","
	case tokenTypeIdent:
		return *t.value
	case tokenTypeColon:
		return ":"
	default:
		return ""
	}
}

var constTokens = [][]rune{
	[]rune("ARRAY"),
	[]rune("MAP"),
	[]rune("FUNC"),
	[]rune("STRING"),
	[]rune("INT"),
	[]rune("FLOAT"),
	[]rune("BOOL"),
	[]rune("NULL"),
	[]rune("STRUCT"),
}

var tokenStarters = map[rune]struct{}{}

func init() {
	for _, v := range constTokens {
		tokenStarters[[]rune(v)[0]] = struct{}{}
	}
}

func tokenize(val []rune) ([]token, error) {
	tokens := make([]token, 0)
	for i := 0; i < len(val); i++ {
		char := val[i]
		switch char {
		case '{':
			tokens = append(tokens, token{tokenTypeLBrack, nil})

		case '}':
			tokens = append(tokens, token{tokenTypeRBrack, nil})

		case ',':
			tokens = append(tokens, token{tokenTypeComma, nil})

		case ':':
			tokens = append(tokens, token{tokenTypeColon, nil})

		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '_':
			// ident
			var ident string
			for i < len(val) && (char == 'a' || char == 'b' || char == 'c' || char == 'd' || char == 'e' || char == 'f' || char == 'g' || char == 'h' || char == 'i' || char == 'j' || char == 'k' || char == 'l' || char == 'm' || char == 'n' || char == 'o' || char == 'p' || char == 'q' || char == 'r' || char == 's' || char == 't' || char == 'u' || char == 'v' || char == 'w' || char == 'x' || char == 'y' || char == 'z' || char == '0' || char == '1' || char == '2' || char == '3' || char == '4' || char == '5' || char == '6' || char == '7' || char == '8' || char == '9' || char == '_' || char == 'A' || char == 'B' || char == 'C' || char == 'D' || char == 'E' || char == 'F' || char == 'G' || char == 'H' || char == 'I' || char == 'J' || char == 'K' || char == 'L' || char == 'M' || char == 'N' || char == 'O' || char == 'P' || char == 'Q' || char == 'R' || char == 'S' || char == 'T' || char == 'U' || char == 'V' || char == 'W' || char == 'X' || char == 'Y' || char == 'Z') {
				ident += string(char)
				if i >= len(val) {
					break
				}
				i++
				char = val[i]
			}
			tokens = append(tokens, token{tokenTypeIdent, &ident})
			i--

		default:
			if _, ok := tokenStarters[char]; ok {
				// Check if token is a const
				for _, v := range constTokens {
					s := string(v)
					if strings.HasPrefix(string(val[i:]), s) {
						tokens = append(tokens, token{tokenTypeConst, &s})
						i += len(v) - 1
						break
					}
				}
			} else {
				return nil, fmt.Errorf("unexpected character in type: %s", string(char))
			}
		}
	}

	return tokens, nil
}
