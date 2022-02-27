package interpreter

type scope struct {
	vars []int
}

type Scope struct {
	scopes []scope
}

func (s *Scope) Push() {
	s.scopes = append(s.scopes, scope{
		vars: make([]int, 0),
	})
}

func (s *Scope) AddVar(vr int) {
	s.scopes[len(s.scopes)-1].vars = append(s.scopes[len(s.scopes)-1].vars, vr)
}

func (s *Scope) Pop() []int {
	val := s.scopes[len(s.scopes)-1].vars
	s.scopes = s.scopes[:len(s.scopes)-1]
	return val
}

func NewScope() *Scope {
	return &Scope{
		scopes: make([]scope, 0),
	}
}
