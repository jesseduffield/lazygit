// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lzma

import (
	"bytes"
	"io"
	"math/rand"
	"strings"
	"testing"

	"github.com/ulikunitz/xz/internal/randtxt"
)

func TestWriter2(t *testing.T) {
	var buf bytes.Buffer
	w, err := Writer2Config{DictCap: 4096}.NewWriter2(&buf)
	if err != nil {
		t.Fatalf("NewWriter error %s", err)
	}
	n, err := w.Write([]byte{'a'})
	if err != nil {
		t.Fatalf("w.Write([]byte{'a'}) error %s", err)
	}
	if n != 1 {
		t.Fatalf("w.Write([]byte{'a'}) returned %d; want %d", n, 1)
	}
	if err = w.Flush(); err != nil {
		t.Fatalf("w.Flush() error %s", err)
	}
	// check that double Flush doesn't write another chunk
	if err = w.Flush(); err != nil {
		t.Fatalf("w.Flush() error %s", err)
	}
	if err = w.Close(); err != nil {
		t.Fatalf("w.Close() error %s", err)
	}
	p := buf.Bytes()
	want := []byte{1, 0, 0, 'a', 0}
	if !bytes.Equal(p, want) {
		t.Fatalf("bytes written %#v; want %#v", p, want)
	}
}

func TestCycle1(t *testing.T) {
	var buf bytes.Buffer
	w, err := Writer2Config{DictCap: 4096}.NewWriter2(&buf)
	if err != nil {
		t.Fatalf("NewWriter error %s", err)
	}
	n, err := w.Write([]byte{'a'})
	if err != nil {
		t.Fatalf("w.Write([]byte{'a'}) error %s", err)
	}
	if n != 1 {
		t.Fatalf("w.Write([]byte{'a'}) returned %d; want %d", n, 1)
	}
	if err = w.Close(); err != nil {
		t.Fatalf("w.Close() error %s", err)
	}
	r, err := Reader2Config{DictCap: 4096}.NewReader2(&buf)
	if err != nil {
		t.Fatalf("NewReader error %s", err)
	}
	p := make([]byte, 3)
	n, err = r.Read(p)
	t.Logf("n %d error %v", n, err)
}

func TestCycle2(t *testing.T) {
	buf := new(bytes.Buffer)
	w, err := Writer2Config{DictCap: 4096}.NewWriter2(buf)
	if err != nil {
		t.Fatalf("NewWriter error %s", err)
	}
	// const txtlen = 1024
	const txtlen = 2100000
	io.CopyN(buf, randtxt.NewReader(rand.NewSource(42)), txtlen)
	txt := buf.String()
	buf.Reset()
	n, err := io.Copy(w, strings.NewReader(txt))
	if err != nil {
		t.Fatalf("Compressing copy error %s", err)
	}
	if n != txtlen {
		t.Fatalf("Compressing data length %d; want %d", n, txtlen)
	}
	if err = w.Close(); err != nil {
		t.Fatalf("w.Close error %s", err)
	}
	t.Logf("buf.Len() %d", buf.Len())
	r, err := Reader2Config{DictCap: 4096}.NewReader2(buf)
	if err != nil {
		t.Fatalf("NewReader error %s", err)
	}
	out := new(bytes.Buffer)
	n, err = io.Copy(out, r)
	if err != nil {
		t.Fatalf("Decompressing copy error %s after %d bytes", err, n)
	}
	if n != txtlen {
		t.Fatalf("Decompression data length %d; want %d", n, txtlen)
	}
	if txt != out.String() {
		t.Fatal("decompressed data differs from original")
	}
}
