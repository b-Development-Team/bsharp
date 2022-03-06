package optimize

type variableInfo struct {
	Value   *Result
	Defines int
}

type scope struct {
	Variables map[int]*variableInfo
}

type Scope struct {
	Scopes []scope
}

func (s *Scope) Push() {
	s.Scopes = append(s.Scopes, scope{
		Variables: make(map[int]*variableInfo),
	})
}

func (s *Scope) Pop() {
	s.Scopes = s.Scopes[:len(s.Scopes)-1]
}

func (s *Scope) SetVal(id int, val *Result) {
	for i, scope := range s.Scopes { // Remove the variable const val from all above scopes
		if i != len(s.Scopes)-1 {
			delete(scope.Variables, id)
		}
	}

	sc := s.Scopes[len(s.Scopes)-1]
	if sc.Variables[id] == nil {
		sc.Variables[id] = &variableInfo{}
	}
	sc.Variables[id].Value = val
}

func (s *Scope) GetVal(id int) *Result {
	var out *Result = nil
	for _, scope := range s.Scopes {
		v, exists := scope.Variables[id]
		if exists && v.Value != nil {
			out = v.Value
		}
	}
	return out
}
