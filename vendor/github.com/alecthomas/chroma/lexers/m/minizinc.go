package m

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// MiniZinc lexer.
var MZN = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "MiniZinc",
		Aliases:   []string{"minizinc", "MZN", "mzn"},
		Filenames: []string{"*.mzn", "*.dzn", "*.fzn"},
		MimeTypes: []string{"text/minizinc"},
	},
	mznRules,
))

func mznRules() Rules {
	return Rules{
		"root": {
			{`\n`, Text, nil},
			{`\s+`, Text, nil},
			{`\\\n`, Text, nil},
			{`\%(.*?)\n`, CommentSingle, nil},
			{`/(\\\n)?[*](.|\n)*?[*](\\\n)?/`, CommentMultiline, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralString, nil},
			{Words(`\b`, `\b`, `ann`, `annotation`, `any`, `constraint`, `function`, `include`, `list`, `of`, `op`, `output`, `minimize`, `maximize`, `par`, `predicate`, `record`, `satisfy`, `solve`, `test`, `type`, `var`), Keyword, nil},
			{Words(`\b`, `\b`, `array`, `set`, `bool`, `enum`, `float`, `int`, `string`, `tuple`), KeywordType, nil},
			{Words(`\b`, `\b`, `for`, `forall`, `if`, `then`, `else`, `endif`, `where`), Keyword, nil},
			{Words(`\b`, `\b`, `abort`, `abs`, `acosh`, `array_intersect`, `array_union`, `array1d`, `array2d`, `array3d`, `array4d`, `array5d`, `array6d`, `asin`, `assert`, `atan`, `bool2int`, `card`, `ceil`, `concat`, `cos`, `cosh`, `dom`, `dom_array`, `dom_size`, `fix`, `exp`, `floor`, `index_set`, `index_set_1of2`, `index_set_2of2`, `index_set_1of3`, `index_set_2of3`, `index_set_3of3`, `int2float`, `is_fixed`, `join`, `lb`, `lb_array`, `length`, `ln`, `log`, `log2`, `log10`, `min`, `max`, `pow`, `product`, `round`, `set2array`, `show`, `show_int`, `show_float`, `sin`, `sinh`, `sqrt`, `sum`, `tan`, `tanh`, `trace`, `ub`, `ub_array`), NameBuiltin, nil},
			{`(not|<->|->|<-|\\/|xor|/\\)`, Operator, nil},
			{`(<|>|<=|>=|==|=|!=)`, Operator, nil},
			{`(\+|-|\*|/|div|mod)`, Operator, nil},
			{Words(`\b`, `\b`, `in`, `subset`, `superset`, `union`, `diff`, `symdiff`, `intersect`), Operator, nil},
			{`(\\|\.\.|\+\+)`, Operator, nil},
			{`[|()\[\]{},:;]`, Punctuation, nil},
			{`(true|false)\b`, KeywordConstant, nil},
			{`([+-]?)\d+(\.(?!\.)\d*)?([eE][-+]?\d+)?`, LiteralNumber, nil},
			{`::\s*([^\W\d]\w*)(\s*\([^\)]*\))?`, NameDecorator, nil},
			{`\b([^\W\d]\w*)\b(\()`, ByGroups(NameFunction, Punctuation), nil},
			{`[^\W\d]\w*`, NameOther, nil},
		},
	}
}
