package a

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Al lexer.
var Al = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "AL",
		Aliases:         []string{"al"},
		Filenames:       []string{"*.al", "*.dal"},
		MimeTypes:       []string{"text/x-al"},
		DotAll:          true,
		CaseInsensitive: true,
	},
	alRules,
))

// https://github.com/microsoft/AL/blob/master/grammar/alsyntax.tmlanguage
func alRules() Rules {
	return Rules{
		"root": {
			{`\s+`, TextWhitespace, nil},
			{`(?s)\/\*.*?\\*\*\/`, CommentMultiline, nil},
			{`(?s)//.*?\n`, CommentSingle, nil},
			{`\"([^\"])*\"`, Text, nil},
			{`'([^'])*'`, LiteralString, nil},
			{`\b(?i:(ARRAY|ASSERTERROR|BEGIN|BREAK|CASE|DO|DOWNTO|ELSE|END|EVENT|EXIT|FOR|FOREACH|FUNCTION|IF|IMPLEMENTS|IN|INDATASET|INTERFACE|INTERNAL|LOCAL|OF|PROCEDURE|PROGRAM|PROTECTED|REPEAT|RUNONCLIENT|SECURITYFILTERING|SUPPRESSDISPOSE|TEMPORARY|THEN|TO|TRIGGER|UNTIL|VAR|WHILE|WITH|WITHEVENTS))\b`, Keyword, nil},
			{`\b(?i:(AND|DIV|MOD|NOT|OR|XOR))\b`, OperatorWord, nil},
			{`\b(?i:(AVERAGE|CONST|COUNT|EXIST|FIELD|FILTER|LOOKUP|MAX|MIN|ORDER|SORTING|SUM|TABLEDATA|UPPERLIMIT|WHERE|ASCENDING|DESCENDING))\b`, Keyword, nil},
			{`\b(?i:(CODEUNIT|PAGE|PAGEEXTENSION|PAGECUSTOMIZATION|DOTNET|ENUM|ENUMEXTENSION|VALUE|QUERY|REPORT|TABLE|TABLEEXTENSION|XMLPORT|PROFILE|CONTROLADDIN|REPORTEXTENSION|INTERFACE|PERMISSIONSET|PERMISSIONSETEXTENSION|ENTITLEMENT))\b`, Keyword, nil},
			{`\b(?i:(Action|Array|Automation|BigInteger|BigText|Blob|Boolean|Byte|Char|ClientType|Code|Codeunit|CompletionTriggerErrorLevel|ConnectionType|Database|DataClassification|DataScope|Date|DateFormula|DateTime|Decimal|DefaultLayout|Dialog|Dictionary|DotNet|DotNetAssembly|DotNetTypeDeclaration|Duration|Enum|ErrorInfo|ErrorType|ExecutionContext|ExecutionMode|FieldClass|FieldRef|FieldType|File|FilterPageBuilder|Guid|InStream|Integer|Joker|KeyRef|List|ModuleDependencyInfo|ModuleInfo|None|Notification|NotificationScope|ObjectType|Option|OutStream|Page|PageResult|Query|Record|RecordId|RecordRef|Report|ReportFormat|SecurityFilter|SecurityFiltering|Table|TableConnectionType|TableFilter|TestAction|TestField|TestFilterField|TestPage|TestPermissions|TestRequestPage|Text|TextBuilder|TextConst|TextEncoding|Time|TransactionModel|TransactionType|Variant|Verbosity|Version|XmlPort|HttpContent|HttpHeaders|HttpClient|HttpRequestMessage|HttpResponseMessage|JsonToken|JsonValue|JsonArray|JsonObject|View|Views|XmlAttribute|XmlAttributeCollection|XmlComment|XmlCData|XmlDeclaration|XmlDocument|XmlDocumentType|XmlElement|XmlNamespaceManager|XmlNameTable|XmlNode|XmlNodeList|XmlProcessingInstruction|XmlReadOptions|XmlText|XmlWriteOptions|WebServiceActionContext|WebServiceActionResultCode|SessionSettings))\b`, Keyword, nil},
			{`\b([<>]=|<>|<|>)\b?`, Operator, nil},
			{`\b(\-|\+|\/|\*)\b`, Operator, nil},
			{`\s*(\:=|\+=|-=|\/=|\*=)\s*?`, Operator, nil},
			{`\b(?i:(ADD|ADDFIRST|ADDLAST|ADDAFTER|ADDBEFORE|ACTION|ACTIONS|AREA|ASSEMBLY|CHARTPART|CUEGROUP|CUSTOMIZES|COLUMN|DATAITEM|DATASET|ELEMENTS|EXTENDS|FIELD|FIELDGROUP|FIELDATTRIBUTE|FIELDELEMENT|FIELDGROUPS|FIELDS|FILTER|FIXED|GRID|GROUP|MOVEAFTER|MOVEBEFORE|KEY|KEYS|LABEL|LABELS|LAYOUT|MODIFY|MOVEFIRST|MOVELAST|MOVEBEFORE|MOVEAFTER|PART|REPEATER|USERCONTROL|REQUESTPAGE|SCHEMA|SEPARATOR|SYSTEMPART|TABLEELEMENT|TEXTATTRIBUTE|TEXTELEMENT|TYPE))\b`, Keyword, nil},
			{`\s*[(\.\.)&\|]\s*`, Operator, nil},
			{`\b((0(x|X)[0-9a-fA-F]*)|(([0-9]+\.?[0-9]*)|(\.[0-9]+))((e|E)(\+|-)?[0-9]+)?)(L|l|UL|ul|u|U|F|f|ll|LL|ull|ULL)?\b`, LiteralNumber, nil},
			{`[;:,]`, Punctuation, nil},
			{`#[ \t]*(if|else|elif|endif|define|undef|region|endregion|pragma)\b.*?\n`, CommentPreproc, nil},
			{`\w+`, Text, nil},
			{`.`, Text, nil},
		},
	}
}
