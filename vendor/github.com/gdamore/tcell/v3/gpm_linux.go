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

// This file implements a client for the GPM (general purpose mouse daemon)
// used to orchestrate mouse access and deliver mouse events to applications
// using the Linux console.  This deliberately speaks the small public GPM ABI
// directly rather than linking against libgpm, so tcell does not acquire GPL
// licensing obligations from the library.

//go:build linux

package tcell

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const (
	gpmConnectSize = 16
	gpmEventSize   = 28

	gpmButtonDown   = 32
	gpmButtonUp     = 16
	gpmButtonFourth = 8
	gpmButtonLeft   = 4
	gpmButtonMiddle = 2
	gpmButtonRight  = 1

	gpmModShift = 1
	gpmModAltGr = 2
	gpmModCtrl  = 4
	gpmModAlt   = 8

	gpmEventMove   = 1
	gpmEventDrag   = 2
	gpmEventDown   = 4
	gpmEventUp     = 8
	gpmEventSingle = 16
	gpmEventDouble = 32
	gpmEventTriple = 64
	gpmEventMFlag  = 128
	gpmEventHard   = 256
	gpmEventEnter  = 512
	gpmEventLeave  = 1024

	gpmMaxInt = 1<<31 - 1
)

var (
	gpmCtlPath = "/dev/gpmctl"
	gpmFDLink  = func(fd string) (string, error) {
		return os.Readlink(filepath.Join("/proc/self/fd", fd))
	}
)

type gpmConsoleMouse struct {
	sock   net.Conn
	eventQ chan<- Event

	mu      sync.Mutex
	flags   MouseFlags
	pressed ButtonMask
	w       int
	h       int

	stopOnce sync.Once
	stopQ    chan struct{}
	done     chan struct{}
}

func init() {
	newConsoleMouse = newGPMConsoleMouse
}

func newGPMConsoleMouse(eventQ chan<- Event) (consoleMouse, error) {
	tty, err := gpmConsoleTTY()
	if err != nil {
		return nil, err
	}
	sock, err := net.Dial("unix", gpmCtlPath)
	if err != nil {
		return nil, fmt.Errorf("failed opening gpmctl: %w", err)
	}

	cm := &gpmConsoleMouse{sock: sock, eventQ: eventQ}
	if _, err = sock.Write(gpmConnectPacket(tty)); err != nil {
		_ = sock.Close()
		return nil, fmt.Errorf("failed sending gpmctl connect: %w", err)
	}
	return cm, nil
}

func gpmConsoleTTY() (uint32, error) {
	var last error
	for _, fd := range []string{"2", "0", "1"} {
		link, err := gpmFDLink(fd)
		if err != nil {
			last = err
			continue
		}
		if tty, ok := parseGPMConsoleTTY(link); ok {
			return tty, nil
		}
		last = fmt.Errorf("not a virtual console: %s", link)
	}
	if last == nil {
		last = errNoConsoleMouse
	}
	return 0, fmt.Errorf("cannot determine Linux virtual console: %w", last)
}

func parseGPMConsoleTTY(link string) (uint32, bool) {
	link = filepath.Clean(link)
	base := filepath.Base(link)
	dir := filepath.Dir(link)
	var s string
	switch {
	case dir == "/dev" && strings.HasPrefix(base, "tty"):
		s = strings.TrimPrefix(base, "tty")
	case filepath.Base(dir) == "vc" && filepath.Dir(dir) == "/dev":
		s = base
	default:
		return 0, false
	}
	tty, err := strconv.ParseUint(s, 10, 32)
	if err != nil || tty == 0 || tty > gpmMaxInt {
		return 0, false
	}
	return uint32(tty), true
}

func gpmConnectPacket(tty uint32) []byte {
	b := make([]byte, gpmConnectSize)
	ne := binary.NativeEndian
	ne.PutUint16(b[0:], 0xffff)         // eventMask
	ne.PutUint16(b[2:], 0xffff)         // defaultMask
	ne.PutUint16(b[4:], 0)              // minMod
	ne.PutUint16(b[6:], 0xffff)         // maxMod
	ne.PutUint32(b[8:], gpmProcessID()) // pid
	ne.PutUint32(b[12:], tty)           // vc
	return b
}

func gpmProcessID() uint32 {
	pid := os.Getpid()
	if pid < 0 || pid > gpmMaxInt {
		return 0
	}
	return uint32(pid)
}

func (cm *gpmConsoleMouse) Start(flags MouseFlags) error {
	if flags&(MouseButtonEvents|MouseDragEvents|MouseMotionEvents) == 0 {
		return errNoConsoleMouse
	}
	cm.mu.Lock()
	cm.flags = flags
	cm.stopQ = make(chan struct{})
	cm.done = make(chan struct{})
	cm.mu.Unlock()

	go cm.loop()
	return nil
}

