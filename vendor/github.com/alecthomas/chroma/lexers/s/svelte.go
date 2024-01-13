package s

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/h"
	"github.com/alecthomas/chroma/lexers/internal"
	"github.com/alecthomas/chroma/lexers/t"
)

// Svelte lexer.
var Svelte = internal.Register(DelegatingLexer(h.HTML, MustNewLazyLexer(
	&Config{
		Name:      "Svelte",
		Aliases:   []string{"svelte"},
		Filenames: []string{"*.svelte"},
		MimeTypes: []string{"application/x-svelte"},
		DotAll:    true,
	},
	svelteRules,
)))

func svelteRules() Rules {
	return Rules{
		"root": {
			// Let HTML handle the comments, including comments containing script and style tags
			{`<!--`, Other, Push("comment")},
			{
				// Highlight script and style tags based on lang attribute
				// and allow attributes besides lang
				`(<\s*(?:script|style).*?lang\s*=\s*['"])` +
					`(.+?)(['"].*?>)` +
					`(.+?)` +
					`(<\s*/\s*(?:script|style)\s*>)`,
				UsingByGroup(internal.Get, 2, 4, Other, Other, Other, Other, Other),
				nil,
			},
			{
				// Make sure `{` is not inside script or style tags
				`(?<!<\s*(?:script|style)(?:(?!(?:script|style)\s*>).)*?)` +
					`{` +
					`(?!(?:(?!<\s*(?:script|style)).)*?(?:script|style)\s*>)`,
				Punctuation,
				Push("templates"),
			},
			// on:submit|preventDefault
			{`(?<=\s+on:\w+(?:\|\w+)*)\|(?=\w+)`, Operator, nil},
			{`.+?`, Other, nil},
		},
		"comment": {
			{`-->`, Other, Pop(1)},
			{`.+?`, Other, nil},
		},
		"templates": {
			{`}`, Punctuation, Pop(1)},
			// Let TypeScript handle strings and the curly braces inside them
			{`(?<!(?<!\\)\\)(['"` + "`])" + `.*?(?<!(?<!\\)\\)\1`, Using(t.TypeScript), nil},
			// If there is another opening curly brace push to templates again
			{"{", Punctuation, Push("templates")},
			{`@(debug|html)\b`, Keyword, nil},
			{
				`(#await)(\s+)(\w+)(\s+)(then|catch)(\s+)(\w+)`,
				ByGroups(Keyword, Text, Using(t.TypeScript), Text,
					Keyword, Text, Using(t.TypeScript),
				),
				nil,
			},
			{`(#|/)(await|each|if|key)\b`, Keyword, nil},
			{`(:else)(\s+)(if)?\b`, ByGroups(Keyword, Text, Keyword), nil},
			{`:(catch|then)\b`, Keyword, nil},
			{`[^{}]+`, Using(t.TypeScript), nil},
		},
	}
}
