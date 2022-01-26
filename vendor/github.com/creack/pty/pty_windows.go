//go:build windows
// +build windows

package pty

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	// NOTE(security): as noted by the comment of syscall.NewLazyDLL and syscall.LoadDLL
	// 	user need to call internal/syscall/windows/sysdll.Add("kernel32.dll") to make sure
	//  the kernel32.dll is loaded from windows system path
	//
	// ref: https://pkg.go.dev/syscall@go1.13?GOOS=windows#LoadDLL
	kernel32DLL = syscall.NewLazyDLL("kernel32.dll")

	// https://docs.microsoft.com/en-us/windows/console/createpseudoconsole
	createPseudoConsole = kernel32DLL.NewProc("CreatePseudoConsole")
	closePseudoConsole  = kernel32DLL.NewProc("ClosePseudoConsole")

	deleteProcThreadAttributeList     = kernel32DLL.NewProc("DeleteProcThreadAttributeList")
	initializeProcThreadAttributeList = kernel32DLL.NewProc("InitializeProcThreadAttributeList")

	// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-updateprocthreadattribute
	updateProcThreadAttribute = kernel32DLL.NewProc("UpdateProcThreadAttribute")

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

	defer func() {
		if err != nil {
			_ = consoleW.Close()
			_ = pr.Close()

			_ = pw.Close()
			_ = consoleR.Close()
		}
	}()

	err = createPseudoConsole.Find()
	if err != nil {
		return nil, nil, err
	}

	var consoleHandle syscall.Handle
	r1, _, err := createPseudoConsole.Call(
		(windowsCoord{X: 80, Y: 30}).Pack(),     // size: default 80x30 window
		consoleR.Fd(),                           // console input
		consoleW.Fd(),                           // console output
		0,                                       // console flags, currently only PSEUDOCONSOLE_INHERIT_CURSOR supported
		uintptr(unsafe.Pointer(&consoleHandle)), // console handler value return
	)
	if r1 != 0 {
		// S_OK: 0
		return nil, nil, os.NewSyscallError("CreatePseudoConsole", err)
	}

	return &WindowsPty{
			handle:   uintptr(consoleHandle),
			r:        pr,
			w:        pw,
			consoleR: consoleR,
			consoleW: consoleW,
		}, &WindowsTty{
			handle: uintptr(consoleHandle),
			r:      consoleR,
			w:      consoleW,
		}, nil
}

var _ Pty = (*WindowsPty)(nil)

type WindowsPty struct {
	handle uintptr
	r, w   *os.File

	consoleR, consoleW *os.File
}

func (p *WindowsPty) Fd() uintptr {
	return p.handle
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

func (p *WindowsPty) InputPipe() *os.File {
	return p.w
}

func (p *WindowsPty) OutputPipe() *os.File {
	return p.r
}

func (p *WindowsPty) Close() error {
	_ = p.r.Close()
	_ = p.w.Close()

	_ = p.consoleR.Close()
	_ = p.consoleW.Close()

	err := closePseudoConsole.Find()
	if err != nil {
		return err
	}

	_, _, err = closePseudoConsole.Call(p.handle)
	return err
}

var _ Tty = (*WindowsTty)(nil)

type WindowsTty struct {
	handle uintptr
	r, w   *os.File
}

func (t *WindowsTty) Fd() uintptr {
	return t.handle
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
