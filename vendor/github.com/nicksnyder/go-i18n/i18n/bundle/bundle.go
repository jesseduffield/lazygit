// Package bundle manages translations for multiple languages.
package bundle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"sync"
	"unicode"

	"github.com/nicksnyder/go-i18n/i18n/language"
	"github.com/nicksnyder/go-i18n/i18n/translation"
	toml "github.com/pelletier/go-toml"
	"gopkg.in/yaml.v2"
)

// TranslateFunc is a copy of i18n.TranslateFunc to avoid a circular dependency.
type TranslateFunc func(translationID string, args ...interface{}) string

// Bundle stores the translations for multiple languages.
type Bundle struct {
	// The primary translations for a language tag and translation id.
	translations map[string]map[string]translation.Translation

	// Translations that can be used when an exact language match is not possible.
	fallbackTranslations map[string]map[string]translation.Translation

	sync.RWMutex
}

// New returns an empty bundle.
func New() *Bundle {
	return &Bundle{
		translations:         make(map[string]map[string]translation.Translation),
		fallbackTranslations: make(map[string]map[string]translation.Translation),
	}
}

// MustLoadTranslationFile is similar to LoadTranslationFile
// except it panics if an error happens.
func (b *Bundle) MustLoadTranslationFile(filename string) {
	if err := b.LoadTranslationFile(filename); err != nil {
		panic(err)
	}
}

// LoadTranslationFile loads the translations from filename into memory.
//
// The language that the translations are associated with is parsed from the filename (e.g. en-US.json).
//
// Generally you should load translation files once during your program's initialization.
func (b *Bundle) LoadTranslationFile(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return b.ParseTranslationFileBytes(filename, buf)
}

// ParseTranslationFileBytes is similar to LoadTranslationFile except it parses the bytes in buf.
//
// It is useful for parsing translation files embedded with go-bindata.
func (b *Bundle) ParseTranslationFileBytes(filename string, buf []byte) error {
	basename := filepath.Base(filename)
	langs := language.Parse(basename)
	switch l := len(langs); {
	case l == 0:
		return fmt.Errorf("no language found in %q", basename)
	case l > 1:
		return fmt.Errorf("multiple languages found in filename %q: %v; expected one", basename, langs)
	}
	translations, err := parseTranslations(filename, buf)
	if err != nil {
		return err
	}
	b.AddTranslation(langs[0], translations...)
	return nil
}

func parseTranslations(filename string, buf []byte) ([]translation.Translation, error) {
	if len(buf) == 0 {
		return []translation.Translation{}, nil
	}

	ext := filepath.Ext(filename)

	// `github.com/pelletier/go-toml` lacks an Unmarshal function,
	// so we should parse TOML separately.
	if ext == ".toml" {
		tree, err := toml.LoadReader(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}

		m := make(map[string]map[string]interface{})
		for k, v := range tree.ToMap() {
			m[k] = v.(map[string]interface{})
		}

		return parseFlatFormat(m)
	}

	// Then parse other formats.
	if isStandardFormat(ext, buf) {
		var standardFormat []map[string]interface{}
		if err := unmarshal(ext, buf, &standardFormat); err != nil {
			return nil, fmt.Errorf("failed to unmarshal %v: %v", filename, err)
		}
		return parseStandardFormat(standardFormat)
	}
	var flatFormat map[string]map[string]interface{}
	if err := unmarshal(ext, buf, &flatFormat); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %v: %v", filename, err)
	}
	return parseFlatFormat(flatFormat)
}

func isStandardFormat(ext string, buf []byte) bool {
	buf = deleteLeadingComments(ext, buf)
	firstRune := rune(buf[0])
	return (ext == ".json" && firstRune == '[') || (ext == ".yaml" && firstRune == '-')
}

// deleteLeadingComments deletes leading newlines and comments in buf.
// It only works for ext == ".yaml".
func deleteLeadingComments(ext string, buf []byte) []byte {
	if ext != ".yaml" {
		return buf
	}

	for {
		buf = bytes.TrimLeftFunc(buf, unicode.IsSpace)
		if buf[0] == '#' {
			buf = deleteLine(buf)
		} else {
			break
		}
	}

	return buf
}

