package r

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// R/S lexer.
var R = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "R",
		Aliases:   []string{"splus", "s", "r"},
		Filenames: []string{"*.S", "*.R", "*.r", ".Rhistory", ".Rprofile", ".Renviron"},
		MimeTypes: []string{"text/S-plus", "text/S", "text/x-r-source", "text/x-r", "text/x-R", "text/x-r-history", "text/x-r-profile"},
	},
	rRules,
))

func rRules() Rules {
	return Rules{
		"comments": {
			{`#.*$`, CommentSingle, nil},
		},
		"valid_name": {
			{"(?:`[^`\\\\]*(?:\\\\.[^`\\\\]*)*`)|(?:(?:[a-zA-z]|[_.][^0-9])[\\w_.]*)", Name, nil},
		},
		"punctuation": {
			{`\[{1,2}|\]{1,2}|\(|\)|;|,`, Punctuation, nil},
		},
		"keywords": {
			{`(if|else|for|while|repeat|in|next|break|return|switch|function)(?![\w.])`, KeywordReserved, nil},
		},
		"operators": {
			{`<<?-|->>?|-|==|<=|>=|<|>|&&?|!=|\|\|?|\?`, Operator, nil},
			{`\*|\+|\^|/|!|%[^%]*%|=|~|\$|@|:{1,3}`, Operator, nil},
		},
		"builtin_symbols": {
			{`(NULL|NA(_(integer|real|complex|character)_)?|letters|LETTERS|Inf|TRUE|FALSE|NaN|pi|\.\.(\.|[0-9]+))(?![\w.])`, KeywordConstant, nil},
			{`(T|F)\b`, NameBuiltinPseudo, nil},
		},
		"numbers": {
			{`0[xX][a-fA-F0-9]+([pP][0-9]+)?[Li]?`, LiteralNumberHex, nil},
			{`[+-]?([0-9]+(\.[0-9]+)?|\.[0-9]+|\.)([eE][+-]?[0-9]+)?[Li]?`, LiteralNumber, nil},
		},
		"statements": {
			Include("comments"),
			{`\s+`, Text, nil},
			{`\'`, LiteralString, Push("string_squote")},
			{`\"`, LiteralString, Push("string_dquote")},
			Include("builtin_symbols"),
			Include("valid_name"),
			Include("numbers"),
			Include("keywords"),
			Include("punctuation"),
			Include("operators"),
		},
		"root": {
			{"((?:`[^`\\\\]*(?:\\\\.[^`\\\\]*)*`)|(?:(?:[a-zA-z]|[_.][^0-9])[\\w_.]*))\\s*(?=\\()", NameFunction, nil},
			Include("statements"),
			{`\{|\}`, Punctuation, nil},
			{`.`, Text, nil},
		},
		"string_squote": {
			{`([^\'\\]|\\.)*\'`, LiteralString, Pop(1)},
		},
		"string_dquote": {
			{`([^"\\]|\\.)*"`, LiteralString, Pop(1)},
		},
	}
}
