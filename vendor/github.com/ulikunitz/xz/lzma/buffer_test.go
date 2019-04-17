// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lzma

import (
	"bytes"
	"io"
	"testing"
)

func TestBuffer_Write(t *testing.T) {
	buf := newBuffer(10)
	b := []byte("1234567890")
	for i := range b {
		n, err := buf.Write(b[i : i+1])
		if err != nil {
			t.Fatalf("buf.Write(b[%d:%d]) error %s", i, i+1, err)
		}
		if n != 1 {
			t.Fatalf("buf.Write(b[%d:%d]) returned %d; want %d",
				i, i+1, n, 1)
		}
	}
	const c = 8
	n, err := buf.Discard(c)
	if err != nil {
		t.Fatalf("Discard error %s", err)
	}
	if n != c {
		t.Fatalf("Discard returned %d; want %d", n, c)
	}
	n, err = buf.Write(b)
	if err == nil {
		t.Fatalf("Write length exceed returned no error; n %d", n)
	}
	if n != c {
		t.Fatalf("Write length exceeding returned %d; want %d", n, c)
	}
	n, err = buf.Discard(4)
	if err != nil {
		t.Fatalf("Discard error %s", err)
	}
	if n != 4 {
		t.Fatalf("Discard returned %d; want %d", n, 4)
	}
	n, err = buf.Write(b[:3])
	if err != nil {
		t.Fatalf("buf.Write(b[:3]) error %s; n %d", err, n)
	}
	if n != 3 {
		t.Fatalf("buf.Write(b[:3]) returned %d; want %d", n, 3)
	}
}

func TestBuffer_Buffered_Available(t *testing.T) {
	buf := newBuffer(19)
	b := []byte("0123456789")
	var err error
	if _, err = buf.Write(b); err != nil {
		t.Fatalf("buf.Write(b) error %s", err)
	}
	if n := buf.Buffered(); n != 10 {
		t.Fatalf("buf.Buffered() returns %d; want %d", n, 10)
	}
	if _, err = buf.Discard(8); err != nil {
		t.Fatalf("buf.Discard(8) error %s", err)
	}
	if _, err = buf.Write(b[:7]); err != nil {
		t.Fatalf("buf.Write(b[:7]) error %s", err)
	}
	if n := buf.Buffered(); n != 9 {
		t.Fatalf("buf.Buffered() returns %d; want %d", n, 9)
	}
}

func TestBuffer_Read(t *testing.T) {
	buf := newBuffer(10)
	b := []byte("0123456789")
	var err error
	if _, err = buf.Write(b); err != nil {
		t.Fatalf("buf.Write(b) error %s", err)
	}
	p := make([]byte, 8)
	n, err := buf.Read(p)
	if err != nil {
		t.Fatalf("buf.Read(p) error %s", err)
	}
	if n != len(p) {
		t.Fatalf("buf.Read(p) returned %d; want %d", n, len(p))
	}
	if !bytes.Equal(p, b[:8]) {
		t.Fatalf("buf.Read(p) put %s into p; want %s", p, b[:8])
	}
	if _, err = buf.Write(b[:7]); err != nil {
		t.Fatalf("buf.Write(b[:7]) error %s", err)
	}
	q := make([]byte, 7)
	n, err = buf.Read(q)
	if err != nil {
		t.Fatalf("buf.Read(q) error %s", err)
	}
	if n != len(q) {
		t.Fatalf("buf.Read(q) returns %d; want %d", n, len(q))
	}
	c := []byte("8901234")
	if !bytes.Equal(q, c) {
		t.Fatalf("buf.Read(q) put %s into q; want %s", q, c)
	}
	if _, err := buf.Write(b[7:]); err != nil {
		t.Fatalf("buf.Write(b[7:]) error %s", err)
	}
	if _, err := buf.Write(b[:2]); err != nil {
		t.Fatalf("buf.Write(b[:2]) error %s", err)
	}
	t.Logf("buf.rear %d buf.front %d", buf.rear, buf.front)
	r := make([]byte, 2)
	n, err = buf.Read(r)
	if err != nil {
		t.Fatalf("buf.Read(r) error %s", err)
	}
	if n != len(r) {
		t.Fatalf("buf.Read(r) returns %d; want %d", n, len(r))
	}
	d := []byte("56")
	if !bytes.Equal(r, d) {
		t.Fatalf("buf.Read(r) put %s into r; want %s", r, d)
	}
}

