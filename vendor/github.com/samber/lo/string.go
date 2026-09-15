package lo

import (
	"math"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/samber/lo/internal/xrand"
)

var (
	//nolint:revive
	LowerCaseLettersCharset = []rune("abcdefghijklmnopqrstuvwxyz")
	UpperCaseLettersCharset = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	LettersCharset          = append(LowerCaseLettersCharset, UpperCaseLettersCharset...)
	NumbersCharset          = []rune("0123456789")
	AlphanumericCharset     = append(LettersCharset, NumbersCharset...)
	SpecialCharset          = []rune("!@#$%^&*()_+-=[]{}|;':\",./<>?")
	AllCharset              = append(AlphanumericCharset, SpecialCharset...)

	// bearer:disable go_lang_permissive_regex_validation
	splitWordReg = regexp.MustCompile(`([a-z])([A-Z0-9])|([a-zA-Z])([0-9])|([0-9])([a-zA-Z])|([A-Z])([A-Z])([a-z])`)
	// bearer:disable go_lang_permissive_regex_validation
	splitNumberLetterReg = regexp.MustCompile(`([0-9])([a-zA-Z])`)
	maximumCapacity      = math.MaxInt>>1 + 1
)

// RandomString return a random string.
// Play: https://go.dev/play/p/rRseOQVVum4
func RandomString(size int, charset []rune) string {
	if size <= 0 {
		panic("lo.RandomString: size must be greater than 0")
	}
	if len(charset) == 0 {
		panic("lo.RandomString: charset must not be empty")
	}

	// see https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	var sb strings.Builder
	sb.Grow(size)

	if len(charset) == 1 {
		// Edge case, because if the charset is a single character,
		// it will panic below (divide by zero).
		// -> https://github.com/samber/lo/issues/679
		for i := 0; i < size; i++ {
			sb.WriteRune(charset[0])
		}
		return sb.String()
	}

	// Calculate the number of bits required to represent the charset,
	// e.g., for 62 characters, it would need 6 bits (since 62 -> 64 = 2^6)
	letterIDBits := int(math.Log2(float64(nearestPowerOfTwo(len(charset)))))
	// Determine the corresponding bitmask,
	// e.g., for 62 characters, the bitmask would be 111111.
	var letterIDMask int64 = 1<<letterIDBits - 1
	// Available count, since xrand.Int64() returns a non-negative number, the first bit is fixed, so there are 63 random bits
	// e.g., for 62 characters, this value is 10 (63 / 6).
	letterIDMax := 63 / letterIDBits
	// Generate the random string in a loop.
	for i, cache, remain := size-1, xrand.Int64(), letterIDMax; i >= 0; {
		// Regenerate the random number if all available bits have been used
		if remain == 0 {
			cache, remain = xrand.Int64(), letterIDMax
		}
		// Select a character from the charset
		if idx := int(cache & letterIDMask); idx < len(charset) {
			sb.WriteRune(charset[idx])
			i--
		}
		// Shift the bits to the right to prepare for the next character selection,
		// e.g., for 62 characters, shift by 6 bits.
		cache >>= letterIDBits
		// Decrease the remaining number of uses for the current random number.
		remain--
	}
	return sb.String()
}

// nearestPowerOfTwo returns the nearest power of two.
func nearestPowerOfTwo(capacity int) int {
	n := capacity - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	if n < 0 {
		return 1
	}
	if n >= maximumCapacity {
		return maximumCapacity
	}
	return n + 1
}

// Substring extracts a substring from a string with Unicode character (rune) awareness.
// offset - starting position of the substring (can be positive, negative, or zero)
// length - number of characters to extract
// With positive offset, counting starts from the beginning of the string
// With negative offset, counting starts from the end of the string
// Play: https://go.dev/play/p/TQlxQi82Lu1
func Substring[T ~string](str T, offset int, length uint) T {
	str = substring(str, offset, length)

	// Validate UTF-8 and fix invalid sequences
	if !utf8.ValidString(string(str)) {
		// Convert to []rune to replicate behavior with duplicated �
		str = T([]rune(str))
	}

	// Remove null bytes from result
	return T(strings.ReplaceAll(string(str), "\x00", ""))
}

