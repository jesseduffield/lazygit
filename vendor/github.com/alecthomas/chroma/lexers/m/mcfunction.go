package m

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// mcfunction lexer.
var MCFunction = internal.Register(MustNewLazyLexer(
	&Config{
		Name:         "mcfunction",
		Aliases:      []string{"mcfunction"},
		Filenames:    []string{"*.mcfunction"},
		MimeTypes:    []string{},
		NotMultiline: true,
		DotAll:       true,
	},
	func() Rules {
		return Rules{
			"simplevalue": {
				{`(true|false)`, KeywordConstant, nil},
				{`[01]b`, LiteralNumber, nil},
				{`-?(0|[1-9]\d*)(\.\d+[eE](\+|-)?\d+|[eE](\+|-)?\d+|\.\d+)`, LiteralNumberFloat, nil},
				{`(-?\d+)(\.\.)(-?\d+)`, ByGroups(LiteralNumberInteger, Punctuation, LiteralNumberInteger), nil},
				{`-?(0|[1-9]\d*)`, LiteralNumberInteger, nil},
				{`"(\\\\|\\"|[^"])*"`, LiteralStringDouble, nil},
				{`'[^']+'`, LiteralStringSingle, nil},
				{`([!#]?)(\w+)`, ByGroups(Punctuation, Text), nil},
			},
			"nbtobjectattribute": {
				Include("nbtvalue"),
				{`:`, Punctuation, nil},
				{`,`, Punctuation, Pop(1)},
				{`\}`, Punctuation, Pop(2)},
			},
			"nbtobjectvalue": {
				{`("(\\\\|\\"|[^"])*"|[a-zA-Z0-9_]+)`, NameTag, Push("nbtobjectattribute")},
				{`\}`, Punctuation, Pop(1)},
			},
			"nbtarrayvalue": {
				Include("nbtvalue"),
				{`,`, Punctuation, nil},
				{`\]`, Punctuation, Pop(1)},
			},
			"nbtvalue": {
				Include("simplevalue"),
				{`\{`, Punctuation, Push("nbtobjectvalue")},
				{`\[`, Punctuation, Push("nbtarrayvalue")},
			},
			"argumentvalue": {
				Include("simplevalue"),
				{`,`, Punctuation, Pop(1)},
				{`[}\]]`, Punctuation, Pop(2)},
			},
			"argumentlist": {
				{`(nbt)(={)`, ByGroups(NameAttribute, Punctuation), Push("nbtobjectvalue")},
				{`([A-Za-z0-9/_!]+)(={)`, ByGroups(NameAttribute, Punctuation), Push("argumentlist")},
				{`([A-Za-z0-9/_!]+)(=)`, ByGroups(NameAttribute, Punctuation), Push("argumentvalue")},
				Include("simplevalue"),
				{`,`, Punctuation, nil},
				{`[}\]]`, Punctuation, Pop(1)},
			},
			"root": {
				{`#.*?\n`, CommentSingle, nil},
				{Words(`/?`, `\b`, `ability`, `attributes`, `advancement`,
					`ban`, `ban-ip`, `banlist`, `bossbar`,
					`camerashake`, `classroommode`, `clear`,
					`clearspawnpoint`, `clone`, `code`, `collect`,
					`createagent`, `data`, `datapack`, `debug`,
					`defaultgamemode`, `deop`, `destroy`, `detect`,
					`detectredstone`, `difficulty`, `dropall`,
					`effect`, `enchant`, `event`, `execute`,
					`experience`, `fill`, `flog`, `forceload`,
					`function`, `gamemode`, `gamerule`,
					`geteduclientinfo`, `give`, `help`, `item`,
					`immutableworld`, `kick`, `kill`, `list`,
					`locate`, `locatebiome`, `loot`, `me`, `mixer`,
					`mobevent`, `move`, `msg`, `music`, `op`,
					`pardon`, `particle`, `playanimation`,
					`playsound`, `position`, `publish`,
					`raytracefog`, `recipe`, `reload`, `remove`,
					`replaceitem`, `ride`, `save`, `save-all`,
					`save-off`, `save-on`, `say`, `schedule`,
					`scoreboard`, `seed`, `setblock`,
					`setidletimeout`, `setmaxplayers`,
					`setworldspawn`, `spawnpoint`, `spectate`,
					`spreadplayers`, `stop`, `stopsound`,
					`structure`, `summon`, `tag`, `team`, `teammsg`,
					`teleport`, `tell`, `tellraw`, `testfor`,
					`testforblock`, `testforblocks`, `tickingarea`,
					`time`, `title`, `toggledownfall`, `tp`,
					`tpagent`, `transfer`, `transferserver`,
					`trigger`, `turn`, `w`, `weather`, `whitelist`,
					`worldborder`, `worldbuilder`, `wsserver`, `xp`,
				), KeywordReserved, nil},
				{Words(``, ``, `@p`, `@r`, `@a`, `@e`, `@s`, `@c`, `@v`),
					KeywordConstant, nil},
				{`\[`, Punctuation, Push("argumentlist")},
				{`{`, Punctuation, Push("nbtobjectvalue")},
				{`~`, NameBuiltin, nil},
				{`([a-zA-Z_]+:)?[a-zA-Z_]+\b`, Text, nil},
				{`([a-z]+)(\.)([0-9]+)\b`, ByGroups(Text, Punctuation, LiteralNumber), nil},
				{`([<>=]|<=|>=)`, Punctuation, nil},
				Include("simplevalue"),
				{`\s+`, TextWhitespace, nil},
			},
		}
	},
))
