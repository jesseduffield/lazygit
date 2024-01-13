// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// tokenType identifies the type of lexical tokens.
type tokenType int

// String returns a string representation of the token type.
func (t tokenType) String() string {
	return tokenNames[t]
}

// Token represents a token and the corresponding string.
type Token struct {
	Type   tokenType
	Value  string
	Line   int
	Column int
}

// String returns a string representation of the token.
func (t *Token) String() string {
	if len(t.Value) > 10 {
		return fmt.Sprintf("%s (line: %d, column: %d): %.10q...",
			t.Type, t.Line, t.Column, t.Value)
	}
	return fmt.Sprintf("%s (line: %d, column: %d): %q",
		t.Type, t.Line, t.Column, t.Value)
}

// All tokens -----------------------------------------------------------------

// The complete list of tokens in CSS3.
const (
	// Scanner flags.
	TokenError tokenType = iota
	TokenEOF
	// From now on, only tokens from the CSS specification.
	TokenIdent
	TokenAtKeyword
	TokenString
	TokenHash
	TokenNumber
	TokenPercentage
	TokenDimension
	TokenURI
	TokenUnicodeRange
	TokenCDO
	TokenCDC
	TokenS
	TokenComment
	TokenFunction
	TokenIncludes
	TokenDashMatch
	TokenPrefixMatch
	TokenSuffixMatch
	TokenSubstringMatch
	TokenChar
	TokenBOM
)

// tokenNames maps tokenType's to their names. Used for conversion to string.
var tokenNames = map[tokenType]string{
	TokenError:          "error",
	TokenEOF:            "EOF",
	TokenIdent:          "IDENT",
	TokenAtKeyword:      "ATKEYWORD",
	TokenString:         "STRING",
	TokenHash:           "HASH",
	TokenNumber:         "NUMBER",
	TokenPercentage:     "PERCENTAGE",
	TokenDimension:      "DIMENSION",
	TokenURI:            "URI",
	TokenUnicodeRange:   "UNICODE-RANGE",
	TokenCDO:            "CDO",
	TokenCDC:            "CDC",
	TokenS:              "S",
	TokenComment:        "COMMENT",
	TokenFunction:       "FUNCTION",
	TokenIncludes:       "INCLUDES",
	TokenDashMatch:      "DASHMATCH",
	TokenPrefixMatch:    "PREFIXMATCH",
	TokenSuffixMatch:    "SUFFIXMATCH",
	TokenSubstringMatch: "SUBSTRINGMATCH",
	TokenChar:           "CHAR",
	TokenBOM:            "BOM",
}

// Macros and productions -----------------------------------------------------
// http://www.w3.org/TR/css3-syntax/#tokenization

var macroRegexp = regexp.MustCompile(`\{[a-z]+\}`)

// macros maps macro names to patterns to be expanded.
var macros = map[string]string{
	// must be escaped: `\.+*?()|[]{}^$`
	"ident":      `-?{nmstart}{nmchar}*`,
	"name":       `{nmchar}+`,
	"nmstart":    `[a-zA-Z_]|{nonascii}|{escape}`,
	"nonascii":   "[\u0080-\uD7FF\uE000-\uFFFD\U00010000-\U0010FFFF]",
	"unicode":    `\\[0-9a-fA-F]{1,6}{wc}?`,
	"escape":     "{unicode}|\\\\[\u0020-\u007E\u0080-\uD7FF\uE000-\uFFFD\U00010000-\U0010FFFF]",
	"nmchar":     `[a-zA-Z0-9_-]|{nonascii}|{escape}`,
	"num":        `[0-9]*\.[0-9]+|[0-9]+`,
	"string":     `"(?:{stringchar}|')*"|'(?:{stringchar}|")*'`,
	"stringchar": `{urlchar}|[ ]|\\{nl}`,
	"nl":         `[\n\r\f]|\r\n`,
	"w":          `{wc}*`,
	"wc":         `[\t\n\f\r ]`,

	// urlchar should accept [(ascii characters minus those that need escaping)|{nonascii}|{escape}]
	// ASCII characters range = `[\u0020-\u007e]`
	// Skip space \u0020 = `[\u0021-\u007e]`
	// Skip quotation mark \0022 = `[\u0021\u0023-\u007e]`
	// Skip apostrophe \u0027 = `[\u0021\u0023-\u0026\u0028-\u007e]`
	// Skip reverse solidus \u005c = `[\u0021\u0023-\u0026\u0028-\u005b\u005d\u007e]`
	// Finally, the left square bracket (\u005b) and right (\u005d) needs escaping themselves
	"urlchar": "[\u0021\u0023-\u0026\u0028-\\\u005b\\\u005d-\u007E]|{nonascii}|{escape}",
}

