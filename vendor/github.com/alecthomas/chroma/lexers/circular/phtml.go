package circular

import (
	"strings"

	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/h"
	"github.com/alecthomas/chroma/lexers/internal"
)

// PHTML lexer is PHP in HTML.
var PHTML = internal.Register(DelegatingLexer(h.HTML, MustNewLazyLexer(
	&Config{
		Name:            "PHTML",
		Aliases:         []string{"phtml"},
		Filenames:       []string{"*.phtml", "*.php", "*.php[345]", "*.inc"},
		MimeTypes:       []string{"application/x-php", "application/x-httpd-php", "application/x-httpd-php3", "application/x-httpd-php4", "application/x-httpd-php5", "text/x-php"},
		DotAll:          true,
		CaseInsensitive: true,
		EnsureNL:        true,
		Priority:        2,
	},
	phtmlRules,
).SetAnalyser(func(text string) float32 {
	if strings.Contains(text, "<?php") {
		return 0.5
	}
	return 0.0
})))

func phtmlRules() Rules {
	return Rules{
		"root": {
			{`<\?(php)?`, CommentPreproc, Push("php")},
			{`[^<]+`, Other, nil},
			{`<`, Other, nil},
		},
	}.Merge(phpCommonRules())
}
