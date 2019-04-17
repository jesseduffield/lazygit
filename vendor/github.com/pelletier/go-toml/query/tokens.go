package query

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"strconv"
	"unicode"
)

// Define tokens
type tokenType int

const (
	eof = -(iota + 1)
)

const (
	tokenError tokenType = iota
	tokenEOF
	tokenKey
	tokenString
	tokenInteger
	tokenFloat
	tokenLeftBracket
	tokenRightBracket
	tokenLeftParen
	tokenRightParen
	tokenComma
	tokenColon
	tokenDollar
	tokenStar
	tokenQuestion
	tokenDot
	tokenDotDot
)

var tokenTypeNames = []string{
	"Error",
	"EOF",
	"Key",
	"String",
	"Integer",
	"Float",
	"[",
	"]",
	"(",
	")",
	",",
	":",
	"$",
	"*",
	"?",
	".",
	"..",
}

type token struct {
	toml.Position
	typ tokenType
	val string
}

func (tt tokenType) String() string {
	idx := int(tt)
	if idx < len(tokenTypeNames) {
		return tokenTypeNames[idx]
	}
	return "Unknown"
}

func (t token) Int() int {
	if result, err := strconv.Atoi(t.val); err != nil {
		panic(err)
	} else {
		return result
	}
}

func (t token) String() string {
	switch t.typ {
	case tokenEOF:
		return "EOF"
	case tokenError:
		return t.val
	}

	return fmt.Sprintf("%q", t.val)
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isAlphanumeric(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isDigit(r rune) bool {
	return unicode.IsNumber(r)
}

func isHexDigit(r rune) bool {
	return isDigit(r) ||
		(r >= 'a' && r <= 'f') ||
		(r >= 'A' && r <= 'F')
}
