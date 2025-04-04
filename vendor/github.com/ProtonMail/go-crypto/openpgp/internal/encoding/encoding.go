// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package encoding implements openpgp packet field encodings as specified in
// RFC 4880 and 6637.
package encoding

import "io"

// Field is an encoded field of an openpgp packet.
type Field interface {
	// Bytes returns the decoded data.
	Bytes() []byte

	// BitLength is the size in bits of the decoded data.
	BitLength() uint16

	// EncodedBytes returns the encoded data.
	EncodedBytes() []byte

	// EncodedLength is the size in bytes of the encoded data.
	EncodedLength() uint16

	// ReadFrom reads the next Field from r.
	ReadFrom(r io.Reader) (int64, error)
}
