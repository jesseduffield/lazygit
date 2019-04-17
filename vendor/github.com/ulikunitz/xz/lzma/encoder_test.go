// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lzma

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"

	"github.com/ulikunitz/xz/internal/randtxt"
)

var testString = `LZMA decoder test example
=========================
! LZMA ! Decoder ! TEST !
=========================
! TEST ! LZMA ! Decoder !
=========================
---- Test Line 1 --------
=========================
---- Test Line 2 --------
=========================
=== End of test file ====
=========================
`

func cycle(t *testing.T, n int) {
	t.Logf("cycle(t,%d)", n)
	if n > len(testString) {
		t.Fatalf("cycle: n=%d larger than len(testString)=%d", n,
			len(testString))
	}
	const dictCap = MinDictCap
	m, err := newHashTable(dictCap, 4)
	if err != nil {
		t.Fatal(err)
	}
	encoderDict, err := newEncoderDict(dictCap, dictCap+1024, m)
	if err != nil {
		t.Fatal(err)
	}
	props := Properties{2, 0, 2}
	if err := props.verify(); err != nil {
		t.Fatalf("properties error %s", err)
	}
	state := newState(props)
	var buf bytes.Buffer
	w, err := newEncoder(&buf, state, encoderDict, eosMarker)
	if err != nil {
		t.Fatalf("newEncoder error %s", err)
	}
	orig := []byte(testString)[:n]
	t.Logf("len(orig) %d", len(orig))
	k, err := w.Write(orig)
	if err != nil {
		t.Fatalf("w.Write error %s", err)
	}
	if k != len(orig) {
		t.Fatalf("w.Write returned %d; want %d", k, len(orig))
	}
	if err = w.Close(); err != nil {
		t.Fatalf("w.Close error %s", err)
	}
	t.Logf("buf.Len() %d len(orig) %d", buf.Len(), len(orig))
	decoderDict, err := newDecoderDict(dictCap)
	if err != nil {
		t.Fatalf("newDecoderDict error %s", err)
	}
	state.Reset()
	r, err := newDecoder(&buf, state, decoderDict, -1)
	if err != nil {
		t.Fatalf("newDecoder error %s", err)
	}
	decoded, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ReadAll(lr) error %s", err)
	}
	t.Logf("decoded: %s", decoded)
	if len(orig) != len(decoded) {
		t.Fatalf("length decoded is %d; want %d", len(decoded),
			len(orig))
	}
	if !bytes.Equal(orig, decoded) {
		t.Fatalf("decoded file differs from original")
	}
}

func TestEncoderCycle1(t *testing.T) {
	cycle(t, len(testString))
}

func TestEncoderCycle2(t *testing.T) {
	buf := new(bytes.Buffer)
	const txtlen = 50000
	io.CopyN(buf, randtxt.NewReader(rand.NewSource(42)), txtlen)
	txt := buf.String()
	buf.Reset()
	const dictCap = MinDictCap
	m, err := newHashTable(dictCap, 4)
	if err != nil {
		t.Fatal(err)
	}
	encoderDict, err := newEncoderDict(dictCap, dictCap+1024, m)
	if err != nil {
		t.Fatal(err)
	}
	props := Properties{3, 0, 2}
	if err := props.verify(); err != nil {
		t.Fatalf("properties error %s", err)
	}
	state := newState(props)
	lbw := &LimitedByteWriter{BW: buf, N: 100}
	w, err := newEncoder(lbw, state, encoderDict, 0)
	if err != nil {
		t.Fatalf("NewEncoder error %s", err)
	}
	_, err = io.WriteString(w, txt)
	if err != nil && err != ErrLimit {
		t.Fatalf("WriteString error %s", err)
	}
	if err = w.Close(); err != nil {
		t.Fatalf("w.Close error %s", err)
	}
	n := w.Compressed()
	txt = txt[:n]
	decoderDict, err := newDecoderDict(dictCap)
	if err != nil {
		t.Fatalf("NewDecoderDict error %s", err)
	}
	state.Reset()
	r, err := newDecoder(buf, state, decoderDict, n)
	if err != nil {
		t.Fatalf("NewDecoder error %s", err)
	}
	out := new(bytes.Buffer)
	if _, err = io.Copy(out, r); err != nil {
		t.Fatalf("decompress copy error %s", err)
	}
	got := out.String()
	t.Logf("%s", got)
	if len(got) != int(n) {
		t.Fatalf("len(got) %d; want %d", len(got), n)
	}
	if got != txt {
		t.Fatalf("got and txt differ")
	}
}
