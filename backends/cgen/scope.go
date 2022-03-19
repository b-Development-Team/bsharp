package cgen

type scope struct {
	vals []string
}

type stack struct {
	vals []scope
}

func (s *stack) Push() {
	s.vals = append(s.vals, scope{})
}

func (s *stack) Pop() {
	s.vals = s.vals[:len(s.vals)-1]
}

func (s *stack) FreeCode() string {
	return JoinCode(s.vals[len(s.vals)-1].vals...)
}

func (s *stack) Add(val string) {
	s.vals[len(s.vals)-1].vals = append(s.vals[len(s.vals)-1].vals, val)
}
