package tokens

import "fmt"

type TokenType int

const (
	TokenTypeIdent TokenType = iota
	TokenTypeNumber
	TokenTypeString
	TokenTypeBrack
)

func (t TokenType) String() string {
	return [...]string{"Ident", "Number", "String", "Brack"}[t]
}

type Token struct {
	Typ   TokenType
	Value string
	Pos   *Pos
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%s, \"%s\", %s)", t.Typ.String(), t.Value, t.Pos.String())
}

type Tokenizer struct {
	Tokens []Token
	s      *Stream
}

func NewTokenizer(s *Stream) *Tokenizer {
	return &Tokenizer{
		Tokens: make([]Token, 0),
		s:      s,
	}
}
