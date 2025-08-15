// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"crypto"
	"crypto/rand"
	"io"
	"math/big"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp/s2k"
)

var (
	defaultRejectPublicKeyAlgorithms = map[PublicKeyAlgorithm]bool{
		PubKeyAlgoElGamal: true,
		PubKeyAlgoDSA:     true,
	}
	defaultRejectHashAlgorithms = map[crypto.Hash]bool{
		crypto.MD5:       true,
		crypto.RIPEMD160: true,
	}
	defaultRejectMessageHashAlgorithms = map[crypto.Hash]bool{
		crypto.SHA1:      true,
		crypto.MD5:       true,
		crypto.RIPEMD160: true,
	}
	defaultRejectCurves = map[Curve]bool{
		CurveSecP256k1: true,
	}
)

// A global feature flag to indicate v5 support.
// Can be set via a build tag, e.g.: `go build -tags v5 ./...`
// If the build tag is missing config_v5.go will set it to true.
//
// Disables parsing of v5 keys and v5 signatures.
// These are non-standard entities, which in the crypto-refresh have been superseded
// by v6 keys, v6 signatures and SEIPDv2 encrypted data, respectively.
var V5Disabled = false

// Config collects a number of parameters along with sensible defaults.
// A nil *Config is valid and results in all default values.
type Config struct {
	// Rand provides the source of entropy.
	// If nil, the crypto/rand Reader is used.
	Rand io.Reader
	// DefaultHash is the default hash function to be used.
	// If zero, SHA-256 is used.
	DefaultHash crypto.Hash
	// DefaultCipher is the cipher to be used.
	// If zero, AES-128 is used.
	DefaultCipher CipherFunction
	// Time returns the current time as the number of seconds since the
	// epoch. If Time is nil, time.Now is used.
	Time func() time.Time
	// DefaultCompressionAlgo is the compression algorithm to be
	// applied to the plaintext before encryption. If zero, no
	// compression is done.
	DefaultCompressionAlgo CompressionAlgo
	// CompressionConfig configures the compression settings.
	CompressionConfig *CompressionConfig
	// S2K (String to Key) config, used for key derivation in the context of secret key encryption
	// and password-encrypted data.
	// If nil, the default configuration is used
	S2KConfig *s2k.Config
	// Iteration count for Iterated S2K (String to Key).
	// Only used if sk2.Mode is nil.
	// This value is duplicated here from s2k.Config for backwards compatibility.
	// It determines the strength of the passphrase stretching when
	// the said passphrase is hashed to produce a key. S2KCount
	// should be between 65536 and 65011712, inclusive. If Config
	// is nil or S2KCount is 0, the value 16777216 used. Not all
	// values in the above range can be represented. S2KCount will
	// be rounded up to the next representable value if it cannot
	// be encoded exactly. When set, it is strongly encrouraged to
	// use a value that is at least 65536. See RFC 4880 Section
	// 3.7.1.3.
	//
	// Deprecated: SK2Count should be configured in S2KConfig instead.
	S2KCount int
	// RSABits is the number of bits in new RSA keys made with NewEntity.
	// If zero, then 2048 bit keys are created.
	RSABits int
	// The public key algorithm to use - will always create a signing primary
	// key and encryption subkey.
	Algorithm PublicKeyAlgorithm
	// Some known primes that are optionally prepopulated by the caller
	RSAPrimes []*big.Int
	// Curve configures the desired packet.Curve if the Algorithm is PubKeyAlgoECDSA,
	// PubKeyAlgoEdDSA, or PubKeyAlgoECDH. If empty Curve25519 is used.
	Curve Curve
	// AEADConfig configures the use of the new AEAD Encrypted Data Packet,
	// defined in the draft of the next version of the OpenPGP specification.
	// If a non-nil AEADConfig is passed, usage of this packet is enabled. By
	// default, it is disabled. See the documentation of AEADConfig for more
	// configuration options related to AEAD.
	// **Note: using this option may break compatibility with other OpenPGP
	// implementations, as well as future versions of this library.**
	AEADConfig *AEADConfig
	// V6Keys configures version 6 key generation. If false, this package still
	// supports version 6 keys, but produces version 4 keys.
	V6Keys bool
	// Minimum RSA key size allowed for key generation and message signing, verification and encryption.
	MinRSABits uint16
	// Reject insecure algorithms, only works with v2 api
	RejectPublicKeyAlgorithms   map[PublicKeyAlgorithm]bool
	RejectHashAlgorithms        map[crypto.Hash]bool
	RejectMessageHashAlgorithms map[crypto.Hash]bool
	RejectCurves                map[Curve]bool
	// "The validity period of the key.  This is the number of seconds after
	// the key creation time that the key expires.  If this is not present
	// or has a value of zero, the key never expires.  This is found only on
	// a self-signature.""
	// https://tools.ietf.org/html/rfc4880#section-5.2.3.6
	KeyLifetimeSecs uint32
	// "The validity period of the signature.  This is the number of seconds
	// after the signature creation time that the signature expires.  If
	// this is not present or has a value of zero, it never expires."
	// https://tools.ietf.org/html/rfc4880#section-5.2.3.10
	SigLifetimeSecs uint32
	// SigningKeyId is used to specify the signing key to use (by Key ID).
	// By default, the signing key is selected automatically, preferring
	// signing subkeys if available.
	SigningKeyId uint64
	// SigningIdentity is used to specify a user ID (packet Signer's User ID, type 28)
	// when producing a generic certification signature onto an existing user ID.
	// The identity must be present in the signer Entity.
	SigningIdentity string
	// InsecureAllowUnauthenticatedMessages controls, whether it is tolerated to read
	// encrypted messages without Modification Detection Code (MDC).
	// MDC is mandated by the IETF OpenPGP Crypto Refresh draft and has long been implemented
	// in most OpenPGP implementations. Messages without MDC are considered unnecessarily
	// insecure and should be prevented whenever possible.
	// In case one needs to deal with messages from very old OpenPGP implementations, there
	// might be no other way than to tolerate the missing MDC. Setting this flag, allows this
	// mode of operation. It should be considered a measure of last resort.
	InsecureAllowUnauthenticatedMessages bool
	// InsecureAllowDecryptionWithSigningKeys allows decryption with keys marked as signing keys in the v2 API.
	// This setting is potentially insecure, but it is needed as some libraries
	// ignored key flags when selecting a key for encryption.
	// Not relevant for the v1 API, as all keys were allowed in decryption.
	InsecureAllowDecryptionWithSigningKeys bool
	// KnownNotations is a map of Notation Data names to bools, which controls
	// the notation names that are allowed to be present in critical Notation Data
	// signature subpackets.
	KnownNotations map[string]bool
	// SignatureNotations is a list of Notations to be added to any signatures.
	SignatureNotations []*Notation
	// CheckIntendedRecipients controls, whether the OpenPGP Intended Recipient Fingerprint feature
	// should be enabled for encryption and decryption.
	// (See https://www.ietf.org/archive/id/draft-ietf-openpgp-crypto-refresh-12.html#name-intended-recipient-fingerpr).
	// When the flag is set, encryption produces Intended Recipient Fingerprint signature sub-packets and decryption
	// checks whether the key it was encrypted to is one of the included fingerprints in the signature.
	// If the flag is disabled, no Intended Recipient Fingerprint sub-packets are created or checked.
	// The default behavior, when the config or flag is nil, is to enable the feature.
	CheckIntendedRecipients *bool
	// CacheSessionKey controls if decryption should return the session key used for decryption.
	// If the flag is set, the session key is cached in the message details struct.
	CacheSessionKey bool
	// CheckPacketSequence is a flag that controls if the pgp message reader should strictly check
	// that the packet sequence conforms with the grammar mandated by rfc4880.
	// The default behavior, when the config or flag is nil, is to check the packet sequence.
	CheckPacketSequence *bool
	// NonDeterministicSignaturesViaNotation is a flag to enable randomization of signatures.
	// If true, a salt notation is used to randomize signatures generated by v4 and v5 keys
	// (v6 signatures are always non-deterministic, by design).
	// This protects EdDSA signatures from potentially leaking the secret key in case of faults (i.e. bitflips) which, in principle, could occur
	// during the signing computation. It is added to signatures of any algo for simplicity, and as it may also serve as protection in case of
	// weaknesses in the hash algo, potentially hindering e.g. some chosen-prefix attacks.
	// The default behavior, when the config or flag is nil, is to enable the feature.
	NonDeterministicSignaturesViaNotation *bool

	// InsecureAllowAllKeyFlagsWhenMissing determines how a key without valid key flags is handled.
	// When set to true, a key without flags is treated as if all flags are enabled.
	// This behavior is consistent with GPG.
	InsecureAllowAllKeyFlagsWhenMissing bool
}

