//go:build windows
// +build windows

package oscommands

import (
	"bytes"
	"io"
	"os/exec"

	"github.com/sasha-s/go-deadlock"
)

type Buffer struct {
	b bytes.Buffer
	m deadlock.Mutex
}

func (b *Buffer) Read(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Read(p)
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Write(p)
}

// TODO: Remove this hack and replace it with a proper way to run commands live on windows. We still have an issue where if a password is requested, the request for a password is written straight to stdout because we can't control the stdout of a subprocess of a subprocess. Keep an eye on https://github.com/creack/pty/pull/109
func (self *cmdObjRunner) getCmdHandler(cmd *exec.Cmd) (*cmdHandler, error) {
	stdoutReader, stdoutWriter := io.Pipe()
	cmd.Stdout = stdoutWriter

	buf := &Buffer{}
	cmd.Stdin = buf

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// because we don't yet have windows support for a pty, we instead just
	// pass our standard stream handlers and because there's no pty to close
	// we pass a no-op function for that.
	return &cmdHandler{
		stdoutPipe: stdoutReader,
		stdinPipe:  buf,
		close:      func() error { return nil },
	}, nil
}
