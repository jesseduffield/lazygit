package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func TestExtract(t *testing.T) {

	tests := []struct {
		name       string
		file       string
		messages   []*i18n.Message
		activeFile []byte
	}{
		{
			name:     "no translations",
			file:     `package main`,
			messages: nil,
		},
		{
			name: "global declaration",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			var m = &i18n.Message{
				ID: "Plural ID",
			}
			`,
			messages: []*i18n.Message{
				{
					ID: "Plural ID",
				},
			},
		},
		{
			name: "short form id only",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			func main() {
				bundle := &i18n.Bundle{}
				l := i18n.NewLocalizer(bundle, "en")
				l.Localize(&i18n.LocalizeConfig{MessageID: "Plural ID"})
			}
			`,
			messages: []*i18n.Message{
				{
					ID: "Plural ID",
				},
			},
		},
		{
			name: "must short form id only",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			func main() {
				bundle := &i18n.Bundle{}
				l := i18n.NewLocalizer(bundle, "en")
				l.MustLocalize(&i18n.LocalizeConfig{MessageID: "Plural ID"})
			}
			`,
			messages: []*i18n.Message{
				{
					ID: "Plural ID",
				},
			},
		},
		{
			name: "custom package name",
			file: `package main

			import bar "github.com/nicksnyder/go-i18n/v2/i18n"

			func main() {
				_ := &bar.Message{
					ID:          "Plural ID",
				}
			}
			`,
			messages: []*i18n.Message{
				{
					ID: "Plural ID",
				},
			},
		},
		{
			name: "exhaustive plural translation",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			func main() {
				_ := &i18n.Message{
					ID:          "Plural ID",
					Description: "Plural description",
					Zero:        "Zero translation",
					One:         "One translation",
					Two:         "Two translation",
					Few:         "Few translation",
					Many:        "Many translation",
					Other:       "Other translation",
				}
			}
			`,
			messages: []*i18n.Message{
				{
					ID:          "Plural ID",
					Description: "Plural description",
					Zero:        "Zero translation",
					One:         "One translation",
					Two:         "Two translation",
					Few:         "Few translation",
					Many:        "Many translation",
					Other:       "Other translation",
				},
			},
			activeFile: []byte(`["Plural ID"]
description = "Plural description"
few = "Few translation"
many = "Many translation"
one = "One translation"
other = "Other translation"
two = "Two translation"
zero = "Zero translation"
`),
		},
		{
			name: "concat id",
			file: `package main

			import "github.com/nicksnyder/go-i18n/v2/i18n"

			func main() {
				_ := &i18n.Message{
					ID: "Plural" +
						" " +
						"ID",
				}
			}
			`,
			messages: []*i18n.Message{
				{
					ID: "Plural ID",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name+" messages", func(t *testing.T) {
			actualMessages, err := extractMessages([]byte(test.file))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(actualMessages, test.messages) {
				t.Fatalf("file:\n%s\nexpected: %s\n     got: %s", test.file, marshalTest(test.messages), marshalTest(actualMessages))
			}
		})
		t.Run(test.name+" active file", func(t *testing.T) {
			indir := mustTempDir("TestExtractCommandIn")
			defer os.RemoveAll(indir)
			outdir := mustTempDir("TestExtractCommandOut")
			defer os.RemoveAll(outdir)

			inpath := filepath.Join(indir, "file.go")
			if err := ioutil.WriteFile(inpath, []byte(test.file), 0666); err != nil {
				t.Fatal(err)
			}

			if code := testableMain([]string{"extract", "-outdir", outdir, indir}); code != 0 {
				t.Fatalf("expected exit code 0; got %d\n", code)
			}

			files, err := ioutil.ReadDir(outdir)
			if err != nil {
				t.Fatal(err)
			}
			if len(files) != 1 {
				t.Fatalf("expected 1 file; got %#v", files)
			}
			actualFile := files[0]
			expectedName := "active.en.toml"
			if actualFile.Name() != expectedName {
				t.Fatalf("expected %s; got %s", expectedName, actualFile.Name())
			}

			outpath := filepath.Join(outdir, actualFile.Name())
			actual, err := ioutil.ReadFile(outpath)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(actual, test.activeFile) {
				t.Fatalf("\nexpected:\n%s\n\ngot:\n%s", test.activeFile, actual)
			}
		})
	}
}

func TestExtractCommand(t *testing.T) {
	outdir, err := ioutil.TempDir("", "TestExtractCommand")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(outdir)
	if code := testableMain([]string{"extract", "-outdir", outdir, "../example/"}); code != 0 {
		t.Fatalf("expected exit code 0; got %d", code)
	}
	actual, err := ioutil.ReadFile(filepath.Join(outdir, "active.en.toml"))
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte(`HelloPerson = "Hello {{.Name}}"

[MyUnreadEmails]
description = "The number of unread emails I have"
one = "I have {{.PluralCount}} unread email."
other = "I have {{.PluralCount}} unread emails."

[PersonUnreadEmails]
description = "The number of unread emails a person has"
one = "{{.Name}} has {{.UnreadEmailCount}} unread email."
other = "{{.Name}} has {{.UnreadEmailCount}} unread emails."
`)
	if !bytes.Equal(actual, expected) {
		t.Fatalf("files not equal\nactual:\n%s\nexpected:\n%s", actual, expected)
	}
}

func marshalTest(value interface{}) string {
	buf, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(buf)
}
