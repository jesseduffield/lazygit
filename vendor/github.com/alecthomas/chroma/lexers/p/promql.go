package p

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Promql lexer.
var Promql = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "PromQL",
		Aliases:   []string{"promql"},
		Filenames: []string{"*.promql"},
		MimeTypes: []string{},
	},
	promqlRules,
))

func promqlRules() Rules {
	return Rules{
		"root": {
			{`\n`, TextWhitespace, nil},
			{`\s+`, TextWhitespace, nil},
			{`,`, Punctuation, nil},
			{Words(``, `\b`, `bool`, `by`, `group_left`, `group_right`, `ignoring`, `offset`, `on`, `without`), Keyword, nil},
			{Words(``, `\b`, `sum`, `min`, `max`, `avg`, `group`, `stddev`, `stdvar`, `count`, `count_values`, `bottomk`, `topk`, `quantile`), Keyword, nil},
			{Words(``, `\b`, `abs`, `absent`, `absent_over_time`, `avg_over_time`, `ceil`, `changes`, `clamp_max`, `clamp_min`, `count_over_time`, `day_of_month`, `day_of_week`, `days_in_month`, `delta`, `deriv`, `exp`, `floor`, `histogram_quantile`, `holt_winters`, `hour`, `idelta`, `increase`, `irate`, `label_join`, `label_replace`, `ln`, `log10`, `log2`, `max_over_time`, `min_over_time`, `minute`, `month`, `predict_linear`, `quantile_over_time`, `rate`, `resets`, `round`, `scalar`, `sort`, `sort_desc`, `sqrt`, `stddev_over_time`, `stdvar_over_time`, `sum_over_time`, `time`, `timestamp`, `vector`, `year`), KeywordReserved, nil},
			{`[1-9][0-9]*[smhdwy]`, LiteralString, nil},
			{`-?[0-9]+\.[0-9]+`, LiteralNumberFloat, nil},
			{`-?[0-9]+`, LiteralNumberInteger, nil},
			{`#.*?$`, CommentSingle, nil},
			{`(\+|\-|\*|\/|\%|\^)`, Operator, nil},
			{`==|!=|>=|<=|<|>`, Operator, nil},
			{`and|or|unless`, OperatorWord, nil},
			{`[_a-zA-Z][a-zA-Z0-9_]+`, NameVariable, nil},
			{`(["\'])(.*?)(["\'])`, ByGroups(Punctuation, LiteralString, Punctuation), nil},
			{`\(`, Operator, Push("function")},
			{`\)`, Operator, nil},
			{`\{`, Punctuation, Push("labels")},
			{`\[`, Punctuation, Push("range")},
		},
		"labels": {
			{`\}`, Punctuation, Pop(1)},
			{`\n`, TextWhitespace, nil},
			{`\s+`, TextWhitespace, nil},
			{`,`, Punctuation, nil},
			{`([_a-zA-Z][a-zA-Z0-9_]*?)(\s*?)(=~|!=|=|!~)(\s*?)("|')(.*?)("|')`, ByGroups(NameLabel, TextWhitespace, Operator, TextWhitespace, Punctuation, LiteralString, Punctuation), nil},
		},
		"range": {
			{`\]`, Punctuation, Pop(1)},
			{`[1-9][0-9]*[smhdwy]`, LiteralString, nil},
		},
		"function": {
			{`\)`, Operator, Pop(1)},
			{`\(`, Operator, Push()},
			Default(Pop(1)),
		},
	}
}
