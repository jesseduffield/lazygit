package i18n

import (
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/internal"
	"golang.org/x/text/language"
	yaml "gopkg.in/yaml.v2"
)

var simpleMessage = internal.MustNewMessage(map[string]string{
	"id":    "simple",
	"other": "simple translation",
})

var detailMessage = internal.MustNewMessage(map[string]string{
	"id":          "detail",
	"description": "detail description",
	"other":       "detail translation",
})

var everythingMessage = internal.MustNewMessage(map[string]string{
	"id":          "everything",
	"description": "everything description",
	"zero":        "zero translation",
	"one":         "one translation",
	"two":         "two translation",
	"few":         "few translation",
	"many":        "many translation",
	"other":       "other translation",
})

func TestPseudoLanguage(t *testing.T) {
	bundle := &Bundle{DefaultLanguage: language.English}
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	expected := "simple simple"
	bundle.MustParseMessageFileBytes([]byte(`
# Comment
simple = "simple simple"
`), "en-double.toml")
	localizer := NewLocalizer(bundle, "en-double")
	localized, err := localizer.Localize(&LocalizeConfig{MessageID: "simple"})
	if err != nil {
		t.Fatal(err)
	}
	if localized != expected {
		t.Fatalf("expected %q\ngot %q", expected, localized)
	}
}

func TestPseudoLanguagePlural(t *testing.T) {
	bundle := &Bundle{DefaultLanguage: language.English}
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
[everything]
few = "few translation"
many = "many translation"
one = "one translation"
other = "other translation"
two = "two translation"
zero = "zero translation"
`), "en-double.toml")
	localizer := NewLocalizer(bundle, "en-double")
	{
		expected := "other translation"
		localized, err := localizer.Localize(&LocalizeConfig{MessageID: "everything", PluralCount: 2})
		if err != nil {
			t.Fatal(err)
		}
		if localized != expected {
			t.Fatalf("expected %q\ngot %q", expected, localized)
		}
	}
	{
		expected := "one translation"
		localized, err := localizer.Localize(&LocalizeConfig{MessageID: "everything", PluralCount: 1})
		if err != nil {
			t.Fatal(err)
		}
		if localized != expected {
			t.Fatalf("expected %q\ngot %q", expected, localized)
		}
	}
}

func TestJSON(t *testing.T) {
	var bundle Bundle
	bundle.MustParseMessageFileBytes([]byte(`{
	"simple": "simple translation",
	"detail": {
		"description": "detail description",
		"other": "detail translation"
	},
	"everything": {
		"description": "everything description",
		"zero": "zero translation",
		"one": "one translation",
		"two": "two translation",
		"few": "few translation",
		"many": "many translation",
		"other": "other translation"
	}
}`), "en-US.json")

	expectMessage(t, bundle, language.AmericanEnglish, "simple", simpleMessage)
	expectMessage(t, bundle, language.AmericanEnglish, "detail", detailMessage)
	expectMessage(t, bundle, language.AmericanEnglish, "everything", everythingMessage)
}

func TestYAML(t *testing.T) {
	var bundle Bundle
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
# Comment
simple: simple translation

# Comment
detail:
  description: detail description 
  other: detail translation

# Comment
everything:
  description: everything description
  zero: zero translation
  one: one translation
  two: two translation
  few: few translation
  many: many translation
  other: other translation
`), "en-US.yaml")

	expectMessage(t, bundle, language.AmericanEnglish, "simple", simpleMessage)
	expectMessage(t, bundle, language.AmericanEnglish, "detail", detailMessage)
	expectMessage(t, bundle, language.AmericanEnglish, "everything", everythingMessage)
}

func TestTOML(t *testing.T) {
	var bundle Bundle
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustParseMessageFileBytes([]byte(`
# Comment
simple = "simple translation"

# Comment
[detail]
description = "detail description"
other = "detail translation"

# Comment
[everything]
description = "everything description"
zero = "zero translation"
one = "one translation"
two = "two translation"
few = "few translation"
many = "many translation"
other = "other translation"
`), "en-US.toml")

	expectMessage(t, bundle, language.AmericanEnglish, "simple", simpleMessage)
	expectMessage(t, bundle, language.AmericanEnglish, "detail", detailMessage)
	expectMessage(t, bundle, language.AmericanEnglish, "everything", everythingMessage)
}

func TestV1Format(t *testing.T) {
	var bundle Bundle
	bundle.MustParseMessageFileBytes([]byte(`[
	{
		"id": "simple",
		"translation": "simple translation"
	},
	{
		"id": "everything",
		"translation": {
			"zero": "zero translation",
			"one": "one translation",
			"two": "two translation",
			"few": "few translation",
			"many": "many translation",
			"other": "other translation"
		}
	}
]
`), "en-US.json")

	expectMessage(t, bundle, language.AmericanEnglish, "simple", simpleMessage)
	e := *everythingMessage
	e.Description = ""
	expectMessage(t, bundle, language.AmericanEnglish, "everything", &e)
}

func TestV1FlatFormat(t *testing.T) {
	var bundle Bundle
	bundle.MustParseMessageFileBytes([]byte(`{
	"simple": {
		"other": "simple translation"
	},
	"everything": {
		"zero": "zero translation",
		"one": "one translation",
		"two": "two translation",
		"few": "few translation",
		"many": "many translation",
		"other": "other translation"
	}
}
`), "en-US.json")

	expectMessage(t, bundle, language.AmericanEnglish, "simple", simpleMessage)
	e := *everythingMessage
	e.Description = ""
	expectMessage(t, bundle, language.AmericanEnglish, "everything", &e)
}

func expectMessage(t *testing.T, bundle Bundle, tag language.Tag, messageID string, message *Message) {
	expected := internal.NewMessageTemplate(message)
	actual := bundle.messageTemplates[tag][messageID]
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("bundle.MessageTemplates[%q][%q] = %#v; want %#v", tag, messageID, actual, expected)
	}
}
