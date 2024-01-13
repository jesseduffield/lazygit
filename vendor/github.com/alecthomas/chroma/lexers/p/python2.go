package p

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Python2 lexer.
var Python2 = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Python 2",
		Aliases:   []string{"python2", "py2"},
		Filenames: []string{},
		MimeTypes: []string{"text/x-python2", "application/x-python2"},
	},
	python2Rules,
))

func python2Rules() Rules {
	return Rules{
		"root": {
			{`\n`, Text, nil},
			{`^(\s*)([rRuUbB]{,2})("""(?:.|\n)*?""")`, ByGroups(Text, LiteralStringAffix, LiteralStringDoc), nil},
			{`^(\s*)([rRuUbB]{,2})('''(?:.|\n)*?''')`, ByGroups(Text, LiteralStringAffix, LiteralStringDoc), nil},
			{`[^\S\n]+`, Text, nil},
			{`\A#!.+$`, CommentHashbang, nil},
			{`#.*$`, CommentSingle, nil},
			{`[]{}:(),;[]`, Punctuation, nil},
			{`\\\n`, Text, nil},
			{`\\`, Text, nil},
			{`(in|is|and|or|not)\b`, OperatorWord, nil},
			{`!=|==|<<|>>|[-~+/*%=<>&^|.]`, Operator, nil},
			Include("keywords"),
			{`(def)((?:\s|\\\s)+)`, ByGroups(Keyword, Text), Push("funcname")},
			{`(class)((?:\s|\\\s)+)`, ByGroups(Keyword, Text), Push("classname")},
			{`(from)((?:\s|\\\s)+)`, ByGroups(KeywordNamespace, Text), Push("fromimport")},
			{`(import)((?:\s|\\\s)+)`, ByGroups(KeywordNamespace, Text), Push("import")},
			Include("builtins"),
			Include("magicfuncs"),
			Include("magicvars"),
			Include("backtick"),
			{`([rR]|[uUbB][rR]|[rR][uUbB])(""")`, ByGroups(LiteralStringAffix, LiteralStringDouble), Push("tdqs")},
			{`([rR]|[uUbB][rR]|[rR][uUbB])(''')`, ByGroups(LiteralStringAffix, LiteralStringSingle), Push("tsqs")},
			{`([rR]|[uUbB][rR]|[rR][uUbB])(")`, ByGroups(LiteralStringAffix, LiteralStringDouble), Push("dqs")},
			{`([rR]|[uUbB][rR]|[rR][uUbB])(')`, ByGroups(LiteralStringAffix, LiteralStringSingle), Push("sqs")},
			{`([uUbB]?)(""")`, ByGroups(LiteralStringAffix, LiteralStringDouble), Combined("stringescape", "tdqs")},
			{`([uUbB]?)(''')`, ByGroups(LiteralStringAffix, LiteralStringSingle), Combined("stringescape", "tsqs")},
			{`([uUbB]?)(")`, ByGroups(LiteralStringAffix, LiteralStringDouble), Combined("stringescape", "dqs")},
			{`([uUbB]?)(')`, ByGroups(LiteralStringAffix, LiteralStringSingle), Combined("stringescape", "sqs")},
			Include("name"),
			Include("numbers"),
		},
		"keywords": {
			{Words(``, `\b`, `assert`, `break`, `continue`, `del`, `elif`, `else`, `except`, `exec`, `finally`, `for`, `global`, `if`, `lambda`, `pass`, `print`, `raise`, `return`, `try`, `while`, `yield`, `yield from`, `as`, `with`), Keyword, nil},
		},
		"builtins": {
			{Words(`(?<!\.)`, `\b`, `__import__`, `abs`, `all`, `any`, `apply`, `basestring`, `bin`, `bool`, `buffer`, `bytearray`, `bytes`, `callable`, `chr`, `classmethod`, `cmp`, `coerce`, `compile`, `complex`, `delattr`, `dict`, `dir`, `divmod`, `enumerate`, `eval`, `execfile`, `exit`, `file`, `filter`, `float`, `frozenset`, `getattr`, `globals`, `hasattr`, `hash`, `hex`, `id`, `input`, `int`, `intern`, `isinstance`, `issubclass`, `iter`, `len`, `list`, `locals`, `long`, `map`, `max`, `min`, `next`, `object`, `oct`, `open`, `ord`, `pow`, `property`, `range`, `raw_input`, `reduce`, `reload`, `repr`, `reversed`, `round`, `set`, `setattr`, `slice`, `sorted`, `staticmethod`, `str`, `sum`, `super`, `tuple`, `type`, `unichr`, `unicode`, `vars`, `xrange`, `zip`), NameBuiltin, nil},
			{`(?<!\.)(self|None|Ellipsis|NotImplemented|False|True|cls)\b`, NameBuiltinPseudo, nil},
			{Words(`(?<!\.)`, `\b`, `ArithmeticError`, `AssertionError`, `AttributeError`, `BaseException`, `DeprecationWarning`, `EOFError`, `EnvironmentError`, `Exception`, `FloatingPointError`, `FutureWarning`, `GeneratorExit`, `IOError`, `ImportError`, `ImportWarning`, `IndentationError`, `IndexError`, `KeyError`, `KeyboardInterrupt`, `LookupError`, `MemoryError`, `NameError`, `NotImplementedError`, `OSError`, `OverflowError`, `OverflowWarning`, `PendingDeprecationWarning`, `ReferenceError`, `RuntimeError`, `RuntimeWarning`, `StandardError`, `StopIteration`, `SyntaxError`, `SyntaxWarning`, `SystemError`, `SystemExit`, `TabError`, `TypeError`, `UnboundLocalError`, `UnicodeDecodeError`, `UnicodeEncodeError`, `UnicodeError`, `UnicodeTranslateError`, `UnicodeWarning`, `UserWarning`, `ValueError`, `VMSError`, `Warning`, `WindowsError`, `ZeroDivisionError`), NameException, nil},
		},
		"magicfuncs": {
			{Words(``, `\b`, `__abs__`, `__add__`, `__and__`, `__call__`, `__cmp__`, `__coerce__`, `__complex__`, `__contains__`, `__del__`, `__delattr__`, `__delete__`, `__delitem__`, `__delslice__`, `__div__`, `__divmod__`, `__enter__`, `__eq__`, `__exit__`, `__float__`, `__floordiv__`, `__ge__`, `__get__`, `__getattr__`, `__getattribute__`, `__getitem__`, `__getslice__`, `__gt__`, `__hash__`, `__hex__`, `__iadd__`, `__iand__`, `__idiv__`, `__ifloordiv__`, `__ilshift__`, `__imod__`, `__imul__`, `__index__`, `__init__`, `__instancecheck__`, `__int__`, `__invert__`, `__iop__`, `__ior__`, `__ipow__`, `__irshift__`, `__isub__`, `__iter__`, `__itruediv__`, `__ixor__`, `__le__`, `__len__`, `__long__`, `__lshift__`, `__lt__`, `__missing__`, `__mod__`, `__mul__`, `__ne__`, `__neg__`, `__new__`, `__nonzero__`, `__oct__`, `__op__`, `__or__`, `__pos__`, `__pow__`, `__radd__`, `__rand__`, `__rcmp__`, `__rdiv__`, `__rdivmod__`, `__repr__`, `__reversed__`, `__rfloordiv__`, `__rlshift__`, `__rmod__`, `__rmul__`, `__rop__`, `__ror__`, `__rpow__`, `__rrshift__`, `__rshift__`, `__rsub__`, `__rtruediv__`, `__rxor__`, `__set__`, `__setattr__`, `__setitem__`, `__setslice__`, `__str__`, `__sub__`, `__subclasscheck__`, `__truediv__`, `__unicode__`, `__xor__`), NameFunctionMagic, nil},
		},
		"magicvars": {
			{Words(``, `\b`, `__bases__`, `__class__`, `__closure__`, `__code__`, `__defaults__`, `__dict__`, `__doc__`, `__file__`, `__func__`, `__globals__`, `__metaclass__`, `__module__`, `__mro__`, `__name__`, `__self__`, `__slots__`, `__weakref__`), NameVariableMagic, nil},
		},
		"numbers": {
			{`(\d+\.\d*|\d*\.\d+)([eE][+-]?[0-9]+)?j?`, LiteralNumberFloat, nil},
			{`\d+[eE][+-]?[0-9]+j?`, LiteralNumberFloat, nil},
			{`0[0-7]+j?`, LiteralNumberOct, nil},
			{`0[bB][01]+`, LiteralNumberBin, nil},
			{`0[xX][a-fA-F0-9]+`, LiteralNumberHex, nil},
			{`\d+L`, LiteralNumberIntegerLong, nil},
			{`\d+j?`, LiteralNumberInteger, nil},
		},
		"backtick": {
			{"`.*?`", LiteralStringBacktick, nil},
		},
		"name": {
			{`@[\w.]+`, NameDecorator, nil},
			{`[a-zA-Z_]\w*`, Name, nil},
		},
		"funcname": {
			Include("magicfuncs"),
			{`[a-zA-Z_]\w*`, NameFunction, Pop(1)},
			Default(Pop(1)),
		},
		"classname": {
			{`[a-zA-Z_]\w*`, NameClass, Pop(1)},
		},
		"import": {
			{`(?:[ \t]|\\\n)+`, Text, nil},
			{`as\b`, KeywordNamespace, nil},
			{`,`, Operator, nil},
			{`[a-zA-Z_][\w.]*`, NameNamespace, nil},
			Default(Pop(1)),
		},
		"fromimport": {
			{`(?:[ \t]|\\\n)+`, Text, nil},
			{`import\b`, KeywordNamespace, Pop(1)},
			{`None\b`, NameBuiltinPseudo, Pop(1)},
			{`[a-zA-Z_.][\w.]*`, NameNamespace, nil},
			Default(Pop(1)),
		},
		"stringescape": {
			{`\\([\\abfnrtv"\']|\n|N\{.*?\}|u[a-fA-F0-9]{4}|U[a-fA-F0-9]{8}|x[a-fA-F0-9]{2}|[0-7]{1,3})`, LiteralStringEscape, nil},
		},
		"strings-single": {
			{`%(\(\w+\))?[-#0 +]*([0-9]+|[*])?(\.([0-9]+|[*]))?[hlL]?[E-GXc-giorsux%]`, LiteralStringInterpol, nil},
			{`[^\\\'"%\n]+`, LiteralStringSingle, nil},
			{`[\'"\\]`, LiteralStringSingle, nil},
			{`%`, LiteralStringSingle, nil},
		},
		"strings-double": {
			{`%(\(\w+\))?[-#0 +]*([0-9]+|[*])?(\.([0-9]+|[*]))?[hlL]?[E-GXc-giorsux%]`, LiteralStringInterpol, nil},
			{`[^\\\'"%\n]+`, LiteralStringDouble, nil},
			{`[\'"\\]`, LiteralStringDouble, nil},
			{`%`, LiteralStringDouble, nil},
		},
		"dqs": {
			{`"`, LiteralStringDouble, Pop(1)},
			{`\\\\|\\"|\\\n`, LiteralStringEscape, nil},
			Include("strings-double"),
		},
		"sqs": {
			{`'`, LiteralStringSingle, Pop(1)},
			{`\\\\|\\'|\\\n`, LiteralStringEscape, nil},
			Include("strings-single"),
		},
		"tdqs": {
			{`"""`, LiteralStringDouble, Pop(1)},
			Include("strings-double"),
			{`\n`, LiteralStringDouble, nil},
		},
		"tsqs": {
			{`'''`, LiteralStringSingle, Pop(1)},
			Include("strings-single"),
			{`\n`, LiteralStringSingle, nil},
		},
	}
}
