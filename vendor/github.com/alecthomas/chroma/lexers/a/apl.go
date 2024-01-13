package a

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Apl lexer.
var Apl = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "APL",
		Aliases:   []string{"apl"},
		Filenames: []string{"*.apl"},
		MimeTypes: []string{},
	},
	aplRules,
))

func aplRules() Rules {
	return Rules{
		"root": {
			{`\s+`, Text, nil},
			{`[⍝#].*$`, CommentSingle, nil},
			{`\'((\'\')|[^\'])*\'`, LiteralStringSingle, nil},
			{`"(("")|[^"])*"`, LiteralStringDouble, nil},
			{`[⋄◇()]`, Punctuation, nil},
			{`[\[\];]`, LiteralStringRegex, nil},
			{`⎕[A-Za-zΔ∆⍙][A-Za-zΔ∆⍙_¯0-9]*`, NameFunction, nil},
			{`[A-Za-zΔ∆⍙_][A-Za-zΔ∆⍙_¯0-9]*`, NameVariable, nil},
			{`¯?(0[Xx][0-9A-Fa-f]+|[0-9]*\.?[0-9]+([Ee][+¯]?[0-9]+)?|¯|∞)([Jj]¯?(0[Xx][0-9A-Fa-f]+|[0-9]*\.?[0-9]+([Ee][+¯]?[0-9]+)?|¯|∞))?`, LiteralNumber, nil},
			{`[\.\\/⌿⍀¨⍣⍨⍠⍤∘⍥@⌺⌶⍢]`, NameAttribute, nil},
			{`[+\-×÷⌈⌊∣|⍳?*⍟○!⌹<≤=>≥≠≡≢∊⍷∪∩~∨∧⍱⍲⍴,⍪⌽⊖⍉↑↓⊂⊃⌷⍋⍒⊤⊥⍕⍎⊣⊢⍁⍂≈⌸⍯↗⊆⍸]`, Operator, nil},
			{`⍬`, NameConstant, nil},
			{`[⎕⍞]`, NameVariableGlobal, nil},
			{`[←→]`, KeywordDeclaration, nil},
			{`[⍺⍵⍶⍹∇:]`, NameBuiltinPseudo, nil},
			{`[{}]`, KeywordType, nil},
		},
	}
}
