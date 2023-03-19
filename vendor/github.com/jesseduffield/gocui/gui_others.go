// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !windows
// +build !windows

package gocui

import (
	"os"
	"os/signal"
	"syscall"
	"unsafe"

	"github.com/go-errors/errors"
)

// getTermWindowSize is get terminal window size on linux or unix.
// When gocui run inside the docker contaienr need to check and get the window size.
func (g *Gui) getTermWindowSize() (int, int, error) {
	var sz struct {
		rows uint16
		cols uint16
		_    [2]uint16 // to match underlying syscall; see https://github.com/awesome-gocui/gocui/issues/33
	}

	var termw, termh int

	out, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return 0, 0, err
	}
	defer out.Close()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGWINCH, syscall.SIGINT)

	for {
		_, _, _ = syscall.Syscall(syscall.SYS_IOCTL,
			out.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&sz)))

		// check terminal window size
		termw, termh = int(sz.cols), int(sz.rows)
		if termw > 0 && termh > 0 {
			return termw, termh, nil
		}

		select {
		case signal := <-signalCh:
			switch signal {
			// when the terminal window size is changed
			case syscall.SIGWINCH:
				continue
			// ctrl + c to cancel
			case syscall.SIGINT:
				return 0, 0, errors.New("stop to get term window size")
			}
		}
	}
}
