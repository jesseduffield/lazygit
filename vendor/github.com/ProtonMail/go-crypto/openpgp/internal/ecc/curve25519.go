// Package ecc implements a generic interface for ECDH, ECDSA, and EdDSA.
package ecc

import (
	"crypto/subtle"
	"io"

	"github.com/ProtonMail/go-crypto/openpgp/errors"
	x25519lib "github.com/cloudflare/circl/dh/x25519"
)

type curve25519 struct{}

func NewCurve25519() *curve25519 {
	return &curve25519{}
}

func (c *curve25519) GetCurveName() string {
	return "curve25519"
}

// MarshalBytePoint encodes the public point from native format, adding the prefix.
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-5.5.5.6
func (c *curve25519) MarshalBytePoint(point []byte) []byte {
	return append([]byte{0x40}, point...)
}

// UnmarshalBytePoint decodes the public point to native format, removing the prefix.
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-5.5.5.6
func (c *curve25519) UnmarshalBytePoint(point []byte) []byte {
	if len(point) != x25519lib.Size+1 {
		return nil
	}

	// Remove prefix
	return point[1:]
}

// MarshalByteSecret encodes the secret scalar from native format.
// Note that the EC secret scalar differs from the definition of public keys in
// [Curve25519] in two ways: (1) the byte-ordering is big-endian, which is
// more uniform with how big integers are represented in OpenPGP, and (2) the
// leading zeros are truncated.
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-5.5.5.6.1.1
// Note that leading zero bytes are stripped later when encoding as an MPI.
func (c *curve25519) MarshalByteSecret(secret []byte) []byte {
	d := make([]byte, x25519lib.Size)
	copyReversed(d, secret)

	// The following ensures that the private key is a number of the form
	// 2^{254} + 8 * [0, 2^{251}), in order to avoid the small subgroup of
	// the curve.
	//
	// This masking is done internally in the underlying lib and so is unnecessary
	// for security, but OpenPGP implementations require that private keys be
	// pre-masked.
	d[0] &= 127
	d[0] |= 64
	d[31] &= 248

	return d
}

// UnmarshalByteSecret decodes the secret scalar from native format.
// Note that the EC secret scalar differs from the definition of public keys in
// [Curve25519] in two ways: (1) the byte-ordering is big-endian, which is
// more uniform with how big integers are represented in OpenPGP, and (2) the
// leading zeros are truncated.
// See https://datatracker.ietf.org/doc/html/draft-ietf-openpgp-crypto-refresh-06#section-5.5.5.6.1.1
func (c *curve25519) UnmarshalByteSecret(d []byte) []byte {
	if len(d) > x25519lib.Size {
		return nil
	}

	// Ensure truncated leading bytes are re-added
	secret := make([]byte, x25519lib.Size)
	copyReversed(secret, d)

	return secret
}

// generateKeyPairBytes Generates a private-public key-pair.
// 'priv' is a private key; a little-endian scalar belonging to the set
// 2^{254} + 8 * [0, 2^{251}), in order to avoid the small subgroup of the
// curve. 'pub' is simply 'priv' * G where G is the base point.
// See https://cr.yp.to/ecdh.html and RFC7748, sec 5.
func (c *curve25519) generateKeyPairBytes(rand io.Reader) (priv, pub x25519lib.Key, err error) {
	_, err = io.ReadFull(rand, priv[:])
	if err != nil {
		return
	}

	x25519lib.KeyGen(&pub, &priv)
	return
}

func (c *curve25519) GenerateECDH(rand io.Reader) (point []byte, secret []byte, err error) {
	priv, pub, err := c.generateKeyPairBytes(rand)
	if err != nil {
		return
	}

	return pub[:], priv[:], nil
}

func (c *genericCurve) MaskSecret(secret []byte) []byte {
	return secret
}

func (c *curve25519) Encaps(rand io.Reader, point []byte) (ephemeral, sharedSecret []byte, err error) {
	// RFC6637 §8: "Generate an ephemeral key pair {v, V=vG}"
	// ephemeralPrivate corresponds to `v`.
	// ephemeralPublic corresponds to `V`.
	ephemeralPrivate, ephemeralPublic, err := c.generateKeyPairBytes(rand)
	if err != nil {
		return nil, nil, err
	}

	// RFC6637 §8: "Obtain the authenticated recipient public key R"
	// pubKey corresponds to `R`.
	var pubKey x25519lib.Key
	copy(pubKey[:], point)

	// RFC6637 §8: "Compute the shared point S = vR"
	//	"VB = convert point V to the octet string"
	// sharedPoint corresponds to `VB`.
	var sharedPoint x25519lib.Key
	x25519lib.Shared(&sharedPoint, &ephemeralPrivate, &pubKey)

	return ephemeralPublic[:], sharedPoint[:], nil
}

func (c *curve25519) Decaps(vsG, secret []byte) (sharedSecret []byte, err error) {
	var ephemeralPublic, decodedPrivate, sharedPoint x25519lib.Key
	// RFC6637 §8: "The decryption is the inverse of the method given."
	// All quoted descriptions in comments below describe encryption, and
	// the reverse is performed.
	// vsG corresponds to `VB` in RFC6637 §8 .

	// RFC6637 §8: "VB = convert point V to the octet string"
	copy(ephemeralPublic[:], vsG)

	// decodedPrivate corresponds to `r` in RFC6637 §8 .
	copy(decodedPrivate[:], secret)

	// RFC6637 §8: "Note that the recipient obtains the shared secret by calculating
	//   S = rV = rvG, where (r,R) is the recipient's key pair."
	// sharedPoint corresponds to `S`.
	x25519lib.Shared(&sharedPoint, &decodedPrivate, &ephemeralPublic)

	return sharedPoint[:], nil
}

func (c *curve25519) ValidateECDH(point []byte, secret []byte) (err error) {
	var pk, sk x25519lib.Key
	copy(sk[:], secret)
	x25519lib.KeyGen(&pk, &sk)

	if subtle.ConstantTimeCompare(point, pk[:]) == 0 {
		return errors.KeyInvalidError("ecc: invalid curve25519 public point")
	}

	return nil
}

func copyReversed(out []byte, in []byte) {
	l := len(in)
	for i := 0; i < l; i++ {
		out[i] = in[l-i-1]
	}
}
