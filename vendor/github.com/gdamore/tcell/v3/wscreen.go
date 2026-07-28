// Copyright 2026 The TCell Authors
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

//go:build js && wasm
// +build js,wasm

package tcell

import (
	"errors"
	"io"
	"sync"
	"syscall/js"

	"github.com/gdamore/tcell/v3/tty"
)

// initialize installs the browser-backed TTY used by tScreen on js/wasm.
func (t *tScreen) initialize() error {
	if t.tty == nil {
		t.tty = newBrowserTty()
	}
	if t.term == "" {
		t.term = "ghostty-truecolor"
	}
	return nil
}

func getCharset() string {
	return "UTF-8"
}

type browserTty struct {
	mu      sync.Mutex
	cond    *sync.Cond
	started bool
	drained bool
	closed  bool
	input   []byte
	resizeQ chan<- bool

	writeFunc  js.Value
	sizeFunc   js.Value
	closeFuncs []js.Func
}

func newBrowserTty() *browserTty {
	t := &browserTty{}
	t.cond = sync.NewCond(&t.mu)
	return t
}

func (t *browserTty) Start() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.started {
		return nil
	}
	global := js.Global()
	t.writeFunc = global.Get("tcellWrite")
	t.sizeFunc = global.Get("tcellWindowSize")
	if t.writeFunc.Type() != js.TypeFunction || t.sizeFunc.Type() != js.TypeFunction {
		return errors.New("tcell wasm terminal host is not installed")
	}

	onData := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return nil
		}
		if args[0].InstanceOf(global.Get("Uint8Array")) {
			data := make([]byte, args[0].Get("byteLength").Int())
			js.CopyBytesToGo(data, args[0])
			t.enqueue(data)
		} else {
			t.enqueue([]byte(args[0].String()))
		}
		return nil
	})
	onResize := js.FuncOf(func(this js.Value, args []js.Value) any {
		t.mu.Lock()
		resizeQ := t.resizeQ
		t.mu.Unlock()
		if resizeQ != nil {
			select {
			case resizeQ <- true:
			default:
			}
		}
		return nil
	})
	t.closeFuncs = []js.Func{onData, onResize}
	global.Set("tcellRead", onData)
	global.Set("tcellResize", onResize)

	t.started = true
	t.drained = false
	t.closed = false
	return nil
}

func (t *browserTty) Stop() error {
	t.mu.Lock()
	t.started = false
	t.drained = false
	funcs := t.closeFuncs
	t.closeFuncs = nil
	t.cond.Broadcast()
	t.mu.Unlock()

	js.Global().Set("tcellRead", js.Undefined())
	js.Global().Set("tcellResize", js.Undefined())
	for _, fn := range funcs {
		fn.Release()
	}
	return nil
}

func (t *browserTty) Drain() error {
	t.mu.Lock()
	t.input = nil
	t.drained = true
	t.cond.Broadcast()
	t.mu.Unlock()
	return nil
}

func (t *browserTty) NotifyResize(resizeQ chan<- bool) {
	t.mu.Lock()
	t.resizeQ = resizeQ
	t.mu.Unlock()
}

func (t *browserTty) WindowSize() (tty.WindowSize, error) {
	var ws tty.WindowSize
	t.mu.Lock()
	sizeFunc := t.sizeFunc
	t.mu.Unlock()
	if sizeFunc.Type() != js.TypeFunction {
		ws.Width = 80
		ws.Height = 24
		return ws, nil
	}
	size := sizeFunc.Invoke()
	ws.Width = size.Get("cols").Int()
	ws.Height = size.Get("rows").Int()
	ws.PixelWidth = size.Get("pixelWidth").Int()
	ws.PixelHeight = size.Get("pixelHeight").Int()
	if ws.Width == 0 {
		ws.Width = 80
	}
	if ws.Height == 0 {
		ws.Height = 24
	}
	return ws, nil
}

func (t *browserTty) Read(b []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for len(t.input) == 0 && t.started && !t.drained && !t.closed {
		t.cond.Wait()
	}
	if t.closed {
		return 0, io.EOF
	}
	if (!t.started || t.drained) && len(t.input) == 0 {
		return 0, io.EOF
	}
	n := copy(b, t.input)
	t.input = t.input[n:]
	return n, nil
}

func (t *browserTty) Write(b []byte) (int, error) {
	t.mu.Lock()
	writeFunc := t.writeFunc
	started := t.started
	t.mu.Unlock()
	if !started || writeFunc.Type() != js.TypeFunction {
		return 0, io.ErrClosedPipe
	}

	data := js.Global().Get("Uint8Array").New(len(b))
	js.CopyBytesToJS(data, b)
	writeFunc.Invoke(data)
	return len(b), nil
}

func (t *browserTty) Close() error {
	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return nil
	}
	t.closed = true
	t.started = false
	t.drained = false
	t.cond.Broadcast()
	t.mu.Unlock()

	return t.Stop()
}

func (t *browserTty) enqueue(data []byte) {
	t.mu.Lock()
	if t.started && !t.closed {
		t.input = append(t.input, data...)
		t.cond.Broadcast()
	}
	t.mu.Unlock()
}
