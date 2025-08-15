// Copyright (C) 2019 ProtonTech AG

// Package ocb provides an implementation of the OCB (offset codebook) mode of
// operation, as described in RFC-7253 of the IRTF and in Rogaway, Bellare,
// Black and Krovetz - OCB: A BLOCK-CIPHER MODE OF OPERATION FOR EFFICIENT
// AUTHENTICATED ENCRYPTION (2003).
// Security considerations (from RFC-7253): A private key MUST NOT be used to
// encrypt more than 2^48 blocks. Tag length should be at least 12 bytes (a
// brute-force forging adversary succeeds after 2^{tag length} attempts). A
// single key SHOULD NOT be used to decrypt ciphertext with different tag
// lengths. Nonces need not be secret, but MUST NOT be reused.
// This package only supports underlying block ciphers with 128-bit blocks,
// such as AES-{128, 192, 256}, but may be extended to other sizes.
package ocb

import (
	"bytes"
	"crypto/cipher"
	"crypto/subtle"
	"errors"
	"math/bits"

	"github.com/ProtonMail/go-crypto/internal/byteutil"
)

type ocb struct {
	block     cipher.Block
	tagSize   int
	nonceSize int
	mask      mask
	// Optimized en/decrypt: For each nonce N used to en/decrypt, the 'Ktop'
	// internal variable can be reused for en/decrypting with nonces sharing
	// all but the last 6 bits with N. The prefix of the first nonce used to
	// compute the new Ktop, and the Ktop value itself, are stored in
	// reusableKtop. If using incremental nonces, this saves one block cipher
	// call every 63 out of 64 OCB encryptions, and stores one nonce and one
	// output of the block cipher in memory only.
	reusableKtop reusableKtop
}

type mask struct {
	// L_*, L_$, (L_i)_{i ∈ N}
	lAst []byte
	lDol []byte
	L    [][]byte
}

type reusableKtop struct {
	noncePrefix []byte
	Ktop        []byte
}

const (
	defaultTagSize   = 16
	defaultNonceSize = 15
)

const (
	enc = iota
	dec
)

func (o *ocb) NonceSize() int {
	return o.nonceSize
}

func (o *ocb) Overhead() int {
	return o.tagSize
}

// NewOCB returns an OCB instance with the given block cipher and default
// tag and nonce sizes.
func NewOCB(block cipher.Block) (cipher.AEAD, error) {
	return NewOCBWithNonceAndTagSize(block, defaultNonceSize, defaultTagSize)
}

// NewOCBWithNonceAndTagSize returns an OCB instance with the given block
// cipher, nonce length, and tag length. Panics on zero nonceSize and
// exceedingly long tag size.
//
// It is recommended to use at least 12 bytes as tag length.
func NewOCBWithNonceAndTagSize(
	block cipher.Block, nonceSize, tagSize int) (cipher.AEAD, error) {
	if block.BlockSize() != 16 {
		return nil, ocbError("Block cipher must have 128-bit blocks")
	}
	if nonceSize < 1 {
		return nil, ocbError("Incorrect nonce length")
	}
	if nonceSize >= block.BlockSize() {
		return nil, ocbError("Nonce length exceeds blocksize - 1")
	}
	if tagSize > block.BlockSize() {
		return nil, ocbError("Custom tag length exceeds blocksize")
	}
	return &ocb{
		block:     block,
		tagSize:   tagSize,
		nonceSize: nonceSize,
		mask:      initializeMaskTable(block),
		reusableKtop: reusableKtop{
			noncePrefix: nil,
			Ktop:        nil,
		},
	}, nil
}

func (o *ocb) Seal(dst, nonce, plaintext, adata []byte) []byte {
	if len(nonce) > o.nonceSize {
		panic("crypto/ocb: Incorrect nonce length given to OCB")
	}
	sep := len(plaintext)
	ret, out := byteutil.SliceForAppend(dst, sep+o.tagSize)
	tag := o.crypt(enc, out[:sep], nonce, adata, plaintext)
	copy(out[sep:], tag)
	return ret
}

