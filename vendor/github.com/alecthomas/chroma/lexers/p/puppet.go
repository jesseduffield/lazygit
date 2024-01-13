package p

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Puppet lexer.
var Puppet = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Puppet",
		Aliases:   []string{"puppet"},
		Filenames: []string{"*.pp"},
		MimeTypes: []string{},
	},
	puppetRules,
))

func puppetRules() Rules {
	return Rules{
		"root": {
			Include("comments"),
			Include("keywords"),
			Include("names"),
			Include("numbers"),
			Include("operators"),
			Include("strings"),
			{`[]{}:(),;[]`, Punctuation, nil},
			{`[^\S\n]+`, Text, nil},
		},
		"comments": {
			{`\s*#.*$`, Comment, nil},
			{`/(\\\n)?[*](.|\n)*?[*](\\\n)?/`, CommentMultiline, nil},
		},
		"operators": {
			{`(=>|\?|<|>|=|\+|-|/|\*|~|!|\|)`, Operator, nil},
			{`(in|and|or|not)\b`, OperatorWord, nil},
		},
		"names": {
			{`[a-zA-Z_]\w*`, NameAttribute, nil},
			{`(\$\S+)(\[)(\S+)(\])`, ByGroups(NameVariable, Punctuation, LiteralString, Punctuation), nil},
			{`\$\S+`, NameVariable, nil},
		},
		"numbers": {
			{`(\d+\.\d*|\d*\.\d+)([eE][+-]?[0-9]+)?j?`, LiteralNumberFloat, nil},
			{`\d+[eE][+-]?[0-9]+j?`, LiteralNumberFloat, nil},
			{`0[0-7]+j?`, LiteralNumberOct, nil},
			{`0[xX][a-fA-F0-9]+`, LiteralNumberHex, nil},
			{`\d+L`, LiteralNumberIntegerLong, nil},
			{`\d+j?`, LiteralNumberInteger, nil},
		},
		"keywords": {
			{Words(`(?i)`, `\b`, `absent`, `alert`, `alias`, `audit`, `augeas`, `before`, `case`, `check`, `class`, `computer`, `configured`, `contained`, `create_resources`, `crit`, `cron`, `debug`, `default`, `define`, `defined`, `directory`, `else`, `elsif`, `emerg`, `err`, `exec`, `extlookup`, `fail`, `false`, `file`, `filebucket`, `fqdn_rand`, `generate`, `host`, `if`, `import`, `include`, `info`, `inherits`, `inline_template`, `installed`, `interface`, `k5login`, `latest`, `link`, `loglevel`, `macauthorization`, `mailalias`, `maillist`, `mcx`, `md5`, `mount`, `mounted`, `nagios_command`, `nagios_contact`, `nagios_contactgroup`, `nagios_host`, `nagios_hostdependency`, `nagios_hostescalation`, `nagios_hostextinfo`, `nagios_hostgroup`, `nagios_service`, `nagios_servicedependency`, `nagios_serviceescalation`, `nagios_serviceextinfo`, `nagios_servicegroup`, `nagios_timeperiod`, `node`, `noop`, `notice`, `notify`, `package`, `present`, `purged`, `realize`, `regsubst`, `resources`, `role`, `router`, `running`, `schedule`, `scheduled_task`, `search`, `selboolean`, `selmodule`, `service`, `sha1`, `shellquote`, `split`, `sprintf`, `ssh_authorized_key`, `sshkey`, `stage`, `stopped`, `subscribe`, `tag`, `tagged`, `template`, `tidy`, `true`, `undef`, `unmounted`, `user`, `versioncmp`, `vlan`, `warning`, `yumrepo`, `zfs`, `zone`, `zpool`), Keyword, nil},
		},
		"strings": {
			{`"([^"])*"`, LiteralString, nil},
			{`'(\\'|[^'])*'`, LiteralString, nil},
		},
	}
}
