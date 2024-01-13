package s

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Sas lexer.
var Sas = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "SAS",
		Aliases:         []string{"sas"},
		Filenames:       []string{"*.SAS", "*.sas"},
		MimeTypes:       []string{"text/x-sas", "text/sas", "application/x-sas"},
		CaseInsensitive: true,
	},
	sasRules,
))

func sasRules() Rules {
	return Rules{
		"root": {
			Include("comments"),
			Include("proc-data"),
			Include("cards-datalines"),
			Include("logs"),
			Include("general"),
			{`.`, Text, nil},
			{`\\\n`, Text, nil},
			{`\n`, Text, nil},
		},
		"comments": {
			{`^\s*\*.*?;`, Comment, nil},
			{`/\*.*?\*/`, Comment, nil},
			{`^\s*\*(.|\n)*?;`, CommentMultiline, nil},
			{`/[*](.|\n)*?[*]/`, CommentMultiline, nil},
		},
		"proc-data": {
			{`(^|;)\s*(proc \w+|data|run|quit)[\s;]`, KeywordReserved, nil},
		},
		"cards-datalines": {
			{`^\s*(datalines|cards)\s*;\s*$`, Keyword, Push("data")},
		},
		"data": {
			{`(.|\n)*^\s*;\s*$`, Other, Pop(1)},
		},
		"logs": {
			{`\n?^\s*%?put `, Keyword, Push("log-messages")},
		},
		"log-messages": {
			{`NOTE(:|-).*`, Generic, Pop(1)},
			{`WARNING(:|-).*`, GenericEmph, Pop(1)},
			{`ERROR(:|-).*`, GenericError, Pop(1)},
			Include("general"),
		},
		"general": {
			Include("keywords"),
			Include("vars-strings"),
			Include("special"),
			Include("numbers"),
		},
		"keywords": {
			{Words(`\b`, `\b`, `abort`, `array`, `attrib`, `by`, `call`, `cards`, `cards4`, `catname`, `continue`, `datalines`, `datalines4`, `delete`, `delim`, `delimiter`, `display`, `dm`, `drop`, `endsas`, `error`, `file`, `filename`, `footnote`, `format`, `goto`, `in`, `infile`, `informat`, `input`, `keep`, `label`, `leave`, `length`, `libname`, `link`, `list`, `lostcard`, `merge`, `missing`, `modify`, `options`, `output`, `out`, `page`, `put`, `redirect`, `remove`, `rename`, `replace`, `retain`, `return`, `select`, `set`, `skip`, `startsas`, `stop`, `title`, `update`, `waitsas`, `where`, `window`, `x`, `systask`), Keyword, nil},
			{Words(`\b`, `\b`, `add`, `and`, `alter`, `as`, `cascade`, `check`, `create`, `delete`, `describe`, `distinct`, `drop`, `foreign`, `from`, `group`, `having`, `index`, `insert`, `into`, `in`, `key`, `like`, `message`, `modify`, `msgtype`, `not`, `null`, `on`, `or`, `order`, `primary`, `references`, `reset`, `restrict`, `select`, `set`, `table`, `unique`, `update`, `validate`, `view`, `where`), Keyword, nil},
			{Words(`\b`, `\b`, `do`, `if`, `then`, `else`, `end`, `until`, `while`), Keyword, nil},
			{Words(`%`, `\b`, `bquote`, `nrbquote`, `cmpres`, `qcmpres`, `compstor`, `datatyp`, `display`, `do`, `else`, `end`, `eval`, `global`, `goto`, `if`, `index`, `input`, `keydef`, `label`, `left`, `length`, `let`, `local`, `lowcase`, `macro`, `mend`, `nrquote`, `nrstr`, `put`, `qleft`, `qlowcase`, `qscan`, `qsubstr`, `qsysfunc`, `qtrim`, `quote`, `qupcase`, `scan`, `str`, `substr`, `superq`, `syscall`, `sysevalf`, `sysexec`, `sysfunc`, `sysget`, `syslput`, `sysprod`, `sysrc`, `sysrput`, `then`, `to`, `trim`, `unquote`, `until`, `upcase`, `verify`, `while`, `window`), NameBuiltin, nil},
			{Words(`\b`, `\(`, `abs`, `addr`, `airy`, `arcos`, `arsin`, `atan`, `attrc`, `attrn`, `band`, `betainv`, `blshift`, `bnot`, `bor`, `brshift`, `bxor`, `byte`, `cdf`, `ceil`, `cexist`, `cinv`, `close`, `cnonct`, `collate`, `compbl`, `compound`, `compress`, `cos`, `cosh`, `css`, `curobs`, `cv`, `daccdb`, `daccdbsl`, `daccsl`, `daccsyd`, `dacctab`, `dairy`, `date`, `datejul`, `datepart`, `datetime`, `day`, `dclose`, `depdb`, `depdbsl`, `depsl`, `depsyd`, `deptab`, `dequote`, `dhms`, `dif`, `digamma`, `dim`, `dinfo`, `dnum`, `dopen`, `doptname`, `doptnum`, `dread`, `dropnote`, `dsname`, `erf`, `erfc`, `exist`, `exp`, `fappend`, `fclose`, `fcol`, `fdelete`, `fetch`, `fetchobs`, `fexist`, `fget`, `fileexist`, `filename`, `fileref`, `finfo`, `finv`, `fipname`, `fipnamel`, `fipstate`, `floor`, `fnonct`, `fnote`, `fopen`, `foptname`, `foptnum`, `fpoint`, `fpos`, `fput`, `fread`, `frewind`, `frlen`, `fsep`, `fuzz`, `fwrite`, `gaminv`, `gamma`, `getoption`, `getvarc`, `getvarn`, `hbound`, `hms`, `hosthelp`, `hour`, `ibessel`, `index`, `indexc`, `indexw`, `input`, `inputc`, `inputn`, `int`, `intck`, `intnx`, `intrr`, `irr`, `jbessel`, `juldate`, `kurtosis`, `lag`, `lbound`, `left`, `length`, `lgamma`, `libname`, `libref`, `log`, `log10`, `log2`, `logpdf`, `logpmf`, `logsdf`, `lowcase`, `max`, `mdy`, `mean`, `min`, `minute`, `mod`, `month`, `mopen`, `mort`, `n`, `netpv`, `nmiss`, `normal`, `note`, `npv`, `open`, `ordinal`, `pathname`, `pdf`, `peek`, `peekc`, `pmf`, `point`, `poisson`, `poke`, `probbeta`, `probbnml`, `probchi`, `probf`, `probgam`, `probhypr`, `probit`, `probnegb`, `probnorm`, `probt`, `put`, `putc`, `putn`, `qtr`, `quote`, `ranbin`, `rancau`, `ranexp`, `rangam`, `range`, `rank`, `rannor`, `ranpoi`, `rantbl`, `rantri`, `ranuni`, `repeat`, `resolve`, `reverse`, `rewind`, `right`, `round`, `saving`, `scan`, `sdf`, `second`, `sign`, `sin`, `sinh`, `skewness`, `soundex`, `spedis`, `sqrt`, `std`, `stderr`, `stfips`, `stname`, `stnamel`, `substr`, `sum`, `symget`, `sysget`, `sysmsg`, `sysprod`, `sysrc`, `system`, `tan`, `tanh`, `time`, `timepart`, `tinv`, `tnonct`, `today`, `translate`, `tranwrd`, `trigamma`, `trim`, `trimn`, `trunc`, `uniform`, `upcase`, `uss`, `var`, `varfmt`, `varinfmt`, `varlabel`, `varlen`, `varname`, `varnum`, `varray`, `varrayx`, `vartype`, `verify`, `vformat`, `vformatd`, `vformatdx`, `vformatn`, `vformatnx`, `vformatw`, `vformatwx`, `vformatx`, `vinarray`, `vinarrayx`, `vinformat`, `vinformatd`, `vinformatdx`, `vinformatn`, `vinformatnx`, `vinformatw`, `vinformatwx`, `vinformatx`, `vlabel`, `vlabelx`, `vlength`, `vlengthx`, `vname`, `vnamex`, `vtype`, `vtypex`, `weekday`, `year`, `yyq`, `zipfips`, `zipname`, `zipnamel`, `zipstate`), NameBuiltin, nil},
		},
		"vars-strings": {
			{`&[a-z_]\w{0,31}\.?`, NameVariable, nil},
			{`%[a-z_]\w{0,31}`, NameFunction, nil},
			{`\'`, LiteralString, Push("string_squote")},
			{`"`, LiteralString, Push("string_dquote")},
		},
		"string_squote": {
			{`'`, LiteralString, Pop(1)},
			{`\\\\|\\"|\\\n`, LiteralStringEscape, nil},
			{`[^$\'\\]+`, LiteralString, nil},
			{`[$\'\\]`, LiteralString, nil},
		},
		"string_dquote": {
			{`"`, LiteralString, Pop(1)},
			{`\\\\|\\"|\\\n`, LiteralStringEscape, nil},
			{`&`, NameVariable, Push("validvar")},
			{`[^$&"\\]+`, LiteralString, nil},
			{`[$"\\]`, LiteralString, nil},
		},
		"validvar": {
			{`[a-z_]\w{0,31}\.?`, NameVariable, Pop(1)},
		},
		"numbers": {
			{`\b[+-]?([0-9]+(\.[0-9]+)?|\.[0-9]+|\.)(E[+-]?[0-9]+)?i?\b`, LiteralNumber, nil},
		},
		"special": {
			{`(null|missing|_all_|_automatic_|_character_|_n_|_infile_|_name_|_null_|_numeric_|_user_|_webout_)`, KeywordConstant, nil},
		},
	}
}
