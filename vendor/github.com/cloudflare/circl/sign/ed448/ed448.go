// Package ed448 implements Ed448 signature scheme as described in RFC-8032.
//
// This package implements two signature variants.
//
//	| Scheme Name | Sign Function     | Verification  | Context           |
//	|-------------|-------------------|---------------|-------------------|
//	| Ed448       | Sign              | Verify        | Yes, can be empty |
//	| Ed448Ph     | SignPh            | VerifyPh      | Yes, can be empty |
//	| All above   | (PrivateKey).Sign | VerifyAny     | As above          |
//
// Specific functions for sign and verify are defined. A generic signing
// function for all schemes is available through the crypto.Signer interface,
// which is implemented by the PrivateKey type. A correspond all-in-one
// verification method is provided by the VerifyAny function.
//
// Both schemes require a context string for domain separation. This parameter
// is passed using a SignerOptions struct defined in this package.
//
// References:
//
//   - RFC8032: https://rfc-editor.org/rfc/rfc8032.txt
//   - EdDSA for more curves: https://eprint.iacr.org/2015/677
//   - High-speed high-security signatures: https://doi.org/10.1007/s13389-012-0027-1
package ed448

import (
	"bytes"
	"crypto"
	cryptoRand "crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/cloudflare/circl/ecc/goldilocks"
	"github.com/cloudflare/circl/internal/sha3"
	"github.com/cloudflare/circl/sign"
)

const (
	// ContextMaxSize is the maximum length (in bytes) allowed for context.
	ContextMaxSize = 255
	// PublicKeySize is the length in bytes of Ed448 public keys.
	PublicKeySize = 57
	// PrivateKeySize is the length in bytes of Ed448 private keys.
	PrivateKeySize = 114
	// SignatureSize is the length in bytes of signatures.
	SignatureSize = 114
	// SeedSize is the size, in bytes, of private key seeds. These are the private key representations used by RFC 8032.
	SeedSize = 57
)

const (
	paramB   = 456 / 8    // Size of keys in bytes.
	hashSize = 2 * paramB // Size of the hash function's output.
)

// SignerOptions implements crypto.SignerOpts and augments with parameters
// that are specific to the Ed448 signature schemes.
type SignerOptions struct {
	// Hash must be crypto.Hash(0) for both Ed448 and Ed448Ph.
	crypto.Hash

	// Context is an optional domain separation string for signing.
	// Its length must be less or equal than 255 bytes.
	Context string

	// Scheme is an identifier for choosing a signature scheme.
	Scheme SchemeID
}

// SchemeID is an identifier for each signature scheme.
type SchemeID uint

const (
	ED448 SchemeID = iota
	ED448Ph
)

// PublicKey is the type of Ed448 public keys.
type PublicKey []byte

// Equal reports whether pub and x have the same value.
func (pub PublicKey) Equal(x crypto.PublicKey) bool {
	xx, ok := x.(PublicKey)
	return ok && bytes.Equal(pub, xx)
}

// PrivateKey is the type of Ed448 private keys. It implements crypto.Signer.
type PrivateKey []byte

// Equal reports whether priv and x have the same value.
func (priv PrivateKey) Equal(x crypto.PrivateKey) bool {
	xx, ok := x.(PrivateKey)
	return ok && subtle.ConstantTimeCompare(priv, xx) == 1
}

// Public returns the PublicKey corresponding to priv.
func (priv PrivateKey) Public() crypto.PublicKey {
	publicKey := make([]byte, PublicKeySize)
	copy(publicKey, priv[SeedSize:])
	return PublicKey(publicKey)
}

// Seed returns the private key seed corresponding to priv. It is provided for
// interoperability with RFC 8032. RFC 8032's private keys correspond to seeds
// in this package.
func (priv PrivateKey) Seed() []byte {
	seed := make([]byte, SeedSize)
	copy(seed, priv[:SeedSize])
	return seed
}

func (priv PrivateKey) Scheme() sign.Scheme { return sch }

func (pub PublicKey) Scheme() sign.Scheme { return sch }

func (priv PrivateKey) MarshalBinary() (data []byte, err error) {
	privateKey := make(PrivateKey, PrivateKeySize)
	copy(privateKey, priv)
	return privateKey, nil
}

func (pub PublicKey) MarshalBinary() (data []byte, err error) {
	publicKey := make(PublicKey, PublicKeySize)
	copy(publicKey, pub)
	return publicKey, nil
}

