package optimize

type scope struct {
	variables map[int]*Result
}

type Scope struct {
	scopes []scope
}

func (s *Scope) Push() {
	s.scopes = append(s.scopes, scope{
		variables: make(map[int]*Result),
	})
}

func (s *Scope) Pop() {
	s.scopes = s.scopes[:len(s.scopes)-1]
}

func (s *Scope) SetVar(id int, val *Result) {
	for i, scope := range s.scopes { // Remove the variable const val from all above scopes
		if i != len(s.scopes)-1 {
			delete(scope.variables, id)
		}
	}

	s.scopes[len(s.scopes)-1].variables[id] = val
}

func (s *Scope) GetVar(id int) *Result {
	var out *Result = nil
	for _, scope := range s.scopes {
		v, exists := scope.variables[id]
		if exists {
			out = v
		}
	}
	return out
}
