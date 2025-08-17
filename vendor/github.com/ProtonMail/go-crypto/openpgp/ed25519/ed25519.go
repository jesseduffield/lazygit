// Package ed25519 implements the ed25519 signature algorithm for OpenPGP
// as defined in the Open PGP crypto refresh.
package ed25519

import (
	"crypto/subtle"
	"io"

	"github.com/ProtonMail/go-crypto/openpgp/errors"
	ed25519lib "github.com/cloudflare/circl/sign/ed25519"
)

const (
	// PublicKeySize is the size, in bytes, of public keys in this package.
	PublicKeySize = ed25519lib.PublicKeySize
	// SeedSize is the size, in bytes, of private key seeds.
	// The private key representation used by RFC 8032.
	SeedSize = ed25519lib.SeedSize
	// SignatureSize is the size, in bytes, of signatures generated and verified by this package.
	SignatureSize = ed25519lib.SignatureSize
)

type PublicKey struct {
	// Point represents the elliptic curve point of the public key.
	Point []byte
}

type PrivateKey struct {
	PublicKey
	// Key the private key representation by RFC 8032,
	// encoded as seed | pub key point.
	Key []byte
}

// NewPublicKey creates a new empty ed25519 public key.
func NewPublicKey() *PublicKey {
	return &PublicKey{}
}

// NewPrivateKey creates a new empty private key referencing the public key.
func NewPrivateKey(key PublicKey) *PrivateKey {
	return &PrivateKey{
		PublicKey: key,
	}
}

// Seed returns the ed25519 private key secret seed.
// The private key representation by RFC 8032.
func (pk *PrivateKey) Seed() []byte {
	return pk.Key[:SeedSize]
}

// MarshalByteSecret returns the underlying 32 byte seed of the private key.
func (pk *PrivateKey) MarshalByteSecret() []byte {
	return pk.Seed()
}

// UnmarshalByteSecret computes the private key from the secret seed
// and stores it in the private key object.
func (sk *PrivateKey) UnmarshalByteSecret(seed []byte) error {
	sk.Key = ed25519lib.NewKeyFromSeed(seed)
	return nil
}

// GenerateKey generates a fresh private key with the provided randomness source.
func GenerateKey(rand io.Reader) (*PrivateKey, error) {
	publicKey, privateKey, err := ed25519lib.GenerateKey(rand)
	if err != nil {
		return nil, err
	}
	privateKeyOut := new(PrivateKey)
	privateKeyOut.PublicKey.Point = publicKey[:]
	privateKeyOut.Key = privateKey[:]
	return privateKeyOut, nil
}

// Sign signs a message with the ed25519 algorithm.
// priv MUST be a valid key! Check this with Validate() before use.
func Sign(priv *PrivateKey, message []byte) ([]byte, error) {
	return ed25519lib.Sign(priv.Key, message), nil
}

// Verify verifies an ed25519 signature.
func Verify(pub *PublicKey, message []byte, signature []byte) bool {
	return ed25519lib.Verify(pub.Point, message, signature)
}

// Validate checks if the ed25519 private key is valid.
func Validate(priv *PrivateKey) error {
	expectedPrivateKey := ed25519lib.NewKeyFromSeed(priv.Seed())
	if subtle.ConstantTimeCompare(priv.Key, expectedPrivateKey) == 0 {
		return errors.KeyInvalidError("ed25519: invalid ed25519 secret")
	}
	if subtle.ConstantTimeCompare(priv.PublicKey.Point, expectedPrivateKey[SeedSize:]) == 0 {
		return errors.KeyInvalidError("ed25519: invalid ed25519 public key")
	}
	return nil
}

// ENCODING/DECODING signature:

// WriteSignature encodes and writes an ed25519 signature to writer.
func WriteSignature(writer io.Writer, signature []byte) error {
	_, err := writer.Write(signature)
	return err
}

// ReadSignature decodes an ed25519 signature from a reader.
func ReadSignature(reader io.Reader) ([]byte, error) {
	signature := make([]byte, SignatureSize)
	if _, err := io.ReadFull(reader, signature); err != nil {
		return nil, err
	}
	return signature, nil
}
