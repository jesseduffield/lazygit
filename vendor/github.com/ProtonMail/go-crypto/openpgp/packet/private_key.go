// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"bytes"
	"crypto"
	"crypto/cipher"
	"crypto/dsa"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp/ecdh"
	"github.com/ProtonMail/go-crypto/openpgp/ecdsa"
	"github.com/ProtonMail/go-crypto/openpgp/ed25519"
	"github.com/ProtonMail/go-crypto/openpgp/ed448"
	"github.com/ProtonMail/go-crypto/openpgp/eddsa"
	"github.com/ProtonMail/go-crypto/openpgp/elgamal"
	"github.com/ProtonMail/go-crypto/openpgp/errors"
	"github.com/ProtonMail/go-crypto/openpgp/internal/encoding"
	"github.com/ProtonMail/go-crypto/openpgp/s2k"
	"github.com/ProtonMail/go-crypto/openpgp/x25519"
	"github.com/ProtonMail/go-crypto/openpgp/x448"
	"golang.org/x/crypto/hkdf"
)

// PrivateKey represents a possibly encrypted private key. See RFC 4880,
// section 5.5.3.
type PrivateKey struct {
	PublicKey
	Encrypted     bool // if true then the private key is unavailable until Decrypt has been called.
	encryptedData []byte
	cipher        CipherFunction
	s2k           func(out, in []byte)
	aead          AEADMode // only relevant if S2KAEAD is enabled
	// An *{rsa|dsa|elgamal|ecdh|ecdsa|ed25519|ed448}.PrivateKey or
	// crypto.Signer/crypto.Decrypter (Decryptor RSA only).
	PrivateKey interface{}
	iv         []byte

	// Type of encryption of the S2K packet
	// Allowed values are 0 (Not encrypted), 253 (AEAD), 254 (SHA1), or
	// 255 (2-byte checksum)
	s2kType S2KType
	// Full parameters of the S2K packet
	s2kParams *s2k.Params
}

// S2KType s2k packet type
type S2KType uint8

const (
	// S2KNON unencrypt
	S2KNON S2KType = 0
	// S2KAEAD use authenticated encryption
	S2KAEAD S2KType = 253
	// S2KSHA1 sha1 sum check
	S2KSHA1 S2KType = 254
	// S2KCHECKSUM sum check
	S2KCHECKSUM S2KType = 255
)

func NewRSAPrivateKey(creationTime time.Time, priv *rsa.PrivateKey) *PrivateKey {
	pk := new(PrivateKey)
	pk.PublicKey = *NewRSAPublicKey(creationTime, &priv.PublicKey)
	pk.PrivateKey = priv
	return pk
}

func NewDSAPrivateKey(creationTime time.Time, priv *dsa.PrivateKey) *PrivateKey {
	pk := new(PrivateKey)
	pk.PublicKey = *NewDSAPublicKey(creationTime, &priv.PublicKey)
	pk.PrivateKey = priv
	return pk
}

func NewElGamalPrivateKey(creationTime time.Time, priv *elgamal.PrivateKey) *PrivateKey {
	pk := new(PrivateKey)
	pk.PublicKey = *NewElGamalPublicKey(creationTime, &priv.PublicKey)
	pk.PrivateKey = priv
	return pk
}

func NewECDSAPrivateKey(creationTime time.Time, priv *ecdsa.PrivateKey) *PrivateKey {
	pk := new(PrivateKey)
	pk.PublicKey = *NewECDSAPublicKey(creationTime, &priv.PublicKey)
	pk.PrivateKey = priv
	return pk
}

func NewEdDSAPrivateKey(creationTime time.Time, priv *eddsa.PrivateKey) *PrivateKey {
	pk := new(PrivateKey)
	pk.PublicKey = *NewEdDSAPublicKey(creationTime, &priv.PublicKey)
	pk.PrivateKey = priv
	return pk
}

func NewECDHPrivateKey(creationTime time.Time, priv *ecdh.PrivateKey) *PrivateKey {
	pk := new(PrivateKey)
	pk.PublicKey = *NewECDHPublicKey(creationTime, &priv.PublicKey)
	pk.PrivateKey = priv
	return pk
}

func NewX25519PrivateKey(creationTime time.Time, priv *x25519.PrivateKey) *PrivateKey {
	pk := new(PrivateKey)
	pk.PublicKey = *NewX25519PublicKey(creationTime, &priv.PublicKey)
	pk.PrivateKey = priv
	return pk
}

func NewX448PrivateKey(creationTime time.Time, priv *x448.PrivateKey) *PrivateKey {
	pk := new(PrivateKey)
	pk.PublicKey = *NewX448PublicKey(creationTime, &priv.PublicKey)
	pk.PrivateKey = priv
	return pk
}

func NewEd25519PrivateKey(creationTime time.Time, priv *ed25519.PrivateKey) *PrivateKey {
	pk := new(PrivateKey)
	pk.PublicKey = *NewEd25519PublicKey(creationTime, &priv.PublicKey)
	pk.PrivateKey = priv
	return pk
}

func NewEd448PrivateKey(creationTime time.Time, priv *ed448.PrivateKey) *PrivateKey {
	pk := new(PrivateKey)
	pk.PublicKey = *NewEd448PublicKey(creationTime, &priv.PublicKey)
	pk.PrivateKey = priv
	return pk
}

