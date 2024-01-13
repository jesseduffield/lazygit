package s

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Systemverilog lexer.
var Systemverilog = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "systemverilog",
		Aliases:   []string{"systemverilog", "sv"},
		Filenames: []string{"*.sv", "*.svh"},
		MimeTypes: []string{"text/x-systemverilog"},
		EnsureNL:  true,
	},
	systemvarilogRules,
))

func systemvarilogRules() Rules {
	return Rules{
		"root": {
			{"^\\s*`define", CommentPreproc, Push("macro")},
			{`^(\s*)(package)(\s+)`, ByGroups(Text, KeywordNamespace, Text), nil},
			{`^(\s*)(import)(\s+)("DPI(?:-C)?")(\s+)`, ByGroups(Text, KeywordNamespace, Text, LiteralString, Text), nil},
			{`^(\s*)(import)(\s+)`, ByGroups(Text, KeywordNamespace, Text), Push("import")},
			{`\n`, Text, nil},
			{`\s+`, Text, nil},
			{`\\\n`, Text, nil},
			{`/(\\\n)?/(\n|(.|\n)*?[^\\]\n)`, CommentSingle, nil},
			{`/(\\\n)?[*](.|\n)*?[*](\\\n)?/`, CommentMultiline, nil},
			{`[{}#@]`, Punctuation, nil},
			{`L?"`, LiteralString, Push("string")},
			{`L?'(\\.|\\[0-7]{1,3}|\\x[a-fA-F0-9]{1,2}|[^\\\'\n])'`, LiteralStringChar, nil},
			{`(\d+\.\d*|\.\d+|\d+)[eE][+-]?\d+[lL]?`, LiteralNumberFloat, nil},
			{`(\d+\.\d*|\.\d+|\d+[fF])[fF]?`, LiteralNumberFloat, nil},
			{`([0-9]+)|(\'h)[0-9a-fA-F]+`, LiteralNumberHex, nil},
			{`([0-9]+)|(\'b)[01]+`, LiteralNumberBin, nil},
			{`([0-9]+)|(\'d)[0-9]+`, LiteralNumberInteger, nil},
			{`([0-9]+)|(\'o)[0-7]+`, LiteralNumberOct, nil},
			{`\'[01xz]`, LiteralNumber, nil},
			{`\d+[Ll]?`, LiteralNumberInteger, nil},
			{`\*/`, Error, nil},
			{`[~!%^&*+=|?:<>/-]`, Operator, nil},
			{`[()\[\],.;\']`, Punctuation, nil},
			{"`[a-zA-Z_]\\w*", NameConstant, nil},
			{Words(``, `\b`, `accept_on`, `alias`, `always`, `always_comb`, `always_ff`, `always_latch`, `and`, `assert`, `assign`, `assume`, `automatic`, `before`, `begin`, `bind`, `bins`, `binsof`, `bit`, `break`, `buf`, `bufif0`, `bufif1`, `byte`, `case`, `casex`, `casez`, `cell`, `chandle`, `checker`, `class`, `clocking`, `cmos`, `config`, `const`, `constraint`, `context`, `continue`, `cover`, `covergroup`, `coverpoint`, `cross`, `deassign`, `default`, `defparam`, `design`, `disable`, `dist`, `do`, `edge`, `else`, `end`, `endcase`, `endchecker`, `endclass`, `endclocking`, `endconfig`, `endfunction`, `endgenerate`, `endgroup`, `endinterface`, `endmodule`, `endpackage`, `endprimitive`, `endprogram`, `endproperty`, `endsequence`, `endspecify`, `endtable`, `endtask`, `enum`, `event`, `eventually`, `expect`, `export`, `extends`, `extern`, `final`, `first_match`, `for`, `force`, `foreach`, `forever`, `fork`, `forkjoin`, `function`, `generate`, `genvar`, `global`, `highz0`, `highz1`, `if`, `iff`, `ifnone`, `ignore_bins`, `illegal_bins`, `implies`, `import`, `incdir`, `include`, `initial`, `inout`, `input`, `inside`, `instance`, `int`, `integer`, `interface`, `intersect`, `join`, `join_any`, `join_none`, `large`, `let`, `liblist`, `library`, `local`, `localparam`, `logic`, `longint`, `macromodule`, `matches`, `medium`, `modport`, `module`, `nand`, `negedge`, `new`, `nexttime`, `nmos`, `nor`, `noshowcancelled`, `not`, `notif0`, `notif1`, `null`, `or`, `output`, `package`, `packed`, `parameter`, `pmos`, `posedge`, `primitive`, `priority`, `program`, `property`, `protected`, `pull0`, `pull1`, `pulldown`, `pullup`, `pulsestyle_ondetect`, `pulsestyle_onevent`, `pure`, `rand`, `randc`, `randcase`, `randsequence`, `rcmos`, `real`, `realtime`, `ref`, `reg`, `reject_on`, `release`, `repeat`, `restrict`, `return`, `rnmos`, `rpmos`, `rtran`, `rtranif0`, `rtranif1`, `s_always`, `s_eventually`, `s_nexttime`, `s_until`, `s_until_with`, `scalared`, `sequence`, `shortint`, `shortreal`, `showcancelled`, `signed`, `small`, `solve`, `specify`, `specparam`, `static`, `string`, `strong`, `strong0`, `strong1`, `struct`, `super`, `supply0`, `supply1`, `sync_accept_on`, `sync_reject_on`, `table`, `tagged`, `task`, `this`, `throughout`, `time`, `timeprecision`, `timeunit`, `tran`, `tranif0`, `tranif1`, `tri`, `tri0`, `tri1`, `triand`, `trior`, `trireg`, `type`, `typedef`, `union`, `unique`, `unique0`, `unsigned`, `until`, `until_with`, `untyped`, `use`, `uwire`, `var`, `vectored`, `virtual`, `void`, `wait`, `wait_order`, `wand`, `weak`, `weak0`, `weak1`, `while`, `wildcard`, `wire`, `with`, `within`, `wor`, `xnor`, `xor`), Keyword, nil},
			{Words(``, `\b`, "`__FILE__", "`__LINE__", "`begin_keywords", "`celldefine", "`default_nettype", "`define", "`else", "`elsif", "`end_keywords", "`endcelldefine", "`endif", "`ifdef", "`ifndef", "`include", "`line", "`nounconnected_drive", "`pragma", "`resetall", "`timescale", "`unconnected_drive", "`undef", "`undefineall"), CommentPreproc, nil},
			{Words(``, `\b`, `$display`, `$displayb`, `$displayh`, `$displayo`, `$dumpall`, `$dumpfile`, `$dumpflush`, `$dumplimit`, `$dumpoff`, `$dumpon`, `$dumpports`, `$dumpportsall`, `$dumpportsflush`, `$dumpportslimit`, `$dumpportsoff`, `$dumpportson`, `$dumpvars`, `$fclose`, `$fdisplay`, `$fdisplayb`, `$fdisplayh`, `$fdisplayo`, `$feof`, `$ferror`, `$fflush`, `$fgetc`, `$fgets`, `$finish`, `$fmonitor`, `$fmonitorb`, `$fmonitorh`, `$fmonitoro`, `$fopen`, `$fread`, `$fscanf`, `$fseek`, `$fstrobe`, `$fstrobeb`, `$fstrobeh`, `$fstrobeo`, `$ftell`, `$fwrite`, `$fwriteb`, `$fwriteh`, `$fwriteo`, `$monitor`, `$monitorb`, `$monitorh`, `$monitoro`, `$monitoroff`, `$monitoron`, `$plusargs`, `$random`, `$readmemb`, `$readmemh`, `$rewind`, `$sformat`, `$sformatf`, `$sscanf`, `$strobe`, `$strobeb`, `$strobeh`, `$strobeo`, `$swrite`, `$swriteb`, `$swriteh`, `$swriteo`, `$test`, `$ungetc`, `$value$plusargs`, `$write`, `$writeb`, `$writeh`, `$writememb`, `$writememh`, `$writeo`), NameBuiltin, nil},
			{`(class)(\s+)`, ByGroups(Keyword, Text), Push("classname")},
			{Words(``, `\b`, `byte`, `shortint`, `int`, `longint`, `integer`, `time`, `bit`, `logic`, `reg`, `supply0`, `supply1`, `tri`, `triand`, `trior`, `tri0`, `tri1`, `trireg`, `uwire`, `wire`, `wand`, `woshortreal`, `real`, `realtime`), KeywordType, nil},
			{`[a-zA-Z_]\w*:(?!:)`, NameLabel, nil},
			{`\$?[a-zA-Z_]\w*`, Name, nil},
		},
		"classname": {
			{`[a-zA-Z_]\w*`, NameClass, Pop(1)},
		},
		"string": {
			{`"`, LiteralString, Pop(1)},
			{`\\([\\abfnrtv"\']|x[a-fA-F0-9]{2,4}|[0-7]{1,3})`, LiteralStringEscape, nil},
			{`[^\\"\n]+`, LiteralString, nil},
			{`\\\n`, LiteralString, nil},
			{`\\`, LiteralString, nil},
		},
		"macro": {
			{`[^/\n]+`, CommentPreproc, nil},
			{`/[*](.|\n)*?[*]/`, CommentMultiline, nil},
			{`//.*?\n`, CommentSingle, Pop(1)},
			{`/`, CommentPreproc, nil},
			{`(?<=\\)\n`, CommentPreproc, nil},
			{`\n`, CommentPreproc, Pop(1)},
		},
		"import": {
			{`[\w:]+\*?`, NameNamespace, Pop(1)},
		},
	}
}
