// Copyright 2026 The TCell Authors
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

//go:build windows
// +build windows

package tcell

import (
	"encoding/binary"
	"errors"
	"sync"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"
)

var (
	k32 = syscall.NewLazyDLL("kernel32.dll")
)

var (
	procReadConsoleInput              = k32.NewProc("ReadConsoleInputW")
	procGetNumberOfConsoleInputEvents = k32.NewProc("GetNumberOfConsoleInputEvents")
	procFlushConsoleInputBuffer       = k32.NewProc("FlushConsoleInputBuffer")
	procWaitForMultipleObjects        = k32.NewProc("WaitForMultipleObjects")
	procSetConsoleMode                = k32.NewProc("SetConsoleMode")
	procGetConsoleMode                = k32.NewProc("GetConsoleMode")
	procGetConsoleScreenBufferInfo    = k32.NewProc("GetConsoleScreenBufferInfo")
	procCreateEvent                   = k32.NewProc("CreateEventW")
	procSetEvent                      = k32.NewProc("SetEvent")
)

const (
	keyEvent    uint16 = 1
	mouseEvent  uint16 = 2
	resizeEvent uint16 = 4
	menuEvent   uint16 = 8 // don't use
	focusEvent  uint16 = 16
)

type inputRecord struct {
	typ  uint16
	_    uint16
	data [16]byte
}

type winTty struct {
	buf        chan byte
	out        syscall.Handle
	in         syscall.Handle
	cancelFlag syscall.Handle
	running    bool
	stopQ      chan struct{}
	resizeCb   func()
	cols       uint16
	rows       uint16
	pair       []uint16 // for surrogate pairs (UTF-16)
	oimode     uint32   // original input mode
	oomode     uint32   // original output mode
	oscreen    consoleInfo
	wg         sync.WaitGroup
	surrogate  rune
	sync.Mutex
}

func (w *winTty) Read(b []byte) (int, error) {
	// first character read blocks
	var num int
	select {
	case c := <-w.buf:
		b[0] = c
		num++
	case <-w.stopQ:
		// stopping, so make sure we eat everything, which might require
		// very short sleeps to ensure all buffered data is consumed.
		break
	}

	// second character read is non-blocking
	for ; num < len(b); num++ {
		select {
		case c := <-w.buf:
			b[num] = c
		case <-time.After(time.Millisecond * 10):
			return num, nil
		}
	}
	return num, nil
}

func (w *winTty) Write(b []byte) (int, error) {
	esc := utf16.Encode([]rune(string(b)))
	if len(esc) > 0 {
		err := syscall.WriteConsole(w.out, &esc[0], uint32(len(esc)), nil, nil)
		if err != nil {
			return 0, err
		}
	}
	return len(b), nil
}

func (w *winTty) Close() error {
	_ = syscall.Close(w.in)
	_ = syscall.Close(w.out)
	return nil
}

func (w *winTty) Drain() error {
	close(w.stopQ)
	time.Sleep(time.Millisecond * 10)
	_, _, _ = procSetEvent.Call(uintptr(w.cancelFlag))
	return nil
}

