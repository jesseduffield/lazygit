package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextIndex(t *testing.T) {
	type scenario struct {
		testName string
		list     []int
		element  int
		expected int
	}

	scenarios := []scenario{
		{
			// I'm not really fussed about how it behaves here
			"no elements",
			[]int{},
			1,
			-1,
		},
		{
			"one element",
			[]int{1},
			1,
			0,
		},
		{
			"two elements",
			[]int{1, 2},
			1,
			1,
		},
		{
			"two elements, giving second one",
			[]int{1, 2},
			2,
			1,
		},
		{
			"three elements, giving second one",
			[]int{1, 2, 3},
			2,
			2,
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			assert.EqualValues(t, s.expected, NextIndex(s.list, s.element))
		})
	}
}

func TestPrevIndex(t *testing.T) {
	type scenario struct {
		testName string
		list     []int
		element  int
		expected int
	}

	scenarios := []scenario{
		{
			// I'm not really fussed about how it behaves here
			"no elements",
			[]int{},
			1,
			0,
		},
		{
			"one element",
			[]int{1},
			1,
			0,
		},
		{
			"two elements",
			[]int{1, 2},
			1,
			0,
		},
		{
			"three elements, giving second one",
			[]int{1, 2, 3},
			2,
			0,
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			assert.EqualValues(t, s.expected, PrevIndex(s.list, s.element))
		})
	}
}

func TestEscapeSpecialChars(t *testing.T) {
	type scenario struct {
		testName string
		input    string
		expected string
	}

	scenarios := []scenario{
		{
			"normal string",
			"ab",
			"ab",
		},
		{
			"string with a special char",
			"a\nb",
			"a\\nb",
		},
		{
			"multiple special chars",
			"\n\r\t\b\f\v",
			"\\n\\r\\t\\b\\f\\v",
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			assert.EqualValues(t, s.expected, EscapeSpecialChars(s.input))
		})
	}
}

func TestLimit(t *testing.T) {
	for _, test := range []struct {
		values []string
		limit  int
		want   []string
	}{
		{
			values: []string{"a", "b", "c"},
			limit:  3,
			want:   []string{"a", "b", "c"},
		},
		{
			values: []string{"a", "b", "c"},
			limit:  4,
			want:   []string{"a", "b", "c"},
		},
		{
			values: []string{"a", "b", "c"},
			limit:  2,
			want:   []string{"a", "b"},
		},
		{
			values: []string{"a", "b", "c"},
			limit:  1,
			want:   []string{"a"},
		},
		{
			values: []string{"a", "b", "c"},
			limit:  0,
			want:   []string{},
		},
		{
			values: []string{},
			limit:  0,
			want:   []string{},
		},
	} {
		if got := Limit(test.values, test.limit); !assert.EqualValues(t, got, test.want) {
			t.Errorf("Limit(%v, %d) = %v; want %v", test.values, test.limit, got, test.want)
		}
	}
}

func TestLimitStr(t *testing.T) {
	for _, test := range []struct {
		values string
		limit  int
		want   string
	}{
		{
			values: "",
			limit:  10,
			want:   "",
		},
		{
			values: "",
			limit:  0,
			want:   "",
		},
		{
			values: "a",
			limit:  1,
			want:   "a",
		},
		{
			values: "ab",
			limit:  2,
			want:   "ab",
		},
		{
			values: "abc",
			limit:  3,
			want:   "abc",
		},
		{
			values: "abcd",
			limit:  3,
			want:   "abc",
		},
		{
			values: "abcde",
			limit:  3,
			want:   "abc",
		},
		{
			values: "あいう",
			limit:  1,
			want:   "あ",
		},
		{
			values: "あいう",
			limit:  2,
			want:   "あい",
		},
	} {
		if got := LimitStr(test.values, test.limit); !assert.EqualValues(t, got, test.want) {
			t.Errorf("LimitString(%v, %d) = %v; want %v", test.values, test.limit, got, test.want)
		}
	}
}
