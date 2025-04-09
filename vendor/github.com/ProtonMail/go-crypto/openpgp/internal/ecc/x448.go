// Package ecc implements a generic interface for ECDH, ECDSA, and EdDSA.
package ecc

import (
	"crypto/subtle"
	"io"

	"github.com/ProtonMail/go-crypto/openpgp/errors"
	x448lib "github.com/cloudflare/circl/dh/x448"
)

type x448 struct{}

func NewX448() *x448 {
	return &x448{}
}

func (c *x448) GetCurveName() string {
	return "x448"
}

// MarshalBytePoint encodes the public point from native format, adding the prefix.
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-5.5.5.6
func (c *x448) MarshalBytePoint(point []byte) []byte {
	return append([]byte{0x40}, point...)
}

// UnmarshalBytePoint decodes a point from prefixed format to native.
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-5.5.5.6
func (c *x448) UnmarshalBytePoint(point []byte) []byte {
	if len(point) != x448lib.Size+1 {
		return nil
	}

	return point[1:]
}

// MarshalByteSecret encoded a scalar from native format to prefixed.
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-5.5.5.6.1.2
func (c *x448) MarshalByteSecret(d []byte) []byte {
	return append([]byte{0x40}, d...)
}

// UnmarshalByteSecret decodes a scalar from prefixed format to native.
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-5.5.5.6.1.2
func (c *x448) UnmarshalByteSecret(d []byte) []byte {
	if len(d) != x448lib.Size+1 {
		return nil
	}

	// Store without prefix
	return d[1:]
}

func (c *x448) generateKeyPairBytes(rand io.Reader) (sk, pk x448lib.Key, err error) {
	if _, err = rand.Read(sk[:]); err != nil {
		return
	}

	x448lib.KeyGen(&pk, &sk)
	return
}

func (c *x448) GenerateECDH(rand io.Reader) (point []byte, secret []byte, err error) {
	priv, pub, err := c.generateKeyPairBytes(rand)
	if err != nil {
		return
	}

	return pub[:], priv[:], nil
}

func (c *x448) Encaps(rand io.Reader, point []byte) (ephemeral, sharedSecret []byte, err error) {
	var pk, ss x448lib.Key
	seed, e, err := c.generateKeyPairBytes(rand)
	if err != nil {
		return nil, nil, err
	}
	copy(pk[:], point)
	x448lib.Shared(&ss, &seed, &pk)

	return e[:], ss[:], nil
}

func (c *x448) Decaps(ephemeral, secret []byte) (sharedSecret []byte, err error) {
	var ss, sk, e x448lib.Key

	copy(sk[:], secret)
	copy(e[:], ephemeral)
	x448lib.Shared(&ss, &sk, &e)

	return ss[:], nil
}

func (c *x448) ValidateECDH(point []byte, secret []byte) error {
	var sk, pk, expectedPk x448lib.Key

	copy(pk[:], point)
	copy(sk[:], secret)
	x448lib.KeyGen(&expectedPk, &sk)

	if subtle.ConstantTimeCompare(expectedPk[:], pk[:]) == 0 {
		return errors.KeyInvalidError("ecc: invalid curve25519 public point")
	}

	return nil
}
