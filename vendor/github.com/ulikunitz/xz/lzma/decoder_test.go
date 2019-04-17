// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lzma

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestDecoder(t *testing.T) {
	filename := "fox.lzma"
	want := "The quick brown fox jumps over the lazy dog.\n"
	for i := 0; i < 2; i++ {
		f, err := os.Open(filename)
		if err != nil {
			t.Fatalf("os.Open(%q) error %s", filename, err)
		}
		p := make([]byte, 13)
		_, err = io.ReadFull(f, p)
		if err != nil {
			t.Fatalf("io.ReadFull error %s", err)
		}
		props, err := PropertiesForCode(p[0])
		if err != nil {
			t.Fatalf("p[0] error %s", err)
		}
		state := newState(props)
		const capacity = 0x800000
		dict, err := newDecoderDict(capacity)
		if err != nil {
			t.Fatalf("newDecoderDict: error %s", err)
		}
		size := int64(-1)
		if i > 0 {
			size = int64(len(want))
		}
		br := bufio.NewReader(f)
		r, err := newDecoder(br, state, dict, size)
		if err != nil {
			t.Fatalf("newDecoder error %s", err)
		}
		bytes, err := ioutil.ReadAll(r)
		if err != nil {
			t.Fatalf("[%d] ReadAll error %s", i, err)
		}
		if err = f.Close(); err != nil {
			t.Fatalf("Close error %s", err)
		}
		got := string(bytes)
		if got != want {
			t.Fatalf("read %q; but want %q", got, want)
		}
	}
}
