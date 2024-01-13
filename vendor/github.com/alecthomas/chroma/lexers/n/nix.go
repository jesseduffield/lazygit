package n

import (
	"strings"

	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Nix lexer.
var Nix = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Nix",
		Aliases:   []string{"nixos", "nix"},
		Filenames: []string{"*.nix"},
		MimeTypes: []string{"text/x-nix"},
	},
	nixRules,
))

func nixRules() Rules {
	// nixb matches right boundary of a nix word. Use it instead of \b.
	const nixb = `(?![a-zA-Z0-9_'-])`

	return Rules{
		"root": {
			Include("keywords"),
			Include("builtins"),
			// "./path" and ".float" literals have to be above "." operator
			Include("literals"),
			Include("operators"),
			{`#.*$`, CommentSingle, nil},
			{`/\*`, CommentMultiline, Push("comment")},
			{`\(`, Punctuation, Push("paren")},
			{`\[`, Punctuation, Push("list")},
			{`"`, StringDouble, Push("qstring")},
			{`''`, StringSingle, Push("istring")},
			{`{`, Punctuation, Push("scope")},
			{`let` + nixb, Keyword, Push("scope")},
			Include("id"),
			Include("space"),
		},
		"keywords": {
			{`import` + nixb, KeywordNamespace, nil},
			{Words(``, nixb, strings.Fields("rec inherit with if then else assert")...), Keyword, nil},
		},
		"builtins": {
			{`throw` + nixb, NameException, nil},
			{Words(``, nixb, strings.Fields("abort baseNameOf builtins currentTime dependencyClosure derivation dirOf fetchTarball filterSource getAttr getEnv hasAttr isNull map removeAttrs toString toXML")...), NameBuiltin, nil},
		},
		"literals": {
			{Words(``, nixb, strings.Fields("true false null")...), NameConstant, nil},
			Include("uri"),
			Include("path"),
			Include("int"),
			Include("float"),
		},
		"operators": {
			{` [/-] `, Operator, nil},
			{`(\.)(\${)`, ByGroups(Operator, StringInterpol), Push("interpol")},
			{`(\?)(\s*)(\${)`, ByGroups(Operator, Text, StringInterpol), Push("interpol")},
			{Words(``, ``, strings.Fields("@ . ? ++ + != ! // == && || -> <= < >= > *")...), Operator, nil},
			{`[;:]`, Punctuation, nil},
		},
		"comment": {
			{`\*/`, CommentMultiline, Pop(1)},
			{`.|\n`, CommentMultiline, nil},
		},
		"paren": {
			{`\)`, Punctuation, Pop(1)},
			Include("root"),
		},
		"list": {
			{`\]`, Punctuation, Pop(1)},
			Include("root"),
		},
		"qstring": {
			{`"`, StringDouble, Pop(1)},
			{`\${`, StringInterpol, Push("interpol")},
			{`\\.`, StringEscape, nil},
			{`.|\n`, StringDouble, nil},
		},
		"istring": {
			{`''\$`, StringEscape, nil},  // "$"
			{`'''`, StringEscape, nil},   // "''"
			{`''\\.`, StringEscape, nil}, // "\."
			{`''`, StringSingle, Pop(1)},
			{`\${`, StringInterpol, Push("interpol")},
			// The next rule is important: "$" escapes any symbol except "{"!
			{`\$.`, StringSingle, nil}, // "$."
			{`.|\n`, StringSingle, nil},
		},
		"scope": {
			{`}:`, Punctuation, Pop(1)},
			{`}`, Punctuation, Pop(1)},
			{`in` + nixb, Keyword, Pop(1)},
			{`\${`, StringInterpol, Push("interpol")},
			Include("root"), // "==" has to be above "="
			{Words(``, ``, strings.Fields("= ? ,")...), Operator, nil},
		},
		"interpol": {
			{`}`, StringInterpol, Pop(1)},
			Include("root"),
		},
		"id": {
			{`[a-zA-Z_][a-zA-Z0-9_'-]*`, Name, nil},
		},
		"uri": {
			{`[a-zA-Z][a-zA-Z0-9+.-]*:[a-zA-Z0-9%/?:@&=+$,_.!~*'-]+`, StringDoc, nil},
		},
		"path": {
			{`[a-zA-Z0-9._+-]*(/[a-zA-Z0-9._+-]+)+`, StringRegex, nil},
			{`~(/[a-zA-Z0-9._+-]+)+/?`, StringRegex, nil},
			{`<[a-zA-Z0-9._+-]+(/[a-zA-Z0-9._+-]+)*>`, StringRegex, nil},
		},
		"int": {
			{`-?[0-9]+` + nixb, NumberInteger, nil},
		},
		"float": {
			{`-?(([1-9][0-9]*\.[0-9]*)|(0?\.[0-9]+))([Ee][+-]?[0-9]+)?` + nixb, NumberFloat, nil},
		},
		"space": {
			{`[ \t\r\n]+`, Text, nil},
		},
	}
}
