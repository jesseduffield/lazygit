// Copyright (C) 2019 ProtonTech AG
// This file contains necessary tools for the aex and ocb packages.
//
// These functions SHOULD NOT be used elsewhere, since they are optimized for
// specific input nature in the EAX and OCB modes of operation.

package byteutil

// GfnDouble computes 2 * input in the field of 2^n elements.
// The irreducible polynomial in the finite field for n=128 is
// x^128 + x^7 + x^2 + x + 1 (equals 0x87)
// Constant-time execution in order to avoid side-channel attacks
func GfnDouble(input []byte) []byte {
	if len(input) != 16 {
		panic("Doubling in GFn only implemented for n = 128")
	}
	// If the first bit is zero, return 2L = L << 1
	// Else return (L << 1) xor 0^120 10000111
	shifted := ShiftBytesLeft(input)
	shifted[15] ^= ((input[0] >> 7) * 0x87)
	return shifted
}

// ShiftBytesLeft outputs the byte array corresponding to x << 1 in binary.
func ShiftBytesLeft(x []byte) []byte {
	l := len(x)
	dst := make([]byte, l)
	for i := 0; i < l-1; i++ {
		dst[i] = (x[i] << 1) | (x[i+1] >> 7)
	}
	dst[l-1] = x[l-1] << 1
	return dst
}

// ShiftNBytesLeft puts in dst the byte array corresponding to x << n in binary.
func ShiftNBytesLeft(dst, x []byte, n int) {
	// Erase first n / 8 bytes
	copy(dst, x[n/8:])

	// Shift the remaining n % 8 bits
	bits := uint(n % 8)
	l := len(dst)
	for i := 0; i < l-1; i++ {
		dst[i] = (dst[i] << bits) | (dst[i+1] >> uint(8-bits))
	}
	dst[l-1] = dst[l-1] << bits

	// Append trailing zeroes
	dst = append(dst, make([]byte, n/8)...)
}

// XorBytesMut replaces X with X XOR Y. len(X) must be >= len(Y).
func XorBytesMut(X, Y []byte) {
	for i := 0; i < len(Y); i++ {
		X[i] ^= Y[i]
	}
}

// XorBytes puts X XOR Y into Z. len(Z) and len(X) must be >= len(Y).
func XorBytes(Z, X, Y []byte) {
	for i := 0; i < len(Y); i++ {
		Z[i] = X[i] ^ Y[i]
	}
}

// RightXor XORs smaller input (assumed Y) at the right of the larger input (assumed X)
func RightXor(X, Y []byte) []byte {
	offset := len(X) - len(Y)
	xored := make([]byte, len(X))
	copy(xored, X)
	for i := 0; i < len(Y); i++ {
		xored[offset+i] ^= Y[i]
	}
	return xored
}

// SliceForAppend takes a slice and a requested number of bytes. It returns a
// slice with the contents of the given slice followed by that many bytes and a
// second slice that aliases into it and contains only the extra bytes. If the
// original slice has sufficient capacity then no allocation is performed.
func SliceForAppend(in []byte, n int) (head, tail []byte) {
	if total := len(in) + n; cap(in) >= total {
		head = in[:total]
	} else {
		head = make([]byte, total)
		copy(head, in)
	}
	tail = head[len(in):]
	return
}
