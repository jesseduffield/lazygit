package t

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Terraform lexer.
var Terraform = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Terraform",
		Aliases:   []string{"terraform", "tf"},
		Filenames: []string{"*.tf"},
		MimeTypes: []string{"application/x-tf", "application/x-terraform"},
	},
	terraformRules,
))

func terraformRules() Rules {
	return Rules{
		"root": {
			{`[\[\](),.{}]`, Punctuation, nil},
			{`-?[0-9]+`, LiteralNumber, nil},
			{`=>`, Punctuation, nil},
			{Words(``, `\b`, `true`, `false`), KeywordConstant, nil},
			{`/(?s)\*(((?!\*/).)*)\*/`, CommentMultiline, nil},
			{`\s*(#|//).*\n`, CommentSingle, nil},
			{`([a-zA-Z]\w*)(\s*)(=(?!>))`, ByGroups(NameAttribute, Text, Text), nil},
			{Words(`^\s*`, `\b`, `variable`, `data`, `resource`, `provider`, `provisioner`, `module`, `output`), KeywordReserved, nil},
			{Words(``, `\b`, `for`, `in`), Keyword, nil},
			{Words(``, ``, `count`, `data`, `var`, `module`, `each`), NameBuiltin, nil},
			{Words(``, `\b`, `abs`, `ceil`, `floor`, `log`, `max`, `min`, `parseint`, `pow`, `signum`), NameBuiltin, nil},
			{Words(``, `\b`, `chomp`, `format`, `formatlist`, `indent`, `join`, `lower`, `regex`, `regexall`, `replace`, `split`, `strrev`, `substr`, `title`, `trim`, `trimprefix`, `trimsuffix`, `trimspace`, `upper`), NameBuiltin, nil},
			{Words(`[^.]`, `\b`, `chunklist`, `coalesce`, `coalescelist`, `compact`, `concat`, `contains`, `distinct`, `element`, `flatten`, `index`, `keys`, `length`, `list`, `lookup`, `map`, `matchkeys`, `merge`, `range`, `reverse`, `setintersection`, `setproduct`, `setsubtract`, `setunion`, `slice`, `sort`, `transpose`, `values`, `zipmap`), NameBuiltin, nil},
			{Words(`[^.]`, `\b`, `base64decode`, `base64encode`, `base64gzip`, `csvdecode`, `jsondecode`, `jsonencode`, `urlencode`, `yamldecode`, `yamlencode`), NameBuiltin, nil},
			{Words(``, `\b`, `abspath`, `dirname`, `pathexpand`, `basename`, `file`, `fileexists`, `fileset`, `filebase64`, `templatefile`), NameBuiltin, nil},
			{Words(``, `\b`, `formatdate`, `timeadd`, `timestamp`), NameBuiltin, nil},
			{Words(``, `\b`, `base64sha256`, `base64sha512`, `bcrypt`, `filebase64sha256`, `filebase64sha512`, `filemd5`, `filesha1`, `filesha256`, `filesha512`, `md5`, `rsadecrypt`, `sha1`, `sha256`, `sha512`, `uuid`, `uuidv5`), NameBuiltin, nil},
			{Words(``, `\b`, `cidrhost`, `cidrnetmask`, `cidrsubnet`), NameBuiltin, nil},
			{Words(``, `\b`, `can`, `tobool`, `tolist`, `tomap`, `tonumber`, `toset`, `tostring`, `try`), NameBuiltin, nil},
			{`=(?!>)|\+|-|\*|\/|:|!|%|>|<(?!<)|>=|<=|==|!=|&&|\||\?`, Operator, nil},
			{`\n|\s+|\\\n`, Text, nil},
			{`[a-zA-Z]\w*`, NameOther, nil},
			{`"`, LiteralStringDouble, Push("string")},
			{`(?s)(<<-?)(\w+)(\n\s*(?:(?!\2).)*\s*\n\s*)(\2)`, ByGroups(Operator, Operator, String, Operator), nil},
		},
		"declaration": {
			{`(\s*)("(?:\\\\|\\"|[^"])*")(\s*)`, ByGroups(Text, NameVariable, Text), nil},
			{`\{`, Punctuation, Pop(1)},
		},
		"string": {
			{`"`, LiteralStringDouble, Pop(1)},
			{`\\\\`, LiteralStringDouble, nil},
			{`\\\\"`, LiteralStringDouble, nil},
			{`\$\{`, LiteralStringInterpol, Push("interp-inside")},
			{`\$`, LiteralStringDouble, nil},
			{`[^"\\\\$]+`, LiteralStringDouble, nil},
		},
		"interp-inside": {
			{`\}`, LiteralStringInterpol, Pop(1)},
			Include("root"),
		},
	}
}