func TestBuffer_Discard(t *testing.T) {
	buf := newBuffer(10)
	b := []byte("0123456789")
	var err error
	if _, err = buf.Write(b); err != nil {
		t.Fatalf("buf.Write(b) error %s", err)
	}
	n, err := buf.Discard(11)
	if err == nil {
		t.Fatalf("buf.Discard(11) didn't return error")
	}
	if n != 10 {
		t.Fatalf("buf.Discard(11) returned %d; want %d", n, 10)
	}
	if _, err := buf.Write(b); err != nil {
		t.Fatalf("buf.Write(b) #2 error %s", err)
	}
	n, err = buf.Discard(10)
	if err != nil {
		t.Fatalf("buf.Discard(10) error %s", err)
	}
	if n != 10 {
		t.Fatalf("buf.Discard(11) returned %d; want %d", n, 10)
	}
	if _, err := buf.Write(b[:4]); err != nil {
		t.Fatalf("buf.Write(b[:4]) error %s", err)
	}
	n, err = buf.Discard(1)
	if err != nil {
		t.Fatalf("buf.Discard(1) error %s", err)
	}
	if n != 1 {
		t.Fatalf("buf.Discard(1) returned %d; want %d", n, 1)
	}
}

func TestBuffer_Discard_error(t *testing.T) {
	buf := newBuffer(10)
	n, err := buf.Discard(-1)
	if err == nil {
		t.Fatal("buf.Discard(-1) didn't return an error")
	}
	if n != 0 {
		t.Fatalf("buf.Discard(-1) returned %d; want %d", n, 0)
	}
}

func TestPrefixLen(t *testing.T) {
	tests := []struct {
		a, b []byte
		k    int
	}{
		{[]byte("abcde"), []byte("abc"), 3},
		{[]byte("abc"), []byte("uvw"), 0},
		{[]byte(""), []byte("uvw"), 0},
		{[]byte("abcde"), []byte("abcuvw"), 3},
	}
	for _, c := range tests {
		k := prefixLen(c.a, c.b)
		if k != c.k {
			t.Errorf("prefixLen(%q,%q) returned %d; want %d",
				c.a, c.b, k, c.k)
		}
		k = prefixLen(c.b, c.a)
		if k != c.k {
			t.Errorf("prefixLen(%q,%q) returned %d; want %d",
				c.b, c.a, k, c.k)
		}
	}
}

func TestMatchLen(t *testing.T) {
	buf := newBuffer(13)
	const s = "abcaba"
	_, err := io.WriteString(buf, s)
	if err != nil {
		t.Fatalf("WriteString error %s", err)
	}
	_, err = io.WriteString(buf, s)
	if err != nil {
		t.Fatalf("WriteString error %s", err)
	}
	if _, err = buf.Discard(12); err != nil {
		t.Fatalf("buf.Discard(6) error %s", err)
	}
	_, err = io.WriteString(buf, s)
	if err != nil {
		t.Fatalf("WriteString error %s", err)
	}
	tests := []struct{ d, n int }{{1, 1}, {3, 2}, {6, 6}, {5, 0}, {2, 0}}
	for _, c := range tests {
		n := buf.matchLen(c.d, []byte(s))
		if n != c.n {
			t.Errorf(
				"MatchLen(%d,[]byte(%q)) returned %d; want %d",
				c.d, s, n, c.n)
		}
	}
}
