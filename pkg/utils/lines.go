package utils

import (
	"bytes"
	"strings"

	"github.com/mattn/go-runewidth"
)

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

func SplitNul(str string) []string {
	if str == "" {
		return make([]string, 0)
	}
	str = strings.TrimSuffix(str, "\x00")
	return strings.Split(str, "\x00")
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

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

// ScanLinesAndTruncateWhenLongerThanBuffer returns a split function that can be
// used with bufio.Scanner.Split(). It is very similar to bufio.ScanLines,
// except that it will truncate lines that are longer than the scanner's read
// buffer (whereas bufio.ScanLines will return an error in that case, which is
// often difficult to handle).
//
// If you are using your own buffer for the scanner, you must set maxBufferSize
// to the same value as the max parameter that you passed to scanner.Buffer().
// Otherwise, maxBufferSize must be set to bufio.MaxScanTokenSize.
func ScanLinesAndTruncateWhenLongerThanBuffer(maxBufferSize int) func(data []byte, atEOF bool) (int, []byte, error) {
	skipOverRemainderOfLongLine := false

	return func(data []byte, atEOF bool) (int, []byte, error) {
		if atEOF && len(data) == 0 {
			// Done
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			if skipOverRemainderOfLongLine {
				skipOverRemainderOfLongLine = false
				return i + 1, nil, nil
			}
			return i + 1, dropCR(data[0:i]), nil
		}
		if atEOF {
			if skipOverRemainderOfLongLine {
				return len(data), nil, nil
			}

			return len(data), dropCR(data), nil
		}

		// Buffer is full, so we can't get more data
		if len(data) >= maxBufferSize {
			if skipOverRemainderOfLongLine {
				return len(data), nil, nil
			}

			skipOverRemainderOfLongLine = true
			return len(data), data, nil
		}

		// Request more data.
		return 0, nil, nil
	}
}

// Wrap lines to a given width, and return:
// - the wrapped lines
// - the line indices of the wrapped lines, indexed by the original line indices
// - the line indices of the original lines, indexed by the wrapped line indices
// If wrap is false, the text is returned as is.
// This code needs to behave the same as `gocui.lineWrap` does.
func WrapViewLinesToWidth(wrap bool, text string, width int) ([]string, []int, []int) {
	text = strings.TrimSuffix(text, "\n")
	lines := strings.Split(text, "\n")
	if !wrap {
		indices := make([]int, len(lines))
		for i := range lines {
			indices[i] = i
		}
		return lines, indices, indices
	}

	wrappedLines := make([]string, 0, len(lines))
	wrappedLineIndices := make([]int, 0, len(lines))
	originalLineIndices := make([]int, 0, len(lines))

	for originalLineIdx, line := range lines {
		wrappedLineIndices = append(wrappedLineIndices, len(wrappedLines))

		// convert tabs to spaces
		for i := 0; i < len(line); i++ {
			if line[i] == '\t' {
				numSpaces := 4 - (i % 4)
				line = line[:i] + "    "[:numSpaces] + line[i+1:]
				i += numSpaces - 1
			}
		}

		appendWrappedLine := func(str string) {
			wrappedLines = append(wrappedLines, str)
			originalLineIndices = append(originalLineIndices, originalLineIdx)
		}

		n := 0
		offset := 0
		lastWhitespaceIndex := -1
		for i, currChr := range line {
			rw := runewidth.RuneWidth(currChr)
			n += rw

			if n > width {
				if currChr == ' ' {
					appendWrappedLine(line[offset:i])
					offset = i + 1
					n = 0
				} else if currChr == '-' {
					appendWrappedLine(line[offset:i])
					offset = i
					n = rw
				} else if lastWhitespaceIndex != -1 {
					if line[lastWhitespaceIndex] == '-' {
						appendWrappedLine(line[offset : lastWhitespaceIndex+1])
					} else {
						appendWrappedLine(line[offset:lastWhitespaceIndex])
					}
					offset = lastWhitespaceIndex + 1
					n = runewidth.StringWidth(line[offset : i+1])
				} else {
					appendWrappedLine(line[offset:i])
					offset = i
					n = rw
				}
				lastWhitespaceIndex = -1
			} else if currChr == ' ' || currChr == '-' {
				lastWhitespaceIndex = i
			}
		}

		appendWrappedLine(line[offset:])
	}

	return wrappedLines, wrappedLineIndices, originalLineIndices
}
