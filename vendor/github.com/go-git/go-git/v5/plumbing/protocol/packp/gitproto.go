package packp

import (
	"fmt"
	"io"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/pktline"
)

var (
	// ErrInvalidGitProtoRequest is returned by Decode if the input is not a
	// valid git protocol request.
	ErrInvalidGitProtoRequest = fmt.Errorf("invalid git protocol request")
)

// GitProtoRequest is a command request for the git protocol.
// It is used to send the command, endpoint, and extra parameters to the
// remote.
// See https://git-scm.com/docs/pack-protocol#_git_transport
type GitProtoRequest struct {
	RequestCommand string
	Pathname       string

	// Optional
	Host string

	// Optional
	ExtraParams []string
}

// validate validates the request.
func (g *GitProtoRequest) validate() error {
	if g.RequestCommand == "" {
		return fmt.Errorf("%w: empty request command", ErrInvalidGitProtoRequest)
	}

	if g.Pathname == "" {
		return fmt.Errorf("%w: empty pathname", ErrInvalidGitProtoRequest)
	}

	return nil
}

// Encode encodes the request into the writer.
func (g *GitProtoRequest) Encode(w io.Writer) error {
	if w == nil {
		return ErrNilWriter
	}

	if err := g.validate(); err != nil {
		return err
	}

	p := pktline.NewEncoder(w)
	req := fmt.Sprintf("%s %s\x00", g.RequestCommand, g.Pathname)
	if host := g.Host; host != "" {
		req += fmt.Sprintf("host=%s\x00", host)
	}

	if len(g.ExtraParams) > 0 {
		req += "\x00"
		for _, param := range g.ExtraParams {
			req += param + "\x00"
		}
	}

	if err := p.Encode([]byte(req)); err != nil {
		return err
	}

	return nil
}

// Decode decodes the request from the reader.
func (g *GitProtoRequest) Decode(r io.Reader) error {
	s := pktline.NewScanner(r)
	if !s.Scan() {
		err := s.Err()
		if err == nil {
			return ErrInvalidGitProtoRequest
		}
		return err
	}

	line := string(s.Bytes())
	if len(line) == 0 {
		return io.EOF
	}

	if line[len(line)-1] != 0 {
		return fmt.Errorf("%w: missing null terminator", ErrInvalidGitProtoRequest)
	}

	parts := strings.SplitN(line, " ", 2)
	if len(parts) != 2 {
		return fmt.Errorf("%w: short request", ErrInvalidGitProtoRequest)
	}

	g.RequestCommand = parts[0]
	params := strings.Split(parts[1], string(null))
	if len(params) < 1 {
		return fmt.Errorf("%w: missing pathname", ErrInvalidGitProtoRequest)
	}

	g.Pathname = params[0]
	if len(params) > 1 {
		g.Host = strings.TrimPrefix(params[1], "host=")
	}

	if len(params) > 2 {
		for _, param := range params[2:] {
			if param != "" {
				g.ExtraParams = append(g.ExtraParams, param)
			}
		}
	}

	return nil
}
