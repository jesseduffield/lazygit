// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package randtxt

import (
	"bufio"
	"io"
	"unicode"
)

// GroupReader groups the incoming text in groups of 5, whereby the
// number of groups per line can be controlled.
type GroupReader struct {
	R             io.ByteReader
	GroupsPerLine int
	off           int64
	eof           bool
}

// NewGroupReader creates a new group reader.
func NewGroupReader(r io.Reader) *GroupReader {
	return &GroupReader{R: bufio.NewReader(r)}
}

// Read formats the data provided by the internal reader in groups of 5
// characters. If GroupsPerLine hasn't been initialized 8 groups per
// line will be produced.
func (r *GroupReader) Read(p []byte) (n int, err error) {
	if r.eof {
		return 0, io.EOF
	}
	groupsPerLine := r.GroupsPerLine
	if groupsPerLine < 1 {
		groupsPerLine = 8
	}
	lineLen := int64(groupsPerLine * 6)
	var c byte
	for i := range p {
		switch {
		case r.off%lineLen == lineLen-1:
			if i+1 == len(p) && len(p) > 1 {
				return i, nil
			}
			c = '\n'
		case r.off%6 == 5:
			if i+1 == len(p) && len(p) > 1 {
				return i, nil
			}
			c = ' '
		default:
			c, err = r.R.ReadByte()
			if err == io.EOF {
				r.eof = true
				if i > 0 {
					switch p[i-1] {
					case ' ':
						p[i-1] = '\n'
						fallthrough
					case '\n':
						return i, io.EOF
					}
				}
				p[i] = '\n'
				return i + 1, io.EOF
			}
			if err != nil {
				return i, err
			}
			switch {
			case c == ' ':
				c = '_'
			case !unicode.IsPrint(rune(c)):
				c = '-'
			}
		}
		p[i] = c
		r.off++
	}
	return len(p), nil
}
