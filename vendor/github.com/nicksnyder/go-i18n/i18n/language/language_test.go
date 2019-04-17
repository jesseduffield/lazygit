package language

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		src  string
		lang []*Language
	}{
		{"en", []*Language{{"en", pluralSpecs["en"]}}},
		{"en-US", []*Language{{"en-us", pluralSpecs["en"]}}},
		{"en_US", []*Language{{"en-us", pluralSpecs["en"]}}},
		{"en-GB", []*Language{{"en-gb", pluralSpecs["en"]}}},
		{"zh-CN", []*Language{{"zh-cn", pluralSpecs["zh"]}}},
		{"zh-TW", []*Language{{"zh-tw", pluralSpecs["zh"]}}},
		{"pt-BR", []*Language{{"pt-br", pluralSpecs["pt"]}}},
		{"pt_BR", []*Language{{"pt-br", pluralSpecs["pt"]}}},
		{"pt-PT", []*Language{{"pt-pt", pluralSpecs["pt"]}}},
		{"pt_PT", []*Language{{"pt-pt", pluralSpecs["pt"]}}},
		{"zh-Hans-CN", []*Language{{"zh-hans-cn", pluralSpecs["zh"]}}},
		{"zh-Hant-TW", []*Language{{"zh-hant-tw", pluralSpecs["zh"]}}},
		{"en-US-en-US", []*Language{{"en-us-en-us", pluralSpecs["en"]}}},
		{".en-US..en-US.", []*Language{{"en-us", pluralSpecs["en"]}}},
		{
			"it, xx-zz, xx-ZZ, zh, en-gb;q=0.8, en;q=0.7, es-ES;q=0.6, de-xx",
			[]*Language{
				{"it", pluralSpecs["it"]},
				{"zh", pluralSpecs["zh"]},
				{"en-gb", pluralSpecs["en"]},
				{"en", pluralSpecs["en"]},
				{"es-es", pluralSpecs["es"]},
				{"de-xx", pluralSpecs["de"]},
			},
		},
		{
			"it-qq,xx,xx-zz,xx-ZZ,zh,en-gb;q=0.8,en;q=0.7,es-ES;q=0.6,de-xx",
			[]*Language{
				{"it-qq", pluralSpecs["it"]},
				{"zh", pluralSpecs["zh"]},
				{"en-gb", pluralSpecs["en"]},
				{"en", pluralSpecs["en"]},
				{"es-es", pluralSpecs["es"]},
				{"de-xx", pluralSpecs["de"]},
			},
		},
		{"en.json", []*Language{{"en", pluralSpecs["en"]}}},
		{"en-US.json", []*Language{{"en-us", pluralSpecs["en"]}}},
		{"en-us.json", []*Language{{"en-us", pluralSpecs["en"]}}},
		{"en-xx.json", []*Language{{"en-xx", pluralSpecs["en"]}}},
		{"xx-Yyen-US", nil},
		{"en US", nil},
		{"", nil},
		{"-", nil},
		{"_", nil},
		{"-en", nil},
		{"_en", nil},
		{"-en-", nil},
		{"_en_", nil},
		{"xx", nil},
	}
	for _, test := range tests {
		lang := Parse(test.src)
		if !reflect.DeepEqual(lang, test.lang) {
			t.Errorf("Parse(%q) = %s expected %s", test.src, lang, test.lang)
		}
	}
}

func TestMatchingTags(t *testing.T) {
	tests := []struct {
		lang    *Language
		matches []string
	}{
		{&Language{"zh-hans-cn", nil}, []string{"zh", "zh-hans", "zh-hans-cn"}},
		{&Language{"foo", nil}, []string{"foo"}},
	}
	for _, test := range tests {
		if actual := test.lang.MatchingTags(); !reflect.DeepEqual(test.matches, actual) {
			t.Errorf("matchingTags(%q) = %q expected %q", test.lang.Tag, actual, test.matches)
		}
	}
}
