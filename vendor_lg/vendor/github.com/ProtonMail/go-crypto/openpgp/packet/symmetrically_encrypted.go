// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"io"

	"github.com/ProtonMail/go-crypto/openpgp/errors"
)

const aeadSaltSize = 32

// SymmetricallyEncrypted represents a symmetrically encrypted byte string. The
// encrypted Contents will consist of more OpenPGP packets. See RFC 4880,
// sections 5.7 and 5.13.
type SymmetricallyEncrypted struct {
	Version            int
	Contents           io.Reader // contains tag for version 2
	IntegrityProtected bool      // If true it is type 18 (with MDC or AEAD). False is packet type 9

	// Specific to version 1
	prefix []byte

	// Specific to version 2
	Cipher        CipherFunction
	Mode          AEADMode
	ChunkSizeByte byte
	Salt          [aeadSaltSize]byte
}

const (
	symmetricallyEncryptedVersionMdc  = 1
	symmetricallyEncryptedVersionAead = 2
)

func (se *SymmetricallyEncrypted) parse(r io.Reader) error {
	if se.IntegrityProtected {
		// See RFC 4880, section 5.13.
		var buf [1]byte
		_, err := readFull(r, buf[:])
		if err != nil {
			return err
		}

		switch buf[0] {
		case symmetricallyEncryptedVersionMdc:
			se.Version = symmetricallyEncryptedVersionMdc
		case symmetricallyEncryptedVersionAead:
			se.Version = symmetricallyEncryptedVersionAead
			if err := se.parseAead(r); err != nil {
				return err
			}
		default:
			return errors.UnsupportedError("unknown SymmetricallyEncrypted version")
		}
	}
	se.Contents = r
	return nil
}

// Decrypt returns a ReadCloser, from which the decrypted Contents of the
// packet can be read. An incorrect key will only be detected after trying
// to decrypt the entire data.
func (se *SymmetricallyEncrypted) Decrypt(c CipherFunction, key []byte) (io.ReadCloser, error) {
	if se.Version == symmetricallyEncryptedVersionAead {
		return se.decryptAead(key)
	}

	return se.decryptMdc(c, key)
}

// SerializeSymmetricallyEncrypted serializes a symmetrically encrypted packet
// to w and returns a WriteCloser to which the to-be-encrypted packets can be
// written.
// If aeadSupported is set to true, SEIPDv2 is used with the indicated CipherSuite.
// Otherwise, SEIPDv1 is used with the indicated CipherFunction.
// Note: aeadSupported MUST match the value passed to SerializeEncryptedKeyAEAD
// and/or SerializeSymmetricKeyEncryptedAEADReuseKey.
// If config is nil, sensible defaults will be used.
func SerializeSymmetricallyEncrypted(w io.Writer, c CipherFunction, aeadSupported bool, cipherSuite CipherSuite, key []byte, config *Config) (Contents io.WriteCloser, err error) {
	writeCloser := noOpCloser{w}
	ciphertext, err := serializeStreamHeader(writeCloser, packetTypeSymmetricallyEncryptedIntegrityProtected)
	if err != nil {
		return
	}

	if aeadSupported {
		return serializeSymmetricallyEncryptedAead(ciphertext, cipherSuite, config.AEADConfig.ChunkSizeByte(), config.Random(), key)
	}

	return serializeSymmetricallyEncryptedMdc(ciphertext, c, key, config)
}
