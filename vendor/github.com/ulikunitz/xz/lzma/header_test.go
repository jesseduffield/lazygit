// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lzma

import "testing"

func TestHeaderMarshalling(t *testing.T) {
	tests := []header{
		{properties: Properties{3, 0, 2}, dictCap: 8 * 1024 * 1024,
			size: -1},
		{properties: Properties{4, 3, 3}, dictCap: 4096,
			size: 10},
	}
	for _, h := range tests {
		data, err := h.marshalBinary()
		if err != nil {
			t.Fatalf("marshalBinary error %s", err)
		}
		var g header
		if err = g.unmarshalBinary(data); err != nil {
			t.Fatalf("unmarshalBinary error %s", err)
		}
		if h != g {
			t.Errorf("got header %#v; want %#v", g, h)
		}
	}
}

func TestValidHeader(t *testing.T) {
	tests := []header{
		{properties: Properties{3, 0, 2}, dictCap: 8 * 1024 * 1024,
			size: -1},
		{properties: Properties{4, 3, 3}, dictCap: 4096,
			size: 10},
	}
	for _, h := range tests {
		data, err := h.marshalBinary()
		if err != nil {
			t.Fatalf("marshalBinary error %s", err)
		}
		if !ValidHeader(data) {
			t.Errorf("ValidHeader returns false for header %v;"+
				" want true", h)
		}
	}
	const a = "1234567890123"
	if ValidHeader([]byte(a)) {
		t.Errorf("ValidHeader returns true for %s; want false", a)
	}
}
