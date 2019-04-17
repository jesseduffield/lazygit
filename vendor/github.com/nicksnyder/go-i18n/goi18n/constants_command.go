package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/nicksnyder/go-i18n/i18n/bundle"
	"github.com/nicksnyder/go-i18n/i18n/language"
	"github.com/nicksnyder/go-i18n/i18n/translation"
)

type constantsCommand struct {
	translationFiles []string
	packageName      string
	outdir           string
}

type templateConstants struct {
	ID       string
	Name     string
	Comments []string
}

type templateHeader struct {
	PackageName string
	Constants   []templateConstants
}

var constTemplate = template.Must(template.New("").Parse(`// DON'T CHANGE THIS FILE MANUALLY
// This file was generated using the command:
// $ goi18n constants

package {{.PackageName}}
{{range .Constants}}
// {{.Name}} is the identifier for the following localizable string template(s):{{range .Comments}}
// {{.}}{{end}}
const {{.Name}} = "{{.ID}}"
{{end}}`))

func (cc *constantsCommand) execute() error {
	if len(cc.translationFiles) != 1 {
		return fmt.Errorf("need one translation file")
	}

	bundle := bundle.New()

	if err := bundle.LoadTranslationFile(cc.translationFiles[0]); err != nil {
		return fmt.Errorf("failed to load translation file %s because %s\n", cc.translationFiles[0], err)
	}

	translations := bundle.Translations()
	lang := translations[bundle.LanguageTags()[0]]

	// create an array of id to organize
	keys := make([]string, len(lang))
	i := 0

	for id := range lang {
		keys[i] = id
		i++
	}
	sort.Strings(keys)

	tmpl := &templateHeader{
		PackageName: cc.packageName,
		Constants:   make([]templateConstants, len(keys)),
	}

	for i, id := range keys {
		tmpl.Constants[i].ID = id
		tmpl.Constants[i].Name = toCamelCase(id)
		tmpl.Constants[i].Comments = toComments(lang[id])
	}

	filename := filepath.Join(cc.outdir, cc.packageName+".go")
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s because %s", filename, err)
	}

	defer f.Close()

	if err = constTemplate.Execute(f, tmpl); err != nil {
		return fmt.Errorf("failed to write file %s because %s", filename, err)
	}

	return nil
}

func (cc *constantsCommand) parse(arguments []string) {
	flags := flag.NewFlagSet("constants", flag.ExitOnError)
	flags.Usage = usageConstants

	packageName := flags.String("package", "R", "")
	outdir := flags.String("outdir", ".", "")

	flags.Parse(arguments)

	cc.translationFiles = flags.Args()
	cc.packageName = *packageName
	cc.outdir = *outdir
}

func (cc *constantsCommand) SetArgs(args []string) {
	cc.translationFiles = args
}

func usageConstants() {
	fmt.Printf(`Generate constant file from translation file.

Usage:

    goi18n constants [options] [file]

Translation files:

    A translation file contains the strings and translations for a single language.

    Translation file names must have a suffix of a supported format (e.g. .json) and
    contain a valid language tag as defined by RFC 5646 (e.g. en-us, fr, zh-hant, etc.).

Options:

    -package name
        goi18n generates the constant file under the package name.
        Default: R

    -outdir directory
        goi18n writes the constant file to this directory.
        Default: .

`)
}

// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
// https://github.com/golang/lint/blob/master/lint.go
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XSRF":  true,
	"XSS":   true,
}

func toCamelCase(id string) string {
	var result string

	r := regexp.MustCompile(`[\-\.\_\s]`)
	words := r.Split(id, -1)

	for _, w := range words {
		upper := strings.ToUpper(w)
		if commonInitialisms[upper] {
			result += upper
			continue
		}

		if len(w) > 0 {
			u := []rune(w)
			u[0] = unicode.ToUpper(u[0])
			result += string(u)
		}
	}
	return result
}

func toComments(trans translation.Translation) []string {
	var result []string
	data := trans.MarshalInterface().(map[string]interface{})

	t := data["translation"]

	switch v := reflect.ValueOf(t); v.Kind() {
	case reflect.Map:
		for _, k := range []language.Plural{"zero", "one", "two", "few", "many", "other"} {
			vt := v.MapIndex(reflect.ValueOf(k))
			if !vt.IsValid() {
				continue
			}
			result = append(result, string(k)+": "+strconv.Quote(fmt.Sprint(vt.Interface())))
		}
	default:
		result = append(result, strconv.Quote(fmt.Sprint(t)))
	}

	return result
}