func deleteLine(buf []byte) []byte {
	index := bytes.IndexRune(buf, '\n')
	if index == -1 { // If there is only one line without newline ...
		return nil // ... delete it and return nothing.
	}
	if index == len(buf)-1 { // If there is only one line with newline ...
		return nil // ... do the same as above.
	}
	return buf[index+1:]
}

// unmarshal finds an appropriate unmarshal function for ext
// (extension of filename) and unmarshals buf to out. out must be a pointer.
func unmarshal(ext string, buf []byte, out interface{}) error {
	switch ext {
	case ".json":
		return json.Unmarshal(buf, out)
	case ".yaml":
		return yaml.Unmarshal(buf, out)
	}

	return fmt.Errorf("unsupported file extension %v", ext)
}

func parseStandardFormat(data []map[string]interface{}) ([]translation.Translation, error) {
	translations := make([]translation.Translation, 0, len(data))
	for i, translationData := range data {
		t, err := translation.NewTranslation(translationData)
		if err != nil {
			return nil, fmt.Errorf("unable to parse translation #%d because %s\n%v", i, err, translationData)
		}
		translations = append(translations, t)
	}
	return translations, nil
}

// parseFlatFormat just converts data from flat format to standard format
// and passes it to parseStandardFormat.
//
// Flat format logic:
// key of data must be a string and data[key] must be always map[string]interface{},
// but if there is only "other" key in it then it is non-plural, else plural.
func parseFlatFormat(data map[string]map[string]interface{}) ([]translation.Translation, error) {
	var standardFormatData []map[string]interface{}
	for id, translationData := range data {
		dataObject := make(map[string]interface{})
		dataObject["id"] = id
		if len(translationData) == 1 { // non-plural form
			_, otherExists := translationData["other"]
			if otherExists {
				dataObject["translation"] = translationData["other"]
			}
		} else { // plural form
			dataObject["translation"] = translationData
		}

		standardFormatData = append(standardFormatData, dataObject)
	}

	return parseStandardFormat(standardFormatData)
}

// AddTranslation adds translations for a language.
//
// It is useful if your translations are in a format not supported by LoadTranslationFile.
func (b *Bundle) AddTranslation(lang *language.Language, translations ...translation.Translation) {
	b.Lock()
	defer b.Unlock()
	if b.translations[lang.Tag] == nil {
		b.translations[lang.Tag] = make(map[string]translation.Translation, len(translations))
	}
	currentTranslations := b.translations[lang.Tag]
	for _, newTranslation := range translations {
		if currentTranslation := currentTranslations[newTranslation.ID()]; currentTranslation != nil {
			currentTranslations[newTranslation.ID()] = currentTranslation.Merge(newTranslation)
		} else {
			currentTranslations[newTranslation.ID()] = newTranslation
		}
	}

	// lang can provide translations for less specific language tags.
	for _, tag := range lang.MatchingTags() {
		b.fallbackTranslations[tag] = currentTranslations
	}
}

// Translations returns all translations in the bundle.
func (b *Bundle) Translations() map[string]map[string]translation.Translation {
	t := make(map[string]map[string]translation.Translation)
	b.RLock()
	for tag, translations := range b.translations {
		t[tag] = make(map[string]translation.Translation)
		for id, translation := range translations {
			t[tag][id] = translation
		}
	}
	b.RUnlock()
	return t
}

// LanguageTags returns the tags of all languages that that have been added.
func (b *Bundle) LanguageTags() []string {
	var tags []string
	b.RLock()
	for k := range b.translations {
		tags = append(tags, k)
	}
	b.RUnlock()
	return tags
}

// LanguageTranslationIDs returns the ids of all translations that have been added for a given language.
func (b *Bundle) LanguageTranslationIDs(languageTag string) []string {
	var ids []string
	b.RLock()
	for id := range b.translations[languageTag] {
		ids = append(ids, id)
	}
	b.RUnlock()
	return ids
}

// MustTfunc is similar to Tfunc except it panics if an error happens.
func (b *Bundle) MustTfunc(pref string, prefs ...string) TranslateFunc {
	tfunc, err := b.Tfunc(pref, prefs...)
	if err != nil {
		panic(err)
	}
	return tfunc
}

// MustTfuncAndLanguage is similar to TfuncAndLanguage except it panics if an error happens.
func (b *Bundle) MustTfuncAndLanguage(pref string, prefs ...string) (TranslateFunc, *language.Language) {
	tfunc, language, err := b.TfuncAndLanguage(pref, prefs...)
	if err != nil {
		panic(err)
	}
	return tfunc, language
}

