package q

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Qbasic lexer.
var Qbasic = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "QBasic",
		Aliases:   []string{"qbasic", "basic"},
		Filenames: []string{"*.BAS", "*.bas"},
		MimeTypes: []string{"text/basic"},
	},
	qbasicRules,
))

func qbasicRules() Rules {
	return Rules{
		"root": {
			{`\n+`, Text, nil},
			{`\s+`, TextWhitespace, nil},
			{`^(\s*)(\d*)(\s*)(REM .*)$`, ByGroups(TextWhitespace, NameLabel, TextWhitespace, CommentSingle), nil},
			{`^(\s*)(\d+)(\s*)`, ByGroups(TextWhitespace, NameLabel, TextWhitespace), nil},
			{`(?=[\s]*)(\w+)(?=[\s]*=)`, NameVariableGlobal, nil},
			{`(?=[^"]*)\'.*$`, CommentSingle, nil},
			{`"[^\n"]*"`, LiteralStringDouble, nil},
			{`(END)(\s+)(FUNCTION|IF|SELECT|SUB)`, ByGroups(KeywordReserved, TextWhitespace, KeywordReserved), nil},
			{`(DECLARE)(\s+)([A-Z]+)(\s+)(\S+)`, ByGroups(KeywordDeclaration, TextWhitespace, NameVariable, TextWhitespace, Name), nil},
			{`(DIM)(\s+)(SHARED)(\s+)([^\s(]+)`, ByGroups(KeywordDeclaration, TextWhitespace, NameVariable, TextWhitespace, NameVariableGlobal), nil},
			{`(DIM)(\s+)([^\s(]+)`, ByGroups(KeywordDeclaration, TextWhitespace, NameVariableGlobal), nil},
			{`^(\s*)([a-zA-Z_]+)(\s*)(\=)`, ByGroups(TextWhitespace, NameVariableGlobal, TextWhitespace, Operator), nil},
			{`(GOTO|GOSUB)(\s+)(\w+\:?)`, ByGroups(KeywordReserved, TextWhitespace, NameLabel), nil},
			{`(SUB)(\s+)(\w+\:?)`, ByGroups(KeywordReserved, TextWhitespace, NameLabel), nil},
			Include("declarations"),
			Include("functions"),
			Include("metacommands"),
			Include("operators"),
			Include("statements"),
			Include("keywords"),
			{`[a-zA-Z_]\w*[$@#&!]`, NameVariableGlobal, nil},
			{`[a-zA-Z_]\w*\:`, NameLabel, nil},
			{`\-?\d*\.\d+[@|#]?`, LiteralNumberFloat, nil},
			{`\-?\d+[@|#]`, LiteralNumberFloat, nil},
			{`\-?\d+#?`, LiteralNumberIntegerLong, nil},
			{`\-?\d+#?`, LiteralNumberInteger, nil},
			{`!=|==|:=|\.=|<<|>>|[-~+/\\*%=<>&^|?:!.]`, Operator, nil},
			{`[\[\]{}(),;]`, Punctuation, nil},
			{`[\w]+`, NameVariableGlobal, nil},
		},
		"declarations": {
			{`\b(DATA|LET)(?=\(|\b)`, KeywordDeclaration, nil},
		},
		"functions": {
			{`\b(ABS|ASC|ATN|CDBL|CHR\$|CINT|CLNG|COMMAND\$|COS|CSNG|CSRLIN|CVD|CVDMBF|CVI|CVL|CVS|CVSMBF|DATE\$|ENVIRON\$|EOF|ERDEV|ERDEV\$|ERL|ERR|EXP|FILEATTR|FIX|FRE|FREEFILE|HEX\$|INKEY\$|INP|INPUT\$|INSTR|INT|IOCTL\$|LBOUND|LCASE\$|LEFT\$|LEN|LOC|LOF|LOG|LPOS|LTRIM\$|MID\$|MKD\$|MKDMBF\$|MKI\$|MKL\$|MKS\$|MKSMBF\$|OCT\$|PEEK|PEN|PLAY|PMAP|POINT|POS|RIGHT\$|RND|RTRIM\$|SADD|SCREEN|SEEK|SETMEM|SGN|SIN|SPACE\$|SPC|SQR|STICK|STR\$|STRIG|STRING\$|TAB|TAN|TIME\$|TIMER|UBOUND|UCASE\$|VAL|VARPTR|VARPTR\$|VARSEG)(?=\(|\b)`, KeywordReserved, nil},
		},
		"metacommands": {
			{`\b(\$DYNAMIC|\$INCLUDE|\$STATIC)(?=\(|\b)`, KeywordConstant, nil},
		},
		"operators": {
			{`\b(AND|EQV|IMP|NOT|OR|XOR)(?=\(|\b)`, OperatorWord, nil},
		},
		"statements": {
			{`\b(BEEP|BLOAD|BSAVE|CALL|CALL\ ABSOLUTE|CALL\ INTERRUPT|CALLS|CHAIN|CHDIR|CIRCLE|CLEAR|CLOSE|CLS|COLOR|COM|COMMON|CONST|DATA|DATE\$|DECLARE|DEF\ FN|DEF\ SEG|DEFDBL|DEFINT|DEFLNG|DEFSNG|DEFSTR|DEF|DIM|DO|LOOP|DRAW|END|ENVIRON|ERASE|ERROR|EXIT|FIELD|FILES|FOR|NEXT|FUNCTION|GET|GOSUB|GOTO|IF|THEN|INPUT|INPUT\ \#|IOCTL|KEY|KEY|KILL|LET|LINE|LINE\ INPUT|LINE\ INPUT\ \#|LOCATE|LOCK|UNLOCK|LPRINT|LSET|MID\$|MKDIR|NAME|ON\ COM|ON\ ERROR|ON\ KEY|ON\ PEN|ON\ PLAY|ON\ STRIG|ON\ TIMER|ON\ UEVENT|ON|OPEN|OPEN\ COM|OPTION\ BASE|OUT|PAINT|PALETTE|PCOPY|PEN|PLAY|POKE|PRESET|PRINT|PRINT\ \#|PRINT\ USING|PSET|PUT|PUT|RANDOMIZE|READ|REDIM|REM|RESET|RESTORE|RESUME|RETURN|RMDIR|RSET|RUN|SCREEN|SEEK|SELECT\ CASE|SHARED|SHELL|SLEEP|SOUND|STATIC|STOP|STRIG|SUB|SWAP|SYSTEM|TIME\$|TIMER|TROFF|TRON|TYPE|UEVENT|UNLOCK|VIEW|WAIT|WHILE|WEND|WIDTH|WINDOW|WRITE)\b`, KeywordReserved, nil},
		},
		"keywords": {
			{`\b(ACCESS|ALIAS|ANY|APPEND|AS|BASE|BINARY|BYVAL|CASE|CDECL|DOUBLE|ELSE|ELSEIF|ENDIF|INTEGER|IS|LIST|LOCAL|LONG|LOOP|MOD|NEXT|OFF|ON|OUTPUT|RANDOM|SIGNAL|SINGLE|STEP|STRING|THEN|TO|UNTIL|USING|WEND)\b`, Keyword, nil},
		},
	}
}
