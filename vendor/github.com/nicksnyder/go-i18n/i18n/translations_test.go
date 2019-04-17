package i18n

import (
	"testing"

	"github.com/nicksnyder/go-i18n/i18n/bundle"
)

var bobMap = map[string]interface{}{"Person": "Bob"}
var bobStruct = struct{ Person string }{Person: "Bob"}

var testCases = []struct {
	id   string
	arg  interface{}
	want string
}{
	{"program_greeting", nil, "Hello world"},
	{"person_greeting", bobMap, "Hello Bob"},
	{"person_greeting", bobStruct, "Hello Bob"},

	{"your_unread_email_count", 0, "You have 0 unread emails."},
	{"your_unread_email_count", 1, "You have 1 unread email."},
	{"your_unread_email_count", 2, "You have 2 unread emails."},
	{"my_height_in_meters", "1.7", "I am 1.7 meters tall."},

	{"person_unread_email_count", []interface{}{0, bobMap}, "Bob has 0 unread emails."},
	{"person_unread_email_count", []interface{}{1, bobMap}, "Bob has 1 unread email."},
	{"person_unread_email_count", []interface{}{2, bobMap}, "Bob has 2 unread emails."},
	{"person_unread_email_count", []interface{}{0, bobStruct}, "Bob has 0 unread emails."},
	{"person_unread_email_count", []interface{}{1, bobStruct}, "Bob has 1 unread email."},
	{"person_unread_email_count", []interface{}{2, bobStruct}, "Bob has 2 unread emails."},

	{"person_unread_email_count_timeframe", []interface{}{3, map[string]interface{}{
		"Person":    "Bob",
		"Timeframe": "0 days",
	}}, "Bob has 3 unread emails in the past 0 days."},
	{"person_unread_email_count_timeframe", []interface{}{3, map[string]interface{}{
		"Person":    "Bob",
		"Timeframe": "1 day",
	}}, "Bob has 3 unread emails in the past 1 day."},
	{"person_unread_email_count_timeframe", []interface{}{3, map[string]interface{}{
		"Person":    "Bob",
		"Timeframe": "2 days",
	}}, "Bob has 3 unread emails in the past 2 days."},
}

func testFile(t *testing.T, path string) {
	b := bundle.New()
	b.MustLoadTranslationFile(path)

	T, err := b.Tfunc("en-US")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		var args []interface{}
		if _, ok := tc.arg.([]interface{}); ok {
			args = tc.arg.([]interface{})
		} else {
			args = []interface{}{tc.arg}
		}

		got := T(tc.id, args...)
		if got != tc.want {
			t.Errorf("got: %v; want: %v", got, tc.want)
		}
	}
}

func TestJSONParse(t *testing.T) {
	testFile(t, "../goi18n/testdata/expected/en-us.all.json")
}

func TestYAMLParse(t *testing.T) {
	testFile(t, "../goi18n/testdata/en-us.yaml")
}

func TestJSONFlatParse(t *testing.T) {
	testFile(t, "../goi18n/testdata/en-us.flat.json")
}

func TestYAMLFlatParse(t *testing.T) {
	testFile(t, "../goi18n/testdata/en-us.flat.yaml")
}

func TestTOMLFlatParse(t *testing.T) {
	testFile(t, "../goi18n/testdata/en-us.flat.toml")
}
