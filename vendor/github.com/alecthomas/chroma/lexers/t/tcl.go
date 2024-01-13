package t

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Tcl lexer.
var Tcl = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Tcl",
		Aliases:   []string{"tcl"},
		Filenames: []string{"*.tcl", "*.rvt"},
		MimeTypes: []string{"text/x-tcl", "text/x-script.tcl", "application/x-tcl"},
	},
	tclRules,
))

func tclRules() Rules {
	return Rules{
		"root": {
			Include("command"),
			Include("basic"),
			Include("data"),
			{`\}`, Keyword, nil},
		},
		"command": {
			{Words(`\b`, `\b`, `after`, `apply`, `array`, `break`, `catch`, `continue`, `elseif`, `else`, `error`, `eval`, `expr`, `for`, `foreach`, `global`, `if`, `namespace`, `proc`, `rename`, `return`, `set`, `switch`, `then`, `trace`, `unset`, `update`, `uplevel`, `upvar`, `variable`, `vwait`, `while`), Keyword, Push("params")},
			{Words(`\b`, `\b`, `append`, `bgerror`, `binary`, `cd`, `chan`, `clock`, `close`, `concat`, `dde`, `dict`, `encoding`, `eof`, `exec`, `exit`, `fblocked`, `fconfigure`, `fcopy`, `file`, `fileevent`, `flush`, `format`, `gets`, `glob`, `history`, `http`, `incr`, `info`, `interp`, `join`, `lappend`, `lassign`, `lindex`, `linsert`, `list`, `llength`, `load`, `loadTk`, `lrange`, `lrepeat`, `lreplace`, `lreverse`, `lsearch`, `lset`, `lsort`, `mathfunc`, `mathop`, `memory`, `msgcat`, `open`, `package`, `pid`, `pkg::create`, `pkg_mkIndex`, `platform`, `platform::shell`, `puts`, `pwd`, `re_syntax`, `read`, `refchan`, `regexp`, `registry`, `regsub`, `scan`, `seek`, `socket`, `source`, `split`, `string`, `subst`, `tell`, `time`, `tm`, `unknown`, `unload`), NameBuiltin, Push("params")},
			{`([\w.-]+)`, NameVariable, Push("params")},
			{`#`, Comment, Push("comment")},
		},
		"command-in-brace": {
			{Words(`\b`, `\b`, `after`, `apply`, `array`, `break`, `catch`, `continue`, `elseif`, `else`, `error`, `eval`, `expr`, `for`, `foreach`, `global`, `if`, `namespace`, `proc`, `rename`, `return`, `set`, `switch`, `then`, `trace`, `unset`, `update`, `uplevel`, `upvar`, `variable`, `vwait`, `while`), Keyword, Push("params-in-brace")},
			{Words(`\b`, `\b`, `append`, `bgerror`, `binary`, `cd`, `chan`, `clock`, `close`, `concat`, `dde`, `dict`, `encoding`, `eof`, `exec`, `exit`, `fblocked`, `fconfigure`, `fcopy`, `file`, `fileevent`, `flush`, `format`, `gets`, `glob`, `history`, `http`, `incr`, `info`, `interp`, `join`, `lappend`, `lassign`, `lindex`, `linsert`, `list`, `llength`, `load`, `loadTk`, `lrange`, `lrepeat`, `lreplace`, `lreverse`, `lsearch`, `lset`, `lsort`, `mathfunc`, `mathop`, `memory`, `msgcat`, `open`, `package`, `pid`, `pkg::create`, `pkg_mkIndex`, `platform`, `platform::shell`, `puts`, `pwd`, `re_syntax`, `read`, `refchan`, `regexp`, `registry`, `regsub`, `scan`, `seek`, `socket`, `source`, `split`, `string`, `subst`, `tell`, `time`, `tm`, `unknown`, `unload`), NameBuiltin, Push("params-in-brace")},
			{`([\w.-]+)`, NameVariable, Push("params-in-brace")},
			{`#`, Comment, Push("comment")},
		},
		"command-in-bracket": {
			{Words(`\b`, `\b`, `after`, `apply`, `array`, `break`, `catch`, `continue`, `elseif`, `else`, `error`, `eval`, `expr`, `for`, `foreach`, `global`, `if`, `namespace`, `proc`, `rename`, `return`, `set`, `switch`, `then`, `trace`, `unset`, `update`, `uplevel`, `upvar`, `variable`, `vwait`, `while`), Keyword, Push("params-in-bracket")},
			{Words(`\b`, `\b`, `append`, `bgerror`, `binary`, `cd`, `chan`, `clock`, `close`, `concat`, `dde`, `dict`, `encoding`, `eof`, `exec`, `exit`, `fblocked`, `fconfigure`, `fcopy`, `file`, `fileevent`, `flush`, `format`, `gets`, `glob`, `history`, `http`, `incr`, `info`, `interp`, `join`, `lappend`, `lassign`, `lindex`, `linsert`, `list`, `llength`, `load`, `loadTk`, `lrange`, `lrepeat`, `lreplace`, `lreverse`, `lsearch`, `lset`, `lsort`, `mathfunc`, `mathop`, `memory`, `msgcat`, `open`, `package`, `pid`, `pkg::create`, `pkg_mkIndex`, `platform`, `platform::shell`, `puts`, `pwd`, `re_syntax`, `read`, `refchan`, `regexp`, `registry`, `regsub`, `scan`, `seek`, `socket`, `source`, `split`, `string`, `subst`, `tell`, `time`, `tm`, `unknown`, `unload`), NameBuiltin, Push("params-in-bracket")},
			{`([\w.-]+)`, NameVariable, Push("params-in-bracket")},
			{`#`, Comment, Push("comment")},
		},
		"command-in-paren": {
			{Words(`\b`, `\b`, `after`, `apply`, `array`, `break`, `catch`, `continue`, `elseif`, `else`, `error`, `eval`, `expr`, `for`, `foreach`, `global`, `if`, `namespace`, `proc`, `rename`, `return`, `set`, `switch`, `then`, `trace`, `unset`, `update`, `uplevel`, `upvar`, `variable`, `vwait`, `while`), Keyword, Push("params-in-paren")},
			{Words(`\b`, `\b`, `append`, `bgerror`, `binary`, `cd`, `chan`, `clock`, `close`, `concat`, `dde`, `dict`, `encoding`, `eof`, `exec`, `exit`, `fblocked`, `fconfigure`, `fcopy`, `file`, `fileevent`, `flush`, `format`, `gets`, `glob`, `history`, `http`, `incr`, `info`, `interp`, `join`, `lappend`, `lassign`, `lindex`, `linsert`, `list`, `llength`, `load`, `loadTk`, `lrange`, `lrepeat`, `lreplace`, `lreverse`, `lsearch`, `lset`, `lsort`, `mathfunc`, `mathop`, `memory`, `msgcat`, `open`, `package`, `pid`, `pkg::create`, `pkg_mkIndex`, `platform`, `platform::shell`, `puts`, `pwd`, `re_syntax`, `read`, `refchan`, `regexp`, `registry`, `regsub`, `scan`, `seek`, `socket`, `source`, `split`, `string`, `subst`, `tell`, `time`, `tm`, `unknown`, `unload`), NameBuiltin, Push("params-in-paren")},
			{`([\w.-]+)`, NameVariable, Push("params-in-paren")},
			{`#`, Comment, Push("comment")},
		},
		"basic": {
			{`\(`, Keyword, Push("paren")},
			{`\[`, Keyword, Push("bracket")},
			{`\{`, Keyword, Push("brace")},
			{`"`, LiteralStringDouble, Push("string")},
			{`(eq|ne|in|ni)\b`, OperatorWord, nil},
			{`!=|==|<<|>>|<=|>=|&&|\|\||\*\*|[-+~!*/%<>&^|?:]`, Operator, nil},
		},
		"data": {
			{`\s+`, Text, nil},
			{`0x[a-fA-F0-9]+`, LiteralNumberHex, nil},
			{`0[0-7]+`, LiteralNumberOct, nil},
			{`\d+\.\d+`, LiteralNumberFloat, nil},
			{`\d+`, LiteralNumberInteger, nil},
			{`\$([\w.:-]+)`, NameVariable, nil},
			{`([\w.:-]+)`, Text, nil},
		},
		"params": {
			{`;`, Keyword, Pop(1)},
			{`\n`, Text, Pop(1)},
			{`(else|elseif|then)\b`, Keyword, nil},
			Include("basic"),
			Include("data"),
		},
		"params-in-brace": {
			{`\}`, Keyword, Push("#pop", "#pop")},
			Include("params"),
		},
		"params-in-paren": {
			{`\)`, Keyword, Push("#pop", "#pop")},
			Include("params"),
		},
		"params-in-bracket": {
			{`\]`, Keyword, Push("#pop", "#pop")},
			Include("params"),
		},
		"string": {
			{`\[`, LiteralStringDouble, Push("string-square")},
			{`(?s)(\\\\|\\[0-7]+|\\.|[^"\\])`, LiteralStringDouble, nil},
			{`"`, LiteralStringDouble, Pop(1)},
		},
		"string-square": {
			{`\[`, LiteralStringDouble, Push("string-square")},
			{`(?s)(\\\\|\\[0-7]+|\\.|\\\n|[^\]\\])`, LiteralStringDouble, nil},
			{`\]`, LiteralStringDouble, Pop(1)},
		},
		"brace": {
			{`\}`, Keyword, Pop(1)},
			Include("command-in-brace"),
			Include("basic"),
			Include("data"),
		},
		"paren": {
			{`\)`, Keyword, Pop(1)},
			Include("command-in-paren"),
			Include("basic"),
			Include("data"),
		},
		"bracket": {
			{`\]`, Keyword, Pop(1)},
			Include("command-in-bracket"),
			Include("basic"),
			Include("data"),
		},
		"comment": {
			{`.*[^\\]\n`, Comment, Pop(1)},
			{`.*\\\n`, Comment, nil},
		},
	}
}