func (cm *gpmConsoleMouse) Stop() {
	cm.mu.Lock()
	stopQ := cm.stopQ
	done := cm.done
	cm.mu.Unlock()
	cm.stopOnce.Do(func() {
		if stopQ != nil {
			close(stopQ)
		}
		_ = cm.sock.Close()
	})
	if done != nil {
		<-done
	}
}

func (cm *gpmConsoleMouse) SetSize(w, h int) {
	cm.mu.Lock()
	cm.w = w
	cm.h = h
	cm.mu.Unlock()
}

func (cm *gpmConsoleMouse) loop() {
	defer func() {
		cm.mu.Lock()
		done := cm.done
		cm.mu.Unlock()
		if done != nil {
			close(done)
		}
	}()

	cm.mu.Lock()
	stopQ := cm.stopQ
	cm.mu.Unlock()

	var b [gpmEventSize]byte
	for {
		if _, err := io.ReadFull(cm.sock, b[:]); err != nil {
			return
		}
		if ev, ok := cm.eventFromPacket(b[:]); ok {
			select {
			case cm.eventQ <- ev:
			case <-stopQ:
				return
			}
		}
	}
}

func (cm *gpmConsoleMouse) eventFromPacket(b []byte) (*EventMouse, bool) {
	if len(b) < gpmEventSize {
		return nil, false
	}

	ne := binary.NativeEndian
	buttons := b[0]
	mods := b[1]
	x := int(int16(ne.Uint16(b[8:]))) - 1
	y := int(int16(ne.Uint16(b[10:]))) - 1
	typ := int(ne.Uint32(b[12:]))
	wdx := int(int16(ne.Uint16(b[24:])))
	wdy := int(int16(ne.Uint16(b[26:])))

	cm.mu.Lock()
	defer cm.mu.Unlock()

	button := gpmButtonMask(buttons, wdx, wdy)
	bare := typ & (gpmEventMove | gpmEventDrag | gpmEventDown | gpmEventUp)
	flags := cm.flags

	if button&(WheelUp|WheelDown|WheelLeft|WheelRight) != 0 {
		if flags&(MouseButtonEvents|MouseDragEvents|MouseMotionEvents) == 0 {
			return nil, false
		}
		return NewEventMouse(cm.clipX(x), cm.clipY(y), button, gpmModMask(mods)), true
	}

	switch bare {
	case gpmEventMove:
		if flags&MouseMotionEvents == 0 {
			return nil, false
		}
		cm.pressed = ButtonNone
	case gpmEventDrag:
		if flags&(MouseDragEvents|MouseMotionEvents) == 0 {
			return nil, false
		}
		if button != ButtonNone {
			cm.pressed |= button
		}
	case gpmEventDown:
		if flags&(MouseButtonEvents|MouseDragEvents|MouseMotionEvents) == 0 {
			return nil, false
		}
		cm.pressed |= button
	case gpmEventUp:
		if flags&(MouseButtonEvents|MouseDragEvents|MouseMotionEvents) == 0 {
			return nil, false
		}
		cm.pressed &^= button
	default:
		return nil, false
	}

	return NewEventMouse(cm.clipX(x), cm.clipY(y), cm.pressed, gpmModMask(mods)), true
}

func (cm *gpmConsoleMouse) clipX(x int) int {
	if cm.w <= 0 {
		return max(x, 0)
	}
	return max(min(x, cm.w-1), 0)
}

func (cm *gpmConsoleMouse) clipY(y int) int {
	if cm.h <= 0 {
		return max(y, 0)
	}
	return max(min(y, cm.h-1), 0)
}

func gpmButtonMask(buttons byte, wdx, wdy int) ButtonMask {
	var mask ButtonMask
	if buttons&gpmButtonLeft != 0 {
		mask |= Button1
	}
	if buttons&gpmButtonRight != 0 {
		mask |= Button2
	}
	if buttons&gpmButtonMiddle != 0 {
		mask |= Button3
	}
	if buttons&gpmButtonFourth != 0 {
		mask |= Button4
	}
	if buttons&gpmButtonUp != 0 || wdy > 0 {
		mask |= WheelUp
	}
	if buttons&gpmButtonDown != 0 || wdy < 0 {
		mask |= WheelDown
	}
	if wdx < 0 {
		mask |= WheelLeft
	}
	if wdx > 0 {
		mask |= WheelRight
	}
	return mask
}

func gpmModMask(mods byte) ModMask {
	var mask ModMask
	if mods&gpmModShift != 0 {
		mask |= ModShift
	}
	if mods&gpmModAltGr != 0 {
		mask |= ModAlt
	}
	if mods&gpmModCtrl != 0 {
		mask |= ModCtrl
	}
	if mods&gpmModAlt != 0 {
		mask |= ModAlt
	}
	return mask
}
