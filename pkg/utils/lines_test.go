package utils

import (
	"bufio"
	"strings"
	"testing"

	"github.com/jesseduffield/gocui"
	"github.com/stretchr/testify/assert"
)

// TestSplitLines is a function.
func TestSplitLines(t *testing.T) {
	type scenario struct {
		multilineString string
		expected        []string
	}

	scenarios := []scenario{
		{
			"",
			[]string{},
		},
		{
			"\n",
			[]string{},
		},
		{
			"hello world !\nhello universe !\n",
			[]string{
				"hello world !",
				"hello universe !",
			},
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, SplitLines(s.multilineString))
	}
}

func TestSplitNul(t *testing.T) {
	type scenario struct {
		multilineString string
		expected        []string
	}

	scenarios := []scenario{
		{
			"",
			[]string{},
		},
		{
			"\x00",
			[]string{
				"",
			},
		},
		{
			"hello world !\x00hello universe !\x00",
			[]string{
				"hello world !",
				"hello universe !",
			},
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, SplitNul(s.multilineString))
	}
}

// TestNormalizeLinefeeds is a function.
func TestNormalizeLinefeeds(t *testing.T) {
	type scenario struct {
		byteArray []byte
		expected  []byte
	}
	scenarios := []scenario{
		{
			// \r\n
			[]byte{97, 115, 100, 102, 13, 10},
			[]byte{97, 115, 100, 102, 10},
		},
		{
			// bash\r\nblah
			[]byte{97, 115, 100, 102, 13, 10, 97, 115, 100, 102},
			[]byte{97, 115, 100, 102, 10, 97, 115, 100, 102},
		},
		{
			// \r
			[]byte{97, 115, 100, 102, 13},
			[]byte{97, 115, 100, 102},
		},
		{
			// \n
			[]byte{97, 115, 100, 102, 10},
			[]byte{97, 115, 100, 102, 10},
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, string(s.expected), NormalizeLinefeeds(string(s.byteArray)))
	}
}

func TestScanLinesAndTruncateWhenLongerThanBuffer(t *testing.T) {
	type scenario struct {
		input         string
		expectedLines []string
	}

	scenarios := []scenario{
		{
			"",
			[]string{},
		},
		{
			"\n",
			[]string{""},
		},
		{
			"abc",
			[]string{"abc"},
		},
		{
			"abc\ndef",
			[]string{"abc", "def"},
		},
		{
			"abc\n\ndef",
			[]string{"abc", "", "def"},
		},
		{
			"abc\r\ndef\r",
			[]string{"abc", "def"},
		},
		{
			"abcdef",
			[]string{"abcde"},
		},
		{
			"abcdef\n",
			[]string{"abcde"},
		},
		{
			"abcdef\nghijkl\nx",
			[]string{"abcde", "ghijk", "x"},
		},
		{
			"abc\ndefghijklmnopqrstuvw\nx",
			[]string{"abc", "defgh", "x"},
		},
	}

	for _, s := range scenarios {
		scanner := bufio.NewScanner(strings.NewReader(s.input))
		scanner.Buffer(make([]byte, 5), 5)
		scanner.Split(ScanLinesAndTruncateWhenLongerThanBuffer(5))
		result := []string{}
		for scanner.Scan() {
			result = append(result, scanner.Text())
		}
		assert.NoError(t, scanner.Err())
		assert.EqualValues(t, s.expectedLines, result)
	}
}

