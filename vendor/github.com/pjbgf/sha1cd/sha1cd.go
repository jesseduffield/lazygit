// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sha1cd implements collision detection based on the whitepaper
// Counter-cryptanalysis from Marc Stevens. The original ubc implementation
// was done by Marc Stevens and Dan Shumow, and can be found at:
// https://github.com/cr-marcstevens/sha1collisiondetection
package sha1cd

// This SHA1 implementation is based on Go's generic SHA1.
// Original: https://github.com/golang/go/blob/master/src/crypto/sha1/sha1.go

import (
	"crypto"
	"encoding/binary"
	"errors"
	"hash"

	shared "github.com/pjbgf/sha1cd/internal"
)

//go:generate go run -C asm . -out ../sha1cdblock_amd64.s -pkg $GOPACKAGE

func init() {
	crypto.RegisterHash(crypto.SHA1, New)
}

// The size of a SHA-1 checksum in bytes.
const Size = shared.Size

// The blocksize of SHA-1 in bytes.
const BlockSize = shared.Chunk

// digest represents the partial evaluation of a checksum.
type digest struct {
	h   [shared.WordBuffers]uint32
	x   [shared.Chunk]byte
	nx  int
	len uint64

	// col defines whether a collision has been found.
	col       bool
	blockFunc func(dig *digest, p []byte)
}

func (d *digest) MarshalBinary() ([]byte, error) {
	b := make([]byte, 0, shared.MarshaledSize)
	b = append(b, shared.Magic...)
	b = appendUint32(b, d.h[0])
	b = appendUint32(b, d.h[1])
	b = appendUint32(b, d.h[2])
	b = appendUint32(b, d.h[3])
	b = appendUint32(b, d.h[4])
	b = append(b, d.x[:d.nx]...)
	b = b[:len(b)+len(d.x)-d.nx] // already zero
	b = appendUint64(b, d.len)
	return b, nil
}

func appendUint32(b []byte, v uint32) []byte {
	return append(b,
		byte(v>>24),
		byte(v>>16),
		byte(v>>8),
		byte(v),
	)
}

func appendUint64(b []byte, v uint64) []byte {
	return append(b,
		byte(v>>56),
		byte(v>>48),
		byte(v>>40),
		byte(v>>32),
		byte(v>>24),
		byte(v>>16),
		byte(v>>8),
		byte(v),
	)
}

func (d *digest) UnmarshalBinary(b []byte) error {
	if len(b) < len(shared.Magic) || string(b[:len(shared.Magic)]) != shared.Magic {
		return errors.New("crypto/sha1: invalid hash state identifier")
	}
	if len(b) != shared.MarshaledSize {
		return errors.New("crypto/sha1: invalid hash state size")
	}
	b = b[len(shared.Magic):]
	b, d.h[0] = consumeUint32(b)
	b, d.h[1] = consumeUint32(b)
	b, d.h[2] = consumeUint32(b)
	b, d.h[3] = consumeUint32(b)
	b, d.h[4] = consumeUint32(b)
	b = b[copy(d.x[:], b):]
	b, d.len = consumeUint64(b)
	d.nx = int(d.len % shared.Chunk)
	return nil
}

func consumeUint64(b []byte) ([]byte, uint64) {
	_ = b[7]
	x := uint64(b[7]) | uint64(b[6])<<8 | uint64(b[shared.WordBuffers])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
	return b[8:], x
}

func consumeUint32(b []byte) ([]byte, uint32) {
	_ = b[3]
	x := uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
	return b[4:], x
}

func (d *digest) Reset() {
	d.h[0] = shared.Init0
	d.h[1] = shared.Init1
	d.h[2] = shared.Init2
	d.h[3] = shared.Init3
	d.h[4] = shared.Init4
	d.nx = 0
	d.len = 0

	d.col = false
}

// New returns a new hash.Hash computing the SHA1 checksum. The Hash also
// implements encoding.BinaryMarshaler and encoding.BinaryUnmarshaler to
// marshal and unmarshal the internal state of the hash.
func New() hash.Hash {
	d := new(digest)

	d.blockFunc = block
	d.Reset()
	return d
}

// NewGeneric is equivalent to New but uses the Go generic implementation,
// avoiding any processor-specific optimizations.
func NewGeneric() hash.Hash {
	d := new(digest)

	d.blockFunc = blockGeneric
	d.Reset()
	return d
}

func (d *digest) Size() int { return Size }

func (d *digest) BlockSize() int { return BlockSize }

func (d *digest) Write(p []byte) (nn int, err error) {
	if len(p) == 0 {
		return
	}

	nn = len(p)
	d.len += uint64(nn)
	if d.nx > 0 {
		n := copy(d.x[d.nx:], p)
		d.nx += n
		if d.nx == shared.Chunk {
			d.blockFunc(d, d.x[:])
			d.nx = 0
		}
		p = p[n:]
	}
	if len(p) >= shared.Chunk {
		n := len(p) &^ (shared.Chunk - 1)
		d.blockFunc(d, p[:n])
		p = p[n:]
	}
	if len(p) > 0 {
		d.nx = copy(d.x[:], p)
	}
	return
}

func (d *digest) Sum(in []byte) []byte {
	// Make a copy of d so that caller can keep writing and summing.
	d0 := *d
	hash := d0.checkSum()
	return append(in, hash[:]...)
}

func (d *digest) checkSum() [Size]byte {
	len := d.len
	// Padding.  Add a 1 bit and 0 bits until 56 bytes mod 64.
	var tmp [64]byte
	tmp[0] = 0x80
	if len%64 < 56 {
		d.Write(tmp[0 : 56-len%64])
	} else {
		d.Write(tmp[0 : 64+56-len%64])
	}

	// Length in bits.
	len <<= 3
	binary.BigEndian.PutUint64(tmp[:], len)
	d.Write(tmp[0:8])

	if d.nx != 0 {
		panic("d.nx != 0")
	}

	var digest [Size]byte

	binary.BigEndian.PutUint32(digest[0:], d.h[0])
	binary.BigEndian.PutUint32(digest[4:], d.h[1])
	binary.BigEndian.PutUint32(digest[8:], d.h[2])
	binary.BigEndian.PutUint32(digest[12:], d.h[3])
	binary.BigEndian.PutUint32(digest[16:], d.h[4])

	return digest
}

// Sum returns the SHA-1 checksum of the data.
func Sum(data []byte) ([Size]byte, bool) {
	d := New().(*digest)
	d.Write(data)
	return d.checkSum(), d.col
}

func (d *digest) CollisionResistantSum(in []byte) ([]byte, bool) {
	// Make a copy of d so that caller can keep writing and summing.
	d0 := *d
	hash := d0.checkSum()
	return append(in, hash[:]...), d0.col
}
