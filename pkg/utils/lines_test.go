package utils

import (
	"bufio"
	"strings"
	"testing"

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
