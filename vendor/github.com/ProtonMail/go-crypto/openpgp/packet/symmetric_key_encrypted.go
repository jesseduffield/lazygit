// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"bytes"
	"crypto/cipher"
	"crypto/sha256"
	"io"
	"strconv"

	"github.com/ProtonMail/go-crypto/openpgp/errors"
	"github.com/ProtonMail/go-crypto/openpgp/s2k"
	"golang.org/x/crypto/hkdf"
)

// This is the largest session key that we'll support. Since at most 256-bit cipher
// is supported in OpenPGP, this is large enough to contain also the auth tag.
const maxSessionKeySizeInBytes = 64

// SymmetricKeyEncrypted represents a passphrase protected session key. See RFC
// 4880, section 5.3.
type SymmetricKeyEncrypted struct {
	Version      int
	CipherFunc   CipherFunction
	Mode         AEADMode
	s2k          func(out, in []byte)
	iv           []byte
	encryptedKey []byte // Contains also the authentication tag for AEAD
}

// parse parses an SymmetricKeyEncrypted packet as specified in
// https://www.ietf.org/archive/id/draft-ietf-openpgp-crypto-refresh-07.html#name-symmetric-key-encrypted-ses
func (ske *SymmetricKeyEncrypted) parse(r io.Reader) error {
	var buf [1]byte

	// Version
	if _, err := readFull(r, buf[:]); err != nil {
		return err
	}
	ske.Version = int(buf[0])
	if ske.Version != 4 && ske.Version != 5 && ske.Version != 6 {
		return errors.UnsupportedError("unknown SymmetricKeyEncrypted version")
	}

	if V5Disabled && ske.Version == 5 {
		return errors.UnsupportedError("support for parsing v5 entities is disabled; build with `-tags v5` if needed")
	}

	if ske.Version > 5 {
		// Scalar octet count
		if _, err := readFull(r, buf[:]); err != nil {
			return err
		}
	}

	// Cipher function
	if _, err := readFull(r, buf[:]); err != nil {
		return err
	}
	ske.CipherFunc = CipherFunction(buf[0])
	if !ske.CipherFunc.IsSupported() {
		return errors.UnsupportedError("unknown cipher: " + strconv.Itoa(int(buf[0])))
	}

	if ske.Version >= 5 {
		// AEAD mode
		if _, err := readFull(r, buf[:]); err != nil {
			return errors.StructuralError("cannot read AEAD octet from packet")
		}
		ske.Mode = AEADMode(buf[0])
	}

	if ske.Version > 5 {
		// Scalar octet count
		if _, err := readFull(r, buf[:]); err != nil {
			return err
		}
	}

	var err error
	if ske.s2k, err = s2k.Parse(r); err != nil {
		if _, ok := err.(errors.ErrDummyPrivateKey); ok {
			return errors.UnsupportedError("missing key GNU extension in session key")
		}
		return err
	}

	if ske.Version >= 5 {
		// AEAD IV
		iv := make([]byte, ske.Mode.IvLength())
		_, err := readFull(r, iv)
		if err != nil {
			return errors.StructuralError("cannot read AEAD IV")
		}

		ske.iv = iv
	}

	encryptedKey := make([]byte, maxSessionKeySizeInBytes)
	// The session key may follow. We just have to try and read to find
	// out. If it exists then we limit it to maxSessionKeySizeInBytes.
	n, err := readFull(r, encryptedKey)
	if err != nil && err != io.ErrUnexpectedEOF {
		return err
	}

	if n != 0 {
		if n == maxSessionKeySizeInBytes {
			return errors.UnsupportedError("oversized encrypted session key")
		}
		ske.encryptedKey = encryptedKey[:n]
	}
	return nil
}