func (c *Config) Random() io.Reader {
	if c == nil || c.Rand == nil {
		return rand.Reader
	}
	return c.Rand
}

func (c *Config) Hash() crypto.Hash {
	if c == nil || uint(c.DefaultHash) == 0 {
		return crypto.SHA256
	}
	return c.DefaultHash
}

func (c *Config) Cipher() CipherFunction {
	if c == nil || uint8(c.DefaultCipher) == 0 {
		return CipherAES128
	}
	return c.DefaultCipher
}

func (c *Config) Now() time.Time {
	if c == nil || c.Time == nil {
		return time.Now().Truncate(time.Second)
	}
	return c.Time().Truncate(time.Second)
}

// KeyLifetime returns the validity period of the key.
func (c *Config) KeyLifetime() uint32 {
	if c == nil {
		return 0
	}
	return c.KeyLifetimeSecs
}

// SigLifetime returns the validity period of the signature.
func (c *Config) SigLifetime() uint32 {
	if c == nil {
		return 0
	}
	return c.SigLifetimeSecs
}

func (c *Config) Compression() CompressionAlgo {
	if c == nil {
		return CompressionNone
	}
	return c.DefaultCompressionAlgo
}

func (c *Config) RSAModulusBits() int {
	if c == nil || c.RSABits == 0 {
		return 2048
	}
	return c.RSABits
}

