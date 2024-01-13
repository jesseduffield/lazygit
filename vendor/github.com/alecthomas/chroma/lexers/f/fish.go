package f

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Fish lexer.
var Fish = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Fish",
		Aliases:   []string{"fish", "fishshell"},
		Filenames: []string{"*.fish", "*.load"},
		MimeTypes: []string{"application/x-fish"},
	},
	fishRules,
))

func fishRules() Rules {
	keywords := []string{
		`begin`, `end`, `if`, `else`, `while`, `break`, `for`, `return`, `function`, `block`,
		`case`, `continue`, `switch`, `not`, `and`, `or`, `set`, `echo`, `exit`, `pwd`, `true`,
		`false`, `cd`, `cdh`, `count`, `test`,
	}
	keywordsPattern := Words(`\b`, `\b`, keywords...)

	builtins := []string{
		`alias`, `bg`, `bind`, `breakpoint`, `builtin`, `argparse`, `abbr`, `string`, `command`,
		`commandline`, `complete`, `contains`, `dirh`, `dirs`, `disown`, `emit`, `eval`, `exec`,
		`fg`, `fish`, `fish_add_path`, `fish_breakpoint_prompt`, `fish_command_not_found`,
		`fish_config`, `fish_git_prompt`, `fish_greeting`, `fish_hg_prompt`, `fish_indent`,
		`fish_is_root_user`, `fish_key_reader`, `fish_mode_prompt`, `fish_opt`, `fish_pager`,
		`fish_prompt`, `fish_right_prompt`, `fish_status_to_signal`, `fish_svn_prompt`,
		`fish_title`, `fish_update_completions`, `fish_vcs_prompt`, `fishd`, `funced`,
		`funcsave`, `functions`, `help`, `history`, `isatty`, `jobs`, `math`, `mimedb`, `nextd`,
		`open`, `prompt_pwd`, `realpath`, `popd`, `prevd`, `psub`, `pushd`, `random`, `read`,
		`set_color`, `source`, `status`, `suspend`, `trap`, `type`, `ulimit`, `umask`, `vared`,
		`fc`, `getopts`, `hash`, `kill`, `printf`, `time`, `wait`,
	}

	return Rules{
		"root": {
			Include("basic"),
			Include("interp"),
			Include("data"),
		},
		"interp": {
			{`\$\(\(`, Keyword, Push("math")},
			{`\(`, Keyword, Push("paren")},
			{`\$#?(\w+|.)`, NameVariable, nil},
		},
		"basic": {
			{Words(`(?<=(?:^|\A|;|&&|\|\||\||`+keywordsPattern+`)\s*)`, `(?=;?\b)`, keywords...), Keyword, nil},
			{`(?<=for\s+\S+\s+)in\b`, Keyword, nil},
			{Words(`\b`, `\s*\b(?!\.)`, builtins...), NameBuiltin, nil},
			{`#!.*\n`, CommentHashbang, nil},
			{`#.*\n`, Comment, nil},
			{`\\[\w\W]`, LiteralStringEscape, nil},
			{`(\b\w+)(\s*)(=)`, ByGroups(NameVariable, Text, Operator), nil},
			{`[\[\]()={}]`, Operator, nil},
			{`(?<=\[[^\]]+)\.\.|-(?=[^\[]+\])`, Operator, nil},
			{`<<-?\s*(\'?)\\?(\w+)[\w\W]+?\2`, LiteralString, nil},
			{`(?<=set\s+(?:--?[^\d\W][\w-]*\s+)?)\w+`, NameVariable, nil},
			{`(?<=for\s+)\w[\w-]*(?=\s+in)`, NameVariable, nil},
			{`(?<=function\s+)\w(?:[^\n])*?(?= *[-\n])`, NameFunction, nil},
			{`(?<=(?:^|\b(?:and|or|sudo)\b|;|\|\||&&|\||\(|(?:\b\w+\s*=\S+\s)) *)\w[\w-]*`, NameFunction, nil},
		},
		"data": {
			{`(?s)\$?"(\\\\|\\[0-7]+|\\.|[^"\\$])*"`, LiteralStringDouble, nil},
			{`"`, LiteralStringDouble, Push("string")},
			{`(?s)\$'(\\\\|\\[0-7]+|\\.|[^'\\])*'`, LiteralStringSingle, nil},
			{`(?s)'.*?'`, LiteralStringSingle, nil},
			{`;`, Punctuation, nil},
			{`&&|\|\||&|\||\^|<|>`, Operator, nil},
			{`\s+`, Text, nil},
			{`\b\d+\b`, LiteralNumber, nil},
			{`(?<=\s+)--?[^\d][\w-]*`, NameAttribute, nil},
			{".+?", Text, nil},
		},
		"string": {
			{`"`, LiteralStringDouble, Pop(1)},
			{`(?s)(\\\\|\\[0-7]+|\\.|[^"\\$])+`, LiteralStringDouble, nil},
			Include("interp"),
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
	}
}
