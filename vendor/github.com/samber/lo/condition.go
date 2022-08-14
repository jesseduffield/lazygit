package lo

// Ternary is a 1 line if/else statement.
func Ternary[T any](condition bool, ifOutput T, elseOutput T) T {
	if condition {
		return ifOutput
	}

	return elseOutput
}

type ifElse[T any] struct {
	result T
	done   bool
}

// If.
func If[T any](condition bool, result T) *ifElse[T] {
	if condition {
		return &ifElse[T]{result, true}
	}

	var t T
	return &ifElse[T]{t, false}
}

// ElseIf.
func (i *ifElse[T]) ElseIf(condition bool, result T) *ifElse[T] {
	if !i.done && condition {
		i.result = result
		i.done = true
	}

	return i
}

// Else.
func (i *ifElse[T]) Else(result T) T {
	if i.done {
		return i.result
	}

	return result
}

type switchCase[T comparable, R any] struct {
	predicate T
	result    R
	done      bool
}

// Switch is a pure functional switch/case/default statement.
func Switch[T comparable, R any](predicate T) *switchCase[T, R] {
	var result R

	return &switchCase[T, R]{
		predicate,
		result,
		false,
	}
}

// Case.
func (s *switchCase[T, R]) Case(val T, result R) *switchCase[T, R] {
	if !s.done && s.predicate == val {
		s.result = result
		s.done = true
	}

	return s
}

// CaseF.
func (s *switchCase[T, R]) CaseF(val T, cb func() R) *switchCase[T, R] {
	if !s.done && s.predicate == val {
		s.result = cb()
		s.done = true
	}

	return s
}

// Default.
func (s *switchCase[T, R]) Default(result R) R {
	if !s.done {
		s.result = result
	}

	return s.result
}

// DefaultF.
func (s *switchCase[T, R]) DefaultF(cb func() R) R {
	if !s.done {
		s.result = cb()
	}

	return s.result
}
