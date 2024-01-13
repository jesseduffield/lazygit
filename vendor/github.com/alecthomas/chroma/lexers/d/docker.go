package d

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/b"
	"github.com/alecthomas/chroma/lexers/internal"
	"github.com/alecthomas/chroma/lexers/j"
)

// Docker lexer.
var Docker = internal.Register(MustNewLazyLexer(
	&Config{
		Name:            "Docker",
		Aliases:         []string{"docker", "dockerfile"},
		Filenames:       []string{"Dockerfile", "*.docker"},
		MimeTypes:       []string{"text/x-dockerfile-config"},
		CaseInsensitive: true,
	},
	dockerRules,
))

func dockerRules() Rules {
	return Rules{
		"root": {
			{`#.*`, Comment, nil},
			{`(ONBUILD)((?:\s*\\?\s*))`, ByGroups(Keyword, Using(b.Bash)), nil},
			{`(HEALTHCHECK)(((?:\s*\\?\s*)--\w+=\w+(?:\s*\\?\s*))*)`, ByGroups(Keyword, Using(b.Bash)), nil},
			{`(VOLUME|ENTRYPOINT|CMD|SHELL)((?:\s*\\?\s*))(\[.*?\])`, ByGroups(Keyword, Using(b.Bash), Using(j.JSON)), nil},
			{`(LABEL|ENV|ARG)((?:(?:\s*\\?\s*)\w+=\w+(?:\s*\\?\s*))*)`, ByGroups(Keyword, Using(b.Bash)), nil},
			{`((?:FROM|MAINTAINER|EXPOSE|WORKDIR|USER|STOPSIGNAL)|VOLUME)\b(.*)`, ByGroups(Keyword, LiteralString), nil},
			{`((?:RUN|CMD|ENTRYPOINT|ENV|ARG|LABEL|ADD|COPY))`, Keyword, nil},
			{`(.*\\\n)*.+`, Using(b.Bash), nil},
		},
	}
}
