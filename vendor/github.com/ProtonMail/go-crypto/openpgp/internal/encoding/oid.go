// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encoding

import (
	"io"

	"github.com/ProtonMail/go-crypto/openpgp/errors"
)

// OID is used to store a variable-length field with a one-octet size
// prefix. See https://tools.ietf.org/html/rfc6637#section-9.
type OID struct {
	bytes []byte
}

const (
	// maxOID is the maximum number of bytes in a OID.
	maxOID = 254
	// reservedOIDLength1 and reservedOIDLength2 are OID lengths that the RFC
	// specifies are reserved.
	reservedOIDLength1 = 0
	reservedOIDLength2 = 0xff
)

// NewOID returns a OID initialized with bytes.
func NewOID(bytes []byte) *OID {
	switch len(bytes) {
	case reservedOIDLength1, reservedOIDLength2:
		panic("encoding: NewOID argument length is reserved")
	default:
		if len(bytes) > maxOID {
			panic("encoding: NewOID argument too large")
		}
	}

	return &OID{
		bytes: bytes,
	}
}

// Bytes returns the decoded data.
func (o *OID) Bytes() []byte {
	return o.bytes
}

// BitLength is the size in bits of the decoded data.
func (o *OID) BitLength() uint16 {
	return uint16(len(o.bytes) * 8)
}

// EncodedBytes returns the encoded data.
func (o *OID) EncodedBytes() []byte {
	return append([]byte{byte(len(o.bytes))}, o.bytes...)
}

// EncodedLength is the size in bytes of the encoded data.
func (o *OID) EncodedLength() uint16 {
	return uint16(1 + len(o.bytes))
}

// ReadFrom reads into b the next OID from r.
func (o *OID) ReadFrom(r io.Reader) (int64, error) {
	var buf [1]byte
	n, err := io.ReadFull(r, buf[:])
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return int64(n), err
	}

	switch buf[0] {
	case reservedOIDLength1, reservedOIDLength2:
		return int64(n), errors.UnsupportedError("reserved for future extensions")
	}

	o.bytes = make([]byte, buf[0])

	nn, err := io.ReadFull(r, o.bytes)
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}

	return int64(n) + int64(nn), err
}
