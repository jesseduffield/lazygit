package p

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Python lexer.
var Python = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Python",
		Aliases:   []string{"python", "py", "sage", "python3", "py3"},
		Filenames: []string{"*.py", "*.pyi", "*.pyw", "*.jy", "*.sage", "*.sc", "SConstruct", "SConscript", "*.bzl", "BUCK", "BUILD", "BUILD.bazel", "WORKSPACE", "*.tac"},
		MimeTypes: []string{"text/x-python", "application/x-python", "text/x-python3", "application/x-python3"},
	},
	pythonRules,
))

func pythonRules() Rules {
	const pythonIdentifier = `[_\p{L}][_\p{L}\p{N}]*`

	return Rules{
		"root": {
			{`\n`, Text, nil},
			{`^(\s*)([rRuUbB]{,2})("""(?:.|\n)*?""")`, ByGroups(Text, LiteralStringAffix, LiteralStringDoc), nil},
			{`^(\s*)([rRuUbB]{,2})('''(?:.|\n)*?''')`, ByGroups(Text, LiteralStringAffix, LiteralStringDoc), nil},
			{`\A#!.+$`, CommentHashbang, nil},
			{`#.*$`, CommentSingle, nil},
			{`\\\n`, Text, nil},
			{`\\`, Text, nil},
			Include("keywords"),
			{`(def)((?:\s|\\\s)+)`, ByGroups(Keyword, Text), Push("funcname")},
			{`(class)((?:\s|\\\s)+)`, ByGroups(Keyword, Text), Push("classname")},
			{`(from)((?:\s|\\\s)+)`, ByGroups(KeywordNamespace, Text), Push("fromimport")},
			{`(import)((?:\s|\\\s)+)`, ByGroups(KeywordNamespace, Text), Push("import")},
			Include("expr"),
		},
		"expr": {
			{`(?i)(rf|fr)(""")`, ByGroups(LiteralStringAffix, LiteralStringDouble), Combined("rfstringescape", "tdqf")},
			{`(?i)(rf|fr)(''')`, ByGroups(LiteralStringAffix, LiteralStringSingle), Combined("rfstringescape", "tsqf")},
			{`(?i)(rf|fr)(")`, ByGroups(LiteralStringAffix, LiteralStringDouble), Combined("rfstringescape", "dqf")},
			{`(?i)(rf|fr)(')`, ByGroups(LiteralStringAffix, LiteralStringSingle), Combined("rfstringescape", "sqf")},
			{`([fF])(""")`, ByGroups(LiteralStringAffix, LiteralStringDouble), Combined("fstringescape", "tdqf")},
			{`([fF])(''')`, ByGroups(LiteralStringAffix, LiteralStringSingle), Combined("fstringescape", "tsqf")},
			{`([fF])(")`, ByGroups(LiteralStringAffix, LiteralStringDouble), Combined("fstringescape", "dqf")},
			{`([fF])(')`, ByGroups(LiteralStringAffix, LiteralStringSingle), Combined("fstringescape", "sqf")},
			{`(?i)(rb|br|r)(""")`, ByGroups(LiteralStringAffix, LiteralStringDouble), Push("tdqs")},
			{`(?i)(rb|br|r)(''')`, ByGroups(LiteralStringAffix, LiteralStringSingle), Push("tsqs")},
			{`(?i)(rb|br|r)(")`, ByGroups(LiteralStringAffix, LiteralStringDouble), Push("dqs")},
			{`(?i)(rb|br|r)(')`, ByGroups(LiteralStringAffix, LiteralStringSingle), Push("sqs")},
			{`([uUbB]?)(""")`, ByGroups(LiteralStringAffix, LiteralStringDouble), Combined("stringescape", "tdqs")},
			{`([uUbB]?)(''')`, ByGroups(LiteralStringAffix, LiteralStringSingle), Combined("stringescape", "tsqs")},
			{`([uUbB]?)(")`, ByGroups(LiteralStringAffix, LiteralStringDouble), Combined("stringescape", "dqs")},
			{`([uUbB]?)(')`, ByGroups(LiteralStringAffix, LiteralStringSingle), Combined("stringescape", "sqs")},
			{`[^\S\n]+`, Text, nil},
			Include("numbers"),
			{`!=|==|<<|>>|:=|[-~+/*%=<>&^|.]`, Operator, nil},
			{`[]{}:(),;[]`, Punctuation, nil},
			{`(in|is|and|or|not)\b`, OperatorWord, nil},
			Include("expr-keywords"),
			Include("builtins"),
			Include("magicfuncs"),
			Include("magicvars"),
			Include("name"),
		},
		"expr-inside-fstring": {
			{`[{([]`, Punctuation, Push("expr-inside-fstring-inner")},
			{`(=\s*)?(\![sraf])?\}`, LiteralStringInterpol, Pop(1)},
			{`(=\s*)?(\![sraf])?:`, LiteralStringInterpol, Pop(1)},
			{`\s+`, Text, nil},
			Include("expr"),
		},
		"expr-inside-fstring-inner": {
			{`[{([]`, Punctuation, Push("expr-inside-fstring-inner")},
			{`[])}]`, Punctuation, Pop(1)},
			{`\s+`, Text, nil},
			Include("expr"),
		},
		"expr-keywords": {
			{Words(``, `\b`, `async for`, `await`, `else`, `for`, `if`, `lambda`, `yield`, `yield from`), Keyword, nil},
			{Words(``, `\b`, `True`, `False`, `None`), KeywordConstant, nil},
		},
		"keywords": {
			{Words(``, `\b`, `assert`, `async`, `await`, `break`, `continue`, `del`, `elif`, `else`, `except`, `finally`, `for`, `global`, `if`, `lambda`, `pass`, `raise`, `nonlocal`, `return`, `try`, `while`, `yield`, `yield from`, `as`, `with`), Keyword, nil},
			{Words(``, `\b`, `True`, `False`, `None`), KeywordConstant, nil},
		},
		"builtins": {
			{Words(`(?<!\.)`, `\b`, `__import__`, `abs`, `all`, `any`, `bin`, `bool`, `bytearray`, `bytes`, `chr`, `classmethod`, `compile`, `complex`, `delattr`, `dict`, `dir`, `divmod`, `enumerate`, `eval`, `filter`, `float`, `format`, `frozenset`, `getattr`, `globals`, `hasattr`, `hash`, `hex`, `id`, `input`, `int`, `isinstance`, `issubclass`, `iter`, `len`, `list`, `locals`, `map`, `max`, `memoryview`, `min`, `next`, `object`, `oct`, `open`, `ord`, `pow`, `print`, `property`, `range`, `repr`, `reversed`, `round`, `set`, `setattr`, `slice`, `sorted`, `staticmethod`, `str`, `sum`, `super`, `tuple`, `type`, `vars`, `zip`), NameBuiltin, nil},
			{`(?<!\.)(self|Ellipsis|NotImplemented|cls)\b`, NameBuiltinPseudo, nil},
			{Words(`(?<!\.)`, `\b`, `ArithmeticError`, `AssertionError`, `AttributeError`, `BaseException`, `BufferError`, `BytesWarning`, `DeprecationWarning`, `EOFError`, `EnvironmentError`, `Exception`, `FloatingPointError`, `FutureWarning`, `GeneratorExit`, `IOError`, `ImportError`, `ImportWarning`, `IndentationError`, `IndexError`, `KeyError`, `KeyboardInterrupt`, `LookupError`, `MemoryError`, `NameError`, `NotImplementedError`, `OSError`, `OverflowError`, `PendingDeprecationWarning`, `ReferenceError`, `ResourceWarning`, `RuntimeError`, `RuntimeWarning`, `StopIteration`, `SyntaxError`, `SyntaxWarning`, `SystemError`, `SystemExit`, `TabError`, `TypeError`, `UnboundLocalError`, `UnicodeDecodeError`, `UnicodeEncodeError`, `UnicodeError`, `UnicodeTranslateError`, `UnicodeWarning`, `UserWarning`, `ValueError`, `VMSError`, `Warning`, `WindowsError`, `ZeroDivisionError`, `BlockingIOError`, `ChildProcessError`, `ConnectionError`, `BrokenPipeError`, `ConnectionAbortedError`, `ConnectionRefusedError`, `ConnectionResetError`, `FileExistsError`, `FileNotFoundError`, `InterruptedError`, `IsADirectoryError`, `NotADirectoryError`, `PermissionError`, `ProcessLookupError`, `TimeoutError`, `StopAsyncIteration`, `ModuleNotFoundError`, `RecursionError`), NameException, nil},
		},
		"magicfuncs": {
			{Words(``, `\b`, `__abs__`, `__add__`, `__aenter__`, `__aexit__`, `__aiter__`, `__and__`, `__anext__`, `__await__`, `__bool__`, `__bytes__`, `__call__`, `__complex__`, `__contains__`, `__del__`, `__delattr__`, `__delete__`, `__delitem__`, `__dir__`, `__divmod__`, `__enter__`, `__eq__`, `__exit__`, `__float__`, `__floordiv__`, `__format__`, `__ge__`, `__get__`, `__getattr__`, `__getattribute__`, `__getitem__`, `__gt__`, `__hash__`, `__iadd__`, `__iand__`, `__ifloordiv__`, `__ilshift__`, `__imatmul__`, `__imod__`, `__imul__`, `__index__`, `__init__`, `__instancecheck__`, `__int__`, `__invert__`, `__ior__`, `__ipow__`, `__irshift__`, `__isub__`, `__iter__`, `__itruediv__`, `__ixor__`, `__le__`, `__len__`, `__length_hint__`, `__lshift__`, `__lt__`, `__matmul__`, `__missing__`, `__mod__`, `__mul__`, `__ne__`, `__neg__`, `__new__`, `__next__`, `__or__`, `__pos__`, `__pow__`, `__prepare__`, `__radd__`, `__rand__`, `__rdivmod__`, `__repr__`, `__reversed__`, `__rfloordiv__`, `__rlshift__`, `__rmatmul__`, `__rmod__`, `__rmul__`, `__ror__`, `__round__`, `__rpow__`, `__rrshift__`, `__rshift__`, `__rsub__`, `__rtruediv__`, `__rxor__`, `__set__`, `__setattr__`, `__setitem__`, `__str__`, `__sub__`, `__subclasscheck__`, `__truediv__`, `__xor__`), NameFunctionMagic, nil},
		},
		"magicvars": {
			{Words(``, `\b`, `__annotations__`, `__bases__`, `__class__`, `__closure__`, `__code__`, `__defaults__`, `__dict__`, `__doc__`, `__file__`, `__func__`, `__globals__`, `__kwdefaults__`, `__module__`, `__mro__`, `__name__`, `__objclass__`, `__qualname__`, `__self__`, `__slots__`, `__weakref__`), NameVariableMagic, nil},
		},
		"numbers": {
			{`(\d(?:_?\d)*\.(?:\d(?:_?\d)*)?|(?:\d(?:_?\d)*)?\.\d(?:_?\d)*)([eE][+-]?\d(?:_?\d)*)?`, LiteralNumberFloat, nil},
			{`\d(?:_?\d)*[eE][+-]?\d(?:_?\d)*j?`, LiteralNumberFloat, nil},
			{`0[oO](?:_?[0-7])+`, LiteralNumberOct, nil},
			{`0[bB](?:_?[01])+`, LiteralNumberBin, nil},
			{`0[xX](?:_?[a-fA-F0-9])+`, LiteralNumberHex, nil},
			{`\d(?:_?\d)*`, LiteralNumberInteger, nil},
		},
		"name": {
			{`@` + pythonIdentifier, NameDecorator, nil},
			{`@`, Operator, nil},
			{pythonIdentifier, Name, nil},
		},
		"funcname": {
			Include("magicfuncs"),
			{pythonIdentifier, NameFunction, Pop(1)},
			Default(Pop(1)),
		},
		"classname": {
			{pythonIdentifier, NameClass, Pop(1)},
		},
		"import": {
			{`(\s+)(as)(\s+)`, ByGroups(Text, Keyword, Text), nil},
			{`\.`, NameNamespace, nil},
			{pythonIdentifier, NameNamespace, nil},
			{`(\s*)(,)(\s*)`, ByGroups(Text, Operator, Text), nil},
			Default(Pop(1)),
		},
		"fromimport": {
			{`(\s+)(import)\b`, ByGroups(Text, KeywordNamespace), Pop(1)},
			{`\.`, NameNamespace, nil},
			{`None\b`, NameBuiltinPseudo, Pop(1)},
			{pythonIdentifier, NameNamespace, nil},
			Default(Pop(1)),
		},
		"rfstringescape": {
			{`\{\{`, LiteralStringEscape, nil},
			{`\}\}`, LiteralStringEscape, nil},
		},
		"fstringescape": {
			Include("rfstringescape"),
			Include("stringescape"),
		},
		"stringescape": {
			{`\\([\\abfnrtv"\']|\n|N\{.*?\}|u[a-fA-F0-9]{4}|U[a-fA-F0-9]{8}|x[a-fA-F0-9]{2}|[0-7]{1,3})`, LiteralStringEscape, nil},
		},
		"fstrings-single": {
			{`\}`, LiteralStringInterpol, nil},
			{`\{`, LiteralStringInterpol, Push("expr-inside-fstring")},
			{`[^\\\'"{}\n]+`, LiteralStringSingle, nil},
			{`[\'"\\]`, LiteralStringSingle, nil},
		},
		"fstrings-double": {
			{`\}`, LiteralStringInterpol, nil},
			{`\{`, LiteralStringInterpol, Push("expr-inside-fstring")},
			{`[^\\\'"{}\n]+`, LiteralStringDouble, nil},
			{`[\'"\\]`, LiteralStringDouble, nil},
		},
		"strings-single": {
			{`%(\(\w+\))?[-#0 +]*([0-9]+|[*])?(\.([0-9]+|[*]))?[hlL]?[E-GXc-giorsaux%]`, LiteralStringInterpol, nil},
			{`\{((\w+)((\.\w+)|(\[[^\]]+\]))*)?(\![sra])?(\:(.?[<>=\^])?[-+ ]?#?0?(\d+)?,?(\.\d+)?[E-GXb-gnosx%]?)?\}`, LiteralStringInterpol, nil},
			{`[^\\\'"%{\n]+`, LiteralStringSingle, nil},
			{`[\'"\\]`, LiteralStringSingle, nil},
			{`%|(\{{1,2})`, LiteralStringSingle, nil},
		},
		"strings-double": {
			{`%(\(\w+\))?[-#0 +]*([0-9]+|[*])?(\.([0-9]+|[*]))?[hlL]?[E-GXc-giorsaux%]`, LiteralStringInterpol, nil},
			{`\{((\w+)((\.\w+)|(\[[^\]]+\]))*)?(\![sra])?(\:(.?[<>=\^])?[-+ ]?#?0?(\d+)?,?(\.\d+)?[E-GXb-gnosx%]?)?\}`, LiteralStringInterpol, nil},
			{`[^\\\'"%{\n]+`, LiteralStringDouble, nil},
			{`[\'"\\]`, LiteralStringDouble, nil},
			{`%|(\{{1,2})`, LiteralStringDouble, nil},
		},
		"dqf": {
			{`"`, LiteralStringDouble, Pop(1)},
			{`\\\\|\\"|\\\n`, LiteralStringEscape, nil},
			Include("fstrings-double"),
		},
		"sqf": {
			{`'`, LiteralStringSingle, Pop(1)},
			{`\\\\|\\'|\\\n`, LiteralStringEscape, nil},
			Include("fstrings-single"),
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
		"tdqf": {
			{`"""`, LiteralStringDouble, Pop(1)},
			Include("fstrings-double"),
			{`\n`, LiteralStringDouble, nil},
		},
		"tsqf": {
			{`'''`, LiteralStringSingle, Pop(1)},
			Include("fstrings-single"),
			{`\n`, LiteralStringSingle, nil},
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
