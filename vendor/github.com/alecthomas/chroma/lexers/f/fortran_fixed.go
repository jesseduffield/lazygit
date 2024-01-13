package f

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// FortranFixed lexer.
var FortranFixed = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "FortranFixed",
		Aliases:         []string{"fortranfixed"},
		Filenames:       []string{"*.f", "*.F"},
		MimeTypes:       []string{"text/x-fortran"},
		NotMultiline:    true,
		CaseInsensitive: true,
	},
	func() Rules {
		return Rules{
			"root": {
				{`[C*].*\n`, Comment, nil},
				{`#.*\n`, CommentPreproc, nil},
				{`[\t ]*!.*\n`, Comment, nil},
				{`(.{5})`, NameLabel, Push("cont-char")},
				{`.*\n`, Using(Fortran), nil},
			},
			"cont-char": {
				{` `, Text, Push("code")},
				{`0`, Comment, Push("code")},
				{`.`, GenericStrong, Push("code")},
			},
			"code": {
				{`(.{66})(.*)(\n)`, ByGroups(Using(Fortran), Comment, Text), Push("root")},
				{`.*\n`, Using(Fortran), Push("root")},
				Default(Push("root")),
			},
		}
	},
))
