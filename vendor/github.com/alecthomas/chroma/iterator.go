package chroma

import "strings"

// An Iterator across tokens.
//
// EOF will be returned at the end of the Token stream.
//
// If an error occurs within an Iterator, it may propagate this in a panic. Formatters should recover.
type Iterator func() Token

// Tokens consumes all tokens from the iterator and returns them as a slice.
func (i Iterator) Tokens() []Token {
	var out []Token
	for t := i(); t != EOF; t = i() {
		out = append(out, t)
	}
	return out
}

// Concaterator concatenates tokens from a series of iterators.
func Concaterator(iterators ...Iterator) Iterator {
	return func() Token {
		for len(iterators) > 0 {
			t := iterators[0]()
			if t != EOF {
				return t
			}
			iterators = iterators[1:]
		}
		return EOF
	}
}

// Literator converts a sequence of literal Tokens into an Iterator.
func Literator(tokens ...Token) Iterator {
	return func() Token {
		if len(tokens) == 0 {
			return EOF
		}
		token := tokens[0]
		tokens = tokens[1:]
		return token
	}
}

// SplitTokensIntoLines splits tokens containing newlines in two.
func SplitTokensIntoLines(tokens []Token) (out [][]Token) {
	var line []Token // nolint: prealloc
	for _, token := range tokens {
		for strings.Contains(token.Value, "\n") {
			parts := strings.SplitAfterN(token.Value, "\n", 2)
			// Token becomes the tail.
			token.Value = parts[1]

			// Append the head to the line and flush the line.
			clone := token.Clone()
			clone.Value = parts[0]
			line = append(line, clone)
			out = append(out, line)
			line = nil
		}
		line = append(line, token)
	}
	if len(line) > 0 {
		out = append(out, line)
	}
	// Strip empty trailing token line.
	if len(out) > 0 {
		last := out[len(out)-1]
		if len(last) == 1 && last[0].Value == "" {
			out = out[:len(out)-1]
		}
	}
	return
}
