package h

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Hy lexer.
var Hy = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Hy",
		Aliases:   []string{"hylang"},
		Filenames: []string{"*.hy"},
		MimeTypes: []string{"text/x-hy", "application/x-hy"},
	},
	hyRules,
))

func hyRules() Rules {
	return Rules{
		"root": {
			{`;.*$`, CommentSingle, nil},
			{`[,\s]+`, Text, nil},
			{`-?\d+\.\d+`, LiteralNumberFloat, nil},
			{`-?\d+`, LiteralNumberInteger, nil},
			{`0[0-7]+j?`, LiteralNumberOct, nil},
			{`0[xX][a-fA-F0-9]+`, LiteralNumberHex, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralString, nil},
			{`'(?!#)[\w!$%*+<=>?/.#-]+`, LiteralStringSymbol, nil},
			{`\\(.|[a-z]+)`, LiteralStringChar, nil},
			{`^(\s*)([rRuU]{,2}"""(?:.|\n)*?""")`, ByGroups(Text, LiteralStringDoc), nil},
			{`^(\s*)([rRuU]{,2}'''(?:.|\n)*?''')`, ByGroups(Text, LiteralStringDoc), nil},
			{`::?(?!#)[\w!$%*+<=>?/.#-]+`, LiteralStringSymbol, nil},
			{"~@|[`\\'#^~&@]", Operator, nil},
			Include("py-keywords"),
			Include("py-builtins"),
			{Words(``, ` `, `cond`, `for`, `->`, `->>`, `car`, `cdr`, `first`, `rest`, `let`, `when`, `unless`, `import`, `do`, `progn`, `get`, `slice`, `assoc`, `with-decorator`, `,`, `list_comp`, `kwapply`, `~`, `is`, `in`, `is-not`, `not-in`, `quasiquote`, `unquote`, `unquote-splice`, `quote`, `|`, `<<=`, `>>=`, `foreach`, `while`, `eval-and-compile`, `eval-when-compile`), Keyword, nil},
			{Words(``, ` `, `def`, `defn`, `defun`, `defmacro`, `defclass`, `lambda`, `fn`, `setv`), KeywordDeclaration, nil},
			{Words(``, ` `, `cycle`, `dec`, `distinct`, `drop`, `even?`, `filter`, `inc`, `instance?`, `iterable?`, `iterate`, `iterator?`, `neg?`, `none?`, `nth`, `numeric?`, `odd?`, `pos?`, `remove`, `repeat`, `repeatedly`, `take`, `take_nth`, `take_while`, `zero?`), NameBuiltin, nil},
			{`(?<=\()(?!#)[\w!$%*+<=>?/.#-]+`, NameFunction, nil},
			{`(?!#)[\w!$%*+<=>?/.#-]+`, NameVariable, nil},
			{`(\[|\])`, Punctuation, nil},
			{`(\{|\})`, Punctuation, nil},
			{`(\(|\))`, Punctuation, nil},
		},
		"py-keywords": {
			{Words(``, `\b`, `assert`, `break`, `continue`, `del`, `elif`, `else`, `except`, `exec`, `finally`, `for`, `global`, `if`, `lambda`, `pass`, `print`, `raise`, `return`, `try`, `while`, `yield`, `yield from`, `as`, `with`), Keyword, nil},
		},
		"py-builtins": {
			{Words(`(?<!\.)`, `\b`, `__import__`, `abs`, `all`, `any`, `apply`, `basestring`, `bin`, `bool`, `buffer`, `bytearray`, `bytes`, `callable`, `chr`, `classmethod`, `cmp`, `coerce`, `compile`, `complex`, `delattr`, `dict`, `dir`, `divmod`, `enumerate`, `eval`, `execfile`, `exit`, `file`, `filter`, `float`, `frozenset`, `getattr`, `globals`, `hasattr`, `hash`, `hex`, `id`, `input`, `int`, `intern`, `isinstance`, `issubclass`, `iter`, `len`, `list`, `locals`, `long`, `map`, `max`, `min`, `next`, `object`, `oct`, `open`, `ord`, `pow`, `property`, `range`, `raw_input`, `reduce`, `reload`, `repr`, `reversed`, `round`, `set`, `setattr`, `slice`, `sorted`, `staticmethod`, `str`, `sum`, `super`, `tuple`, `type`, `unichr`, `unicode`, `vars`, `xrange`, `zip`), NameBuiltin, nil},
			{`(?<!\.)(self|None|Ellipsis|NotImplemented|False|True|cls)\b`, NameBuiltinPseudo, nil},
			{Words(`(?<!\.)`, `\b`, `ArithmeticError`, `AssertionError`, `AttributeError`, `BaseException`, `DeprecationWarning`, `EOFError`, `EnvironmentError`, `Exception`, `FloatingPointError`, `FutureWarning`, `GeneratorExit`, `IOError`, `ImportError`, `ImportWarning`, `IndentationError`, `IndexError`, `KeyError`, `KeyboardInterrupt`, `LookupError`, `MemoryError`, `NameError`, `NotImplemented`, `NotImplementedError`, `OSError`, `OverflowError`, `OverflowWarning`, `PendingDeprecationWarning`, `ReferenceError`, `RuntimeError`, `RuntimeWarning`, `StandardError`, `StopIteration`, `SyntaxError`, `SyntaxWarning`, `SystemError`, `SystemExit`, `TabError`, `TypeError`, `UnboundLocalError`, `UnicodeDecodeError`, `UnicodeEncodeError`, `UnicodeError`, `UnicodeTranslateError`, `UnicodeWarning`, `UserWarning`, `ValueError`, `VMSError`, `Warning`, `WindowsError`, `ZeroDivisionError`), NameException, nil},
		},
	}
}
