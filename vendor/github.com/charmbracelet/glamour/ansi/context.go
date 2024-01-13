package ansi

import (
	"html"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

// RenderContext holds the current rendering options and state.
type RenderContext struct {
	options Options

	blockStack *BlockStack
	table      *TableElement

	stripper *bluemonday.Policy
}

// NewRenderContext returns a new RenderContext.
func NewRenderContext(options Options) RenderContext {
	return RenderContext{
		options:    options,
		blockStack: &BlockStack{},
		table:      &TableElement{},
		stripper:   bluemonday.StrictPolicy(),
	}
}

// SanitizeHTML sanitizes HTML content.
func (ctx RenderContext) SanitizeHTML(s string, trimSpaces bool) string {
	s = ctx.stripper.Sanitize(s)
	if trimSpaces {
		s = strings.TrimSpace(s)
	}

	return html.UnescapeString(s)
}
