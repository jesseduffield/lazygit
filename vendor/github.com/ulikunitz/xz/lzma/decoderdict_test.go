// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lzma

import (
	"fmt"
	"testing"
)

func peek(d *decoderDict) []byte {
	p := make([]byte, d.buffered())
	k, err := d.peek(p)
	if err != nil {
		panic(fmt.Errorf("peek: "+
			"Read returned unexpected error %s", err))
	}
	if k != len(p) {
		panic(fmt.Errorf("peek: "+
			"Read returned %d; wanted %d", k, len(p)))
	}
	return p
}

func TestNewDecoderDict(t *testing.T) {
	if _, err := newDecoderDict(0); err == nil {
		t.Fatalf("no error for zero dictionary capacity")
	}
	if _, err := newDecoderDict(8); err != nil {
		t.Fatalf("error %s", err)
	}
}
