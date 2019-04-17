// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chacha20

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"
)

func TestCore(t *testing.T) {
	// This is just a smoke test that checks the example from
	// https://tools.ietf.org/html/rfc7539#section-2.3.2. The
	// chacha20poly1305 package contains much more extensive tests of this
	// code.
	var key [32]byte
	for i := range key {
		key[i] = byte(i)
	}

	var input [16]byte
	input[0] = 1
	input[7] = 9
	input[11] = 0x4a

	var out [64]byte
	XORKeyStream(out[:], out[:], &input, &key)
	const expected = "10f1e7e4d13b5915500fdd1fa32071c4c7d1f4c733c068030422aa9ac3d46c4ed2826446079faa0914c2d705d98b02a2b5129cd1de164eb9cbd083e8a2503c4e"
	if result := hex.EncodeToString(out[:]); result != expected {
		t.Errorf("wanted %x but got %x", expected, result)
	}
}

// Run the test cases with the input and output in different buffers.
func TestNoOverlap(t *testing.T) {
	for _, c := range testVectors {
		s := New(c.key, c.nonce)
		input, err := hex.DecodeString(c.input)
		if err != nil {
			t.Fatalf("cannot decode input %#v: %v", c.input, err)
		}
		output := make([]byte, c.length)
		s.XORKeyStream(output, input)
		got := hex.EncodeToString(output)
		if got != c.output {
			t.Errorf("length=%v: got %#v, want %#v", c.length, got, c.output)
		}
	}
}

// Run the test cases with the input and output overlapping entirely.
func TestOverlap(t *testing.T) {
	for _, c := range testVectors {
		s := New(c.key, c.nonce)
		data, err := hex.DecodeString(c.input)
		if err != nil {
			t.Fatalf("cannot decode input %#v: %v", c.input, err)
		}
		s.XORKeyStream(data, data)
		got := hex.EncodeToString(data)
		if got != c.output {
			t.Errorf("length=%v: got %#v, want %#v", c.length, got, c.output)
		}
	}
}

// Run the test cases with various source and destination offsets.
func TestUnaligned(t *testing.T) {
	const max = 8 // max offset (+1) to test
	for _, c := range testVectors {
		input := make([]byte, c.length+max)
		output := make([]byte, c.length+max)
		for i := 0; i < max; i++ { // input offsets
			for j := 0; j < max; j++ { // output offsets
				s := New(c.key, c.nonce)

				input := input[i : i+c.length]
				output := output[j : j+c.length]

				data, err := hex.DecodeString(c.input)
				if err != nil {
					t.Fatalf("cannot decode input %#v: %v", c.input, err)
				}
				copy(input, data)
				s.XORKeyStream(output, input)
				got := hex.EncodeToString(output)
				if got != c.output {
					t.Errorf("length=%v: got %#v, want %#v", c.length, got, c.output)
				}
			}
		}
	}
}

// Run the test cases by calling XORKeyStream multiple times.
func TestStep(t *testing.T) {
	// wide range of step sizes to try and hit edge cases
	steps := [...]int{1, 3, 4, 7, 8, 17, 24, 30, 64, 256}
	rnd := rand.New(rand.NewSource(123))
	for _, c := range testVectors {
		s := New(c.key, c.nonce)
		input, err := hex.DecodeString(c.input)
		if err != nil {
			t.Fatalf("cannot decode input %#v: %v", c.input, err)
		}
		output := make([]byte, c.length)

		// step through the buffers
		i, step := 0, steps[rnd.Intn(len(steps))]
		for i+step < c.length {
			s.XORKeyStream(output[i:i+step], input[i:i+step])
			if i+step < c.length && output[i+step] != 0 {
				t.Errorf("length=%v, i=%v, step=%v: output overwritten", c.length, i, step)
			}
			i += step
			step = steps[rnd.Intn(len(steps))]
		}
		// finish the encryption
		s.XORKeyStream(output[i:], input[i:])

		got := hex.EncodeToString(output)
		if got != c.output {
			t.Errorf("length=%v: got %#v, want %#v", c.length, got, c.output)
		}
	}
}

