// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package curve25519

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"
)

const expectedHex = "89161fde887b2b53de549af483940106ecc114d6982daa98256de23bdf77661a"

func TestBaseScalarMult(t *testing.T) {
	var a, b [32]byte
	in := &a
	out := &b
	a[0] = 1

	for i := 0; i < 200; i++ {
		ScalarBaseMult(out, in)
		in, out = out, in
	}

	result := fmt.Sprintf("%x", in[:])
	if result != expectedHex {
		t.Errorf("incorrect result: got %s, want %s", result, expectedHex)
	}
}

func TestTestVectors(t *testing.T) {
	for _, tv := range testVectors {
		var got [32]byte
		ScalarMult(&got, &tv.In, &tv.Base)
		if !bytes.Equal(got[:], tv.Expect[:]) {
			t.Logf("    in = %x", tv.In)
			t.Logf("  base = %x", tv.Base)
			t.Logf("   got = %x", got)
			t.Logf("expect = %x", tv.Expect)
			t.Fail()
		}
	}
}

// TestHighBitIgnored tests the following requirement in RFC 7748:
//
//	When receiving such an array, implementations of X25519 (but not X448) MUST
//	mask the most significant bit in the final byte.
//
// Regression test for issue #30095.
func TestHighBitIgnored(t *testing.T) {
	var s, u [32]byte
	rand.Read(s[:])
	rand.Read(u[:])

	var hi0, hi1 [32]byte

	u[31] &= 0x7f
	ScalarMult(&hi0, &s, &u)

	u[31] |= 0x80
	ScalarMult(&hi1, &s, &u)

	if !bytes.Equal(hi0[:], hi1[:]) {
		t.Errorf("high bit of group point should not affect result")
	}
}

func BenchmarkScalarBaseMult(b *testing.B) {
	var in, out [32]byte
	in[0] = 1

	b.SetBytes(32)
	for i := 0; i < b.N; i++ {
		ScalarBaseMult(&out, &in)
	}
}
