package c

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// CassandraCQL lexer.
var CassandraCQL = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "Cassandra CQL",
		Aliases:         []string{"cassandra", "cql"},
		Filenames:       []string{"*.cql"},
		MimeTypes:       []string{"text/x-cql"},
		NotMultiline:    true,
		CaseInsensitive: true,
	},
	cassandraCQLRules,
))

func cassandraCQLRules() Rules {
	return Rules{
		"root": {
			{`\s+`, TextWhitespace, nil},
			{`(--|\/\/).*\n?`, CommentSingle, nil},
			{`/\*`, CommentMultiline, Push("multiline-comments")},
			{`(ascii|bigint|blob|boolean|counter|date|decimal|double|float|frozen|inet|int|list|map|set|smallint|text|time|timestamp|timeuuid|tinyint|tuple|uuid|varchar|varint)\b`, NameBuiltin, nil},
			{Words(``, `\b`, `ADD`, `AGGREGATE`, `ALL`, `ALLOW`, `ALTER`, `AND`, `ANY`, `APPLY`, `AS`, `ASC`, `AUTHORIZE`, `BATCH`, `BEGIN`, `BY`, `CLUSTERING`, `COLUMNFAMILY`, `COMPACT`, `CONSISTENCY`, `COUNT`, `CREATE`, `CUSTOM`, `DELETE`, `DESC`, `DISTINCT`, `DROP`, `EACH_QUORUM`, `ENTRIES`, `EXISTS`, `FILTERING`, `FROM`, `FULL`, `GRANT`, `IF`, `IN`, `INDEX`, `INFINITY`, `INSERT`, `INTO`, `KEY`, `KEYS`, `KEYSPACE`, `KEYSPACES`, `LEVEL`, `LIMIT`, `LOCAL_ONE`, `LOCAL_QUORUM`, `MATERIALIZED`, `MODIFY`, `NAN`, `NORECURSIVE`, `NOSUPERUSER`, `NOT`, `OF`, `ON`, `ONE`, `ORDER`, `PARTITION`, `PASSWORD`, `PER`, `PERMISSION`, `PERMISSIONS`, `PRIMARY`, `QUORUM`, `RENAME`, `REVOKE`, `SCHEMA`, `SELECT`, `STATIC`, `STORAGE`, `SUPERUSER`, `TABLE`, `THREE`, `TO`, `TOKEN`, `TRUNCATE`, `TTL`, `TWO`, `TYPE`, `UNLOGGED`, `UPDATE`, `USE`, `USER`, `USERS`, `USING`, `VALUES`, `VIEW`, `WHERE`, `WITH`, `WRITETIME`, `REPLICATION`, `OR`, `REPLACE`, `FUNCTION`, `CALLED`, `INPUT`, `RETURNS`, `LANGUAGE`, `ROLE`, `ROLES`, `TRIGGER`, `DURABLE_WRITES`, `LOGIN`, `OPTIONS`, `LOGGED`, `SFUNC`, `STYPE`, `FINALFUNC`, `INITCOND`, `IS`, `CONTAINS`, `JSON`, `PAGING`, `OFF`), Keyword, nil},
			{"[+*/<>=~!@#%^&|`?-]+", Operator, nil},
			{
				`(?s)(java|javascript)(\s+)(AS)(\s+)('|\$\$)(.*?)(\5)`,
				UsingByGroup(
					internal.Get,
					1, 6,
					NameBuiltin, TextWhitespace, Keyword, TextWhitespace,
					LiteralStringHeredoc, LiteralStringHeredoc, LiteralStringHeredoc,
				),
				nil,
			},
			{`(true|false|null)\b`, KeywordConstant, nil},
			{`0x[0-9a-f]+`, LiteralNumberHex, nil},
			{`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`, LiteralNumberHex, nil},
			{`\.[0-9]+(e[+-]?[0-9]+)?`, Error, nil},
			{`-?[0-9]+(\.[0-9])?(e[+-]?[0-9]+)?`, LiteralNumberFloat, nil},
			{`[0-9]+`, LiteralNumberInteger, nil},
			{`'`, LiteralStringSingle, Push("string")},
			{`"`, LiteralStringName, Push("quoted-ident")},
			{`\$\$`, LiteralStringHeredoc, Push("dollar-string")},
			{`[a-z_]\w*`, Name, nil},
			{`:(['"]?)[a-z]\w*\b\1`, NameVariable, nil},
			{`[;:()\[\]\{\},.]`, Punctuation, nil},
		},
		"multiline-comments": {
			{`/\*`, CommentMultiline, Push("multiline-comments")},
			{`\*/`, CommentMultiline, Pop(1)},
			{`[^/*]+`, CommentMultiline, nil},
			{`[/*]`, CommentMultiline, nil},
		},
		"string": {
			{`[^']+`, LiteralStringSingle, nil},
			{`''`, LiteralStringSingle, nil},
			{`'`, LiteralStringSingle, Pop(1)},
		},
		"quoted-ident": {
			{`[^"]+`, LiteralStringName, nil},
			{`""`, LiteralStringName, nil},
			{`"`, LiteralStringName, Pop(1)},
		},
		"dollar-string": {
			{`[^\$]+`, LiteralStringHeredoc, nil},
			{`\$\$`, LiteralStringHeredoc, Pop(1)},
		},
	}
}
