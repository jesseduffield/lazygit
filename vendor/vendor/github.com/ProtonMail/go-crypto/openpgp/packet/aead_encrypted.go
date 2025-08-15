// Copyright (C) 2019 ProtonTech AG

package packet

import (
	"io"

	"github.com/ProtonMail/go-crypto/openpgp/errors"
	"github.com/ProtonMail/go-crypto/openpgp/internal/algorithm"
)

// AEADEncrypted represents an AEAD Encrypted Packet.
// See https://www.ietf.org/archive/id/draft-koch-openpgp-2015-rfc4880bis-00.html#name-aead-encrypted-data-packet-t
type AEADEncrypted struct {
	cipher        CipherFunction
	mode          AEADMode
	chunkSizeByte byte
	Contents      io.Reader // Encrypted chunks and tags
	initialNonce  []byte    // Referred to as IV in RFC4880-bis
}

// Only currently defined version
const aeadEncryptedVersion = 1

func (ae *AEADEncrypted) parse(buf io.Reader) error {
	headerData := make([]byte, 4)
	if n, err := io.ReadFull(buf, headerData); n < 4 {
		return errors.AEADError("could not read aead header:" + err.Error())
	}
	// Read initial nonce
	mode := AEADMode(headerData[2])
	nonceLen := mode.IvLength()

	// This packet supports only EAX and OCB
	// https://www.ietf.org/archive/id/draft-koch-openpgp-2015-rfc4880bis-00.html#name-aead-encrypted-data-packet-t
	if nonceLen == 0 || mode > AEADModeOCB {
		return errors.AEADError("unknown mode")
	}

	initialNonce := make([]byte, nonceLen)
	if n, err := io.ReadFull(buf, initialNonce); n < nonceLen {
		return errors.AEADError("could not read aead nonce:" + err.Error())
	}
	ae.Contents = buf
	ae.initialNonce = initialNonce
	c := headerData[1]
	if _, ok := algorithm.CipherById[c]; !ok {
		return errors.UnsupportedError("unknown cipher: " + string(c))
	}
	ae.cipher = CipherFunction(c)
	ae.mode = mode
	ae.chunkSizeByte = headerData[3]
	return nil
}

// Decrypt returns a io.ReadCloser from which decrypted bytes can be read, or
// an error.
func (ae *AEADEncrypted) Decrypt(ciph CipherFunction, key []byte) (io.ReadCloser, error) {
	return ae.decrypt(key)
}

// decrypt prepares an aeadCrypter and returns a ReadCloser from which
// decrypted bytes can be read (see aeadDecrypter.Read()).
func (ae *AEADEncrypted) decrypt(key []byte) (io.ReadCloser, error) {
	blockCipher := ae.cipher.new(key)
	aead := ae.mode.new(blockCipher)
	// Carry the first tagLen bytes
	chunkSize := decodeAEADChunkSize(ae.chunkSizeByte)
	tagLen := ae.mode.TagLength()
	chunkBytes := make([]byte, chunkSize+tagLen*2)
	peekedBytes := chunkBytes[chunkSize+tagLen:]
	n, err := io.ReadFull(ae.Contents, peekedBytes)
	if n < tagLen || (err != nil && err != io.EOF) {
		return nil, errors.AEADError("Not enough data to decrypt:" + err.Error())
	}

	return &aeadDecrypter{
		aeadCrypter: aeadCrypter{
			aead:           aead,
			chunkSize:      chunkSize,
			nonce:          ae.initialNonce,
			associatedData: ae.associatedData(),
			chunkIndex:     make([]byte, 8),
			packetTag:      packetTypeAEADEncrypted,
		},
		reader:      ae.Contents,
		chunkBytes:  chunkBytes,
		peekedBytes: peekedBytes,
	}, nil
}

// associatedData for chunks: tag, version, cipher, mode, chunk size byte
func (ae *AEADEncrypted) associatedData() []byte {
	return []byte{
		0xD4,
		aeadEncryptedVersion,
		byte(ae.cipher),
		byte(ae.mode),
		ae.chunkSizeByte}
}
