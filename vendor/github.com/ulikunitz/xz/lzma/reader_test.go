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
	"os"
	"path/filepath"
	"testing"
	"testing/iotest"
)

func TestNewReader(t *testing.T) {
	f, err := os.Open("examples/a.lzma")
	if err != nil {
		t.Fatalf("open examples/a.lzma: %s", err)
	}
	defer f.Close()
	_, err = NewReader(bufio.NewReader(f))
	if err != nil {
		t.Fatalf("NewReader: %s", err)
	}
}

const (
	dirname  = "examples"
	origname = "a.txt"
)

func readOrigFile(t *testing.T) []byte {
	orig, err := ioutil.ReadFile(filepath.Join(dirname, origname))
	if err != nil {
		t.Fatalf("ReadFile: %s", err)
	}
	return orig
}

func testDecodeFile(t *testing.T, filename string, orig []byte) {
	pathname := filepath.Join(dirname, filename)
	f, err := os.Open(pathname)
	if err != nil {
		t.Fatalf("Open(%q): %s", pathname, err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			t.Fatalf("f.Close() error %s", err)
		}
	}()
	t.Logf("file %s opened", filename)
	l, err := NewReader(bufio.NewReader(f))
	if err != nil {
		t.Fatalf("NewReader: %s", err)
	}
	decoded, err := ioutil.ReadAll(l)
	if err != nil {
		t.Fatalf("ReadAll: %s", err)
	}
	t.Logf("%s", decoded)
	if len(orig) != len(decoded) {
		t.Fatalf("length decoded is %d; want %d",
			len(decoded), len(orig))
	}
	if !bytes.Equal(orig, decoded) {
		t.Fatalf("decoded file differs from original")
	}
}

func TestReaderSimple(t *testing.T) {
	// DebugOn(os.Stderr)
	// defer DebugOff()

	testDecodeFile(t, "a.lzma", readOrigFile(t))
}

func TestReaderAll(t *testing.T) {
	dirname := "examples"
	dir, err := os.Open(dirname)
	if err != nil {
		t.Fatalf("Open: %s", err)
	}
	defer func() {
		if err := dir.Close(); err != nil {
			t.Fatalf("dir.Close() error %s", err)
		}
	}()
	all, err := dir.Readdirnames(0)
	if err != nil {
		t.Fatalf("Readdirnames: %s", err)
	}
	// filter now all file with the pattern "a*.lzma"
	files := make([]string, 0, len(all))
	for _, fn := range all {
		match, err := filepath.Match("a*.lzma", fn)
		if err != nil {
			t.Fatalf("Match: %s", err)
		}
		if match {
			files = append(files, fn)
		}
	}
	t.Log("files:", files)
	orig := readOrigFile(t)
	// actually test the files
	for _, fn := range files {
		testDecodeFile(t, fn, orig)
	}
}

//
func Example_reader() {
	f, err := os.Open("fox.lzma")
	if err != nil {
		log.Fatal(err)
	}
	// no need for defer; Fatal calls os.Exit(1) that doesn't execute deferred functions
	r, err := NewReader(bufio.NewReader(f))
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(os.Stdout, r)
	if err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	// Output:
	// The quick brown fox jumps over the lazy dog.
}

type wrapTest struct {
	name string
	wrap func(io.Reader) io.Reader
}

