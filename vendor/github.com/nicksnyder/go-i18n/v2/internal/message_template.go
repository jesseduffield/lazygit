package internal

import (
	"bytes"

	"text/template"

	"github.com/nicksnyder/go-i18n/v2/internal/plural"
)

// MessageTemplate is an executable template for a message.
type MessageTemplate struct {
	*Message
	PluralTemplates map[plural.Form]*Template
}

// NewMessageTemplate returns a new message template.
func NewMessageTemplate(m *Message) *MessageTemplate {
	pluralTemplates := map[plural.Form]*Template{}
	setPluralTemplate(pluralTemplates, plural.Zero, m.Zero)
	setPluralTemplate(pluralTemplates, plural.One, m.One)
	setPluralTemplate(pluralTemplates, plural.Two, m.Two)
	setPluralTemplate(pluralTemplates, plural.Few, m.Few)
	setPluralTemplate(pluralTemplates, plural.Many, m.Many)
	setPluralTemplate(pluralTemplates, plural.Other, m.Other)
	if len(pluralTemplates) == 0 {
		return nil
	}
	return &MessageTemplate{
		Message:         m,
		PluralTemplates: pluralTemplates,
	}
}

func setPluralTemplate(pluralTemplates map[plural.Form]*Template, pluralForm plural.Form, src string) {
	if src != "" {
		pluralTemplates[pluralForm] = &Template{Src: src}
	}
}

// Execute executes the template for the plural form and template data.
func (mt *MessageTemplate) Execute(pluralForm plural.Form, data interface{}, funcs template.FuncMap) (string, error) {
	t := mt.PluralTemplates[pluralForm]
	if err := t.parse(mt.LeftDelim, mt.RightDelim, funcs); err != nil {
		return "", err
	}
	if t.Template == nil {
		return t.Src, nil
	}
	var buf bytes.Buffer
	if err := t.Template.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
