// Package ecc implements a generic interface for ECDH, ECDSA, and EdDSA.
package ecc

import (
	"bytes"
	"crypto/subtle"
	"io"

	"github.com/ProtonMail/go-crypto/openpgp/errors"
	ed25519lib "github.com/cloudflare/circl/sign/ed25519"
)

const ed25519Size = 32

type ed25519 struct{}

func NewEd25519() *ed25519 {
	return &ed25519{}
}

func (c *ed25519) GetCurveName() string {
	return "ed25519"
}

// MarshalBytePoint encodes the public point from native format, adding the prefix.
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-5.5.5.5
func (c *ed25519) MarshalBytePoint(x []byte) []byte {
	return append([]byte{0x40}, x...)
}

// UnmarshalBytePoint decodes a point from prefixed format to native.
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-5.5.5.5
func (c *ed25519) UnmarshalBytePoint(point []byte) (x []byte) {
	if len(point) != ed25519lib.PublicKeySize+1 {
		return nil
	}

	// Return unprefixed
	return point[1:]
}

// MarshalByteSecret encodes a scalar in native format.
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-5.5.5.5
func (c *ed25519) MarshalByteSecret(d []byte) []byte {
	return d
}

// UnmarshalByteSecret decodes a scalar in native format and re-adds the stripped leading zeroes
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-5.5.5.5
func (c *ed25519) UnmarshalByteSecret(s []byte) (d []byte) {
	if len(s) > ed25519lib.SeedSize {
		return nil
	}

	// Handle stripped leading zeroes
	d = make([]byte, ed25519lib.SeedSize)
	copy(d[ed25519lib.SeedSize-len(s):], s)
	return
}

// MarshalSignature splits a signature in R and S.
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-5.2.3.3.1
func (c *ed25519) MarshalSignature(sig []byte) (r, s []byte) {
	return sig[:ed25519Size], sig[ed25519Size:]
}

// UnmarshalSignature decodes R and S in the native format, re-adding the stripped leading zeroes
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-5.2.3.3.1
func (c *ed25519) UnmarshalSignature(r, s []byte) (sig []byte) {
	// Check size
	if len(r) > 32 || len(s) > 32 {
		return nil
	}

	sig = make([]byte, ed25519lib.SignatureSize)

	// Handle stripped leading zeroes
	copy(sig[ed25519Size-len(r):ed25519Size], r)
	copy(sig[ed25519lib.SignatureSize-len(s):], s)
	return sig
}

func (c *ed25519) GenerateEdDSA(rand io.Reader) (pub, priv []byte, err error) {
	pk, sk, err := ed25519lib.GenerateKey(rand)

	if err != nil {
		return nil, nil, err
	}

	return pk, sk[:ed25519lib.SeedSize], nil
}

func getEd25519Sk(publicKey, privateKey []byte) ed25519lib.PrivateKey {
	privateKeyCap, privateKeyLen, publicKeyLen := cap(privateKey), len(privateKey), len(publicKey)

	if privateKeyCap >= privateKeyLen+publicKeyLen &&
		bytes.Equal(privateKey[privateKeyLen:privateKeyLen+publicKeyLen], publicKey) {
		return privateKey[:privateKeyLen+publicKeyLen]
	}

	return append(privateKey[:privateKeyLen:privateKeyLen], publicKey...)
}

func (c *ed25519) Sign(publicKey, privateKey, message []byte) (sig []byte, err error) {
	sig = ed25519lib.Sign(getEd25519Sk(publicKey, privateKey), message)
	return sig, nil
}

func (c *ed25519) Verify(publicKey, message, sig []byte) bool {
	return ed25519lib.Verify(publicKey, message, sig)
}

func (c *ed25519) ValidateEdDSA(publicKey, privateKey []byte) (err error) {
	priv := getEd25519Sk(publicKey, privateKey)
	expectedPriv := ed25519lib.NewKeyFromSeed(priv.Seed())
	if subtle.ConstantTimeCompare(priv, expectedPriv) == 0 {
		return errors.KeyInvalidError("ecc: invalid ed25519 secret")
	}
	return nil
}
