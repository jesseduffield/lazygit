package internal

import (
	"testing"

	"github.com/nicksnyder/go-i18n/v2/internal/plural"
)

func TestMessageTemplate(t *testing.T) {
	mt := NewMessageTemplate(&Message{ID: "HelloWorld", Other: "Hello World"})
	if mt.PluralTemplates[plural.Other].Src != "Hello World" {
		panic(mt.PluralTemplates)
	}
}

func TestNilMessageTemplate(t *testing.T) {
	if mt := NewMessageTemplate(&Message{ID: "HelloWorld"}); mt != nil {
		panic(mt)
	}
}
