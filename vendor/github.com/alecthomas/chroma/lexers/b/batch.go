package b

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Batchfile lexer.
var Batchfile = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "Batchfile",
		Aliases:         []string{"bat", "batch", "dosbatch", "winbatch"},
		Filenames:       []string{"*.bat", "*.cmd"},
		MimeTypes:       []string{"application/x-dos-batch"},
		CaseInsensitive: true,
	},
	batchfileRules,
))

func batchfileRules() Rules {
	return Rules{
		"root": {
			{`\)((?=\()|(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a]))(?:(?:[^\n\x1a^]|\^[\n\x1a]?[\w\W])*)`, CommentSingle, nil},
			{`(?=((?:(?<=^[^:])|^[^:]?)[\t\v\f\r ,;=\xa0]*)(:))`, Text, Push("follow")},
			{`(?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)`, UsingSelf("text"), nil},
			Include("redirect"),
			{`[\n\x1a]+`, Text, nil},
			{`\(`, Punctuation, Push("root/compound")},
			{`@+`, Punctuation, nil},
			{`((?:for|if|rem)(?:(?=(?:\^[\n\x1a]?)?/)|(?:(?!\^)|(?<=m))(?:(?=\()|(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a]))))((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)?(?:(?:(?:\^[\n\x1a]?)?[^"\n\x1a&<>|\t\v\f\r ,;=\xa0])+)?(?:\^[\n\x1a]?)?/(?:\^[\n\x1a]?)?\?)`, ByGroups(Keyword, UsingSelf("text")), Push("follow")},
			{`(goto(?=(?:\^[\n\x1a]?)?[\t\v\f\r ,;=\xa0+./:[\\\]]|[\n\x1a&<>|(]))((?:(?:"[^\n\x1a"]*(?:"|(?=[\n\x1a])))|(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))|[^"%\n\x1a&<>|])*(?:\^[\n\x1a]?)?/(?:\^[\n\x1a]?)?\?(?:(?:"[^\n\x1a"]*(?:"|(?=[\n\x1a])))|(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))|[^"%\n\x1a&<>|])*)`, ByGroups(Keyword, UsingSelf("text")), Push("follow")},
			{Words(``, `(?=(?:\^[\n\x1a]?)?[\t\v\f\r ,;=\xa0+./:[\\\]]|[\n\x1a&<>|(])`, `assoc`, `break`, `cd`, `chdir`, `cls`, `color`, `copy`, `date`, `del`, `dir`, `dpath`, `echo`, `endlocal`, `erase`, `exit`, `ftype`, `keys`, `md`, `mkdir`, `mklink`, `move`, `path`, `pause`, `popd`, `prompt`, `pushd`, `rd`, `ren`, `rename`, `rmdir`, `setlocal`, `shift`, `start`, `time`, `title`, `type`, `ver`, `verify`, `vol`), Keyword, Push("follow")},
			{`(call)((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)?)(:)`, ByGroups(Keyword, UsingSelf("text"), Punctuation), Push("call")},
			{`call(?=(?:\^[\n\x1a]?)?[\t\v\f\r ,;=\xa0+./:[\\\]]|[\n\x1a&<>|(])`, Keyword, nil},
			{`(for(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a])(?!\^))((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+))(/f(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a]))`, ByGroups(Keyword, UsingSelf("text"), Keyword), Push("for/f", "for")},
			{`(for(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a])(?!\^))((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+))(/l(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a]))`, ByGroups(Keyword, UsingSelf("text"), Keyword), Push("for/l", "for")},
			{`for(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a])(?!\^)`, Keyword, Push("for2", "for")},
			{`(goto(?=(?:\^[\n\x1a]?)?[\t\v\f\r ,;=\xa0+./:[\\\]]|[\n\x1a&<>|(]))((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)?)(:?)`, ByGroups(Keyword, UsingSelf("text"), Punctuation), Push("label")},
			{`(if(?:(?=\()|(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a]))(?!\^))((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)?)((?:/i(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a]))?)((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)?)((?:not(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a]))?)((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)?)`, ByGroups(Keyword, UsingSelf("text"), Keyword, UsingSelf("text"), Keyword, UsingSelf("text")), Push("(?", "if")},
			{`rem(((?=\()|(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a]))(?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)?(?:[&<>|]+|(?:(?:"[^\n\x1a"]*(?:"|(?=[\n\x1a])))|(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))|(?:(?:(?:\^[\n\x1a]?)?[^"\n\x1a&<>|\t\v\f\r ,;=\xa0])+))+)?.*|(?=(?:\^[\n\x1a]?)?[\t\v\f\r ,;=\xa0+./:[\\\]]|[\n\x1a&<>|(])(?:(?:[^\n\x1a^]|\^[\n\x1a]?[\w\W])*))`, CommentSingle, Push("follow")},
			{`(set(?=(?:\^[\n\x1a]?)?[\t\v\f\r ,;=\xa0+./:[\\\]]|[\n\x1a&<>|(]))((?:(?:\^[\n\x1a]?)?[^\S\n])*)(/a)`, ByGroups(Keyword, UsingSelf("text"), Keyword), Push("arithmetic")},
			{`(set(?=(?:\^[\n\x1a]?)?[\t\v\f\r ,;=\xa0+./:[\\\]]|[\n\x1a&<>|(]))((?:(?:\^[\n\x1a]?)?[^\S\n])*)((?:/p)?)((?:(?:\^[\n\x1a]?)?[^\S\n])*)((?:(?:(?:\^[\n\x1a]?)?[^"\n\x1a&<>|^=]|\^[\n\x1a]?[^"=])+)?)((?:(?:\^[\n\x1a]?)?=)?)`, ByGroups(Keyword, UsingSelf("text"), Keyword, UsingSelf("text"), UsingSelf("variable"), Punctuation), Push("follow")},
			Default(Push("follow")),
		},
		"follow": {
			{`((?:(?<=^[^:])|^[^:]?)[\t\v\f\r ,;=\xa0]*)(:)([\t\v\f\r ,;=\xa0]*)((?:(?:[^\n\x1a&<>|\t\v\f\r ,;=\xa0+:^]|\^[\n\x1a]?[\w\W])*))(.*)`, ByGroups(Text, Punctuation, Text, NameLabel, CommentSingle), nil},
			Include("redirect"),
			{`(?=[\n\x1a])`, Text, Pop(1)},
			{`\|\|?|&&?`, Punctuation, Pop(1)},
			Include("text"),
		},
		"arithmetic": {
			{`0[0-7]+`, LiteralNumberOct, nil},
			{`0x[\da-f]+`, LiteralNumberHex, nil},
			{`\d+`, LiteralNumberInteger, nil},
			{`[(),]+`, Punctuation, nil},
			{`([=+\-*/!~]|%|\^\^)+`, Operator, nil},
			{`((?:"[^\n\x1a"]*(?:"|(?=[\n\x1a])))|(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))|(\^[\n\x1a]?)?[^()=+\-*/!~%^"\n\x1a&<>|\t\v\f\r ,;=\xa0]|\^[\n\x1a\t\v\f\r ,;=\xa0]?[\w\W])+`, UsingSelf("variable"), nil},
			{`(?=[\x00|&])`, Text, Pop(1)},
			Include("follow"),
		},
		"call": {
			{`(:?)((?:(?:[^\n\x1a&<>|\t\v\f\r ,;=\xa0+:^]|\^[\n\x1a]?[\w\W])*))`, ByGroups(Punctuation, NameLabel), Pop(1)},
		},
		"label": {
			{`((?:(?:[^\n\x1a&<>|\t\v\f\r ,;=\xa0+:^]|\^[\n\x1a]?[\w\W])*)?)((?:(?:"[^\n\x1a"]*(?:"|(?=[\n\x1a])))|(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))|\^[\n\x1a]?[\w\W]|[^"%^\n\x1a&<>|])*)`, ByGroups(NameLabel, CommentSingle), Pop(1)},
		},
		"redirect": {
			{`((?:(?<=[\n\x1a\t\v\f\r ,;=\xa0])\d)?)(>>?&|<&)([\n\x1a\t\v\f\r ,;=\xa0]*)(\d)`, ByGroups(LiteralNumberInteger, Punctuation, Text, LiteralNumberInteger), nil},
			{`((?:(?<=[\n\x1a\t\v\f\r ,;=\xa0])(?<!\^[\n\x1a])\d)?)(>>?|<)((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)?(?:[&<>|]+|(?:(?:"[^\n\x1a"]*(?:"|(?=[\n\x1a])))|(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))|(?:(?:(?:\^[\n\x1a]?)?[^"\n\x1a&<>|\t\v\f\r ,;=\xa0])+))+))`, ByGroups(LiteralNumberInteger, Punctuation, UsingSelf("text")), nil},
		},
		"root/compound": {
			{`\)`, Punctuation, Pop(1)},
			{`(?=((?:(?<=^[^:])|^[^:]?)[\t\v\f\r ,;=\xa0]*)(:))`, Text, Push("follow/compound")},
			{`(?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)`, UsingSelf("text"), nil},
			Include("redirect/compound"),
			{`[\n\x1a]+`, Text, nil},
			{`\(`, Punctuation, Push("root/compound")},
			{`@+`, Punctuation, nil},
			{`((?:for|if|rem)(?:(?=(?:\^[\n\x1a]?)?/)|(?:(?!\^)|(?<=m))(?:(?=\()|(?:(?=\))|(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a])))))((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)?(?:(?:(?:\^[\n\x1a]?)?[^"\n\x1a&<>|\t\v\f\r ,;=\xa0)])+)?(?:\^[\n\x1a]?)?/(?:\^[\n\x1a]?)?\?)`, ByGroups(Keyword, UsingSelf("text")), Push("follow/compound")},
			{`(goto(?:(?=\))|(?=(?:\^[\n\x1a]?)?[\t\v\f\r ,;=\xa0+./:[\\\]]|[\n\x1a&<>|(])))((?:(?:"[^\n\x1a"]*(?:"|(?=[\n\x1a])))|(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))|[^"%\n\x1a&<>|)])*(?:\^[\n\x1a]?)?/(?:\^[\n\x1a]?)?\?(?:(?:"[^\n\x1a"]*(?:"|(?=[\n\x1a])))|(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))|[^"%\n\x1a&<>|)])*)`, ByGroups(Keyword, UsingSelf("text")), Push("follow/compound")},
			{Words(``, `(?:(?=\))|(?=(?:\^[\n\x1a]?)?[\t\v\f\r ,;=\xa0+./:[\\\]]|[\n\x1a&<>|(]))`, `assoc`, `break`, `cd`, `chdir`, `cls`, `color`, `copy`, `date`, `del`, `dir`, `dpath`, `echo`, `endlocal`, `erase`, `exit`, `ftype`, `keys`, `md`, `mkdir`, `mklink`, `move`, `path`, `pause`, `popd`, `prompt`, `pushd`, `rd`, `ren`, `rename`, `rmdir`, `setlocal`, `shift`, `start`, `time`, `title`, `type`, `ver`, `verify`, `vol`), Keyword, Push("follow/compound")},
			{`(call)((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)?)(:)`, ByGroups(Keyword, UsingSelf("text"), Punctuation), Push("call/compound")},
			{`call(?:(?=\))|(?=(?:\^[\n\x1a]?)?[\t\v\f\r ,;=\xa0+./:[\\\]]|[\n\x1a&<>|(]))`, Keyword, nil},
			{`(for(?:(?=\))|(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a]))(?!\^))((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+))(/f(?:(?=\))|(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a])))`, ByGroups(Keyword, UsingSelf("text"), Keyword), Push("for/f", "for")},
			{`(for(?:(?=\))|(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a]))(?!\^))((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+))(/l(?:(?=\))|(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a])))`, ByGroups(Keyword, UsingSelf("text"), Keyword), Push("for/l", "for")},
			{`for(?:(?=\))|(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a]))(?!\^)`, Keyword, Push("for2", "for")},
			{`(goto(?:(?=\))|(?=(?:\^[\n\x1a]?)?[\t\v\f\r ,;=\xa0+./:[\\\]]|[\n\x1a&<>|(])))((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)?)(:?)`, ByGroups(Keyword, UsingSelf("text"), Punctuation), Push("label/compound")},
			{`(if(?:(?=\()|(?:(?=\))|(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a])))(?!\^))((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)?)((?:/i(?:(?=\))|(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a])))?)((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)?)((?:not(?:(?=\))|(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a])))?)((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)?)`, ByGroups(Keyword, UsingSelf("text"), Keyword, UsingSelf("text"), Keyword, UsingSelf("text")), Push("(?", "if")},
			{`rem(((?=\()|(?:(?=\))|(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a])))(?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)?(?:[&<>|]+|(?:(?:"[^\n\x1a"]*(?:"|(?=[\n\x1a])))|(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))|(?:(?:(?:\^[\n\x1a]?)?[^"\n\x1a&<>|\t\v\f\r ,;=\xa0])+))+)?.*|(?:(?=\))|(?=(?:\^[\n\x1a]?)?[\t\v\f\r ,;=\xa0+./:[\\\]]|[\n\x1a&<>|(]))(?:(?:[^\n\x1a^)]|\^[\n\x1a]?[^)])*))`, CommentSingle, Push("follow/compound")},
			{`(set(?:(?=\))|(?=(?:\^[\n\x1a]?)?[\t\v\f\r ,;=\xa0+./:[\\\]]|[\n\x1a&<>|(])))((?:(?:\^[\n\x1a]?)?[^\S\n])*)(/a)`, ByGroups(Keyword, UsingSelf("text"), Keyword), Push("arithmetic/compound")},
			{`(set(?:(?=\))|(?=(?:\^[\n\x1a]?)?[\t\v\f\r ,;=\xa0+./:[\\\]]|[\n\x1a&<>|(])))((?:(?:\^[\n\x1a]?)?[^\S\n])*)((?:/p)?)((?:(?:\^[\n\x1a]?)?[^\S\n])*)((?:(?:(?:\^[\n\x1a]?)?[^"\n\x1a&<>|^=)]|\^[\n\x1a]?[^"=])+)?)((?:(?:\^[\n\x1a]?)?=)?)`, ByGroups(Keyword, UsingSelf("text"), Keyword, UsingSelf("text"), UsingSelf("variable"), Punctuation), Push("follow/compound")},
			Default(Push("follow/compound")),
		},
		"follow/compound": {
			{`(?=\))`, Text, Pop(1)},
			{`((?:(?<=^[^:])|^[^:]?)[\t\v\f\r ,;=\xa0]*)(:)([\t\v\f\r ,;=\xa0]*)((?:(?:[^\n\x1a&<>|\t\v\f\r ,;=\xa0+:^)]|\^[\n\x1a]?[^)])*))(.*)`, ByGroups(Text, Punctuation, Text, NameLabel, CommentSingle), nil},
			Include("redirect/compound"),
			{`(?=[\n\x1a])`, Text, Pop(1)},
			{`\|\|?|&&?`, Punctuation, Pop(1)},
			Include("text"),
		},
		"arithmetic/compound": {
			{`(?=\))`, Text, Pop(1)},
			{`0[0-7]+`, LiteralNumberOct, nil},
			{`0x[\da-f]+`, LiteralNumberHex, nil},
			{`\d+`, LiteralNumberInteger, nil},
			{`[(),]+`, Punctuation, nil},
			{`([=+\-*/!~]|%|\^\^)+`, Operator, nil},
			{`((?:"[^\n\x1a"]*(?:"|(?=[\n\x1a])))|(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))|(\^[\n\x1a]?)?[^()=+\-*/!~%^"\n\x1a&<>|\t\v\f\r ,;=\xa0]|\^[\n\x1a\t\v\f\r ,;=\xa0]?[^)])+`, UsingSelf("variable"), nil},
			{`(?=[\x00|&])`, Text, Pop(1)},
			Include("follow"),
		},
		"call/compound": {
			{`(?=\))`, Text, Pop(1)},
			{`(:?)((?:(?:[^\n\x1a&<>|\t\v\f\r ,;=\xa0+:^)]|\^[\n\x1a]?[^)])*))`, ByGroups(Punctuation, NameLabel), Pop(1)},
		},
		"label/compound": {
			{`(?=\))`, Text, Pop(1)},
			{`((?:(?:[^\n\x1a&<>|\t\v\f\r ,;=\xa0+:^)]|\^[\n\x1a]?[^)])*)?)((?:(?:"[^\n\x1a"]*(?:"|(?=[\n\x1a])))|(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))|\^[\n\x1a]?[^)]|[^"%^\n\x1a&<>|)])*)`, ByGroups(NameLabel, CommentSingle), Pop(1)},
		},
		"redirect/compound": {
			{`((?:(?<=[\n\x1a\t\v\f\r ,;=\xa0])\d)?)(>>?&|<&)([\n\x1a\t\v\f\r ,;=\xa0]*)(\d)`, ByGroups(LiteralNumberInteger, Punctuation, Text, LiteralNumberInteger), nil},
			{`((?:(?<=[\n\x1a\t\v\f\r ,;=\xa0])(?<!\^[\n\x1a])\d)?)(>>?|<)((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)?(?:[&<>|]+|(?:(?:"[^\n\x1a"]*(?:"|(?=[\n\x1a])))|(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))|(?:(?:(?:\^[\n\x1a]?)?[^"\n\x1a&<>|\t\v\f\r ,;=\xa0)])+))+))`, ByGroups(LiteralNumberInteger, Punctuation, UsingSelf("text")), nil},
		},
		"variable-or-escape": {
			{`(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))`, NameVariable, nil},
			{`%%|\^[\n\x1a]?(\^!|[\w\W])`, LiteralStringEscape, nil},
		},
		"string": {
			{`"`, LiteralStringDouble, Pop(1)},
			{`(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))`, NameVariable, nil},
			{`\^!|%%`, LiteralStringEscape, nil},
			{`[^"%^\n\x1a]+|[%^]`, LiteralStringDouble, nil},
			Default(Pop(1)),
		},
		"sqstring": {
			Include("variable-or-escape"),
			{`[^%]+|%`, LiteralStringSingle, nil},
		},
		"bqstring": {
			Include("variable-or-escape"),
			{`[^%]+|%`, LiteralStringBacktick, nil},
		},
		"text": {
			{`"`, LiteralStringDouble, Push("string")},
			Include("variable-or-escape"),
			{`[^"%^\n\x1a&<>|\t\v\f\r ,;=\xa0\d)]+|.`, Text, nil},
		},
		"variable": {
			{`"`, LiteralStringDouble, Push("string")},
			Include("variable-or-escape"),
			{`[^"%^\n\x1a]+|.`, NameVariable, nil},
		},
		"for": {
			{`((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+))(in)((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+))(\()`, ByGroups(UsingSelf("text"), Keyword, UsingSelf("text"), Punctuation), Pop(1)},
			Include("follow"),
		},
		"for2": {
			{`\)`, Punctuation, nil},
			{`((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+))(do(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a]))`, ByGroups(UsingSelf("text"), Keyword), Pop(1)},
			{`[\n\x1a]+`, Text, nil},
			Include("follow"),
		},
		"for/f": {
			{`(")((?:(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))|[^"])*?")([\n\x1a\t\v\f\r ,;=\xa0]*)(\))`, ByGroups(LiteralStringDouble, UsingSelf("string"), Text, Punctuation), nil},
			{`"`, LiteralStringDouble, Push("#pop", "for2", "string")},
			{`('(?:%%|(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))|[\w\W])*?')([\n\x1a\t\v\f\r ,;=\xa0]*)(\))`, ByGroups(UsingSelf("sqstring"), Text, Punctuation), nil},
			{"(`(?:%%|(?:(?:%(?:\\*|(?:~[a-z]*(?:\\$[^:]+:)?)?\\d|[^%:\\n\\x1a]+(?::(?:~(?:-?\\d+)?(?:,(?:-?\\d+)?)?|(?:[^%\\n\\x1a^]|\\^[^%\\n\\x1a])[^=\\n\\x1a]*=(?:[^%\\n\\x1a^]|\\^[^%\\n\\x1a])*)?)?%))|(?:\\^?![^!:\\n\\x1a]+(?::(?:~(?:-?\\d+)?(?:,(?:-?\\d+)?)?|(?:[^!\\n\\x1a^]|\\^[^!\\n\\x1a])[^=\\n\\x1a]*=(?:[^!\\n\\x1a^]|\\^[^!\\n\\x1a])*)?)?\\^?!))|[\\w\\W])*?`)([\\n\\x1a\\t\\v\\f\\r ,;=\\xa0]*)(\\))", ByGroups(UsingSelf("bqstring"), Text, Punctuation), nil},
			Include("for2"),
		},
		"for/l": {
			{`-?\d+`, LiteralNumberInteger, nil},
			Include("for2"),
		},
		"if": {
			{`((?:cmdextversion|errorlevel)(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a]))((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+))(\d+)`, ByGroups(Keyword, UsingSelf("text"), LiteralNumberInteger), Pop(1)},
			{`(defined(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a]))((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+))((?:[&<>|]+|(?:(?:"[^\n\x1a"]*(?:"|(?=[\n\x1a])))|(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))|(?:(?:(?:\^[\n\x1a]?)?[^"\n\x1a&<>|\t\v\f\r ,;=\xa0])+))+))`, ByGroups(Keyword, UsingSelf("text"), UsingSelf("variable")), Pop(1)},
			{`(exist(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a]))((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)(?:[&<>|]+|(?:(?:"[^\n\x1a"]*(?:"|(?=[\n\x1a])))|(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))|(?:(?:(?:\^[\n\x1a]?)?[^"\n\x1a&<>|\t\v\f\r ,;=\xa0])+))+))`, ByGroups(Keyword, UsingSelf("text")), Pop(1)},
			{`((?:-?(?:0[0-7]+|0x[\da-f]+|\d+)(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a]))(?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+))((?:equ|geq|gtr|leq|lss|neq))((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)(?:-?(?:0[0-7]+|0x[\da-f]+|\d+)(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a])))`, ByGroups(UsingSelf("arithmetic"), OperatorWord, UsingSelf("arithmetic")), Pop(1)},
			{`(?:[&<>|]+|(?:(?:"[^\n\x1a"]*(?:"|(?=[\n\x1a])))|(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))|(?:(?:(?:\^[\n\x1a]?)?[^"\n\x1a&<>|\t\v\f\r ,;=\xa0])+))+)`, UsingSelf("text"), Push("#pop", "if2")},
		},
		"if2": {
			{`((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)?)(==)((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)?(?:[&<>|]+|(?:(?:"[^\n\x1a"]*(?:"|(?=[\n\x1a])))|(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))|(?:(?:(?:\^[\n\x1a]?)?[^"\n\x1a&<>|\t\v\f\r ,;=\xa0])+))+))`, ByGroups(UsingSelf("text"), Operator, UsingSelf("text")), Pop(1)},
			{`((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+))((?:equ|geq|gtr|leq|lss|neq))((?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)(?:[&<>|]+|(?:(?:"[^\n\x1a"]*(?:"|(?=[\n\x1a])))|(?:(?:%(?:\*|(?:~[a-z]*(?:\$[^:]+:)?)?\d|[^%:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^%\n\x1a^]|\^[^%\n\x1a])[^=\n\x1a]*=(?:[^%\n\x1a^]|\^[^%\n\x1a])*)?)?%))|(?:\^?![^!:\n\x1a]+(?::(?:~(?:-?\d+)?(?:,(?:-?\d+)?)?|(?:[^!\n\x1a^]|\^[^!\n\x1a])[^=\n\x1a]*=(?:[^!\n\x1a^]|\^[^!\n\x1a])*)?)?\^?!))|(?:(?:(?:\^[\n\x1a]?)?[^"\n\x1a&<>|\t\v\f\r ,;=\xa0])+))+))`, ByGroups(UsingSelf("text"), OperatorWord, UsingSelf("text")), Pop(1)},
		},
		"(?": {
			{`(?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)`, UsingSelf("text"), nil},
			{`\(`, Punctuation, Push("#pop", "else?", "root/compound")},
			Default(Pop(1)),
		},
		"else?": {
			{`(?:(?:(?:\^[\n\x1a])?[\t\v\f\r ,;=\xa0])+)`, UsingSelf("text"), nil},
			{`else(?=\^?[\t\v\f\r ,;=\xa0]|[&<>|\n\x1a])`, Keyword, Pop(1)},
			Default(Pop(1)),
		},
	}
}