func (w *winTty) getConsoleInput() error {
	// cancelFlag comes first as WaitForMultipleObjects returns the lowest index
	// in the event that both events are signaled.
	waitObjects := []syscall.Handle{w.cancelFlag, w.in}

	// As arrays are contiguous in memory, a pointer to the first object is the
	// same as a pointer to the array itself.
	pWaitObjects := unsafe.Pointer(&waitObjects[0])

	rv, _, er := procWaitForMultipleObjects.Call(
		uintptr(len(waitObjects)),
		uintptr(pWaitObjects),
		uintptr(0),
		w32Infinite)

	// WaitForMultipleObjects returns WAIT_OBJECT_0 + the index.
	switch rv {
	case w32WaitObject0: // w.cancelFlag
		return errors.New("cancelled")
	case w32WaitObject0 + 1: // w.in
		// rec := &inputRecord{}
		var nrec int32
		rv, _, er := procGetNumberOfConsoleInputEvents.Call(
			uintptr(w.in),
			uintptr(unsafe.Pointer(&nrec)))
		rec := make([]inputRecord, nrec)
		rv, _, er = procReadConsoleInput.Call(
			uintptr(w.in),
			uintptr(unsafe.Pointer(&rec[0])),
			uintptr(nrec),
			uintptr(unsafe.Pointer(&nrec)))
		if rv == 0 {
			return er
		}
	loop:
		for i := range nrec {
			ir := rec[i]
			switch ir.typ {
			case keyEvent:
				// we normally only expect to see ascii, but paste data may come in as UTF-16.
				wc := rune(binary.LittleEndian.Uint16(ir.data[10:]))
				if wc >= 0xD800 && wc <= 0xDBFF {
					// if it was a high surrogate, which happens for pasted UTF-16,
					// then save it until we get the low and can decode it.
					w.surrogate = wc
					continue
				} else if wc >= 0xDC00 && wc <= 0xDFFF {
					wc = utf16.DecodeRune(w.surrogate, wc)
				}
				w.surrogate = 0
				for _, chr := range []byte(string(wc)) {
					// We normally expect only to see ASCII (win32-input-mode),
					// but apparently pasted data can arrive in UTF-16 here.
					select {
					case w.buf <- chr:
					case <-w.stopQ:
						break loop
					}
				}

			case resizeEvent:
				w.Lock()
				w.cols = binary.LittleEndian.Uint16(ir.data[0:])
				w.rows = binary.LittleEndian.Uint16(ir.data[2:])
				cb := w.resizeCb
				w.Unlock()
				if cb != nil {
					cb()
				}

			default:
			}
		}
		return nil
	default:
		return er
	}
}

func (w *winTty) scanInput() {
	defer w.wg.Done()
	for {
		if e := w.getConsoleInput(); e != nil {
			return
		}
	}
}

func (w *winTty) Start() error {

	w.Lock()
	defer w.Unlock()

	if w.running {
		return errors.New("already engaged")
	}
	_, _, _ = procFlushConsoleInputBuffer.Call(uintptr(w.in))
	w.stopQ = make(chan struct{})
	cf, _, err := procCreateEvent.Call(
		uintptr(0),
		uintptr(1),
		uintptr(0),
		uintptr(0))
	if cf == uintptr(0) {
		return err
	}
	w.running = true
	w.cancelFlag = syscall.Handle(cf)

	_, _, _ = procSetConsoleMode.Call(uintptr(w.in),
		uintptr(modeVtInput|modeResizeEn|modeExtendFlg))
	_, _, _ = procSetConsoleMode.Call(uintptr(w.out),
		uintptr(modeVtOutput|modeNoAutoNL|modeCookedOut|modeUnderline))

	w.wg.Add(1)
	go w.scanInput()
	return nil
}

func (w *winTty) Stop() error {
	w.wg.Wait()
	w.Lock()
	defer w.Unlock()
	_, _, _ = procSetConsoleMode.Call(uintptr(w.in), uintptr(w.oimode))
	_, _, _ = procSetConsoleMode.Call(uintptr(w.out), uintptr(w.oomode))
	_, _, _ = procFlushConsoleInputBuffer.Call(uintptr(w.in))
	w.running = false

	return nil
}

func (w *winTty) NotifyResize(cb func()) {
	w.resizeCb = cb
}

func (w *winTty) WindowSize() (WindowSize, error) {
	w.Lock()
	defer w.Unlock()
	return WindowSize{Width: int(w.cols), Height: int(w.rows)}, nil
}

func NewDevTty() (Tty, error) {
	w := &winTty{}
	var err error
	w.in, err = syscall.Open("CONIN$", syscall.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	w.out, err = syscall.Open("CONOUT$", syscall.O_RDWR, 0)
	if err != nil {
		_ = syscall.Close(w.in)
		return nil, err
	}
	w.buf = make(chan byte, 128)

	_, _, _ = procGetConsoleScreenBufferInfo.Call(uintptr(w.out), uintptr(unsafe.Pointer(&w.oscreen)))
	_, _, _ = procGetConsoleMode.Call(uintptr(w.out), uintptr(unsafe.Pointer(&w.oomode)))
	_, _, _ = procGetConsoleMode.Call(uintptr(w.in), uintptr(unsafe.Pointer(&w.oimode)))
	w.rows = uint16(w.oscreen.size.y)
	w.cols = uint16(w.oscreen.size.x)

	return w, nil
}