// NewSignerPrivateKey creates a PrivateKey from a crypto.Signer that
// implements RSA, ECDSA or EdDSA.
func NewSignerPrivateKey(creationTime time.Time, signer interface{}) *PrivateKey {
	pk := new(PrivateKey)
	// In general, the public Keys should be used as pointers. We still
	// type-switch on the values, for backwards-compatibility.
	switch pubkey := signer.(type) {
	case *rsa.PrivateKey:
		pk.PublicKey = *NewRSAPublicKey(creationTime, &pubkey.PublicKey)
	case rsa.PrivateKey:
		pk.PublicKey = *NewRSAPublicKey(creationTime, &pubkey.PublicKey)
	case *ecdsa.PrivateKey:
		pk.PublicKey = *NewECDSAPublicKey(creationTime, &pubkey.PublicKey)
	case ecdsa.PrivateKey:
		pk.PublicKey = *NewECDSAPublicKey(creationTime, &pubkey.PublicKey)
	case *eddsa.PrivateKey:
		pk.PublicKey = *NewEdDSAPublicKey(creationTime, &pubkey.PublicKey)
	case eddsa.PrivateKey:
		pk.PublicKey = *NewEdDSAPublicKey(creationTime, &pubkey.PublicKey)
	case *ed25519.PrivateKey:
		pk.PublicKey = *NewEd25519PublicKey(creationTime, &pubkey.PublicKey)
	case ed25519.PrivateKey:
		pk.PublicKey = *NewEd25519PublicKey(creationTime, &pubkey.PublicKey)
	case *ed448.PrivateKey:
		pk.PublicKey = *NewEd448PublicKey(creationTime, &pubkey.PublicKey)
	case ed448.PrivateKey:
		pk.PublicKey = *NewEd448PublicKey(creationTime, &pubkey.PublicKey)
	default:
		panic("openpgp: unknown signer type in NewSignerPrivateKey")
	}
	pk.PrivateKey = signer
	return pk
}

// NewDecrypterPrivateKey creates a PrivateKey from a *{rsa|elgamal|ecdh|x25519|x448}.PrivateKey.
func NewDecrypterPrivateKey(creationTime time.Time, decrypter interface{}) *PrivateKey {
	pk := new(PrivateKey)
	switch priv := decrypter.(type) {
	case *rsa.PrivateKey:
		pk.PublicKey = *NewRSAPublicKey(creationTime, &priv.PublicKey)
	case *elgamal.PrivateKey:
		pk.PublicKey = *NewElGamalPublicKey(creationTime, &priv.PublicKey)
	case *ecdh.PrivateKey:
		pk.PublicKey = *NewECDHPublicKey(creationTime, &priv.PublicKey)
	case *x25519.PrivateKey:
		pk.PublicKey = *NewX25519PublicKey(creationTime, &priv.PublicKey)
	case *x448.PrivateKey:
		pk.PublicKey = *NewX448PublicKey(creationTime, &priv.PublicKey)
	default:
		panic("openpgp: unknown decrypter type in NewDecrypterPrivateKey")
	}
	pk.PrivateKey = decrypter
	return pk
}

