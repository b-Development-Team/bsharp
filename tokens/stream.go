package tokens

import "fmt"

type Pos struct {
	File string
	Line int
	Char int
}

func (p *Pos) Error(format string, args ...any) error {
	return fmt.Errorf("%s: %s", p.String(), fmt.Sprintf(format, args...))
}

func (p *Pos) String() string {
	return fmt.Sprintf("%s:%d:%d", p.File, p.Line+1, p.Char+1)
}

func (p *Pos) Dup() *Pos {
	return &Pos{
		File: p.File,
		Line: p.Line,
		Char: p.Char,
	}
}

type Stream struct {
	code []rune
	pos  int

	line int
	char int
	file string
}

func NewStream(file, code string) *Stream {
	return &Stream{
		code: []rune(code),
		file: file,
	}
}

func (s *Stream) Char() rune {
	return s.code[s.pos]
}

func (s *Stream) Peek(off int) rune {
	return s.code[s.pos+off]
}

func (s *Stream) Pos() *Pos {
	return &Pos{
		File: s.file,
		Line: s.line,
		Char: s.char,
	}
}

func (s *Stream) Eat() {
	if s.Char() == '\n' {
		s.line++
		s.char = 0
	} else {
		s.char++
	}
	s.pos++
}

func (s *Stream) HasNext() bool {
	return s.pos < len(s.code)
}
