package chroma

import (
	"fmt"
	"strings"
)

var (
	defaultOptions = &TokeniseOptions{
		State:    "root",
		EnsureLF: true,
	}
)

// Config for a lexer.
type Config struct {
	// Name of the lexer.
	Name string

	// Shortcuts for the lexer
	Aliases []string

	// File name globs
	Filenames []string

	// Secondary file name globs
	AliasFilenames []string

	// MIME types
	MimeTypes []string

	// Regex matching is case-insensitive.
	CaseInsensitive bool

	// Regex matches all characters.
	DotAll bool

	// Regex does not match across lines ($ matches EOL).
	//
	// Defaults to multiline.
	NotMultiline bool

	// Don't strip leading and trailing newlines from the input.
	// DontStripNL bool

	// Strip all leading and trailing whitespace from the input
	// StripAll bool

	// Make sure that the input ends with a newline. This
	// is required for some lexers that consume input linewise.
	EnsureNL bool

	// If given and greater than 0, expand tabs in the input.
	// TabSize int

	// Priority of lexer.
	//
	// If this is 0 it will be treated as a default of 1.
	Priority float32
}

// Token output to formatter.
type Token struct {
	Type  TokenType `json:"type"`
	Value string    `json:"value"`
}

func (t *Token) String() string   { return t.Value }
func (t *Token) GoString() string { return fmt.Sprintf("&Token{%s, %q}", t.Type, t.Value) }

// Clone returns a clone of the Token.
func (t *Token) Clone() Token {
	return *t
}

// EOF is returned by lexers at the end of input.
var EOF Token

// TokeniseOptions contains options for tokenisers.
type TokeniseOptions struct {
	// State to start tokenisation in. Defaults to "root".
	State string
	// Nested tokenisation.
	Nested bool

	// If true, all EOLs are converted into LF
	// by replacing CRLF and CR
	EnsureLF bool
}

// A Lexer for tokenising source code.
type Lexer interface {
	// Config describing the features of the Lexer.
	Config() *Config
	// Tokenise returns an Iterator over tokens in text.
	Tokenise(options *TokeniseOptions, text string) (Iterator, error)
}

// Lexers is a slice of lexers sortable by name.
type Lexers []Lexer

func (l Lexers) Len() int      { return len(l) }
func (l Lexers) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l Lexers) Less(i, j int) bool {
	return strings.ToLower(l[i].Config().Name) < strings.ToLower(l[j].Config().Name)
}

// PrioritisedLexers is a slice of lexers sortable by priority.
type PrioritisedLexers []Lexer

func (l PrioritisedLexers) Len() int      { return len(l) }
func (l PrioritisedLexers) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l PrioritisedLexers) Less(i, j int) bool {
	ip := l[i].Config().Priority
	if ip == 0 {
		ip = 1
	}
	jp := l[j].Config().Priority
	if jp == 0 {
		jp = 1
	}
	return ip > jp
}

// Analyser determines how appropriate this lexer is for the given text.
type Analyser interface {
	AnalyseText(text string) float32
}
