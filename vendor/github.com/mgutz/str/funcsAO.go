package str

import (
	"fmt"
	"html"
	//"log"
	"regexp"
	"strings"
)

// Verbose flag enables console output for those functions that have
// counterparts in Go's excellent stadard packages.
var Verbose = false
var templateOpen = "{{"
var templateClose = "}}"

var beginEndSpacesRe = regexp.MustCompile("^\\s+|\\s+$")
var camelizeRe = regexp.MustCompile(`(\-|_|\s)+(.)?`)
var camelizeRe2 = regexp.MustCompile(`(\-|_|\s)+`)
var capitalsRe = regexp.MustCompile("([A-Z])")
var dashSpaceRe = regexp.MustCompile(`[-\s]+`)
var dashesRe = regexp.MustCompile("-+")
var isAlphaNumericRe = regexp.MustCompile(`[^0-9a-z\xC0-\xFF]`)
var isAlphaRe = regexp.MustCompile(`[^a-z\xC0-\xFF]`)
var nWhitespaceRe = regexp.MustCompile(`\s+`)
var notDigitsRe = regexp.MustCompile(`[^0-9]`)
var slugifyRe = regexp.MustCompile(`[^\w\s\-]`)
var spaceUnderscoreRe = regexp.MustCompile("[_\\s]+")
var spacesRe = regexp.MustCompile("[\\s\\xA0]+")
var stripPuncRe = regexp.MustCompile(`[^\w\s]|_`)
var templateRe = regexp.MustCompile(`([\-\[\]()*\s])`)
var templateRe2 = regexp.MustCompile(`\$`)
var underscoreRe = regexp.MustCompile(`([a-z\d])([A-Z]+)`)
var whitespaceRe = regexp.MustCompile(`^[\s\xa0]*$`)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Between extracts a string between left and right strings.
func Between(s, left, right string) string {
	l := len(left)
	startPos := strings.Index(s, left)
	if startPos < 0 {
		return ""
	}
	endPos := IndexOf(s, right, startPos+l)
	//log.Printf("%s: left %s right %s start %d end %d", s, left, right, startPos+l, endPos)
	if endPos < 0 {
		return ""
	} else if right == "" {
		return s[endPos:]
	} else {
		return s[startPos+l : endPos]
	}
}

// BetweenF is the filter form for Between.
func BetweenF(left, right string) func(string) string {
	return func(s string) string {
		return Between(s, left, right)
	}
}

// Camelize return new string which removes any underscores or dashes and convert a string into camel casing.
func Camelize(s string) string {
	return camelizeRe.ReplaceAllStringFunc(s, func(val string) string {
		val = strings.ToUpper(val)
		val = camelizeRe2.ReplaceAllString(val, "")
		return val
	})
}

// Capitalize uppercases the first char of s and lowercases the rest.
func Capitalize(s string) string {
	return strings.ToUpper(s[0:1]) + strings.ToLower(s[1:])
}

// CharAt returns a string from the character at the specified position.
func CharAt(s string, index int) string {
	l := len(s)
	shortcut := index < 0 || index > l-1 || l == 0
	if shortcut {
		return ""
	}
	return s[index : index+1]
}

// CharAtF is the filter form of CharAt.
func CharAtF(index int) func(string) string {
	return func(s string) string {
		return CharAt(s, index)
	}
}

// ChompLeft removes prefix at the start of a string.
func ChompLeft(s, prefix string) string {
	if strings.HasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s
}

// ChompLeftF is the filter form of ChompLeft.
func ChompLeftF(prefix string) func(string) string {
	return func(s string) string {
		return ChompLeft(s, prefix)
	}
}

// ChompRight removes suffix from end of s.
func ChompRight(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	}
	return s
}

// ChompRightF is the filter form of ChompRight.
func ChompRightF(suffix string) func(string) string {
	return func(s string) string {
		return ChompRight(s, suffix)
	}
}

// Classify returns a camelized string with the first letter upper cased.
func Classify(s string) string {
	return Camelize("-" + s)
}

// ClassifyF is the filter form of Classify.
func ClassifyF(s string) func(string) string {
	return func(s string) string {
		return Classify(s)
	}
}

// Clean compresses all adjacent whitespace to a single space and trims s.
func Clean(s string) string {
	s = spacesRe.ReplaceAllString(s, " ")
	s = beginEndSpacesRe.ReplaceAllString(s, "")
	return s
}

// Dasherize  converts a camel cased string into a string delimited by dashes.
func Dasherize(s string) string {
	s = strings.TrimSpace(s)
	s = spaceUnderscoreRe.ReplaceAllString(s, "-")
	s = capitalsRe.ReplaceAllString(s, "-$1")
	s = dashesRe.ReplaceAllString(s, "-")
	s = strings.ToLower(s)
	return s
}

