// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encoding

import (
	"io"
	"math/big"
	"math/bits"
)

// An MPI is used to store the contents of a big integer, along with the bit
// length that was specified in the original input. This allows the MPI to be
// reserialized exactly.
type MPI struct {
	bytes     []byte
	bitLength uint16
}

// NewMPI returns a MPI initialized with bytes.
func NewMPI(bytes []byte) *MPI {
	for len(bytes) != 0 && bytes[0] == 0 {
		bytes = bytes[1:]
	}
	if len(bytes) == 0 {
		bitLength := uint16(0)
		return &MPI{bytes, bitLength}
	}
	bitLength := 8*uint16(len(bytes)-1) + uint16(bits.Len8(bytes[0]))
	return &MPI{bytes, bitLength}
}

// Bytes returns the decoded data.
func (m *MPI) Bytes() []byte {
	return m.bytes
}

// BitLength is the size in bits of the decoded data.
func (m *MPI) BitLength() uint16 {
	return m.bitLength
}

// EncodedBytes returns the encoded data.
func (m *MPI) EncodedBytes() []byte {
	return append([]byte{byte(m.bitLength >> 8), byte(m.bitLength)}, m.bytes...)
}

// EncodedLength is the size in bytes of the encoded data.
func (m *MPI) EncodedLength() uint16 {
	return uint16(2 + len(m.bytes))
}

// ReadFrom reads into m the next MPI from r.
func (m *MPI) ReadFrom(r io.Reader) (int64, error) {
	var buf [2]byte
	n, err := io.ReadFull(r, buf[0:])
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return int64(n), err
	}

	m.bitLength = uint16(buf[0])<<8 | uint16(buf[1])
	m.bytes = make([]byte, (int(m.bitLength)+7)/8)

	nn, err := io.ReadFull(r, m.bytes)
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}

	// remove leading zero bytes from malformed GnuPG encoded MPIs:
	// https://bugs.gnupg.org/gnupg/issue1853
	// for _, b := range m.bytes {
	// 	if b != 0 {
	// 		break
	// 	}
	// 	m.bytes = m.bytes[1:]
	// 	m.bitLength -= 8
	// }

	return int64(n) + int64(nn), err
}

// SetBig initializes m with the bits from n.
func (m *MPI) SetBig(n *big.Int) *MPI {
	m.bytes = n.Bytes()
	m.bitLength = uint16(n.BitLen())
	return m
}
