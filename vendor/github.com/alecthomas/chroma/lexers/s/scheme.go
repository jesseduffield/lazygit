package s

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// nolint

// Scheme lexer.
var SchemeLang = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Scheme",
		Aliases:   []string{"scheme", "scm"},
		Filenames: []string{"*.scm", "*.ss"},
		MimeTypes: []string{"text/x-scheme", "application/x-scheme"},
	},
	schemeLangRules,
))

func schemeLangRules() Rules {
	return Rules{
		"root": {
			{`;.*$`, CommentSingle, nil},
			{`#\|`, CommentMultiline, Push("multiline-comment")},
			{`#;\s*\(`, Comment, Push("commented-form")},
			{`#!r6rs`, Comment, nil},
			{`\s+`, Text, nil},
			{`-?\d+\.\d+`, LiteralNumberFloat, nil},
			{`-?\d+`, LiteralNumberInteger, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralString, nil},
			{`'[\w!$%&*+,/:<=>?@^~|-]+`, LiteralStringSymbol, nil},
			{`#\\(alarm|backspace|delete|esc|linefeed|newline|page|return|space|tab|vtab|x[0-9a-zA-Z]{1,5}|.)`, LiteralStringChar, nil},
			{`(#t|#f)`, NameConstant, nil},
			{"('|#|`|,@|,|\\.)", Operator, nil},
			{`(lambda |define |if |else |cond |and |or |case |let |let\* |letrec |begin |do |delay |set\! |\=\> |quote |quasiquote |unquote |unquote\-splicing |define\-syntax |let\-syntax |letrec\-syntax |syntax\-rules )`, Keyword, nil},
			{`(?<='\()[\w!$%&*+,/:<=>?@^~|-]+`, NameVariable, nil},
			{`(?<=#\()[\w!$%&*+,/:<=>?@^~|-]+`, NameVariable, nil},
			{`(?<=\()(\* |\+ |\- |\/ |\< |\<\= |\= |\> |\>\= |abs |acos |angle |append |apply |asin |assoc |assq |assv |atan |boolean\? |caaaar |caaadr |caaar |caadar |caaddr |caadr |caar |cadaar |cadadr |cadar |caddar |cadddr |caddr |cadr |call\-with\-current\-continuation |call\-with\-input\-file |call\-with\-output\-file |call\-with\-values |call\/cc |car |cdaaar |cdaadr |cdaar |cdadar |cdaddr |cdadr |cdar |cddaar |cddadr |cddar |cdddar |cddddr |cdddr |cddr |cdr |ceiling |char\-\>integer |char\-alphabetic\? |char\-ci\<\=\? |char\-ci\<\? |char\-ci\=\? |char\-ci\>\=\? |char\-ci\>\? |char\-downcase |char\-lower\-case\? |char\-numeric\? |char\-ready\? |char\-upcase |char\-upper\-case\? |char\-whitespace\? |char\<\=\? |char\<\? |char\=\? |char\>\=\? |char\>\? |char\? |close\-input\-port |close\-output\-port |complex\? |cons |cos |current\-input\-port |current\-output\-port |denominator |display |dynamic\-wind |eof\-object\? |eq\? |equal\? |eqv\? |eval |even\? |exact\-\>inexact |exact\? |exp |expt |floor |for\-each |force |gcd |imag\-part |inexact\-\>exact |inexact\? |input\-port\? |integer\-\>char |integer\? |interaction\-environment |lcm |length |list |list\-\>string |list\-\>vector |list\-ref |list\-tail |list\? |load |log |magnitude |make\-polar |make\-rectangular |make\-string |make\-vector |map |max |member |memq |memv |min |modulo |negative\? |newline |not |null\-environment |null\? |number\-\>string |number\? |numerator |odd\? |open\-input\-file |open\-output\-file |output\-port\? |pair\? |peek\-char |port\? |positive\? |procedure\? |quotient |rational\? |rationalize |read |read\-char |real\-part |real\? |remainder |reverse |round |scheme\-report\-environment |set\-car\! |set\-cdr\! |sin |sqrt |string |string\-\>list |string\-\>number |string\-\>symbol |string\-append |string\-ci\<\=\? |string\-ci\<\? |string\-ci\=\? |string\-ci\>\=\? |string\-ci\>\? |string\-copy |string\-fill\! |string\-length |string\-ref |string\-set\! |string\<\=\? |string\<\? |string\=\? |string\>\=\? |string\>\? |string\? |substring |symbol\-\>string |symbol\? |tan |transcript\-off |transcript\-on |truncate |values |vector |vector\-\>list |vector\-fill\! |vector\-length |vector\-ref |vector\-set\! |vector\? |with\-input\-from\-file |with\-output\-to\-file |write |write\-char |zero\? )`, NameBuiltin, nil},
			{`(?<=\()[\w!$%&*+,/:<=>?@^~|-]+`, NameFunction, nil},
			{`[\w!$%&*+,/:<=>?@^~|-]+`, NameVariable, nil},
			{`(\(|\))`, Punctuation, nil},
			{`(\[|\])`, Punctuation, nil},
		},
		"multiline-comment": {
			{`#\|`, CommentMultiline, Push()},
			{`\|#`, CommentMultiline, Pop(1)},
			{`[^|#]+`, CommentMultiline, nil},
			{`[|#]`, CommentMultiline, nil},
		},
		"commented-form": {
			{`\(`, Comment, Push()},
			{`\)`, Comment, Pop(1)},
			{`[^()]+`, Comment, nil},
		},
	}
}