// Tfunc is similar to TfuncAndLanguage except is doesn't return the Language.
func (b *Bundle) Tfunc(pref string, prefs ...string) (TranslateFunc, error) {
	tfunc, _, err := b.TfuncAndLanguage(pref, prefs...)
	return tfunc, err
}

// TfuncAndLanguage returns a TranslateFunc for the first Language that
// has a non-zero number of translations in the bundle.
//
// The returned Language matches the the first language preference that could be satisfied,
// but this may not strictly match the language of the translations used to satisfy that preference.
//
// For example, the user may request "zh". If there are no translations for "zh" but there are translations
// for "zh-cn", then the translations for "zh-cn" will be used but the returned Language will be "zh".
//
// It can parse languages from Accept-Language headers (RFC 2616),
// but it assumes weights are monotonically decreasing.
func (b *Bundle) TfuncAndLanguage(pref string, prefs ...string) (TranslateFunc, *language.Language, error) {
	lang := b.supportedLanguage(pref, prefs...)
	var err error
	if lang == nil {
		err = fmt.Errorf("no supported languages found %#v", append(prefs, pref))
	}
	return func(translationID string, args ...interface{}) string {
		return b.translate(lang, translationID, args...)
	}, lang, err
}

// supportedLanguage returns the first language which
// has a non-zero number of translations in the bundle.
func (b *Bundle) supportedLanguage(pref string, prefs ...string) *language.Language {
	lang := b.translatedLanguage(pref)
	if lang == nil {
		for _, pref := range prefs {
			lang = b.translatedLanguage(pref)
			if lang != nil {
				break
			}
		}
	}
	return lang
}

func (b *Bundle) translatedLanguage(src string) *language.Language {
	langs := language.Parse(src)
	b.RLock()
	defer b.RUnlock()
	for _, lang := range langs {
		if len(b.translations[lang.Tag]) > 0 ||
			len(b.fallbackTranslations[lang.Tag]) > 0 {
			return lang
		}
	}
	return nil
}

func (b *Bundle) translate(lang *language.Language, translationID string, args ...interface{}) string {
	if lang == nil {
		return translationID
	}

	translation := b.translation(lang, translationID)
	if translation == nil {
		return translationID
	}

	var data interface{}
	var count interface{}
	if argc := len(args); argc > 0 {
		if isNumber(args[0]) {
			count = args[0]
			if argc > 1 {
				data = args[1]
			}
		} else {
			data = args[0]
		}
	}

	if count != nil {
		if data == nil {
			data = map[string]interface{}{"Count": count}
		} else {
			dataMap := toMap(data)
			dataMap["Count"] = count
			data = dataMap
		}
	} else {
		dataMap := toMap(data)
		if c, ok := dataMap["Count"]; ok {
			count = c
		}
	}

	p, _ := lang.Plural(count)
	template := translation.Template(p)
	if template == nil {
		return translationID
	}

	s := template.Execute(data)
	if s == "" {
		return translationID
	}
	return s
}

func (b *Bundle) translation(lang *language.Language, translationID string) translation.Translation {
	b.RLock()
	defer b.RUnlock()
	translations := b.translations[lang.Tag]
	if translations == nil {
		translations = b.fallbackTranslations[lang.Tag]
		if translations == nil {
			return nil
		}
	}
	return translations[translationID]
}

func isNumber(n interface{}) bool {
	switch n.(type) {
	case int, int8, int16, int32, int64, string:
		return true
	}
	return false
}

func toMap(input interface{}) map[string]interface{} {
	if data, ok := input.(map[string]interface{}); ok {
		return data
	}
	v := reflect.ValueOf(input)
	switch v.Kind() {
	case reflect.Ptr:
		return toMap(v.Elem().Interface())
	case reflect.Struct:
		return structToMap(v)
	default:
		return nil
	}
}

// Converts the top level of a struct to a map[string]interface{}.
// Code inspired by github.com/fatih/structs.
func structToMap(v reflect.Value) map[string]interface{} {
	out := make(map[string]interface{})
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" {
			// unexported field. skip.
			continue
		}
		out[field.Name] = v.FieldByName(field.Name).Interface()
	}
	return out
}
