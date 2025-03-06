package packp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/pktline"
)

const ackLineLen = 44

// ServerResponse object acknowledgement from upload-pack service
type ServerResponse struct {
	ACKs []plumbing.Hash
}

// Decode decodes the response into the struct, isMultiACK should be true, if
// the request was done with multi_ack or multi_ack_detailed capabilities.
func (r *ServerResponse) Decode(reader *bufio.Reader, isMultiACK bool) error {
	s := pktline.NewScanner(reader)

	for s.Scan() {
		line := s.Bytes()

		if err := r.decodeLine(line); err != nil {
			return err
		}

		// we need to detect when the end of a response header and the beginning
		// of a packfile header happened, some requests to the git daemon
		// produces a duplicate ACK header even when multi_ack is not supported.
		stop, err := r.stopReading(reader)
		if err != nil {
			return err
		}

		if stop {
			break
		}
	}

	// isMultiACK is true when the remote server advertises the related
	// capabilities when they are not in transport.UnsupportedCapabilities.
	//
	// Users may decide to remove multi_ack and multi_ack_detailed from the
	// unsupported capabilities list, which allows them to do initial clones
	// from Azure DevOps.
	//
	// Follow-up fetches may error, therefore errors are wrapped with additional
	// information highlighting that this capabilities are not supported by go-git.
	//
	// TODO: Implement support for multi_ack or multi_ack_detailed responses.
	err := s.Err()
	if err != nil && isMultiACK {
		return fmt.Errorf("multi_ack and multi_ack_detailed are not supported: %w", err)
	}

	return err
}

// stopReading detects when a valid command such as ACK or NAK is found to be
// read in the buffer without moving the read pointer.
func (r *ServerResponse) stopReading(reader *bufio.Reader) (bool, error) {
	ahead, err := reader.Peek(7)
	if err == io.EOF {
		return true, nil
	}

	if err != nil {
		return false, err
	}

	if len(ahead) > 4 && r.isValidCommand(ahead[0:3]) {
		return false, nil
	}

	if len(ahead) == 7 && r.isValidCommand(ahead[4:]) {
		return false, nil
	}

	return true, nil
}

func (r *ServerResponse) isValidCommand(b []byte) bool {
	commands := [][]byte{ack, nak}
	for _, c := range commands {
		if bytes.Equal(b, c) {
			return true
		}
	}

	return false
}

func (r *ServerResponse) decodeLine(line []byte) error {
	if len(line) == 0 {
		return fmt.Errorf("unexpected flush")
	}

	if len(line) >= 3 {
		if bytes.Equal(line[0:3], ack) {
			return r.decodeACKLine(line)
		}

		if bytes.Equal(line[0:3], nak) {
			return nil
		}
	}

	return fmt.Errorf("unexpected content %q", string(line))
}

func (r *ServerResponse) decodeACKLine(line []byte) error {
	if len(line) < ackLineLen {
		return fmt.Errorf("malformed ACK %q", line)
	}

	sp := bytes.Index(line, []byte(" "))
	if sp+41 > len(line) {
		return fmt.Errorf("malformed ACK %q", line)
	}
	h := plumbing.NewHash(string(line[sp+1 : sp+41]))
	r.ACKs = append(r.ACKs, h)
	return nil
}

// Encode encodes the ServerResponse into a writer.
func (r *ServerResponse) Encode(w io.Writer, isMultiACK bool) error {
	if len(r.ACKs) > 1 && !isMultiACK {
		// For further information, refer to comments in the Decode func above.
		return errors.New("multi_ack and multi_ack_detailed are not supported")
	}

	e := pktline.NewEncoder(w)
	if len(r.ACKs) == 0 {
		return e.Encodef("%s\n", nak)
	}

	return e.Encodef("%s %s\n", ack, r.ACKs[0].String())
}
