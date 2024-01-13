package y

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

var YANG = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "YANG",
		Aliases:   []string{"yang"},
		Filenames: []string{"*.yang"},
		MimeTypes: []string{"application/yang"},
	},
	yangRules,
))

func yangRules() Rules {
	return Rules{
		"root": {
			{`\s+`, Whitespace, nil},
			{`[\{\}\;]+`, Punctuation, nil},
			{`(?<![\-\w])(and|or|not|\+|\.)(?![\-\w])`, Operator, nil},

			{`"(?:\\"|[^"])*?"`, StringDouble, nil},
			{`'(?:\\'|[^'])*?'`, StringSingle, nil},

			{`/\*`, CommentMultiline, Push("comments")},
			{`//.*?$`, CommentSingle, nil},

			// match BNF stmt for `node-identifier` with [ prefix ":"]
			{`(?:^|(?<=[\s{};]))([\w.-]+)(:)([\w.-]+)(?=[\s{};])`, ByGroups(KeywordNamespace, Punctuation, Text), nil},

			// match BNF stmt `date-arg-str`
			{`([0-9]{4}\-[0-9]{2}\-[0-9]{2})(?=[\s\{\}\;])`, LiteralDate, nil},
			{`([0-9]+\.[0-9]+)(?=[\s\{\}\;])`, NumberFloat, nil},
			{`([0-9]+)(?=[\s\{\}\;])`, NumberInteger, nil},

			// TOP_STMTS_KEYWORDS
			{Words(``, `(?=[^\w\-\:])`, `module`, `submodule`), Keyword, nil},
			// MODULE_HEADER_STMT_KEYWORDS
			{Words(``, `(?=[^\w\-\:])`, `belongs-to`, `namespace`, `prefix`, `yang-version`), Keyword, nil},
			// META_STMT_KEYWORDS
			{Words(``, `(?=[^\w\-\:])`, `contact`, `description`, `organization`, `reference`, `revision`), Keyword, nil},
			// LINKAGE_STMTS_KEYWORDS
			{Words(``, `(?=[^\w\-\:])`, `import`, `include`, `revision-date`), Keyword, nil},
			// BODY_STMT_KEYWORDS
			{Words(``, `(?=[^\w\-\:])`, `action`, `argument`, `augment`, `deviation`, `extension`, `feature`, `grouping`, `identity`, `if-feature`, `input`, `notification`, `output`, `rpc`, `typedef`), Keyword, nil},
			// DATA_DEF_STMT_KEYWORDS
			{Words(``, `(?=[^\w\-\:])`, `anydata`, `anyxml`, `case`, `choice`, `config`, `container`, `deviate`, `leaf`, `leaf-list`, `list`, `must`, `presence`, `refine`, `uses`, `when`), Keyword, nil},
			// TYPE_STMT_KEYWORDS
			{Words(``, `(?=[^\w\-\:])`, `base`, `bit`, `default`, `enum`, `error-app-tag`, `error-message`, `fraction-digits`, `length`, `max-elements`, `min-elements`, `modifier`, `ordered-by`, `path`, `pattern`, `position`, `range`, `require-instance`, `status`, `type`, `units`, `value`, `yin-element`), Keyword, nil},
			// LIST_STMT_KEYWORDS
			{Words(``, `(?=[^\w\-\:])`, `key`, `mandatory`, `unique`), Keyword, nil},

			// CONSTANTS_KEYWORDS - RFC7950 other keywords
			{Words(``, `(?=[^\w\-\:])`, `add`, `current`, `delete`, `deprecated`, `false`, `invert-match`, `max`, `min`, `not-supported`, `obsolete`, `replace`, `true`, `unbounded`, `user`), NameClass, nil},

			// RFC7950 Built-In Types
			{Words(``, `(?=[^\w\-\:])`, `binary`, `bits`, `boolean`, `decimal64`, `empty`, `enumeration`, `identityref`, `instance-identifier`, `int16`, `int32`, `int64`, `int8`, `leafref`, `string`, `uint16`, `uint32`, `uint64`, `uint8`, `union`), NameClass, nil},

			{`[^;{}\s\'\"]+`, Text, nil},
		},
		"comments": {
			{`[^*/]`, CommentMultiline, nil},
			{`/\*`, CommentMultiline, Push("comment")},
			{`\*/`, CommentMultiline, Pop(1)},
			{`[*/]`, CommentMultiline, nil},
		},
	}
}
