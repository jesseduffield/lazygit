package lo

// Ternary is a single line if/else statement.
// Take care to avoid dereferencing potentially nil pointers in your A/B expressions, because they are both evaluated. See TernaryF to avoid this problem.
// Play: https://go.dev/play/p/t-D7WBL44h2
func Ternary[T any](condition bool, ifOutput, elseOutput T) T {
	if condition {
		return ifOutput
	}

	return elseOutput
}

// TernaryF is a single line if/else statement whose options are functions.
// Play: https://go.dev/play/p/AO4VW20JoqM
func TernaryF[T any](condition bool, ifFunc, elseFunc func() T) T {
	if condition {
		return ifFunc()
	}

	return elseFunc()
}

type ifElse[T any] struct {
	result T
	done   bool
}

// If is a single line if/else statement.
// Play: https://go.dev/play/p/WSw3ApMxhyW
func If[T any](condition bool, result T) *ifElse[T] { //nolint:revive
	if condition {
		return &ifElse[T]{result, true}
	}

	var t T
	return &ifElse[T]{t, false}
}

// IfF is a single line if/else statement whose options are functions.
// Play: https://go.dev/play/p/WSw3ApMxhyW
func IfF[T any](condition bool, resultF func() T) *ifElse[T] { //nolint:revive
	if condition {
		return &ifElse[T]{resultF(), true}
	}

	var t T
	return &ifElse[T]{t, false}
}

// ElseIf.
// Play: https://go.dev/play/p/WSw3ApMxhyW
func (i *ifElse[T]) ElseIf(condition bool, result T) *ifElse[T] {
	if !i.done && condition {
		i.result = result
		i.done = true
	}

	return i
}

// ElseIfF.
// Play: https://go.dev/play/p/WSw3ApMxhyW
func (i *ifElse[T]) ElseIfF(condition bool, resultF func() T) *ifElse[T] {
	if !i.done && condition {
		i.result = resultF()
		i.done = true
	}

	return i
}

// Else.
// Play: https://go.dev/play/p/WSw3ApMxhyW
func (i *ifElse[T]) Else(result T) T {
	if i.done {
		return i.result
	}

	return result
}

// ElseF.
// Play: https://go.dev/play/p/WSw3ApMxhyW
func (i *ifElse[T]) ElseF(resultF func() T) T {
	if i.done {
		return i.result
	}

	return resultF()
}

type switchCase[T comparable, R any] struct {
	predicate T
	result    R
	done      bool
}

// Switch is a pure functional switch/case/default statement.
// Play: https://go.dev/play/p/TGbKUMAeRUd
func Switch[T comparable, R any](predicate T) *switchCase[T, R] { //nolint:revive
	var result R

	return &switchCase[T, R]{
		predicate,
		result,
		false,
	}
}

// Case.
// Play: https://go.dev/play/p/TGbKUMAeRUd
func (s *switchCase[T, R]) Case(val T, result R) *switchCase[T, R] {
	if !s.done && s.predicate == val {
		s.result = result
		s.done = true
	}

	return s
}

// CaseF.
// Play: https://go.dev/play/p/TGbKUMAeRUd
func (s *switchCase[T, R]) CaseF(val T, callback func() R) *switchCase[T, R] {
	if !s.done && s.predicate == val {
		s.result = callback()
		s.done = true
	}

	return s
}

// Default.
// Play: https://go.dev/play/p/TGbKUMAeRUd
func (s *switchCase[T, R]) Default(result R) R {
	if !s.done {
		s.result = result
	}

	return s.result
}

// DefaultF.
// Play: https://go.dev/play/p/TGbKUMAeRUd
func (s *switchCase[T, R]) DefaultF(callback func() R) R {
	if !s.done {
		s.result = callback()
	}

	return s.result
}
