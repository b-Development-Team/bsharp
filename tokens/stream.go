package tokens

import "fmt"

type Pos struct {
	File    string
	Line    int
	Char    int
	EndLine int
	EndChar int
}

type PosError struct {
	Msg string
	Pos *Pos
}

func (p *PosError) Error() string {
	return fmt.Sprintf("%s: %s", p.Pos.String(), p.Msg)
}

func (p *Pos) Error(format string, args ...any) error {
	return &PosError{
		Msg: fmt.Sprintf(format, args...),
		Pos: p,
	}
}

func (a *Pos) Extend(b *Pos) *Pos {
	return &Pos{
		File:    a.File,
		Line:    a.Line,
		Char:    a.Char,
		EndLine: b.EndLine,
		EndChar: b.EndChar,
	}
}

func (p *Pos) Contains(b *Pos) bool {
	return p.File == b.File && b.Line == p.Line && b.Char > p.Char && b.Char < p.EndChar
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

func (s *Stream) CanPeek(off int) bool {
	return s.pos+off < len(s.code)
}

func (s *Stream) Peek(off int) rune {
	return s.code[s.pos+off]
}

func (s *Stream) Pos() *Pos {
	return &Pos{
		File: s.file,
		Line: s.line,
		Char: s.char,

		EndLine: s.line,
		EndChar: s.char + 1,
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
