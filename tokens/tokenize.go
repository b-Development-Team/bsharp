package tokens

import (
	"unicode"
)

func (t *Tokenizer) Tokenize() error {
	for t.s.HasNext() {
	Tokens:
		switch t.s.Char() {
		case '[':
			t.addToken(Token{
				Typ:   TokenTypeLBrack,
				Value: "[",
				Pos:   t.s.Pos(),
			})
			t.s.Eat()

		case ']':
			t.addToken(Token{
				Typ:   TokenTypeRBrack,
				Value: "]",
				Pos:   t.s.Pos(),
			})
			t.s.Eat()

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			t.addNum()

		case '"':
			t.addString()

		case '\'':
			err := t.addByte()
			if err != nil {
				return err
			}

		case '#': // Comment
			t.s.Eat()
			for t.s.HasNext() {
				if t.s.Char() == '\n' {
					t.s.Eat()
					break
				}

				if t.s.Char() == '#' {
					t.s.Eat()
					break
				}

				t.s.Eat()
			}

		case ' ', '\n', '\t', '\r':
			// Just ignore
			t.s.Eat()

		case '-':
			if t.s.CanPeek(1) {
				switch t.s.Peek(1) {
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
					t.addNum()
					break Tokens
				}
			}
			fallthrough

		default:
			if !isLetter(t.s.Char()) {
				return t.s.Pos().Error("unexpected character: %s", string(t.s.Char()))
			}
			t.addIdent()
		}
	}

	return nil
}

func (t *Tokenizer) addNum() {
	pos := t.s.Pos()
	val := ""
	for t.s.HasNext() {
		isNum := true
		switch t.s.Char() {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', '-':
		default:
			isNum = false
		}
		if !isNum {
			break
		}

		val += string(t.s.Char())
		pos = pos.Extend(t.s.Pos())
		t.s.Eat()
	}
	t.addToken(Token{
		Typ:   TokenTypeNumber,
		Value: val,
		Pos:   pos,
	})
}

func (t *Tokenizer) addString() {
	pos := t.s.Pos()
	t.s.Eat() // Eat first '"'

	val := ""
	for t.s.HasNext() {
		if t.s.Char() == '"' {
			pos = pos.Extend(t.s.Pos())
			t.s.Eat()
			break
		}

		if t.s.Char() == '\\' && t.s.CanPeek(1) { // escaped characters
			t.s.Eat()
			switch t.s.Char() {
			case 'n':
				val += "\n"

			case '\\':
				val += "\\"

			case '"':
				val += "\""

			case 't':
				val += "\t"
			}
			t.s.Eat()
			continue
		}

		val += string(t.s.Char())
		t.s.Eat()
	}

	t.addToken(Token{
		Typ:   TokenTypeString,
		Value: val,
		Pos:   pos,
	})
}

func (t *Tokenizer) addByte() error {
	pos := t.s.Pos()
	t.s.Eat() // Eat first '
	if !t.s.CanPeek(2) {
		return pos.Error("expected closing single quote")
	}

	var val byte
	c := t.s.Char()
	if c == '\\' { // Escape code
		t.s.Eat() // Eat backslash
		v := t.s.Char()
		switch v {
		case '\'':
			val = '\''

		case 'n':
			val = '\n'

		case '\\':
			val = '\\'

		case 't':
			val = '\t'
		}
		t.s.Eat()
	} else {
		val = byte(c)
		t.s.Eat()
	}

	if t.s.Char() != '\'' {
		return t.s.Pos().Error("expected closing single quote")
	}

	pos = pos.Extend(t.s.Pos())
	t.s.Eat() // Eat closing quote

	t.addToken(Token{
		Typ:   TokenTypeByte,
		Value: string(val),
		Pos:   pos,
	})

	return nil
}

func isLetter(char rune) bool {
	return char == '+' || char == '-' || char == '*' || char == '/' || char == '^' || char == '%' || char == '=' || char == '!' || char == '<' || char == '>' || char == '{' || char == '}' || char == ',' || char == ':' || unicode.IsLetter(char)
}

func (t *Tokenizer) addIdent() {
	pos := t.s.Pos()
	val := ""
	for t.s.HasNext() {
		if !isLetter(t.s.Char()) {
			break
		}

		val += string(t.s.Char())
		pos = pos.Extend(t.s.Pos())
		t.s.Eat()
	}

	t.addToken(Token{
		Typ:   TokenTypeIdent,
		Value: val,
		Pos:   pos,
	})
}

func (t *Tokenizer) addToken(tok Token) {
	t.Tokens = append(t.Tokens, tok)
}