func (c *Config) PublicKeyAlgorithm() PublicKeyAlgorithm {
	if c == nil || c.Algorithm == 0 {
		return PubKeyAlgoRSA
	}
	return c.Algorithm
}

func (c *Config) CurveName() Curve {
	if c == nil || c.Curve == "" {
		return Curve25519
	}
	return c.Curve
}

// Deprecated: The hash iterations should now be queried via the S2K() method.
func (c *Config) PasswordHashIterations() int {
	if c == nil || c.S2KCount == 0 {
		return 0
	}
	return c.S2KCount
}

func (c *Config) S2K() *s2k.Config {
	if c == nil {
		return nil
	}
	// for backwards compatibility
	if c.S2KCount > 0 && c.S2KConfig == nil {
		return &s2k.Config{
			S2KCount: c.S2KCount,
		}
	}
	return c.S2KConfig
}

func (c *Config) AEAD() *AEADConfig {
	if c == nil {
		return nil
	}
	return c.AEADConfig
}

func (c *Config) SigningKey() uint64 {
	if c == nil {
		return 0
	}
	return c.SigningKeyId
}

func (c *Config) SigningUserId() string {
	if c == nil {
		return ""
	}
	return c.SigningIdentity
}

func (c *Config) AllowUnauthenticatedMessages() bool {
	if c == nil {
		return false
	}
	return c.InsecureAllowUnauthenticatedMessages
}

func (c *Config) AllowDecryptionWithSigningKeys() bool {
	if c == nil {
		return false
	}
	return c.InsecureAllowDecryptionWithSigningKeys
}

func (c *Config) KnownNotation(notationName string) bool {
	if c == nil {
		return false
	}
	return c.KnownNotations[notationName]
}

func (c *Config) Notations() []*Notation {
	if c == nil {
		return nil
	}
	return c.SignatureNotations
}

func (c *Config) V6() bool {
	if c == nil {
		return false
	}
	return c.V6Keys
}

func (c *Config) IntendedRecipients() bool {
	if c == nil || c.CheckIntendedRecipients == nil {
		return true
	}
	return *c.CheckIntendedRecipients
}

func (c *Config) RetrieveSessionKey() bool {
	if c == nil {
		return false
	}
	return c.CacheSessionKey
}

func (c *Config) MinimumRSABits() uint16 {
	if c == nil || c.MinRSABits == 0 {
		return 2047
	}
	return c.MinRSABits
}

func (c *Config) RejectPublicKeyAlgorithm(alg PublicKeyAlgorithm) bool {
	var rejectedAlgorithms map[PublicKeyAlgorithm]bool
	if c == nil || c.RejectPublicKeyAlgorithms == nil {
		// Default
		rejectedAlgorithms = defaultRejectPublicKeyAlgorithms
	} else {
		rejectedAlgorithms = c.RejectPublicKeyAlgorithms
	}
	return rejectedAlgorithms[alg]
}

func (c *Config) RejectHashAlgorithm(hash crypto.Hash) bool {
	var rejectedAlgorithms map[crypto.Hash]bool
	if c == nil || c.RejectHashAlgorithms == nil {
		// Default
		rejectedAlgorithms = defaultRejectHashAlgorithms
	} else {
		rejectedAlgorithms = c.RejectHashAlgorithms
	}
	return rejectedAlgorithms[hash]
}

func (c *Config) RejectMessageHashAlgorithm(hash crypto.Hash) bool {
	var rejectedAlgorithms map[crypto.Hash]bool
	if c == nil || c.RejectMessageHashAlgorithms == nil {
		// Default
		rejectedAlgorithms = defaultRejectMessageHashAlgorithms
	} else {
		rejectedAlgorithms = c.RejectMessageHashAlgorithms
	}
	return rejectedAlgorithms[hash]
}

func (c *Config) RejectCurve(curve Curve) bool {
	var rejectedCurve map[Curve]bool
	if c == nil || c.RejectCurves == nil {
		// Default
		rejectedCurve = defaultRejectCurves
	} else {
		rejectedCurve = c.RejectCurves
	}
	return rejectedCurve[curve]
}

func (c *Config) StrictPacketSequence() bool {
	if c == nil || c.CheckPacketSequence == nil {
		return true
	}
	return *c.CheckPacketSequence
}

func (c *Config) RandomizeSignaturesViaNotation() bool {
	if c == nil || c.NonDeterministicSignaturesViaNotation == nil {
		return true
	}
	return *c.NonDeterministicSignaturesViaNotation
}

func (c *Config) AllowAllKeyFlagsWhenMissing() bool {
	if c == nil {
		return false
	}
	return c.InsecureAllowAllKeyFlagsWhenMissing
}

// BoolPointer is a helper function to set a boolean pointer in the Config.
// e.g., config.CheckPacketSequence = BoolPointer(true)
func BoolPointer(value bool) *bool {
	return &value
}
