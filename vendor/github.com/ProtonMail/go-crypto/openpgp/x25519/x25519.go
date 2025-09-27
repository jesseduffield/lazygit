package x25519

import (
	"crypto/sha256"
	"crypto/subtle"
	"io"

	"github.com/ProtonMail/go-crypto/openpgp/aes/keywrap"
	"github.com/ProtonMail/go-crypto/openpgp/errors"
	x25519lib "github.com/cloudflare/circl/dh/x25519"
	"golang.org/x/crypto/hkdf"
)

const (
	hkdfInfo      = "OpenPGP X25519"
	aes128KeySize = 16
	// The size of a public or private key in bytes.
	KeySize = x25519lib.Size
)

type PublicKey struct {
	// Point represents the encoded elliptic curve point of the public key.
	Point []byte
}

type PrivateKey struct {
	PublicKey
	// Secret represents the secret of the private key.
	Secret []byte
}

// NewPrivateKey creates a new empty private key including the public key.
func NewPrivateKey(key PublicKey) *PrivateKey {
	return &PrivateKey{
		PublicKey: key,
	}
}

// Validate validates that the provided public key matches the private key.
func Validate(pk *PrivateKey) (err error) {
	var expectedPublicKey, privateKey x25519lib.Key
	subtle.ConstantTimeCopy(1, privateKey[:], pk.Secret)
	x25519lib.KeyGen(&expectedPublicKey, &privateKey)
	if subtle.ConstantTimeCompare(expectedPublicKey[:], pk.PublicKey.Point) == 0 {
		return errors.KeyInvalidError("x25519: invalid key")
	}
	return nil
}

// GenerateKey generates a new x25519 key pair.
func GenerateKey(rand io.Reader) (*PrivateKey, error) {
	var privateKey, publicKey x25519lib.Key
	privateKeyOut := new(PrivateKey)
	err := generateKey(rand, &privateKey, &publicKey)
	if err != nil {
		return nil, err
	}
	privateKeyOut.PublicKey.Point = publicKey[:]
	privateKeyOut.Secret = privateKey[:]
	return privateKeyOut, nil
}

func generateKey(rand io.Reader, privateKey *x25519lib.Key, publicKey *x25519lib.Key) error {
	maxRounds := 10
	isZero := true
	for round := 0; isZero; round++ {
		if round == maxRounds {
			return errors.InvalidArgumentError("x25519: zero keys only, randomness source might be corrupt")
		}
		_, err := io.ReadFull(rand, privateKey[:])
		if err != nil {
			return err
		}
		isZero = constantTimeIsZero(privateKey[:])
	}
	x25519lib.KeyGen(publicKey, privateKey)
	return nil
}

// Encrypt encrypts a sessionKey with x25519 according to
// the OpenPGP crypto refresh specification section 5.1.6. The function assumes that the
// sessionKey has the correct format and padding according to the specification.
func Encrypt(rand io.Reader, publicKey *PublicKey, sessionKey []byte) (ephemeralPublicKey *PublicKey, encryptedSessionKey []byte, err error) {
	var ephemeralPrivate, ephemeralPublic, staticPublic, shared x25519lib.Key
	// Check that the input static public key has 32 bytes
	if len(publicKey.Point) != KeySize {
		err = errors.KeyInvalidError("x25519: the public key has the wrong size")
		return
	}
	copy(staticPublic[:], publicKey.Point)
	// Generate ephemeral keyPair
	err = generateKey(rand, &ephemeralPrivate, &ephemeralPublic)
	if err != nil {
		return
	}
	// Compute shared key
	ok := x25519lib.Shared(&shared, &ephemeralPrivate, &staticPublic)
	if !ok {
		err = errors.KeyInvalidError("x25519: the public key is a low order point")
		return
	}
	// Derive the encryption key from the shared secret
	encryptionKey := applyHKDF(ephemeralPublic[:], publicKey.Point[:], shared[:])
	ephemeralPublicKey = &PublicKey{
		Point: ephemeralPublic[:],
	}
	// Encrypt the sessionKey with aes key wrapping
	encryptedSessionKey, err = keywrap.Wrap(encryptionKey, sessionKey)
	return
}

