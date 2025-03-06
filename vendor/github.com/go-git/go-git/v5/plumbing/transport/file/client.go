// Package file implements the file transport protocol.
package file

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/internal/common"
	"golang.org/x/sys/execabs"
)

// DefaultClient is the default local client.
var DefaultClient = NewClient(
	transport.UploadPackServiceName,
	transport.ReceivePackServiceName,
)

type runner struct {
	UploadPackBin  string
	ReceivePackBin string
}

// NewClient returns a new local client using the given git-upload-pack and
// git-receive-pack binaries.
func NewClient(uploadPackBin, receivePackBin string) transport.Transport {
	return common.NewClient(&runner{
		UploadPackBin:  uploadPackBin,
		ReceivePackBin: receivePackBin,
	})
}

func prefixExecPath(cmd string) (string, error) {
	// Use `git --exec-path` to find the exec path.
	execCmd := execabs.Command("git", "--exec-path")

	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	stdoutBuf := bufio.NewReader(stdout)

	err = execCmd.Start()
	if err != nil {
		return "", err
	}

	execPathBytes, isPrefix, err := stdoutBuf.ReadLine()
	if err != nil {
		return "", err
	}
	if isPrefix {
		return "", errors.New("couldn't read exec-path line all at once")
	}

	err = execCmd.Wait()
	if err != nil {
		return "", err
	}
	execPath := string(execPathBytes)
	execPath = strings.TrimSpace(execPath)
	cmd = filepath.Join(execPath, cmd)

	// Make sure it actually exists.
	_, err = execabs.LookPath(cmd)
	if err != nil {
		return "", err
	}
	return cmd, nil
}

func (r *runner) Command(cmd string, ep *transport.Endpoint, auth transport.AuthMethod,
) (common.Command, error) {

	switch cmd {
	case transport.UploadPackServiceName:
		cmd = r.UploadPackBin
	case transport.ReceivePackServiceName:
		cmd = r.ReceivePackBin
	}

	_, err := execabs.LookPath(cmd)
	if err != nil {
		if e, ok := err.(*execabs.Error); ok && e.Err == execabs.ErrNotFound {
			cmd, err = prefixExecPath(cmd)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &command{cmd: execabs.Command(cmd, adjustPathForWindows(ep.Path))}, nil
}

func isDriveLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// On Windows, the path that results from a file: URL has a leading slash. This
// has to be removed if there's a drive letter
func adjustPathForWindows(p string) string {
	if runtime.GOOS != "windows" {
		return p
	}
	if len(p) >= 3 && p[0] == '/' && isDriveLetter(p[1]) && p[2] == ':' {
		return p[1:]
	}
	return p
}

type command struct {
	cmd          *execabs.Cmd
	stderrCloser io.Closer
	closed       bool
}

func (c *command) Start() error {
	return c.cmd.Start()
}

func (c *command) StderrPipe() (io.Reader, error) {
	// Pipe returned by Command.StderrPipe has a race with Read + Command.Wait.
	// We use an io.Pipe and close it after the command finishes.
	r, w := io.Pipe()
	c.cmd.Stderr = w
	c.stderrCloser = r
	return r, nil
}

func (c *command) StdinPipe() (io.WriteCloser, error) {
	return c.cmd.StdinPipe()
}

func (c *command) StdoutPipe() (io.Reader, error) {
	return c.cmd.StdoutPipe()
}

func (c *command) Kill() error {
	c.cmd.Process.Kill()
	return c.Close()
}

// Close waits for the command to exit.
func (c *command) Close() error {
	if c.closed {
		return nil
	}

	defer func() {
		c.closed = true
		_ = c.stderrCloser.Close()

	}()

	err := c.cmd.Wait()
	if _, ok := err.(*os.PathError); ok {
		return nil
	}

	// When a repository does not exist, the command exits with code 128.
	if _, ok := err.(*execabs.ExitError); ok {
		return nil
	}

	return err
}
