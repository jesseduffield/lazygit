package m

import (
	"regexp"

	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

var (
	mysqlAnalyserNameBetweenBacktickRe = regexp.MustCompile("`[a-zA-Z_]\\w*`")
	mysqlAnalyserNameBetweenBracketRe  = regexp.MustCompile(`\[[a-zA-Z_]\w*\]`)
)

// MySQL lexer.
var MySQL = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "MySQL",
		Aliases:         []string{"mysql"},
		Filenames:       []string{"*.sql"},
		MimeTypes:       []string{"text/x-mysql"},
		NotMultiline:    true,
		CaseInsensitive: true,
	},
	mySQLRules,
).SetAnalyser(func(text string) float32 {
	nameBetweenBacktickCount := len(mysqlAnalyserNameBetweenBacktickRe.FindAllString(text, -1))
	nameBetweenBracketCount := len(mysqlAnalyserNameBetweenBracketRe.FindAllString(text, -1))

	var result float32

	// Same logic as above in the TSQL analysis.
	dialectNameCount := nameBetweenBacktickCount + nameBetweenBracketCount
	if dialectNameCount >= 1 && nameBetweenBacktickCount >= (2*nameBetweenBracketCount) {
		// Found at least twice as many `name` as [name].
		result += 0.5
	} else if nameBetweenBacktickCount > nameBetweenBracketCount {
		result += 0.2
	} else if nameBetweenBacktickCount > 0 {
		result += 0.1
	}

	return result
}))

func mySQLRules() Rules {
	return Rules{
		"root": {
			{`\s+`, TextWhitespace, nil},
			{`(#|--\s+).*\n?`, CommentSingle, nil},
			{`/\*`, CommentMultiline, Push("multiline-comments")},
			{`[0-9]+`, LiteralNumberInteger, nil},
			{`[0-9]*\.[0-9]+(e[+-][0-9]+)`, LiteralNumberFloat, nil},
			{`((?:_[a-z0-9]+)?)(')`, ByGroups(LiteralStringAffix, LiteralStringSingle), Push("string")},
			{`((?:_[a-z0-9]+)?)(")`, ByGroups(LiteralStringAffix, LiteralStringDouble), Push("double-string")},
			{"[+*/<>=~!@#%^&|`?-]", Operator, nil},
			{`\b(tinyint|smallint|mediumint|int|integer|bigint|date|datetime|time|bit|bool|tinytext|mediumtext|longtext|text|tinyblob|mediumblob|longblob|blob|float|double|double\s+precision|real|numeric|dec|decimal|timestamp|year|char|varchar|varbinary|varcharacter|enum|set)(\b\s*)(\()?`, ByGroups(KeywordType, TextWhitespace, Punctuation), nil},
			{`\b(add|all|alter|analyze|and|as|asc|asensitive|before|between|bigint|binary|blob|both|by|call|cascade|case|change|char|character|check|collate|column|condition|constraint|continue|convert|create|cross|current_date|current_time|current_timestamp|current_user|cursor|database|databases|day_hour|day_microsecond|day_minute|day_second|dec|decimal|declare|default|delayed|delete|desc|describe|deterministic|distinct|distinctrow|div|double|drop|dual|each|else|elseif|enclosed|escaped|exists|exit|explain|fetch|flush|float|float4|float8|for|force|foreign|from|fulltext|grant|group|having|high_priority|hour_microsecond|hour_minute|hour_second|identified|if|ignore|in|index|infile|inner|inout|insensitive|insert|int|int1|int2|int3|int4|int8|integer|interval|into|is|iterate|join|key|keys|kill|leading|leave|left|like|limit|lines|load|localtime|localtimestamp|lock|long|loop|low_priority|match|minute_microsecond|minute_second|mod|modifies|natural|no_write_to_binlog|not|numeric|on|optimize|option|optionally|or|order|out|outer|outfile|precision|primary|privileges|procedure|purge|raid0|read|reads|real|references|regexp|release|rename|repeat|replace|require|restrict|return|revoke|right|rlike|schema|schemas|second_microsecond|select|sensitive|separator|set|show|smallint|soname|spatial|specific|sql|sql_big_result|sql_calc_found_rows|sql_small_result|sqlexception|sqlstate|sqlwarning|ssl|starting|straight_join|table|terminated|then|to|trailing|trigger|undo|union|unique|unlock|unsigned|update|usage|use|user|using|utc_date|utc_time|utc_timestamp|values|varying|when|where|while|with|write|x509|xor|year_month|zerofill)\b`, Keyword, nil},
			{`\b(auto_increment|engine|charset|tables)\b`, KeywordPseudo, nil},
			{`(true|false|null)`, NameConstant, nil},
			{`([a-z_]\w*)(\s*)(\()`, ByGroups(NameFunction, TextWhitespace, Punctuation), nil},
			{`[a-z_]\w*`, Name, nil},
			{`@[a-z0-9]*[._]*[a-z0-9]*`, NameVariable, nil},
			{`[;:()\[\],.]`, Punctuation, nil},
		},
		"multiline-comments": {
			{`/\*`, CommentMultiline, Push("multiline-comments")},
			{`\*/`, CommentMultiline, Pop(1)},
			{`[^/*]+`, CommentMultiline, nil},
			{`[/*]`, CommentMultiline, nil},
		},
		"string": {
			{`[^']+`, LiteralStringSingle, nil},
			{`''`, LiteralStringSingle, nil},
			{`'`, LiteralStringSingle, Pop(1)},
		},
		"double-string": {
			{`[^"]+`, LiteralStringDouble, nil},
			{`""`, LiteralStringDouble, nil},
			{`"`, LiteralStringDouble, Pop(1)},
		},
	}
}
