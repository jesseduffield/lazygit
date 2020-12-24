// +build linux

// Copyright 2019 The TCell Authors
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

package tcell

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"golang.org/x/sys/unix"
)

type termiosPrivate struct {
	tio *unix.Termios
}

func (t *tScreen) termioInit() error {
	var e error
	var raw *unix.Termios
	var tio *unix.Termios

	if t.in, e = os.OpenFile("/dev/tty", os.O_RDONLY, 0); e != nil {
		goto failed
	}
	if t.out, e = os.OpenFile("/dev/tty", os.O_WRONLY, 0); e != nil {
		goto failed
	}

	tio, e = unix.IoctlGetTermios(int(t.out.Fd()), unix.TCGETS)
	if e != nil {
		goto failed
	}

	t.tiosp = &termiosPrivate{tio: tio}

	// make a local copy, to make it raw
	raw = &unix.Termios{
		Cflag: tio.Cflag,
		Oflag: tio.Oflag,
		Iflag: tio.Iflag,
		Lflag: tio.Lflag,
		Cc:    tio.Cc,
	}
	raw.Iflag &^= (unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP |
		unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON)
	raw.Oflag &^= unix.OPOST
	raw.Lflag &^= (unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG |
		unix.IEXTEN)
	raw.Cflag &^= (unix.CSIZE | unix.PARENB)
	raw.Cflag |= unix.CS8

	// This is setup for blocking reads.  In the past we attempted to
	// use non-blocking reads, but now a separate input loop and timer
	// copes with the problems we had on some systems (BSD/Darwin)
	// where close hung forever.
	raw.Cc[unix.VMIN] = 1
	raw.Cc[unix.VTIME] = 0

	e = unix.IoctlSetTermios(int(t.out.Fd()), unix.TCSETS, raw)
	if e != nil {
		goto failed
	}

	signal.Notify(t.sigwinch, syscall.SIGWINCH)

	if w, h, e := t.getWinSize(); e == nil && w != 0 && h != 0 {
		t.cells.Resize(w, h)
	}

	return nil

failed:
	if t.in != nil {
		t.in.Close()
	}
	if t.out != nil {
		t.out.Close()
	}
	return e
}

func (t *tScreen) termioFini() {

	signal.Stop(t.sigwinch)

	<-t.indoneq

	if t.out != nil && t.tiosp != nil {
		unix.IoctlSetTermios(int(t.out.Fd()), unix.TCSETSF, t.tiosp.tio)
		t.out.Close()
	}

	if t.in != nil {
		t.in.Close()
	}
}

func (t *tScreen) getWinSize() (int, int, error) {

	wsz, err := unix.IoctlGetWinsize(int(t.out.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return -1, -1, err
	}
	cols := int(wsz.Col)
	rows := int(wsz.Row)
	if cols == 0 {
		colsEnv := os.Getenv("COLUMNS")
		if colsEnv != "" {
			if cols, err = strconv.Atoi(colsEnv); err != nil {
				return -1, -1, err
			}
		} else {
			cols = t.ti.Columns
		}
	}
	if rows == 0 {
		rowsEnv := os.Getenv("LINES")
		if rowsEnv != "" {
			if rows, err = strconv.Atoi(rowsEnv); err != nil {
				return -1, -1, err
			}
		} else {
			rows = t.ti.Lines
		}
	}
	return cols, rows, nil
}

func (t *tScreen) Beep() error {
	t.writeString(string(byte(7)))
	return nil
}