// Sign creates a signature of a message given a key pair.
// This function supports all the two signature variants defined in RFC-8032,
// namely Ed448 (or pure EdDSA) and Ed448Ph.
// The opts.HashFunc() must return zero to the specify Ed448 variant. This can
// be achieved by passing crypto.Hash(0) as the value for opts.
// Use an Options struct to pass a bool indicating that the ed448Ph variant
// should be used.
// The struct can also be optionally used to pass a context string for signing.
func (priv PrivateKey) Sign(
	rand io.Reader,
	message []byte,
	opts crypto.SignerOpts,
) (signature []byte, err error) {
	var ctx string
	var scheme SchemeID

	if o, ok := opts.(SignerOptions); ok {
		ctx = o.Context
		scheme = o.Scheme
	}

	switch true {
	case scheme == ED448 && opts.HashFunc() == crypto.Hash(0):
		return Sign(priv, message, ctx), nil
	case scheme == ED448Ph && opts.HashFunc() == crypto.Hash(0):
		return SignPh(priv, message, ctx), nil
	default:
		return nil, errors.New("ed448: bad hash algorithm")
	}
}

// GenerateKey generates a public/private key pair using entropy from rand.
// If rand is nil, crypto/rand.Reader will be used.
func GenerateKey(rand io.Reader) (PublicKey, PrivateKey, error) {
	if rand == nil {
		rand = cryptoRand.Reader
	}

	seed := make(PrivateKey, SeedSize)
	if _, err := io.ReadFull(rand, seed); err != nil {
		return nil, nil, err
	}

	privateKey := NewKeyFromSeed(seed)
	publicKey := make([]byte, PublicKeySize)
	copy(publicKey, privateKey[SeedSize:])

	return publicKey, privateKey, nil
}

// NewKeyFromSeed calculates a private key from a seed. It will panic if
// len(seed) is not SeedSize. This function is provided for interoperability
// with RFC 8032. RFC 8032's private keys correspond to seeds in this
// package.
func NewKeyFromSeed(seed []byte) PrivateKey {
	privateKey := make([]byte, PrivateKeySize)
	newKeyFromSeed(privateKey, seed)
	return privateKey
}

func newKeyFromSeed(privateKey, seed []byte) {
	if l := len(seed); l != SeedSize {
		panic("ed448: bad seed length: " + strconv.Itoa(l))
	}

	var h [hashSize]byte
	H := sha3.NewShake256()
	_, _ = H.Write(seed)
	_, _ = H.Read(h[:])
	s := &goldilocks.Scalar{}
	deriveSecretScalar(s, h[:paramB])

	copy(privateKey[:SeedSize], seed)
	_ = goldilocks.Curve{}.ScalarBaseMult(s).ToBytes(privateKey[SeedSize:])
}

func signAll(signature []byte, privateKey PrivateKey, message, ctx []byte, preHash bool) {
	if len(ctx) > ContextMaxSize {
		panic(fmt.Errorf("ed448: bad context length: %v", len(ctx)))
	}

	H := sha3.NewShake256()
	var PHM []byte

	if preHash {
		var h [64]byte
		_, _ = H.Write(message)
		_, _ = H.Read(h[:])
		PHM = h[:]
		H.Reset()
	} else {
		PHM = message
	}

	// 1.  Hash the 57-byte private key using SHAKE256(x, 114).
	var h [hashSize]byte
	_, _ = H.Write(privateKey[:SeedSize])
	_, _ = H.Read(h[:])
	s := &goldilocks.Scalar{}
	deriveSecretScalar(s, h[:paramB])
	prefix := h[paramB:]

	// 2.  Compute SHAKE256(dom4(F, C) || prefix || PH(M), 114).
	var rPM [hashSize]byte
	H.Reset()

	writeDom(&H, ctx, preHash)

	_, _ = H.Write(prefix)
	_, _ = H.Write(PHM)
	_, _ = H.Read(rPM[:])

	// 3.  Compute the point [r]B.
	r := &goldilocks.Scalar{}
	r.FromBytes(rPM[:])
	R := (&[paramB]byte{})[:]
	if err := (goldilocks.Curve{}.ScalarBaseMult(r).ToBytes(R)); err != nil {
		panic(err)
	}
	// 4.  Compute SHAKE256(dom4(F, C) || R || A || PH(M), 114)
	var hRAM [hashSize]byte
	H.Reset()

	writeDom(&H, ctx, preHash)

	_, _ = H.Write(R)
	_, _ = H.Write(privateKey[SeedSize:])
	_, _ = H.Write(PHM)
	_, _ = H.Read(hRAM[:])

	// 5.  Compute S = (r + k * s) mod order.
	k := &goldilocks.Scalar{}
	k.FromBytes(hRAM[:])
	S := &goldilocks.Scalar{}
	S.Mul(k, s)
	S.Add(S, r)

	// 6.  The signature is the concatenation of R and S.
	copy(signature[:paramB], R[:])
	copy(signature[paramB:], S[:])
}

// Sign signs the message with privateKey and returns a signature.
// This function supports the signature variant defined in RFC-8032: Ed448,
// also known as the pure version of EdDSA.
// It will panic if len(privateKey) is not PrivateKeySize.
func Sign(priv PrivateKey, message []byte, ctx string) []byte {
	signature := make([]byte, SignatureSize)
	signAll(signature, priv, message, []byte(ctx), false)
	return signature
}

