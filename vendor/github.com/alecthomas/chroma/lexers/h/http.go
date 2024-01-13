package h

import (
	"strings"

	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// HTTP lexer.
var HTTP = internal.Register(httpBodyContentTypeLexer(MustNewLazyLexer(
	&Config{
		Name:         "HTTP",
		Aliases:      []string{"http"},
		Filenames:    []string{},
		MimeTypes:    []string{},
		NotMultiline: true,
		DotAll:       true,
	},
	httpRules,
)))

func httpRules() Rules {
	return Rules{
		"root": {
			{`(GET|POST|PUT|DELETE|HEAD|OPTIONS|TRACE|PATCH|CONNECT)( +)([^ ]+)( +)(HTTP)(/)([12]\.[01])(\r?\n|\Z)`, ByGroups(NameFunction, Text, NameNamespace, Text, KeywordReserved, Operator, LiteralNumber, Text), Push("headers")},
			{`(HTTP)(/)([12]\.[01])( +)(\d{3})( +)([^\r\n]+)(\r?\n|\Z)`, ByGroups(KeywordReserved, Operator, LiteralNumber, Text, LiteralNumber, Text, NameException, Text), Push("headers")},
		},
		"headers": {
			{`([^\s:]+)( *)(:)( *)([^\r\n]+)(\r?\n|\Z)`, EmitterFunc(httpHeaderBlock), nil},
			{`([\t ]+)([^\r\n]+)(\r?\n|\Z)`, EmitterFunc(httpContinuousHeaderBlock), nil},
			{`\r?\n`, Text, Push("content")},
		},
		"content": {
			{`.+`, EmitterFunc(httpContentBlock), nil},
		},
	}
}

func httpContentBlock(groups []string, state *LexerState) Iterator {
	tokens := []Token{
		{Generic, groups[0]},
	}
	return Literator(tokens...)
}

func httpHeaderBlock(groups []string, state *LexerState) Iterator {
	tokens := []Token{
		{Name, groups[1]},
		{Text, groups[2]},
		{Operator, groups[3]},
		{Text, groups[4]},
		{Literal, groups[5]},
		{Text, groups[6]},
	}
	return Literator(tokens...)
}

func httpContinuousHeaderBlock(groups []string, state *LexerState) Iterator {
	tokens := []Token{
		{Text, groups[1]},
		{Literal, groups[2]},
		{Text, groups[3]},
	}
	return Literator(tokens...)
}

func httpBodyContentTypeLexer(lexer Lexer) Lexer { return &httpBodyContentTyper{lexer} }

type httpBodyContentTyper struct{ Lexer }

func (d *httpBodyContentTyper) Tokenise(options *TokeniseOptions, text string) (Iterator, error) { // nolint: gocognit
	var contentType string
	var isContentType bool
	var subIterator Iterator

	it, err := d.Lexer.Tokenise(options, text)
	if err != nil {
		return nil, err
	}

	return func() Token {
		token := it()

		if token == EOF {
			if subIterator != nil {
				return subIterator()
			}
			return EOF
		}

		switch {
		case token.Type == Name && strings.ToLower(token.Value) == "content-type":
			{
				isContentType = true
			}
		case token.Type == Literal && isContentType:
			{
				isContentType = false
				contentType = strings.TrimSpace(token.Value)
				pos := strings.Index(contentType, ";")
				if pos > 0 {
					contentType = strings.TrimSpace(contentType[:pos])
				}
			}
		case token.Type == Generic && contentType != "":
			{
				lexer := internal.MatchMimeType(contentType)

				// application/calendar+xml can be treated as application/xml
				// if there's not a better match.
				if lexer == nil && strings.Contains(contentType, "+") {
					slashPos := strings.Index(contentType, "/")
					plusPos := strings.LastIndex(contentType, "+")
					contentType = contentType[:slashPos+1] + contentType[plusPos+1:]
					lexer = internal.MatchMimeType(contentType)
				}

				if lexer == nil {
					token.Type = Text
				} else {
					subIterator, err = lexer.Tokenise(nil, token.Value)
					if err != nil {
						panic(err)
					}
					return EOF
				}
			}
		}
		return token
	}, nil
}
