package b

import (
	"regexp"

	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// TODO(moorereason): can this be factored away?
var bashAnalyserRe = regexp.MustCompile(`(?m)^#!.*/bin/(?:env |)(?:bash|zsh|sh|ksh)`)

// Bash lexer.
var Bash = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Bash",
		Aliases:   []string{"bash", "sh", "ksh", "zsh", "shell"},
		Filenames: []string{"*.sh", "*.ksh", "*.bash", "*.ebuild", "*.eclass", ".env", "*.env", "*.exheres-0", "*.exlib", "*.zsh", "*.zshrc", ".bashrc", "bashrc", ".bash_*", "bash_*", "zshrc", ".zshrc", "PKGBUILD"},
		MimeTypes: []string{"application/x-sh", "application/x-shellscript"},
	},
	bashRules,
).SetAnalyser(func(text string) float32 {
	if bashAnalyserRe.FindString(text) != "" {
		return 1.0
	}
	return 0.0
}))

func bashRules() Rules {
	return Rules{
		"root": {
			Include("basic"),
			{"`", LiteralStringBacktick, Push("backticks")},
			Include("data"),
			Include("interp"),
		},
		"interp": {
			{`\$\(\(`, Keyword, Push("math")},
			{`\$\(`, Keyword, Push("paren")},
			{`\$\{#?`, LiteralStringInterpol, Push("curly")},
			{`\$[a-zA-Z_]\w*`, NameVariable, nil},
			{`\$(?:\d+|[#$?!_*@-])`, NameVariable, nil},
			{`\$`, Text, nil},
		},
		"basic": {
			{`\b(if|fi|else|while|do|done|for|then|return|function|case|select|continue|until|esac|elif)(\s*)\b`, ByGroups(Keyword, Text), nil},
			{"\\b(alias|bg|bind|break|builtin|caller|cd|command|compgen|complete|declare|dirs|disown|echo|enable|eval|exec|exit|export|false|fc|fg|getopts|hash|help|history|jobs|kill|let|local|logout|popd|printf|pushd|pwd|read|readonly|set|shift|shopt|source|suspend|test|time|times|trap|true|type|typeset|ulimit|umask|unalias|unset|wait)(?=[\\s)`])", NameBuiltin, nil},
			{`\A#!.+\n`, CommentPreproc, nil},
			{`#.*(\S|$)`, CommentSingle, nil},
			{`\\[\w\W]`, LiteralStringEscape, nil},
			{`(\b\w+)(\s*)(\+?=)`, ByGroups(NameVariable, Text, Operator), nil},
			{`[\[\]{}()=]`, Operator, nil},
			{`<<<`, Operator, nil},
			{`<<-?\s*(\'?)\\?(\w+)[\w\W]+?\2`, LiteralString, nil},
			{`&&|\|\|`, Operator, nil},
		},
		"data": {
			{`(?s)\$?"(\\\\|\\[0-7]+|\\.|[^"\\$])*"`, LiteralStringDouble, nil},
			{`"`, LiteralStringDouble, Push("string")},
			{`(?s)\$'(\\\\|\\[0-7]+|\\.|[^'\\])*'`, LiteralStringSingle, nil},
			{`(?s)'.*?'`, LiteralStringSingle, nil},
			{`;`, Punctuation, nil},
			{`&`, Punctuation, nil},
			{`\|`, Punctuation, nil},
			{`\s+`, Text, nil},
			{`\d+(?= |$)`, LiteralNumber, nil},
			{"[^=\\s\\[\\]{}()$\"\\'`\\\\<&|;]+", Text, nil},
			{`<`, Text, nil},
		},
		"string": {
			{`"`, LiteralStringDouble, Pop(1)},
			{`(?s)(\\\\|\\[0-7]+|\\.|[^"\\$])+`, LiteralStringDouble, nil},
			Include("interp"),
		},
		"curly": {
			{`\}`, LiteralStringInterpol, Pop(1)},
			{`:-`, Keyword, nil},
			{`\w+`, NameVariable, nil},
			{"[^}:\"\\'`$\\\\]+", Punctuation, nil},
			{`:`, Punctuation, nil},
			Include("root"),
		},
		"paren": {
			{`\)`, Keyword, Pop(1)},
			Include("root"),
		},
		"math": {
			{`\)\)`, Keyword, Pop(1)},
			{`[-+*/%^|&]|\*\*|\|\|`, Operator, nil},
			{`\d+#\d+`, LiteralNumber, nil},
			{`\d+#(?! )`, LiteralNumber, nil},
			{`\d+`, LiteralNumber, nil},
			Include("root"),
		},
		"backticks": {
			{"`", LiteralStringBacktick, Pop(1)},
			Include("root"),
		},
	}
}