// EscapeHTML is alias for html.EscapeString.
func EscapeHTML(s string) string {
	if Verbose {
		fmt.Println("Use html.EscapeString instead of EscapeHTML")
	}
	return html.EscapeString(s)
}

// DecodeHTMLEntities decodes HTML entities into their proper string representation.
// DecodeHTMLEntities is an alias for html.UnescapeString
func DecodeHTMLEntities(s string) string {
	if Verbose {
		fmt.Println("Use html.UnescapeString instead of DecodeHTMLEntities")
	}
	return html.UnescapeString(s)
}

// EnsurePrefix ensures s starts with prefix.
func EnsurePrefix(s, prefix string) string {
	if strings.HasPrefix(s, prefix) {
		return s
	}
	return prefix + s
}

// EnsurePrefixF is the filter form of EnsurePrefix.
func EnsurePrefixF(prefix string) func(string) string {
	return func(s string) string {
		return EnsurePrefix(s, prefix)
	}
}

// EnsureSuffix ensures s ends with suffix.
func EnsureSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s
	}
	return s + suffix
}

// EnsureSuffixF is the filter form of EnsureSuffix.
func EnsureSuffixF(suffix string) func(string) string {
	return func(s string) string {
		return EnsureSuffix(s, suffix)
	}
}

// Humanize transforms s into a human friendly form.
func Humanize(s string) string {
	if s == "" {
		return s
	}
	s = Underscore(s)
	var humanizeRe = regexp.MustCompile(`_id$`)
	s = humanizeRe.ReplaceAllString(s, "")
	s = strings.Replace(s, "_", " ", -1)
	s = strings.TrimSpace(s)
	s = Capitalize(s)
	return s
}

// Iif is short for immediate if. If condition is true return truthy else falsey.
func Iif(condition bool, truthy string, falsey string) string {
	if condition {
		return truthy
	}
	return falsey
}

// IndexOf finds the index of needle in s starting from start.
func IndexOf(s string, needle string, start int) int {
	l := len(s)
	if needle == "" {
		if start < 0 {
			return 0
		} else if start < l {
			return start
		} else {
			return l
		}
	}
	if start < 0 || start > l-1 {
		return -1
	}
	pos := strings.Index(s[start:], needle)
	if pos == -1 {
		return -1
	}
	return start + pos
}

// IsAlpha returns true if a string contains only letters from ASCII (a-z,A-Z). Other letters from other languages are not supported.
func IsAlpha(s string) bool {
	return !isAlphaRe.MatchString(strings.ToLower(s))
}

// IsAlphaNumeric returns true if a string contains letters and digits.
func IsAlphaNumeric(s string) bool {
	return !isAlphaNumericRe.MatchString(strings.ToLower(s))
}

// IsLower returns true if s comprised of all lower case characters.
func IsLower(s string) bool {
	return IsAlpha(s) && s == strings.ToLower(s)
}

// IsNumeric returns true if a string contains only digits from 0-9. Other digits not in Latin (such as Arabic) are not currently supported.
func IsNumeric(s string) bool {
	return !notDigitsRe.MatchString(s)
}

// IsUpper returns true if s contains all upper case chracters.
func IsUpper(s string) bool {
	return IsAlpha(s) && s == strings.ToUpper(s)
}

// IsEmpty returns true if the string is solely composed of whitespace.
func IsEmpty(s string) bool {
	if s == "" {
		return true
	}
	return whitespaceRe.MatchString(s)
}

// Left returns the left substring of length n.
func Left(s string, n int) string {
	if n < 0 {
		return Right(s, -n)
	}
	return Substr(s, 0, n)
}

// LeftF is the filter form of Left.
func LeftF(n int) func(string) string {
	return func(s string) string {
		return Left(s, n)
	}
}

// LeftOf returns the substring left of needle.
func LeftOf(s string, needle string) string {
	return Between(s, "", needle)
}

// Letters returns an array of runes as strings so it can be indexed into.
func Letters(s string) []string {
	result := []string{}
	for _, r := range s {
		result = append(result, string(r))
	}
	return result
}

// Lines convert windows newlines to unix newlines then convert to an Array of lines.
func Lines(s string) []string {
	s = strings.Replace(s, "\r\n", "\n", -1)
	return strings.Split(s, "\n")
}

// Map maps an array's iitem through an iterator.
func Map(arr []string, iterator func(string) string) []string {
	r := []string{}
	for _, item := range arr {
		r = append(r, iterator(item))
	}
	return r
}

// Match returns true if patterns matches the string
func Match(s, pattern string) bool {
	r := regexp.MustCompile(pattern)
	return r.MatchString(s)
}
