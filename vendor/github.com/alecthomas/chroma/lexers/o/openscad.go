package o

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

var OpenSCAD = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "OpenSCAD",
		Aliases:   []string{"openscad"},
		Filenames: []string{"*.scad"},
		MimeTypes: []string{"text/x-scad"},
	},
	openSCADRules,
))

func openSCADRules() Rules {
	return Rules{
		"root": {
			{`[^\S\n]+`, Text, nil},
			{`\n`, Text, nil},
			{`//(\n|[\w\W]*?[^\\]\n)`, CommentSingle, nil},
			{`/(\\\n)?[*][\w\W]*?[*](\\\n)?/`, CommentMultiline, nil},
			{`/(\\\n)?[*][\w\W]*`, CommentMultiline, nil},
			{`[{}\[\]\(\),;:]`, Punctuation, nil},
			{`[*!#%\-+=?/]`, Operator, nil},
			{`<|<=|==|!=|>=|>|&&|\|\|`, Operator, nil},
			{`\$(f[asn]|t|vp[rtd]|children)`, NameVariableMagic, nil},
			{Words(``, `\b`, `PI`, `undef`), KeywordConstant, nil},
			{`(use|include)((?:\s|\\\\s)+)`, ByGroups(KeywordNamespace, Text), Push("includes")},
			{`(module)(\s*)([^\s\(]+)`, ByGroups(KeywordNamespace, Text, NameNamespace), nil},
			{`(function)(\s*)([^\s\(]+)`, ByGroups(KeywordDeclaration, Text, NameFunction), nil},
			{`\b(true|false)\b`, Literal, nil},
			{`\b(function|module|include|use|for|intersection_for|if|else|return)\b`, Keyword, nil},
			{`\b(circle|square|polygon|text|sphere|cube|cylinder|polyhedron|translate|rotate|scale|resize|mirror|multmatrix|color|offset|hull|minkowski|union|difference|intersection|abs|sign|sin|cos|tan|acos|asin|atan|atan2|floor|round|ceil|ln|log|pow|sqrt|exp|rands|min|max|concat|lookup|str|chr|search|version|version_num|norm|cross|parent_module|echo|import|import_dxf|dxf_linear_extrude|linear_extrude|rotate_extrude|surface|projection|render|dxf_cross|dxf_dim|let|assign|len)\b`, NameBuiltin, nil},
			{`\bchildren\b`, NameBuiltinPseudo, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralStringDouble, nil},
			{`-?\d+(\.\d+)?(e[+-]?\d+)?`, Number, nil},
			{`[a-zA-Z_]\w*`, Name, nil},
		},
		"includes": {
			{"(<)([^>]*)(>)", ByGroups(Punctuation, CommentPreprocFile, Punctuation), nil},
			Default(Pop(1)),
		},
	}
}
