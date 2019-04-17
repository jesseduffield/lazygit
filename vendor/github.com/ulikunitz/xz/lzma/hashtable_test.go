// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lzma

import (
	"fmt"
	"testing"
)

func TestHashTable(t *testing.T) {
	ht, err := newHashTable(32, 2)
	if err != nil {
		t.Fatalf("newHashTable: error %s", err)
	}
	//    01234567890123456
	s := "abcabcdefghijklmn"
	n, err := ht.Write([]byte(s))
	if err != nil {
		t.Fatalf("ht.Write: error %s", err)
	}
	if n != len(s) {
		t.Fatalf("ht.Write returned %d; want %d", n, len(s))
	}
	tests := []struct {
		s string
		w string
	}{
		{"ab", "[3 0]"},
		{"bc", "[4 1]"},
		{"ca", "[2]"},
		{"xx", "[]"},
		{"gh", "[9]"},
		{"mn", "[15]"},
	}
	distances := make([]int64, 20)
	for _, c := range tests {
		distances := distances[:20]
		k := ht.Matches([]byte(c.s), distances)
		distances = distances[:k]
		o := fmt.Sprintf("%v", distances)
		if o != c.w {
			t.Errorf("%s: offsets %s; want %s", c.s, o, c.w)
		}
	}
}
