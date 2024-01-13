package p

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Postscript lexer.
var Postscript = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "PostScript",
		Aliases:   []string{"postscript", "postscr"},
		Filenames: []string{"*.ps", "*.eps"},
		MimeTypes: []string{"application/postscript"},
	},
	postscriptRules,
))

func postscriptRules() Rules {
	return Rules{
		"root": {
			{`^%!.+\n`, CommentPreproc, nil},
			{`%%.*\n`, CommentSpecial, nil},
			{`(^%.*\n){2,}`, CommentMultiline, nil},
			{`%.*\n`, CommentSingle, nil},
			{`\(`, LiteralString, Push("stringliteral")},
			{`[{}<>\[\]]`, Punctuation, nil},
			{`<[0-9A-Fa-f]+>(?=[()<>\[\]{}/%\s])`, LiteralNumberHex, nil},
			{`[0-9]+\#(\-|\+)?([0-9]+\.?|[0-9]*\.[0-9]+|[0-9]+\.[0-9]*)((e|E)[0-9]+)?(?=[()<>\[\]{}/%\s])`, LiteralNumberOct, nil},
			{`(\-|\+)?([0-9]+\.?|[0-9]*\.[0-9]+|[0-9]+\.[0-9]*)((e|E)[0-9]+)?(?=[()<>\[\]{}/%\s])`, LiteralNumberFloat, nil},
			{`(\-|\+)?[0-9]+(?=[()<>\[\]{}/%\s])`, LiteralNumberInteger, nil},
			{`\/[^()<>\[\]{}/%\s]+(?=[()<>\[\]{}/%\s])`, NameVariable, nil},
			{`[^()<>\[\]{}/%\s]+(?=[()<>\[\]{}/%\s])`, NameFunction, nil},
			{`(false|true)(?=[()<>\[\]{}/%\s])`, KeywordConstant, nil},
			{`(eq|ne|g[et]|l[et]|and|or|not|if(?:else)?|for(?:all)?)(?=[()<>\[\]{}/%\s])`, KeywordReserved, nil},
			{Words(``, `(?=[()<>\[\]{}/%\s])`, `abs`, `add`, `aload`, `arc`, `arcn`, `array`, `atan`, `begin`, `bind`, `ceiling`, `charpath`, `clip`, `closepath`, `concat`, `concatmatrix`, `copy`, `cos`, `currentlinewidth`, `currentmatrix`, `currentpoint`, `curveto`, `cvi`, `cvs`, `def`, `defaultmatrix`, `dict`, `dictstackoverflow`, `div`, `dtransform`, `dup`, `end`, `exch`, `exec`, `exit`, `exp`, `fill`, `findfont`, `floor`, `get`, `getinterval`, `grestore`, `gsave`, `gt`, `identmatrix`, `idiv`, `idtransform`, `index`, `invertmatrix`, `itransform`, `length`, `lineto`, `ln`, `load`, `log`, `loop`, `matrix`, `mod`, `moveto`, `mul`, `neg`, `newpath`, `pathforall`, `pathbbox`, `pop`, `print`, `pstack`, `put`, `quit`, `rand`, `rangecheck`, `rcurveto`, `repeat`, `restore`, `rlineto`, `rmoveto`, `roll`, `rotate`, `round`, `run`, `save`, `scale`, `scalefont`, `setdash`, `setfont`, `setgray`, `setlinecap`, `setlinejoin`, `setlinewidth`, `setmatrix`, `setrgbcolor`, `shfill`, `show`, `showpage`, `sin`, `sqrt`, `stack`, `stringwidth`, `stroke`, `strokepath`, `sub`, `syntaxerror`, `transform`, `translate`, `truncate`, `typecheck`, `undefined`, `undefinedfilename`, `undefinedresult`), NameBuiltin, nil},
			{`\s+`, Text, nil},
		},
		"stringliteral": {
			{`[^()\\]+`, LiteralString, nil},
			{`\\`, LiteralStringEscape, Push("escape")},
			{`\(`, LiteralString, Push()},
			{`\)`, LiteralString, Pop(1)},
		},
		"escape": {
			{`[0-8]{3}|n|r|t|b|f|\\|\(|\)`, LiteralStringEscape, Pop(1)},
			Default(Pop(1)),
		},
	}
}
