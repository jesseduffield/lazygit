// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lzma

import (
	"bytes"
	"fmt"
	"testing"
)

func TestChunkTypeString(t *testing.T) {
	tests := [...]struct {
		c chunkType
		s string
	}{
		{cEOS, "EOS"},
		{cUD, "UD"},
		{cU, "U"},
		{cL, "L"},
		{cLR, "LR"},
		{cLRN, "LRN"},
		{cLRND, "LRND"},
	}
	for _, c := range tests {
		s := fmt.Sprintf("%v", c.c)
		if s != c.s {
			t.Errorf("got %s; want %s", s, c.s)
		}
	}
}

func TestHeaderChunkType(t *testing.T) {
	tests := []struct {
		h byte
		c chunkType
	}{
		{h: 0, c: cEOS},
		{h: 1, c: cUD},
		{h: 2, c: cU},
		{h: 1<<7 | 0x1f, c: cL},
		{h: 1<<7 | 1<<5 | 0x1f, c: cLR},
		{h: 1<<7 | 1<<6 | 0x1f, c: cLRN},
		{h: 1<<7 | 1<<6 | 1<<5 | 0x1f, c: cLRND},
		{h: 1<<7 | 1<<6 | 1<<5, c: cLRND},
	}
	if _, err := headerChunkType(3); err == nil {
		t.Fatalf("headerChunkType(%d) got %v; want %v",
			3, err, errHeaderByte)
	}
	for _, tc := range tests {
		c, err := headerChunkType(tc.h)
		if err != nil {
			t.Fatalf("headerChunkType error %s", err)
		}
		if c != tc.c {
			t.Errorf("got %s; want %s", c, tc.c)
		}
	}
}

func TestHeaderLen(t *testing.T) {
	tests := []struct {
		c chunkType
		n int
	}{
		{cEOS, 1}, {cU, 3}, {cUD, 3}, {cL, 5}, {cLR, 5}, {cLRN, 6},
		{cLRND, 6},
	}
	for _, tc := range tests {
		n := headerLen(tc.c)
		if n != tc.n {
			t.Errorf("header length for %s %d; want %d",
				tc.c, n, tc.n)
		}
	}
}

func chunkHeaderSamples(t *testing.T) []chunkHeader {
	props := Properties{LC: 3, LP: 0, PB: 2}
	headers := make([]chunkHeader, 0, 12)
	for c := cEOS; c <= cLRND; c++ {
		var h chunkHeader
		h.ctype = c
		if c >= cUD {
			h.uncompressed = 0x0304
		}
		if c >= cL {
			h.compressed = 0x0201
		}
		if c >= cLRN {
			h.props = props
		}
		headers = append(headers, h)
	}
	return headers
}

func TestChunkHeaderMarshalling(t *testing.T) {
	for _, h := range chunkHeaderSamples(t) {
		data, err := h.MarshalBinary()
		if err != nil {
			t.Fatalf("MarshalBinary for %v error %s", h, err)
		}
		var g chunkHeader
		if err = g.UnmarshalBinary(data); err != nil {
			t.Fatalf("UnmarshalBinary error %s", err)
		}
		if g != h {
			t.Fatalf("got %v; want %v", g, h)
		}
	}
}

func TestReadChunkHeader(t *testing.T) {
	for _, h := range chunkHeaderSamples(t) {
		data, err := h.MarshalBinary()
		if err != nil {
			t.Fatalf("MarshalBinary for %v error %s", h, err)
		}
		r := bytes.NewReader(data)
		g, err := readChunkHeader(r)
		if err != nil {
			t.Fatalf("readChunkHeader for %v error %s", h, err)
		}
		if *g != h {
			t.Fatalf("got %v; want %v", g, h)
		}
	}
}

func TestReadEOS(t *testing.T) {
	var b [1]byte
	r := bytes.NewReader(b[:])
	h, err := readChunkHeader(r)
	if err != nil {
		t.Fatalf("readChunkHeader error %s", err)
	}
	if h.ctype != cEOS {
		t.Errorf("ctype got %s; want %s", h.ctype, cEOS)
	}
	if h.compressed != 0 {
		t.Errorf("compressed got %d; want %d", h.compressed, 0)
	}
	if h.uncompressed != 0 {
		t.Errorf("uncompressed got %d; want %d", h.uncompressed, 0)
	}
	wantProps := Properties{}
	if h.props != wantProps {
		t.Errorf("props got %v; want %v", h.props, wantProps)
	}
}
