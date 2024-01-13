package g

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// GLSL lexer.
var GLSL = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "GLSL",
		Aliases:   []string{"glsl"},
		Filenames: []string{"*.vert", "*.frag", "*.geo"},
		MimeTypes: []string{"text/x-glslsrc"},
	},
	glslRules,
))

func glslRules() Rules {
	return Rules{
		"root": {
			{`^#.*`, CommentPreproc, nil},
			{`//.*`, CommentSingle, nil},
			{`/(\\\n)?[*](.|\n)*?[*](\\\n)?/`, CommentMultiline, nil},
			{`\+|-|~|!=?|\*|/|%|<<|>>|<=?|>=?|==?|&&?|\^|\|\|?`, Operator, nil},
			{`[?:]`, Operator, nil},
			{`\bdefined\b`, Operator, nil},
			{`[;{}(),\[\]]`, Punctuation, nil},
			{`[+-]?\d*\.\d+([eE][-+]?\d+)?`, LiteralNumberFloat, nil},
			{`[+-]?\d+\.\d*([eE][-+]?\d+)?`, LiteralNumberFloat, nil},
			{`0[xX][0-9a-fA-F]*`, LiteralNumberHex, nil},
			{`0[0-7]*`, LiteralNumberOct, nil},
			{`[1-9][0-9]*`, LiteralNumberInteger, nil},
			{Words(`\b`, `\b`, `attribute`, `const`, `uniform`, `varying`, `centroid`, `break`, `continue`, `do`, `for`, `while`, `if`, `else`, `in`, `out`, `inout`, `float`, `int`, `void`, `bool`, `true`, `false`, `invariant`, `discard`, `return`, `mat2`, `mat3mat4`, `mat2x2`, `mat3x2`, `mat4x2`, `mat2x3`, `mat3x3`, `mat4x3`, `mat2x4`, `mat3x4`, `mat4x4`, `vec2`, `vec3`, `vec4`, `ivec2`, `ivec3`, `ivec4`, `bvec2`, `bvec3`, `bvec4`, `sampler1D`, `sampler2D`, `sampler3DsamplerCube`, `sampler1DShadow`, `sampler2DShadow`, `struct`), Keyword, nil},
			{Words(`\b`, `\b`, `asm`, `class`, `union`, `enum`, `typedef`, `template`, `this`, `packed`, `goto`, `switch`, `default`, `inline`, `noinline`, `volatile`, `public`, `static`, `extern`, `external`, `interface`, `long`, `short`, `double`, `half`, `fixed`, `unsigned`, `lowp`, `mediump`, `highp`, `precision`, `input`, `output`, `hvec2`, `hvec3`, `hvec4`, `dvec2`, `dvec3`, `dvec4`, `fvec2`, `fvec3`, `fvec4`, `sampler2DRect`, `sampler3DRect`, `sampler2DRectShadow`, `sizeof`, `cast`, `namespace`, `using`), Keyword, nil},
			{`[a-zA-Z_]\w*`, Name, nil},
			{`\.`, Punctuation, nil},
			{`\s+`, Text, nil},
		},
	}
}
