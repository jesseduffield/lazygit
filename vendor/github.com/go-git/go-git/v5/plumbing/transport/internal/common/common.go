// Package common implements the git pack protocol with a pluggable transport.
// This is a low-level package to implement new transports. Use a concrete
// implementation instead (e.g. http, file, ssh).
//
// A simple example of usage can be found in the file package.
package common

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/plumbing/format/pktline"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/capability"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/sideband"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/utils/ioutil"
)

const (
	readErrorSecondsTimeout = 10
)

var (
	ErrTimeoutExceeded = errors.New("timeout exceeded")
	// stdErrSkipPattern is used for skipping lines from a command's stderr output.
	// Any line matching this pattern will be skipped from further
	// processing and not be returned to calling code.
	stdErrSkipPattern = regexp.MustCompile("^remote:( =*){0,1}$")
)

// Commander creates Command instances. This is the main entry point for
// transport implementations.
type Commander interface {
	// Command creates a new Command for the given git command and
	// endpoint. cmd can be git-upload-pack or git-receive-pack. An
	// error should be returned if the endpoint is not supported or the
	// command cannot be created (e.g. binary does not exist, connection
	// cannot be established).
	Command(cmd string, ep *transport.Endpoint, auth transport.AuthMethod) (Command, error)
}

// Command is used for a single command execution.
// This interface is modeled after exec.Cmd and ssh.Session in the standard
// library.
type Command interface {
	// StderrPipe returns a pipe that will be connected to the command's
	// standard error when the command starts. It should not be called after
	// Start.
	StderrPipe() (io.Reader, error)
	// StdinPipe returns a pipe that will be connected to the command's
	// standard input when the command starts. It should not be called after
	// Start. The pipe should be closed when no more input is expected.
	StdinPipe() (io.WriteCloser, error)
	// StdoutPipe returns a pipe that will be connected to the command's
	// standard output when the command starts. It should not be called after
	// Start.
	StdoutPipe() (io.Reader, error)
	// Start starts the specified command. It does not wait for it to
	// complete.
	Start() error
	// Close closes the command and releases any resources used by it. It
	// will block until the command exits.
	Close() error
}

// CommandKiller expands the Command interface, enabling it for being killed.
type CommandKiller interface {
	// Kill and close the session whatever the state it is. It will block until
	// the command is terminated.
	Kill() error
}

type client struct {
	cmdr Commander
}

// NewClient creates a new client using the given Commander.
func NewClient(runner Commander) transport.Transport {
	return &client{runner}
}

// NewUploadPackSession creates a new UploadPackSession.
func (c *client) NewUploadPackSession(ep *transport.Endpoint, auth transport.AuthMethod) (
	transport.UploadPackSession, error) {

	return c.newSession(transport.UploadPackServiceName, ep, auth)
}

// NewReceivePackSession creates a new ReceivePackSession.
func (c *client) NewReceivePackSession(ep *transport.Endpoint, auth transport.AuthMethod) (
	transport.ReceivePackSession, error) {

	return c.newSession(transport.ReceivePackServiceName, ep, auth)
}

type session struct {
	Stdin   io.WriteCloser
	Stdout  io.Reader
	Command Command

	isReceivePack bool
	advRefs       *packp.AdvRefs
	packRun       bool
	finished      bool
	firstErrLine  chan string
}

