// Copyright 2021 The TCell Authors
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

// stdIoTty is an implementation of the Tty API based upon stdin/stdout.
type stdIoTty struct {
	fd      int
	in      *os.File
	out     *os.File
	saved   *term.State
	sig     chan os.Signal
	started bool
}

func (tty *stdIoTty) Read(b []byte) (int, error) {
	return tty.in.Read(b)
}

func (tty *stdIoTty) Write(b []byte) (int, error) {
	return tty.out.Write(b)
}

func (tty *stdIoTty) Close() error {
	return nil
}

func (tty *stdIoTty) Start() error {
	if tty.started {
		return nil
	}

	var err error
	tty.in = os.Stdin
	tty.out = os.Stdout
	tty.fd = int(tty.in.Fd())

	if !term.IsTerminal(tty.fd) {
		return errors.New("device is not a terminal")
	}

	_ = tty.in.SetReadDeadline(time.Time{})
	saved, err := term.MakeRaw(tty.fd) // also sets vMin and vTime
	if err != nil {
		return err
	}
	if err = tcFlushInput(tty.fd); err != nil {
		_ = term.Restore(tty.fd, saved)
		return err
	}
	tty.saved = saved
	tty.started = true

	return nil
}

func (tty *stdIoTty) Drain() error {
	_ = tty.in.SetReadDeadline(time.Now())
	if err := tcSetBufParams(tty.fd, 0, 0); err != nil {
		return err
	}
	return nil
}

func (tty *stdIoTty) Stop() error {
	if err := term.Restore(tty.fd, tty.saved); err != nil {
		return err
	}
	_ = tty.in.SetReadDeadline(time.Now())

	tty.NotifyResize(nil)

	tty.started = false

	return nil
}

func (tty *stdIoTty) WindowSize() (WindowSize, error) {
	size := WindowSize{}
	ws, err := unix.IoctlGetWinsize(tty.fd, unix.TIOCGWINSZ)
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

func (tty *stdIoTty) NotifyResize(resizeQ chan<- bool) {

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

// NewStdioTty opens a tty using standard input/output.
func NewStdIoTty() (Tty, error) {
	tty := &stdIoTty{
		sig: make(chan os.Signal),
		in:  os.Stdin,
		out: os.Stdout,
	}
	var err error
	tty.fd = int(tty.in.Fd())
	if !term.IsTerminal(tty.fd) {
		return nil, errors.New("not a terminal")
	}
	if tty.saved, err = term.GetState(tty.fd); err != nil {
		return nil, fmt.Errorf("failed to get state: %w", err)
	}
	return tty, nil
}
