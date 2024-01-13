package p

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Pig lexer.
var Pig = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "Pig",
		Aliases:         []string{"pig"},
		Filenames:       []string{"*.pig"},
		MimeTypes:       []string{"text/x-pig"},
		CaseInsensitive: true,
	},
	pigRules,
))

func pigRules() Rules {
	return Rules{
		"root": {
			{`\s+`, Text, nil},
			{`--.*`, Comment, nil},
			{`/\*[\w\W]*?\*/`, CommentMultiline, nil},
			{`\\\n`, Text, nil},
			{`\\`, Text, nil},
			{`\'(?:\\[ntbrf\\\']|\\u[0-9a-f]{4}|[^\'\\\n\r])*\'`, LiteralString, nil},
			Include("keywords"),
			Include("types"),
			Include("builtins"),
			Include("punct"),
			Include("operators"),
			{`[0-9]*\.[0-9]+(e[0-9]+)?[fd]?`, LiteralNumberFloat, nil},
			{`0x[0-9a-f]+`, LiteralNumberHex, nil},
			{`[0-9]+L?`, LiteralNumberInteger, nil},
			{`\n`, Text, nil},
			{`([a-z_]\w*)(\s*)(\()`, ByGroups(NameFunction, Text, Punctuation), nil},
			{`[()#:]`, Text, nil},
			{`[^(:#\'")\s]+`, Text, nil},
			{`\S+\s+`, Text, nil},
		},
		"keywords": {
			{`(assert|and|any|all|arrange|as|asc|bag|by|cache|CASE|cat|cd|cp|%declare|%default|define|dense|desc|describe|distinct|du|dump|eval|exex|explain|filter|flatten|foreach|full|generate|group|help|if|illustrate|import|inner|input|into|is|join|kill|left|limit|load|ls|map|matches|mkdir|mv|not|null|onschema|or|order|outer|output|parallel|pig|pwd|quit|register|returns|right|rm|rmf|rollup|run|sample|set|ship|split|stderr|stdin|stdout|store|stream|through|union|using|void)\b`, Keyword, nil},
		},
		"builtins": {
			{`(AVG|BinStorage|cogroup|CONCAT|copyFromLocal|copyToLocal|COUNT|cross|DIFF|MAX|MIN|PigDump|PigStorage|SIZE|SUM|TextLoader|TOKENIZE)\b`, NameBuiltin, nil},
		},
		"types": {
			{`(bytearray|BIGINTEGER|BIGDECIMAL|chararray|datetime|double|float|int|long|tuple)\b`, KeywordType, nil},
		},
		"punct": {
			{`[;(){}\[\]]`, Punctuation, nil},
		},
		"operators": {
			{`[#=,./%+\-?]`, Operator, nil},
			{`(eq|gt|lt|gte|lte|neq|matches)\b`, Operator, nil},
			{`(==|<=|<|>=|>|!=)`, Operator, nil},
		},
	}
}
