// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package algorithm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"

	"golang.org/x/crypto/cast5"
)

// Cipher is an official symmetric key cipher algorithm. See RFC 4880,
// section 9.2.
type Cipher interface {
	// Id returns the algorithm ID, as a byte, of the cipher.
	Id() uint8
	// KeySize returns the key size, in bytes, of the cipher.
	KeySize() int
	// BlockSize returns the block size, in bytes, of the cipher.
	BlockSize() int
	// New returns a fresh instance of the given cipher.
	New(key []byte) cipher.Block
}

// The following constants mirror the OpenPGP standard (RFC 4880).
const (
	TripleDES = CipherFunction(2)
	CAST5     = CipherFunction(3)
	AES128    = CipherFunction(7)
	AES192    = CipherFunction(8)
	AES256    = CipherFunction(9)
)

// CipherById represents the different block ciphers specified for OpenPGP. See
// http://www.iana.org/assignments/pgp-parameters/pgp-parameters.xhtml#pgp-parameters-13
var CipherById = map[uint8]Cipher{
	TripleDES.Id(): TripleDES,
	CAST5.Id():     CAST5,
	AES128.Id():    AES128,
	AES192.Id():    AES192,
	AES256.Id():    AES256,
}

type CipherFunction uint8

// ID returns the algorithm Id, as a byte, of cipher.
func (sk CipherFunction) Id() uint8 {
	return uint8(sk)
}

// KeySize returns the key size, in bytes, of cipher.
func (cipher CipherFunction) KeySize() int {
	switch cipher {
	case CAST5:
		return cast5.KeySize
	case AES128:
		return 16
	case AES192, TripleDES:
		return 24
	case AES256:
		return 32
	}
	return 0
}

// BlockSize returns the block size, in bytes, of cipher.
func (cipher CipherFunction) BlockSize() int {
	switch cipher {
	case TripleDES:
		return des.BlockSize
	case CAST5:
		return 8
	case AES128, AES192, AES256:
		return 16
	}
	return 0
}

// New returns a fresh instance of the given cipher.
func (cipher CipherFunction) New(key []byte) (block cipher.Block) {
	var err error
	switch cipher {
	case TripleDES:
		block, err = des.NewTripleDESCipher(key)
	case CAST5:
		block, err = cast5.NewCipher(key)
	case AES128, AES192, AES256:
		block, err = aes.NewCipher(key)
	}
	if err != nil {
		panic(err.Error())
	}
	return
}
