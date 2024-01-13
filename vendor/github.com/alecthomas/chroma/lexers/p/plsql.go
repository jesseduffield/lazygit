package p

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Pl/Pgsql lexer.
var PLpgSQL = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "PL/pgSQL",
		Aliases:         []string{"plpgsql"},
		Filenames:       []string{},
		MimeTypes:       []string{"text/x-plpgsql"},
		NotMultiline:    true,
		CaseInsensitive: true,
	},
	plpgSQLRules,
))

func plpgSQLRules() Rules {
	return Rules{
		"root": {
			{`\%[a-z]\w*\b`, NameBuiltin, nil},
			{`:=`, Operator, nil},
			{`\<\<[a-z]\w*\>\>`, NameLabel, nil},
			{`\#[a-z]\w*\b`, KeywordPseudo, nil},
			{`\s+`, TextWhitespace, nil},
			{`--.*\n?`, CommentSingle, nil},
			{`/\*`, CommentMultiline, Push("multiline-comments")},
			{`(bigint|bigserial|bit|bit\s+varying|bool|boolean|box|bytea|char|character|character\s+varying|cidr|circle|date|decimal|double\s+precision|float4|float8|inet|int|int2|int4|int8|integer|interval|json|jsonb|line|lseg|macaddr|money|numeric|path|pg_lsn|point|polygon|real|serial|serial2|serial4|serial8|smallint|smallserial|text|time|timestamp|timestamptz|timetz|tsquery|tsvector|txid_snapshot|uuid|varbit|varchar|with\s+time\s+zone|without\s+time\s+zone|xml|anyarray|anyelement|anyenum|anynonarray|anyrange|cstring|fdw_handler|internal|language_handler|opaque|record|void)\b`, NameBuiltin, nil},
			{Words(``, `\b`, `ABORT`, `ABSOLUTE`, `ACCESS`, `ACTION`, `ADD`, `ADMIN`, `AFTER`, `AGGREGATE`, `ALL`, `ALSO`, `ALTER`, `ALWAYS`, `ANALYSE`, `ANALYZE`, `AND`, `ANY`, `ARRAY`, `AS`, `ASC`, `ASSERTION`, `ASSIGNMENT`, `ASYMMETRIC`, `AT`, `ATTRIBUTE`, `AUTHORIZATION`, `BACKWARD`, `BEFORE`, `BEGIN`, `BETWEEN`, `BIGINT`, `BINARY`, `BIT`, `BOOLEAN`, `BOTH`, `BY`, `CACHE`, `CALLED`, `CASCADE`, `CASCADED`, `CASE`, `CAST`, `CATALOG`, `CHAIN`, `CHAR`, `CHARACTER`, `CHARACTERISTICS`, `CHECK`, `CHECKPOINT`, `CLASS`, `CLOSE`, `CLUSTER`, `COALESCE`, `COLLATE`, `COLLATION`, `COLUMN`, `COMMENT`, `COMMENTS`, `COMMIT`, `COMMITTED`, `CONCURRENTLY`, `CONFIGURATION`, `CONNECTION`, `CONSTRAINT`, `CONSTRAINTS`, `CONTENT`, `CONTINUE`, `CONVERSION`, `COPY`, `COST`, `CREATE`, `CROSS`, `CSV`, `CURRENT`, `CURRENT_CATALOG`, `CURRENT_DATE`, `CURRENT_ROLE`, `CURRENT_SCHEMA`, `CURRENT_TIME`, `CURRENT_TIMESTAMP`, `CURRENT_USER`, `CURSOR`, `CYCLE`, `DATA`, `DATABASE`, `DAY`, `DEALLOCATE`, `DEC`, `DECIMAL`, `DECLARE`, `DEFAULT`, `DEFAULTS`, `DEFERRABLE`, `DEFERRED`, `DEFINER`, `DELETE`, `DELIMITER`, `DELIMITERS`, `DESC`, `DICTIONARY`, `DISABLE`, `DISCARD`, `DISTINCT`, `DO`, `DOCUMENT`, `DOMAIN`, `DOUBLE`, `DROP`, `EACH`, `ELSE`, `ENABLE`, `ENCODING`, `ENCRYPTED`, `END`, `ENUM`, `ESCAPE`, `EVENT`, `EXCEPT`, `EXCLUDE`, `EXCLUDING`, `EXCLUSIVE`, `EXECUTE`, `EXISTS`, `EXPLAIN`, `EXTENSION`, `EXTERNAL`, `EXTRACT`, `FALSE`, `FAMILY`, `FETCH`, `FILTER`, `FIRST`, `FLOAT`, `FOLLOWING`, `FOR`, `FORCE`, `FOREIGN`, `FORWARD`, `FREEZE`, `FROM`, `FULL`, `FUNCTION`, `FUNCTIONS`, `GLOBAL`, `GRANT`, `GRANTED`, `GREATEST`, `GROUP`, `HANDLER`, `HAVING`, `HEADER`, `HOLD`, `HOUR`, `IDENTITY`, `IF`, `ILIKE`, `IMMEDIATE`, `IMMUTABLE`, `IMPLICIT`, `IN`, `INCLUDING`, `INCREMENT`, `INDEX`, `INDEXES`, `INHERIT`, `INHERITS`, `INITIALLY`, `INLINE`, `INNER`, `INOUT`, `INPUT`, `INSENSITIVE`, `INSERT`, `INSTEAD`, `INT`, `INTEGER`, `INTERSECT`, `INTERVAL`, `INTO`, `INVOKER`, `IS`, `ISNULL`, `ISOLATION`, `JOIN`, `KEY`, `LABEL`, `LANGUAGE`, `LARGE`, `LAST`, `LATERAL`, `LC_COLLATE`, `LC_CTYPE`, `LEADING`, `LEAKPROOF`, `LEAST`, `LEFT`, `LEVEL`, `LIKE`, `LIMIT`, `LISTEN`, `LOAD`, `LOCAL`, `LOCALTIME`, `LOCALTIMESTAMP`, `LOCATION`, `LOCK`, `MAPPING`, `MATCH`, `MATERIALIZED`, `MAXVALUE`, `MINUTE`, `MINVALUE`, `MODE`, `MONTH`, `MOVE`, `NAME`, `NAMES`, `NATIONAL`, `NATURAL`, `NCHAR`, `NEXT`, `NO`, `NONE`, `NOT`, `NOTHING`, `NOTIFY`, `NOTNULL`, `NOWAIT`, `NULL`, `NULLIF`, `NULLS`, `NUMERIC`, `OBJECT`, `OF`, `OFF`, `OFFSET`, `OIDS`, `ON`, `ONLY`, `OPERATOR`, `OPTION`, `OPTIONS`, `OR`, `ORDER`, `ORDINALITY`, `OUT`, `OUTER`, `OVER`, `OVERLAPS`, `OVERLAY`, `OWNED`, `OWNER`, `PARSER`, `PARTIAL`, `PARTITION`, `PASSING`, `PASSWORD`, `PLACING`, `PLANS`, `POLICY`, `POSITION`, `PRECEDING`, `PRECISION`, `PREPARE`, `PREPARED`, `PRESERVE`, `PRIMARY`, `PRIOR`, `PRIVILEGES`, `PROCEDURAL`, `PROCEDURE`, `PROGRAM`, `QUOTE`, `RANGE`, `READ`, `REAL`, `REASSIGN`, `RECHECK`, `RECURSIVE`, `REF`, `REFERENCES`, `REFRESH`, `REINDEX`, `RELATIVE`, `RELEASE`, `RENAME`, `REPEATABLE`, `REPLACE`, `REPLICA`, `RESET`, `RESTART`, `RESTRICT`, `RETURNING`, `RETURNS`, `REVOKE`, `RIGHT`, `ROLE`, `ROLLBACK`, `ROW`, `ROWS`, `RULE`, `SAVEPOINT`, `SCHEMA`, `SCROLL`, `SEARCH`, `SECOND`, `SECURITY`, `SELECT`, `SEQUENCE`, `SEQUENCES`, `SERIALIZABLE`, `SERVER`, `SESSION`, `SESSION_USER`, `SET`, `SETOF`, `SHARE`, `SHOW`, `SIMILAR`, `SIMPLE`, `SMALLINT`, `SNAPSHOT`, `SOME`, `STABLE`, `STANDALONE`, `START`, `STATEMENT`, `STATISTICS`, `STDIN`, `STDOUT`, `STORAGE`, `STRICT`, `STRIP`, `SUBSTRING`, `SYMMETRIC`, `SYSID`, `SYSTEM`, `TABLE`, `TABLES`, `TABLESPACE`, `TEMP`, `TEMPLATE`, `TEMPORARY`, `TEXT`, `THEN`, `TIME`, `TIMESTAMP`, `TO`, `TRAILING`, `TRANSACTION`, `TREAT`, `TRIGGER`, `TRIM`, `TRUE`, `TRUNCATE`, `TRUSTED`, `TYPE`, `TYPES`, `UNBOUNDED`, `UNCOMMITTED`, `UNENCRYPTED`, `UNION`, `UNIQUE`, `UNKNOWN`, `UNLISTEN`, `UNLOGGED`, `UNTIL`, `UPDATE`, `USER`, `USING`, `VACUUM`, `VALID`, `VALIDATE`, `VALIDATOR`, `VALUE`, `VALUES`, `VARCHAR`, `VARIADIC`, `VARYING`, `VERBOSE`, `VERSION`, `VIEW`, `VIEWS`, `VOLATILE`, `WHEN`, `WHERE`, `WHITESPACE`, `WINDOW`, `WITH`, `WITHIN`, `WITHOUT`, `WORK`, `WRAPPER`, `WRITE`, `XML`, `XMLATTRIBUTES`, `XMLCONCAT`, `XMLELEMENT`, `XMLEXISTS`, `XMLFOREST`, `XMLPARSE`, `XMLPI`, `XMLROOT`, `XMLSERIALIZE`, `YEAR`, `YES`, `ZONE`, `ALIAS`, `CONSTANT`, `DIAGNOSTICS`, `ELSIF`, `EXCEPTION`, `EXIT`, `FOREACH`, `GET`, `LOOP`, `NOTICE`, `OPEN`, `PERFORM`, `QUERY`, `RAISE`, `RETURN`, `REVERSE`, `SQLSTATE`, `WHILE`), Keyword, nil},
			{"[+*/<>=~!@#%^&|`?-]+", Operator, nil},
			{`::`, Operator, nil},
			{`\$\d+`, NameVariable, nil},
			{`([0-9]*\.[0-9]*|[0-9]+)(e[+-]?[0-9]+)?`, LiteralNumberFloat, nil},
			{`[0-9]+`, LiteralNumberInteger, nil},
			{`((?:E|U&)?)(')`, ByGroups(LiteralStringAffix, LiteralStringSingle), Push("string")},
			{`((?:U&)?)(")`, ByGroups(LiteralStringAffix, LiteralStringName), Push("quoted-ident")},
			// { `(?s)(\$)([^$]*)(\$)(.*?)(\$)(\2)(\$)`, ?? <function language_callback at 0x108964400> ??, nil },
			{`[a-z_]\w*`, Name, nil},
			{`:(['"]?)[a-z]\w*\b\1`, NameVariable, nil},
			{`[;:()\[\]{},.]`, Punctuation, nil},
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
	}
}