// productions maps the list of tokens to patterns to be expanded.
var productions = map[tokenType]string{
	// Unused regexps (matched using other methods) are commented out.
	TokenIdent:        `{ident}`,
	TokenAtKeyword:    `@{ident}`,
	TokenString:       `{string}`,
	TokenHash:         `#{name}`,
	TokenNumber:       `{num}`,
	TokenPercentage:   `{num}%`,
	TokenDimension:    `{num}{ident}`,
	TokenURI:          `url\({w}(?:{string}|{urlchar}*?){w}\)`,
	TokenUnicodeRange: `U\+[0-9A-F\?]{1,6}(?:-[0-9A-F]{1,6})?`,
	//TokenCDO:            `<!--`,
	TokenCDC:      `-->`,
	TokenS:        `{wc}+`,
	TokenComment:  `/\*[^\*]*[\*]+(?:[^/][^\*]*[\*]+)*/`,
	TokenFunction: `{ident}\(`,
	//TokenIncludes:       `~=`,
	//TokenDashMatch:      `\|=`,
	//TokenPrefixMatch:    `\^=`,
	//TokenSuffixMatch:    `\$=`,
	//TokenSubstringMatch: `\*=`,
	//TokenChar:           `[^"']`,
	//TokenBOM:            "\uFEFF",
}

// matchers maps the list of tokens to compiled regular expressions.
//
// The map is filled on init() using the macros and productions defined in
// the CSS specification.
var matchers = map[tokenType]*regexp.Regexp{}

// matchOrder is the order to test regexps when first-char shortcuts
// can't be used.
var matchOrder = []tokenType{
	TokenURI,
	TokenFunction,
	TokenUnicodeRange,
	TokenIdent,
	TokenDimension,
	TokenPercentage,
	TokenNumber,
	TokenCDC,
}

func init() {
	// replace macros and compile regexps for productions.
	replaceMacro := func(s string) string {
		return "(?:" + macros[s[1:len(s)-1]] + ")"
	}
	for t, s := range productions {
		for macroRegexp.MatchString(s) {
			s = macroRegexp.ReplaceAllStringFunc(s, replaceMacro)
		}
		matchers[t] = regexp.MustCompile("^(?:" + s + ")")
	}
}

// Scanner --------------------------------------------------------------------

// New returns a new CSS scanner for the given input.
func New(input string) *Scanner {
	// Normalize newlines.
	input = strings.Replace(input, "\r\n", "\n", -1)
	return &Scanner{
		input: input,
		row:   1,
		col:   1,
	}
}

// Scanner scans an input and emits tokens following the CSS3 specification.
type Scanner struct {
	input string
	pos   int
	row   int
	col   int
	err   *Token
}

