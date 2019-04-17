package plural

import (
	"testing"

	"golang.org/x/text/language"
)

func TestRules(t *testing.T) {
	expectedRule := &Rule{}

	testCases := []struct {
		name  string
		rules Rules
		tag   language.Tag
		rule  *Rule
	}{
		{
			name: "exact match",
			rules: Rules{
				language.English: expectedRule,
				language.Spanish: &Rule{},
			},
			tag:  language.English,
			rule: expectedRule,
		},
		{
			name: "inexact match",
			rules: Rules{
				language.English: expectedRule,
			},
			tag:  language.AmericanEnglish,
			rule: expectedRule,
		},
		{
			name: "portuguese doesn't match european portuguese",
			rules: Rules{
				language.EuropeanPortuguese: &Rule{},
			},
			tag:  language.Portuguese,
			rule: nil,
		},
		{
			name: "european portuguese preferred",
			rules: Rules{
				language.Portuguese:         &Rule{},
				language.EuropeanPortuguese: expectedRule,
			},
			tag:  language.EuropeanPortuguese,
			rule: expectedRule,
		},
		{
			name: "zh-Hans",
			rules: Rules{
				language.Chinese: expectedRule,
			},
			tag:  language.SimplifiedChinese,
			rule: expectedRule,
		},
		{
			name: "zh-Hant",
			rules: Rules{
				language.Chinese: expectedRule,
			},
			tag:  language.TraditionalChinese,
			rule: expectedRule,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if rule := testCase.rules.Rule(testCase.tag); rule != testCase.rule {
				panic(rule)
			}
		})
	}
}