func substring[T ~string](str T, offset int, length uint) T {
	switch {
	// Empty length or offset beyond string bounds - return empty string
	case length == 0, offset >= len(str):
		return ""

	// Positive offset - count from the beginning
	case offset > 0:
		// Skip offset runes from the start
		for i, r := range str {
			if offset--; offset == 0 {
				str = str[i+utf8.RuneLen(r):]
				break
			}
		}

		// If couldn't skip enough runes - string is shorter than offset
		if offset != 0 {
			return ""
		}

		// If remaining string is shorter than or equal to length - return it entirely
		if uint(len(str)) <= length {
			return str
		}

		// Otherwise proceed to trimming by length
		fallthrough

	// Zero offset or offset less than minus string length - start from beginning
	case offset < -len(str), offset == 0:
		// Count length runes from the start
		for i := range str {
			if length == 0 {
				return str[:i]
			}
			length--
		}

		return str

	// Negative offset - count from the end of string
	default: // -len(str) < offset < 0
		// Helper function to move backward through runes
		backwardPos := func(end int, count uint) (start int) {
			for {
				_, i := utf8.DecodeLastRuneInString(string(str[:end]))
				end -= i

				if count--; count == 0 || end == 0 {
					return end
				}
			}
		}

		offset := uint(-offset)

		// If offset is less than or equal to length - take from position to end
		if offset <= length {
			start := backwardPos(len(str), offset)
			return str[start:]
		}

		// Otherwise calculate start and end positions
		end := backwardPos(len(str), offset-length)
		start := backwardPos(end, length)

		return str[start:end]
	}
}

// ChunkString returns a slice of strings split into groups of length size. If the string can't be split evenly,
// the final chunk will be the remaining characters.
// Play: https://go.dev/play/p/__FLTuJVz54
//
// Note: lo.ChunkString and lo.Chunk functions behave inconsistently for empty input: lo.ChunkString("", n) returns [""] instead of [].
// See https://github.com/samber/lo/issues/788
func ChunkString[T ~string](str T, size int) []T {
	if size <= 0 {
		panic("lo.ChunkString: size must be greater than 0")
	}

	if size >= len(str) {
		return []T{str}
	}

	chunks := make([]T, 0, ((len(str)-1)/size)+1)
	currentLen := 0
	currentStart := 0
	for i := range str {
		if currentLen == size {
			chunks = append(chunks, str[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, str[currentStart:])
	return chunks
}

// RuneLength is an alias to utf8.RuneCountInString which returns the number of runes in string.
// Play: https://go.dev/play/p/tuhgW_lWY8l
func RuneLength(str string) int {
	return utf8.RuneCountInString(str)
}

// PascalCase converts string to pascal case.
// Play: https://go.dev/play/p/Dy_V_6DUYhe
func PascalCase(str string) string {
	items := Words(str)
	for i := range items {
		items[i] = Capitalize(items[i])
	}
	return strings.Join(items, "")
}

// CamelCase converts string to camel case.
// Play: https://go.dev/play/p/Go6aKwUiq59
func CamelCase(str string) string {
	items := Words(str)
	for i, item := range items {
		item = strings.ToLower(item)
		if i > 0 {
			item = Capitalize(item)
		}
		items[i] = item
	}
	return strings.Join(items, "")
}

// KebabCase converts string to kebab case.
// Play: https://go.dev/play/p/96gT_WZnTVP
func KebabCase(str string) string {
	items := Words(str)
	for i := range items {
		items[i] = strings.ToLower(items[i])
	}
	return strings.Join(items, "-")
}

// SnakeCase converts string to snake case.
// Play: https://go.dev/play/p/ziB0V89IeVH
func SnakeCase(str string) string {
	items := Words(str)
	for i := range items {
		items[i] = strings.ToLower(items[i])
	}
	return strings.Join(items, "_")
}

// Words splits string into a slice of its words.
// Play: https://go.dev/play/p/-f3VIQqiaVw
func Words(str string) []string {
	str = splitWordReg.ReplaceAllString(str, `$1$3$5$7 $2$4$6$8$9`)
	// example: Int8Value => Int 8Value => Int 8 Value
	str = splitNumberLetterReg.ReplaceAllString(str, "$1 $2")
	var result strings.Builder
	result.Grow(len(str))
	for _, r := range str {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		} else {
			result.WriteRune(' ')
		}
	}
	return strings.Fields(result.String())
}

// Capitalize converts the first character of string to upper case and the remaining to lower case.
// Play: https://go.dev/play/p/uLTZZQXqnsa
func Capitalize(str string) string {
	return cases.Title(language.English).String(str)
}

// Ellipsis trims and truncates a string to a specified length in runes and appends an ellipsis
// if truncated. The length parameter counts Unicode code points (runes), not bytes, so multi-byte
// characters such as emoji or CJK ideographs are never split in the middle.
// Play: https://go.dev/play/p/qE93rgqe1TW
func Ellipsis(str string, length int) string {
	str = strings.TrimSpace(str)

	const ellipsis = "..."

	cutPosition := 0
	for i := range str {
		if length == len(ellipsis) {
			cutPosition = i
		}

		if length--; length < 0 {
			return strings.TrimSpace(str[:cutPosition]) + ellipsis
		}
	}

	return str
}
