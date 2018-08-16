package i18n

import (
	"fmt"
	"io/ioutil"

	"github.com/nicksnyder/go-i18n/v2/internal"
	"github.com/nicksnyder/go-i18n/v2/internal/plural"

	"golang.org/x/text/language"
)

// UnmarshalFunc unmarshals data into v.
type UnmarshalFunc = internal.UnmarshalFunc

// Bundle stores a set of messages and pluralization rules.
// Most applications only need a single bundle
// that is initialized early in the application's lifecycle.
type Bundle struct {
	// DefaultLanguage is the default language of the bundle.
	DefaultLanguage language.Tag

	// UnmarshalFuncs is a map of file extensions to UnmarshalFuncs.
	UnmarshalFuncs map[string]UnmarshalFunc

	messageTemplates map[language.Tag]map[string]*internal.MessageTemplate
	pluralRules      plural.Rules
	tags             []language.Tag
	matcher          language.Matcher
}

func (b *Bundle) init() {
	if b.pluralRules == nil {
		b.pluralRules = plural.DefaultRules()
	}
	b.addTag(b.DefaultLanguage)
}

// RegisterUnmarshalFunc registers an UnmarshalFunc for format.
func (b *Bundle) RegisterUnmarshalFunc(format string, unmarshalFunc UnmarshalFunc) {
	if b.UnmarshalFuncs == nil {
		b.UnmarshalFuncs = make(map[string]UnmarshalFunc)
	}
	b.UnmarshalFuncs[format] = unmarshalFunc
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

// MessageFile represents a parsed message file.
type MessageFile = internal.MessageFile

// ParseMessageFileBytes parses the bytes in buf to add translations to the bundle.
//
// The format of the file is everything after the last ".".
//
// The language tag of the file is everything after the second to last "." or after the last path separator, but before the format.
func (b *Bundle) ParseMessageFileBytes(buf []byte, path string) (*MessageFile, error) {
	messageFile, err := internal.ParseMessageFileBytes(buf, path, b.UnmarshalFuncs)
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
	b.init()
	pluralRule := b.pluralRules.Rule(tag)
	if pluralRule == nil {
		return fmt.Errorf("no plural rule registered for %s", tag)
	}
	if b.messageTemplates == nil {
		b.messageTemplates = map[language.Tag]map[string]*internal.MessageTemplate{}
	}
	if b.messageTemplates[tag] == nil {
		b.messageTemplates[tag] = map[string]*internal.MessageTemplate{}
		b.addTag(tag)
	}
	for _, m := range messages {
		b.messageTemplates[tag][m.ID] = internal.NewMessageTemplate(m)
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