func (w *wrapTest) testFile(t *testing.T, filename string, orig []byte) {
	pathname := filepath.Join(dirname, filename)
	f, err := os.Open(pathname)
	if err != nil {
		t.Fatalf("Open(\"%s\"): %s", pathname, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	t.Logf("%s file %s opened", w.name, filename)
	l, err := NewReader(w.wrap(f))
	if err != nil {
		t.Fatalf("%s NewReader: %s", w.name, err)
	}
	decoded, err := ioutil.ReadAll(l)
	if err != nil {
		t.Fatalf("%s ReadAll: %s", w.name, err)
	}
	t.Logf("%s", decoded)
	if len(orig) != len(decoded) {
		t.Fatalf("%s length decoded is %d; want %d",
			w.name, len(decoded), len(orig))
	}
	if !bytes.Equal(orig, decoded) {
		t.Fatalf("%s decoded file differs from original", w.name)
	}
}

func TestReaderWrap(t *testing.T) {
	tests := [...]wrapTest{
		{"DataErrReader", iotest.DataErrReader},
		{"HalfReader", iotest.HalfReader},
		{"OneByteReader", iotest.OneByteReader},
		// TimeOutReader would require buffer
	}
	orig := readOrigFile(t)
	for _, tst := range tests {
		tst.testFile(t, "a.lzma", orig)
	}
}

func TestReaderBadFiles(t *testing.T) {
	dirname := "examples"
	dir, err := os.Open(dirname)
	if err != nil {
		t.Fatalf("Open: %s", err)
	}
	defer func() {
		if err := dir.Close(); err != nil {
			t.Fatalf("dir.Close() error %s", err)
		}
	}()
	all, err := dir.Readdirnames(0)
	if err != nil {
		t.Fatalf("Readdirnames: %s", err)
	}
	// filter now all file with the pattern "bad*.lzma"
	files := make([]string, 0, len(all))
	for _, fn := range all {
		match, err := filepath.Match("bad*.lzma", fn)
		if err != nil {
			t.Fatalf("Match: %s", err)
		}
		if match {
			files = append(files, fn)
		}
	}
	t.Log("files:", files)
	for _, filename := range files {
		pathname := filepath.Join(dirname, filename)
		f, err := os.Open(pathname)
		if err != nil {
			t.Fatalf("Open(\"%s\"): %s", pathname, err)
		}
		defer func(f *os.File) {
			if err := f.Close(); err != nil {
				t.Fatalf("f.Close() error %s", err)
			}
		}(f)
		t.Logf("file %s opened", filename)
		l, err := NewReader(f)
		if err != nil {
			t.Fatalf("NewReader: %s", err)
		}
		decoded, err := ioutil.ReadAll(l)
		if err == nil {
			t.Errorf("ReadAll for %s: no error", filename)
			t.Logf("%s", decoded)
			continue
		}
		t.Logf("%s: error %s", filename, err)
	}
}

type repReader byte

func (r repReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = byte(r)
	}
	return len(p), nil
}

func newRepReader(c byte, n int64) *io.LimitedReader {
	return &io.LimitedReader{R: repReader(c), N: n}
}

func newCodeReader(r io.Reader) *io.PipeReader {
	pr, pw := io.Pipe()
	go func() {
		bw := bufio.NewWriter(pw)
		lw, err := NewWriter(bw)
		if err != nil {
			log.Fatalf("NewWriter error %s", err)
		}
		if _, err = io.Copy(lw, r); err != nil {
			log.Fatalf("io.Copy error %s", err)
		}
		if err = lw.Close(); err != nil {
			log.Fatalf("lw.Close error %s", err)
		}
		if err = bw.Flush(); err != nil {
			log.Fatalf("bw.Flush() error %s", err)
		}
		if err = pw.CloseWithError(io.EOF); err != nil {
			log.Fatalf("pw.CloseWithError(io.EOF) error %s", err)
		}
	}()
	return pr
}

func TestReaderErrAgain(t *testing.T) {
	lengths := []int64{0, 128, 1024, 4095, 4096, 4097, 8191, 8192, 8193}
	buf := make([]byte, 128)
	const c = 'A'
	for _, n := range lengths {
		t.Logf("n: %d", n)
		pr := newCodeReader(newRepReader(c, n))
		r, err := NewReader(pr)
		if err != nil {
			t.Fatalf("NewReader(pr) error %s", err)
		}
		k := int64(0)
		for {
			m, err := r.Read(buf)
			k += int64(m)
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Errorf("r.Read(buf) error %s", err)
				break
			}
			if m > len(buf) {
				t.Fatalf("r.Read(buf) %d; want <= %d", m,
					len(buf))
			}
			for i, b := range buf[:m] {
				if b != c {
					t.Fatalf("buf[%d]=%c; want %c", i, b,
						c)
				}
			}
		}
		if k != n {
			t.Errorf("Read %d bytes; want %d", k, n)
		}
	}
}
