package v

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

const vbName = `[_\w][\w]*`

// VB.Net lexer.
var VBNet = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "VB.net",
		Aliases:         []string{"vb.net", "vbnet"},
		Filenames:       []string{"*.vb", "*.bas"},
		MimeTypes:       []string{"text/x-vbnet", "text/x-vba"},
		CaseInsensitive: true,
	},
	vbNetRules,
))

func vbNetRules() Rules {
	return Rules{
		"root": {
			{`^\s*<.*?>`, NameAttribute, nil},
			{`\s+`, Text, nil},
			{`\n`, Text, nil},
			{`rem\b.*?\n`, Comment, nil},
			{`'.*?\n`, Comment, nil},
			{`#If\s.*?\sThen|#ElseIf\s.*?\sThen|#Else|#End\s+If|#Const|#ExternalSource.*?\n|#End\s+ExternalSource|#Region.*?\n|#End\s+Region|#ExternalChecksum`, CommentPreproc, nil},
			{`[(){}!#,.:]`, Punctuation, nil},
			{`Option\s+(Strict|Explicit|Compare)\s+(On|Off|Binary|Text)`, KeywordDeclaration, nil},
			{Words(`(?<!\.)`, `\b`, `AddHandler`, `Alias`, `ByRef`, `ByVal`, `Call`, `Case`, `Catch`, `CBool`, `CByte`, `CChar`, `CDate`, `CDec`, `CDbl`, `CInt`, `CLng`, `CObj`, `Continue`, `CSByte`, `CShort`, `CSng`, `CStr`, `CType`, `CUInt`, `CULng`, `CUShort`, `Declare`, `Default`, `Delegate`, `DirectCast`, `Do`, `Each`, `Else`, `ElseIf`, `EndIf`, `Erase`, `Error`, `Event`, `Exit`, `False`, `Finally`, `For`, `Friend`, `Get`, `Global`, `GoSub`, `GoTo`, `Handles`, `If`, `Implements`, `Inherits`, `Interface`, `Let`, `Lib`, `Loop`, `Me`, `MustInherit`, `MustOverride`, `MyBase`, `MyClass`, `Narrowing`, `New`, `Next`, `Not`, `Nothing`, `NotInheritable`, `NotOverridable`, `Of`, `On`, `Operator`, `Option`, `Optional`, `Overloads`, `Overridable`, `Overrides`, `ParamArray`, `Partial`, `Private`, `Protected`, `Public`, `RaiseEvent`, `ReadOnly`, `ReDim`, `RemoveHandler`, `Resume`, `Return`, `Select`, `Set`, `Shadows`, `Shared`, `Single`, `Static`, `Step`, `Stop`, `SyncLock`, `Then`, `Throw`, `To`, `True`, `Try`, `TryCast`, `Wend`, `Using`, `When`, `While`, `Widening`, `With`, `WithEvents`, `WriteOnly`), Keyword, nil},
			{`(?<!\.)End\b`, Keyword, Push("end")},
			{`(?<!\.)(Dim|Const)\b`, Keyword, Push("dim")},
			{`(?<!\.)(Function|Sub|Property)(\s+)`, ByGroups(Keyword, Text), Push("funcname")},
			{`(?<!\.)(Class|Structure|Enum)(\s+)`, ByGroups(Keyword, Text), Push("classname")},
			{`(?<!\.)(Module|Namespace|Imports)(\s+)`, ByGroups(Keyword, Text), Push("namespace")},
			{`(?<!\.)(Boolean|Byte|Char|Date|Decimal|Double|Integer|Long|Object|SByte|Short|Single|String|Variant|UInteger|ULong|UShort)\b`, KeywordType, nil},
			{`(?<!\.)(AddressOf|And|AndAlso|As|GetType|In|Is|IsNot|Like|Mod|Or|OrElse|TypeOf|Xor)\b`, OperatorWord, nil},
			{`&=|[*]=|/=|\\=|\^=|\+=|-=|<<=|>>=|<<|>>|:=|<=|>=|<>|[-&*/\\^+=<>\[\]]`, Operator, nil},
			{`"`, LiteralString, Push("string")},
			{`_\n`, Text, nil},
			{vbName, Name, nil},
			{`#.*?#`, LiteralDate, nil},
			{`(\d+\.\d*|\d*\.\d+)(F[+-]?[0-9]+)?`, LiteralNumberFloat, nil},
			{`\d+([SILDFR]|US|UI|UL)?`, LiteralNumberInteger, nil},
			{`&H[0-9a-f]+([SILDFR]|US|UI|UL)?`, LiteralNumberInteger, nil},
			{`&O[0-7]+([SILDFR]|US|UI|UL)?`, LiteralNumberInteger, nil},
		},
		"string": {
			{`""`, LiteralString, nil},
			{`"C?`, LiteralString, Pop(1)},
			{`[^"]+`, LiteralString, nil},
		},
		"dim": {
			{vbName, NameVariable, Pop(1)},
			Default(Pop(1)),
		},
		"funcname": {
			{vbName, NameFunction, Pop(1)},
		},
		"classname": {
			{vbName, NameClass, Pop(1)},
		},
		"namespace": {
			{vbName, NameNamespace, nil},
			{`\.`, NameNamespace, nil},
			Default(Pop(1)),
		},
		"end": {
			{`\s+`, Text, nil},
			{`(Function|Sub|Property|Class|Structure|Enum|Module|Namespace)\b`, Keyword, Pop(1)},
			Default(Pop(1)),
		},
	}
}
