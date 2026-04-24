package oscommands

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

type winPty struct {
	hpc     windows.Handle
	inWrite *os.File
	outRead *os.File

	closeMu sync.Mutex
	closed  bool
}

func (p *winPty) Read(buf []byte) (int, error)  { return p.outRead.Read(buf) }
func (p *winPty) Write(buf []byte) (int, error) { return p.inWrite.Write(buf) }

func (p *winPty) Resize(cols, rows uint16) error {
	return windows.ResizePseudoConsole(p.hpc, windows.Coord{X: int16(cols), Y: int16(rows)})
}

func (p *winPty) Close() error {
	p.closeMu.Lock()
	defer p.closeMu.Unlock()
	if p.closed {
		return nil
	}
	p.closed = true
	// Closing the pseudoconsole breaks the pipes; the child's next write
	// fails and it exits. Then we close our ends of the pipes.
	windows.ClosePseudoConsole(p.hpc)
	p.inWrite.Close()
	p.outRead.Close()
	return nil
}

func waitForProcess(proc *os.Process) func() error {
	return func() error {
		state, err := proc.Wait()
		if err != nil {
			return err
		}
		if !state.Success() {
			return fmt.Errorf("exit status %d", state.ExitCode())
		}
		return nil
	}
}

func StartPty(cmd *exec.Cmd, cols, rows uint16) (sp StartedPty, err error) {
	// Two pipes: one for the child's stdin (we never write to it, but ConPTY
	// needs a handle), one for the child's stdout/stderr multiplexed through
	// the pseudoconsole.
	var inRead, inWrite, outRead, outWrite windows.Handle
	if err = windows.CreatePipe(&inRead, &inWrite, nil, 0); err != nil {
		return StartedPty{}, fmt.Errorf("CreatePipe (in): %w", err)
	}
	defer func() {
		if err != nil {
			windows.CloseHandle(inWrite)
		}
	}()
	if err = windows.CreatePipe(&outRead, &outWrite, nil, 0); err != nil {
		windows.CloseHandle(inRead)
		return StartedPty{}, fmt.Errorf("CreatePipe (out): %w", err)
	}
	defer func() {
		if err != nil {
			windows.CloseHandle(outRead)
		}
	}()

	// CreatePseudoConsole dupes the handles it needs internally; we release
	// our references to the child-side ends immediately after.
	var hpc windows.Handle
	size := windows.Coord{X: int16(cols), Y: int16(rows)}
	if err = windows.CreatePseudoConsole(size, inRead, outWrite, 0, &hpc); err != nil {
		windows.CloseHandle(inRead)
		windows.CloseHandle(outWrite)
		return StartedPty{}, fmt.Errorf("CreatePseudoConsole: %w", err)
	}
	windows.CloseHandle(inRead)
	windows.CloseHandle(outWrite)
	defer func() {
		if err != nil {
			windows.ClosePseudoConsole(hpc)
		}
	}()

	// Attach the pseudoconsole to the child via a process attribute list.
	attrList, err := windows.NewProcThreadAttributeList(1)
	if err != nil {
		return StartedPty{}, fmt.Errorf("NewProcThreadAttributeList: %w", err)
	}
	defer attrList.Delete()
	// PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE wants the HPCON value itself as
	// the attribute value, not a pointer to it — an HPCON is already a
	// pointer-sized handle. (go vet's unsafeptr check flags this, but it's
	// correct per Microsoft's ConPTY sample; the check isn't run by `go
	// test` anyway.)
	if err = attrList.Update(
		windows.PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE,
		unsafe.Pointer(hpc),
		unsafe.Sizeof(hpc),
	); err != nil {
		return StartedPty{}, fmt.Errorf("UpdateProcThreadAttribute: %w", err)
	}

	var si windows.StartupInfoEx
	si.Cb = uint32(unsafe.Sizeof(si))
	si.ProcThreadAttributeList = attrList.List()

	var appNamePtr *uint16
	if cmd.Path != "" {
		if appNamePtr, err = windows.UTF16PtrFromString(cmd.Path); err != nil {
			return StartedPty{}, err
		}
	}
	cmdLinePtr, err := windows.UTF16PtrFromString(windows.ComposeCommandLine(cmd.Args))
	if err != nil {
		return StartedPty{}, err
	}
	var dirPtr *uint16
	if cmd.Dir != "" {
		if dirPtr, err = windows.UTF16PtrFromString(cmd.Dir); err != nil {
			return StartedPty{}, err
		}
	}
	envBlock, err := createEnvBlock(cmd.Env)
	if err != nil {
		return StartedPty{}, err
	}
	var envPtr *uint16
	if envBlock != nil {
		envPtr = &envBlock[0]
	}

	var pi windows.ProcessInformation
	err = windows.CreateProcess(
		appNamePtr,
		cmdLinePtr,
		nil, // process security
		nil, // thread security
		false,
		windows.EXTENDED_STARTUPINFO_PRESENT|windows.CREATE_UNICODE_ENVIRONMENT,
		envPtr,
		dirPtr,
		&si.StartupInfo,
		&pi,
	)
	if err != nil {
		return StartedPty{}, fmt.Errorf("CreateProcess: %w", err)
	}
	windows.CloseHandle(pi.Thread)
	windows.CloseHandle(pi.Process)

	proc, err := os.FindProcess(int(pi.ProcessId))
	if err != nil {
		return StartedPty{}, err
	}

	return StartedPty{
		Pty: &winPty{
			hpc:     hpc,
			inWrite: os.NewFile(uintptr(inWrite), "conpty-in"),
			outRead: os.NewFile(uintptr(outRead), "conpty-out"),
		},
		Process: proc,
		Wait:    waitForProcess(proc),
	}, nil
}

// createEnvBlock packs env vars into the UTF-16 double-null-terminated block
// that CreateProcess expects. Returns nil if env is empty, which tells
// CreateProcess to inherit the parent's environment.
func createEnvBlock(env []string) ([]uint16, error) {
	if len(env) == 0 {
		return nil, nil
	}
	var block []uint16
	for _, s := range env {
		utf16s, err := windows.UTF16FromString(s)
		if err != nil {
			return nil, err
		}
		block = append(block, utf16s...)
	}
	block = append(block, 0)
	return block, nil
}
