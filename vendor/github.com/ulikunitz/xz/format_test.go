// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xz

import (
	"bytes"
	"testing"
)

func TestHeader(t *testing.T) {
	h := header{flags: CRC32}
	data, err := h.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary error %s", err)
	}
	var g header
	if err = g.UnmarshalBinary(data); err != nil {
		t.Fatalf("UnmarshalBinary error %s", err)
	}
	if g != h {
		t.Fatalf("unmarshalled %#v; want %#v", g, h)
	}
}

func TestFooter(t *testing.T) {
	f := footer{indexSize: 64, flags: CRC32}
	data, err := f.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary error %s", err)
	}
	var g footer
	if err = g.UnmarshalBinary(data); err != nil {
		t.Fatalf("UnmarshalBinary error %s", err)
	}
	if g != f {
		t.Fatalf("unmarshalled %#v; want %#v", g, f)
	}
}

func TestRecord(t *testing.T) {
	r := record{1234567, 10000}
	p, err := r.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary error %s", err)
	}
	n := len(p)
	buf := bytes.NewReader(p)
	g, m, err := readRecord(buf)
	if err != nil {
		t.Fatalf("readFrom error %s", err)
	}
	if m != n {
		t.Fatalf("read %d bytes; wrote %d", m, n)
	}
	if g.unpaddedSize != r.unpaddedSize {
		t.Fatalf("got unpaddedSize %d; want %d", g.unpaddedSize,
			r.unpaddedSize)
	}
	if g.uncompressedSize != r.uncompressedSize {
		t.Fatalf("got uncompressedSize %d; want %d", g.uncompressedSize,
			r.uncompressedSize)
	}
}

func TestIndex(t *testing.T) {
	records := []record{{1234, 1}, {2345, 2}}

	var buf bytes.Buffer
	n, err := writeIndex(&buf, records)
	if err != nil {
		t.Fatalf("writeIndex error %s", err)
	}
	if n != int64(buf.Len()) {
		t.Fatalf("writeIndex returned %d; want %d", n, buf.Len())
	}

	// indicator
	c, err := buf.ReadByte()
	if err != nil {
		t.Fatalf("buf.ReadByte error %s", err)
	}
	if c != 0 {
		t.Fatalf("indicator %d; want %d", c, 0)
	}

	g, m, err := readIndexBody(&buf)
	if err != nil {
		for i, r := range g {
			t.Logf("records[%d] %v", i, r)
		}
		t.Fatalf("readIndexBody error %s", err)
	}
	if m != n-1 {
		t.Fatalf("readIndexBody returned %d; want %d", m, n-1)
	}
	for i, rec := range records {
		if g[i] != rec {
			t.Errorf("records[%d] is %v; want %v", i, g[i], rec)
		}
	}
}

func TestBlockHeader(t *testing.T) {
	h := blockHeader{
		compressedSize:   1234,
		uncompressedSize: -1,
		filters:          []filter{&lzmaFilter{4096}},
	}
	data, err := h.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary error %s", err)
	}

	r := bytes.NewReader(data)
	g, n, err := readBlockHeader(r)
	if err != nil {
		t.Fatalf("readBlockHeader error %s", err)
	}
	if n != len(data) {
		t.Fatalf("readBlockHeader returns %d bytes; want %d", n,
			len(data))
	}
	if g.compressedSize != h.compressedSize {
		t.Errorf("got compressedSize %d; want %d",
			g.compressedSize, h.compressedSize)
	}
	if g.uncompressedSize != h.uncompressedSize {
		t.Errorf("got uncompressedSize %d; want %d",
			g.uncompressedSize, h.uncompressedSize)
	}
	if len(g.filters) != len(h.filters) {
		t.Errorf("got len(filters) %d; want %d",
			len(g.filters), len(h.filters))
	}
	glf := g.filters[0].(*lzmaFilter)
	hlf := h.filters[0].(*lzmaFilter)
	if glf.dictCap != hlf.dictCap {
		t.Errorf("got dictCap %d; want %d", glf.dictCap, hlf.dictCap)
	}
}