func (pk *PrivateKey) parse(r io.Reader) (err error) {
	err = (&pk.PublicKey).parse(r)
	if err != nil {
		return
	}
	v5 := pk.PublicKey.Version == 5
	v6 := pk.PublicKey.Version == 6

	if V5Disabled && v5 {
		return errors.UnsupportedError("support for parsing v5 entities is disabled; build with `-tags v5` if needed")
	}

	var buf [1]byte
	_, err = readFull(r, buf[:])
	if err != nil {
		return
	}
	pk.s2kType = S2KType(buf[0])
	var optCount [1]byte
	if v5 || (v6 && pk.s2kType != S2KNON) {
		if _, err = readFull(r, optCount[:]); err != nil {
			return
		}
	}

	switch pk.s2kType {
	case S2KNON:
		pk.s2k = nil
		pk.Encrypted = false
	case S2KSHA1, S2KCHECKSUM, S2KAEAD:
		if (v5 || v6) && pk.s2kType == S2KCHECKSUM {
			return errors.StructuralError(fmt.Sprintf("wrong s2k identifier for version %d", pk.Version))
		}
		_, err = readFull(r, buf[:])
		if err != nil {
			return
		}
		pk.cipher = CipherFunction(buf[0])
		if pk.cipher != 0 && !pk.cipher.IsSupported() {
			return errors.UnsupportedError("unsupported cipher function in private key")
		}
		// [Optional] If string-to-key usage octet was 253,
		// a one-octet AEAD algorithm.
		if pk.s2kType == S2KAEAD {
			_, err = readFull(r, buf[:])
			if err != nil {
				return
			}
			pk.aead = AEADMode(buf[0])
			if !pk.aead.IsSupported() {
				return errors.UnsupportedError("unsupported aead mode in private key")
			}
		}

		// [Optional] Only for a version 6 packet,
		// and if string-to-key usage octet was 255, 254, or 253,
		// an one-octet count of the following field.
		if v6 {
			_, err = readFull(r, buf[:])
			if err != nil {
				return
			}
		}

		pk.s2kParams, err = s2k.ParseIntoParams(r)
		if err != nil {
			return
		}
		if pk.s2kParams.Dummy() {
			return
		}
		if pk.s2kParams.Mode() == s2k.Argon2S2K && pk.s2kType != S2KAEAD {
			return errors.StructuralError("using Argon2 S2K without AEAD is not allowed")
		}
		if pk.s2kParams.Mode() == s2k.SimpleS2K && pk.Version == 6 {
			return errors.StructuralError("using Simple S2K with version 6 keys is not allowed")
		}
		pk.s2k, err = pk.s2kParams.Function()
		if err != nil {
			return
		}
		pk.Encrypted = true
	default:
		return errors.UnsupportedError("deprecated s2k function in private key")
	}

	if pk.Encrypted {
		var ivSize int
		// If the S2K usage octet was 253, the IV is of the size expected by the AEAD mode,
		// unless it's a version 5 key, in which case it's the size of the symmetric cipher's block size.
		// For all other S2K modes, it's always the block size.
		if !v5 && pk.s2kType == S2KAEAD {
			ivSize = pk.aead.IvLength()
		} else {
			ivSize = pk.cipher.blockSize()
		}

		if ivSize == 0 {
			return errors.UnsupportedError("unsupported cipher in private key: " + strconv.Itoa(int(pk.cipher)))
		}
		pk.iv = make([]byte, ivSize)
		_, err = readFull(r, pk.iv)
		if err != nil {
			return
		}
		if v5 && pk.s2kType == S2KAEAD {
			pk.iv = pk.iv[:pk.aead.IvLength()]
		}
	}

	var privateKeyData []byte
	if v5 {
		var n [4]byte /* secret material four octet count */
		_, err = readFull(r, n[:])
		if err != nil {
			return
		}
		count := uint32(uint32(n[0])<<24 | uint32(n[1])<<16 | uint32(n[2])<<8 | uint32(n[3]))
		if !pk.Encrypted {
			count = count + 2 /* two octet checksum */
		}
		privateKeyData = make([]byte, count)
		_, err = readFull(r, privateKeyData)
		if err != nil {
			return
		}
	} else {
		privateKeyData, err = io.ReadAll(r)
		if err != nil {
			return
		}
	}
	if !pk.Encrypted {
		if len(privateKeyData) < 2 {
			return errors.StructuralError("truncated private key data")
		}
		if pk.Version != 6 {
			// checksum
			var sum uint16
			for i := 0; i < len(privateKeyData)-2; i++ {
				sum += uint16(privateKeyData[i])
			}
			if privateKeyData[len(privateKeyData)-2] != uint8(sum>>8) ||
				privateKeyData[len(privateKeyData)-1] != uint8(sum) {
				return errors.StructuralError("private key checksum failure")
			}
			privateKeyData = privateKeyData[:len(privateKeyData)-2]
			return pk.parsePrivateKey(privateKeyData)
		} else {
			// No checksum
			return pk.parsePrivateKey(privateKeyData)
		}
	}

	pk.encryptedData = privateKeyData
	return
}

// Dummy returns true if the private key is a dummy key. This is a GNU extension.
func (pk *PrivateKey) Dummy() bool {
	return pk.s2kParams.Dummy()
}

func mod64kHash(d []byte) uint16 {
	var h uint16
	for _, b := range d {
		h += uint16(b)
	}
	return h
}

func (pk *PrivateKey) Serialize(w io.Writer) (err error) {
	contents := bytes.NewBuffer(nil)
	err = pk.PublicKey.serializeWithoutHeaders(contents)
	if err != nil {
		return
	}
	if _, err = contents.Write([]byte{uint8(pk.s2kType)}); err != nil {
		return
	}

	optional := bytes.NewBuffer(nil)
	if pk.Encrypted || pk.Dummy() {
		// [Optional] If string-to-key usage octet was 255, 254, or 253,
		// a one-octet symmetric encryption algorithm.
		if _, err = optional.Write([]byte{uint8(pk.cipher)}); err != nil {
			return
		}
		// [Optional] If string-to-key usage octet was 253,
		// a one-octet AEAD algorithm.
		if pk.s2kType == S2KAEAD {
			if _, err = optional.Write([]byte{uint8(pk.aead)}); err != nil {
				return
			}
		}

		s2kBuffer := bytes.NewBuffer(nil)
		if err := pk.s2kParams.Serialize(s2kBuffer); err != nil {
			return err
		}
		// [Optional] Only for a version 6 packet, and if string-to-key
		// usage octet was 255, 254, or 253, an one-octet
		// count of the following field.
		if pk.Version == 6 {
			if _, err = optional.Write([]byte{uint8(s2kBuffer.Len())}); err != nil {
				return
			}
		}
		// [Optional] If string-to-key usage octet was 255, 254, or 253,
		// a string-to-key (S2K) specifier. The length of the string-to-key specifier
		// depends on its type
		if _, err = io.Copy(optional, s2kBuffer); err != nil {
			return
		}

		// IV
		if pk.Encrypted {
			if _, err = optional.Write(pk.iv); err != nil {
				return
			}
			if pk.Version == 5 && pk.s2kType == S2KAEAD {
				// Add padding for version 5
				padding := make([]byte, pk.cipher.blockSize()-len(pk.iv))
				if _, err = optional.Write(padding); err != nil {
					return
				}
			}
		}
	}
	if pk.Version == 5 || (pk.Version == 6 && pk.s2kType != S2KNON) {
		contents.Write([]byte{uint8(optional.Len())})
	}

	if _, err := io.Copy(contents, optional); err != nil {
		return err
	}

	if !pk.Dummy() {
		l := 0
		var priv []byte
		if !pk.Encrypted {
			buf := bytes.NewBuffer(nil)
			err = pk.serializePrivateKey(buf)
			if err != nil {
				return err
			}
			l = buf.Len()
			if pk.Version != 6 {
				checksum := mod64kHash(buf.Bytes())
				buf.Write([]byte{byte(checksum >> 8), byte(checksum)})
			}
			priv = buf.Bytes()
		} else {
			priv, l = pk.encryptedData, len(pk.encryptedData)
		}

		if pk.Version == 5 {
			contents.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)})
		}
		contents.Write(priv)
	}

	ptype := packetTypePrivateKey
	if pk.IsSubkey {
		ptype = packetTypePrivateSubkey
	}
	err = serializeHeader(w, ptype, contents.Len())
	if err != nil {
		return
	}
	_, err = io.Copy(w, contents)
	if err != nil {
		return
	}
	return
}

