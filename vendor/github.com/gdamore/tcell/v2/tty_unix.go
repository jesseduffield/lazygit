// Copyright 2021 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
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

package tcell

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

// devTty is an implementation of the Tty API based upon /dev/tty.
type devTty struct {
	fd    int
	f     *os.File
	of    *os.File // the first open of /dev/tty
	saved *term.State
	sig   chan os.Signal
	cb    func()
	stopQ chan struct{}
	dev   string
	wg    sync.WaitGroup
	l     sync.Mutex
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
	tty.l.Lock()
	defer tty.l.Unlock()

	// We open another copy of /dev/tty.  This is a workaround for unusual behavior
	// observed in macOS, apparently caused when a subshell (for example) closes our
	// own tty device (when it exits for example).  Getting a fresh new one seems to
	// resolve the problem.  (We believe this is a bug in the macOS tty driver that
	// fails to account for dup() references to the same file before applying close()
	// related behaviors to the tty.)  We're also holding the original copy we opened
	// since closing that might have deleterious effects as well.  The upshot is that
	// we will have up to two separate file handles open on /dev/tty.  (Note that when
	// using stdin/stdout instead of /dev/tty this problem is not observed.)
	var err error
	if tty.f, err = os.OpenFile(tty.dev, os.O_RDWR, 0); err != nil {
		return err
	}

	if !term.IsTerminal(tty.fd) {
		return errors.New("device is not a terminal")
	}

	_ = tty.f.SetReadDeadline(time.Time{})
	saved, err := term.MakeRaw(tty.fd) // also sets vMin and vTime
	if err != nil {
		return err
	}
	tty.saved = saved

	tty.stopQ = make(chan struct{})
	tty.wg.Add(1)
	go func(stopQ chan struct{}) {
		defer tty.wg.Done()
		for {
			select {
			case <-tty.sig:
				tty.l.Lock()
				cb := tty.cb
				tty.l.Unlock()
				if cb != nil {
					cb()
				}
			case <-stopQ:
				return
			}
		}
	}(tty.stopQ)

	signal.Notify(tty.sig, syscall.SIGWINCH)
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
	tty.l.Lock()
	if err := term.Restore(tty.fd, tty.saved); err != nil {
		tty.l.Unlock()
		return err
	}
	_ = tty.f.SetReadDeadline(time.Now())

	signal.Stop(tty.sig)
	close(tty.stopQ)
	tty.l.Unlock()

	tty.wg.Wait()

	// close our tty device -- we'll get another one if we Start again later.
	_ = tty.f.Close()

	return nil
}

func (tty *devTty) WindowSize() (WindowSize, error) {
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

func (tty *devTty) NotifyResize(cb func()) {
	tty.l.Lock()
	tty.cb = cb
	tty.l.Unlock()
}

// NewDevTty opens a /dev/tty based Tty.
func NewDevTty() (Tty, error) {
	return NewDevTtyFromDev("/dev/tty")
}

// NewDevTtyFromDev opens a tty device given a path.  This can be useful to bind to other nodes.
func NewDevTtyFromDev(dev string) (Tty, error) {
	tty := &devTty{
		dev: dev,
		sig: make(chan os.Signal),
	}
	var err error
	if tty.of, err = os.OpenFile(dev, os.O_RDWR, 0); err != nil {
		return nil, err
	}
	tty.fd = int(tty.of.Fd())
	if !term.IsTerminal(tty.fd) {
		_ = tty.f.Close()
		return nil, errors.New("not a terminal")
	}
	if tty.saved, err = term.GetState(tty.fd); err != nil {
		_ = tty.f.Close()
		return nil, fmt.Errorf("failed to get state: %w", err)
	}
	return tty, nil
}
