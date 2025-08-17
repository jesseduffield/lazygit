// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sha3

// This file defines the ShakeHash interface, and provides
// functions for creating SHAKE and cSHAKE instances, as well as utility
// functions for hashing bytes to arbitrary-length output.
//
//
// SHAKE implementation is based on FIPS PUB 202 [1]
// cSHAKE implementations is based on NIST SP 800-185 [2]
//
// [1] https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.202.pdf
// [2] https://doi.org/10.6028/NIST.SP.800-185

import (
	"io"
)

// ShakeHash defines the interface to hash functions that
// support arbitrary-length output.
type ShakeHash interface {
	// Write absorbs more data into the hash's state. It panics if input is
	// written to it after output has been read from it.
	io.Writer

	// Read reads more output from the hash; reading affects the hash's
	// state. (ShakeHash.Read is thus very different from Hash.Sum)
	// It never returns an error.
	io.Reader

	// Clone returns a copy of the ShakeHash in its current state.
	Clone() ShakeHash

	// Reset resets the ShakeHash to its initial state.
	Reset()
}

// Consts for configuring initial SHA-3 state
const (
	dsbyteShake = 0x1f
	rate128     = 168
	rate256     = 136
)

// Clone returns copy of SHAKE context within its current state.
func (d *State) Clone() ShakeHash {
	return d.clone()
}

// NewShake128 creates a new SHAKE128 variable-output-length ShakeHash.
// Its generic security strength is 128 bits against all attacks if at
// least 32 bytes of its output are used.
func NewShake128() State {
	return State{rate: rate128, dsbyte: dsbyteShake}
}

// NewTurboShake128 creates a new TurboSHAKE128 variable-output-length ShakeHash.
// Its generic security strength is 128 bits against all attacks if at
// least 32 bytes of its output are used.
// D is the domain separation byte and must be between 0x01 and 0x7f inclusive.
func NewTurboShake128(D byte) State {
	if D == 0 || D > 0x7f {
		panic("turboshake: D out of range")
	}
	return State{rate: rate128, dsbyte: D, turbo: true}
}

// NewShake256 creates a new SHAKE256 variable-output-length ShakeHash.
// Its generic security strength is 256 bits against all attacks if
// at least 64 bytes of its output are used.
func NewShake256() State {
	return State{rate: rate256, dsbyte: dsbyteShake}
}

// NewTurboShake256 creates a new TurboSHAKE256 variable-output-length ShakeHash.
// Its generic security strength is 256 bits against all attacks if
// at least 64 bytes of its output are used.
// D is the domain separation byte and must be between 0x01 and 0x7f inclusive.
func NewTurboShake256(D byte) State {
	if D == 0 || D > 0x7f {
		panic("turboshake: D out of range")
	}
	return State{rate: rate256, dsbyte: D, turbo: true}
}

// ShakeSum128 writes an arbitrary-length digest of data into hash.
func ShakeSum128(hash, data []byte) {
	h := NewShake128()
	_, _ = h.Write(data)
	_, _ = h.Read(hash)
}

// ShakeSum256 writes an arbitrary-length digest of data into hash.
func ShakeSum256(hash, data []byte) {
	h := NewShake256()
	_, _ = h.Write(data)
	_, _ = h.Read(hash)
}

// TurboShakeSum128 writes an arbitrary-length digest of data into hash.
func TurboShakeSum128(hash, data []byte, D byte) {
	h := NewTurboShake128(D)
	_, _ = h.Write(data)
	_, _ = h.Read(hash)
}

// TurboShakeSum256 writes an arbitrary-length digest of data into hash.
func TurboShakeSum256(hash, data []byte, D byte) {
	h := NewTurboShake256(D)
	_, _ = h.Write(data)
	_, _ = h.Read(hash)
}

func (d *State) SwitchDS(D byte) {
	d.dsbyte = D
}
