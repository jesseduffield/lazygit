package r

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Rust lexer.
var Rust = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Rust",
		Aliases:   []string{"rust", "rs"},
		Filenames: []string{"*.rs", "*.rs.in"},
		MimeTypes: []string{"text/rust", "text/x-rust"},
		EnsureNL:  true,
	},
	rustRules,
))

func rustRules() Rules {
	return Rules{
		"root": {
			{`#![^[\r\n].*$`, CommentPreproc, nil},
			Default(Push("base")),
		},
		"base": {
			{`\n`, TextWhitespace, nil},
			{`\s+`, TextWhitespace, nil},
			{`//!.*?\n`, LiteralStringDoc, nil},
			{`///(\n|[^/].*?\n)`, LiteralStringDoc, nil},
			{`//(.*?)\n`, CommentSingle, nil},
			{`/\*\*(\n|[^/*])`, LiteralStringDoc, Push("doccomment")},
			{`/\*!`, LiteralStringDoc, Push("doccomment")},
			{`/\*`, CommentMultiline, Push("comment")},
			{`r#*"(?:\\.|[^\\;])*"#*`, LiteralString, nil},
			{`"(?:\\.|[^\\"])*"`, LiteralString, nil},
			{`\$([a-zA-Z_]\w*|\(,?|\),?|,?)`, CommentPreproc, nil},
			{Words(``, `\b`, `as`, `async`, `await`, `box`, `const`, `crate`, `dyn`, `else`, `extern`, `for`, `if`, `impl`, `in`, `loop`, `match`, `move`, `mut`, `pub`, `ref`, `return`, `static`, `super`, `trait`, `unsafe`, `use`, `where`, `while`), Keyword, nil},
			{Words(``, `\b`, `abstract`, `become`, `do`, `final`, `macro`, `override`, `priv`, `typeof`, `try`, `unsized`, `virtual`, `yield`), KeywordReserved, nil},
			{`(true|false)\b`, KeywordConstant, nil},
			{`self\b`, NameBuiltinPseudo, nil},
			{`mod\b`, Keyword, Push("modname")},
			{`let\b`, KeywordDeclaration, nil},
			{`fn\b`, Keyword, Push("funcname")},
			{`(struct|enum|type|union)\b`, Keyword, Push("typename")},
			{`(default)(\s+)(type|fn)\b`, ByGroups(Keyword, Text, Keyword), nil},
			{Words(``, `\b`, `u8`, `u16`, `u32`, `u64`, `u128`, `i8`, `i16`, `i32`, `i64`, `i128`, `usize`, `isize`, `f32`, `f64`, `char`, `str`, `bool`), KeywordType, nil},
			{`[sS]elf\b`, NameBuiltinPseudo, nil},
			{Words(``, `\b`, `Copy`, `Send`, `Sized`, `Sync`, `Unpin`, `Drop`, `Fn`, `FnMut`, `FnOnce`, `drop`, `Box`, `ToOwned`, `Clone`, `PartialEq`, `PartialOrd`, `Eq`, `Ord`, `AsRef`, `AsMut`, `Into`, `From`, `Default`, `Iterator`, `Extend`, `IntoIterator`, `DoubleEndedIterator`, `ExactSizeIterator`, `Option`, `Some`, `None`, `Result`, `Ok`, `Err`, `String`, `ToString`, `Vec`), NameBuiltin, nil},
			{Words(``, `!`, `asm`, `assert`, `assert_eq`, `assert_ne`, `cfg`, `column`, `compile_error`, `concat`, `concat_idents`, `dbg`, `debug_assert`, `debug_assert_eq`, `debug_assert_ne`, `env`, `eprint`, `eprintln`, `file`, `format`, `format_args`, `format_args_nl`, `global_asm`, `include`, `include_bytes`, `include_str`, `is_aarch64_feature_detected`, `is_arm_feature_detected`, `is_mips64_feature_detected`, `is_mips_feature_detected`, `is_powerpc64_feature_detected`, `is_powerpc_feature_detected`, `is_x86_feature_detected`, `line`, `llvm_asm`, `log_syntax`, `macro_rules`, `matches`, `module_path`, `option_env`, `panic`, `print`, `println`, `stringify`, `thread_local`, `todo`, `trace_macros`, `unimplemented`, `unreachable`, `vec`, `write`, `writeln`), NameFunctionMagic, nil},
			{`::\b`, Text, nil},
			{`(?::|->)`, Text, Push("typename")},
			{`(break|continue)(\b\s*)(\'[A-Za-z_]\w*)?`, ByGroups(Keyword, TextWhitespace, NameLabel), nil},
			{`'(\\['"\\nrt]|\\x[0-7][0-9a-fA-F]|\\0|\\u\{[0-9a-fA-F]{1,6}\}|.)'`, LiteralStringChar, nil},
			{`b'(\\['"\\nrt]|\\x[0-9a-fA-F]{2}|\\0|\\u\{[0-9a-fA-F]{1,6}\}|.)'`, LiteralStringChar, nil},
			{`0b[01_]+`, LiteralNumberBin, Push("number_lit")},
			{`0o[0-7_]+`, LiteralNumberOct, Push("number_lit")},
			{`0[xX][0-9a-fA-F_]+`, LiteralNumberHex, Push("number_lit")},
			{`[0-9][0-9_]*(\.[0-9_]+[eE][+\-]?[0-9_]+|\.[0-9_]*(?!\.)|[eE][+\-]?[0-9_]+)`, LiteralNumberFloat, Push("number_lit")},
			{`[0-9][0-9_]*`, LiteralNumberInteger, Push("number_lit")},
			{`b"`, LiteralString, Push("bytestring")},
			{`(?s)b?r(#*)".*?"\1`, LiteralString, nil},
			{`'`, Operator, Push("lifetime")},
			{`\.\.=?`, Operator, nil},
			{`[{}()\[\],.;]`, Punctuation, nil},
			{`[+\-*/%&|<>^!~@=:?]`, Operator, nil},
			{`(r#)?[a-zA-Z_]\w*`, Name, nil},
			{`r#[a-zA-Z_]\w*`, Name, nil},
			{`#!?\[`, CommentPreproc, Push("attribute[")},
			{`#`, Text, nil},
		},
		"comment": {
			{`[^*/]+`, CommentMultiline, nil},
			{`/\*`, CommentMultiline, Push()},
			{`\*/`, CommentMultiline, Pop(1)},
			{`[*/]`, CommentMultiline, nil},
		},
		"doccomment": {
			{`[^*/]+`, LiteralStringDoc, nil},
			{`/\*`, LiteralStringDoc, Push()},
			{`\*/`, LiteralStringDoc, Pop(1)},
			{`[*/]`, LiteralStringDoc, nil},
		},
		"modname": {
			{`\s+`, Text, nil},
			{`[a-zA-Z_]\w*`, NameNamespace, Pop(1)},
			Default(Pop(1)),
		},
		"funcname": {
			{`\s+`, Text, nil},
			{`[a-zA-Z_]\w*`, NameFunction, Pop(1)},
			Default(Pop(1)),
		},
		"typename": {
			{`\s+`, Text, nil},
			{`&`, KeywordPseudo, nil},
			{`'`, Operator, Push("lifetime")},
			{Words(``, `\b`, `Copy`, `Send`, `Sized`, `Sync`, `Unpin`, `Drop`, `Fn`, `FnMut`, `FnOnce`, `drop`, `Box`, `ToOwned`, `Clone`, `PartialEq`, `PartialOrd`, `Eq`, `Ord`, `AsRef`, `AsMut`, `Into`, `From`, `Default`, `Iterator`, `Extend`, `IntoIterator`, `DoubleEndedIterator`, `ExactSizeIterator`, `Option`, `Some`, `None`, `Result`, `Ok`, `Err`, `String`, `ToString`, `Vec`), NameBuiltin, nil},
			{Words(``, `\b`, `u8`, `u16`, `u32`, `u64`, `u128`, `i8`, `i16`, `i32`, `i64`, `i128`, `usize`, `isize`, `f32`, `f64`, `char`, `str`, `bool`), KeywordType, nil},
			{`[a-zA-Z_]\w*`, NameClass, Pop(1)},
			Default(Pop(1)),
		},
		"lifetime": {
			{`(static|_)`, NameBuiltin, nil},
			{`[a-zA-Z_]+\w*`, NameAttribute, nil},
			Default(Pop(1)),
		},
		"number_lit": {
			{`[ui](8|16|32|64|size)`, Keyword, Pop(1)},
			{`f(32|64)`, Keyword, Pop(1)},
			Default(Pop(1)),
		},
		"string": {
			{`"`, LiteralString, Pop(1)},
			{`\\['"\\nrt]|\\x[0-7][0-9a-fA-F]|\\0|\\u\{[0-9a-fA-F]{1,6}\}`, LiteralStringEscape, nil},
			{`[^\\"]+`, LiteralString, nil},
			{`\\`, LiteralString, nil},
		},
		"bytestring": {
			{`\\x[89a-fA-F][0-9a-fA-F]`, LiteralStringEscape, nil},
			Include("string"),
		},
		"attribute_common": {
			{`"`, LiteralString, Push("string")},
			{`\[`, CommentPreproc, Push("attribute[")},
		},
		"attribute[": {
			Include("attribute_common"),
			{`\]`, CommentPreproc, Pop(1)},
			{`[^"\]\[]+`, CommentPreproc, nil},
		},
	}
}
