package m

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Meson lexer.
var Meson = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Meson",
		Aliases:   []string{"meson", "meson.build"},
		Filenames: []string{"meson.build", "meson_options.txt"},
		MimeTypes: []string{"text/x-meson"},
	},
	func() Rules {
		return Rules{
			"root": {
				{`#.*?$`, Comment, nil},
				{`'''.*'''`, LiteralStringSingle, nil},
				{`[1-9][0-9]*`, LiteralNumberInteger, nil},
				{`0o[0-7]+`, LiteralNumberOct, nil},
				{`0x[a-fA-F0-9]+`, LiteralNumberHex, nil},
				Include("string"),
				Include("keywords"),
				Include("expr"),
				{`[a-zA-Z_][a-zA-Z_0-9]*`, Name, nil},
				{`\s+`, TextWhitespace, nil},
			},
			"string": {
				{`[']{3}([']{0,2}([^\\']|\\(.|\n)))*[']{3}`, LiteralString, nil},
				{`'.*?(?<!\\)(\\\\)*?'`, LiteralString, nil},
			},
			"keywords": {
				{Words(``, `\b`, `if`, `elif`, `else`, `endif`, `foreach`, `endforeach`, `break`, `continue`), Keyword, nil},
			},
			"expr": {
				{`(in|and|or|not)\b`, OperatorWord, nil},
				{`(\*=|/=|%=|\+]=|-=|==|!=|\+|-|=)`, Operator, nil},
				{`[\[\]{}:().,?]`, Punctuation, nil},
				{Words(``, `\b`, `true`, `false`), KeywordConstant, nil},
				Include("builtins"),
				{Words(``, `\b`, `meson`, `build_machine`, `host_machine`, `target_machine`), NameVariableMagic, nil},
			},
			"builtins": {
				{Words(`(?<!\.)`, `\b`, `add_global_arguments`, `add_global_link_arguments`, `add_languages`, `add_project_arguments`, `add_project_link_arguments`, `add_test_setup`, `assert`, `benchmark`, `both_libraries`, `build_target`, `configuration_data`, `configure_file`, `custom_target`, `declare_dependency`, `dependency`, `disabler`, `environment`, `error`, `executable`, `files`, `find_library`, `find_program`, `generator`, `get_option`, `get_variable`, `include_directories`, `install_data`, `install_headers`, `install_man`, `install_subdir`, `is_disabler`, `is_variable`, `jar`, `join_paths`, `library`, `message`, `project`, `range`, `run_command`, `set_variable`, `shared_library`, `shared_module`, `static_library`, `subdir`, `subdir_done`, `subproject`, `summary`, `test`, `vcs_tag`, `warning`), NameBuiltin, nil},
				{`(?<!\.)import\b`, NameNamespace, nil},
			},
		}
	},
))
