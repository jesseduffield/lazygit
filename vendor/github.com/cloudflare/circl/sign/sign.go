// Package sign provides unified interfaces for signature schemes.
//
// A register of schemes is available in the package
//
//	github.com/cloudflare/circl/sign/schemes
package sign

import (
	"crypto"
	"encoding"
	"errors"
)

type SignatureOpts struct {
	// If non-empty, includes the given context in the signature if supported
	// and will cause an error during signing otherwise.
	Context string
}

// A public key is used to verify a signature set by the corresponding private
// key.
type PublicKey interface {
	// Returns the signature scheme for this public key.
	Scheme() Scheme
	Equal(crypto.PublicKey) bool
	encoding.BinaryMarshaler
	crypto.PublicKey
}

// A private key allows one to create signatures.
type PrivateKey interface {
	// Returns the signature scheme for this private key.
	Scheme() Scheme
	Equal(crypto.PrivateKey) bool
	// For compatibility with Go standard library
	crypto.Signer
	crypto.PrivateKey
	encoding.BinaryMarshaler
}

// A private key that retains the seed with which it was generated.
type Seeded interface {
	// returns the seed if retained, otherwise nil
	Seed() []byte
}

// A Scheme represents a specific instance of a signature scheme.
type Scheme interface {
	// Name of the scheme.
	Name() string

	// GenerateKey creates a new key-pair.
	GenerateKey() (PublicKey, PrivateKey, error)

	// Creates a signature using the PrivateKey on the given message and
	// returns the signature. opts are additional options which can be nil.
	//
	// Panics if key is nil or wrong type or opts context is not supported.
	Sign(sk PrivateKey, message []byte, opts *SignatureOpts) []byte

	// Checks whether the given signature is a valid signature set by
	// the private key corresponding to the given public key on the
	// given message. opts are additional options which can be nil.
	//
	// Panics if key is nil or wrong type or opts context is not supported.
	Verify(pk PublicKey, message []byte, signature []byte, opts *SignatureOpts) bool

	// Deterministically derives a keypair from a seed. If you're unsure,
	// you're better off using GenerateKey().
	//
	// Panics if seed is not of length SeedSize().
	DeriveKey(seed []byte) (PublicKey, PrivateKey)

	// Unmarshals a PublicKey from the provided buffer.
	UnmarshalBinaryPublicKey([]byte) (PublicKey, error)

	// Unmarshals a PublicKey from the provided buffer.
	UnmarshalBinaryPrivateKey([]byte) (PrivateKey, error)

	// Size of binary marshalled public keys.
	PublicKeySize() int

	// Size of binary marshalled public keys.
	PrivateKeySize() int

	// Size of signatures.
	SignatureSize() int

	// Size of seeds.
	SeedSize() int

	// Returns whether contexts are supported.
	SupportsContext() bool
}

var (
	// ErrTypeMismatch is the error used if types of, for instance, private
	// and public keys don't match.
	ErrTypeMismatch = errors.New("types mismatch")

	// ErrSeedSize is the error used if the provided seed is of the wrong
	// size.
	ErrSeedSize = errors.New("wrong seed size")

	// ErrPubKeySize is the error used if the provided public key is of
	// the wrong size.
	ErrPubKeySize = errors.New("wrong size for public key")

	// ErrPrivKeySize is the error used if the provided private key is of
	// the wrong size.
	ErrPrivKeySize = errors.New("wrong size for private key")

	// ErrContextNotSupported is the error used if a context is not
	// supported.
	ErrContextNotSupported = errors.New("context not supported")

	// ErrContextTooLong is the error used if the context string is too long.
	ErrContextTooLong = errors.New("context string too long")
)
