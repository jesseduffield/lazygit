package h

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// HLB lexer.
var HLB = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "HLB",
		Aliases:   []string{"hlb"},
		Filenames: []string{"*.hlb"},
		MimeTypes: []string{},
	},
	hlbRules,
))

func hlbRules() Rules {
	return Rules{
		"root": {
			{`(#.*)`, ByGroups(CommentSingle), nil},
			{`((\b(0(b|B|o|O|x|X)[a-fA-F0-9]+)\b)|(\b(0|[1-9][0-9]*)\b))`, ByGroups(LiteralNumber), nil},
			{`((\b(true|false)\b))`, ByGroups(NameBuiltin), nil},
			{`(\bstring\b|\bint\b|\bbool\b|\bfs\b|\boption\b)`, ByGroups(KeywordType), nil},
			{`(\b[a-zA-Z_][a-zA-Z0-9]*\b)(\()`, ByGroups(NameFunction, Punctuation), Push("params")},
			{`(\{)`, ByGroups(Punctuation), Push("block")},
			{`(\n|\r|\r\n)`, Text, nil},
			{`.`, Text, nil},
		},
		"string": {
			{`"`, LiteralString, Pop(1)},
			{`\\"`, LiteralString, nil},
			{`[^\\"]+`, LiteralString, nil},
		},
		"block": {
			{`(\})`, ByGroups(Punctuation), Pop(1)},
			{`(#.*)`, ByGroups(CommentSingle), nil},
			{`((\b(0(b|B|o|O|x|X)[a-fA-F0-9]+)\b)|(\b(0|[1-9][0-9]*)\b))`, ByGroups(LiteralNumber), nil},
			{`((\b(true|false)\b))`, ByGroups(KeywordConstant), nil},
			{`"`, LiteralString, Push("string")},
			{`(with)`, ByGroups(KeywordReserved), nil},
			{`(as)([\t ]+)(\b[a-zA-Z_][a-zA-Z0-9]*\b)`, ByGroups(KeywordReserved, Text, NameFunction), nil},
			{`(\bstring\b|\bint\b|\bbool\b|\bfs\b|\boption\b)([\t ]+)(\{)`, ByGroups(KeywordType, Text, Punctuation), Push("block")},
			{`(?!\b(?:scratch|image|resolve|http|checksum|chmod|filename|git|keepGitDir|local|includePatterns|excludePatterns|followPaths|generate|frontendInput|shell|run|readonlyRootfs|env|dir|user|network|security|host|ssh|secret|mount|target|localPath|uid|gid|mode|readonly|tmpfs|sourcePath|cache|mkdir|createParents|chown|createdTime|mkfile|rm|allowNotFound|allowWildcards|copy|followSymlinks|contentsOnly|unpack|createDestPath)\b)(\b[a-zA-Z_][a-zA-Z0-9]*\b)`, ByGroups(NameOther), nil},
			{`(\n|\r|\r\n)`, Text, nil},
			{`.`, Text, nil},
		},
		"params": {
			{`(\))`, ByGroups(Punctuation), Pop(1)},
			{`(variadic)`, ByGroups(Keyword), nil},
			{`(\bstring\b|\bint\b|\bbool\b|\bfs\b|\boption\b)`, ByGroups(KeywordType), nil},
			{`(\b[a-zA-Z_][a-zA-Z0-9]*\b)`, ByGroups(NameOther), nil},
			{`(\n|\r|\r\n)`, Text, nil},
			{`.`, Text, nil},
		},
	}
}
