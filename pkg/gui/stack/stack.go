package stack

type Stack struct {
	stack []string
}

func (s *Stack) Push(contextKey string) {
	s.stack = append(s.stack, contextKey)
}

func (s *Stack) Pop() (string, bool) {
	if len(s.stack) == 0 {
		return "", false
	}

	n := len(s.stack) - 1
	value := s.stack[n]
	s.stack = s.stack[:n]

	return value, true
}
