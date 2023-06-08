//go:build windows
// +build windows

package pty

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
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
	// NOTE(security): as noted by the comment of syscall.NewLazyDLL and syscall.LoadDLL
	// 	user need to call internal/syscall/windows/sysdll.Add("kernel32.dll") to make sure
	//  the kernel32.dll is loaded from windows system path
	//
	// ref: https://pkg.go.dev/syscall@go1.13?GOOS=windows#LoadDLL
	kernel32DLL = windows.NewLazyDLL("kernel32.dll")

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
		_ = consoleW.Close()
		_ = pr.Close()
		return nil, nil, err
	}

	var consoleHandle windows.Handle

	err = procCreatePseudoConsole(windows.Handle(consoleR.Fd()), windows.Handle(consoleW.Fd()),
		0, &consoleHandle)
	if err != nil {
		_ = consoleW.Close()
		_ = pr.Close()
		_ = pw.Close()
		_ = consoleR.Close()
		return nil, nil, err
	}

	// These pipes can be closed here without any worry
	err = consoleW.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to close pseudo console handle: %w", err)
	}

	err = consoleR.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to close pseudo console handle: %w", err)
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
	var err error

	if err = attrList.Update(
		PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE,
		unsafe.Pointer(p.handle),
		unsafe.Sizeof(p.handle),
	); err != nil {
		return fmt.Errorf("failed to update proc thread attributes for pseudo console: %w", err)
	}

	return nil
}

func (p *WindowsPty) Close() error {
	_ = p.r.Close()
	_ = p.w.Close()

	err := closePseudoConsole.Find()
	if err != nil {
		return err
	}

	_, _, err = closePseudoConsole.Call(uintptr(p.handle))
	return err
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
	_ = t.r.Close()
	return t.w.Close()
}

func procCreatePseudoConsole(hInput windows.Handle, hOutput windows.Handle, dwFlags uint32, consoleHandle *windows.Handle) error {
	var r0 uintptr
	var err error

	err = createPseudoConsole.Find()
	if err != nil {
		return err
	}

	r0, _, err = createPseudoConsole.Call(
		(windowsCoord{X: 80, Y: 30}).Pack(),    // size: default 80x30 window
		uintptr(hInput),                        // console input
		uintptr(hOutput),                       // console output
		uintptr(dwFlags),                       // console flags, currently only PSEUDOCONSOLE_INHERIT_CURSOR supported
		uintptr(unsafe.Pointer(consoleHandle)), // console handler value return
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
