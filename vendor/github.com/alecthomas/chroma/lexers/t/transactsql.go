package t

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// TransactSQL lexer.
var TransactSQL = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "Transact-SQL",
		Aliases:         []string{"tsql", "t-sql"},
		MimeTypes:       []string{"text/x-tsql"},
		NotMultiline:    true,
		CaseInsensitive: true,
	},
	transactSQLRules,
))

func transactSQLRules() Rules {
	return Rules{
		"root": {
			{`\s+`, TextWhitespace, nil},
			{`--(?m).*?$\n?`, CommentSingle, nil},
			{`/\*`, CommentMultiline, Push("multiline-comments")},
			{`'`, LiteralStringSingle, Push("string")},
			{`"`, LiteralStringName, Push("quoted-ident")},
			{Words(``, ``, `!<`, `!=`, `!>`, `<`, `<=`, `<>`, `=`, `>`, `>=`, `+`, `+=`, `-`, `-=`, `*`, `*=`, `/`, `/=`, `%`, `%=`, `&`, `&=`, `|`, `|=`, `^`, `^=`, `~`, `::`), Operator, nil},
			{Words(``, `\b`, `all`, `and`, `any`, `between`, `except`, `exists`, `in`, `intersect`, `like`, `not`, `or`, `some`, `union`), OperatorWord, nil},
			{Words(``, `\b`, `bigint`, `binary`, `bit`, `char`, `cursor`, `date`, `datetime`, `datetime2`, `datetimeoffset`, `decimal`, `float`, `hierarchyid`, `image`, `int`, `money`, `nchar`, `ntext`, `numeric`, `nvarchar`, `real`, `smalldatetime`, `smallint`, `smallmoney`, `sql_variant`, `table`, `text`, `time`, `timestamp`, `tinyint`, `uniqueidentifier`, `varbinary`, `varchar`, `xml`), NameClass, nil},
			{Words(``, `\b`, `$partition`, `abs`, `acos`, `app_name`, `applock_mode`, `applock_test`, `ascii`, `asin`, `assemblyproperty`, `atan`, `atn2`, `avg`, `binary_checksum`, `cast`, `ceiling`, `certencoded`, `certprivatekey`, `char`, `charindex`, `checksum`, `checksum_agg`, `choose`, `col_length`, `col_name`, `columnproperty`, `compress`, `concat`, `connectionproperty`, `context_info`, `convert`, `cos`, `cot`, `count`, `count_big`, `current_request_id`, `current_timestamp`, `current_transaction_id`, `current_user`, `cursor_status`, `database_principal_id`, `databasepropertyex`, `dateadd`, `datediff`, `datediff_big`, `datefromparts`, `datename`, `datepart`, `datetime2fromparts`, `datetimefromparts`, `datetimeoffsetfromparts`, `day`, `db_id`, `db_name`, `decompress`, `degrees`, `dense_rank`, `difference`, `eomonth`, `error_line`, `error_message`, `error_number`, `error_procedure`, `error_severity`, `error_state`, `exp`, `file_id`, `file_idex`, `file_name`, `filegroup_id`, `filegroup_name`, `filegroupproperty`, `fileproperty`, `floor`, `format`, `formatmessage`, `fulltextcatalogproperty`, `fulltextserviceproperty`, `get_filestream_transaction_context`, `getansinull`, `getdate`, `getutcdate`, `grouping`, `grouping_id`, `has_perms_by_name`, `host_id`, `host_name`, `iif`, `index_col`, `indexkey_property`, `indexproperty`, `is_member`, `is_rolemember`, `is_srvrolemember`, `isdate`, `isjson`, `isnull`, `isnumeric`, `json_modify`, `json_query`, `json_value`, `left`, `len`, `log`, `log10`, `lower`, `ltrim`, `max`, `min`, `min_active_rowversion`, `month`, `nchar`, `newid`, `newsequentialid`, `ntile`, `object_definition`, `object_id`, `object_name`, `object_schema_name`, `objectproperty`, `objectpropertyex`, `opendatasource`, `openjson`, `openquery`, `openrowset`, `openxml`, `original_db_name`, `original_login`, `parse`, `parsename`, `patindex`, `permissions`, `pi`, `power`, `pwdcompare`, `pwdencrypt`, `quotename`, `radians`, `rand`, `rank`, `replace`, `replicate`, `reverse`, `right`, `round`, `row_number`, `rowcount_big`, `rtrim`, `schema_id`, `schema_name`, `scope_identity`, `serverproperty`, `session_context`, `session_user`, `sign`, `sin`, `smalldatetimefromparts`, `soundex`, `sp_helplanguage`, `space`, `sqrt`, `square`, `stats_date`, `stdev`, `stdevp`, `str`, `string_escape`, `string_split`, `stuff`, `substring`, `sum`, `suser_id`, `suser_name`, `suser_sid`, `suser_sname`, `switchoffset`, `sysdatetime`, `sysdatetimeoffset`, `system_user`, `sysutcdatetime`, `tan`, `textptr`, `textvalid`, `timefromparts`, `todatetimeoffset`, `try_cast`, `try_convert`, `try_parse`, `type_id`, `type_name`, `typeproperty`, `unicode`, `upper`, `user_id`, `user_name`, `var`, `varp`, `xact_state`, `year`), NameFunction, nil},
			{`(goto)(\s+)(\w+\b)`, ByGroups(Keyword, TextWhitespace, NameLabel), nil},
			{Words(``, `\b`, `absolute`, `action`, `ada`, `add`, `admin`, `after`, `aggregate`, `alias`, `all`, `allocate`, `alter`, `and`, `any`, `are`, `array`, `as`, `asc`, `asensitive`, `assertion`, `asymmetric`, `at`, `atomic`, `authorization`, `avg`, `backup`, `before`, `begin`, `between`, `binary`, `bit`, `bit_length`, `blob`, `boolean`, `both`, `breadth`, `break`, `browse`, `bulk`, `by`, `call`, `called`, `cardinality`, `cascade`, `cascaded`, `case`, `cast`, `catalog`, `catch`, `char`, `char_length`, `character`, `character_length`, `check`, `checkpoint`, `class`, `clob`, `close`, `clustered`, `coalesce`, `collate`, `collation`, `collect`, `column`, `commit`, `completion`, `compute`, `condition`, `connect`, `connection`, `constraint`, `constraints`, `constructor`, `contains`, `containstable`, `continue`, `convert`, `corr`, `corresponding`, `count`, `covar_pop`, `covar_samp`, `create`, `cross`, `cube`, `cume_dist`, `current`, `current_catalog`, `current_date`, `current_default_transform_group`, `current_path`, `current_role`, `current_schema`, `current_time`, `current_timestamp`, `current_transform_group_for_type`, `current_user`, `cursor`, `cycle`, `data`, `database`, `date`, `day`, `dbcc`, `deallocate`, `dec`, `decimal`, `declare`, `default`, `deferrable`, `deferred`, `delete`, `deny`, `depth`, `deref`, `desc`, `describe`, `descriptor`, `destroy`, `destructor`, `deterministic`, `diagnostics`, `dictionary`, `disconnect`, `disk`, `distinct`, `distributed`, `domain`, `double`, `drop`, `dump`, `dynamic`, `each`, `element`, `else`, `end`, `end-exec`, `equals`, `errlvl`, `escape`, `every`, `except`, `exception`, `exec`, `execute`, `exists`, `exit`, `external`, `extract`, `false`, `fetch`, `file`, `fillfactor`, `filter`, `first`, `float`, `for`, `foreign`, `fortran`, `found`, `free`, `freetext`, `freetexttable`, `from`, `full`, `fulltexttable`, `function`, `fusion`, `general`, `get`, `global`, `go`, `goto`, `grant`, `group`, `grouping`, `having`, `hold`, `holdlock`, `host`, `hour`, `identity`, `identity_insert`, `identitycol`, `if`, `ignore`, `immediate`, `in`, `include`, `index`, `indicator`, `initialize`, `initially`, `inner`, `inout`, `input`, `insensitive`, `insert`, `int`, `integer`, `intersect`, `intersection`, `interval`, `into`, `is`, `isolation`, `iterate`, `join`, `key`, `kill`, `language`, `large`, `last`, `lateral`, `leading`, `left`, `less`, `level`, `like`, `like_regex`, `limit`, `lineno`, `ln`, `load`, `local`, `localtime`, `localtimestamp`, `locator`, `lower`, `map`, `match`, `max`, `member`, `merge`, `method`, `min`, `minute`, `mod`, `modifies`, `modify`, `module`, `month`, `multiset`, `names`, `national`, `natural`, `nchar`, `nclob`, `new`, `next`, `no`, `nocheck`, `nonclustered`, `none`, `normalize`, `not`, `null`, `nullif`, `numeric`, `object`, `occurrences_regex`, `octet_length`, `of`, `off`, `offsets`, `old`, `on`, `only`, `open`, `opendatasource`, `openquery`, `openrowset`, `openxml`, `operation`, `option`, `or`, `order`, `ordinality`, `out`, `outer`, `output`, `over`, `overlaps`, `overlay`, `pad`, `parameter`, `parameters`, `partial`, `partition`, `pascal`, `path`, `percent`, `percent_rank`, `percentile_cont`, `percentile_disc`, `pivot`, `plan`, `position`, `position_regex`, `postfix`, `precision`, `prefix`, `preorder`, `prepare`, `preserve`, `primary`, `print`, `prior`, `privileges`, `proc`, `procedure`, `public`, `raiserror`, `range`, `read`, `reads`, `readtext`, `real`, `reconfigure`, `recursive`, `ref`, `references`, `referencing`, `regr_avgx`, `regr_avgy`, `regr_count`, `regr_intercept`, `regr_r2`, `regr_slope`, `regr_sxx`, `regr_sxy`, `regr_syy`, `relative`, `release`, `replication`, `restore`, `restrict`, `result`, `return`, `returns`, `revert`, `revoke`, `right`, `role`, `rollback`, `rollup`, `routine`, `row`, `rowcount`, `rowguidcol`, `rows`, `rule`, `save`, `savepoint`, `schema`, `scope`, `scroll`, `search`, `second`, `section`, `securityaudit`, `select`, `semantickeyphrasetable`, `semanticsimilaritydetailstable`, `semanticsimilaritytable`, `sensitive`, `sequence`, `session`, `session_user`, `set`, `sets`, `setuser`, `shutdown`, `similar`, `size`, `smallint`, `some`, `space`, `specific`, `specifictype`, `sql`, `sqlca`, `sqlcode`, `sqlerror`, `sqlexception`, `sqlstate`, `sqlwarning`, `start`, `state`, `statement`, `static`, `statistics`, `stddev_pop`, `stddev_samp`, `structure`, `submultiset`, `substring`, `substring_regex`, `sum`, `symmetric`, `system`, `system_user`, `table`, `tablesample`, `temporary`, `terminate`, `textsize`, `than`, `then`, `throw`, `time`, `timestamp`, `timezone_hour`, `timezone_minute`, `to`, `top`, `trailing`, `tran`, `transaction`, `translate`, `translate_regex`, `translation`, `treat`, `trigger`, `trim`, `true`, `truncate`, `try`, `try_convert`, `tsequal`, `uescape`, `under`, `union`, `unique`, `unknown`, `unnest`, `unpivot`, `update`, `updatetext`, `upper`, `usage`, `use`, `user`, `using`, `value`, `values`, `var_pop`, `var_samp`, `varchar`, `variable`, `varying`, `view`, `waitfor`, `when`, `whenever`, `where`, `while`, `width_bucket`, `window`, `with`, `within`, `without`, `work`, `write`, `writetext`, `xmlagg`, `xmlattributes`, `xmlbinary`, `xmlcast`, `xmlcomment`, `xmlconcat`, `xmldocument`, `xmlelement`, `xmlexists`, `xmlforest`, `xmliterate`, `xmlnamespaces`, `xmlparse`, `xmlpi`, `xmlquery`, `xmlserialize`, `xmltable`, `xmltext`, `xmlvalidate`, `year`, `zone`), Keyword, nil},
			{`(\[)([^]]+)(\])`, ByGroups(Operator, Name, Operator), nil},
			{`0x[0-9a-f]+`, LiteralNumberHex, nil},
			{`[0-9]+\.[0-9]*(e[+-]?[0-9]+)?`, LiteralNumberFloat, nil},
			{`\.[0-9]+(e[+-]?[0-9]+)?`, LiteralNumberFloat, nil},
			{`[0-9]+e[+-]?[0-9]+`, LiteralNumberFloat, nil},
			{`[0-9]+`, LiteralNumberInteger, nil},
			{`[;(),.]`, Punctuation, nil},
			{`@@\w+`, NameBuiltin, nil},
			{`@\w+`, NameVariable, nil},
			{`(\w+)(:)`, ByGroups(NameLabel, Punctuation), nil},
			{`#?#?\w+`, Name, nil},
			{`\?`, NameVariableMagic, nil},
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