func serializeRSAPrivateKey(w io.Writer, priv *rsa.PrivateKey) error {
	if _, err := w.Write(new(encoding.MPI).SetBig(priv.D).EncodedBytes()); err != nil {
		return err
	}
	if _, err := w.Write(new(encoding.MPI).SetBig(priv.Primes[1]).EncodedBytes()); err != nil {
		return err
	}
	if _, err := w.Write(new(encoding.MPI).SetBig(priv.Primes[0]).EncodedBytes()); err != nil {
		return err
	}
	_, err := w.Write(new(encoding.MPI).SetBig(priv.Precomputed.Qinv).EncodedBytes())
	return err
}

func serializeDSAPrivateKey(w io.Writer, priv *dsa.PrivateKey) error {
	_, err := w.Write(new(encoding.MPI).SetBig(priv.X).EncodedBytes())
	return err
}

func serializeElGamalPrivateKey(w io.Writer, priv *elgamal.PrivateKey) error {
	_, err := w.Write(new(encoding.MPI).SetBig(priv.X).EncodedBytes())
	return err
}

func serializeECDSAPrivateKey(w io.Writer, priv *ecdsa.PrivateKey) error {
	_, err := w.Write(encoding.NewMPI(priv.MarshalIntegerSecret()).EncodedBytes())
	return err
}

func serializeEdDSAPrivateKey(w io.Writer, priv *eddsa.PrivateKey) error {
	_, err := w.Write(encoding.NewMPI(priv.MarshalByteSecret()).EncodedBytes())
	return err
}

func serializeECDHPrivateKey(w io.Writer, priv *ecdh.PrivateKey) error {
	_, err := w.Write(encoding.NewMPI(priv.MarshalByteSecret()).EncodedBytes())
	return err
}

func serializeX25519PrivateKey(w io.Writer, priv *x25519.PrivateKey) error {
	_, err := w.Write(priv.Secret)
	return err
}

func serializeX448PrivateKey(w io.Writer, priv *x448.PrivateKey) error {
	_, err := w.Write(priv.Secret)
	return err
}

func serializeEd25519PrivateKey(w io.Writer, priv *ed25519.PrivateKey) error {
	_, err := w.Write(priv.MarshalByteSecret())
	return err
}

func serializeEd448PrivateKey(w io.Writer, priv *ed448.PrivateKey) error {
	_, err := w.Write(priv.MarshalByteSecret())
	return err
}

// decrypt decrypts an encrypted private key using a decryption key.
func (pk *PrivateKey) decrypt(decryptionKey []byte) error {
	if pk.Dummy() {
		return errors.ErrDummyPrivateKey("dummy key found")
	}
	if !pk.Encrypted {
		return nil
	}
	block := pk.cipher.new(decryptionKey)
	var data []byte
	switch pk.s2kType {
	case S2KAEAD:
		aead := pk.aead.new(block)
		additionalData, err := pk.additionalData()
		if err != nil {
			return err
		}
		// Decrypt the encrypted key material with aead
		data, err = aead.Open(nil, pk.iv, pk.encryptedData, additionalData)
		if err != nil {
			return err
		}
	case S2KSHA1, S2KCHECKSUM:
		cfb := cipher.NewCFBDecrypter(block, pk.iv)
		data = make([]byte, len(pk.encryptedData))
		cfb.XORKeyStream(data, pk.encryptedData)
		if pk.s2kType == S2KSHA1 {
			if len(data) < sha1.Size {
				return errors.StructuralError("truncated private key data")
			}
			h := sha1.New()
			h.Write(data[:len(data)-sha1.Size])
			sum := h.Sum(nil)
			if !bytes.Equal(sum, data[len(data)-sha1.Size:]) {
				return errors.StructuralError("private key checksum failure")
			}
			data = data[:len(data)-sha1.Size]
		} else {
			if len(data) < 2 {
				return errors.StructuralError("truncated private key data")
			}
			var sum uint16
			for i := 0; i < len(data)-2; i++ {
				sum += uint16(data[i])
			}
			if data[len(data)-2] != uint8(sum>>8) ||
				data[len(data)-1] != uint8(sum) {
				return errors.StructuralError("private key checksum failure")
			}
			data = data[:len(data)-2]
		}
	default:
		return errors.InvalidArgumentError("invalid s2k type")
	}

	err := pk.parsePrivateKey(data)
	if _, ok := err.(errors.KeyInvalidError); ok {
		return errors.KeyInvalidError("invalid key parameters")
	}
	if err != nil {
		return err
	}

	// Mark key as unencrypted
	pk.s2kType = S2KNON
	pk.s2k = nil
	pk.Encrypted = false
	pk.encryptedData = nil
	return nil
}

