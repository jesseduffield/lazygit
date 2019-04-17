// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

import (
	"testing"
	"time"
)

func TestIntComparator(t *testing.T) {

	// i1,i2,expected
	tests := [][]interface{}{
		{1, 1, 0},
		{1, 2, -1},
		{2, 1, 1},
		{11, 22, -1},
		{0, 0, 0},
		{1, 0, 1},
		{0, 1, -1},
	}

	for _, test := range tests {
		actual := IntComparator(test[0], test[1])
		expected := test[2]
		if actual != expected {
			t.Errorf("Got %v expected %v", actual, expected)
		}
	}
}

func TestStringComparator(t *testing.T) {

	// s1,s2,expected
	tests := [][]interface{}{
		{"a", "a", 0},
		{"a", "b", -1},
		{"b", "a", 1},
		{"aa", "aab", -1},
		{"", "", 0},
		{"a", "", 1},
		{"", "a", -1},
		{"", "aaaaaaa", -1},
	}

	for _, test := range tests {
		actual := StringComparator(test[0], test[1])
		expected := test[2]
		if actual != expected {
			t.Errorf("Got %v expected %v", actual, expected)
		}
	}
}

func TestTimeComparator(t *testing.T) {

	now := time.Now()

	// i1,i2,expected
	tests := [][]interface{}{
		{now, now, 0},
		{now.Add(24 * 7 * 2 * time.Hour), now, 1},
		{now, now.Add(24 * 7 * 2 * time.Hour), -1},
	}

	for _, test := range tests {
		actual := TimeComparator(test[0], test[1])
		expected := test[2]
		if actual != expected {
			t.Errorf("Got %v expected %v", actual, expected)
		}
	}
}

func TestCustomComparator(t *testing.T) {

	type Custom struct {
		id   int
		name string
	}

	byID := func(a, b interface{}) int {
		c1 := a.(Custom)
		c2 := b.(Custom)
		switch {
		case c1.id > c2.id:
			return 1
		case c1.id < c2.id:
			return -1
		default:
			return 0
		}
	}

	// o1,o2,expected
	tests := [][]interface{}{
		{Custom{1, "a"}, Custom{1, "a"}, 0},
		{Custom{1, "a"}, Custom{2, "b"}, -1},
		{Custom{2, "b"}, Custom{1, "a"}, 1},
		{Custom{1, "a"}, Custom{1, "b"}, 0},
	}

	for _, test := range tests {
		actual := byID(test[0], test[1])
		expected := test[2]
		if actual != expected {
			t.Errorf("Got %v expected %v", actual, expected)
		}
	}
}