// Decrypt decrypts a session key stored in ciphertext with the provided x25519
// private key and ephemeral public key.
func Decrypt(privateKey *PrivateKey, ephemeralPublicKey *PublicKey, ciphertext []byte) (encodedSessionKey []byte, err error) {
	var ephemeralPublic, staticPrivate, shared x25519lib.Key
	// Check that the input ephemeral public key has 32 bytes
	if len(ephemeralPublicKey.Point) != KeySize {
		err = errors.KeyInvalidError("x25519: the public key has the wrong size")
		return
	}
	copy(ephemeralPublic[:], ephemeralPublicKey.Point)
	subtle.ConstantTimeCopy(1, staticPrivate[:], privateKey.Secret)
	// Compute shared key
	ok := x25519lib.Shared(&shared, &staticPrivate, &ephemeralPublic)
	if !ok {
		err = errors.KeyInvalidError("x25519: the ephemeral public key is a low order point")
		return
	}
	// Derive the encryption key from the shared secret
	encryptionKey := applyHKDF(ephemeralPublicKey.Point[:], privateKey.PublicKey.Point[:], shared[:])
	// Decrypt the session key with aes key wrapping
	encodedSessionKey, err = keywrap.Unwrap(encryptionKey, ciphertext)
	return
}

func applyHKDF(ephemeralPublicKey []byte, publicKey []byte, sharedSecret []byte) []byte {
	inputKey := make([]byte, 3*KeySize)
	// ephemeral public key | recipient public key | shared secret
	subtle.ConstantTimeCopy(1, inputKey[:KeySize], ephemeralPublicKey)
	subtle.ConstantTimeCopy(1, inputKey[KeySize:2*KeySize], publicKey)
	subtle.ConstantTimeCopy(1, inputKey[2*KeySize:], sharedSecret)
	hkdfReader := hkdf.New(sha256.New, inputKey, []byte{}, []byte(hkdfInfo))
	encryptionKey := make([]byte, aes128KeySize)
	_, _ = io.ReadFull(hkdfReader, encryptionKey)
	return encryptionKey
}

func constantTimeIsZero(bytes []byte) bool {
	isZero := byte(0)
	for _, b := range bytes {
		isZero |= b
	}
	return isZero == 0
}

// ENCODING/DECODING ciphertexts:

// EncodeFieldsLength returns the length of the ciphertext encoding
// given the encrypted session key.
func EncodedFieldsLength(encryptedSessionKey []byte, v6 bool) int {
	lenCipherFunction := 0
	if !v6 {
		lenCipherFunction = 1
	}
	return KeySize + 1 + len(encryptedSessionKey) + lenCipherFunction
}

// EncodeField encodes x25519 session key encryption fields as
// ephemeral x25519 public key | follow byte length | cipherFunction (v3 only) | encryptedSessionKey
// and writes it to writer.
func EncodeFields(writer io.Writer, ephemeralPublicKey *PublicKey, encryptedSessionKey []byte, cipherFunction byte, v6 bool) (err error) {
	lenAlgorithm := 0
	if !v6 {
		lenAlgorithm = 1
	}
	if _, err = writer.Write(ephemeralPublicKey.Point); err != nil {
		return err
	}
	if _, err = writer.Write([]byte{byte(len(encryptedSessionKey) + lenAlgorithm)}); err != nil {
		return err
	}
	if !v6 {
		if _, err = writer.Write([]byte{cipherFunction}); err != nil {
			return err
		}
	}
	_, err = writer.Write(encryptedSessionKey)
	return err
}

// DecodeField decodes a x25519 session key encryption as
// ephemeral x25519 public key | follow byte length | cipherFunction (v3 only) | encryptedSessionKey.
func DecodeFields(reader io.Reader, v6 bool) (ephemeralPublicKey *PublicKey, encryptedSessionKey []byte, cipherFunction byte, err error) {
	var buf [1]byte
	ephemeralPublicKey = &PublicKey{
		Point: make([]byte, KeySize),
	}
	// 32 octets representing an ephemeral x25519 public key.
	if _, err = io.ReadFull(reader, ephemeralPublicKey.Point); err != nil {
		return nil, nil, 0, err
	}
	// A one-octet size of the following fields.
	if _, err = io.ReadFull(reader, buf[:]); err != nil {
		return nil, nil, 0, err
	}
	followingLen := buf[0]
	// The one-octet algorithm identifier, if it was passed (in the case of a v3 PKESK packet).
	if !v6 {
		if _, err = io.ReadFull(reader, buf[:]); err != nil {
			return nil, nil, 0, err
		}
		cipherFunction = buf[0]
		followingLen -= 1
	}
	// The encrypted session key.
	encryptedSessionKey = make([]byte, followingLen)
	if _, err = io.ReadFull(reader, encryptedSessionKey); err != nil {
		return nil, nil, 0, err
	}
	return ephemeralPublicKey, encryptedSessionKey, cipherFunction, nil
}