func (pk *PrivateKey) decryptWithCache(passphrase []byte, keyCache *s2k.Cache) error {
	if pk.Dummy() {
		return errors.ErrDummyPrivateKey("dummy key found")
	}
	if !pk.Encrypted {
		return nil
	}

	key, err := keyCache.GetOrComputeDerivedKey(passphrase, pk.s2kParams, pk.cipher.KeySize())
	if err != nil {
		return err
	}
	if pk.s2kType == S2KAEAD {
		key = pk.applyHKDF(key)
	}
	return pk.decrypt(key)
}

// Decrypt decrypts an encrypted private key using a passphrase.
func (pk *PrivateKey) Decrypt(passphrase []byte) error {
	if pk.Dummy() {
		return errors.ErrDummyPrivateKey("dummy key found")
	}
	if !pk.Encrypted {
		return nil
	}

	key := make([]byte, pk.cipher.KeySize())
	pk.s2k(key, passphrase)
	if pk.s2kType == S2KAEAD {
		key = pk.applyHKDF(key)
	}
	return pk.decrypt(key)
}

// DecryptPrivateKeys decrypts all encrypted keys with the given config and passphrase.
// Avoids recomputation of similar s2k key derivations.
func DecryptPrivateKeys(keys []*PrivateKey, passphrase []byte) error {
	// Create a cache to avoid recomputation of key derviations for the same passphrase.
	s2kCache := &s2k.Cache{}
	for _, key := range keys {
		if key != nil && !key.Dummy() && key.Encrypted {
			err := key.decryptWithCache(passphrase, s2kCache)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// encrypt encrypts an unencrypted private key.
func (pk *PrivateKey) encrypt(key []byte, params *s2k.Params, s2kType S2KType, cipherFunction CipherFunction, rand io.Reader) error {
	if pk.Dummy() {
		return errors.ErrDummyPrivateKey("dummy key found")
	}
	if pk.Encrypted {
		return nil
	}
	// check if encryptionKey has the correct size
	if len(key) != cipherFunction.KeySize() {
		return errors.InvalidArgumentError("supplied encryption key has the wrong size")
	}

	if params.Mode() == s2k.Argon2S2K && s2kType != S2KAEAD {
		return errors.InvalidArgumentError("using Argon2 S2K without AEAD is not allowed")
	}
	if params.Mode() != s2k.Argon2S2K && params.Mode() != s2k.IteratedSaltedS2K &&
		params.Mode() != s2k.SaltedS2K { // only allowed for high-entropy passphrases
		return errors.InvalidArgumentError("insecure S2K mode")
	}

	priv := bytes.NewBuffer(nil)
	err := pk.serializePrivateKey(priv)
	if err != nil {
		return err
	}

	pk.cipher = cipherFunction
	pk.s2kParams = params
	pk.s2k, err = pk.s2kParams.Function()
	if err != nil {
		return err
	}

	privateKeyBytes := priv.Bytes()
	pk.s2kType = s2kType
	block := pk.cipher.new(key)
	switch s2kType {
	case S2KAEAD:
		if pk.aead == 0 {
			return errors.StructuralError("aead mode is not set on key")
		}
		aead := pk.aead.new(block)
		additionalData, err := pk.additionalData()
		if err != nil {
			return err
		}
		pk.iv = make([]byte, aead.NonceSize())
		_, err = io.ReadFull(rand, pk.iv)
		if err != nil {
			return err
		}
		// Decrypt the encrypted key material with aead
		pk.encryptedData = aead.Seal(nil, pk.iv, privateKeyBytes, additionalData)
	case S2KSHA1, S2KCHECKSUM:
		pk.iv = make([]byte, pk.cipher.blockSize())
		_, err = io.ReadFull(rand, pk.iv)
		if err != nil {
			return err
		}
		cfb := cipher.NewCFBEncrypter(block, pk.iv)
		if s2kType == S2KSHA1 {
			h := sha1.New()
			h.Write(privateKeyBytes)
			sum := h.Sum(nil)
			privateKeyBytes = append(privateKeyBytes, sum...)
		} else {
			var sum uint16
			for _, b := range privateKeyBytes {
				sum += uint16(b)
			}
			privateKeyBytes = append(privateKeyBytes, []byte{uint8(sum >> 8), uint8(sum)}...)
		}
		pk.encryptedData = make([]byte, len(privateKeyBytes))
		cfb.XORKeyStream(pk.encryptedData, privateKeyBytes)
	default:
		return errors.InvalidArgumentError("invalid s2k type for encryption")
	}

	pk.Encrypted = true
	pk.PrivateKey = nil
	return err
}

// EncryptWithConfig encrypts an unencrypted private key using the passphrase and the config.
func (pk *PrivateKey) EncryptWithConfig(passphrase []byte, config *Config) error {
	params, err := s2k.Generate(config.Random(), config.S2K())
	if err != nil {
		return err
	}
	// Derive an encryption key with the configured s2k function.
	key := make([]byte, config.Cipher().KeySize())
	s2k, err := params.Function()
	if err != nil {
		return err
	}
	s2k(key, passphrase)
	s2kType := S2KSHA1
	if config.AEAD() != nil {
		s2kType = S2KAEAD
		pk.aead = config.AEAD().Mode()
		pk.cipher = config.Cipher()
		key = pk.applyHKDF(key)
	}
	// Encrypt the private key with the derived encryption key.
	return pk.encrypt(key, params, s2kType, config.Cipher(), config.Random())
}

// EncryptPrivateKeys encrypts all unencrypted keys with the given config and passphrase.
// Only derives one key from the passphrase, which is then used to encrypt each key.
func EncryptPrivateKeys(keys []*PrivateKey, passphrase []byte, config *Config) error {
	params, err := s2k.Generate(config.Random(), config.S2K())
	if err != nil {
		return err
	}
	// Derive an encryption key with the configured s2k function.
	encryptionKey := make([]byte, config.Cipher().KeySize())
	s2k, err := params.Function()
	if err != nil {
		return err
	}
	s2k(encryptionKey, passphrase)
	for _, key := range keys {
		if key != nil && !key.Dummy() && !key.Encrypted {
			s2kType := S2KSHA1
			if config.AEAD() != nil {
				s2kType = S2KAEAD
				key.aead = config.AEAD().Mode()
				key.cipher = config.Cipher()
				derivedKey := key.applyHKDF(encryptionKey)
				err = key.encrypt(derivedKey, params, s2kType, config.Cipher(), config.Random())
			} else {
				err = key.encrypt(encryptionKey, params, s2kType, config.Cipher(), config.Random())
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Encrypt encrypts an unencrypted private key using a passphrase.
func (pk *PrivateKey) Encrypt(passphrase []byte) error {
	// Default config of private key encryption
	config := &Config{
		S2KConfig: &s2k.Config{
			S2KMode:  s2k.IteratedSaltedS2K,
			S2KCount: 65536,
			Hash:     crypto.SHA256,
		},
		DefaultCipher: CipherAES256,
	}
	return pk.EncryptWithConfig(passphrase, config)
}

func (pk *PrivateKey) serializePrivateKey(w io.Writer) (err error) {
	switch priv := pk.PrivateKey.(type) {
	case *rsa.PrivateKey:
		err = serializeRSAPrivateKey(w, priv)
	case *dsa.PrivateKey:
		err = serializeDSAPrivateKey(w, priv)
	case *elgamal.PrivateKey:
		err = serializeElGamalPrivateKey(w, priv)
	case *ecdsa.PrivateKey:
		err = serializeECDSAPrivateKey(w, priv)
	case *eddsa.PrivateKey:
		err = serializeEdDSAPrivateKey(w, priv)
	case *ecdh.PrivateKey:
		err = serializeECDHPrivateKey(w, priv)
	case *x25519.PrivateKey:
		err = serializeX25519PrivateKey(w, priv)
	case *x448.PrivateKey:
		err = serializeX448PrivateKey(w, priv)
	case *ed25519.PrivateKey:
		err = serializeEd25519PrivateKey(w, priv)
	case *ed448.PrivateKey:
		err = serializeEd448PrivateKey(w, priv)
	default:
		err = errors.InvalidArgumentError("unknown private key type")
	}
	return
}

func (pk *PrivateKey) parsePrivateKey(data []byte) (err error) {
	switch pk.PublicKey.PubKeyAlgo {
	case PubKeyAlgoRSA, PubKeyAlgoRSASignOnly, PubKeyAlgoRSAEncryptOnly:
		return pk.parseRSAPrivateKey(data)
	case PubKeyAlgoDSA:
		return pk.parseDSAPrivateKey(data)
	case PubKeyAlgoElGamal:
		return pk.parseElGamalPrivateKey(data)
	case PubKeyAlgoECDSA:
		return pk.parseECDSAPrivateKey(data)
	case PubKeyAlgoECDH:
		return pk.parseECDHPrivateKey(data)
	case PubKeyAlgoEdDSA:
		return pk.parseEdDSAPrivateKey(data)
	case PubKeyAlgoX25519:
		return pk.parseX25519PrivateKey(data)
	case PubKeyAlgoX448:
		return pk.parseX448PrivateKey(data)
	case PubKeyAlgoEd25519:
		return pk.parseEd25519PrivateKey(data)
	case PubKeyAlgoEd448:
		return pk.parseEd448PrivateKey(data)
	default:
		err = errors.StructuralError("unknown private key type")
		return
	}
}

func (pk *PrivateKey) parseRSAPrivateKey(data []byte) (err error) {
	rsaPub := pk.PublicKey.PublicKey.(*rsa.PublicKey)
	rsaPriv := new(rsa.PrivateKey)
	rsaPriv.PublicKey = *rsaPub

	buf := bytes.NewBuffer(data)
	d := new(encoding.MPI)
	if _, err := d.ReadFrom(buf); err != nil {
		return err
	}

	p := new(encoding.MPI)
	if _, err := p.ReadFrom(buf); err != nil {
		return err
	}

	q := new(encoding.MPI)
	if _, err := q.ReadFrom(buf); err != nil {
		return err
	}

	rsaPriv.D = new(big.Int).SetBytes(d.Bytes())
	rsaPriv.Primes = make([]*big.Int, 2)
	rsaPriv.Primes[0] = new(big.Int).SetBytes(p.Bytes())
	rsaPriv.Primes[1] = new(big.Int).SetBytes(q.Bytes())
	if err := rsaPriv.Validate(); err != nil {
		return errors.KeyInvalidError(err.Error())
	}
	rsaPriv.Precompute()
	pk.PrivateKey = rsaPriv

	return nil
}

func (pk *PrivateKey) parseDSAPrivateKey(data []byte) (err error) {
	dsaPub := pk.PublicKey.PublicKey.(*dsa.PublicKey)
	dsaPriv := new(dsa.PrivateKey)
	dsaPriv.PublicKey = *dsaPub

	buf := bytes.NewBuffer(data)
	x := new(encoding.MPI)
	if _, err := x.ReadFrom(buf); err != nil {
		return err
	}

	dsaPriv.X = new(big.Int).SetBytes(x.Bytes())
	if err := validateDSAParameters(dsaPriv); err != nil {
		return err
	}
	pk.PrivateKey = dsaPriv

	return nil
}

func (pk *PrivateKey) parseElGamalPrivateKey(data []byte) (err error) {
	pub := pk.PublicKey.PublicKey.(*elgamal.PublicKey)
	priv := new(elgamal.PrivateKey)
	priv.PublicKey = *pub

	buf := bytes.NewBuffer(data)
	x := new(encoding.MPI)
	if _, err := x.ReadFrom(buf); err != nil {
		return err
	}

	priv.X = new(big.Int).SetBytes(x.Bytes())
	if err := validateElGamalParameters(priv); err != nil {
		return err
	}
	pk.PrivateKey = priv

	return nil
}

func (pk *PrivateKey) parseECDSAPrivateKey(data []byte) (err error) {
	ecdsaPub := pk.PublicKey.PublicKey.(*ecdsa.PublicKey)
	ecdsaPriv := ecdsa.NewPrivateKey(*ecdsaPub)

	buf := bytes.NewBuffer(data)
	d := new(encoding.MPI)
	if _, err := d.ReadFrom(buf); err != nil {
		return err
	}

	if err := ecdsaPriv.UnmarshalIntegerSecret(d.Bytes()); err != nil {
		return err
	}
	if err := ecdsa.Validate(ecdsaPriv); err != nil {
		return err
	}
	pk.PrivateKey = ecdsaPriv

	return nil
}

func (pk *PrivateKey) parseECDHPrivateKey(data []byte) (err error) {
	ecdhPub := pk.PublicKey.PublicKey.(*ecdh.PublicKey)
	ecdhPriv := ecdh.NewPrivateKey(*ecdhPub)

	buf := bytes.NewBuffer(data)
	d := new(encoding.MPI)
	if _, err := d.ReadFrom(buf); err != nil {
		return err
	}

	if err := ecdhPriv.UnmarshalByteSecret(d.Bytes()); err != nil {
		return err
	}

	if err := ecdh.Validate(ecdhPriv); err != nil {
		return err
	}

	pk.PrivateKey = ecdhPriv

	return nil
}

func (pk *PrivateKey) parseX25519PrivateKey(data []byte) (err error) {
	publicKey := pk.PublicKey.PublicKey.(*x25519.PublicKey)
	privateKey := x25519.NewPrivateKey(*publicKey)
	privateKey.PublicKey = *publicKey

	privateKey.Secret = make([]byte, x25519.KeySize)

	if len(data) != x25519.KeySize {
		err = errors.StructuralError("wrong x25519 key size")
		return err
	}
	subtle.ConstantTimeCopy(1, privateKey.Secret, data)
	if err = x25519.Validate(privateKey); err != nil {
		return err
	}
	pk.PrivateKey = privateKey
	return nil
}

func (pk *PrivateKey) parseX448PrivateKey(data []byte) (err error) {
	publicKey := pk.PublicKey.PublicKey.(*x448.PublicKey)
	privateKey := x448.NewPrivateKey(*publicKey)
	privateKey.PublicKey = *publicKey

	privateKey.Secret = make([]byte, x448.KeySize)

	if len(data) != x448.KeySize {
		err = errors.StructuralError("wrong x448 key size")
		return err
	}
	subtle.ConstantTimeCopy(1, privateKey.Secret, data)
	if err = x448.Validate(privateKey); err != nil {
		return err
	}
	pk.PrivateKey = privateKey
	return nil
}

func (pk *PrivateKey) parseEd25519PrivateKey(data []byte) (err error) {
	publicKey := pk.PublicKey.PublicKey.(*ed25519.PublicKey)
	privateKey := ed25519.NewPrivateKey(*publicKey)
	privateKey.PublicKey = *publicKey

	if len(data) != ed25519.SeedSize {
		err = errors.StructuralError("wrong ed25519 key size")
		return err
	}
	err = privateKey.UnmarshalByteSecret(data)
	if err != nil {
		return err
	}
	err = ed25519.Validate(privateKey)
	if err != nil {
		return err
	}
	pk.PrivateKey = privateKey
	return nil
}

func (pk *PrivateKey) parseEd448PrivateKey(data []byte) (err error) {
	publicKey := pk.PublicKey.PublicKey.(*ed448.PublicKey)
	privateKey := ed448.NewPrivateKey(*publicKey)
	privateKey.PublicKey = *publicKey

	if len(data) != ed448.SeedSize {
		err = errors.StructuralError("wrong ed448 key size")
		return err
	}
	err = privateKey.UnmarshalByteSecret(data)
	if err != nil {
		return err
	}
	err = ed448.Validate(privateKey)
	if err != nil {
		return err
	}
	pk.PrivateKey = privateKey
	return nil
}

func (pk *PrivateKey) parseEdDSAPrivateKey(data []byte) (err error) {
	eddsaPub := pk.PublicKey.PublicKey.(*eddsa.PublicKey)
	eddsaPriv := eddsa.NewPrivateKey(*eddsaPub)
	eddsaPriv.PublicKey = *eddsaPub

	buf := bytes.NewBuffer(data)
	d := new(encoding.MPI)
	if _, err := d.ReadFrom(buf); err != nil {
		return err
	}

	if err = eddsaPriv.UnmarshalByteSecret(d.Bytes()); err != nil {
		return err
	}

	if err := eddsa.Validate(eddsaPriv); err != nil {
		return err
	}

	pk.PrivateKey = eddsaPriv

	return nil
}

func (pk *PrivateKey) additionalData() ([]byte, error) {
	additionalData := bytes.NewBuffer(nil)
	// Write additional data prefix based on packet type
	var packetByte byte
	if pk.PublicKey.IsSubkey {
		packetByte = 0xc7
	} else {
		packetByte = 0xc5
	}
	// Write public key to additional data
	_, err := additionalData.Write([]byte{packetByte})
	if err != nil {
		return nil, err
	}
	err = pk.PublicKey.serializeWithoutHeaders(additionalData)
	if err != nil {
		return nil, err
	}
	return additionalData.Bytes(), nil
}

func (pk *PrivateKey) applyHKDF(inputKey []byte) []byte {
	var packetByte byte
	if pk.PublicKey.IsSubkey {
		packetByte = 0xc7
	} else {
		packetByte = 0xc5
	}
	associatedData := []byte{packetByte, byte(pk.Version), byte(pk.cipher), byte(pk.aead)}
	hkdfReader := hkdf.New(sha256.New, inputKey, []byte{}, associatedData)
	encryptionKey := make([]byte, pk.cipher.KeySize())
	_, _ = readFull(hkdfReader, encryptionKey)
	return encryptionKey
}

func validateDSAParameters(priv *dsa.PrivateKey) error {
	p := priv.P // group prime
	q := priv.Q // subgroup order
	g := priv.G // g has order q mod p
	x := priv.X // secret
	y := priv.Y // y == g**x mod p
	one := big.NewInt(1)
	// expect g, y >= 2 and g < p
	if g.Cmp(one) <= 0 || y.Cmp(one) <= 0 || g.Cmp(p) > 0 {
		return errors.KeyInvalidError("dsa: invalid group")
	}
	// expect p > q
	if p.Cmp(q) <= 0 {
		return errors.KeyInvalidError("dsa: invalid group prime")
	}
	// q should be large enough and divide p-1
	pSub1 := new(big.Int).Sub(p, one)
	if q.BitLen() < 150 || new(big.Int).Mod(pSub1, q).Cmp(big.NewInt(0)) != 0 {
		return errors.KeyInvalidError("dsa: invalid order")
	}
	// confirm that g has order q mod p
	if !q.ProbablyPrime(32) || new(big.Int).Exp(g, q, p).Cmp(one) != 0 {
		return errors.KeyInvalidError("dsa: invalid order")
	}
	// check y
	if new(big.Int).Exp(g, x, p).Cmp(y) != 0 {
		return errors.KeyInvalidError("dsa: mismatching values")
	}

	return nil
}

func validateElGamalParameters(priv *elgamal.PrivateKey) error {
	p := priv.P // group prime
	g := priv.G // g has order p-1 mod p
	x := priv.X // secret
	y := priv.Y // y == g**x mod p
	one := big.NewInt(1)
	// Expect g, y >= 2 and g < p
	if g.Cmp(one) <= 0 || y.Cmp(one) <= 0 || g.Cmp(p) > 0 {
		return errors.KeyInvalidError("elgamal: invalid group")
	}
	if p.BitLen() < 1024 {
		return errors.KeyInvalidError("elgamal: group order too small")
	}
	pSub1 := new(big.Int).Sub(p, one)
	if new(big.Int).Exp(g, pSub1, p).Cmp(one) != 0 {
		return errors.KeyInvalidError("elgamal: invalid group")
	}
	// Since p-1 is not prime, g might have a smaller order that divides p-1.
	// We cannot confirm the exact order of g, but we make sure it is not too small.
	gExpI := new(big.Int).Set(g)
	i := 1
	threshold := 2 << 17 // we want order > threshold
	for i < threshold {
		i++ // we check every order to make sure key validation is not easily bypassed by guessing y'
		gExpI.Mod(new(big.Int).Mul(gExpI, g), p)
		if gExpI.Cmp(one) == 0 {
			return errors.KeyInvalidError("elgamal: order too small")
		}
	}
	// Check y
	if new(big.Int).Exp(g, x, p).Cmp(y) != 0 {
		return errors.KeyInvalidError("elgamal: mismatching values")
	}

	return nil
}
