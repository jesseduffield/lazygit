// Copyright 2023 Proton AG. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"io"
	"strconv"

	"github.com/ProtonMail/go-crypto/openpgp/errors"
	"golang.org/x/crypto/hkdf"
)

// parseAead parses a V2 SEIPD packet (AEAD) as specified in
// https://www.ietf.org/archive/id/draft-ietf-openpgp-crypto-refresh-07.html#section-5.13.2
func (se *SymmetricallyEncrypted) parseAead(r io.Reader) error {
	headerData := make([]byte, 3)
	if n, err := io.ReadFull(r, headerData); n < 3 {
		return errors.StructuralError("could not read aead header: " + err.Error())
	}

	// Cipher
	se.Cipher = CipherFunction(headerData[0])
	// cipherFunc must have block size 16 to use AEAD
	if se.Cipher.blockSize() != 16 {
		return errors.UnsupportedError("invalid aead cipher: " + strconv.Itoa(int(se.Cipher)))
	}

	// Mode
	se.Mode = AEADMode(headerData[1])
	if se.Mode.TagLength() == 0 {
		return errors.UnsupportedError("unknown aead mode: " + strconv.Itoa(int(se.Mode)))
	}

	// Chunk size
	se.ChunkSizeByte = headerData[2]
	if se.ChunkSizeByte > 16 {
		return errors.UnsupportedError("invalid aead chunk size byte: " + strconv.Itoa(int(se.ChunkSizeByte)))
	}

	// Salt
	if n, err := io.ReadFull(r, se.Salt[:]); n < aeadSaltSize {
		return errors.StructuralError("could not read aead salt: " + err.Error())
	}

	return nil
}

// associatedData for chunks: tag, version, cipher, mode, chunk size byte
func (se *SymmetricallyEncrypted) associatedData() []byte {
	return []byte{
		0xD2,
		symmetricallyEncryptedVersionAead,
		byte(se.Cipher),
		byte(se.Mode),
		se.ChunkSizeByte,
	}
}

// decryptAead decrypts a V2 SEIPD packet (AEAD) as specified in
// https://www.ietf.org/archive/id/draft-ietf-openpgp-crypto-refresh-07.html#section-5.13.2
func (se *SymmetricallyEncrypted) decryptAead(inputKey []byte) (io.ReadCloser, error) {
	if se.Cipher.KeySize() != len(inputKey) {
		return nil, errors.StructuralError(fmt.Sprintf("invalid session key length for cipher: got %d bytes, but expected %d bytes", len(inputKey), se.Cipher.KeySize()))
	}

	aead, nonce := getSymmetricallyEncryptedAeadInstance(se.Cipher, se.Mode, inputKey, se.Salt[:], se.associatedData())
	// Carry the first tagLen bytes
	chunkSize := decodeAEADChunkSize(se.ChunkSizeByte)
	tagLen := se.Mode.TagLength()
	chunkBytes := make([]byte, chunkSize+tagLen*2)
	peekedBytes := chunkBytes[chunkSize+tagLen:]
	n, err := io.ReadFull(se.Contents, peekedBytes)
	if n < tagLen || (err != nil && err != io.EOF) {
		return nil, errors.StructuralError("not enough data to decrypt:" + err.Error())
	}

	return &aeadDecrypter{
		aeadCrypter: aeadCrypter{
			aead:           aead,
			chunkSize:      decodeAEADChunkSize(se.ChunkSizeByte),
			nonce:          nonce,
			associatedData: se.associatedData(),
			chunkIndex:     nonce[len(nonce)-8:],
			packetTag:      packetTypeSymmetricallyEncryptedIntegrityProtected,
		},
		reader:      se.Contents,
		chunkBytes:  chunkBytes,
		peekedBytes: peekedBytes,
	}, nil
}

// serializeSymmetricallyEncryptedAead encrypts to a writer a V2 SEIPD packet (AEAD) as specified in
// https://www.ietf.org/archive/id/draft-ietf-openpgp-crypto-refresh-07.html#section-5.13.2
func serializeSymmetricallyEncryptedAead(ciphertext io.WriteCloser, cipherSuite CipherSuite, chunkSizeByte byte, rand io.Reader, inputKey []byte) (Contents io.WriteCloser, err error) {
	// cipherFunc must have block size 16 to use AEAD
	if cipherSuite.Cipher.blockSize() != 16 {
		return nil, errors.InvalidArgumentError("invalid aead cipher function")
	}

	if cipherSuite.Cipher.KeySize() != len(inputKey) {
		return nil, errors.InvalidArgumentError("error in aead serialization: bad key length")
	}

	// Data for en/decryption: tag, version, cipher, aead mode, chunk size
	prefix := []byte{
		0xD2,
		symmetricallyEncryptedVersionAead,
		byte(cipherSuite.Cipher),
		byte(cipherSuite.Mode),
		chunkSizeByte,
	}

	// Write header (that correspond to prefix except first byte)
	n, err := ciphertext.Write(prefix[1:])
	if err != nil || n < 4 {
		return nil, err
	}

	// Random salt
	salt := make([]byte, aeadSaltSize)
	if _, err := io.ReadFull(rand, salt); err != nil {
		return nil, err
	}

	if _, err := ciphertext.Write(salt); err != nil {
		return nil, err
	}

	aead, nonce := getSymmetricallyEncryptedAeadInstance(cipherSuite.Cipher, cipherSuite.Mode, inputKey, salt, prefix)

	chunkSize := decodeAEADChunkSize(chunkSizeByte)
	tagLen := aead.Overhead()
	chunkBytes := make([]byte, chunkSize+tagLen)
	return &aeadEncrypter{
		aeadCrypter: aeadCrypter{
			aead:           aead,
			chunkSize:      chunkSize,
			associatedData: prefix,
			nonce:          nonce,
			chunkIndex:     nonce[len(nonce)-8:],
			packetTag:      packetTypeSymmetricallyEncryptedIntegrityProtected,
		},
		writer:     ciphertext,
		chunkBytes: chunkBytes,
	}, nil
}

func getSymmetricallyEncryptedAeadInstance(c CipherFunction, mode AEADMode, inputKey, salt, associatedData []byte) (aead cipher.AEAD, nonce []byte) {
	hkdfReader := hkdf.New(sha256.New, inputKey, salt, associatedData)

	encryptionKey := make([]byte, c.KeySize())
	_, _ = readFull(hkdfReader, encryptionKey)

	nonce = make([]byte, mode.IvLength())

	// Last 64 bits of nonce are the counter
	_, _ = readFull(hkdfReader, nonce[:len(nonce)-8])

	blockCipher := c.new(encryptionKey)
	aead = mode.new(blockCipher)

	return
}
