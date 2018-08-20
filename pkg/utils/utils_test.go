package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestWithPadding(t *testing.T) {
	type scenario struct {
		str      string
		padding  int
		expected string
	}

	scenarios := []scenario{
		{
			"hello world !",
			1,
			"hello world !",
		},
		{
			"hello world !",
			14,
			"hello world ! ",
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, WithPadding(s.str, s.padding))
	}
}

func TestTrimTrailingNewline(t *testing.T) {
	type scenario struct {
		str      string
		expected string
	}

	scenarios := []scenario{
		{
			"hello world !\n",
			"hello world !",
		},
		{
			"hello world !",
			"hello world !",
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, TrimTrailingNewline(s.str))
	}
}

var testCases = []struct {
	Input    []byte
	Expected []byte
}{
	{
		// \r\n
		Input:    []byte{97, 115, 100, 102, 13, 10},
		Expected: []byte{97, 115, 100, 102},
	},
	{
		// \r
		Input:    []byte{97, 115, 100, 102, 13},
		Expected: []byte{97, 115, 100, 102},
	},
	{
		// \n
		Input:    []byte{97, 115, 100, 102, 10},
		Expected: []byte{97, 115, 100, 102, 10},
	},
}

func TestNormalizeLinefeeds(t *testing.T) {
	for _, tc := range testCases {
		input := NormalizeLinefeeds(string(tc.Input))
		expected := string(tc.Expected)
		if input != expected {
			t.Error("Expected " + expected + ", got " + input)
		}
}
