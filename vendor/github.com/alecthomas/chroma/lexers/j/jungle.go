package j

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

var Jungle = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Jungle",
		Aliases:   []string{"jungle"},
		Filenames: []string{"*.jungle"},
		MimeTypes: []string{"text/x-jungle"},
	},
	jungleRules,
))

func jungleRules() Rules {
	return Rules{
		"root": {
			{`[^\S\n]+`, Text, nil},
			{`\n`, Text, nil},
			{`#(\n|[\w\W]*?[^#]\n)`, CommentSingle, nil},
			{`^(?=\S)`, None, Push("instruction")},
			{`[\.;\[\]\(\)\$]`, Punctuation, nil},
			{`[a-zA-Z_]\w*`, Name, nil},
		},
		"instruction": {
			{`[^\S\n]+`, Text, nil},
			{`=`, Operator, Push("value")},
			{`(?=\S)`, None, Push("var")},
			Default(Pop(1)),
		},
		"value": {
			{`[^\S\n]+`, Text, nil},
			{`\$\(`, Punctuation, Push("var")},
			{`[;\[\]\(\)\$]`, Punctuation, nil},
			{`#(\n|[\w\W]*?[^#]\n)`, CommentSingle, nil},
			{`[\w_\-\.\/\\]+`, Text, nil},
			Default(Pop(1)),
		},
		"var": {
			{`[^\S\n]+`, Text, nil},
			{`\b(((re)?source|barrel)Path|excludeAnnotations|annotations|lang)\b`, NameBuiltin, nil},
			{`\bbase\b`, NameConstant, nil},
			{`\b(ind|zsm|hrv|ces|dan|dut|eng|fin|fre|deu|gre|hun|ita|nob|po[lr]|rus|sl[ov]|spa|swe|ara|heb|zh[st]|jpn|kor|tha|vie|bul|tur)`, NameConstant, nil},
			{`\b((semi)?round|rectangle)(-\d+x\d+)?\b`, NameConstant, nil},
			{`[\.;\[\]\(\$]`, Punctuation, nil},
			{`\)`, Punctuation, Pop(1)},
			{`[a-zA-Z_]\w*`, Name, nil},
			Default(Pop(1)),
		},
	}
}