func (c *client) newSession(s string, ep *transport.Endpoint, auth transport.AuthMethod) (*session, error) {
	cmd, err := c.cmdr.Command(s, ep, auth)
	if err != nil {
		return nil, err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &session{
		Stdin:         stdin,
		Stdout:        stdout,
		Command:       cmd,
		firstErrLine:  c.listenFirstError(stderr),
		isReceivePack: s == transport.ReceivePackServiceName,
	}, nil
}

func (c *client) listenFirstError(r io.Reader) chan string {
	if r == nil {
		return nil
	}

	errLine := make(chan string, 1)
	go func() {
		s := bufio.NewScanner(r)
		for {
			if s.Scan() {
				line := s.Text()
				if !stdErrSkipPattern.MatchString(line) {
					errLine <- line
					break
				}
			} else {
				close(errLine)
				break
			}
		}

		_, _ = io.Copy(io.Discard, r)
	}()

	return errLine
}

func (s *session) AdvertisedReferences() (*packp.AdvRefs, error) {
	return s.AdvertisedReferencesContext(context.TODO())
}

// AdvertisedReferences retrieves the advertised references from the server.
func (s *session) AdvertisedReferencesContext(ctx context.Context) (*packp.AdvRefs, error) {
	if s.advRefs != nil {
		return s.advRefs, nil
	}

	ar := packp.NewAdvRefs()
	if err := ar.Decode(s.StdoutContext(ctx)); err != nil {
		if err := s.handleAdvRefDecodeError(err); err != nil {
			return nil, err
		}
	}

	// Some servers like jGit, announce capabilities instead of returning an
	// packp message with a flush. This verifies that we received a empty
	// adv-refs, even it contains capabilities.
	if !s.isReceivePack && ar.IsEmpty() {
		return nil, transport.ErrEmptyRemoteRepository
	}

	transport.FilterUnsupportedCapabilities(ar.Capabilities)
	s.advRefs = ar
	return ar, nil
}

func (s *session) handleAdvRefDecodeError(err error) error {
	var errLine *pktline.ErrorLine
	if errors.As(err, &errLine) {
		if isRepoNotFoundError(errLine.Text) {
			return transport.ErrRepositoryNotFound
		}

		return errLine
	}

	// If repository is not found, we get empty stdout and server writes an
	// error to stderr.
	if errors.Is(err, packp.ErrEmptyInput) {
		// TODO:(v6): handle this error in a better way.
		// Instead of checking the stderr output for a specific error message,
		// define an ExitError and embed the stderr output and exit (if one
		// exists) in the error struct. Just like exec.ExitError.
		s.finished = true
		if err := s.checkNotFoundError(); err != nil {
			return err
		}

		return io.ErrUnexpectedEOF
	}

	// For empty (but existing) repositories, we get empty advertised-references
	// message. But valid. That is, it includes at least a flush.
	if err == packp.ErrEmptyAdvRefs {
		// Empty repositories are valid for git-receive-pack.
		if s.isReceivePack {
			return nil
		}

		if err := s.finish(); err != nil {
			return err
		}

		return transport.ErrEmptyRemoteRepository
	}

	// Some server sends the errors as normal content (git protocol), so when
	// we try to decode it fails, we need to check the content of it, to detect
	// not found errors
	if uerr, ok := err.(*packp.ErrUnexpectedData); ok {
		if isRepoNotFoundError(string(uerr.Data)) {
			return transport.ErrRepositoryNotFound
		}
	}

	return err
}

// UploadPack performs a request to the server to fetch a packfile. A reader is
// returned with the packfile content. The reader must be closed after reading.
func (s *session) UploadPack(ctx context.Context, req *packp.UploadPackRequest) (*packp.UploadPackResponse, error) {
	if req.IsEmpty() {
		// XXX: IsEmpty means haves are a subset of wants, in that case we have
		// everything we asked for. Close the connection and return nil.
		if err := s.finish(); err != nil {
			return nil, err
		}
		// TODO:(v6) return nil here
		return nil, transport.ErrEmptyUploadPackRequest
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	if _, err := s.AdvertisedReferencesContext(ctx); err != nil {
		return nil, err
	}

	s.packRun = true

	in := s.StdinContext(ctx)
	out := s.StdoutContext(ctx)

	if err := uploadPack(in, out, req); err != nil {
		return nil, err
	}

	r, err := ioutil.NonEmptyReader(out)
	if err == ioutil.ErrEmptyReader {
		if c, ok := s.Stdout.(io.Closer); ok {
			_ = c.Close()
		}

		return nil, transport.ErrEmptyUploadPackRequest
	}

	if err != nil {
		return nil, err
	}

	rc := ioutil.NewReadCloser(r, s)
	return DecodeUploadPackResponse(rc, req)
}

func (s *session) StdinContext(ctx context.Context) io.WriteCloser {
	return ioutil.NewWriteCloserOnError(
		ioutil.NewContextWriteCloser(ctx, s.Stdin),
		s.onError,
	)
}

func (s *session) StdoutContext(ctx context.Context) io.Reader {
	return ioutil.NewReaderOnError(
		ioutil.NewContextReader(ctx, s.Stdout),
		s.onError,
	)
}

func (s *session) onError(err error) {
	if k, ok := s.Command.(CommandKiller); ok {
		_ = k.Kill()
	}

	_ = s.Close()
}

func (s *session) ReceivePack(ctx context.Context, req *packp.ReferenceUpdateRequest) (*packp.ReportStatus, error) {
	if _, err := s.AdvertisedReferences(); err != nil {
		return nil, err
	}

	s.packRun = true

	w := s.StdinContext(ctx)
	if err := req.Encode(w); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	if !req.Capabilities.Supports(capability.ReportStatus) {
		// If we don't have report-status, we can only
		// check return value error.
		return nil, s.Command.Close()
	}

	r := s.StdoutContext(ctx)

	var d *sideband.Demuxer
	if req.Capabilities.Supports(capability.Sideband64k) {
		d = sideband.NewDemuxer(sideband.Sideband64k, r)
	} else if req.Capabilities.Supports(capability.Sideband) {
		d = sideband.NewDemuxer(sideband.Sideband, r)
	}
	if d != nil {
		d.Progress = req.Progress
		r = d
	}

	report := packp.NewReportStatus()
	if err := report.Decode(r); err != nil {
		return nil, err
	}

	if err := report.Error(); err != nil {
		defer s.Close()
		return report, err
	}

	return report, s.Command.Close()
}

func (s *session) finish() error {
	if s.finished {
		return nil
	}

	s.finished = true

	// If we did not run a upload/receive-pack, we close the connection
	// gracefully by sending a flush packet to the server. If the server
	// operates correctly, it will exit with status 0.
	if !s.packRun {
		_, err := s.Stdin.Write(pktline.FlushPkt)
		return err
	}

	return nil
}

func (s *session) Close() (err error) {
	err = s.finish()

	defer ioutil.CheckClose(s.Command, &err)
	return
}

func (s *session) checkNotFoundError() error {
	t := time.NewTicker(time.Second * readErrorSecondsTimeout)
	defer t.Stop()

	select {
	case <-t.C:
		return ErrTimeoutExceeded
	case line, ok := <-s.firstErrLine:
		if !ok || len(line) == 0 {
			return nil
		}

		if isRepoNotFoundError(line) {
			return transport.ErrRepositoryNotFound
		}

		// TODO:(v6): return server error just as it is without a prefix
		return fmt.Errorf("unknown error: %s", line)
	}
}

const (
	githubRepoNotFoundErr      = "Repository not found."
	bitbucketRepoNotFoundErr   = "repository does not exist."
	localRepoNotFoundErr       = "does not appear to be a git repository"
	gitProtocolNotFoundErr     = "Repository not found."
	gitProtocolNoSuchErr       = "no such repository"
	gitProtocolAccessDeniedErr = "access denied"
	gogsAccessDeniedErr        = "Repository does not exist or you do not have access"
	gitlabRepoNotFoundErr      = "The project you were looking for could not be found"
)

func isRepoNotFoundError(s string) bool {
	for _, err := range []string{
		githubRepoNotFoundErr,
		bitbucketRepoNotFoundErr,
		localRepoNotFoundErr,
		gitProtocolNotFoundErr,
		gitProtocolNoSuchErr,
		gitProtocolAccessDeniedErr,
		gogsAccessDeniedErr,
		gitlabRepoNotFoundErr,
	} {
		if strings.Contains(s, err) {
			return true
		}
	}

	return false
}

// uploadPack implements the git-upload-pack protocol.
func uploadPack(w io.WriteCloser, _ io.Reader, req *packp.UploadPackRequest) error {
	// TODO support multi_ack mode
	// TODO support multi_ack_detailed mode
	// TODO support acks for common objects
	// TODO build a proper state machine for all these processing options

	if err := req.UploadRequest.Encode(w); err != nil {
		return fmt.Errorf("sending upload-req message: %s", err)
	}

	if err := req.UploadHaves.Encode(w, true); err != nil {
		return fmt.Errorf("sending haves message: %s", err)
	}

	if err := sendDone(w); err != nil {
		return fmt.Errorf("sending done message: %s", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("closing input: %s", err)
	}

	return nil
}

func sendDone(w io.Writer) error {
	e := pktline.NewEncoder(w)

	return e.Encodef("done\n")
}

// DecodeUploadPackResponse decodes r into a new packp.UploadPackResponse
func DecodeUploadPackResponse(r io.ReadCloser, req *packp.UploadPackRequest) (
	*packp.UploadPackResponse, error,
) {
	res := packp.NewUploadPackResponse(req)
	if err := res.Decode(r); err != nil {
		return nil, fmt.Errorf("error decoding upload-pack response: %s", err)
	}

	return res, nil
}