// Test that Advance() discards bytes until a block boundary is hit.
func TestAdvance(t *testing.T) {
	for _, c := range testVectors {
		for i := 0; i < 63; i++ {
			s := New(c.key, c.nonce)
			z := New(c.key, c.nonce)
			input, err := hex.DecodeString(c.input)
			if err != nil {
				t.Fatalf("cannot decode input %#v: %v", c.input, err)
			}
			zeros, discard := make([]byte, 64), make([]byte, 64)
			so, zo := make([]byte, c.length), make([]byte, c.length)
			for j := 0; j < c.length; j += 64 {
				lim := j + i
				if lim > c.length {
					lim = c.length
				}
				s.XORKeyStream(so[j:lim], input[j:lim])
				// calling s.Advance() multiple times should have no effect
				for k := 0; k < i%3+1; k++ {
					s.Advance()
				}
				z.XORKeyStream(zo[j:lim], input[j:lim])
				if lim < c.length {
					end := 64 - i
					if c.length-lim < end {
						end = c.length - lim
					}
					z.XORKeyStream(discard[:], zeros[:end])
				}
			}

			got := hex.EncodeToString(so)
			want := hex.EncodeToString(zo)
			if got != want {
				t.Errorf("length=%v: got %#v, want %#v", c.length, got, want)
			}
		}
	}
}

func BenchmarkChaCha20(b *testing.B) {
	sizes := []int{32, 63, 64, 256, 1024, 1350, 65536}
	for _, size := range sizes {
		s := size
		b.Run(fmt.Sprint(s), func(b *testing.B) {
			k := [32]byte{}
			c := [16]byte{}
			src := make([]byte, s)
			dst := make([]byte, s)
			b.SetBytes(int64(s))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				XORKeyStream(dst, src, &c, &k)
			}
		})
	}
}

func TestHChaCha20(t *testing.T) {
	// See draft-paragon-paseto-rfc-00 §7.2.1.
	key := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f}
	nonce := []byte{0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00, 0x4a,
		0x00, 0x00, 0x00, 0x00, 0x31, 0x41, 0x59, 0x27}
	expected := []byte{0x82, 0x41, 0x3b, 0x42, 0x27, 0xb2, 0x7b, 0xfe,
		0xd3, 0x0e, 0x42, 0x50, 0x8a, 0x87, 0x7d, 0x73,
		0xa0, 0xf9, 0xe4, 0xd5, 0x8a, 0x74, 0xa8, 0x53,
		0xc1, 0x2e, 0xc4, 0x13, 0x26, 0xd3, 0xec, 0xdc,
	}
	result := HChaCha20(&[8]uint32{
		binary.LittleEndian.Uint32(key[0:4]),
		binary.LittleEndian.Uint32(key[4:8]),
		binary.LittleEndian.Uint32(key[8:12]),
		binary.LittleEndian.Uint32(key[12:16]),
		binary.LittleEndian.Uint32(key[16:20]),
		binary.LittleEndian.Uint32(key[20:24]),
		binary.LittleEndian.Uint32(key[24:28]),
		binary.LittleEndian.Uint32(key[28:32]),
	}, &[4]uint32{
		binary.LittleEndian.Uint32(nonce[0:4]),
		binary.LittleEndian.Uint32(nonce[4:8]),
		binary.LittleEndian.Uint32(nonce[8:12]),
		binary.LittleEndian.Uint32(nonce[12:16]),
	})
	for i := 0; i < 8; i++ {
		want := binary.LittleEndian.Uint32(expected[i*4 : (i+1)*4])
		if got := result[i]; got != want {
			t.Errorf("word %d incorrect: want 0x%x, got 0x%x", i, want, got)
		}
	}
}
