package v

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
	. "github.com/alecthomas/chroma/lexers/p" // nolint
)

// Viml lexer.
var Viml = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "VimL",
		Aliases:   []string{"vim"},
		Filenames: []string{"*.vim", ".vimrc", ".exrc", ".gvimrc", "_vimrc", "_exrc", "_gvimrc", "vimrc", "gvimrc"},
		MimeTypes: []string{"text/x-vim"},
	},
	vimlRules,
))

func vimlRules() Rules {
	return Rules{
		"root": {
			{`^([ \t:]*)(py(?:t(?:h(?:o(?:n)?)?)?)?)([ \t]*)(<<)([ \t]*)(.*)((?:\n|.)*)(\6)`, ByGroups(UsingSelf("root"), Keyword, Text, Operator, Text, Text, Using(Python), Text), nil},
			{`^([ \t:]*)(py(?:t(?:h(?:o(?:n)?)?)?)?)([ \t])(.*)`, ByGroups(UsingSelf("root"), Keyword, Text, Using(Python)), nil},
			{`^\s*".*`, Comment, nil},
			{`[ \t]+`, Text, nil},
			{`/(\\\\|\\/|[^\n/])*/`, LiteralStringRegex, nil},
			{`"(\\\\|\\"|[^\n"])*"`, LiteralStringDouble, nil},
			{`'(''|[^\n'])*'`, LiteralStringSingle, nil},
			{`(?<=\s)"[^\-:.%#=*].*`, Comment, nil},
			{`-?\d+`, LiteralNumber, nil},
			{`#[0-9a-f]{6}`, LiteralNumberHex, nil},
			{`^:`, Punctuation, nil},
			{`[()<>+=!|,~-]`, Punctuation, nil},
			{`\b(let|if|else|endif|elseif|fun|function|endfunction)\b`, Keyword, nil},
			{`\b(NONE|bold|italic|underline|dark|light)\b`, NameBuiltin, nil},
			{`\b\w+\b`, NameOther, nil},
			{`.`, Text, nil},
		},
	}
}
