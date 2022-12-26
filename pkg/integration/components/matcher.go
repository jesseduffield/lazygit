package components

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/samber/lo"
)

// for making assertions on string values
type matcher struct {
	rules []matcherRule

	// this is printed when there's an error so that it's clear what the context of the assertion is
	prefix string
}

type matcherRule struct {
	// e.g. "contains 'foo'"
	name string
	// returns a bool that says whether the test passed and if it returns false, it
	// also returns a string of the error message
	testFn func(string) (bool, string)
}

func NewMatcher(name string, testFn func(string) (bool, string)) *matcher {
	rules := []matcherRule{{name: name, testFn: testFn}}
	return &matcher{rules: rules}
}

func (self *matcher) name() string {
	if len(self.rules) == 0 {
		return "anything"
	}

	return strings.Join(
		lo.Map(self.rules, func(rule matcherRule, _ int) string { return rule.name }),
		", ",
	)
}

func (self *matcher) test(value string) (bool, string) {
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

func (self *matcher) Contains(target string) *matcher {
	return self.appendRule(matcherRule{
		name: fmt.Sprintf("contains '%s'", target),
		testFn: func(value string) (bool, string) {
			return strings.Contains(value, target), fmt.Sprintf("Expected '%s' to be found in '%s'", target, value)
		},
	})
}

func (self *matcher) DoesNotContain(target string) *matcher {
	return self.appendRule(matcherRule{
		name: fmt.Sprintf("does not contain '%s'", target),
		testFn: func(value string) (bool, string) {
			return !strings.Contains(value, target), fmt.Sprintf("Expected '%s' to NOT be found in '%s'", target, value)
		},
	})
}

func (self *matcher) MatchesRegexp(target string) *matcher {
	return self.appendRule(matcherRule{
		name: fmt.Sprintf("matches regular expression '%s'", target),
		testFn: func(value string) (bool, string) {
			matched, err := regexp.MatchString(target, value)
			if err != nil {
				return false, fmt.Sprintf("Unexpected error parsing regular expression '%s': %s", target, err.Error())
			}
			return matched, fmt.Sprintf("Expected '%s' to match regular expression '%s'", value, target)
		},
	})
}

func (self *matcher) Equals(target string) *matcher {
	return self.appendRule(matcherRule{
		name: fmt.Sprintf("equals '%s'", target),
		testFn: func(value string) (bool, string) {
			return target == value, fmt.Sprintf("Expected '%s' to equal '%s'", value, target)
		},
	})
}

const IS_SELECTED_RULE_NAME = "is selected"

// special rule that is only to be used in the TopLines and Lines methods, as a way of
// asserting that a given line is selected.
func (self *matcher) IsSelected() *matcher {
	return self.appendRule(matcherRule{
		name: IS_SELECTED_RULE_NAME,
		testFn: func(value string) (bool, string) {
			panic("Special IsSelected matcher is not supposed to have its testFn method called. This rule should only be used within the .Lines() and .TopLines() method on a ViewAsserter.")
		},
	})
}

func (self *matcher) appendRule(rule matcherRule) *matcher {
	self.rules = append(self.rules, rule)

	return self
}

// adds context so that if the matcher test(s) fails, we understand what we were trying to test.
// E.g. prefix: "Unexpected content in view 'files'."
func (self *matcher) context(prefix string) *matcher {
	self.prefix = prefix

	return self
}

// if the matcher has an `IsSelected` rule, it returns true, along with the matcher after that rule has been removed
func (self *matcher) checkIsSelected() (bool, *matcher) {
	check := lo.ContainsBy(self.rules, func(rule matcherRule) bool { return rule.name == IS_SELECTED_RULE_NAME })

	self.rules = lo.Filter(self.rules, func(rule matcherRule, _ int) bool { return rule.name != IS_SELECTED_RULE_NAME })

	return check, self
}

// this matcher has no rules meaning it always passes the test. Use this
// when you don't care what value you're dealing with.
func Anything() *matcher {
	return &matcher{}
}

func Contains(target string) *matcher {
	return Anything().Contains(target)
}

func DoesNotContain(target string) *matcher {
	return Anything().DoesNotContain(target)
}

func MatchesRegexp(target string) *matcher {
	return Anything().MatchesRegexp(target)
}

func Equals(target string) *matcher {
	return Anything().Equals(target)
}
