package e

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Erlang lexer.
var Erlang = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Erlang",
		Aliases:   []string{"erlang"},
		Filenames: []string{"*.erl", "*.hrl", "*.es", "*.escript"},
		MimeTypes: []string{"text/x-erlang"},
	},
	erlangRules,
))

func erlangRules() Rules {
	return Rules{
		"root": {
			{`\s+`, Text, nil},
			{`%.*\n`, Comment, nil},
			{Words(``, `\b`, `after`, `begin`, `case`, `catch`, `cond`, `end`, `fun`, `if`, `let`, `of`, `query`, `receive`, `try`, `when`), Keyword, nil},
			{Words(``, `\b`, `abs`, `append_element`, `apply`, `atom_to_list`, `binary_to_list`, `bitstring_to_list`, `binary_to_term`, `bit_size`, `bump_reductions`, `byte_size`, `cancel_timer`, `check_process_code`, `delete_module`, `demonitor`, `disconnect_node`, `display`, `element`, `erase`, `exit`, `float`, `float_to_list`, `fun_info`, `fun_to_list`, `function_exported`, `garbage_collect`, `get`, `get_keys`, `group_leader`, `hash`, `hd`, `integer_to_list`, `iolist_to_binary`, `iolist_size`, `is_atom`, `is_binary`, `is_bitstring`, `is_boolean`, `is_builtin`, `is_float`, `is_function`, `is_integer`, `is_list`, `is_number`, `is_pid`, `is_port`, `is_process_alive`, `is_record`, `is_reference`, `is_tuple`, `length`, `link`, `list_to_atom`, `list_to_binary`, `list_to_bitstring`, `list_to_existing_atom`, `list_to_float`, `list_to_integer`, `list_to_pid`, `list_to_tuple`, `load_module`, `localtime_to_universaltime`, `make_tuple`, `md5`, `md5_final`, `md5_update`, `memory`, `module_loaded`, `monitor`, `monitor_node`, `node`, `nodes`, `open_port`, `phash`, `phash2`, `pid_to_list`, `port_close`, `port_command`, `port_connect`, `port_control`, `port_call`, `port_info`, `port_to_list`, `process_display`, `process_flag`, `process_info`, `purge_module`, `put`, `read_timer`, `ref_to_list`, `register`, `resume_process`, `round`, `send`, `send_after`, `send_nosuspend`, `set_cookie`, `setelement`, `size`, `spawn`, `spawn_link`, `spawn_monitor`, `spawn_opt`, `split_binary`, `start_timer`, `statistics`, `suspend_process`, `system_flag`, `system_info`, `system_monitor`, `system_profile`, `term_to_binary`, `tl`, `trace`, `trace_delivered`, `trace_info`, `trace_pattern`, `trunc`, `tuple_size`, `tuple_to_list`, `universaltime_to_localtime`, `unlink`, `unregister`, `whereis`), NameBuiltin, nil},
			{Words(``, `\b`, `and`, `andalso`, `band`, `bnot`, `bor`, `bsl`, `bsr`, `bxor`, `div`, `not`, `or`, `orelse`, `rem`, `xor`), OperatorWord, nil},
			{`^-`, Punctuation, Push("directive")},
			{`(\+\+?|--?|\*|/|<|>|/=|=:=|=/=|=<|>=|==?|<-|!|\?)`, Operator, nil},
			{`"`, LiteralString, Push("string")},
			{`<<`, NameLabel, nil},
			{`>>`, NameLabel, nil},
			{`((?:[a-z]\w*|'[^\n']*[^\\]'))(:)`, ByGroups(NameNamespace, Punctuation), nil},
			{`(?:^|(?<=:))((?:[a-z]\w*|'[^\n']*[^\\]'))(\s*)(\()`, ByGroups(NameFunction, Text, Punctuation), nil},
			{`[+-]?(?:[2-9]|[12][0-9]|3[0-6])#[0-9a-zA-Z]+`, LiteralNumberInteger, nil},
			{`[+-]?\d+`, LiteralNumberInteger, nil},
			{`[+-]?\d+.\d+`, LiteralNumberFloat, nil},
			{`[]\[:_@\".{}()|;,]`, Punctuation, nil},
			{`(?:[A-Z_]\w*)`, NameVariable, nil},
			{`(?:[a-z]\w*|'[^\n']*[^\\]')`, Name, nil},
			{`\?(?:(?:[A-Z_]\w*)|(?:[a-z]\w*|'[^\n']*[^\\]'))`, NameConstant, nil},
			{`\$(?:(?:\\(?:[bdefnrstv\'"\\]|[0-7][0-7]?[0-7]?|(?:x[0-9a-fA-F]{2}|x\{[0-9a-fA-F]+\})|\^[a-zA-Z]))|\\[ %]|[^\\])`, LiteralStringChar, nil},
			{`#(?:[a-z]\w*|'[^\n']*[^\\]')(:?\.(?:[a-z]\w*|'[^\n']*[^\\]'))?`, NameLabel, nil},
			{`\A#!.+\n`, CommentHashbang, nil},
			{`#\{`, Punctuation, Push("map_key")},
		},
		"string": {
			{`(?:\\(?:[bdefnrstv\'"\\]|[0-7][0-7]?[0-7]?|(?:x[0-9a-fA-F]{2}|x\{[0-9a-fA-F]+\})|\^[a-zA-Z]))`, LiteralStringEscape, nil},
			{`"`, LiteralString, Pop(1)},
			{`~[0-9.*]*[~#+BPWXb-ginpswx]`, LiteralStringInterpol, nil},
			{`[^"\\~]+`, LiteralString, nil},
			{`~`, LiteralString, nil},
		},
		"directive": {
			{`(define)(\s*)(\()((?:(?:[A-Z_]\w*)|(?:[a-z]\w*|'[^\n']*[^\\]')))`, ByGroups(NameEntity, Text, Punctuation, NameConstant), Pop(1)},
			{`(record)(\s*)(\()((?:(?:[A-Z_]\w*)|(?:[a-z]\w*|'[^\n']*[^\\]')))`, ByGroups(NameEntity, Text, Punctuation, NameLabel), Pop(1)},
			{`(?:[a-z]\w*|'[^\n']*[^\\]')`, NameEntity, Pop(1)},
		},
		"map_key": {
			Include("root"),
			{`=>`, Punctuation, Push("map_val")},
			{`:=`, Punctuation, Push("map_val")},
			{`\}`, Punctuation, Pop(1)},
		},
		"map_val": {
			Include("root"),
			{`,`, Punctuation, Pop(1)},
			{`(?=\})`, Punctuation, Pop(1)},
		},
	}
}
