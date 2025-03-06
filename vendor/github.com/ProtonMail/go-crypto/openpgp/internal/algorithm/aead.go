// Copyright (C) 2019 ProtonTech AG

package algorithm

import (
	"crypto/cipher"
	"github.com/ProtonMail/go-crypto/eax"
	"github.com/ProtonMail/go-crypto/ocb"
)

// AEADMode defines the Authenticated Encryption with Associated Data mode of
// operation.
type AEADMode uint8

// Supported modes of operation (see RFC4880bis [EAX] and RFC7253)
const (
	AEADModeEAX = AEADMode(1)
	AEADModeOCB = AEADMode(2)
	AEADModeGCM = AEADMode(3)
)

// TagLength returns the length in bytes of authentication tags.
func (mode AEADMode) TagLength() int {
	switch mode {
	case AEADModeEAX:
		return 16
	case AEADModeOCB:
		return 16
	case AEADModeGCM:
		return 16
	default:
		return 0
	}
}

// NonceLength returns the length in bytes of nonces.
func (mode AEADMode) NonceLength() int {
	switch mode {
	case AEADModeEAX:
		return 16
	case AEADModeOCB:
		return 15
	case AEADModeGCM:
		return 12
	default:
		return 0
	}
}

// New returns a fresh instance of the given mode
func (mode AEADMode) New(block cipher.Block) (alg cipher.AEAD) {
	var err error
	switch mode {
	case AEADModeEAX:
		alg, err = eax.NewEAX(block)
	case AEADModeOCB:
		alg, err = ocb.NewOCB(block)
	case AEADModeGCM:
		alg, err = cipher.NewGCM(block)
	}
	if err != nil {
		panic(err.Error())
	}
	return alg
}
