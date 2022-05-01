package types

import (
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
	[]rune("BYTE"),
	[]rune("FLOAT"),
	[]rune("BOOL"),
	[]rune("NULL"),
	[]rune("STRUCT"),
	[]rune("ANY"),
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

		case ' ', '\t', '\n', '\r':
			continue

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
				// add ident
				ident := ""
			loop:
				for i < len(val) {
					switch char {
					case ',', ':', '}', '{', ' ', '\t', '\n':
						break loop
					}
					ident += string(char)
					if i >= len(val) {
						break
					}
					i++
					if i < len(val) {
						char = val[i]
					}
				}

				tokens = append(tokens, token{tokenTypeIdent, &ident})
				i--
			}
		}
	}

	return tokens, nil
}
