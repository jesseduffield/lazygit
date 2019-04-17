// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xz

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestReaderSimple(t *testing.T) {
	const file = "fox.xz"
	xz, err := os.Open(file)
	if err != nil {
		t.Fatalf("os.Open(%q) error %s", file, err)
	}
	r, err := NewReader(xz)
	if err != nil {
		t.Fatalf("NewReader error %s", err)
	}
	var buf bytes.Buffer
	if _, err = io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy error %s", err)
	}
}

func TestReaderSingleStream(t *testing.T) {
	data, err := ioutil.ReadFile("fox.xz")
	if err != nil {
		t.Fatalf("ReadFile error %s", err)
	}
	xz := bytes.NewReader(data)
	rc := ReaderConfig{SingleStream: true}
	r, err := rc.NewReader(xz)
	if err != nil {
		t.Fatalf("NewReader error %s", err)
	}
	var buf bytes.Buffer
	if _, err = io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy error %s", err)
	}
	buf.Reset()
	data = append(data, 0)
	xz = bytes.NewReader(data)
	r, err = rc.NewReader(xz)
	if err != nil {
		t.Fatalf("NewReader error %s", err)
	}
	if _, err = io.Copy(&buf, r); err != errUnexpectedData {
		t.Fatalf("io.Copy returned %v; want %v", err, errUnexpectedData)
	}
}

func TestReaaderMultipleStreams(t *testing.T) {
	data, err := ioutil.ReadFile("fox.xz")
	if err != nil {
		t.Fatalf("ReadFile error %s", err)
	}
	m := make([]byte, 0, 4*len(data)+4*4)
	m = append(m, data...)
	m = append(m, data...)
	m = append(m, 0, 0, 0, 0)
	m = append(m, data...)
	m = append(m, 0, 0, 0, 0)
	m = append(m, 0, 0, 0, 0)
	m = append(m, data...)
	m = append(m, 0, 0, 0, 0)
	xz := bytes.NewReader(m)
	r, err := NewReader(xz)
	if err != nil {
		t.Fatalf("NewReader error %s", err)
	}
	var buf bytes.Buffer
	if _, err = io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy error %s", err)
	}
}
