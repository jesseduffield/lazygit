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
	// if there are no rules, then we pass the test by default
	if len(self.rules) == 0 {
		return true, ""
	}

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
	rule := matcherRule{
		name: fmt.Sprintf("contains '%s'", target),
		testFn: func(value string) (bool, string) {
			return strings.Contains(value, target), fmt.Sprintf("Expected '%s' to be found in '%s'", target, value)
		},
	}

	self.rules = append(self.rules, rule)

	return self
}

func (self *matcher) DoesNotContain(target string) *matcher {
	rule := matcherRule{
		name: fmt.Sprintf("does not contain '%s'", target),
		testFn: func(value string) (bool, string) {
			return !strings.Contains(value, target), fmt.Sprintf("Expected '%s' to NOT be found in '%s'", target, value)
		},
	}

	self.rules = append(self.rules, rule)

	return self
}

func (self *matcher) MatchesRegexp(target string) *matcher {
	rule := matcherRule{
		name: fmt.Sprintf("matches regular expression '%s'", target),
		testFn: func(value string) (bool, string) {
			matched, err := regexp.MatchString(target, value)
			if err != nil {
				return false, fmt.Sprintf("Unexpected error parsing regular expression '%s': %s", target, err.Error())
			}
			return matched, fmt.Sprintf("Expected '%s' to match regular expression '%s'", value, target)
		},
	}

	self.rules = append(self.rules, rule)

	return self
}

func (self *matcher) Equals(target string) *matcher {
	rule := matcherRule{
		name: fmt.Sprintf("equals '%s'", target),
		testFn: func(value string) (bool, string) {
			return target == value, fmt.Sprintf("Expected '%s' to equal '%s'", value, target)
		},
	}

	self.rules = append(self.rules, rule)

	return self
}

func (self *matcher) context(prefix string) *matcher {
	self.prefix = prefix

	return self
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
