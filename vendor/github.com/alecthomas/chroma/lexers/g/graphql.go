package g

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Go lexer.
var Graphql = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "GraphQL",
		Aliases:   []string{"graphql", "graphqls", "gql"},
		Filenames: []string{"*.graphql", "*.graphqls"},
	},
	graphqlRules,
))

func graphqlRules() Rules {
	return Rules{
		"root": {
			{`(query|mutation|subscription|fragment|scalar|implements|interface|union|enum|input|type)`, KeywordDeclaration, Push("type")},
			{`(on|extend|schema|directive|\.\.\.)`, KeywordDeclaration, nil},
			{`(QUERY|MUTATION|SUBSCRIPTION|FIELD|FRAGMENT_DEFINITION|FRAGMENT_SPREAD|INLINE_FRAGMENT|SCHEMA|SCALAR|OBJECT|FIELD_DEFINITION|ARGUMENT_DEFINITION|INTERFACE|UNION|ENUM|ENUM_VALUE|INPUT_OBJECT|INPUT_FIELD_DEFINITION)\b`, KeywordConstant, nil},
			{`[^\W\d]\w*`, NameProperty, nil},
			{`\@\w+`, NameDecorator, nil},
			{`:`, Punctuation, Push("type")},
			{`[\(\)\{\}\[\],!\|=]`, Punctuation, nil},
			{`\$\w+`, NameVariable, nil},
			{`\d+i`, LiteralNumber, nil},
			{`\d+\.\d*([Ee][-+]\d+)?i`, LiteralNumber, nil},
			{`\.\d+([Ee][-+]\d+)?i`, LiteralNumber, nil},
			{`\d+[Ee][-+]\d+i`, LiteralNumber, nil},
			{`\d+(\.\d+[eE][+\-]?\d+|\.\d*|[eE][+\-]?\d+)`, LiteralNumberFloat, nil},
			{`\.\d+([eE][+\-]?\d+)?`, LiteralNumberFloat, nil},
			{`(0|[1-9][0-9]*)`, LiteralNumberInteger, nil},
			{`"""[\x00-\x7F]*?"""`, LiteralString, nil},
			{`"(\\["\\abfnrtv]|\\x[0-9a-fA-F]{2}|\\[0-7]{1,3}|\\u[0-9a-fA-F]{4}|\\U[0-9a-fA-F]{8}|[^\\])"`, LiteralStringChar, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralString, nil},
			{`"(true|false|null)*"`, Literal, nil},
			{`[\r\n\s]+`, Whitespace, nil},
			{`#[^\r\n]*`, Comment, nil},
		},
		// Treats the next word as a class, default rules it would be a property
		"type": {
			{`[^\W\d]\w*`, NameClass, Pop(1)},
			Include("root"),
		},
	}
}
