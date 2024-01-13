package d

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// D lexer. https://dlang.org/spec/lex.html
var D = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "D",
		Aliases:   []string{"d"},
		Filenames: []string{"*.d", "*.di"},
		MimeTypes: []string{"text/x-d"},
		EnsureNL:  true,
	},
	dRules,
))

func dRules() Rules {
	return Rules{
		"root": {
			{`[^\S\n]+`, Text, nil},

			// https://dlang.org/spec/lex.html#comment
			{`//.*?\n`, CommentSingle, nil},
			{`/\*.*?\*/`, CommentMultiline, nil},
			{`/\+.*?\+/`, CommentMultiline, nil},

			// https://dlang.org/spec/lex.html#keywords
			{`(asm|assert|body|break|case|cast|catch|continue|default|debug|delete|deprecated|do|else|finally|for|foreach|foreach_reverse|goto|if|in|invariant|is|macro|mixin|new|out|pragma|return|super|switch|this|throw|try|version|while|with)\b`, Keyword, nil},
			{`__(FILE|FILE_FULL_PATH|MODULE|LINE|FUNCTION|PRETTY_FUNCTION|DATE|EOF|TIME|TIMESTAMP|VENDOR|VERSION)__\b`, NameBuiltin, nil},
			{`__(traits|vector|parameters)\b`, NameBuiltin, nil},
			{`((?:(?:[^\W\d]|\$)[\w.\[\]$<>]*\s+)+?)((?:[^\W\d]|\$)[\w$]*)(\s*)(\()`, ByGroups(UsingSelf("root"), NameFunction, Text, Operator), nil},

			// https://dlang.org/spec/attribute.html#uda
			{`@[\w.]*`, NameDecorator, nil},
			{`(abstract|auto|alias|align|const|delegate|enum|export|final|function|inout|lazy|nothrow|override|package|private|protected|public|pure|static|synchronized|template|volatile|__gshared)\b`, KeywordDeclaration, nil},

			// https://dlang.org/spec/type.html#basic-data-types
			{`(void|bool|byte|ubyte|short|ushort|int|uint|long|ulong|cent|ucent|float|double|real|ifloat|idouble|ireal|cfloat|cdouble|creal|char|wchar|dchar|string|wstring|dstring)\b`, KeywordType, nil},
			{`(module)(\s+)`, ByGroups(KeywordNamespace, Text), Push("import")},
			{`(true|false|null)\b`, KeywordConstant, nil},
			{`(class|interface|struct|template|union)(\s+)`, ByGroups(KeywordDeclaration, Text), Push("class")},
			{`(import)(\s+)`, ByGroups(KeywordNamespace, Text), Push("import")},

			// https://dlang.org/spec/lex.html#string_literals
			// TODO support delimited strings
			{`[qr]?"(\\\\|\\"|[^"])*"[cwd]?`, LiteralString, nil},
			{"(`)([^`]*)(`)[cwd]?", LiteralString, nil},
			{`'\\.'|'[^\\]'|'\\u[0-9a-fA-F]{4}'`, LiteralStringChar, nil},
			{`(\.)((?:[^\W\d]|\$)[\w$]*)`, ByGroups(Operator, NameAttribute), nil},
			{`^\s*([^\W\d]|\$)[\w$]*:`, NameLabel, nil},

			// https://dlang.org/spec/lex.html#floatliteral
			{`([0-9][0-9_]*\.([0-9][0-9_]*)?|\.[0-9][0-9_]*)([eE][+\-]?[0-9][0-9_]*)?[fFL]?i?|[0-9][eE][+\-]?[0-9][0-9_]*[fFL]?|[0-9]([eE][+\-]?[0-9][0-9_]*)?[fFL]|0[xX]([0-9a-fA-F][0-9a-fA-F_]*\.?|([0-9a-fA-F][0-9a-fA-F_]*)?\.[0-9a-fA-F][0-9a-fA-F_]*)[pP][+\-]?[0-9][0-9_]*[fFL]?`, LiteralNumberFloat, nil},
			// https://dlang.org/spec/lex.html#integerliteral
			{`0[xX][0-9a-fA-F][0-9a-fA-F_]*[lL]?`, LiteralNumberHex, nil},
			{`0[bB][01][01_]*[lL]?`, LiteralNumberBin, nil},
			{`0[0-7_]+[lL]?`, LiteralNumberOct, nil},
			{`0|[1-9][0-9_]*[lL]?`, LiteralNumberInteger, nil},
			{`([~^*!%&\[\](){}<>|+=:;,./?-]|q{)`, Operator, nil},
			{`([^\W\d]|\$)[\w$]*`, Name, nil},
			{`\n`, Text, nil},
		},
		"class": {
			{`([^\W\d]|\$)[\w$]*`, NameClass, Pop(1)},
		},
		"import": {
			{`[\w.]+\*?`, NameNamespace, Pop(1)},
		},
	}
}
