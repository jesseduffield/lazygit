package terminfo

type stack []interface{}

func (s *stack) push(v interface{}) {
	*s = append(*s, v)
}

func (s *stack) pop() interface{} {
	if len(*s) == 0 {
		return nil
	}
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}

func (s *stack) popInt() int {
	if i, ok := s.pop().(int); ok {
		return i
	}
	return 0
}

func (s *stack) popBool() bool {
	if b, ok := s.pop().(bool); ok {
		return b
	}
	return false
}

func (s *stack) popByte() byte {
	if b, ok := s.pop().(byte); ok {
		return b
	}
	return 0
}

func (s *stack) popString() string {
	if a, ok := s.pop().(string); ok {
		return a
	}
	return ""
}

func (s *stack) reset() {
	*s = (*s)[:0]
}
