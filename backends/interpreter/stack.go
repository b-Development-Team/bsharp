package interpreter

type scope struct {
	vals map[int]*Value
}

type stack struct {
	vals []scope
}

func (s *stack) Push() {
	s.vals = append(s.vals, scope{
		vals: make(map[int]*Value),
	})
}

func (s *stack) Pop() {
	s.vals = s.vals[:len(s.vals)-1]
}

func (s *stack) Set(id int, val *Value, redefine bool) {
	if redefine {
		s.vals[len(s.vals)-1].vals[id] = val
	} else {
		for i := len(s.vals) - 1; i >= 0; i-- {
			if _, exists := s.vals[i].vals[id]; exists {
				s.vals[i].vals[id] = val
				return
			}
		}

		// Not found, put in current scope
		s.vals[len(s.vals)-1].vals[id] = val
	}
}

func (s *stack) Get(id int) *Value {
	for i := len(s.vals) - 1; i >= 0; i-- {
		if val, exists := s.vals[i].vals[id]; exists {
			return val
		}
	}

	return nil
}
