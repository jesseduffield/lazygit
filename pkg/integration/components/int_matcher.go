package components

import (
	"fmt"
)

type IntMatcher struct {
	*Matcher[int]
}

func (self *IntMatcher) EqualsInt(target int) *IntMatcher {
	self.appendRule(matcherRule[int]{
		name: fmt.Sprintf("equals %d", target),
		testFn: func(value int) (bool, string) {
			return value == target, fmt.Sprintf("Expected %d to equal %d", value, target)
		},
	})

	return self
}

func (self *IntMatcher) GreaterThan(target int) *IntMatcher {
	self.appendRule(matcherRule[int]{
		name: fmt.Sprintf("greater than %d", target),
		testFn: func(value int) (bool, string) {
			return value > target, fmt.Sprintf("Expected %d to greater than %d", value, target)
		},
	})

	return self
}

func (self *IntMatcher) LessThan(target int) *IntMatcher {
	self.appendRule(matcherRule[int]{
		name: fmt.Sprintf("less than %d", target),
		testFn: func(value int) (bool, string) {
			return value < target, fmt.Sprintf("Expected %d to less than %d", value, target)
		},
	})

	return self
}

func AnyInt() *IntMatcher {
	return &IntMatcher{Matcher: &Matcher[int]{}}
}

func EqualsInt(target int) *IntMatcher {
	return AnyInt().EqualsInt(target)
}

func GreaterThan(target int) *IntMatcher {
	return AnyInt().GreaterThan(target)
}

func LessThan(target int) *IntMatcher {
	return AnyInt().LessThan(target)
}
