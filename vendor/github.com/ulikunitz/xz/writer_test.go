// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xz

import (
	"bytes"
	"io"
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/ulikunitz/xz/internal/randtxt"
)

func TestWriter(t *testing.T) {
	const text = "The quick brown fox jumps over the lazy dog."
	var buf bytes.Buffer
	w, err := NewWriter(&buf)
	if err != nil {
		t.Fatalf("NewWriter error %s", err)
	}
	n, err := io.WriteString(w, text)
	if err != nil {
		t.Fatalf("WriteString error %s", err)
	}
	if n != len(text) {
		t.Fatalf("Writestring wrote %d bytes; want %d", n, len(text))
	}
	if err = w.Close(); err != nil {
		t.Fatalf("w.Close error %s", err)
	}
	var out bytes.Buffer
	r, err := NewReader(&buf)
	if err != nil {
		t.Fatalf("NewReader error %s", err)
	}
	if _, err = io.Copy(&out, r); err != nil {
		t.Fatalf("io.Copy error %s", err)
	}
	s := out.String()
	if s != text {
		t.Fatalf("reader decompressed to %q; want %q", s, text)
	}
}

func TestIssue12(t *testing.T) {
	var buf bytes.Buffer
	w, err := NewWriter(&buf)
	if err != nil {
		t.Fatalf("NewWriter error %s", err)
	}
	if err = w.Close(); err != nil {
		t.Fatalf("w.Close error %s", err)
	}
	r, err := NewReader(&buf)
	if err != nil {
		t.Fatalf("NewReader error %s", err)
	}
	var out bytes.Buffer
	if _, err = io.Copy(&out, r); err != nil {
		t.Fatalf("io.Copy error %s", err)
	}
	s := out.String()
	if s != "" {
		t.Fatalf("reader decompressed to %q; want %q", s, "")
	}
}

func Example() {
	const text = "The quick brown fox jumps over the lazy dog."
	var buf bytes.Buffer

	// compress text
	w, err := NewWriter(&buf)
	if err != nil {
		log.Fatalf("NewWriter error %s", err)
	}
	if _, err := io.WriteString(w, text); err != nil {
		log.Fatalf("WriteString error %s", err)
	}
	if err := w.Close(); err != nil {
		log.Fatalf("w.Close error %s", err)
	}

	// decompress buffer and write result to stdout
	r, err := NewReader(&buf)
	if err != nil {
		log.Fatalf("NewReader error %s", err)
	}
	if _, err = io.Copy(os.Stdout, r); err != nil {
		log.Fatalf("io.Copy error %s", err)
	}

	// Output:
	// The quick brown fox jumps over the lazy dog.
}

func TestWriter2(t *testing.T) {
	const txtlen = 1023
	var buf bytes.Buffer
	io.CopyN(&buf, randtxt.NewReader(rand.NewSource(41)), txtlen)
	txt := buf.String()

	buf.Reset()
	w, err := NewWriter(&buf)
	if err != nil {
		t.Fatalf("NewWriter error %s", err)
	}
	n, err := io.WriteString(w, txt)
	if err != nil {
		t.Fatalf("WriteString error %s", err)
	}
	if n != len(txt) {
		t.Fatalf("WriteString wrote %d bytes; want %d", n, len(txt))
	}
	if err = w.Close(); err != nil {
		t.Fatalf("Close error %s", err)
	}
	t.Logf("buf.Len() %d", buf.Len())
	r, err := NewReader(&buf)
	if err != nil {
		t.Fatalf("NewReader error %s", err)
	}
	var out bytes.Buffer
	k, err := io.Copy(&out, r)
	if err != nil {
		t.Fatalf("Decompressing copy error %s after %d bytes", err, n)
	}
	if k != txtlen {
		t.Fatalf("Decompression data length %d; want %d", k, txtlen)
	}
	if txt != out.String() {
		t.Fatal("decompressed data differs from original")
	}
}
