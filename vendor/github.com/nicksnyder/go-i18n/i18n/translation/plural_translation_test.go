package translation

import (
	"reflect"
	"testing"

	"github.com/nicksnyder/go-i18n/i18n/language"
)

func mustTemplate(t *testing.T, src string) *template {
	tmpl, err := newTemplate(src)
	if err != nil {
		t.Fatal(err)
	}
	return tmpl
}

func pluralTranslationFixture(t *testing.T, id string, pluralCategories ...language.Plural) *pluralTranslation {
	templates := make(map[language.Plural]*template, len(pluralCategories))
	for _, pc := range pluralCategories {
		templates[pc] = mustTemplate(t, string(pc))
	}
	return &pluralTranslation{id, templates}
}

func verifyDeepEqual(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\n%#v\nnot equal to expected value\n%#v", actual, expected)
	}
}

func TestPluralTranslationMerge(t *testing.T) {
	pt := pluralTranslationFixture(t, "id", language.One, language.Other)
	oneTemplate, otherTemplate := pt.templates[language.One], pt.templates[language.Other]

	pt.Merge(pluralTranslationFixture(t, "id"))
	verifyDeepEqual(t, pt.templates, map[language.Plural]*template{
		language.One:   oneTemplate,
		language.Other: otherTemplate,
	})

	pt2 := pluralTranslationFixture(t, "id", language.One, language.Two)
	pt.Merge(pt2)
	verifyDeepEqual(t, pt.templates, map[language.Plural]*template{
		language.One:   pt2.templates[language.One],
		language.Two:   pt2.templates[language.Two],
		language.Other: otherTemplate,
	})
}

