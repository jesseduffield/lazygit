package v

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Verilog lexer.
var Verilog = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "verilog",
		Aliases:   []string{"verilog", "v"},
		Filenames: []string{"*.v"},
		MimeTypes: []string{"text/x-verilog"},
		EnsureNL:  true,
	},
	verilogRules,
))

func verilogRules() Rules {
	return Rules{
		"root": {
			{"^\\s*`define", CommentPreproc, Push("macro")},
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
			{`^(\s*)(package)(\s+)`, ByGroups(Text, KeywordNamespace, Text), nil},
			{`^(\s*)(import)(\s+)`, ByGroups(Text, KeywordNamespace, Text), Push("import")},
			{Words(``, `\b`, `always`, `always_comb`, `always_ff`, `always_latch`, `and`, `assign`, `automatic`, `begin`, `break`, `buf`, `bufif0`, `bufif1`, `case`, `casex`, `casez`, `cmos`, `const`, `continue`, `deassign`, `default`, `defparam`, `disable`, `do`, `edge`, `else`, `end`, `endcase`, `endfunction`, `endgenerate`, `endmodule`, `endpackage`, `endprimitive`, `endspecify`, `endtable`, `endtask`, `enum`, `event`, `final`, `for`, `force`, `forever`, `fork`, `function`, `generate`, `genvar`, `highz0`, `highz1`, `if`, `initial`, `inout`, `input`, `integer`, `join`, `large`, `localparam`, `macromodule`, `medium`, `module`, `nand`, `negedge`, `nmos`, `nor`, `not`, `notif0`, `notif1`, `or`, `output`, `packed`, `parameter`, `pmos`, `posedge`, `primitive`, `pull0`, `pull1`, `pulldown`, `pullup`, `rcmos`, `ref`, `release`, `repeat`, `return`, `rnmos`, `rpmos`, `rtran`, `rtranif0`, `rtranif1`, `scalared`, `signed`, `small`, `specify`, `specparam`, `strength`, `string`, `strong0`, `strong1`, `struct`, `table`, `task`, `tran`, `tranif0`, `tranif1`, `type`, `typedef`, `unsigned`, `var`, `vectored`, `void`, `wait`, `weak0`, `weak1`, `while`, `xnor`, `xor`), Keyword, nil},
			{Words("`", `\b`, `accelerate`, `autoexpand_vectornets`, `celldefine`, `default_nettype`, `else`, `elsif`, `endcelldefine`, `endif`, `endprotect`, `endprotected`, `expand_vectornets`, `ifdef`, `ifndef`, `include`, `noaccelerate`, `noexpand_vectornets`, `noremove_gatenames`, `noremove_netnames`, `nounconnected_drive`, `protect`, `protected`, `remove_gatenames`, `remove_netnames`, `resetall`, `timescale`, `unconnected_drive`, `undef`), CommentPreproc, nil},
			{Words(`\$`, `\b`, `bits`, `bitstoreal`, `bitstoshortreal`, `countdrivers`, `display`, `fclose`, `fdisplay`, `finish`, `floor`, `fmonitor`, `fopen`, `fstrobe`, `fwrite`, `getpattern`, `history`, `incsave`, `input`, `itor`, `key`, `list`, `log`, `monitor`, `monitoroff`, `monitoron`, `nokey`, `nolog`, `printtimescale`, `random`, `readmemb`, `readmemh`, `realtime`, `realtobits`, `reset`, `reset_count`, `reset_value`, `restart`, `rtoi`, `save`, `scale`, `scope`, `shortrealtobits`, `showscopes`, `showvariables`, `showvars`, `sreadmemb`, `sreadmemh`, `stime`, `stop`, `strobe`, `time`, `timeformat`, `write`), NameBuiltin, nil},
			{Words(``, `\b`, `byte`, `shortint`, `int`, `longint`, `integer`, `time`, `bit`, `logic`, `reg`, `supply0`, `supply1`, `tri`, `triand`, `trior`, `tri0`, `tri1`, `trireg`, `uwire`, `wire`, `wand`, `woshortreal`, `real`, `realtime`), KeywordType, nil},
			{`[a-zA-Z_]\w*:(?!:)`, NameLabel, nil},
			{`\$?[a-zA-Z_]\w*`, Name, nil},
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
