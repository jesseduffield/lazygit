package i18n_test

import (
	"testing"

	"golang.org/x/text/language"
)

var matcher language.Matcher

func BenchmarkNewMatcher(b *testing.B) {
	langs := []language.Tag{
		language.English,
		language.AmericanEnglish,
		language.BritishEnglish,
		language.Spanish,
		language.EuropeanSpanish,
		language.Portuguese,
		language.French,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matcher = language.NewMatcher(langs)
	}
}

func BenchmarkMatchStrings(b *testing.B) {
	langs := []language.Tag{
		language.English,
		language.AmericanEnglish,
		language.BritishEnglish,
		language.Spanish,
		language.EuropeanSpanish,
		language.Portuguese,
		language.French,
	}
	matcher := language.NewMatcher(langs)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		language.MatchStrings(matcher, "en-US,en;q=0.9")
	}
}
