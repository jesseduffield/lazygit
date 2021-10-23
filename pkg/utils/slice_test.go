package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIncludesString is a function.
func TestIncludesString(t *testing.T) {
	type scenario struct {
		list     []string
		element  string
		expected bool
	}

	scenarios := []scenario{
		{
			[]string{"a", "b"},
			"a",
			true,
		},
		{
			[]string{"a", "b"},
			"c",
			false,
		},
		{
			[]string{"a", "b"},
			"",
			false,
		},
		{
			[]string{""},
			"",
			true,
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, IncludesString(s.list, s.element))
	}
}

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
		t.Run(s.testName, func(t *testing.T) {
			assert.EqualValues(t, s.expected, EscapeSpecialChars(s.input))
		})
	}
}

func TestUniq(t *testing.T) {
	for _, test := range []struct {
		values []string
		want   []string
	}{
		{
			values: []string{"a", "b", "c"},
			want:   []string{"a", "b", "c"},
		},
		{
			values: []string{"a", "b", "a", "b", "c"},
			want:   []string{"a", "b", "c"},
		},
	} {
		if got := Uniq(test.values); !assert.EqualValues(t, got, test.want) {
			t.Errorf("Uniq(%v) = %v; want %v", test.values, got, test.want)
		}
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

func TestReverse(t *testing.T) {
	for _, test := range []struct {
		values []string
		want   []string
	}{
		{
			values: []string{"a", "b", "c"},
			want:   []string{"c", "b", "a"},
		},
		{
			values: []string{},
			want:   []string{},
		},
	} {
		if got := Reverse(test.values); !assert.EqualValues(t, got, test.want) {
			t.Errorf("Reverse(%v) = %v; want %v", test.values, got, test.want)
		}
	}
}
