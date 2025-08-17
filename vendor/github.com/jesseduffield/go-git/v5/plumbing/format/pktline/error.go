package pktline

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

var (
	// ErrInvalidErrorLine is returned by Decode when the packet line is not an
	// error line.
	ErrInvalidErrorLine = errors.New("expected an error-line")

	errPrefix = []byte("ERR ")
)

// ErrorLine is a packet line that contains an error message.
// Once this packet is sent by client or server, the data transfer process is
// terminated.
// See https://git-scm.com/docs/pack-protocol#_pkt_line_format
type ErrorLine struct {
	Text string
}

// Error implements the error interface.
func (e *ErrorLine) Error() string {
	return e.Text
}

// Encode encodes the ErrorLine into a packet line.
func (e *ErrorLine) Encode(w io.Writer) error {
	p := NewEncoder(w)
	return p.Encodef("%s%s\n", string(errPrefix), e.Text)
}

// Decode decodes a packet line into an ErrorLine.
func (e *ErrorLine) Decode(r io.Reader) error {
	s := NewScanner(r)
	if !s.Scan() {
		return s.Err()
	}

	line := s.Bytes()
	if !bytes.HasPrefix(line, errPrefix) {
		return ErrInvalidErrorLine
	}

	e.Text = strings.TrimSpace(string(line[4:]))
	return nil
}
