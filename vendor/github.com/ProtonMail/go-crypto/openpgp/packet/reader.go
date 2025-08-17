// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"io"

	"github.com/ProtonMail/go-crypto/openpgp/errors"
)

type PacketReader interface {
	Next() (p Packet, err error)
	Push(reader io.Reader) (err error)
	Unread(p Packet)
}

// Reader reads packets from an io.Reader and allows packets to be 'unread' so
// that they result from the next call to Next.
type Reader struct {
	q       []Packet
	readers []io.Reader
}

// New io.Readers are pushed when a compressed or encrypted packet is processed
// and recursively treated as a new source of packets. However, a carefully
// crafted packet can trigger an infinite recursive sequence of packets. See
// http://mumble.net/~campbell/misc/pgp-quine
// https://web.nvd.nist.gov/view/vuln/detail?vulnId=CVE-2013-4402
// This constant limits the number of recursive packets that may be pushed.
const maxReaders = 32

// Next returns the most recently unread Packet, or reads another packet from
// the top-most io.Reader. Unknown/unsupported/Marker packet types are skipped.
func (r *Reader) Next() (p Packet, err error) {
	for {
		p, err := r.read()
		if err == io.EOF {
			break
		} else if err != nil {
			if _, ok := err.(errors.UnknownPacketTypeError); ok {
				continue
			}
			if _, ok := err.(errors.UnsupportedError); ok {
				switch p.(type) {
				case *SymmetricallyEncrypted, *AEADEncrypted, *Compressed, *LiteralData:
					return nil, err
				}
				continue
			}
			return nil, err
		} else {
			//A marker packet MUST be ignored when received
			switch p.(type) {
			case *Marker:
				continue
			}
			return p, nil
		}
	}
	return nil, io.EOF
}

// Next returns the most recently unread Packet, or reads another packet from
// the top-most io.Reader. Unknown/Marker packet types are skipped while unsupported
// packets are returned as UnsupportedPacket type.
func (r *Reader) NextWithUnsupported() (p Packet, err error) {
	for {
		p, err = r.read()
		if err == io.EOF {
			break
		} else if err != nil {
			if _, ok := err.(errors.UnknownPacketTypeError); ok {
				continue
			}
			if casteErr, ok := err.(errors.UnsupportedError); ok {
				return &UnsupportedPacket{
					IncompletePacket: p,
					Error:            casteErr,
				}, nil
			}
			return
		} else {
			//A marker packet MUST be ignored when received
			switch p.(type) {
			case *Marker:
				continue
			}
			return
		}
	}
	return nil, io.EOF
}

func (r *Reader) read() (p Packet, err error) {
	if len(r.q) > 0 {
		p = r.q[len(r.q)-1]
		r.q = r.q[:len(r.q)-1]
		return
	}
	for len(r.readers) > 0 {
		p, err = Read(r.readers[len(r.readers)-1])
		if err == io.EOF {
			r.readers = r.readers[:len(r.readers)-1]
			continue
		}
		return p, err
	}
	return nil, io.EOF
}

// Push causes the Reader to start reading from a new io.Reader. When an EOF
// error is seen from the new io.Reader, it is popped and the Reader continues
// to read from the next most recent io.Reader. Push returns a StructuralError
// if pushing the reader would exceed the maximum recursion level, otherwise it
// returns nil.
func (r *Reader) Push(reader io.Reader) (err error) {
	if len(r.readers) >= maxReaders {
		return errors.StructuralError("too many layers of packets")
	}
	r.readers = append(r.readers, reader)
	return nil
}

// Unread causes the given Packet to be returned from the next call to Next.
func (r *Reader) Unread(p Packet) {
	r.q = append(r.q, p)
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		q:       nil,
		readers: []io.Reader{r},
	}
}

// CheckReader is similar to Reader but additionally
// uses the pushdown automata to verify the read packet sequence.
type CheckReader struct {
	Reader
	verifier  *SequenceVerifier
	fullyRead bool
}

// Next returns the most recently unread Packet, or reads another packet from
// the top-most io.Reader. Unknown packet types are skipped.
// If the read packet sequence does not conform to the packet composition
// rules in rfc4880, it returns an error.
func (r *CheckReader) Next() (p Packet, err error) {
	if r.fullyRead {
		return nil, io.EOF
	}
	if len(r.q) > 0 {
		p = r.q[len(r.q)-1]
		r.q = r.q[:len(r.q)-1]
		return
	}
	var errMsg error
	for len(r.readers) > 0 {
		p, errMsg, err = ReadWithCheck(r.readers[len(r.readers)-1], r.verifier)
		if errMsg != nil {
			err = errMsg
			return
		}
		if err == nil {
			return
		}
		if err == io.EOF {
			r.readers = r.readers[:len(r.readers)-1]
			continue
		}
		//A marker packet MUST be ignored when received
		switch p.(type) {
		case *Marker:
			continue
		}
		if _, ok := err.(errors.UnknownPacketTypeError); ok {
			continue
		}
		if _, ok := err.(errors.UnsupportedError); ok {
			switch p.(type) {
			case *SymmetricallyEncrypted, *AEADEncrypted, *Compressed, *LiteralData:
				return nil, err
			}
			continue
		}
		return nil, err
	}
	if errMsg = r.verifier.Next(EOSSymbol); errMsg != nil {
		return nil, errMsg
	}
	if errMsg = r.verifier.AssertValid(); errMsg != nil {
		return nil, errMsg
	}
	r.fullyRead = true
	return nil, io.EOF
}

func NewCheckReader(r io.Reader) *CheckReader {
	return &CheckReader{
		Reader: Reader{
			q:       nil,
			readers: []io.Reader{r},
		},
		verifier:  NewSequenceVerifier(),
		fullyRead: false,
	}
}
