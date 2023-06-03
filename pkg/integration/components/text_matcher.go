package components

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/samber/lo"
)

type TextMatcher struct {
	*Matcher[string]
}

func (self *TextMatcher) Contains(target string) *TextMatcher {
	self.appendRule(matcherRule[string]{
		name: fmt.Sprintf("contains '%s'", target),
		testFn: func(value string) (bool, string) {
			// everything contains the empty string so we unconditionally return true here
			if target == "" {
				return true, ""
			}

			return strings.Contains(value, target), fmt.Sprintf("Expected '%s' to be found in '%s'", target, value)
		},
	})

	return self
}

func (self *TextMatcher) DoesNotContain(target string) *TextMatcher {
	self.appendRule(matcherRule[string]{
		name: fmt.Sprintf("does not contain '%s'", target),
		testFn: func(value string) (bool, string) {
			return !strings.Contains(value, target), fmt.Sprintf("Expected '%s' to NOT be found in '%s'", target, value)
		},
	})

	return self
}

func (self *TextMatcher) MatchesRegexp(target string) *TextMatcher {
	self.appendRule(matcherRule[string]{
		name: fmt.Sprintf("matches regular expression '%s'", target),
		testFn: func(value string) (bool, string) {
			matched, err := regexp.MatchString(target, value)
			if err != nil {
				return false, fmt.Sprintf("Unexpected error parsing regular expression '%s': %s", target, err.Error())
			}
			return matched, fmt.Sprintf("Expected '%s' to match regular expression /%s/", value, target)
		},
	})

	return self
}

func (self *TextMatcher) Equals(target string) *TextMatcher {
	self.appendRule(matcherRule[string]{
		name: fmt.Sprintf("equals '%s'", target),
		testFn: func(value string) (bool, string) {
			return target == value, fmt.Sprintf("Expected '%s' to equal '%s'", value, target)
		},
	})

	return self
}

const IS_SELECTED_RULE_NAME = "is selected"

// special rule that is only to be used in the TopLines and Lines methods, as a way of
// asserting that a given line is selected.
func (self *TextMatcher) IsSelected() *TextMatcher {
	self.appendRule(matcherRule[string]{
		name: IS_SELECTED_RULE_NAME,
		testFn: func(value string) (bool, string) {
			panic("Special IsSelected matcher is not supposed to have its testFn method called. This rule should only be used within the .Lines() and .TopLines() method on a ViewAsserter.")
		},
	})

	return self
}

// if the matcher has an `IsSelected` rule, it returns true, along with the matcher after that rule has been removed
func (self *TextMatcher) checkIsSelected() (bool, *TextMatcher) {
	// copying into a new matcher in case we want to re-use the original later
	newMatcher := &TextMatcher{}
	*newMatcher = *self

	check := lo.ContainsBy(newMatcher.rules, func(rule matcherRule[string]) bool { return rule.name == IS_SELECTED_RULE_NAME })

	newMatcher.rules = lo.Filter(newMatcher.rules, func(rule matcherRule[string], _ int) bool { return rule.name != IS_SELECTED_RULE_NAME })

	return check, newMatcher
}

// this matcher has no rules meaning it always passes the test. Use this
// when you don't care what value you're dealing with.
func AnyString() *TextMatcher {
	return &TextMatcher{Matcher: &Matcher[string]{}}
}

func Contains(target string) *TextMatcher {
	return AnyString().Contains(target)
}

func DoesNotContain(target string) *TextMatcher {
	return AnyString().DoesNotContain(target)
}

func MatchesRegexp(target string) *TextMatcher {
	return AnyString().MatchesRegexp(target)
}

func Equals(target string) *TextMatcher {
	return AnyString().Equals(target)
}
