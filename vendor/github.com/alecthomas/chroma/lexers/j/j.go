package j

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// J lexer.
var J = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "J",
		Aliases:   []string{"j"},
		Filenames: []string{"*.ijs"},
		MimeTypes: []string{"text/x-j"},
	},
	jRules,
))

func jRules() Rules {
	return Rules{
		"root": {
			{`#!.*$`, CommentPreproc, nil},
			{`NB\..*`, CommentSingle, nil},
			{`\n+\s*Note`, CommentMultiline, Push("comment")},
			{`\s*Note.*`, CommentSingle, nil},
			{`\s+`, Text, nil},
			{`'`, LiteralString, Push("singlequote")},
			{`0\s+:\s*0|noun\s+define\s*$`, NameEntity, Push("nounDefinition")},
			{`(([1-4]|13)\s+:\s*0|(adverb|conjunction|dyad|monad|verb)\s+define)\b`, NameFunction, Push("explicitDefinition")},
			{Words(``, `\b[a-zA-Z]\w*\.`, `for_`, `goto_`, `label_`), NameLabel, nil},
			{Words(``, `\.`, `assert`, `break`, `case`, `catch`, `catchd`, `catcht`, `continue`, `do`, `else`, `elseif`, `end`, `fcase`, `for`, `if`, `return`, `select`, `throw`, `try`, `while`, `whilst`), NameLabel, nil},
			{`\b[a-zA-Z]\w*`, NameVariable, nil},
			{Words(``, ``, `ARGV`, `CR`, `CRLF`, `DEL`, `Debug`, `EAV`, `EMPTY`, `FF`, `JVERSION`, `LF`, `LF2`, `Note`, `TAB`, `alpha17`, `alpha27`, `apply`, `bind`, `boxopen`, `boxxopen`, `bx`, `clear`, `cutLF`, `cutopen`, `datatype`, `def`, `dfh`, `drop`, `each`, `echo`, `empty`, `erase`, `every`, `evtloop`, `exit`, `expand`, `fetch`, `file2url`, `fixdotdot`, `fliprgb`, `getargs`, `getenv`, `hfd`, `inv`, `inverse`, `iospath`, `isatty`, `isutf8`, `items`, `leaf`, `list`, `nameclass`, `namelist`, `names`, `nc`, `nl`, `on`, `pick`, `rows`, `script`, `scriptd`, `sign`, `sminfo`, `smoutput`, `sort`, `split`, `stderr`, `stdin`, `stdout`, `table`, `take`, `timespacex`, `timex`, `tmoutput`, `toCRLF`, `toHOST`, `toJ`, `tolower`, `toupper`, `type`, `ucp`, `ucpcount`, `usleep`, `utf8`, `uucp`), NameFunction, nil},
			{`=[.:]`, Operator, nil},
			{"[-=+*#$%@!~`^&\";:.,<>{}\\[\\]\\\\|/]", Operator, nil},
			{`[abCdDeEfHiIjLMoprtT]\.`, KeywordReserved, nil},
			{`[aDiLpqsStux]\:`, KeywordReserved, nil},
			{`(_[0-9])\:`, KeywordConstant, nil},
			{`\(`, Punctuation, Push("parentheses")},
			Include("numbers"),
		},
		"comment": {
			{`[^)]`, CommentMultiline, nil},
			{`^\)`, CommentMultiline, Pop(1)},
			{`[)]`, CommentMultiline, nil},
		},
		"explicitDefinition": {
			{`\b[nmuvxy]\b`, NameDecorator, nil},
			Include("root"),
			{`[^)]`, Name, nil},
			{`^\)`, NameLabel, Pop(1)},
			{`[)]`, Name, nil},
		},
		"numbers": {
			{`\b_{1,2}\b`, LiteralNumber, nil},
			{`_?\d+(\.\d+)?(\s*[ejr]\s*)_?\d+(\.?=\d+)?`, LiteralNumber, nil},
			{`_?\d+\.(?=\d+)`, LiteralNumberFloat, nil},
			{`_?\d+x`, LiteralNumberIntegerLong, nil},
			{`_?\d+`, LiteralNumberInteger, nil},
		},
		"nounDefinition": {
			{`[^)]`, LiteralString, nil},
			{`^\)`, NameLabel, Pop(1)},
			{`[)]`, LiteralString, nil},
		},
		"parentheses": {
			{`\)`, Punctuation, Pop(1)},
			Include("explicitDefinition"),
			Include("root"),
		},
		"singlequote": {
			{`[^']`, LiteralString, nil},
			{`''`, LiteralString, nil},
			{`'`, LiteralString, Pop(1)},
		},
	}
}
