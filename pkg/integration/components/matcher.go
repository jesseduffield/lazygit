package components

import (
	"strings"

	"github.com/samber/lo"
)

// for making assertions on string values
type Matcher[T any] struct {
	rules []matcherRule[T]

	// this is printed when there's an error so that it's clear what the context of the assertion is
	prefix string
}

type matcherRule[T any] struct {
	// e.g. "contains 'foo'"
	name string
	// returns a bool that says whether the test passed and if it returns false, it
	// also returns a string of the error message
	testFn func(T) (bool, string)
}

func (self *Matcher[T]) name() string {
	if len(self.rules) == 0 {
		return "anything"
	}

	return strings.Join(
		lo.Map(self.rules, func(rule matcherRule[T], _ int) string { return rule.name }),
		", ",
	)
}

func (self *Matcher[T]) test(value T) (bool, string) {
	for _, rule := range self.rules {
		ok, message := rule.testFn(value)
		if ok {
			continue
		}

		if self.prefix != "" {
			return false, self.prefix + " " + message
		}

		return false, message
	}

	return true, ""
}

func (self *Matcher[T]) appendRule(rule matcherRule[T]) *Matcher[T] {
	self.rules = append(self.rules, rule)

	return self
}

// adds context so that if the matcher test(s) fails, we understand what we were trying to test.
// E.g. prefix: "Unexpected content in view 'files'."
func (self *Matcher[T]) context(prefix string) *Matcher[T] {
	self.prefix = prefix

	return self
}
