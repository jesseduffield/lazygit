package c

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Cap'N'Proto Proto lexer.
var CapNProto = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Cap'n Proto",
		Aliases:   []string{"capnp"},
		Filenames: []string{"*.capnp"},
		MimeTypes: []string{},
	},
	capNProtoRules,
))

func capNProtoRules() Rules {
	return Rules{
		"root": {
			{`#.*?$`, CommentSingle, nil},
			{`@[0-9a-zA-Z]*`, NameDecorator, nil},
			{`=`, Literal, Push("expression")},
			{`:`, NameClass, Push("type")},
			{`\$`, NameAttribute, Push("annotation")},
			{`(struct|enum|interface|union|import|using|const|annotation|extends|in|of|on|as|with|from|fixed)\b`, Keyword, nil},
			{`[\w.]+`, Name, nil},
			{`[^#@=:$\w]+`, Text, nil},
		},
		"type": {
			{`[^][=;,(){}$]+`, NameClass, nil},
			{`[[(]`, NameClass, Push("parentype")},
			Default(Pop(1)),
		},
		"parentype": {
			{`[^][;()]+`, NameClass, nil},
			{`[[(]`, NameClass, Push()},
			{`[])]`, NameClass, Pop(1)},
			Default(Pop(1)),
		},
		"expression": {
			{`[^][;,(){}$]+`, Literal, nil},
			{`[[(]`, Literal, Push("parenexp")},
			Default(Pop(1)),
		},
		"parenexp": {
			{`[^][;()]+`, Literal, nil},
			{`[[(]`, Literal, Push()},
			{`[])]`, Literal, Pop(1)},
			Default(Pop(1)),
		},
		"annotation": {
			{`[^][;,(){}=:]+`, NameAttribute, nil},
			{`[[(]`, NameAttribute, Push("annexp")},
			Default(Pop(1)),
		},
		"annexp": {
			{`[^][;()]+`, NameAttribute, nil},
			{`[[(]`, NameAttribute, Push()},
			{`[])]`, NameAttribute, Pop(1)},
			Default(Pop(1)),
		},
	}
}