func (o *ocb) Open(dst, nonce, ciphertext, adata []byte) ([]byte, error) {
	if len(nonce) > o.nonceSize {
		panic("Nonce too long for this instance")
	}
	if len(ciphertext) < o.tagSize {
		return nil, ocbError("Ciphertext shorter than tag length")
	}
	sep := len(ciphertext) - o.tagSize
	ret, out := byteutil.SliceForAppend(dst, sep)
	ciphertextData := ciphertext[:sep]
	tag := o.crypt(dec, out, nonce, adata, ciphertextData)
	if subtle.ConstantTimeCompare(tag, ciphertext[sep:]) == 1 {
		return ret, nil
	}
	for i := range out {
		out[i] = 0
	}
	return nil, ocbError("Tag authentication failed")
}

// On instruction enc (resp. dec), crypt is the encrypt (resp. decrypt)
// function. It writes the resulting plain/ciphertext into Y and returns
// the tag.
func (o *ocb) crypt(instruction int, Y, nonce, adata, X []byte) []byte {
	//
	// Consider X as a sequence of 128-bit blocks
	//
	// Note: For encryption (resp. decryption), X is the plaintext (resp., the
	// ciphertext without the tag).
	blockSize := o.block.BlockSize()

	//
	// Nonce-dependent and per-encryption variables
	//
	// Zero out the last 6 bits of the nonce into truncatedNonce to see if Ktop
	// is already computed.
	truncatedNonce := make([]byte, len(nonce))
	copy(truncatedNonce, nonce)
	truncatedNonce[len(truncatedNonce)-1] &= 192
	var Ktop []byte
	if bytes.Equal(truncatedNonce, o.reusableKtop.noncePrefix) {
		Ktop = o.reusableKtop.Ktop
	} else {
		// Nonce = num2str(TAGLEN mod 128, 7) || zeros(120 - bitlen(N)) || 1 || N
		paddedNonce := append(make([]byte, blockSize-1-len(nonce)), 1)
		paddedNonce = append(paddedNonce, truncatedNonce...)
		paddedNonce[0] |= byte(((8 * o.tagSize) % (8 * blockSize)) << 1)
		// Last 6 bits of paddedNonce are already zero. Encrypt into Ktop
		paddedNonce[blockSize-1] &= 192
		Ktop = paddedNonce
		o.block.Encrypt(Ktop, Ktop)
		o.reusableKtop.noncePrefix = truncatedNonce
		o.reusableKtop.Ktop = Ktop
	}

	// Stretch = Ktop || ((lower half of Ktop) XOR (lower half of Ktop << 8))
	xorHalves := make([]byte, blockSize/2)
	byteutil.XorBytes(xorHalves, Ktop[:blockSize/2], Ktop[1:1+blockSize/2])
	stretch := append(Ktop, xorHalves...)
	bottom := int(nonce[len(nonce)-1] & 63)
	offset := make([]byte, len(stretch))
	byteutil.ShiftNBytesLeft(offset, stretch, bottom)
	offset = offset[:blockSize]

	//
	// Process any whole blocks
	//
	// Note: For encryption Y is ciphertext || tag, for decryption Y is
	// plaintext || tag.
	checksum := make([]byte, blockSize)
	m := len(X) / blockSize
	for i := 0; i < m; i++ {
		index := bits.TrailingZeros(uint(i + 1))
		if len(o.mask.L)-1 < index {
			o.mask.extendTable(index)
		}
		byteutil.XorBytesMut(offset, o.mask.L[bits.TrailingZeros(uint(i+1))])
		blockX := X[i*blockSize : (i+1)*blockSize]
		blockY := Y[i*blockSize : (i+1)*blockSize]
		switch instruction {
		case enc:
			byteutil.XorBytesMut(checksum, blockX)
			byteutil.XorBytes(blockY, blockX, offset)
			o.block.Encrypt(blockY, blockY)
			byteutil.XorBytesMut(blockY, offset)
		case dec:
			byteutil.XorBytes(blockY, blockX, offset)
			o.block.Decrypt(blockY, blockY)
			byteutil.XorBytesMut(blockY, offset)
			byteutil.XorBytesMut(checksum, blockY)
		}
	}
	//
	// Process any final partial block and compute raw tag
	//
	tag := make([]byte, blockSize)
	if len(X)%blockSize != 0 {
		byteutil.XorBytesMut(offset, o.mask.lAst)
		pad := make([]byte, blockSize)
		o.block.Encrypt(pad, offset)
		chunkX := X[blockSize*m:]
		chunkY := Y[blockSize*m : len(X)]
		switch instruction {
		case enc:
			byteutil.XorBytesMut(checksum, chunkX)
			checksum[len(chunkX)] ^= 128
			byteutil.XorBytes(chunkY, chunkX, pad[:len(chunkX)])
			// P_* || bit(1) || zeroes(127) - len(P_*)
		case dec:
			byteutil.XorBytes(chunkY, chunkX, pad[:len(chunkX)])
			// P_* || bit(1) || zeroes(127) - len(P_*)
			byteutil.XorBytesMut(checksum, chunkY)
			checksum[len(chunkY)] ^= 128
		}
	}
	byteutil.XorBytes(tag, checksum, offset)
	byteutil.XorBytesMut(tag, o.mask.lDol)
	o.block.Encrypt(tag, tag)
	byteutil.XorBytesMut(tag, o.hash(adata))
	return tag[:o.tagSize]
}

