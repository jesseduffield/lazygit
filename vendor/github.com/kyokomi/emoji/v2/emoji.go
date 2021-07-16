// Package emoji terminal output.
package emoji

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"unicode"
)

//go:generate generateEmojiCodeMap -pkg emoji -o emoji_codemap.go

// Replace Padding character for emoji.
var (
	ReplacePadding = " "
)

// CodeMap gets the underlying map of emoji.
func CodeMap() map[string]string {
	return emojiCode()
}

// RevCodeMap gets the underlying map of emoji.
func RevCodeMap() map[string][]string {
	return emojiRevCode()
}

func AliasList(shortCode string) []string {
	return emojiRevCode()[emojiCode()[shortCode]]
}

// HasAlias flags if the given `shortCode` has multiple aliases with other
// codes.
func HasAlias(shortCode string) bool {
	return len(AliasList(shortCode)) > 1
}

// NormalizeShortCode normalizes a given `shortCode` to a deterministic alias.
func NormalizeShortCode(shortCode string) string {
	shortLists := AliasList(shortCode)
	if len(shortLists) == 0 {
		return shortCode
	}
	return shortLists[0]
}

// regular expression that matches :flag-[countrycode]:
var flagRegexp = regexp.MustCompile(":flag-([a-z]{2}):")

func emojize(x string) string {
	str, ok := emojiCode()[x]
	if ok {
		return str + ReplacePadding
	}
	if match := flagRegexp.FindStringSubmatch(x); len(match) == 2 {
		return regionalIndicator(match[1][0]) + regionalIndicator(match[1][1])
	}
	return x
}

// regionalIndicator maps a lowercase letter to a unicode regional indicator
func regionalIndicator(i byte) string {
	return string('\U0001F1E6' + rune(i) - 'a')
}

func replaseEmoji(input *bytes.Buffer) string {
	emoji := bytes.NewBufferString(":")
	for {
		i, _, err := input.ReadRune()
		if err != nil {
			// not replase
			return emoji.String()
		}

		if i == ':' && emoji.Len() == 1 {
			return emoji.String() + replaseEmoji(input)
		}

		emoji.WriteRune(i)
		switch {
		case unicode.IsSpace(i):
			return emoji.String()
		case i == ':':
			return emojize(emoji.String())
		}
	}
}

func compile(x string) string {
	if x == "" {
		return ""
	}

	input := bytes.NewBufferString(x)
	output := bytes.NewBufferString("")

	for {
		i, _, err := input.ReadRune()
		if err != nil {
			break
		}
		switch i {
		default:
			output.WriteRune(i)
		case ':':
			output.WriteString(replaseEmoji(input))
		}
	}
	return output.String()
}

// Print is fmt.Print which supports emoji
func Print(a ...interface{}) (int, error) {
	return fmt.Print(compile(fmt.Sprint(a...)))
}

// Println is fmt.Println which supports emoji
func Println(a ...interface{}) (int, error) {
	return fmt.Println(compile(fmt.Sprint(a...)))
}

// Printf is fmt.Printf which supports emoji
func Printf(format string, a ...interface{}) (int, error) {
	return fmt.Print(compile(fmt.Sprintf(format, a...)))
}

// Fprint is fmt.Fprint which supports emoji
func Fprint(w io.Writer, a ...interface{}) (int, error) {
	return fmt.Fprint(w, compile(fmt.Sprint(a...)))
}

// Fprintln is fmt.Fprintln which supports emoji
func Fprintln(w io.Writer, a ...interface{}) (int, error) {
	return fmt.Fprintln(w, compile(fmt.Sprint(a...)))
}

// Fprintf is fmt.Fprintf which supports emoji
func Fprintf(w io.Writer, format string, a ...interface{}) (int, error) {
	return fmt.Fprint(w, compile(fmt.Sprintf(format, a...)))
}

// Sprint is fmt.Sprint which supports emoji
func Sprint(a ...interface{}) string {
	return compile(fmt.Sprint(a...))
}

// Sprintf is fmt.Sprintf which supports emoji
func Sprintf(format string, a ...interface{}) string {
	return compile(fmt.Sprintf(format, a...))
}

// Errorf is fmt.Errorf which supports emoji
func Errorf(format string, a ...interface{}) error {
	return errors.New(compile(Sprintf(format, a...)))
}
