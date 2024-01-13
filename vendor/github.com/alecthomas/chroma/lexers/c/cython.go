package c

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Cython lexer.
var Cython = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Cython",
		Aliases:   []string{"cython", "pyx", "pyrex"},
		Filenames: []string{"*.pyx", "*.pxd", "*.pxi"},
		MimeTypes: []string{"text/x-cython", "application/x-cython"},
	},
	cythonRules,
))

func cythonRules() Rules {
	return Rules{
		"root": {
			{`\n`, Text, nil},
			{`^(\s*)("""(?:.|\n)*?""")`, ByGroups(Text, LiteralStringDoc), nil},
			{`^(\s*)('''(?:.|\n)*?''')`, ByGroups(Text, LiteralStringDoc), nil},
			{`[^\S\n]+`, Text, nil},
			{`#.*$`, Comment, nil},
			{`[]{}:(),;[]`, Punctuation, nil},
			{`\\\n`, Text, nil},
			{`\\`, Text, nil},
			{`(in|is|and|or|not)\b`, OperatorWord, nil},
			{`(<)([a-zA-Z0-9.?]+)(>)`, ByGroups(Punctuation, KeywordType, Punctuation), nil},
			{`!=|==|<<|>>|[-~+/*%=<>&^|.?]`, Operator, nil},
			{`(from)(\d+)(<=)(\s+)(<)(\d+)(:)`, ByGroups(Keyword, LiteralNumberInteger, Operator, Name, Operator, Name, Punctuation), nil},
			Include("keywords"),
			{`(def|property)(\s+)`, ByGroups(Keyword, Text), Push("funcname")},
			{`(cp?def)(\s+)`, ByGroups(Keyword, Text), Push("cdef")},
			{`(cdef)(:)`, ByGroups(Keyword, Punctuation), nil},
			{`(class|struct)(\s+)`, ByGroups(Keyword, Text), Push("classname")},
			{`(from)(\s+)`, ByGroups(Keyword, Text), Push("fromimport")},
			{`(c?import)(\s+)`, ByGroups(Keyword, Text), Push("import")},
			Include("builtins"),
			Include("backtick"),
			{`(?:[rR]|[uU][rR]|[rR][uU])"""`, LiteralString, Push("tdqs")},
			{`(?:[rR]|[uU][rR]|[rR][uU])'''`, LiteralString, Push("tsqs")},
			{`(?:[rR]|[uU][rR]|[rR][uU])"`, LiteralString, Push("dqs")},
			{`(?:[rR]|[uU][rR]|[rR][uU])'`, LiteralString, Push("sqs")},
			{`[uU]?"""`, LiteralString, Combined("stringescape", "tdqs")},
			{`[uU]?'''`, LiteralString, Combined("stringescape", "tsqs")},
			{`[uU]?"`, LiteralString, Combined("stringescape", "dqs")},
			{`[uU]?'`, LiteralString, Combined("stringescape", "sqs")},
			Include("name"),
			Include("numbers"),
		},
		"keywords": {
			{Words(``, `\b`, `assert`, `break`, `by`, `continue`, `ctypedef`, `del`, `elif`, `else`, `except`, `except?`, `exec`, `finally`, `for`, `fused`, `gil`, `global`, `if`, `include`, `lambda`, `nogil`, `pass`, `print`, `raise`, `return`, `try`, `while`, `yield`, `as`, `with`), Keyword, nil},
			{`(DEF|IF|ELIF|ELSE)\b`, CommentPreproc, nil},
		},
		"builtins": {
			{Words(`(?<!\.)`, `\b`, `__import__`, `abs`, `all`, `any`, `apply`, `basestring`, `bin`, `bool`, `buffer`, `bytearray`, `bytes`, `callable`, `chr`, `classmethod`, `cmp`, `coerce`, `compile`, `complex`, `delattr`, `dict`, `dir`, `divmod`, `enumerate`, `eval`, `execfile`, `exit`, `file`, `filter`, `float`, `frozenset`, `getattr`, `globals`, `hasattr`, `hash`, `hex`, `id`, `input`, `int`, `intern`, `isinstance`, `issubclass`, `iter`, `len`, `list`, `locals`, `long`, `map`, `max`, `min`, `next`, `object`, `oct`, `open`, `ord`, `pow`, `property`, `range`, `raw_input`, `reduce`, `reload`, `repr`, `reversed`, `round`, `set`, `setattr`, `slice`, `sorted`, `staticmethod`, `str`, `sum`, `super`, `tuple`, `type`, `unichr`, `unicode`, `unsigned`, `vars`, `xrange`, `zip`), NameBuiltin, nil},
			{`(?<!\.)(self|None|Ellipsis|NotImplemented|False|True|NULL)\b`, NameBuiltinPseudo, nil},
			{Words(`(?<!\.)`, `\b`, `ArithmeticError`, `AssertionError`, `AttributeError`, `BaseException`, `DeprecationWarning`, `EOFError`, `EnvironmentError`, `Exception`, `FloatingPointError`, `FutureWarning`, `GeneratorExit`, `IOError`, `ImportError`, `ImportWarning`, `IndentationError`, `IndexError`, `KeyError`, `KeyboardInterrupt`, `LookupError`, `MemoryError`, `NameError`, `NotImplemented`, `NotImplementedError`, `OSError`, `OverflowError`, `OverflowWarning`, `PendingDeprecationWarning`, `ReferenceError`, `RuntimeError`, `RuntimeWarning`, `StandardError`, `StopIteration`, `SyntaxError`, `SyntaxWarning`, `SystemError`, `SystemExit`, `TabError`, `TypeError`, `UnboundLocalError`, `UnicodeDecodeError`, `UnicodeEncodeError`, `UnicodeError`, `UnicodeTranslateError`, `UnicodeWarning`, `UserWarning`, `ValueError`, `Warning`, `ZeroDivisionError`), NameException, nil},
		},
		"numbers": {
			{`(\d+\.?\d*|\d*\.\d+)([eE][+-]?[0-9]+)?`, LiteralNumberFloat, nil},
			{`0\d+`, LiteralNumberOct, nil},
			{`0[xX][a-fA-F0-9]+`, LiteralNumberHex, nil},
			{`\d+L`, LiteralNumberIntegerLong, nil},
			{`\d+`, LiteralNumberInteger, nil},
		},
		"backtick": {
			{"`.*?`", LiteralStringBacktick, nil},
		},
		"name": {
			{`@\w+`, NameDecorator, nil},
			{`[a-zA-Z_]\w*`, Name, nil},
		},
		"funcname": {
			{`[a-zA-Z_]\w*`, NameFunction, Pop(1)},
		},
		"cdef": {
			{`(public|readonly|extern|api|inline)\b`, KeywordReserved, nil},
			{`(struct|enum|union|class)\b`, Keyword, nil},
			{`([a-zA-Z_]\w*)(\s*)(?=[(:#=]|$)`, ByGroups(NameFunction, Text), Pop(1)},
			{`([a-zA-Z_]\w*)(\s*)(,)`, ByGroups(NameFunction, Text, Punctuation), nil},
			{`from\b`, Keyword, Pop(1)},
			{`as\b`, Keyword, nil},
			{`:`, Punctuation, Pop(1)},
			{`(?=["\'])`, Text, Pop(1)},
			{`[a-zA-Z_]\w*`, KeywordType, nil},
			{`.`, Text, nil},
		},
		"classname": {
			{`[a-zA-Z_]\w*`, NameClass, Pop(1)},
		},
		"import": {
			{`(\s+)(as)(\s+)`, ByGroups(Text, Keyword, Text), nil},
			{`[a-zA-Z_][\w.]*`, NameNamespace, nil},
			{`(\s*)(,)(\s*)`, ByGroups(Text, Operator, Text), nil},
			Default(Pop(1)),
		},
		"fromimport": {
			{`(\s+)(c?import)\b`, ByGroups(Text, Keyword), Pop(1)},
			{`[a-zA-Z_.][\w.]*`, NameNamespace, nil},
			Default(Pop(1)),
		},
		"stringescape": {
			{`\\([\\abfnrtv"\']|\n|N\{.*?\}|u[a-fA-F0-9]{4}|U[a-fA-F0-9]{8}|x[a-fA-F0-9]{2}|[0-7]{1,3})`, LiteralStringEscape, nil},
		},
		"strings": {
			{`%(\([a-zA-Z0-9]+\))?[-#0 +]*([0-9]+|[*])?(\.([0-9]+|[*]))?[hlL]?[E-GXc-giorsux%]`, LiteralStringInterpol, nil},
			{`[^\\\'"%\n]+`, LiteralString, nil},
			{`[\'"\\]`, LiteralString, nil},
			{`%`, LiteralString, nil},
		},
		"nl": {
			{`\n`, LiteralString, nil},
		},
		"dqs": {
			{`"`, LiteralString, Pop(1)},
			{`\\\\|\\"|\\\n`, LiteralStringEscape, nil},
			Include("strings"),
		},
		"sqs": {
			{`'`, LiteralString, Pop(1)},
			{`\\\\|\\'|\\\n`, LiteralStringEscape, nil},
			Include("strings"),
		},
		"tdqs": {
			{`"""`, LiteralString, Pop(1)},
			Include("strings"),
			Include("nl"),
		},
		"tsqs": {
			{`'''`, LiteralString, Pop(1)},
			Include("strings"),
			Include("nl"),
		},
	}
}
