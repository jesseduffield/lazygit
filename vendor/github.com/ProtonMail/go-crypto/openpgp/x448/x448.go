package x448

import (
	"crypto/sha512"
	"crypto/subtle"
	"io"

	"github.com/ProtonMail/go-crypto/openpgp/aes/keywrap"
	"github.com/ProtonMail/go-crypto/openpgp/errors"
	x448lib "github.com/cloudflare/circl/dh/x448"
	"golang.org/x/crypto/hkdf"
)

const (
	hkdfInfo      = "OpenPGP X448"
	aes256KeySize = 32
	// The size of a public or private key in bytes.
	KeySize = x448lib.Size
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

// Validate validates that the provided public key matches
// the private key.
func Validate(pk *PrivateKey) (err error) {
	var expectedPublicKey, privateKey x448lib.Key
	subtle.ConstantTimeCopy(1, privateKey[:], pk.Secret)
	x448lib.KeyGen(&expectedPublicKey, &privateKey)
	if subtle.ConstantTimeCompare(expectedPublicKey[:], pk.PublicKey.Point) == 0 {
		return errors.KeyInvalidError("x448: invalid key")
	}
	return nil
}

// GenerateKey generates a new x448 key pair.
func GenerateKey(rand io.Reader) (*PrivateKey, error) {
	var privateKey, publicKey x448lib.Key
	privateKeyOut := new(PrivateKey)
	err := generateKey(rand, &privateKey, &publicKey)
	if err != nil {
		return nil, err
	}
	privateKeyOut.PublicKey.Point = publicKey[:]
	privateKeyOut.Secret = privateKey[:]
	return privateKeyOut, nil
}

func generateKey(rand io.Reader, privateKey *x448lib.Key, publicKey *x448lib.Key) error {
	maxRounds := 10
	isZero := true
	for round := 0; isZero; round++ {
		if round == maxRounds {
			return errors.InvalidArgumentError("x448: zero keys only, randomness source might be corrupt")
		}
		_, err := io.ReadFull(rand, privateKey[:])
		if err != nil {
			return err
		}
		isZero = constantTimeIsZero(privateKey[:])
	}
	x448lib.KeyGen(publicKey, privateKey)
	return nil
}

// Encrypt encrypts a sessionKey with x448 according to
// the OpenPGP crypto refresh specification section 5.1.7. The function assumes that the
// sessionKey has the correct format and padding according to the specification.
func Encrypt(rand io.Reader, publicKey *PublicKey, sessionKey []byte) (ephemeralPublicKey *PublicKey, encryptedSessionKey []byte, err error) {
	var ephemeralPrivate, ephemeralPublic, staticPublic, shared x448lib.Key
	// Check that the input static public key has 56 bytes.
	if len(publicKey.Point) != KeySize {
		err = errors.KeyInvalidError("x448: the public key has the wrong size")
		return nil, nil, err
	}
	copy(staticPublic[:], publicKey.Point)
	// Generate ephemeral keyPair.
	if err = generateKey(rand, &ephemeralPrivate, &ephemeralPublic); err != nil {
		return nil, nil, err
	}
	// Compute shared key.
	ok := x448lib.Shared(&shared, &ephemeralPrivate, &staticPublic)
	if !ok {
		err = errors.KeyInvalidError("x448: the public key is a low order point")
		return nil, nil, err
	}
	// Derive the encryption key from the shared secret.
	encryptionKey := applyHKDF(ephemeralPublic[:], publicKey.Point[:], shared[:])
	ephemeralPublicKey = &PublicKey{
		Point: ephemeralPublic[:],
	}
	// Encrypt the sessionKey with aes key wrapping.
	encryptedSessionKey, err = keywrap.Wrap(encryptionKey, sessionKey)
	if err != nil {
		return nil, nil, err
	}
	return ephemeralPublicKey, encryptedSessionKey, nil
}

// Decrypt decrypts a session key stored in ciphertext with the provided x448
// private key and ephemeral public key.
func Decrypt(privateKey *PrivateKey, ephemeralPublicKey *PublicKey, ciphertext []byte) (encodedSessionKey []byte, err error) {
	var ephemeralPublic, staticPrivate, shared x448lib.Key
	// Check that the input ephemeral public key has 56 bytes.
	if len(ephemeralPublicKey.Point) != KeySize {
		err = errors.KeyInvalidError("x448: the public key has the wrong size")
		return nil, err
	}
	copy(ephemeralPublic[:], ephemeralPublicKey.Point)
	subtle.ConstantTimeCopy(1, staticPrivate[:], privateKey.Secret)
	// Compute shared key.
	ok := x448lib.Shared(&shared, &staticPrivate, &ephemeralPublic)
	if !ok {
		err = errors.KeyInvalidError("x448: the ephemeral public key is a low order point")
		return nil, err
	}
	// Derive the encryption key from the shared secret.
	encryptionKey := applyHKDF(ephemeralPublicKey.Point[:], privateKey.PublicKey.Point[:], shared[:])
	// Decrypt the session key with aes key wrapping.
	encodedSessionKey, err = keywrap.Unwrap(encryptionKey, ciphertext)
	if err != nil {
		return nil, err
	}
	return encodedSessionKey, nil
}

func applyHKDF(ephemeralPublicKey []byte, publicKey []byte, sharedSecret []byte) []byte {
	inputKey := make([]byte, 3*KeySize)
	// ephemeral public key | recipient public key | shared secret.
	subtle.ConstantTimeCopy(1, inputKey[:KeySize], ephemeralPublicKey)
	subtle.ConstantTimeCopy(1, inputKey[KeySize:2*KeySize], publicKey)
	subtle.ConstantTimeCopy(1, inputKey[2*KeySize:], sharedSecret)
	hkdfReader := hkdf.New(sha512.New, inputKey, []byte{}, []byte(hkdfInfo))
	encryptionKey := make([]byte, aes256KeySize)
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

// EncodeField encodes x448 session key encryption fields as
// ephemeral x448 public key | follow byte length | cipherFunction (v3 only) | encryptedSessionKey
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
	if _, err = writer.Write(encryptedSessionKey); err != nil {
		return err
	}
	return nil
}

// DecodeField decodes a x448 session key encryption as
// ephemeral x448 public key | follow byte length | cipherFunction (v3 only) | encryptedSessionKey.
func DecodeFields(reader io.Reader, v6 bool) (ephemeralPublicKey *PublicKey, encryptedSessionKey []byte, cipherFunction byte, err error) {
	var buf [1]byte
	ephemeralPublicKey = &PublicKey{
		Point: make([]byte, KeySize),
	}
	// 56 octets representing an ephemeral x448 public key.
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
