package s

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Sparql lexer.
var Sparql = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "SPARQL",
		Aliases:   []string{"sparql"},
		Filenames: []string{"*.rq", "*.sparql"},
		MimeTypes: []string{"application/sparql-query"},
	},
	sparqlRules,
))

func sparqlRules() Rules {
	return Rules{
		"root": {
			{`\s+`, Text, nil},
			{`((?i)select|construct|describe|ask|where|filter|group\s+by|minus|distinct|reduced|from\s+named|from|order\s+by|desc|asc|limit|offset|bindings|load|clear|drop|create|add|move|copy|insert\s+data|delete\s+data|delete\s+where|delete|insert|using\s+named|using|graph|default|named|all|optional|service|silent|bind|union|not\s+in|in|as|having|to|prefix|base)\b`, Keyword, nil},
			{`(a)\b`, Keyword, nil},
			{"(<(?:[^<>\"{}|^`\\\\\\x00-\\x20])*>)", NameLabel, nil},
			{`(_:[_\p{L}\p{N}](?:[-_.\p{L}\p{N}]*[-_\p{L}\p{N}])?)`, NameLabel, nil},
			{`[?$][_\p{L}\p{N}]+`, NameVariable, nil},
			{`([\p{L}][-_.\p{L}\p{N}]*)?(\:)((?:[_:\p{L}\p{N}]|(?:%[0-9A-Fa-f][0-9A-Fa-f])|(?:\\[ _~.\-!$&"()*+,;=/?#@%]))(?:(?:[-_:.\p{L}\p{N}]|(?:%[0-9A-Fa-f][0-9A-Fa-f])|(?:\\[ _~.\-!$&"()*+,;=/?#@%]))*(?:[-_:\p{L}\p{N}]|(?:%[0-9A-Fa-f][0-9A-Fa-f])|(?:\\[ _~.\-!$&"()*+,;=/?#@%])))?)?`, ByGroups(NameNamespace, Punctuation, NameTag), nil},
			{`((?i)str|lang|langmatches|datatype|bound|iri|uri|bnode|rand|abs|ceil|floor|round|concat|strlen|ucase|lcase|encode_for_uri|contains|strstarts|strends|strbefore|strafter|year|month|day|hours|minutes|seconds|timezone|tz|now|md5|sha1|sha256|sha384|sha512|coalesce|if|strlang|strdt|sameterm|isiri|isuri|isblank|isliteral|isnumeric|regex|substr|replace|exists|not\s+exists|count|sum|min|max|avg|sample|group_concat|separator)\b`, NameFunction, nil},
			{`(true|false)`, KeywordConstant, nil},
			{`[+\-]?(\d+\.\d*[eE][+-]?\d+|\.?\d+[eE][+-]?\d+)`, LiteralNumberFloat, nil},
			{`[+\-]?(\d+\.\d*|\.\d+)`, LiteralNumberFloat, nil},
			{`[+\-]?\d+`, LiteralNumberInteger, nil},
			{`(\|\||&&|=|\*|\-|\+|/|!=|<=|>=|!|<|>)`, Operator, nil},
			{`[(){}.;,:^\[\]]`, Punctuation, nil},
			{`#[^\n]*`, Comment, nil},
			{`"""`, LiteralString, Push("triple-double-quoted-string")},
			{`"`, LiteralString, Push("single-double-quoted-string")},
			{`'''`, LiteralString, Push("triple-single-quoted-string")},
			{`'`, LiteralString, Push("single-single-quoted-string")},
		},
		"triple-double-quoted-string": {
			{`"""`, LiteralString, Push("end-of-string")},
			{`[^\\]+`, LiteralString, nil},
			{`\\`, LiteralString, Push("string-escape")},
		},
		"single-double-quoted-string": {
			{`"`, LiteralString, Push("end-of-string")},
			{`[^"\\\n]+`, LiteralString, nil},
			{`\\`, LiteralString, Push("string-escape")},
		},
		"triple-single-quoted-string": {
			{`'''`, LiteralString, Push("end-of-string")},
			{`[^\\]+`, LiteralString, nil},
			{`\\`, LiteralStringEscape, Push("string-escape")},
		},
		"single-single-quoted-string": {
			{`'`, LiteralString, Push("end-of-string")},
			{`[^'\\\n]+`, LiteralString, nil},
			{`\\`, LiteralString, Push("string-escape")},
		},
		"string-escape": {
			{`u[0-9A-Fa-f]{4}`, LiteralStringEscape, Pop(1)},
			{`U[0-9A-Fa-f]{8}`, LiteralStringEscape, Pop(1)},
			{`.`, LiteralStringEscape, Pop(1)},
		},
		"end-of-string": {
			{`(@)([a-zA-Z]+(?:-[a-zA-Z0-9]+)*)`, ByGroups(Operator, NameFunction), Pop(2)},
			{`\^\^`, Operator, Pop(2)},
			Default(Pop(2)),
		},
	}
}
