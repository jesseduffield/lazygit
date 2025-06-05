// Package ecc implements a generic interface for ECDH, ECDSA, and EdDSA.
package ecc

import (
	"bytes"
	"crypto/subtle"
	"io"

	"github.com/ProtonMail/go-crypto/openpgp/errors"
	ed448lib "github.com/cloudflare/circl/sign/ed448"
)

type ed448 struct{}

func NewEd448() *ed448 {
	return &ed448{}
}

func (c *ed448) GetCurveName() string {
	return "ed448"
}

// MarshalBytePoint encodes the public point from native format, adding the prefix.
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-5.5.5.5
func (c *ed448) MarshalBytePoint(x []byte) []byte {
	// Return prefixed
	return append([]byte{0x40}, x...)
}

// UnmarshalBytePoint decodes a point from prefixed format to native.
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-5.5.5.5
func (c *ed448) UnmarshalBytePoint(point []byte) (x []byte) {
	if len(point) != ed448lib.PublicKeySize+1 {
		return nil
	}

	// Strip prefix
	return point[1:]
}

// MarshalByteSecret encoded a scalar from native format to prefixed.
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-5.5.5.5
func (c *ed448) MarshalByteSecret(d []byte) []byte {
	// Return prefixed
	return append([]byte{0x40}, d...)
}

// UnmarshalByteSecret decodes a scalar from prefixed format to native.
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-5.5.5.5
func (c *ed448) UnmarshalByteSecret(s []byte) (d []byte) {
	// Check prefixed size
	if len(s) != ed448lib.SeedSize+1 {
		return nil
	}

	// Strip prefix
	return s[1:]
}

// MarshalSignature splits a signature in R and S, where R is in prefixed native format and
// S is an MPI with value zero.
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-5.2.3.3.2
func (c *ed448) MarshalSignature(sig []byte) (r, s []byte) {
	return append([]byte{0x40}, sig...), []byte{}
}

// UnmarshalSignature decodes R and S in the native format. Only R is used, in prefixed native format.
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-5.2.3.3.2
func (c *ed448) UnmarshalSignature(r, s []byte) (sig []byte) {
	if len(r) != ed448lib.SignatureSize+1 {
		return nil
	}

	return r[1:]
}

func (c *ed448) GenerateEdDSA(rand io.Reader) (pub, priv []byte, err error) {
	pk, sk, err := ed448lib.GenerateKey(rand)

	if err != nil {
		return nil, nil, err
	}

	return pk, sk[:ed448lib.SeedSize], nil
}

func getEd448Sk(publicKey, privateKey []byte) ed448lib.PrivateKey {
	privateKeyCap, privateKeyLen, publicKeyLen := cap(privateKey), len(privateKey), len(publicKey)

	if privateKeyCap >= privateKeyLen+publicKeyLen &&
		bytes.Equal(privateKey[privateKeyLen:privateKeyLen+publicKeyLen], publicKey) {
		return privateKey[:privateKeyLen+publicKeyLen]
	}

	return append(privateKey[:privateKeyLen:privateKeyLen], publicKey...)
}

func (c *ed448) Sign(publicKey, privateKey, message []byte) (sig []byte, err error) {
	// Ed448 is used with the empty string as a context string.
	// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-13.7
	sig = ed448lib.Sign(getEd448Sk(publicKey, privateKey), message, "")

	return sig, nil
}

func (c *ed448) Verify(publicKey, message, sig []byte) bool {
	// Ed448 is used with the empty string as a context string.
	// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-13.7
	return ed448lib.Verify(publicKey, message, sig, "")
}

func (c *ed448) ValidateEdDSA(publicKey, privateKey []byte) (err error) {
	priv := getEd448Sk(publicKey, privateKey)
	expectedPriv := ed448lib.NewKeyFromSeed(priv.Seed())
	if subtle.ConstantTimeCompare(priv, expectedPriv) == 0 {
		return errors.KeyInvalidError("ecc: invalid ed448 secret")
	}
	return nil
}