// Next returns the next token from the input.
//
// At the end of the input the token type is TokenEOF.
//
// If the input can't be tokenized the token type is TokenError. This occurs
// in case of unclosed quotation marks or comments.
func (s *Scanner) Next() *Token {
	if s.err != nil {
		return s.err
	}
	if s.pos >= len(s.input) {
		s.err = &Token{TokenEOF, "", s.row, s.col}
		return s.err
	}
	if s.pos == 0 {
		// Test BOM only once, at the beginning of the file.
		if strings.HasPrefix(s.input, "\uFEFF") {
			return s.emitSimple(TokenBOM, "\uFEFF")
		}
	}
	// There's a lot we can guess based on the first byte so we'll take a
	// shortcut before testing multiple regexps.
	input := s.input[s.pos:]
	switch input[0] {
	case '\t', '\n', '\f', '\r', ' ':
		// Whitespace.
		return s.emitToken(TokenS, matchers[TokenS].FindString(input))
	case '.':
		// Dot is too common to not have a quick check.
		// We'll test if this is a Char; if it is followed by a number it is a
		// dimension/percentage/number, and this will be matched later.
		if len(input) > 1 && !unicode.IsDigit(rune(input[1])) {
			return s.emitSimple(TokenChar, ".")
		}
	case '#':
		// Another common one: Hash or Char.
		if match := matchers[TokenHash].FindString(input); match != "" {
			return s.emitToken(TokenHash, match)
		}
		return s.emitSimple(TokenChar, "#")
	case '@':
		// Another common one: AtKeyword or Char.
		if match := matchers[TokenAtKeyword].FindString(input); match != "" {
			return s.emitSimple(TokenAtKeyword, match)
		}
		return s.emitSimple(TokenChar, "@")
	case ':', ',', ';', '%', '&', '+', '=', '>', '(', ')', '[', ']', '{', '}':
		// More common chars.
		return s.emitSimple(TokenChar, string(input[0]))
	case '"', '\'':
		// String or error.
		match := matchers[TokenString].FindString(input)
		if match != "" {
			return s.emitToken(TokenString, match)
		}

		s.err = &Token{TokenError, "unclosed quotation mark", s.row, s.col}
		return s.err
	case '/':
		// Comment, error or Char.
		if len(input) > 1 && input[1] == '*' {
			match := matchers[TokenComment].FindString(input)
			if match != "" {
				return s.emitToken(TokenComment, match)
			} else {
				s.err = &Token{TokenError, "unclosed comment", s.row, s.col}
				return s.err
			}
		}
		return s.emitSimple(TokenChar, "/")
	case '~':
		// Includes or Char.
		return s.emitPrefixOrChar(TokenIncludes, "~=")
	case '|':
		// DashMatch or Char.
		return s.emitPrefixOrChar(TokenDashMatch, "|=")
	case '^':
		// PrefixMatch or Char.
		return s.emitPrefixOrChar(TokenPrefixMatch, "^=")
	case '$':
		// SuffixMatch or Char.
		return s.emitPrefixOrChar(TokenSuffixMatch, "$=")
	case '*':
		// SubstringMatch or Char.
		return s.emitPrefixOrChar(TokenSubstringMatch, "*=")
	case '<':
		// CDO or Char.
		return s.emitPrefixOrChar(TokenCDO, "<!--")
	}
	// Test all regexps, in order.
	for _, token := range matchOrder {
		if match := matchers[token].FindString(input); match != "" {
			return s.emitToken(token, match)
		}
	}
	// We already handled unclosed quotation marks and comments,
	// so this can only be a Char.
	r, width := utf8.DecodeRuneInString(input)
	token := &Token{TokenChar, string(r), s.row, s.col}
	s.col += width
	s.pos += width
	return token
}

// updatePosition updates input coordinates based on the consumed text.
func (s *Scanner) updatePosition(text string) {
	width := utf8.RuneCountInString(text)
	lines := strings.Count(text, "\n")
	s.row += lines
	if lines == 0 {
		s.col += width
	} else {
		s.col = utf8.RuneCountInString(text[strings.LastIndex(text, "\n"):])
	}
	s.pos += len(text) // while col is a rune index, pos is a byte index
}

// emitToken returns a Token for the string v and updates the scanner position.
func (s *Scanner) emitToken(t tokenType, v string) *Token {
	token := &Token{t, v, s.row, s.col}
	s.updatePosition(v)
	return token
}

// emitSimple returns a Token for the string v and updates the scanner
// position in a simplified manner.
//
// The string is known to have only ASCII characters and to not have a newline.
func (s *Scanner) emitSimple(t tokenType, v string) *Token {
	token := &Token{t, v, s.row, s.col}
	s.col += len(v)
	s.pos += len(v)
	return token
}

// emitPrefixOrChar returns a Token for type t if the current position
// matches the given prefix. Otherwise it returns a Char token using the
// first character from the prefix.
//
// The prefix is known to have only ASCII characters and to not have a newline.
func (s *Scanner) emitPrefixOrChar(t tokenType, prefix string) *Token {
	if strings.HasPrefix(s.input[s.pos:], prefix) {
		return s.emitSimple(t, prefix)
	}
	return s.emitSimple(TokenChar, string(prefix[0]))
}
