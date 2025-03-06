// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package armor

import (
	"encoding/base64"
	"io"
	"sort"
)

var armorHeaderSep = []byte(": ")
var blockEnd = []byte("\n=")
var newline = []byte("\n")
var armorEndOfLineOut = []byte("-----\n")

const crc24Init = 0xb704ce
const crc24Poly = 0x1864cfb

// crc24 calculates the OpenPGP checksum as specified in RFC 4880, section 6.1
func crc24(crc uint32, d []byte) uint32 {
	for _, b := range d {
		crc ^= uint32(b) << 16
		for i := 0; i < 8; i++ {
			crc <<= 1
			if crc&0x1000000 != 0 {
				crc ^= crc24Poly
			}
		}
	}
	return crc
}

// writeSlices writes its arguments to the given Writer.
func writeSlices(out io.Writer, slices ...[]byte) (err error) {
	for _, s := range slices {
		_, err = out.Write(s)
		if err != nil {
			return err
		}
	}
	return
}

// lineBreaker breaks data across several lines, all of the same byte length
// (except possibly the last). Lines are broken with a single '\n'.
type lineBreaker struct {
	lineLength  int
	line        []byte
	used        int
	out         io.Writer
	haveWritten bool
}

func newLineBreaker(out io.Writer, lineLength int) *lineBreaker {
	return &lineBreaker{
		lineLength: lineLength,
		line:       make([]byte, lineLength),
		used:       0,
		out:        out,
	}
}

func (l *lineBreaker) Write(b []byte) (n int, err error) {
	n = len(b)

	if n == 0 {
		return
	}

	if l.used == 0 && l.haveWritten {
		_, err = l.out.Write([]byte{'\n'})
		if err != nil {
			return
		}
	}

	if l.used+len(b) < l.lineLength {
		l.used += copy(l.line[l.used:], b)
		return
	}

	l.haveWritten = true
	_, err = l.out.Write(l.line[0:l.used])
	if err != nil {
		return
	}
	excess := l.lineLength - l.used
	l.used = 0

	_, err = l.out.Write(b[0:excess])
	if err != nil {
		return
	}

	_, err = l.Write(b[excess:])
	return
}

func (l *lineBreaker) Close() (err error) {
	if l.used > 0 {
		_, err = l.out.Write(l.line[0:l.used])
		if err != nil {
			return
		}
	}

	return
}

// encoding keeps track of a running CRC24 over the data which has been written
// to it and outputs a OpenPGP checksum when closed, followed by an armor
// trailer.
//
// It's built into a stack of io.Writers:
//
//	encoding -> base64 encoder -> lineBreaker -> out
type encoding struct {
	out        io.Writer
	breaker    *lineBreaker
	b64        io.WriteCloser
	crc        uint32
	crcEnabled bool
	blockType  []byte
}

func (e *encoding) Write(data []byte) (n int, err error) {
	if e.crcEnabled {
		e.crc = crc24(e.crc, data)
	}
	return e.b64.Write(data)
}

func (e *encoding) Close() (err error) {
	err = e.b64.Close()
	if err != nil {
		return
	}
	e.breaker.Close()

	if e.crcEnabled {
		var checksumBytes [3]byte
		checksumBytes[0] = byte(e.crc >> 16)
		checksumBytes[1] = byte(e.crc >> 8)
		checksumBytes[2] = byte(e.crc)

		var b64ChecksumBytes [4]byte
		base64.StdEncoding.Encode(b64ChecksumBytes[:], checksumBytes[:])

		return writeSlices(e.out, blockEnd, b64ChecksumBytes[:], newline, armorEnd, e.blockType, armorEndOfLine)
	}
	return writeSlices(e.out, newline, armorEnd, e.blockType, armorEndOfLine)
}

func encode(out io.Writer, blockType string, headers map[string]string, checksum bool) (w io.WriteCloser, err error) {
	bType := []byte(blockType)
	err = writeSlices(out, armorStart, bType, armorEndOfLineOut)
	if err != nil {
		return
	}

	keys := make([]string, len(headers))
	i := 0
	for k := range headers {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		err = writeSlices(out, []byte(k), armorHeaderSep, []byte(headers[k]), newline)
		if err != nil {
			return
		}
	}

	_, err = out.Write(newline)
	if err != nil {
		return
	}

	e := &encoding{
		out:        out,
		breaker:    newLineBreaker(out, 64),
		blockType:  bType,
		crc:        crc24Init,
		crcEnabled: checksum,
	}
	e.b64 = base64.NewEncoder(base64.StdEncoding, e.breaker)
	return e, nil
}

// Encode returns a WriteCloser which will encode the data written to it in
// OpenPGP armor.
func Encode(out io.Writer, blockType string, headers map[string]string) (w io.WriteCloser, err error) {
	return encode(out, blockType, headers, true)
}

// EncodeWithChecksumOption returns a WriteCloser which will encode the data written to it in
// OpenPGP armor and provides the option to include a checksum.
// When forming ASCII Armor, the CRC24 footer SHOULD NOT be generated,
// unless interoperability with implementations that require the CRC24 footer
// to be present is a concern.
func EncodeWithChecksumOption(out io.Writer, blockType string, headers map[string]string, doChecksum bool) (w io.WriteCloser, err error) {
	return encode(out, blockType, headers, doChecksum)
}
