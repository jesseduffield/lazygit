package components

import (
	"fmt"
	"regexp"
	"strings"
)

// for making assertions on string values
type matcher struct {
	// e.g. "contains 'foo'"
	name string
	// returns a bool that says whether the test passed and if it returns false, it
	// also returns a string of the error message
	testFn func(string) (bool, string)
	// this is printed when there's an error so that it's clear what the context of the assertion is
	prefix string
}

func NewMatcher(name string, testFn func(string) (bool, string)) *matcher {
	return &matcher{name: name, testFn: testFn}
}

func (self *matcher) test(value string) (bool, string) {
	ok, message := self.testFn(value)
	if ok {
		return true, ""
	}

	if self.prefix != "" {
		return false, self.prefix + " " + message
	}

	return false, message
}

func (self *matcher) context(prefix string) *matcher {
	self.prefix = prefix

	return self
}

func Contains(target string) *matcher {
	return NewMatcher(
		fmt.Sprintf("contains '%s'", target),
		func(value string) (bool, string) {
			return strings.Contains(value, target), fmt.Sprintf("Expected '%s' to be found in '%s'", target, value)
		},
	)
}

func NotContains(target string) *matcher {
	return NewMatcher(
		fmt.Sprintf("does not contain '%s'", target),
		func(value string) (bool, string) {
			return !strings.Contains(value, target), fmt.Sprintf("Expected '%s' to NOT be found in '%s'", target, value)
		},
	)
}

func MatchesRegexp(target string) *matcher {
	return NewMatcher(
		fmt.Sprintf("matches regular expression '%s'", target),
		func(value string) (bool, string) {
			matched, err := regexp.MatchString(target, value)
			if err != nil {
				return false, fmt.Sprintf("Unexpected error parsing regular expression '%s': %s", target, err.Error())
			}
			return matched, fmt.Sprintf("Expected '%s' to match regular expression '%s'", value, target)
		},
	)
}

func Equals(target string) *matcher {
	return NewMatcher(
		fmt.Sprintf("equals '%s'", target),
		func(value string) (bool, string) {
			return target == value, fmt.Sprintf("Expected '%s' to equal '%s'", value, target)
		},
	)
}
