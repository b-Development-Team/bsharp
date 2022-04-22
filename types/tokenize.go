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
	[]rune("ANY"),
}

var tokenStarters = map[rune]struct{}{}

func init() {
	for _, v := range constTokens {
		tokenStarters[[]rune(v)[0]] = struct{}{}
	}
}

var ValidIdentLetters = map[rune]struct{}{
	'a': {},
	'b': {},
	'c': {},
	'd': {},
	'e': {},
	'f': {},
	'g': {},
	'h': {},
	'i': {},
	'j': {},
	'k': {},
	'l': {},
	'm': {},
	'n': {},
	'o': {},
	'p': {},
	'q': {},
	'r': {},
	's': {},
	't': {},
	'u': {},
	'v': {},
	'w': {},
	'x': {},
	'y': {},
	'z': {},
	'A': {},
	'B': {},
	'C': {},
	'D': {},
	'E': {},
	'F': {},
	'G': {},
	'H': {},
	'I': {},
	'J': {},
	'K': {},
	'L': {},
	'M': {},
	'N': {},
	'O': {},
	'P': {},
	'Q': {},
	'R': {},
	'S': {},
	'T': {},
	'U': {},
	'V': {},
	'W': {},
	'X': {},
	'Y': {},
	'Z': {},
	'0': {},
	'1': {},
	'2': {},
	'3': {},
	'4': {},
	'5': {},
	'6': {},
	'7': {},
	'8': {},
	'9': {},
	'_': {},
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
			for i < len(val) {
				_, exists := ValidIdentLetters[char]
				if !exists {
					break
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