func TestWrapViewLinesToWidth(t *testing.T) {
	tests := []struct {
		name                         string
		wrap                         bool
		text                         string
		width                        int
		expectedWrappedLines         []string
		expectedWrappedLinesIndices  []int
		expectedOriginalLinesIndices []int
	}{
		{
			name:  "Wrap off",
			wrap:  false,
			text:  "1st line\n2nd line\n3rd line",
			width: 5,
			expectedWrappedLines: []string{
				"1st line",
				"2nd line",
				"3rd line",
			},
			expectedWrappedLinesIndices:  []int{0, 1, 2},
			expectedOriginalLinesIndices: []int{0, 1, 2},
		},
		{
			name:  "Wrap on space",
			wrap:  true,
			text:  "Hello World",
			width: 5,
			expectedWrappedLines: []string{
				"Hello",
				"World",
			},
			expectedWrappedLinesIndices:  []int{0},
			expectedOriginalLinesIndices: []int{0, 0},
		},
		{
			name:  "Wrap on hyphen",
			wrap:  true,
			text:  "Hello-World",
			width: 6,
			expectedWrappedLines: []string{
				"Hello-",
				"World",
			},
		},
		{
			name:  "Wrap on hyphen 2",
			wrap:  true,
			text:  "Blah Hello-World",
			width: 12,
			expectedWrappedLines: []string{
				"Blah Hello-",
				"World",
			},
		},
		{
			name:  "Wrap on hyphen 3",
			wrap:  true,
			text:  "Blah Hello-World",
			width: 11,
			expectedWrappedLines: []string{
				"Blah Hello-",
				"World",
			},
		},
		{
			name:  "Wrap on hyphen 4",
			wrap:  true,
			text:  "Blah Hello-World",
			width: 10,
			expectedWrappedLines: []string{
				"Blah Hello",
				"-World",
			},
		},
		{
			name:  "Wrap on space 2",
			wrap:  true,
			text:  "Blah Hello World",
			width: 10,
			expectedWrappedLines: []string{
				"Blah Hello",
				"World",
			},
		},
		{
			name:  "Wrap on space with more words",
			wrap:  true,
			text:  "Longer word here",
			width: 10,
			expectedWrappedLines: []string{
				"Longer",
				"word here",
			},
		},
		{
			name:  "Split word that's too long",
			wrap:  true,
			text:  "ThisWordIsWayTooLong",
			width: 10,
			expectedWrappedLines: []string{
				"ThisWordIs",
				"WayTooLong",
			},
		},
		{
			name:  "Split word that's too long over multiple lines",
			wrap:  true,
			text:  "ThisWordIsWayTooLong",
			width: 5,
			expectedWrappedLines: []string{
				"ThisW",
				"ordIs",
				"WayTo",
				"oLong",
			},
		},
		{
			name:  "Lots of hyphens",
			wrap:  true,
			text:  "one-two-three-four-five",
			width: 8,
			expectedWrappedLines: []string{
				"one-two-",
				"three-",
				"four-",
				"five",
			},
		},
		{
			name:  "Several lines using all the available width",
			wrap:  true,
			text:  "aaa bb cc ddd-ee ff",
			width: 5,
			expectedWrappedLines: []string{
				"aaa",
				"bb cc",
				"ddd-",
				"ee ff",
			},
		},
		{
			name:  "Several lines using all the available width, with multi-cell runes",
			wrap:  true,
			text:  "🐤🐤🐤 🐝🐝 🙉🙉 🦊🦊🦊-🐬🐬 🦢🦢",
			width: 9,
			expectedWrappedLines: []string{
				"🐤🐤🐤",
				"🐝🐝 🙉🙉",
				"🦊🦊🦊-",
				"🐬🐬 🦢🦢",
			},
		},
		{
			name:  "Space in last column",
			wrap:  true,
			text:  "hello world",
			width: 6,
			expectedWrappedLines: []string{
				"hello",
				"world",
			},
		},
		{
			name:  "Hyphen in last column",
			wrap:  true,
			text:  "hello-world",
			width: 6,
			expectedWrappedLines: []string{
				"hello-",
				"world",
			},
		},
		{
			name:  "English text",
			wrap:  true,
			text:  "+The sea reach of the Thames stretched before us like the bedinnind of an interminable waterway. In the offind the sea and the sky were welded todether without a joint, and in the luminous space the tanned sails of the bardes drifting blah blah",
			width: 81,
			expectedWrappedLines: []string{
				"+The sea reach of the Thames stretched before us like the bedinnind of an",
				"interminable waterway. In the offind the sea and the sky were welded todether",
				"without a joint, and in the luminous space the tanned sails of the bardes",
				"drifting blah blah",
			},
		},
		{
			name:  "Tabs",
			wrap:  true,
			text:  "\ta\tbb\tccc\tdddd\teeeee",
			width: 50,
			expectedWrappedLines: []string{
				"    a   bb  ccc dddd    eeeee",
			},
		},
		{
			name:  "Multiple lines",
			wrap:  true,
			text:  "First paragraph\nThe second paragraph is a bit longer.\nThird paragraph\n",
			width: 10,
			expectedWrappedLines: []string{
				"First",
				"paragraph",
				"The second",
				"paragraph",
				"is a bit",
				"longer.",
				"Third",
				"paragraph",
			},
			expectedWrappedLinesIndices:  []int{0, 2, 6},
			expectedOriginalLinesIndices: []int{0, 0, 1, 1, 1, 1, 2, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrappedLines, wrappedLinesIndices, originalLinesIndices := WrapViewLinesToWidth(tt.wrap, tt.text, tt.width)
			assert.Equal(t, tt.expectedWrappedLines, wrappedLines)
			if tt.expectedWrappedLinesIndices != nil {
				assert.Equal(t, tt.expectedWrappedLinesIndices, wrappedLinesIndices)
			}
			if tt.expectedOriginalLinesIndices != nil {
				assert.Equal(t, tt.expectedOriginalLinesIndices, originalLinesIndices)
			}

			// As a sanity check, also test that gocui's line wrapping behaves the same way
			view := gocui.NewView("", 0, 0, tt.width+1, 1000, gocui.OutputNormal)
			assert.Equal(t, tt.width, view.InnerWidth())
			view.Wrap = tt.wrap
			view.SetContent(tt.text)
			assert.Equal(t, wrappedLines, view.ViewBufferLines())
		})
	}
}
