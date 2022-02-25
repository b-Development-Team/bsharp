package tokens

import "fmt"

type TokenType int

const (
	TokenTypeIdent TokenType = iota
	TokenTypeNumber
	TokenTypeString
	TokenTypeLBrack
	TokenTypeRBrack
)

func (t TokenType) String() string {
	return [...]string{"Ident", "Number", "String", "LBrack", "RBrack"}[t]
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
	pos    int
	s      *Stream
}

func NewTokenizer(s *Stream) *Tokenizer {
	return &Tokenizer{
		Tokens: make([]Token, 0),
		s:      s,
	}
}

func (t *Tokenizer) HasNext() bool {
	return t.pos < len(t.Tokens)
}

func (t *Tokenizer) Tok() Token {
	return t.Tokens[t.pos]
}

func (t *Tokenizer) Eat() {
	t.pos++
}

func (t *Tokenizer) Last() *Pos {
	tok := t.Tokens[len(t.Tokens)-1]
	p := tok.Pos.Dup()
	p.Char += len(tok.Value)
	return p
}