// Decrypt attempts to decrypt an encrypted session key and returns the key and
// the cipher to use when decrypting a subsequent Symmetrically Encrypted Data
// packet.
func (ske *SymmetricKeyEncrypted) Decrypt(passphrase []byte) ([]byte, CipherFunction, error) {
	key := make([]byte, ske.CipherFunc.KeySize())
	ske.s2k(key, passphrase)
	if len(ske.encryptedKey) == 0 {
		return key, ske.CipherFunc, nil
	}
	switch ske.Version {
	case 4:
		plaintextKey, cipherFunc, err := ske.decryptV4(key)
		return plaintextKey, cipherFunc, err
	case 5, 6:
		plaintextKey, err := ske.aeadDecrypt(ske.Version, key)
		return plaintextKey, CipherFunction(0), err
	}
	err := errors.UnsupportedError("unknown SymmetricKeyEncrypted version")
	return nil, CipherFunction(0), err
}

func (ske *SymmetricKeyEncrypted) decryptV4(key []byte) ([]byte, CipherFunction, error) {
	// the IV is all zeros
	iv := make([]byte, ske.CipherFunc.blockSize())
	c := cipher.NewCFBDecrypter(ske.CipherFunc.new(key), iv)
	plaintextKey := make([]byte, len(ske.encryptedKey))
	c.XORKeyStream(plaintextKey, ske.encryptedKey)
	cipherFunc := CipherFunction(plaintextKey[0])
	if cipherFunc.blockSize() == 0 {
		return nil, ske.CipherFunc, errors.UnsupportedError(
			"unknown cipher: " + strconv.Itoa(int(cipherFunc)))
	}
	plaintextKey = plaintextKey[1:]
	if len(plaintextKey) != cipherFunc.KeySize() {
		return nil, cipherFunc, errors.StructuralError(
			"length of decrypted key not equal to cipher keysize")
	}
	return plaintextKey, cipherFunc, nil
}

func (ske *SymmetricKeyEncrypted) aeadDecrypt(version int, key []byte) ([]byte, error) {
	adata := []byte{0xc3, byte(version), byte(ske.CipherFunc), byte(ske.Mode)}
	aead := getEncryptedKeyAeadInstance(ske.CipherFunc, ske.Mode, key, adata, version)

	plaintextKey, err := aead.Open(nil, ske.iv, ske.encryptedKey, adata)
	if err != nil {
		return nil, err
	}
	return plaintextKey, nil
}

// SerializeSymmetricKeyEncrypted serializes a symmetric key packet to w.
// The packet contains a random session key, encrypted by a key derived from
// the given passphrase. The session key is returned and must be passed to
// SerializeSymmetricallyEncrypted.
// If config is nil, sensible defaults will be used.
func SerializeSymmetricKeyEncrypted(w io.Writer, passphrase []byte, config *Config) (key []byte, err error) {
	cipherFunc := config.Cipher()

	sessionKey := make([]byte, cipherFunc.KeySize())
	_, err = io.ReadFull(config.Random(), sessionKey)
	if err != nil {
		return
	}

	err = SerializeSymmetricKeyEncryptedReuseKey(w, sessionKey, passphrase, config)
	if err != nil {
		return
	}

	key = sessionKey
	return
}

// SerializeSymmetricKeyEncryptedReuseKey serializes a symmetric key packet to w.
// The packet contains the given session key, encrypted by a key derived from
// the given passphrase. The returned session key must be passed to
// SerializeSymmetricallyEncrypted.
// If config is nil, sensible defaults will be used.
// Deprecated: Use SerializeSymmetricKeyEncryptedAEADReuseKey instead.
func SerializeSymmetricKeyEncryptedReuseKey(w io.Writer, sessionKey []byte, passphrase []byte, config *Config) (err error) {
	return SerializeSymmetricKeyEncryptedAEADReuseKey(w, sessionKey, passphrase, config.AEAD() != nil, config)
}