// SignPh creates a signature of a message given a keypair.
// This function supports the signature variant defined in RFC-8032: Ed448ph,
// meaning it internally hashes the message using SHAKE-256.
// Context could be passed to this function, which length should be no more than
// 255. It can be empty.
func SignPh(priv PrivateKey, message []byte, ctx string) []byte {
	signature := make([]byte, SignatureSize)
	signAll(signature, priv, message, []byte(ctx), true)
	return signature
}

func verify(public PublicKey, message, signature, ctx []byte, preHash bool) bool {
	if len(public) != PublicKeySize ||
		len(signature) != SignatureSize ||
		len(ctx) > ContextMaxSize ||
		!isLessThanOrder(signature[paramB:]) {
		return false
	}

	P, err := goldilocks.FromBytes(public)
	if err != nil {
		return false
	}

	H := sha3.NewShake256()
	var PHM []byte

	if preHash {
		var h [64]byte
		_, _ = H.Write(message)
		_, _ = H.Read(h[:])
		PHM = h[:]
		H.Reset()
	} else {
		PHM = message
	}

	var hRAM [hashSize]byte
	R := signature[:paramB]

	writeDom(&H, ctx, preHash)

	_, _ = H.Write(R)
	_, _ = H.Write(public)
	_, _ = H.Write(PHM)
	_, _ = H.Read(hRAM[:])

	k := &goldilocks.Scalar{}
	k.FromBytes(hRAM[:])
	S := &goldilocks.Scalar{}
	S.FromBytes(signature[paramB:])

	encR := (&[paramB]byte{})[:]
	P.Neg()
	_ = goldilocks.Curve{}.CombinedMult(S, k, P).ToBytes(encR)
	return bytes.Equal(R, encR)
}

// VerifyAny returns true if the signature is valid. Failure cases are invalid
// signature, or when the public key cannot be decoded.
// This function supports all the two signature variants defined in RFC-8032,
// namely Ed448 (or pure EdDSA) and Ed448Ph.
// The opts.HashFunc() must return zero, this can be achieved by passing
// crypto.Hash(0) as the value for opts.
// Use a SignerOptions struct to pass a context string for signing.
func VerifyAny(public PublicKey, message, signature []byte, opts crypto.SignerOpts) bool {
	var ctx string
	var scheme SchemeID
	if o, ok := opts.(SignerOptions); ok {
		ctx = o.Context
		scheme = o.Scheme
	}

	switch true {
	case scheme == ED448 && opts.HashFunc() == crypto.Hash(0):
		return Verify(public, message, signature, ctx)
	case scheme == ED448Ph && opts.HashFunc() == crypto.Hash(0):
		return VerifyPh(public, message, signature, ctx)
	default:
		return false
	}
}

// Verify returns true if the signature is valid. Failure cases are invalid
// signature, or when the public key cannot be decoded.
// This function supports the signature variant defined in RFC-8032: Ed448,
// also known as the pure version of EdDSA.
func Verify(public PublicKey, message, signature []byte, ctx string) bool {
	return verify(public, message, signature, []byte(ctx), false)
}

// VerifyPh returns true if the signature is valid. Failure cases are invalid
// signature, or when the public key cannot be decoded.
// This function supports the signature variant defined in RFC-8032: Ed448ph,
// meaning it internally hashes the message using SHAKE-256.
// Context could be passed to this function, which length should be no more than
// 255. It can be empty.
func VerifyPh(public PublicKey, message, signature []byte, ctx string) bool {
	return verify(public, message, signature, []byte(ctx), true)
}

func deriveSecretScalar(s *goldilocks.Scalar, h []byte) {
	h[0] &= 0xFC        // The two least significant bits of the first octet are cleared,
	h[paramB-1] = 0x00  // all eight bits the last octet are cleared, and
	h[paramB-2] |= 0x80 // the highest bit of the second to last octet is set.
	s.FromBytes(h[:paramB])
}

// isLessThanOrder returns true if 0 <= x < order and if the last byte of x is zero.
func isLessThanOrder(x []byte) bool {
	order := goldilocks.Curve{}.Order()
	i := len(order) - 1
	for i > 0 && x[i] == order[i] {
		i--
	}
	return x[paramB-1] == 0 && x[i] < order[i]
}

func writeDom(h io.Writer, ctx []byte, preHash bool) {
	dom4 := "SigEd448"
	_, _ = h.Write([]byte(dom4))

	if preHash {
		_, _ = h.Write([]byte{byte(0x01), byte(len(ctx))})
	} else {
		_, _ = h.Write([]byte{byte(0x00), byte(len(ctx))})
	}
	_, _ = h.Write(ctx)
}
