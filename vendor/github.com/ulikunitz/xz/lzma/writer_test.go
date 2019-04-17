// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lzma

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/ulikunitz/xz/internal/randtxt"
)

func TestWriterCycle(t *testing.T) {
	orig := readOrigFile(t)
	buf := new(bytes.Buffer)
	w, err := NewWriter(buf)
	if err != nil {
		t.Fatalf("NewWriter: error %s", err)
	}
	n, err := w.Write(orig)
	if err != nil {
		t.Fatalf("w.Write error %s", err)
	}
	if n != len(orig) {
		t.Fatalf("w.Write returned %d; want %d", n, len(orig))
	}
	if err = w.Close(); err != nil {
		t.Fatalf("w.Close error %s", err)
	}
	t.Logf("buf.Len() %d len(orig) %d", buf.Len(), len(orig))
	if buf.Len() > len(orig) {
		t.Errorf("buf.Len()=%d bigger then len(orig)=%d", buf.Len(),
			len(orig))
	}
	lr, err := NewReader(buf)
	if err != nil {
		t.Fatalf("NewReader error %s", err)
	}
	decoded, err := ioutil.ReadAll(lr)
	if err != nil {
		t.Fatalf("ReadAll(lr) error %s", err)
	}
	t.Logf("%s", decoded)
	if len(orig) != len(decoded) {
		t.Fatalf("length decoded is %d; want %d", len(decoded),
			len(orig))
	}
	if !bytes.Equal(orig, decoded) {
		t.Fatalf("decoded file differs from original")
	}
}

func TestWriterLongData(t *testing.T) {
	const (
		seed = 49
		size = 82237
	)
	r := io.LimitReader(randtxt.NewReader(rand.NewSource(seed)), size)
	txt, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ReadAll error %s", err)
	}
	if len(txt) != size {
		t.Fatalf("ReadAll read %d bytes; want %d", len(txt), size)
	}
	buf := &bytes.Buffer{}
	w, err := WriterConfig{DictCap: 0x4000}.NewWriter(buf)
	if err != nil {
		t.Fatalf("WriterConfig.NewWriter error %s", err)
	}
	n, err := w.Write(txt)
	if err != nil {
		t.Fatalf("w.Write error %s", err)
	}
	if n != len(txt) {
		t.Fatalf("w.Write wrote %d bytes; want %d", n, size)
	}
	if err = w.Close(); err != nil {
		t.Fatalf("w.Close error %s", err)
	}
	t.Logf("compressed length %d", buf.Len())
	lr, err := NewReader(buf)
	if err != nil {
		t.Fatalf("NewReader error %s", err)
	}
	txtRead, err := ioutil.ReadAll(lr)
	if err != nil {
		t.Fatalf("ReadAll(lr) error %s", err)
	}
	if len(txtRead) != size {
		t.Fatalf("ReadAll(lr) returned %d bytes; want %d",
			len(txtRead), size)
	}
	if !bytes.Equal(txtRead, txt) {
		t.Fatal("ReadAll(lr) returned txt differs from origin")
	}
}

func TestWriter_Size(t *testing.T) {
	buf := new(bytes.Buffer)
	w, err := WriterConfig{Size: 10, EOSMarker: true}.NewWriter(buf)
	if err != nil {
		t.Fatalf("WriterConfig.NewWriter error %s", err)
	}
	q := []byte{'a'}
	for i := 0; i < 9; i++ {
		n, err := w.Write(q)
		if err != nil {
			t.Fatalf("w.Write error %s", err)
		}
		if n != 1 {
			t.Fatalf("w.Write returned %d; want %d", n, 1)
		}
		q[0]++
	}
	if err := w.Close(); err != errSize {
		t.Fatalf("expected errSize, but got %v", err)
	}
	n, err := w.Write(q)
	if err != nil {
		t.Fatalf("w.Write error %s", err)
	}
	if n != 1 {
		t.Fatalf("w.Write returned %d; want %d", n, 1)
	}
	if err = w.Close(); err != nil {
		t.Fatalf("w.Close error %s", err)
	}
	t.Logf("compressed size %d", buf.Len())
	r, err := NewReader(buf)
	if err != nil {
		t.Fatalf("NewReader error %s", err)
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ReadAll error %s", err)
	}
	s := string(b)
	want := "abcdefghij"
	if s != want {
		t.Fatalf("read %q, want %q", s, want)
	}
}

// The example uses the buffered reader and writer from package bufio.
func Example_writer() {
	pr, pw := io.Pipe()
	go func() {
		bw := bufio.NewWriter(pw)
		w, err := NewWriter(bw)
		if err != nil {
			log.Fatal(err)
		}
		input := []byte("The quick brown fox jumps over the lazy dog.")
		if _, err = w.Write(input); err != nil {
			log.Fatal(err)
		}
		if err = w.Close(); err != nil {
			log.Fatal(err)
		}
		// reader waits for the data
		if err = bw.Flush(); err != nil {
			log.Fatal(err)
		}
	}()
	r, err := NewReader(pr)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(os.Stdout, r)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// The quick brown fox jumps over the lazy dog.
}

func BenchmarkReader(b *testing.B) {
	const (
		seed = 49
		size = 50000
	)
	r := io.LimitReader(randtxt.NewReader(rand.NewSource(seed)), size)
	txt, err := ioutil.ReadAll(r)
	if err != nil {
		b.Fatalf("ReadAll error %s", err)
	}
	buf := &bytes.Buffer{}
	w, err := WriterConfig{DictCap: 0x4000}.NewWriter(buf)
	if err != nil {
		b.Fatalf("WriterConfig{}.NewWriter error %s", err)
	}
	if _, err = w.Write(txt); err != nil {
		b.Fatalf("w.Write error %s", err)
	}
	if err = w.Close(); err != nil {
		b.Fatalf("w.Close error %s", err)
	}
	data, err := ioutil.ReadAll(buf)
	if err != nil {
		b.Fatalf("ReadAll error %s", err)
	}
	b.SetBytes(int64(len(txt)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lr, err := NewReader(bytes.NewReader(data))
		if err != nil {
			b.Fatalf("NewReader error %s", err)
		}
		if _, err = ioutil.ReadAll(lr); err != nil {
			b.Fatalf("ReadAll(lr) error %s", err)
		}
	}
}

func BenchmarkWriter(b *testing.B) {
	const (
		seed = 49
		size = 50000
	)
	r := io.LimitReader(randtxt.NewReader(rand.NewSource(seed)), size)
	txt, err := ioutil.ReadAll(r)
	if err != nil {
		b.Fatalf("ReadAll error %s", err)
	}
	buf := &bytes.Buffer{}
	b.SetBytes(int64(len(txt)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		w, err := WriterConfig{DictCap: 0x4000}.NewWriter(buf)
		if err != nil {
			b.Fatalf("NewWriter error %s", err)
		}
		if _, err = w.Write(txt); err != nil {
			b.Fatalf("w.Write error %s", err)
		}
		if err = w.Close(); err != nil {
			b.Fatalf("w.Close error %s", err)
		}
	}
}
