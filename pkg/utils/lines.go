package utils

import "strings"

// SplitLines takes a multiline string and splits it on newlines
// currently we are also stripping \r's which may have adverse effects for
// windows users (but no issues have been raised yet)
func SplitLines(multilineString string) []string {
	multilineString = strings.Replace(multilineString, "\r", "", -1)
	if multilineString == "" || multilineString == "\n" {
		return make([]string, 0)
	}
	lines := strings.Split(multilineString, "\n")
	if lines[len(lines)-1] == "" {
		return lines[:len(lines)-1]
	}
	return lines
}

// NormalizeLinefeeds - Removes all Windows and Mac style line feeds
func NormalizeLinefeeds(str string) string {
	str = strings.Replace(str, "\r\n", "\n", -1)
	str = strings.Replace(str, "\r", "", -1)
	return str
}

// EscapeSpecialChars - Replaces all special chars like \n with \\n
func EscapeSpecialChars(str string) string {
	return strings.NewReplacer(
		"\n", "\\n",
		"\r", "\\r",
		"\t", "\\t",
		"\b", "\\b",
		"\f", "\\f",
		"\v", "\\v",
	).Replace(str)
}