// This hash function is used to compute the tag. Per design, on empty input it
// returns a slice of zeros, of the same length as the underlying block cipher
// block size.
func (o *ocb) hash(adata []byte) []byte {
	//
	// Consider A as a sequence of 128-bit blocks
	//
	A := make([]byte, len(adata))
	copy(A, adata)
	blockSize := o.block.BlockSize()

	//
	// Process any whole blocks
	//
	sum := make([]byte, blockSize)
	offset := make([]byte, blockSize)
	m := len(A) / blockSize
	for i := 0; i < m; i++ {
		chunk := A[blockSize*i : blockSize*(i+1)]
		index := bits.TrailingZeros(uint(i + 1))
		// If the mask table is too short
		if len(o.mask.L)-1 < index {
			o.mask.extendTable(index)
		}
		byteutil.XorBytesMut(offset, o.mask.L[index])
		byteutil.XorBytesMut(chunk, offset)
		o.block.Encrypt(chunk, chunk)
		byteutil.XorBytesMut(sum, chunk)
	}

	//
	// Process any final partial block; compute final hash value
	//
	if len(A)%blockSize != 0 {
		byteutil.XorBytesMut(offset, o.mask.lAst)
		// Pad block with 1 || 0 ^ 127 - bitlength(a)
		ending := make([]byte, blockSize-len(A)%blockSize)
		ending[0] = 0x80
		encrypted := append(A[blockSize*m:], ending...)
		byteutil.XorBytesMut(encrypted, offset)
		o.block.Encrypt(encrypted, encrypted)
		byteutil.XorBytesMut(sum, encrypted)
	}
	return sum
}

func initializeMaskTable(block cipher.Block) mask {
	//
	// Key-dependent variables
	//
	lAst := make([]byte, block.BlockSize())
	block.Encrypt(lAst, lAst)
	lDol := byteutil.GfnDouble(lAst)
	L := make([][]byte, 1)
	L[0] = byteutil.GfnDouble(lDol)

	return mask{
		lAst: lAst,
		lDol: lDol,
		L:    L,
	}
}

// Extends the L array of mask m up to L[limit], with L[i] = GfnDouble(L[i-1])
func (m *mask) extendTable(limit int) {
	for i := len(m.L); i <= limit; i++ {
		m.L = append(m.L, byteutil.GfnDouble(m.L[i-1]))
	}
}

func ocbError(err string) error {
	return errors.New("crypto/ocb: " + err)
}
