// Copyright 2025 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos
// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris zos

package tty

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

// Use -1 to differentiate from stdin (fd=0).
const uninitializedTtyFd = -1

// devTty is an implementation of the Tty API based upon /dev/tty.
type devTty struct {
	fd      int
	f       *os.File
	saved   *term.State
	sig     chan os.Signal
	dev     string
	started bool
}

func (tty *devTty) Read(b []byte) (int, error) {
	return tty.f.Read(b)
}

func (tty *devTty) Write(b []byte) (int, error) {
	return tty.f.Write(b)
}

func (tty *devTty) Close() error {
	return tty.f.Close()
}

func (tty *devTty) Start() error {

	if tty.started {
		return nil
	}
	// We open another copy of /dev/tty.  This is a workaround for unusual behavior
	// observed in macOS, apparently caused when a subshell (for example) closes our
	// own tty device (when it exits for example).  Getting a fresh new one seems to
	// resolve the problem.  (We believe this is a bug in the macOS tty driver that
	// fails to account for dup() references to the same file before applying close()
	// related behaviors to the tty.)  (Note that when using stdin/stdout instead of
	// /dev/tty this problem is not observed.)
	var err error
	if tty.f, err = os.OpenFile(tty.dev, os.O_RDWR, 0); err != nil {
		return err
	}

	tty.fd = int(tty.f.Fd())

	if !term.IsTerminal(tty.fd) {
		tty.f.Close()
		return errors.New("device is not a terminal")
	}

	_ = tty.f.SetReadDeadline(time.Time{})
	saved, err := term.MakeRaw(tty.fd) // also sets vMin and vTime
	if err != nil {
		tty.f.Close()
		return err
	}
	if err = tcFlushInput(tty.fd); err != nil {
		_ = term.Restore(tty.fd, saved)
		tty.f.Close()
		return err
	}
	tty.saved = saved
	tty.started = true

	return nil
}

func (tty *devTty) Drain() error {
	_ = tty.f.SetReadDeadline(time.Now())
	if err := tcSetBufParams(tty.fd, 0, 0); err != nil {
		return err
	}
	return nil
}

func (tty *devTty) Stop() error {
	// unconditionally set this, because we cannot recover
	// if we fail anyway, so this gives the best hope of
	// picking up the pieces in such a circumstance
	tty.started = false

	if err := term.Restore(tty.fd, tty.saved); err != nil {
		return err
	}
	_ = tty.f.SetReadDeadline(time.Now())

	tty.NotifyResize(nil)

	// close our tty device -- we'll get another one if we Start again later.
	_ = tty.f.Close()
	tty.fd = uninitializedTtyFd

	return nil
}

func (tty *devTty) WindowSize() (WindowSize, error) {
	size := WindowSize{}
	fd := tty.fd
	if tty.fd == uninitializedTtyFd {
		// If WindowSize is called when the tty isn't yet running, the fd for /dev/tty won't be initialized,
		// so open the file just long enough to retrieve the window size.
		f, err := os.OpenFile(tty.dev, os.O_RDWR, 0)
		if err != nil {
			return size, err
		}
		defer func() { _ = f.Close() }()
		fd = int(f.Fd())
	}

	ws, err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	if err != nil {
		return size, err
	}
	w := int(ws.Col)
	h := int(ws.Row)
	if w == 0 {
		w, _ = strconv.Atoi(os.Getenv("COLUMNS"))
	}
	if w == 0 {
		w = 80 // default
	}
	if h == 0 {
		h, _ = strconv.Atoi(os.Getenv("LINES"))
	}
	if h == 0 {
		h = 25 // default
	}
	size.Width = w
	size.Height = h
	size.PixelWidth = int(ws.Xpixel)
	size.PixelHeight = int(ws.Ypixel)
	return size, nil
}

func (tty *devTty) NotifyResize(resizeQ chan<- bool) {

	sigQ := tty.sig
	tty.sig = nil

	if sigQ != nil {
		signal.Stop(sigQ)
		close(sigQ)
	}

	if resizeQ == nil {
		return
	}

	sigQ = make(chan os.Signal, 1)
	signal.Notify(sigQ, syscall.SIGWINCH)

	tty.sig = sigQ

	go func() {
		for range sigQ {
			select {
			case resizeQ <- true:
			default: // queue full, so nvm.
			}
		}
	}()
}

// NewDevTty opens a /dev/tty based Tty.
func NewDevTty() (Tty, error) {
	return NewDevTtyFromDev("/dev/tty")
}

// NewDevTtyFromDev opens a tty device given a path.  This can be useful to bind to other nodes.
func NewDevTtyFromDev(dev string) (Tty, error) {
	tty := &devTty{
		fd:  uninitializedTtyFd,
		dev: dev,
		sig: make(chan os.Signal),
	}
	// Only open the file long enough to check that the device
	// represents a TTY.  We will reopen it in start.  We do collect
	// the terminal state so we can restore it later though.
	if f, err := os.OpenFile(dev, os.O_RDWR, 0); err != nil {
		return nil, err
	} else {
		defer func() { _ = f.Close() }()
		fd := int(f.Fd())
		if !term.IsTerminal(fd) {
			return nil, errors.New("not a terminal")
		}
		if tty.saved, err = term.GetState(fd); err != nil {
			return nil, fmt.Errorf("failed to get state: %w", err)
		}
	}
	return tty, nil
}
