package i18n

import (
	"fmt"

	"text/template"

	"github.com/nicksnyder/go-i18n/v2/internal"
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
	bundle.init()
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

	// Funcs is used to extend the Go template engines built in functions
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

type messageNotFoundErr struct {
	messageID string
}

func (e *messageNotFoundErr) Error() string {
	return fmt.Sprintf("message %q not found", e.messageID)
}

type pluralizeErr struct {
	messageID string
	tag       language.Tag
}

func (e *pluralizeErr) Error() string {
	return fmt.Sprintf("unable to pluralize %q because there no plural rule for %q", e.messageID, e.tag)
}

// Localize returns a localized message.
func (l *Localizer) Localize(lc *LocalizeConfig) (string, error) {
	messageID := lc.MessageID
	if lc.DefaultMessage != nil {
		messageID = lc.DefaultMessage.ID
	}

	var operands *plural.Operands
	templateData := lc.TemplateData
	if lc.PluralCount != nil {
		var err error
		operands, err = plural.NewOperands(lc.PluralCount)
		if err != nil {
			return "", &invalidPluralCountErr{messageID: messageID, pluralCount: lc.PluralCount, err: err}
		}
		if templateData == nil {
			templateData = map[string]interface{}{
				"PluralCount": lc.PluralCount,
			}
		}
	}
	tag, template := l.getTemplate(messageID, lc.DefaultMessage)
	if template == nil {
		return "", &messageNotFoundErr{messageID: messageID}
	}
	pluralForm := l.pluralForm(tag, operands)
	if pluralForm == plural.Invalid {
		return "", &pluralizeErr{messageID: messageID, tag: tag}
	}
	return template.Execute(pluralForm, templateData, lc.Funcs)
}

func (l *Localizer) getTemplate(id string, defaultMessage *Message) (language.Tag, *internal.MessageTemplate) {
	// Fast path.
	// Optimistically assume this message id is defined in each language.
	fastTag, template := l.matchTemplate(id, l.bundle.matcher, l.bundle.tags)
	if template != nil {
		return fastTag, template
	}
	if fastTag == l.bundle.DefaultLanguage {
		if defaultMessage == nil {
			return fastTag, nil
		}
		return fastTag, internal.NewMessageTemplate(defaultMessage)
	}
	if len(l.bundle.tags) > 1 {
		// Slow path.
		// We didn't find a translation for the tag suggested by the default matcher
		// so we need to create a new matcher that contains only the tags in the bundle
		// that have this message.
		foundTags := make([]language.Tag, 0, len(l.bundle.messageTemplates))
		if l.bundle.DefaultLanguage != fastTag {
			foundTags = append(foundTags, l.bundle.DefaultLanguage)
		}
		for t, templates := range l.bundle.messageTemplates {
			if t == fastTag {
				// We already tried this tag in the fast path
				continue
			}
			template := templates[id]
			if template == nil || template.Other == "" {
				continue
			}
			foundTags = append(foundTags, t)
		}
		tag, template := l.matchTemplate(id, language.NewMatcher(foundTags), foundTags)
		if template != nil {
			return tag, template
		}
	}
	if defaultMessage == nil {
		return l.bundle.DefaultLanguage, nil
	}
	return l.bundle.DefaultLanguage, internal.NewMessageTemplate(defaultMessage)
}

func (l *Localizer) matchTemplate(id string, matcher language.Matcher, tags []language.Tag) (language.Tag, *internal.MessageTemplate) {
	_, i, _ := matcher.Match(l.tags...)
	tag := tags[i]
	templates := l.bundle.messageTemplates[tag]
	if templates != nil && templates[id] != nil {
		return tag, templates[id]
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
