package o

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// 1S:Enterprise lexer.
var OnesEnterprise = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "OnesEnterprise",
		Aliases:         []string{"ones", "onesenterprise", "1S", "1S:Enterprise"},
		Filenames:       []string{"*.EPF", "*.epf", "*.ERF", "*.erf"},
		MimeTypes:       []string{"application/octet-stream"},
		CaseInsensitive: true,
	},
	onesRules,
))

func onesRules() Rules {
	return Rules{
		"root": {
			{`\n`, Text, nil},
			{`\s+`, Text, nil},
			{`\\\n`, Text, nil},
			{`[^\S\n]+`, Text, nil},
			{`//(.*?)\n`, Comment, nil},
			{`(#область|#region|#конецобласти|#endregion|#если|#if|#иначе|#else|#конецесли|#endif).*`, CommentPreproc, nil},
			{`(&наклиенте|&atclient|&насервере|&atserver|&насерверебезконтекста|&atservernocontext|&наклиентенасерверебезконтекста|&atclientatservernocontext).*`, CommentPreproc, nil},
			{`(>=|<=|<>|\+|-|=|>|<|\*|/|%)`, Operator, nil},
			{`(;|,|\)|\(|\.)`, Punctuation, nil},
			{Words(``, `\b`, `истина`, `true`, `ложь`, `false`, `и`, `and`, `или`, `or`, `не`, `not`), Operator, nil},
			{Words(``, `\b`, `если`, `if`, `тогда`, `then`, `иначе`, `else`, `иначеесли`, `elsif`, `конецесли`, `endif`), Operator, nil},
			{Words(``, `\b`, `для`, `for`, `каждого`, `each`, `из`, `in`, `цикл`, `do`, `пока`, `while`, `конеццикла`, `enddo`, `по`, `to`), Operator, nil},
			{Words(``, `\b`, `прервать`, `break`, `продолжить`, `continue`, `возврат`, `return`, `перейти`, `goto`), Operator, nil},
			{Words(``, `\b`, `процедура`, `procedure`, `конецпроцедуры`, `endprocedure`, `функция`, `function`, `конецфункции`, `endfunction`), Keyword, nil},
			{Words(``, `\b`, `новый`, `new`, `знач`, `val`, `экспорт`, `export`, `перем`, `var`), Keyword, nil},
			{Words(``, `\b`, `попытка`, `try`, `исключение`, `except`, `вызватьисключение`, `raise`, `конецпопытки`, `endtry`), Keyword, nil},
			{Words(``, `\b`, `выполнить`, `execute`, `вычислить`, `eval`), Keyword, nil},
			{`"`, LiteralString, Push("string")},
			{`[_а-яА-Я0-9][а-яА-Я0-9]*`, Name, nil},
			{`[_\w][\w]*`, Name, nil},
		},
		"string": {
			{`""`, LiteralString, nil},
			{`"C?`, LiteralString, Pop(1)},
			{`[^"]+`, LiteralString, nil},
		},
	}
}
