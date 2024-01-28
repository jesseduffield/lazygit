//go:build windows
// +build windows

package pty

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	// Ref: https://pkg.go.dev/golang.org/x/sys/windows#pkg-constants
	PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE = 0x20016
)

type WindowsPty struct {
	handle windows.Handle
	r, w   *os.File
}

type WindowsTty struct {
	handle windows.Handle
	r, w   *os.File
}

var (
	// NOTE(security): As noted by the comment of syscall.NewLazyDLL and syscall.LoadDLL
	//                 user need to call internal/syscall/windows/sysdll.Add("kernel32.dll") to make sure
	//                 the kernel32.dll is loaded from windows system path.
	//
	// Ref: https://pkg.go.dev/syscall@go1.13?GOOS=windows#LoadDLL
	kernel32DLL = windows.NewLazySystemDLL("kernel32.dll")

	// https://docs.microsoft.com/en-us/windows/console/createpseudoconsole
	createPseudoConsole = kernel32DLL.NewProc("CreatePseudoConsole")
	closePseudoConsole  = kernel32DLL.NewProc("ClosePseudoConsole")

	resizePseudoConsole        = kernel32DLL.NewProc("ResizePseudoConsole")
	getConsoleScreenBufferInfo = kernel32DLL.NewProc("GetConsoleScreenBufferInfo")
)

func open() (_ Pty, _ Tty, err error) {
	pr, consoleW, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	consoleR, pw, err := os.Pipe()
	if err != nil {
		// Closing everything. Best effort.
		_ = consoleW.Close()
		_ = pr.Close()
		return nil, nil, err
	}

	var consoleHandle windows.Handle

	// TODO: As we removed the use of `.Fd()` on Unix (https://github.com/creack/pty/pull/168), we need to check if we should do the same here.
	if err := procCreatePseudoConsole(
		windows.Handle(consoleR.Fd()),
		windows.Handle(consoleW.Fd()),
		0,
		&consoleHandle); err != nil {
		// Closing everything. Best effort.
		_ = consoleW.Close()
		_ = pr.Close()
		_ = pw.Close()
		_ = consoleR.Close()
		return nil, nil, err
	}

	// These pipes can be closed here without any worry.
	if err := consoleW.Close(); err != nil {
		return nil, nil, fmt.Errorf("failed to close pseudo console write handle: %w", err)
	}
	if err := consoleR.Close(); err != nil {
		return nil, nil, fmt.Errorf("failed to close pseudo console read handle: %w", err)
	}

	return &WindowsPty{
			handle: consoleHandle,
			r:      pr,
			w:      pw,
		}, &WindowsTty{
			handle: consoleHandle,
			r:      consoleR,
			w:      consoleW,
		}, nil
}

func (p *WindowsPty) Name() string {
	return p.r.Name()
}

func (p *WindowsPty) Fd() uintptr {
	return uintptr(p.handle)
}

func (p *WindowsPty) Read(data []byte) (int, error) {
	return p.r.Read(data)
}

func (p *WindowsPty) Write(data []byte) (int, error) {
	return p.w.Write(data)
}

func (p *WindowsPty) WriteString(s string) (int, error) {
	return p.w.WriteString(s)
}

func (p *WindowsPty) UpdateProcThreadAttribute(attrList *windows.ProcThreadAttributeListContainer) error {
	if err := attrList.Update(
		PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE,
		unsafe.Pointer(p.handle),
		unsafe.Sizeof(p.handle),
	); err != nil {
		return fmt.Errorf("failed to update proc thread attributes for pseudo console: %w", err)
	}

	return nil
}

func (p *WindowsPty) Close() error {
	// Best effort.
	_ = p.r.Close()
	_ = p.w.Close()

	if err := closePseudoConsole.Find(); err != nil {
		return err
	}

	_, _, err := closePseudoConsole.Call(uintptr(p.handle))
	return err
}

func (p *WindowsPty) SetDeadline(value time.Time) error {
	return os.ErrNoDeadline
}

func (t *WindowsTty) Name() string {
	return t.r.Name()
}

func (t *WindowsTty) Fd() uintptr {
	return uintptr(t.handle)
}

func (t *WindowsTty) Read(p []byte) (int, error) {
	return t.r.Read(p)
}

func (t *WindowsTty) Write(p []byte) (int, error) {
	return t.w.Write(p)
}

func (t *WindowsTty) Close() error {
	_ = t.r.Close() // Best effort.
	return t.w.Close()
}

func (t *WindowsTty) SetDeadline(value time.Time) error {
	return os.ErrNoDeadline
}

func procCreatePseudoConsole(hInput windows.Handle, hOutput windows.Handle, dwFlags uint32, consoleHandle *windows.Handle) error {
	if err := createPseudoConsole.Find(); err != nil {
		return err
	}

	// TODO: Check if it is expected to ignore `err` here.
	r0, _, _ := createPseudoConsole.Call(
		(windowsCoord{X: 80, Y: 30}).Pack(),    // Size: default 80x30 window.
		uintptr(hInput),                        // Console input.
		uintptr(hOutput),                       // Console output.
		uintptr(dwFlags),                       // Console flags, currently only PSEUDOCONSOLE_INHERIT_CURSOR supported.
		uintptr(unsafe.Pointer(consoleHandle)), // Console handler value return.
	)

	if int32(r0) < 0 {
		if r0&0x1fff0000 == 0x00070000 {
			r0 &= 0xffff
		}

		// S_OK: 0
		return syscall.Errno(r0)
	}

	return nil
}
