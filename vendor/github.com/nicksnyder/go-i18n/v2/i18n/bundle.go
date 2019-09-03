package i18n

import (
	"fmt"
	"io/ioutil"

	"github.com/nicksnyder/go-i18n/v2/internal/plural"

	"golang.org/x/text/language"
)

// UnmarshalFunc unmarshals data into v.
type UnmarshalFunc func(data []byte, v interface{}) error

// Bundle stores a set of messages and pluralization rules.
// Most applications only need a single bundle
// that is initialized early in the application's lifecycle.
// It is not goroutine safe to modify the bundle while Localizers
// are reading from it.
type Bundle struct {
	defaultLanguage  language.Tag
	unmarshalFuncs   map[string]UnmarshalFunc
	messageTemplates map[language.Tag]map[string]*MessageTemplate
	pluralRules      plural.Rules
	tags             []language.Tag
	matcher          language.Matcher
}

// artTag is the language tag used for artifical languages
// https://en.wikipedia.org/wiki/Codes_for_constructed_languages
var artTag = language.MustParse("art")

// NewBundle returns a bundle with a default language and a default set of plural rules.
func NewBundle(defaultLanguage language.Tag) *Bundle {
	b := &Bundle{
		defaultLanguage: defaultLanguage,
		pluralRules:     plural.DefaultRules(),
	}
	b.pluralRules[artTag] = b.pluralRules.Rule(language.English)
	b.addTag(defaultLanguage)
	return b
}

// RegisterUnmarshalFunc registers an UnmarshalFunc for format.
func (b *Bundle) RegisterUnmarshalFunc(format string, unmarshalFunc UnmarshalFunc) {
	if b.unmarshalFuncs == nil {
		b.unmarshalFuncs = make(map[string]UnmarshalFunc)
	}
	b.unmarshalFuncs[format] = unmarshalFunc
}

// LoadMessageFile loads the bytes from path
// and then calls ParseMessageFileBytes.
func (b *Bundle) LoadMessageFile(path string) (*MessageFile, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return b.ParseMessageFileBytes(buf, path)
}

// MustLoadMessageFile is similar to LoadTranslationFile
// except it panics if an error happens.
func (b *Bundle) MustLoadMessageFile(path string) {
	if _, err := b.LoadMessageFile(path); err != nil {
		panic(err)
	}
}

// ParseMessageFileBytes parses the bytes in buf to add translations to the bundle.
//
// The format of the file is everything after the last ".".
//
// The language tag of the file is everything after the second to last "." or after the last path separator, but before the format.
func (b *Bundle) ParseMessageFileBytes(buf []byte, path string) (*MessageFile, error) {
	messageFile, err := ParseMessageFileBytes(buf, path, b.unmarshalFuncs)
	if err != nil {
		return nil, err
	}
	if err := b.AddMessages(messageFile.Tag, messageFile.Messages...); err != nil {
		return nil, err
	}
	return messageFile, nil
}

// MustParseMessageFileBytes is similar to ParseMessageFileBytes
// except it panics if an error happens.
func (b *Bundle) MustParseMessageFileBytes(buf []byte, path string) {
	if _, err := b.ParseMessageFileBytes(buf, path); err != nil {
		panic(err)
	}
}

// AddMessages adds messages for a language.
// It is useful if your messages are in a format not supported by ParseMessageFileBytes.
func (b *Bundle) AddMessages(tag language.Tag, messages ...*Message) error {
	pluralRule := b.pluralRules.Rule(tag)
	if pluralRule == nil {
		return fmt.Errorf("no plural rule registered for %s", tag)
	}
	if b.messageTemplates == nil {
		b.messageTemplates = map[language.Tag]map[string]*MessageTemplate{}
	}
	if b.messageTemplates[tag] == nil {
		b.messageTemplates[tag] = map[string]*MessageTemplate{}
		b.addTag(tag)
	}
	for _, m := range messages {
		b.messageTemplates[tag][m.ID] = NewMessageTemplate(m)
	}
	return nil
}

// MustAddMessages is similar to AddMessages except it panics if an error happens.
func (b *Bundle) MustAddMessages(tag language.Tag, messages ...*Message) {
	if err := b.AddMessages(tag, messages...); err != nil {
		panic(err)
	}
}

func (b *Bundle) addTag(tag language.Tag) {
	for _, t := range b.tags {
		if t == tag {
			// Tag already exists
			return
		}
	}
	b.tags = append(b.tags, tag)
	b.matcher = language.NewMatcher(b.tags)
}

// LanguageTags returns the list of language tags
// of all the translations loaded into the bundle
func (b *Bundle) LanguageTags() []language.Tag {
	return b.tags
}
