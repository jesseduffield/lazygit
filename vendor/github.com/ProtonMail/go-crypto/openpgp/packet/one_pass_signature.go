// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"crypto"
	"encoding/binary"
	"io"
	"strconv"

	"github.com/ProtonMail/go-crypto/openpgp/errors"
	"github.com/ProtonMail/go-crypto/openpgp/internal/algorithm"
)

// OnePassSignature represents a one-pass signature packet. See RFC 4880,
// section 5.4.
type OnePassSignature struct {
	Version        int
	SigType        SignatureType
	Hash           crypto.Hash
	PubKeyAlgo     PublicKeyAlgorithm
	KeyId          uint64
	IsLast         bool
	Salt           []byte // v6 only
	KeyFingerprint []byte // v6 only
}

func (ops *OnePassSignature) parse(r io.Reader) (err error) {
	var buf [8]byte
	// Read: version | signature type | hash algorithm | public-key algorithm
	_, err = readFull(r, buf[:4])
	if err != nil {
		return
	}
	if buf[0] != 3 && buf[0] != 6 {
		return errors.UnsupportedError("one-pass-signature packet version " + strconv.Itoa(int(buf[0])))
	}
	ops.Version = int(buf[0])

	var ok bool
	ops.Hash, ok = algorithm.HashIdToHashWithSha1(buf[2])
	if !ok {
		return errors.UnsupportedError("hash function: " + strconv.Itoa(int(buf[2])))
	}

	ops.SigType = SignatureType(buf[1])
	ops.PubKeyAlgo = PublicKeyAlgorithm(buf[3])

	if ops.Version == 6 {
		// Only for v6, a variable-length field containing the salt
		_, err = readFull(r, buf[:1])
		if err != nil {
			return
		}
		saltLength := int(buf[0])
		var expectedSaltLength int
		expectedSaltLength, err = SaltLengthForHash(ops.Hash)
		if err != nil {
			return
		}
		if saltLength != expectedSaltLength {
			err = errors.StructuralError("unexpected salt size for the given hash algorithm")
			return
		}
		salt := make([]byte, expectedSaltLength)
		_, err = readFull(r, salt)
		if err != nil {
			return
		}
		ops.Salt = salt

		// Only for v6 packets, 32 octets of the fingerprint of the signing key.
		fingerprint := make([]byte, 32)
		_, err = readFull(r, fingerprint)
		if err != nil {
			return
		}
		ops.KeyFingerprint = fingerprint
		ops.KeyId = binary.BigEndian.Uint64(ops.KeyFingerprint[:8])
	} else {
		_, err = readFull(r, buf[:8])
		if err != nil {
			return
		}
		ops.KeyId = binary.BigEndian.Uint64(buf[:8])
	}

	_, err = readFull(r, buf[:1])
	if err != nil {
		return
	}
	ops.IsLast = buf[0] != 0
	return
}

// Serialize marshals the given OnePassSignature to w.
func (ops *OnePassSignature) Serialize(w io.Writer) error {
	//v3 length 1+1+1+1+8+1 =
	packetLength := 13
	if ops.Version == 6 {
		// v6 length 1+1+1+1+1+len(salt)+32+1 =
		packetLength = 38 + len(ops.Salt)
	}

	if err := serializeHeader(w, packetTypeOnePassSignature, packetLength); err != nil {
		return err
	}

	var buf [8]byte
	buf[0] = byte(ops.Version)
	buf[1] = uint8(ops.SigType)
	var ok bool
	buf[2], ok = algorithm.HashToHashIdWithSha1(ops.Hash)
	if !ok {
		return errors.UnsupportedError("hash type: " + strconv.Itoa(int(ops.Hash)))
	}
	buf[3] = uint8(ops.PubKeyAlgo)

	_, err := w.Write(buf[:4])
	if err != nil {
		return err
	}

	if ops.Version == 6 {
		// write salt for v6 signatures
		_, err := w.Write([]byte{uint8(len(ops.Salt))})
		if err != nil {
			return err
		}
		_, err = w.Write(ops.Salt)
		if err != nil {
			return err
		}

		// write fingerprint v6 signatures
		_, err = w.Write(ops.KeyFingerprint)
		if err != nil {
			return err
		}
	} else {
		binary.BigEndian.PutUint64(buf[:8], ops.KeyId)
		_, err := w.Write(buf[:8])
		if err != nil {
			return err
		}
	}

	isLast := []byte{byte(0)}
	if ops.IsLast {
		isLast[0] = 1
	}

	_, err = w.Write(isLast)
	return err
}
