package f

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Fennel lexer.
var Fennel = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Fennel",
		Aliases:   []string{"fennel", "fnl"},
		Filenames: []string{"*.fennel"},
		MimeTypes: []string{"text/x-fennel", "application/x-fennel"},
	},
	fennelRules,
))

// Here's some Fennel code used to generate the lists of keywords:
// (local fennel (require :fennel))
//
// (fn member? [t x] (each [_ y (ipairs t)] (when (= y x) (lua "return true"))))
//
// (local declarations [:fn :lambda :λ :local :var :global :macro :macros])
// (local keywords [])
// (local globals [])
//
// (each [name data (pairs (fennel.syntax))]
//   (if (member? declarations name) nil ; already populated
//       data.special? (table.insert keywords name)
//       data.macro? (table.insert keywords name)
//       data.global? (table.insert globals name)))
//
// (fn quoted [tbl]
//   (table.sort tbl)
//   (table.concat (icollect [_ k (ipairs tbl)]
//                   (string.format "`%s`" k)) ", "))
//
// (print :Keyword (quoted keywords))
// (print :KeywordDeclaration (quoted declarations))
// (print :NameBuiltin (quoted globals))

func fennelRules() Rules {
	return Rules{
		"root": {
			{`;.*$`, CommentSingle, nil},
			{`\s+`, Whitespace, nil},
			{`-?\d+\.\d+`, LiteralNumberFloat, nil},
			{`-?\d+`, LiteralNumberInteger, nil},
			{`0x-?[abcdef\d]+`, LiteralNumberHex, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralString, nil},
			{`'(?!#)[\w!$%*+<=>?/.#-]+`, LiteralStringSymbol, nil},
			{`\\(.|[a-z]+)`, LiteralStringChar, nil},
			{`::?#?(?!#)[\w!$%*+<=>?/.#-]+`, LiteralStringSymbol, nil},
			{"~@|[`\\'#^~&@]", Operator, nil},
			{Words(``, ` `, `#`, `%`, `*`, `+`, `-`, `->`, `->>`, `-?>`, `-?>>`, `.`, `..`, `/`, `//`, `:`, `<`, `<=`, `=`, `>`, `>=`, `?.`, `^`, `accumulate`, `and`, `band`, `bnot`, `bor`, `bxor`, `collect`, `comment`, `do`, `doc`, `doto`, `each`, `eval-compiler`, `for`, `hashfn`, `icollect`, `if`, `import-macros`, `include`, `length`, `let`, `lshift`, `lua`, `macrodebug`, `match`, `not`, `not=`, `or`, `partial`, `pick-args`, `pick-values`, `quote`, `require-macros`, `rshift`, `set`, `set-forcibly!`, `tset`, `values`, `when`, `while`, `with-open`, `~=`), Keyword, nil},
			{Words(``, ` `, `fn`, `global`, `lambda`, `local`, `macro`, `macros`, `var`, `λ`), KeywordDeclaration, nil},
			{Words(``, ` `, `_G`, `arg`, `assert`, `bit32`, `bit32.arshift`, `bit32.band`, `bit32.bnot`, `bit32.bor`, `bit32.btest`, `bit32.bxor`, `bit32.extract`, `bit32.lrotate`, `bit32.lshift`, `bit32.replace`, `bit32.rrotate`, `bit32.rshift`, `collectgarbage`, `coroutine`, `coroutine.create`, `coroutine.resume`, `coroutine.running`, `coroutine.status`, `coroutine.wrap`, `coroutine.yield`, `debug`, `debug.debug`, `debug.gethook`, `debug.getinfo`, `debug.getlocal`, `debug.getmetatable`, `debug.getregistry`, `debug.getupvalue`, `debug.getuservalue`, `debug.sethook`, `debug.setlocal`, `debug.setmetatable`, `debug.setupvalue`, `debug.setuservalue`, `debug.traceback`, `debug.upvalueid`, `debug.upvaluejoin`, `dofile`, `error`, `getmetatable`, `io`, `io.close`, `io.flush`, `io.input`, `io.lines`, `io.open`, `io.output`, `io.popen`, `io.read`, `io.tmpfile`, `io.type`, `io.write`, `ipairs`, `load`, `loadfile`, `loadstring`, `math`, `math.abs`, `math.acos`, `math.asin`, `math.atan`, `math.atan2`, `math.ceil`, `math.cos`, `math.cosh`, `math.deg`, `math.exp`, `math.floor`, `math.fmod`, `math.frexp`, `math.ldexp`, `math.log`, `math.log10`, `math.max`, `math.min`, `math.modf`, `math.pow`, `math.rad`, `math.random`, `math.randomseed`, `math.sin`, `math.sinh`, `math.sqrt`, `math.tan`, `math.tanh`, `module`, `next`, `os`, `os.clock`, `os.date`, `os.difftime`, `os.execute`, `os.exit`, `os.getenv`, `os.remove`, `os.rename`, `os.setlocale`, `os.time`, `os.tmpname`, `package`, `package.loadlib`, `package.searchpath`, `package.seeall`, `pairs`, `pcall`, `print`, `rawequal`, `rawget`, `rawlen`, `rawset`, `require`, `select`, `setmetatable`, `string`, `string.byte`, `string.char`, `string.dump`, `string.find`, `string.format`, `string.gmatch`, `string.gsub`, `string.len`, `string.lower`, `string.match`, `string.rep`, `string.reverse`, `string.sub`, `string.upper`, `table`, `table.concat`, `table.insert`, `table.maxn`, `table.pack`, `table.remove`, `table.sort`, `table.unpack`, `tonumber`, `tostring`, `type`, `unpack`, `xpcall`), NameBuiltin, nil},
			{`(?<=\()(?!#)[\w!$%*+<=>?/.#-]+`, NameFunction, nil},
			{`(?!#)[\w!$%*+<=>?/.#-]+`, NameVariable, nil},
			{`(\[|\])`, Punctuation, nil},
			{`(\{|\})`, Punctuation, nil},
			{`(\(|\))`, Punctuation, nil},
		},
	}
}
