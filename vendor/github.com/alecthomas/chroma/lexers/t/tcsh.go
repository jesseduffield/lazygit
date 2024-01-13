package t

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Tcsh lexer.
var Tcsh = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Tcsh",
		Aliases:   []string{"tcsh", "csh"},
		Filenames: []string{"*.tcsh", "*.csh"},
		MimeTypes: []string{"application/x-csh"},
	},
	tcshRules,
))

func tcshRules() Rules {
	return Rules{
		"root": {
			Include("basic"),
			{`\$\(`, Keyword, Push("paren")},
			{`\$\{#?`, Keyword, Push("curly")},
			{"`", LiteralStringBacktick, Push("backticks")},
			Include("data"),
		},
		"basic": {
			{`\b(if|endif|else|while|then|foreach|case|default|continue|goto|breaksw|end|switch|endsw)\s*\b`, Keyword, nil},
			{`\b(alias|alloc|bg|bindkey|break|builtins|bye|caller|cd|chdir|complete|dirs|echo|echotc|eval|exec|exit|fg|filetest|getxvers|glob|getspath|hashstat|history|hup|inlib|jobs|kill|limit|log|login|logout|ls-F|migrate|newgrp|nice|nohup|notify|onintr|popd|printenv|pushd|rehash|repeat|rootnode|popd|pushd|set|shift|sched|setenv|setpath|settc|setty|setxvers|shift|source|stop|suspend|source|suspend|telltc|time|umask|unalias|uncomplete|unhash|universe|unlimit|unset|unsetenv|ver|wait|warp|watchlog|where|which)\s*\b`, NameBuiltin, nil},
			{`#.*`, Comment, nil},
			{`\\[\w\W]`, LiteralStringEscape, nil},
			{`(\b\w+)(\s*)(=)`, ByGroups(NameVariable, Text, Operator), nil},
			{`[\[\]{}()=]+`, Operator, nil},
			{`<<\s*(\'?)\\?(\w+)[\w\W]+?\2`, LiteralString, nil},
			{`;`, Punctuation, nil},
		},
		"data": {
			{`(?s)"(\\\\|\\[0-7]+|\\.|[^"\\])*"`, LiteralStringDouble, nil},
			{`(?s)'(\\\\|\\[0-7]+|\\.|[^'\\])*'`, LiteralStringSingle, nil},
			{`\s+`, Text, nil},
			{"[^=\\s\\[\\]{}()$\"\\'`\\\\;#]+", Text, nil},
			{`\d+(?= |\Z)`, LiteralNumber, nil},
			{`\$#?(\w+|.)`, NameVariable, nil},
		},
		"curly": {
			{`\}`, Keyword, Pop(1)},
			{`:-`, Keyword, nil},
			{`\w+`, NameVariable, nil},
			{"[^}:\"\\'`$]+", Punctuation, nil},
			{`:`, Punctuation, nil},
			Include("root"),
		},
		"paren": {
			{`\)`, Keyword, Pop(1)},
			Include("root"),
		},
		"backticks": {
			{"`", LiteralStringBacktick, Pop(1)},
			Include("root"),
		},
	}
}