/* Test implementations from old idea

func TestCopy(t *testing.T) {
	ls := &LocalizedString{
		ID:          "id",
		Translation: testingTemplate(t, "translation {{.Hello}}"),
		Translations: map[language.Plural]*template{
			language.One:   testingTemplate(t, "plural {{.One}}"),
			language.Other: testingTemplate(t, "plural {{.Other}}"),
		},
	}

	c := ls.Copy()
	delete(c.Translations, language.One)
	if _, ok := ls.Translations[language.One]; !ok {
		t.Errorf("deleting plural translation from copy deleted it from the original")
	}
	c.Translations[language.Two] = testingTemplate(t, "plural {{.Two}}")
	if _, ok := ls.Translations[language.Two]; ok {
		t.Errorf("adding plural translation to copy added it to the original")
	}
}

func TestNormalize(t *testing.T) {
	oneTemplate := testingTemplate(t, "one {{.One}}")
	ls := &LocalizedString{
		Translation: testingTemplate(t, "single {{.Single}}"),
		Translations: map[language.Plural]*template{
			language.One: oneTemplate,
			language.Two: testingTemplate(t, "two {{.Two}}"),
		},
	}
	ls.Normalize(LanguageWithCode("en"))
	if ls.Translation != nil {
		t.Errorf("ls.Translation is %#v; expected nil", ls.Translation)
	}
	if actual := ls.Translations[language.Two]; actual != nil {
		t.Errorf("ls.Translation[language.Two] is %#v; expected nil", actual)
	}
	if actual := ls.Translations[language.One]; actual != oneTemplate {
		t.Errorf("ls.Translations[language.One] is %#v; expected %#v", actual, oneTemplate)
	}
	if _, ok := ls.Translations[language.Other]; !ok {
		t.Errorf("ls.Translations[language.Other] shouldn't be empty")
	}
}

func TestMergeTranslation(t *testing.T) {
	ls := &LocalizedString{}

	translation := testingTemplate(t, "one {{.Hello}}")
	ls.Merge(&LocalizedString{
		Translation: translation,
	})
	if ls.Translation != translation {
		t.Errorf("expected %#v; got %#v", translation, ls.Translation)
	}

	ls.Merge(&LocalizedString{})
	if ls.Translation != translation {
		t.Errorf("expected %#v; got %#v", translation, ls.Translation)
	}

	translation = testingTemplate(t, "two {{.Hello}}")
	ls.Merge(&LocalizedString{
		Translation: translation,
	})
	if ls.Translation != translation {
		t.Errorf("expected %#v; got %#v", translation, ls.Translation)
	}
}

func TestMergeTranslations(t *testing.T) {
	ls := &LocalizedString{}

	oneTemplate := testingTemplate(t, "one {{.One}}")
	otherTemplate := testingTemplate(t, "other {{.Other}}")
	ls.Merge(&LocalizedString{
		Translations: map[language.Plural]*template{
			language.One:   oneTemplate,
			language.Other: otherTemplate,
		},
	})
	if actual := ls.Translations[language.One]; actual != oneTemplate {
		t.Errorf("ls.Translations[language.One] expected %#v; got %#v", oneTemplate, actual)
	}
	if actual := ls.Translations[language.Other]; actual != otherTemplate {
		t.Errorf("ls.Translations[language.Other] expected %#v; got %#v", otherTemplate, actual)
	}

	ls.Merge(&LocalizedString{
		Translations: map[language.Plural]*template{},
	})
	if actual := ls.Translations[language.One]; actual != oneTemplate {
		t.Errorf("ls.Translations[language.One] expected %#v; got %#v", oneTemplate, actual)
	}
	if actual := ls.Translations[language.Other]; actual != otherTemplate {
		t.Errorf("ls.Translations[language.Other] expected %#v; got %#v", otherTemplate, actual)
	}

	twoTemplate := testingTemplate(t, "two {{.Two}}")
	otherTemplate = testingTemplate(t, "second other {{.Other}}")
	ls.Merge(&LocalizedString{
		Translations: map[language.Plural]*template{
			language.Two:   twoTemplate,
			language.Other: otherTemplate,
		},
	})
	if actual := ls.Translations[language.One]; actual != oneTemplate {
		t.Errorf("ls.Translations[language.One] expected %#v; got %#v", oneTemplate, actual)
	}
	if actual := ls.Translations[language.Two]; actual != twoTemplate {
		t.Errorf("ls.Translations[language.Two] expected %#v; got %#v", twoTemplate, actual)
	}
	if actual := ls.Translations[language.Other]; actual != otherTemplate {
		t.Errorf("ls.Translations[language.Other] expected %#v; got %#v", otherTemplate, actual)
	}
}

func TestMissingTranslations(t *testing.T) {
	en := LanguageWithCode("en")

	tests := []struct {
		localizedString *LocalizedString
		language        *Language
		expected        bool
	}{
		{
			&LocalizedString{},
			en,
			true,
		},
		{
			&LocalizedString{Translation: testingTemplate(t, "single {{.Single}}")},
			en,
			false,
		},
		{
			&LocalizedString{
				Translation: testingTemplate(t, "single {{.Single}}"),
				Translations: map[language.Plural]*template{
					language.One: testingTemplate(t, "one {{.One}}"),
				}},
			en,
			true,
		},
		{
			&LocalizedString{Translations: map[language.Plural]*template{
				language.One: testingTemplate(t, "one {{.One}}"),
			}},
			en,
			true,
		},
		{
			&LocalizedString{Translations: map[language.Plural]*template{
				language.One:   nil,
				language.Other: nil,
			}},
			en,
			true,
		},
		{
			&LocalizedString{Translations: map[language.Plural]*template{
				language.One:   testingTemplate(t, ""),
				language.Other: testingTemplate(t, ""),
			}},
			en,
			true,
		},
		{
			&LocalizedString{Translations: map[language.Plural]*template{
				language.One:   testingTemplate(t, "one {{.One}}"),
				language.Other: testingTemplate(t, "other {{.Other}}"),
			}},
			en,
			false,
		},
	}

	for _, tt := range tests {
		if actual := tt.localizedString.MissingTranslations(tt.language); actual != tt.expected {
			t.Errorf("expected %t got %t for %s, %#v",
				tt.expected, actual, tt.language.code, tt.localizedString)
		}
	}
}

func TestHasTranslations(t *testing.T) {
	en := LanguageWithCode("en")

	tests := []struct {
		localizedString *LocalizedString
		language        *Language
		expected        bool
	}{
		{
			&LocalizedString{},
			en,
			false,
		},
		{
			&LocalizedString{Translation: testingTemplate(t, "single {{.Single}}")},
			en,
			true,
		},
		{
			&LocalizedString{
				Translation:  testingTemplate(t, "single {{.Single}}"),
				Translations: map[language.Plural]*template{}},
			en,
			false,
		},
		{
			&LocalizedString{Translations: map[language.Plural]*template{
				language.One: testingTemplate(t, "one {{.One}}"),
			}},
			en,
			true,
		},
		{
			&LocalizedString{Translations: map[language.Plural]*template{
				language.Two: testingTemplate(t, "two {{.Two}}"),
			}},
			en,
			false,
		},
		{
			&LocalizedString{Translations: map[language.Plural]*template{
				language.One: nil,
			}},
			en,
			false,
		},
		{
			&LocalizedString{Translations: map[language.Plural]*template{
				language.One: testingTemplate(t, ""),
			}},
			en,
			false,
		},
	}

	for _, tt := range tests {
		if actual := tt.localizedString.HasTranslations(tt.language); actual != tt.expected {
			t.Errorf("expected %t got %t for %s, %#v",
				tt.expected, actual, tt.language.code, tt.localizedString)
		}
	}
}

func testingTemplate(t *testing.T, src string) *template {
	tmpl, err := newTemplate(src)
	if err != nil {
		t.Fatal(err)
	}
	return tmpl
}
*/
