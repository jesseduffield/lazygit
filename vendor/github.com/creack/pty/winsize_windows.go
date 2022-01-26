package pty

import (
	"os"
	"unsafe"
)

// types from golang.org/x/sys/windows
type (
	// copy of https://pkg.go.dev/golang.org/x/sys/windows#Coord
	windowsCoord struct {
		X int16
		Y int16
	}

	// copy of https://pkg.go.dev/golang.org/x/sys/windows#SmallRect
	windowsSmallRect struct {
		Left   int16
		Top    int16
		Right  int16
		Bottom int16
	}

	// copy of https://pkg.go.dev/golang.org/x/sys/windows#ConsoleScreenBufferInfo
	windowsConsoleScreenBufferInfo struct {
		Size              windowsCoord
		CursorPosition    windowsCoord
		Attributes        uint16
		Window            windowsSmallRect
		MaximumWindowSize windowsCoord
	}
)

func (c windowsCoord) Pack() uintptr {
	return uintptr((int32(c.Y) << 16) | int32(c.X))
}

// Setsize resizes t to ws.
func Setsize(t FdHolder, ws *Winsize) error {
	err := resizePseudoConsole.Find()
	if err != nil {
		return err
	}

	r1, _, err := resizePseudoConsole.Call(
		t.Fd(),
		(windowsCoord{X: int16(ws.Cols), Y: int16(ws.Rows)}).Pack(),
	)
	if r1 != 0 {
		// S_OK: 0
		return os.NewSyscallError("ResizePseudoConsole", err)
	}

	return nil
}

// GetsizeFull returns the full terminal size description.
func GetsizeFull(t FdHolder) (size *Winsize, err error) {
	err = getConsoleScreenBufferInfo.Find()
	if err != nil {
		return nil, err
	}

	var info windowsConsoleScreenBufferInfo
	r1, _, err := getConsoleScreenBufferInfo.Call(t.Fd(), uintptr(unsafe.Pointer(&info)))
	if r1 != 0 {
		// S_OK: 0
		return nil, os.NewSyscallError("GetConsoleScreenBufferInfo", err)
	}

	return &Winsize{
		Rows: uint16(info.Window.Bottom - info.Window.Top + 1),
		Cols: uint16(info.Window.Right - info.Window.Left + 1),
	}, nil
}