// SerializeSymmetricKeyEncryptedAEADReuseKey serializes a symmetric key packet to w.
// The packet contains the given session key, encrypted by a key derived from
// the given passphrase. The returned session key must be passed to
// SerializeSymmetricallyEncrypted.
// If aeadSupported is set, SKESK v6 is used, otherwise v4.
// Note: aeadSupported MUST match the value passed to SerializeSymmetricallyEncrypted.
// If config is nil, sensible defaults will be used.
func SerializeSymmetricKeyEncryptedAEADReuseKey(w io.Writer, sessionKey []byte, passphrase []byte, aeadSupported bool, config *Config) (err error) {
	var version int
	if aeadSupported {
		version = 6
	} else {
		version = 4
	}
	cipherFunc := config.Cipher()
	// cipherFunc must be AES
	if !cipherFunc.IsSupported() || cipherFunc < CipherAES128 || cipherFunc > CipherAES256 {
		return errors.UnsupportedError("unsupported cipher: " + strconv.Itoa(int(cipherFunc)))
	}

	keySize := cipherFunc.KeySize()
	s2kBuf := new(bytes.Buffer)
	keyEncryptingKey := make([]byte, keySize)
	// s2k.Serialize salts and stretches the passphrase, and writes the
	// resulting key to keyEncryptingKey and the s2k descriptor to s2kBuf.
	err = s2k.Serialize(s2kBuf, keyEncryptingKey, config.Random(), passphrase, config.S2K())
	if err != nil {
		return
	}
	s2kBytes := s2kBuf.Bytes()

	var packetLength int
	switch version {
	case 4:
		packetLength = 2 /* header */ + len(s2kBytes) + 1 /* cipher type */ + keySize
	case 5, 6:
		ivLen := config.AEAD().Mode().IvLength()
		tagLen := config.AEAD().Mode().TagLength()
		packetLength = 3 + len(s2kBytes) + ivLen + keySize + tagLen
	}
	if version > 5 {
		packetLength += 2 // additional octet count fields
	}

	err = serializeHeader(w, packetTypeSymmetricKeyEncrypted, packetLength)
	if err != nil {
		return
	}

	// Symmetric Key Encrypted Version
	buf := []byte{byte(version)}

	if version > 5 {
		// Scalar octet count
		buf = append(buf, byte(3+len(s2kBytes)+config.AEAD().Mode().IvLength()))
	}

	// Cipher function
	buf = append(buf, byte(cipherFunc))

	if version >= 5 {
		// AEAD mode
		buf = append(buf, byte(config.AEAD().Mode()))
	}
	if version > 5 {
		// Scalar octet count
		buf = append(buf, byte(len(s2kBytes)))
	}
	_, err = w.Write(buf)
	if err != nil {
		return
	}
	_, err = w.Write(s2kBytes)
	if err != nil {
		return
	}

	switch version {
	case 4:
		iv := make([]byte, cipherFunc.blockSize())
		c := cipher.NewCFBEncrypter(cipherFunc.new(keyEncryptingKey), iv)
		encryptedCipherAndKey := make([]byte, keySize+1)
		c.XORKeyStream(encryptedCipherAndKey, buf[1:])
		c.XORKeyStream(encryptedCipherAndKey[1:], sessionKey)
		_, err = w.Write(encryptedCipherAndKey)
		if err != nil {
			return
		}
	case 5, 6:
		mode := config.AEAD().Mode()
		adata := []byte{0xc3, byte(version), byte(cipherFunc), byte(mode)}
		aead := getEncryptedKeyAeadInstance(cipherFunc, mode, keyEncryptingKey, adata, version)

		// Sample iv using random reader
		iv := make([]byte, config.AEAD().Mode().IvLength())
		_, err = io.ReadFull(config.Random(), iv)
		if err != nil {
			return
		}
		// Seal and write (encryptedData includes auth. tag)

		encryptedData := aead.Seal(nil, iv, sessionKey, adata)
		_, err = w.Write(iv)
		if err != nil {
			return
		}
		_, err = w.Write(encryptedData)
		if err != nil {
			return
		}
	}

	return
}

func getEncryptedKeyAeadInstance(c CipherFunction, mode AEADMode, inputKey, associatedData []byte, version int) (aead cipher.AEAD) {
	var blockCipher cipher.Block
	if version > 5 {
		hkdfReader := hkdf.New(sha256.New, inputKey, []byte{}, associatedData)

		encryptionKey := make([]byte, c.KeySize())
		_, _ = readFull(hkdfReader, encryptionKey)

		blockCipher = c.new(encryptionKey)
	} else {
		blockCipher = c.new(inputKey)
	}
	return mode.new(blockCipher)
}
