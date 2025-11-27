//go:build plan9
// +build plan9

// Copyright 2025 The TCell Authors
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
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// p9Tty implements tcell.Tty using Plan 9's /dev/cons and /dev/consctl.
// Raw mode is enabled by writing "rawon" to /dev/consctl while the fd stays open.
// Resize notifications are read from /dev/wctl: the first read returns geometry,
// subsequent reads block until the window changes (rio(4)).
//
// References:
// - kbdfs(8): cons/consctl rawon|rawoff semantics
// - rio(4): wctl geometry and blocking-on-change behavior
// - vt(1): VT100 emulator typically used for TUI programs on Plan 9
//
// Limitations:
// - We assume VT100-level capabilities (often no colors, no mouse).
// - Window size is conservative: we return 80x24 unless overridden.
//   Set LINES/COLUMNS (or TCELL_LINES/TCELL_COLS) to refine.
// - Mouse and bracketed paste are not wired; terminfo/xterm queries
//   are not attempted because vt(1) may not support them.
type p9Tty struct {
	cons    *os.File // /dev/cons (read+write)
	consctl *os.File // /dev/consctl (write "rawon"/"rawoff")
	wctl    *os.File // /dev/wctl (resize notifications)

	// protect close/stop; Read/Write are serialized by os.File
	mu     sync.Mutex
	closed atomic.Bool

	// resize callback
	onResize atomic.Value // func()
	wg       sync.WaitGroup
	stopCh   chan struct{}
}

func NewDevTty() (Tty, error) { // tcell signature
	return newPlan9TTY()
}

func NewStdIoTty() (Tty, error) { // also required by tcell
	// On Plan 9 there is no POSIX tty discipline on stdin/stdout;
	// use /dev/cons explicitly for robustness.
	return newPlan9TTY()
}

func NewDevTtyFromDev(_ string) (Tty, error) { // required by tcell
	// Plan 9 does not have multiple "ttys" in the POSIX sense;
	// always bind to /dev/cons and /dev/consctl.
	return newPlan9TTY()
}

func newPlan9TTY() (Tty, error) {
	cons, err := os.OpenFile("/dev/cons", os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("open /dev/cons: %w", err)
	}
	consctl, err := os.OpenFile("/dev/consctl", os.O_WRONLY, 0)
	if err != nil {
		_ = cons.Close()
		return nil, fmt.Errorf("open /dev/consctl: %w", err)
	}
	// /dev/wctl may not exist (console without rio); best-effort.
	wctl, _ := os.OpenFile("/dev/wctl", os.O_RDWR, 0)

	t := &p9Tty{
		cons:    cons,
		consctl: consctl,
		wctl:    wctl,
		stopCh:  make(chan struct{}),
	}
	return t, nil
}

func (t *p9Tty) Start() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed.Load() {
		return errors.New("tty closed")
	}

	// Recreate stop channel if absent or closed (supports resume).
	if t.stopCh == nil || isClosed(t.stopCh) {
		t.stopCh = make(chan struct{})
	}

	// Put console into raw mode; remains active while consctl is open.
	if _, err := t.consctl.Write([]byte("rawon")); err != nil {
		return fmt.Errorf("enable raw mode: %w", err)
	}

	// Reopen /dev/wctl on resume; best-effort (system console may lack it).
	if t.wctl == nil {
		if f, err := os.OpenFile("/dev/wctl", os.O_RDWR, 0); err == nil {
			t.wctl = f
		}
	}

	if t.wctl != nil {
		t.wg.Add(1)
		go t.watchResize()
	}
	return nil
}

func (t *p9Tty) Drain() error {
	// Per tcell docs, this may reasonably be a no-op on non-POSIX ttys.
	// Read deadlines are not available on plan9 os.File; we rely on Stop().
	return nil
}

func (t *p9Tty) Stop() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Signal watcher to stop (if not already).
	if t.stopCh != nil && !isClosed(t.stopCh) {
		close(t.stopCh)
	}

	// Exit raw mode first.
	_, _ = t.consctl.Write([]byte("rawoff"))

	// Closing wctl unblocks watchResize; nil it so Start() can reopen later.
	if t.wctl != nil {
		_ = t.wctl.Close()
		t.wctl = nil
	}

	// Ensure watcher goroutine has exited before returning.
	t.wg.Wait()
	return nil
}

func (t *p9Tty) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed.Swap(true) {
		return nil
	}

	if t.stopCh != nil && !isClosed(t.stopCh) {
		close(t.stopCh)
	}
	_, _ = t.consctl.Write([]byte("rawoff"))

	_ = t.cons.Close()
	_ = t.consctl.Close()
	if t.wctl != nil {
		_ = t.wctl.Close()
		t.wctl = nil
	}

	t.wg.Wait()
	return nil
}

func (t *p9Tty) Read(p []byte) (int, error) {
	return t.cons.Read(p)
}

func (t *p9Tty) Write(p []byte) (int, error) {
	return t.cons.Write(p)
}

func (t *p9Tty) NotifyResize(cb func()) {
	if cb == nil {
		t.onResize.Store((func())(nil))
		return
	}
	t.onResize.Store(cb)
}

func (t *p9Tty) WindowSize() (WindowSize, error) {
	// Strategy:
	// 1) honor explicit overrides (TCELL_LINES/TCELL_COLS, LINES/COLUMNS),
	// 2) otherwise return conservative 80x24.
	// Reading /dev/wctl gives pixel geometry, but char cell metrics are
	// not generally available to non-draw clients; vt(1) is fixed-cell.
	lines, cols := envInt("TCELL_LINES"), envInt("TCELL_COLS")
	if lines == 0 {
		lines = envInt("LINES")
	}
	if cols == 0 {
		cols = envInt("COLUMNS")
	}
	if lines <= 0 {
		lines = 24
	}
	if cols <= 0 {
		cols = 80
	}
	return WindowSize{Width: cols, Height: lines}, nil
}

// watchResize blocks on /dev/wctl reads; each read returns when the window
// changes size/position/state, per rio(4). We ignore the parsed geometry and
// just notify tcell to re-query WindowSize().
func (t *p9Tty) watchResize() {
	defer t.wg.Done()

	r := bufio.NewReader(t.wctl)
	for {
		select {
		case <-t.stopCh:
			return
		default:
		}
		// Each read delivers something like:
		// "   minx        miny        maxx        maxy   visible current\n"
		// We don't need to parse here; just signal.
		_, err := r.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			// transient errors: continue
		}
		if cb, _ := t.onResize.Load().(func()); cb != nil {
			cb()
		}
	}
}

func envInt(name string) int {
	if s := strings.TrimSpace(os.Getenv(name)); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			return v
		}
	}
	return 0
}

// helper: safe check if a channel is closed
func isClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}
