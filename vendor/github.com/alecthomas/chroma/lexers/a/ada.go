package a

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Ada lexer.
var Ada = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "Ada",
		Aliases:         []string{"ada", "ada95", "ada2005"},
		Filenames:       []string{"*.adb", "*.ads", "*.ada"},
		MimeTypes:       []string{"text/x-ada"},
		CaseInsensitive: true,
	},
	adaRules,
))

func adaRules() Rules {
	return Rules{
		"root": {
			{`[^\S\n]+`, Text, nil},
			{`--.*?\n`, CommentSingle, nil},
			{`[^\S\n]+`, Text, nil},
			{`function|procedure|entry`, KeywordDeclaration, Push("subprogram")},
			{`(subtype|type)(\s+)(\w+)`, ByGroups(KeywordDeclaration, Text, KeywordType), Push("type_def")},
			{`task|protected`, KeywordDeclaration, nil},
			{`(subtype)(\s+)`, ByGroups(KeywordDeclaration, Text), nil},
			{`(end)(\s+)`, ByGroups(KeywordReserved, Text), Push("end")},
			{`(pragma)(\s+)(\w+)`, ByGroups(KeywordReserved, Text, CommentPreproc), nil},
			{`(true|false|null)\b`, KeywordConstant, nil},
			{Words(``, `\b`, `Address`, `Byte`, `Boolean`, `Character`, `Controlled`, `Count`, `Cursor`, `Duration`, `File_Mode`, `File_Type`, `Float`, `Generator`, `Integer`, `Long_Float`, `Long_Integer`, `Long_Long_Float`, `Long_Long_Integer`, `Natural`, `Positive`, `Reference_Type`, `Short_Float`, `Short_Integer`, `Short_Short_Float`, `Short_Short_Integer`, `String`, `Wide_Character`, `Wide_String`), KeywordType, nil},
			{`(and(\s+then)?|in|mod|not|or(\s+else)|rem)\b`, OperatorWord, nil},
			{`generic|private`, KeywordDeclaration, nil},
			{`package`, KeywordDeclaration, Push("package")},
			{`array\b`, KeywordReserved, Push("array_def")},
			{`(with|use)(\s+)`, ByGroups(KeywordNamespace, Text), Push("import")},
			{`(\w+)(\s*)(:)(\s*)(constant)`, ByGroups(NameConstant, Text, Punctuation, Text, KeywordReserved), nil},
			{`<<\w+>>`, NameLabel, nil},
			{`(\w+)(\s*)(:)(\s*)(declare|begin|loop|for|while)`, ByGroups(NameLabel, Text, Punctuation, Text, KeywordReserved), nil},
			{Words(`\b`, `\b`, `abort`, `abs`, `abstract`, `accept`, `access`, `aliased`, `all`, `array`, `at`, `begin`, `body`, `case`, `constant`, `declare`, `delay`, `delta`, `digits`, `do`, `else`, `elsif`, `end`, `entry`, `exception`, `exit`, `interface`, `for`, `goto`, `if`, `is`, `limited`, `loop`, `new`, `null`, `of`, `or`, `others`, `out`, `overriding`, `pragma`, `protected`, `raise`, `range`, `record`, `renames`, `requeue`, `return`, `reverse`, `select`, `separate`, `subtype`, `synchronized`, `task`, `tagged`, `terminate`, `then`, `type`, `until`, `when`, `while`, `xor`), KeywordReserved, nil},
			{`"[^"]*"`, LiteralString, nil},
			Include("attribute"),
			Include("numbers"),
			{`'[^']'`, LiteralStringChar, nil},
			{`(\w+)(\s*|[(,])`, ByGroups(Name, UsingSelf("root")), nil},
			{`(<>|=>|:=|[()|:;,.'])`, Punctuation, nil},
			{`[*<>+=/&-]`, Operator, nil},
			{`\n+`, Text, nil},
		},
		"numbers": {
			{`[0-9_]+#[0-9a-f]+#`, LiteralNumberHex, nil},
			{`[0-9_]+\.[0-9_]*`, LiteralNumberFloat, nil},
			{`[0-9_]+`, LiteralNumberInteger, nil},
		},
		"attribute": {
			{`(')(\w+)`, ByGroups(Punctuation, NameAttribute), nil},
		},
		"subprogram": {
			{`\(`, Punctuation, Push("#pop", "formal_part")},
			{`;`, Punctuation, Pop(1)},
			{`is\b`, KeywordReserved, Pop(1)},
			{`"[^"]+"|\w+`, NameFunction, nil},
			Include("root"),
		},
		"end": {
			{`(if|case|record|loop|select)`, KeywordReserved, nil},
			{`"[^"]+"|[\w.]+`, NameFunction, nil},
			{`\s+`, Text, nil},
			{`;`, Punctuation, Pop(1)},
		},
		"type_def": {
			{`;`, Punctuation, Pop(1)},
			{`\(`, Punctuation, Push("formal_part")},
			{`with|and|use`, KeywordReserved, nil},
			{`array\b`, KeywordReserved, Push("#pop", "array_def")},
			{`record\b`, KeywordReserved, Push("record_def")},
			{`(null record)(;)`, ByGroups(KeywordReserved, Punctuation), Pop(1)},
			Include("root"),
		},
		"array_def": {
			{`;`, Punctuation, Pop(1)},
			{`(\w+)(\s+)(range)`, ByGroups(KeywordType, Text, KeywordReserved), nil},
			Include("root"),
		},
		"record_def": {
			{`end record`, KeywordReserved, Pop(1)},
			Include("root"),
		},
		"import": {
			{`[\w.]+`, NameNamespace, Pop(1)},
			Default(Pop(1)),
		},
		"formal_part": {
			{`\)`, Punctuation, Pop(1)},
			{`\w+`, NameVariable, nil},
			{`,|:[^=]`, Punctuation, nil},
			{`(in|not|null|out|access)\b`, KeywordReserved, nil},
			Include("root"),
		},
		"package": {
			{`body`, KeywordDeclaration, nil},
			{`is\s+new|renames`, KeywordReserved, nil},
			{`is`, KeywordReserved, Pop(1)},
			{`;`, Punctuation, Pop(1)},
			{`\(`, Punctuation, Push("package_instantiation")},
			{`([\w.]+)`, NameClass, nil},
			Include("root"),
		},
		"package_instantiation": {
			{`("[^"]+"|\w+)(\s+)(=>)`, ByGroups(NameVariable, Text, Punctuation), nil},
			{`[\w.\'"]`, Text, nil},
			{`\)`, Punctuation, Pop(1)},
			Include("root"),
		},
	}
}
