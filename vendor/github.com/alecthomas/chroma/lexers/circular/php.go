package circular

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// PHP lexer for pure PHP code (not embedded in HTML).
var PHP = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "PHP",
		Aliases:         []string{"php", "php3", "php4", "php5"},
		Filenames:       []string{"*.php", "*.php[345]", "*.inc"},
		MimeTypes:       []string{"text/x-php"},
		DotAll:          true,
		CaseInsensitive: true,
		EnsureNL:        true,
	},
	phpRules,
))

func phpRules() Rules {
	return phpCommonRules().Rename("php", "root")
}

func phpCommonRules() Rules {
	return Rules{
		"php": {
			{`\?>`, CommentPreproc, Pop(1)},
			{`(<<<)([\'"]?)((?:[\\_a-z]|[^\x00-\x7f])(?:[\\\w]|[^\x00-\x7f])*)(\2\n.*?\n\s*)(\3)(;?)(\n)`, ByGroups(LiteralString, LiteralString, LiteralStringDelimiter, LiteralString, LiteralStringDelimiter, Punctuation, Text), nil},
			{`\s+`, Text, nil},
			{`#.*?\n`, CommentSingle, nil},
			{`//.*?\n`, CommentSingle, nil},
			{`/\*\*/`, CommentMultiline, nil},
			{`/\*\*.*?\*/`, LiteralStringDoc, nil},
			{`/\*.*?\*/`, CommentMultiline, nil},
			{`(->|::)(\s*)((?:[\\_a-z]|[^\x00-\x7f])(?:[\\\w]|[^\x00-\x7f])*)`, ByGroups(Operator, Text, NameAttribute), nil},
			{`[~!%^&*+=|:.<>/@-]+`, Operator, nil},
			{`\?`, Operator, nil},
			{`[\[\]{}();,]+`, Punctuation, nil},
			{`(class)(\s+)`, ByGroups(Keyword, Text), Push("classname")},
			{`(function)(\s*)(?=\()`, ByGroups(Keyword, Text), nil},
			{`(function)(\s+)(&?)(\s*)`, ByGroups(Keyword, Text, Operator, Text), Push("functionname")},
			{`(const)(\s+)((?:[\\_a-z]|[^\x00-\x7f])(?:[\\\w]|[^\x00-\x7f])*)`, ByGroups(Keyword, Text, NameConstant), nil},
			{`(and|E_PARSE|old_function|E_ERROR|or|as|E_WARNING|parent|eval|PHP_OS|break|exit|case|extends|PHP_VERSION|cfunction|FALSE|print|for|require|continue|foreach|require_once|declare|return|default|static|do|switch|die|stdClass|echo|else|TRUE|elseif|var|empty|if|xor|enddeclare|include|virtual|endfor|include_once|while|endforeach|global|endif|list|endswitch|new|endwhile|not|array|E_ALL|NULL|final|php_user_filter|interface|implements|public|private|protected|abstract|clone|try|catch|throw|this|use|namespace|trait|yield|finally)\b`, Keyword, nil},
			{`(true|false|null)\b`, KeywordConstant, nil},
			Include("magicconstants"),
			{`\$\{\$+(?:[\\_a-z]|[^\x00-\x7f])(?:[\\\w]|[^\x00-\x7f])*\}`, NameVariable, nil},
			{`\$+(?:[\\_a-z]|[^\x00-\x7f])(?:[\\\w]|[^\x00-\x7f])*`, NameVariable, nil},
			{`(?:[\\_a-z]|[^\x00-\x7f])(?:[\\\w]|[^\x00-\x7f])*`, NameOther, nil},
			{`(\d+\.\d*|\d*\.\d+)(e[+-]?[0-9]+)?`, LiteralNumberFloat, nil},
			{`\d+e[+-]?[0-9]+`, LiteralNumberFloat, nil},
			{`0[0-7]+`, LiteralNumberOct, nil},
			{`0x[a-f0-9_]+`, LiteralNumberHex, nil},
			{`\d[\d_]*`, LiteralNumberInteger, nil},
			{`0b[01]+`, LiteralNumberBin, nil},
			{`'([^'\\]*(?:\\.[^'\\]*)*)'`, LiteralStringSingle, nil},
			{"`([^`\\\\]*(?:\\\\.[^`\\\\]*)*)`", LiteralStringBacktick, nil},
			{`"`, LiteralStringDouble, Push("string")},
		},
		"magicfuncs": {
			{Words(``, `\b`, `__construct`, `__destruct`, `__call`, `__callStatic`, `__get`, `__set`, `__isset`, `__unset`, `__sleep`, `__wakeup`, `__toString`, `__invoke`, `__set_state`, `__clone`, `__debugInfo`), NameFunctionMagic, nil},
		},
		"magicconstants": {
			{Words(``, `\b`, `__LINE__`, `__FILE__`, `__DIR__`, `__FUNCTION__`, `__CLASS__`, `__TRAIT__`, `__METHOD__`, `__NAMESPACE__`), NameConstant, nil},
		},
		"classname": {
			{`(?:[\\_a-z]|[^\x00-\x7f])(?:[\\\w]|[^\x00-\x7f])*`, NameClass, Pop(1)},
		},
		"functionname": {
			Include("magicfuncs"),
			{`(?:[\\_a-z]|[^\x00-\x7f])(?:[\\\w]|[^\x00-\x7f])*`, NameFunction, Pop(1)},
			Default(Pop(1)),
		},
		"string": {
			{`"`, LiteralStringDouble, Pop(1)},
			{`[^{$"\\]+`, LiteralStringDouble, nil},
			{`\\([nrt"$\\]|[0-7]{1,3}|x[0-9a-f]{1,2})`, LiteralStringEscape, nil},
			{`\$(?:[\\_a-z]|[^\x00-\x7f])(?:[\\\w]|[^\x00-\x7f])*(\[\S+?\]|->(?:[\\_a-z]|[^\x00-\x7f])(?:[\\\w]|[^\x00-\x7f])*)?`, LiteralStringInterpol, nil},
			{`(\{\$\{)(.*?)(\}\})`, ByGroups(LiteralStringInterpol, UsingSelf("root"), LiteralStringInterpol), nil},
			{`(\{)(\$.*?)(\})`, ByGroups(LiteralStringInterpol, UsingSelf("root"), LiteralStringInterpol), nil},
			{`(\$\{)(\S+)(\})`, ByGroups(LiteralStringInterpol, NameVariable, LiteralStringInterpol), nil},
			{`[${\\]`, LiteralStringDouble, nil},
		},
	}
}
