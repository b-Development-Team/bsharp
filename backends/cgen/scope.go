package cgen

type scope struct {
	function bool
	vals     []string
}

type stack struct {
	vals []scope
}

func (s *stack) Push(function ...bool) {
	if len(function) > 0 {
		s.vals = append(s.vals, scope{function: true})
		return
	}
	s.vals = append(s.vals, scope{})
}

func (s *stack) Pop() {
	s.vals = s.vals[:len(s.vals)-1]
}

func (s *stack) FreeCode() string {
	return JoinCode(s.vals[len(s.vals)-1].vals...)
}

func (s *stack) FnFreeCode() string {
	// Get free code for everything in function
	lastFn := -1
	for i, v := range s.vals {
		if v.function {
			lastFn = i
		}
	}
	out := ""
	for i := lastFn; i < len(s.vals); i++ {
		out = JoinCode(append([]string{out}, s.vals[i].vals...)...)
	}
	return out
}

func (s *stack) Add(val string) {
	s.vals[len(s.vals)-1].vals = append(s.vals[len(s.vals)-1].vals, val)
}
