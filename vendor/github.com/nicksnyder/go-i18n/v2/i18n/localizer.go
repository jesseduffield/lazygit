package i18n

import (
	"fmt"

	"text/template"

	"github.com/nicksnyder/go-i18n/v2/internal/plural"
	"golang.org/x/text/language"
)

// Localizer provides Localize and MustLocalize methods that return localized messages.
type Localizer struct {
	// bundle contains the messages that can be returned by the Localizer.
	bundle *Bundle

	// tags is the list of language tags that the Localizer checks
	// in order when localizing a message.
	tags []language.Tag
}

// NewLocalizer returns a new Localizer that looks up messages
// in the bundle according to the language preferences in langs.
// It can parse Accept-Language headers as defined in http://www.ietf.org/rfc/rfc2616.txt.
func NewLocalizer(bundle *Bundle, langs ...string) *Localizer {
	return &Localizer{
		bundle: bundle,
		tags:   parseTags(langs),
	}
}

func parseTags(langs []string) []language.Tag {
	tags := []language.Tag{}
	for _, lang := range langs {
		t, _, err := language.ParseAcceptLanguage(lang)
		if err != nil {
			continue
		}
		tags = append(tags, t...)
	}
	return tags
}

// LocalizeConfig configures a call to the Localize method on Localizer.
type LocalizeConfig struct {
	// MessageID is the id of the message to lookup.
	// This field is ignored if DefaultMessage is set.
	MessageID string

	// TemplateData is the data passed when executing the message's template.
	// If TemplateData is nil and PluralCount is not nil, then the message template
	// will be executed with data that contains the plural count.
	TemplateData interface{}

	// PluralCount determines which plural form of the message is used.
	PluralCount interface{}

	// DefaultMessage is used if the message is not found in any message files.
	DefaultMessage *Message

	// Funcs is used to extend the Go template engine's built in functions
	Funcs template.FuncMap
}

type invalidPluralCountErr struct {
	messageID   string
	pluralCount interface{}
	err         error
}

func (e *invalidPluralCountErr) Error() string {
	return fmt.Sprintf("invalid plural count %#v for message id %q: %s", e.pluralCount, e.messageID, e.err)
}

// MessageNotFoundErr is returned from Localize when a message could not be found.
type MessageNotFoundErr struct {
	messageID string
}

func (e *MessageNotFoundErr) Error() string {
	return fmt.Sprintf("message %q not found", e.messageID)
}

type pluralizeErr struct {
	messageID string
	tag       language.Tag
}

func (e *pluralizeErr) Error() string {
	return fmt.Sprintf("unable to pluralize %q because there no plural rule for %q", e.messageID, e.tag)
}

type messageIDMismatchErr struct {
	messageID        string
	defaultMessageID string
}

func (e *messageIDMismatchErr) Error() string {
	return fmt.Sprintf("message id %q does not match default message id %q", e.messageID, e.defaultMessageID)
}

// Localize returns a localized message.
func (l *Localizer) Localize(lc *LocalizeConfig) (string, error) {
	msg, _, err := l.LocalizeWithTag(lc)
	return msg, err
}

// Localize returns a localized message.
func (l *Localizer) LocalizeMessage(msg *Message) (string, error) {
	return l.Localize(&LocalizeConfig{
		DefaultMessage: msg,
	})
}

// TODO: uncomment this (and the test) when extract has been updated to extract these call sites too.
// Localize returns a localized message.
// func (l *Localizer) LocalizeMessageID(messageID string) (string, error) {
// 	return l.Localize(&LocalizeConfig{
// 		MessageID: messageID,
// 	})
// }

// LocalizeWithTag returns a localized message and the language tag.
// It may return a best effort localized message even if an error happens.
func (l *Localizer) LocalizeWithTag(lc *LocalizeConfig) (string, language.Tag, error) {
	messageID := lc.MessageID
	if lc.DefaultMessage != nil {
		if messageID != "" && messageID != lc.DefaultMessage.ID {
			return "", language.Und, &messageIDMismatchErr{messageID: messageID, defaultMessageID: lc.DefaultMessage.ID}
		}
		messageID = lc.DefaultMessage.ID
	}

	var operands *plural.Operands
	templateData := lc.TemplateData
	if lc.PluralCount != nil {
		var err error
		operands, err = plural.NewOperands(lc.PluralCount)
		if err != nil {
			return "", language.Und, &invalidPluralCountErr{messageID: messageID, pluralCount: lc.PluralCount, err: err}
		}
		if templateData == nil {
			templateData = map[string]interface{}{
				"PluralCount": lc.PluralCount,
			}
		}
	}

	tag, template := l.getTemplate(messageID, lc.DefaultMessage)
	if template == nil {
		return "", language.Und, &MessageNotFoundErr{messageID: messageID}
	}

	pluralForm := l.pluralForm(tag, operands)
	if pluralForm == plural.Invalid {
		return "", language.Und, &pluralizeErr{messageID: messageID, tag: tag}
	}

	msg, err := template.Execute(pluralForm, templateData, lc.Funcs)
	if err != nil {
		// Attempt to fallback to "Other" pluralization in case translations are incomplete.
		if pluralForm != plural.Other {
			msg2, err2 := template.Execute(plural.Other, templateData, lc.Funcs)
			if err2 == nil {
				return msg2, tag, err
			}
		}
		return "", language.Und, err
	}
	return msg, tag, nil
}

func (l *Localizer) getTemplate(id string, defaultMessage *Message) (language.Tag, *MessageTemplate) {
	// Fast path.
	// Optimistically assume this message id is defined in each language.
	fastTag, template := l.matchTemplate(id, defaultMessage, l.bundle.matcher, l.bundle.tags)
	if template != nil {
		return fastTag, template
	}

	if len(l.bundle.tags) <= 1 {
		return l.bundle.defaultLanguage, nil
	}

	// Slow path.
	// We didn't find a translation for the tag suggested by the default matcher
	// so we need to create a new matcher that contains only the tags in the bundle
	// that have this message.
	foundTags := make([]language.Tag, 0, len(l.bundle.messageTemplates)+1)
	foundTags = append(foundTags, l.bundle.defaultLanguage)

	for t, templates := range l.bundle.messageTemplates {
		template := templates[id]
		if template == nil || template.Other == "" {
			continue
		}
		foundTags = append(foundTags, t)
	}

	return l.matchTemplate(id, defaultMessage, language.NewMatcher(foundTags), foundTags)
}

func (l *Localizer) matchTemplate(id string, defaultMessage *Message, matcher language.Matcher, tags []language.Tag) (language.Tag, *MessageTemplate) {
	_, i, _ := matcher.Match(l.tags...)
	tag := tags[i]
	templates := l.bundle.messageTemplates[tag]
	if templates != nil && templates[id] != nil {
		return tag, templates[id]
	}
	if tag == l.bundle.defaultLanguage && defaultMessage != nil {
		return tag, NewMessageTemplate(defaultMessage)
	}
	return tag, nil
}

func (l *Localizer) pluralForm(tag language.Tag, operands *plural.Operands) plural.Form {
	if operands == nil {
		return plural.Other
	}
	pluralRule := l.bundle.pluralRules.Rule(tag)
	if pluralRule == nil {
		return plural.Invalid
	}
	return pluralRule.PluralFormFunc(operands)
}

// MustLocalize is similar to Localize, except it panics if an error happens.
func (l *Localizer) MustLocalize(lc *LocalizeConfig) string {
	localized, err := l.Localize(lc)
	if err != nil {
		panic(err)
	}
	return localized
}
