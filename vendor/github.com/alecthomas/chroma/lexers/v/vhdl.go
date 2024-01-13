package v

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// VHDL lexer.
var VHDL = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "VHDL",
		Aliases:         []string{"vhdl"},
		Filenames:       []string{"*.vhdl", "*.vhd"},
		MimeTypes:       []string{"text/x-vhdl"},
		CaseInsensitive: true,
	},
	vhdlRules,
))

func vhdlRules() Rules {
	return Rules{
		"root": {
			{`\n`, Text, nil},
			{`\s+`, Text, nil},
			{`\\\n`, Text, nil},
			{`--.*?$`, CommentSingle, nil},
			{`'(U|X|0|1|Z|W|L|H|-)'`, LiteralStringChar, nil},
			{`[~!%^&*+=|?:<>/-]`, Operator, nil},
			{`'[a-z_]\w*`, NameAttribute, nil},
			{`[()\[\],.;\']`, Punctuation, nil},
			{`"[^\n\\"]*"`, LiteralString, nil},
			{`(library)(\s+)([a-z_]\w*)`, ByGroups(Keyword, Text, NameNamespace), nil},
			{`(use)(\s+)(entity)`, ByGroups(Keyword, Text, Keyword), nil},
			{`(use)(\s+)([a-z_][\w.]*\.)(all)`, ByGroups(Keyword, Text, NameNamespace, Keyword), nil},
			{`(use)(\s+)([a-z_][\w.]*)`, ByGroups(Keyword, Text, NameNamespace), nil},
			{`(std|ieee)(\.[a-z_]\w*)`, ByGroups(NameNamespace, NameNamespace), nil},
			{Words(``, `\b`, `std`, `ieee`, `work`), NameNamespace, nil},
			{`(entity|component)(\s+)([a-z_]\w*)`, ByGroups(Keyword, Text, NameClass), nil},
			{`(architecture|configuration)(\s+)([a-z_]\w*)(\s+)(of)(\s+)([a-z_]\w*)(\s+)(is)`, ByGroups(Keyword, Text, NameClass, Text, Keyword, Text, NameClass, Text, Keyword), nil},
			{`([a-z_]\w*)(:)(\s+)(process|for)`, ByGroups(NameClass, Operator, Text, Keyword), nil},
			// This seems to cause a recursive loop.
			// {`(end)(\s+)`, ByGroups(UsingSelf("root"), Text), Push("endblock")},
			{`(end)(\s+)`, ByGroups(Keyword, Text), Push("endblock")},
			Include("types"),
			Include("keywords"),
			Include("numbers"),
			{`[a-z_]\w*`, Name, nil},
		},
		"endblock": {
			Include("keywords"),
			{`[a-z_]\w*`, NameClass, nil},
			{`(\s+)`, Text, nil},
			{`;`, Punctuation, Pop(1)},
		},
		"types": {
			{Words(``, `\b`, `boolean`, `bit`, `character`, `severity_level`, `integer`, `time`, `delay_length`, `natural`, `positive`, `string`, `bit_vector`, `file_open_kind`, `file_open_status`, `std_ulogic`, `std_ulogic_vector`, `std_logic`, `std_logic_vector`, `signed`, `unsigned`), KeywordType, nil},
		},
		"keywords": {
			{Words(``, `\b`, `abs`, `access`, `after`, `alias`, `all`, `and`, `architecture`, `array`, `assert`, `attribute`, `begin`, `block`, `body`, `buffer`, `bus`, `case`, `component`, `configuration`, `constant`, `disconnect`, `downto`, `else`, `elsif`, `end`, `entity`, `exit`, `file`, `for`, `function`, `generate`, `generic`, `group`, `guarded`, `if`, `impure`, `in`, `inertial`, `inout`, `is`, `label`, `library`, `linkage`, `literal`, `loop`, `map`, `mod`, `nand`, `new`, `next`, `nor`, `not`, `null`, `of`, `on`, `open`, `or`, `others`, `out`, `package`, `port`, `postponed`, `procedure`, `process`, `pure`, `range`, `record`, `register`, `reject`, `rem`, `return`, `rol`, `ror`, `select`, `severity`, `signal`, `shared`, `sla`, `sll`, `sra`, `srl`, `subtype`, `then`, `to`, `transport`, `type`, `units`, `until`, `use`, `variable`, `wait`, `when`, `while`, `with`, `xnor`, `xor`), Keyword, nil},
		},
		"numbers": {
			{`\d{1,2}#[0-9a-f_]+#?`, LiteralNumberInteger, nil},
			{`\d+`, LiteralNumberInteger, nil},
			{`(\d+\.\d*|\.\d+|\d+)E[+-]?\d+`, LiteralNumberFloat, nil},
			{`X"[0-9a-f_]+"`, LiteralNumberHex, nil},
			{`O"[0-7_]+"`, LiteralNumberOct, nil},
			{`B"[01_]+"`, LiteralNumberBin, nil},
		},
	}
}
