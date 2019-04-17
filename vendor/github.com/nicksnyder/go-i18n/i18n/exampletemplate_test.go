package i18n_test

import (
	"github.com/nicksnyder/go-i18n/i18n"
	"os"
	"text/template"
)

var funcMap = map[string]interface{}{
	"T": i18n.IdentityTfunc,
}

var tmpl = template.Must(template.New("").Funcs(funcMap).Parse(`
{{T "program_greeting"}}
{{T "person_greeting" .}}
{{T "your_unread_email_count" 0}}
{{T "your_unread_email_count" 1}}
{{T "your_unread_email_count" 2}}
{{T "person_unread_email_count" 0 .}}
{{T "person_unread_email_count" 1 .}}
{{T "person_unread_email_count" 2 .}}
`))

func Example_template() {
	i18n.MustLoadTranslationFile("../goi18n/testdata/expected/en-us.all.json")

	T, _ := i18n.Tfunc("en-US")
	tmpl.Funcs(map[string]interface{}{
		"T": T,
	})

	tmpl.Execute(os.Stdout, map[string]interface{}{
		"Person":    "Bob",
		"Timeframe": T("d_days", 1),
	})

	tmpl.Execute(os.Stdout, struct {
		Person    string
		Timeframe string
	}{
		Person:    "Bob",
		Timeframe: T("d_days", 1),
	})

	// Output:
	// Hello world
	// Hello Bob
	// You have 0 unread emails.
	// You have 1 unread email.
	// You have 2 unread emails.
	// Bob has 0 unread emails.
	// Bob has 1 unread email.
	// Bob has 2 unread emails.
	//
	// Hello world
	// Hello Bob
	// You have 0 unread emails.
	// You have 1 unread email.
	// You have 2 unread emails.
	// Bob has 0 unread emails.
	// Bob has 1 unread email.
	// Bob has 2 unread emails.
}
