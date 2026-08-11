package oscommands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/jesseduffield/lazygit/pkg/utils"
	"golang.org/x/sys/windows"
)

type winPty struct {
	hpc windows.Handle
	// job holds the child and every descendant it spawns; terminating it
	// kills whatever is left of the process tree (see Close).
	job windows.Handle
	// conhost is a handle to the conhost.exe serving this pty, or 0 if it
	// couldn't be identified. Held so that the teardown in Close can reap
	// it on Windows builds whose conhost fails to run down on its own.
	conhost windows.Handle
	inWrite *os.File
	outRead *os.File

	// mu guards hpcClosed, which gates ClosePseudoConsole (it must run
	// exactly once) and also keeps Resize from touching the HPCON once it's
	// been freed: the background waiter in StartPty closes the pseudoconsole
	// on child exit, which would otherwise race a concurrent onResize and
	// hand ResizePseudoConsole a freed handle.
	mu        sync.Mutex
	hpcClosed bool
}

func (p *winPty) Read(buf []byte) (int, error)  { return p.outRead.Read(buf) }
func (p *winPty) Write(buf []byte) (int, error) { return p.inWrite.Write(buf) }

func (p *winPty) Resize(cols, rows uint16) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.hpcClosed {
		// The child already exited and the pseudoconsole was torn down, so
		// there is nothing left to resize.
		return nil
	}
	return windows.ResizePseudoConsole(p.hpc, clampPtySize(cols, rows))
}

// clampPtySize clamps a requested pty size to the minimum that ConPTY
// accepts: CreatePseudoConsole and ResizePseudoConsole reject zero
// dimensions with E_INVALIDARG, but callers legitimately request them — the
// pty is sized after the main view, which is zero-sized while hidden, e.g.
// in full-screen mode with a side panel focused.
func clampPtySize(cols, rows uint16) windows.Coord {
	return windows.Coord{X: int16(max(cols, 1)), Y: int16(max(rows, 1))}
}

// closeHpc closes the pseudoconsole exactly once. Safe to call from multiple
// goroutines and at any time. We need this separately from Close because the
// background waiter in StartPty closes the pseudoconsole as soon as the child
// exits — that's what makes outRead return EOF, matching the Unix behavior
// where the master fd EOFs when the slave closes — while the pipe fds stay
// open until somebody explicitly tears the pty down.
func (p *winPty) closeHpc() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.hpcClosed {
		return
	}
	p.hpcClosed = true
	windows.ClosePseudoConsole(p.hpc)
}

// How long Close waits for the conhost to run itself down after its clients
// are gone, before concluding that it never will (see Close) and reaping it.
const conhostExitTimeout = time.Second

var (
	// ptyTeardowns counts the in-flight teardown goroutines spawned by
	// Close; TerminateLivePtys waits for them when lazygit exits.
	ptyTeardowns sync.WaitGroup
	// ptyQuit is closed by TerminateLivePtys. In-flight teardowns skip the
	// conhost rundown wait once it is closed: the conhost serves nothing
	// once its clients are gone, and the exit must not stall for its sake.
	ptyQuit     = make(chan struct{})
	ptyQuitOnce sync.Once
)

// TerminateLivePtys synchronously terminates the process trees and console
// hosts of all ptys whose teardown hasn't finished yet. Call it when lazygit
// is about to exit: the asynchronous teardowns in Close won't get to finish
// (the conhost rundown wait outlives the process), and while
// KILL_ON_JOB_CLOSE reaps the clients when the job handles are closed at
// process death, nothing would reap the conhosts on the Windows builds that
// need it (see Close). A long diff on screen keeps its git process running
// the whole time it is shown, so quitting with such a teardown in flight is
// the rule, not the exception.
func TerminateLivePtys() {
	ptyQuitOnce.Do(func() { close(ptyQuit) })

	done := make(chan struct{})
	go utils.Safe(func() {
		ptyTeardowns.Wait()
		close(done)
	})
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		// Don't hold up the exit any longer; the job handles' rundown
		// still covers the clients.
	}
}

// Close tears the pty down without waiting for it: the teardown runs on a
// background goroutine and Close returns immediately.
//
// It has to, because ClosePseudoConsole can block for a long time: before
// Windows 11 24H2 it waits for the console host to exit, and since closing
// only delivers CTRL_CLOSE_EVENT to the attached client without terminating
// it, a client that keeps running (git still computing an expensive diff, a
// diff renderer waiting for input) keeps the host — and with it
// ClosePseudoConsole — alive arbitrarily long. Close is called while holding
// the global PtyMutex and while the task's onDone once is executing, where
// blocking wedges every subsequent task for the view (and with it the UI), so
// none of this may happen on the caller's thread.
//
// Within the teardown, the pipe ends must be closed before the
// pseudoconsole, and without holding p.mu: closing the pseudoconsole flushes
// the client's pending output into the out pipe, and with the task stopped
// nobody is reading anymore, so that flush can only complete once the pipe
// is broken. The background waiter's closeHpc may already be wedged in such
// a flush while holding p.mu; closing the pipes is what unblocks it.
//
// Closing the pseudoconsole delivers CTRL_CLOSE_EVENT only to the clients
// attached to it at that moment. A child that is stopped right after being
// spawned is still starting up and not attached yet, so the event misses it
// and it survives, running its command to completion as an orphan — and
// keeping its console host alive with it (#5879); the same holds for
// grandchildren spawned while the console is going down, and for clients
// that ignore the event (the Windows flavor of #5675). The job kill reaps
// all of those. There is no point in delaying it: the close event is not a
// graceful signal worth waiting on — git and the common diff tools leave it
// to the default handler, which calls ExitProcess at whatever instruction
// the process happens to execute — so clients that got the event are
// already dying. Killing at an arbitrary point cannot leak a stale
// index.lock, because pty-rendered commands don't take that lock (see
// withPtyGitConfig in pkg/gui/pty.go).
//
// The pseudoconsole close gets its own goroutine because the kill must not
// wait for it: on builds where ClosePseudoConsole blocks until the console
// host exits (pre-24H2), the host keeps running as long as a surviving
// client does, and that client only goes away through the job kill —
// sequencing the kill after a blocking close would thus deadlock in
// exactly the case the kill exists for.
//
// After the kill, the conhost serving the pty is reaped as well if it
// doesn't exit by itself: a healthy conhost runs down once the reference
// handle is closed and its clients are gone, but conhost builds before
// Windows 11 24H2 fail to complete the rundown when a client attached
// after the close event was delivered and was then killed — the fate of
// exactly the clients the job kill is for — and such a conhost sits
// around forever, serving nothing (#5879). The reap is inert on healthy
// builds: the wait succeeds and only the handle is closed.
//
// When lazygit is quitting, the conhost rundown wait is skipped; see
// TerminateLivePtys.
func (p *winPty) Close() error {
	ptyTeardowns.Add(1)
	go utils.Safe(func() {
		defer ptyTeardowns.Done()

		p.inWrite.Close()
		p.outRead.Close()
		go utils.Safe(p.closeHpc)

		_ = windows.TerminateJobObject(p.job, 1)
		_ = windows.CloseHandle(p.job)

		if p.conhost != 0 {
			timeout := conhostExitTimeout
			select {
			case <-ptyQuit:
				timeout = 0
			default:
			}
			event, err := windows.WaitForSingleObject(p.conhost, uint32(timeout/time.Millisecond))
			if err != nil || event != windows.WAIT_OBJECT_0 {
				_ = windows.TerminateProcess(p.conhost, 1)
			}
			_ = windows.CloseHandle(p.conhost)
		}
	})
	return nil
}

// startWaiter runs proc.Wait in a goroutine and, as soon as the child exits,
// closes the pseudoconsole so that any pending Read on outRead returns EOF
// after buffered output drains. Returns a Wait func that blocks until the
// child has exited and reports its exit status with *exec.Cmd.Wait semantics.
//
// This shape exists because on Unix the master fd EOFs naturally when the
// slave closes on child exit, but ConPTY keeps the pipe alive until we call
// ClosePseudoConsole explicitly. Without doing that on child exit, the
// scanner in pkg/tasks.NewCmdTask would block forever on the next read and
// the post-content view never gets cleared (FlushStaleCells never fires).
func startWaiter(proc *os.Process, p *winPty) func() error {
	done := make(chan struct{})
	var waitErr error
	go func() {
		defer close(done)
		state, err := proc.Wait()
		p.closeHpc()
		if err != nil {
			waitErr = err
			return
		}
		if !state.Success() {
			waitErr = fmt.Errorf("exit status %d", state.ExitCode())
		}
	}()
	return func() error {
		<-done
		return waitErr
	}
}

// conhostScanMu serializes CreatePseudoConsole and the child-process scans
// around it, so that two concurrently starting ptys can't make each other's
// "which conhost is new" diff ambiguous.
var conhostScanMu sync.Mutex

// conhostChildren returns the pids of all conhost.exe processes that are
// direct children of this process. Errors just yield a smaller (possibly
// empty) set; the caller treats identification as best-effort.
func conhostChildren() map[uint32]bool {
	pids := map[uint32]bool{}
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return pids
	}
	defer func() { _ = windows.CloseHandle(snap) }()
	me := uint32(os.Getpid())
	var pe windows.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))
	for err := windows.Process32First(snap, &pe); err == nil; err = windows.Process32Next(snap, &pe) {
		if pe.ParentProcessID == me && strings.EqualFold(windows.UTF16ToString(pe.ExeFile[:]), "conhost.exe") {
			pids[pe.ProcessID] = true
		}
	}
	return pids
}

// openNewConhostChild returns a handle to the single conhost child that
// appeared since the before scan, or 0 if there isn't exactly one candidate
// or it can't be opened.
func openNewConhostChild(before map[uint32]bool) windows.Handle {
	var found []uint32
	for pid := range conhostChildren() {
		if !before[pid] {
			found = append(found, pid)
		}
	}
	if len(found) != 1 {
		return 0
	}
	h, err := windows.OpenProcess(windows.SYNCHRONIZE|windows.PROCESS_TERMINATE, false, found[0])
	if err != nil {
		return 0
	}
	return h
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
			_ = windows.CloseHandle(inWrite)
		}
	}()
	if err = windows.CreatePipe(&outRead, &outWrite, nil, 0); err != nil {
		_ = windows.CloseHandle(inRead)
		return StartedPty{}, fmt.Errorf("CreatePipe (out): %w", err)
	}
	defer func() {
		if err != nil {
			_ = windows.CloseHandle(outRead)
		}
	}()

	// CreatePseudoConsole dupes the handles it needs internally; we release
	// our references to the child-side ends immediately after.
	//
	// It also spawns the conhost.exe serving the console session, as a
	// direct child of this process. The teardown in Close needs a handle to
	// that conhost (see there), but Windows offers no way to obtain one
	// from the HPCON, so identify it by diffing our conhost children around
	// the call. Open a real handle right away so that pid reuse can't later
	// misdirect the teardown's reap. If identification fails, the handle
	// stays 0 and the teardown skips the reap.
	var hpc, conhost windows.Handle
	size := clampPtySize(cols, rows)
	conhostScanMu.Lock()
	conhostsBefore := conhostChildren()
	err = windows.CreatePseudoConsole(size, inRead, outWrite, 0, &hpc)
	if err == nil {
		conhost = openNewConhostChild(conhostsBefore)
	}
	conhostScanMu.Unlock()
	if err != nil {
		_ = windows.CloseHandle(inRead)
		_ = windows.CloseHandle(outWrite)
		return StartedPty{}, fmt.Errorf("CreatePseudoConsole: %w", err)
	}
	_ = windows.CloseHandle(inRead)
	_ = windows.CloseHandle(outWrite)
	defer func() {
		if err != nil {
			windows.ClosePseudoConsole(hpc)
			if conhost != 0 {
				_ = windows.CloseHandle(conhost)
			}
		}
	}()

	// The child goes into a job object so that the teardown in Close can
	// terminate the whole process tree. KILL_ON_JOB_CLOSE makes the OS do
	// that when the last handle to the job is closed, which doubles as a
	// safety net: if lazygit exits without running the teardown, the handle
	// is closed for it and the tree is reaped.
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return StartedPty{}, fmt.Errorf("CreateJobObject: %w", err)
	}
	defer func() {
		if err != nil {
			// Kills the child on error paths where it was already assigned
			// to the job; plain handle cleanup before that.
			_ = windows.CloseHandle(job)
		}
	}()
	limits := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
			LimitFlags: windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
		},
	}
	if _, err = windows.SetInformationJobObject(
		job, windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&limits)), uint32(unsafe.Sizeof(limits)),
	); err != nil {
		return StartedPty{}, fmt.Errorf("SetInformationJobObject: %w", err)
	}

	// Attach the pseudoconsole to the child via a process attribute list.
	attrList, err := windows.NewProcThreadAttributeList(1)
	if err != nil {
		return StartedPty{}, fmt.Errorf("NewProcThreadAttributeList: %w", err)
	}
	defer attrList.Delete()
	// PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE wants the HPCON value itself as
	// the attribute value, not a pointer to it — an HPCON is already a
	// pointer-sized handle, per Microsoft's ConPTY sample. Spelling that as
	// unsafe.Pointer(hpc) trips go vet's unsafeptr check (a uintptr-based
	// type converted straight to unsafe.Pointer), which gopls surfaces in
	// the editor. Reinterpret the handle's bits through its address instead:
	// &hpc is a real pointer, so none of these conversions is the flagged
	// uintptr→unsafe.Pointer cast, while the resulting value is identical.
	if err = attrList.Update(
		windows.PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE,
		*(*unsafe.Pointer)(unsafe.Pointer(&hpc)),
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
		windows.EXTENDED_STARTUPINFO_PRESENT|windows.CREATE_UNICODE_ENVIRONMENT|windows.CREATE_SUSPENDED,
		envPtr,
		dirPtr,
		&si.StartupInfo,
		&pi,
	)
	if err != nil {
		return StartedPty{}, fmt.Errorf("CreateProcess: %w", err)
	}

	// The child was created suspended so that it can be assigned to the job
	// before it runs its first instruction; that way every descendant it
	// ever spawns is in the job from the start.
	if err = windows.AssignProcessToJobObject(job, pi.Process); err != nil {
		// Not in the job yet, so the deferred job-handle close can't reap it.
		_ = windows.TerminateProcess(pi.Process, 1)
		_ = windows.CloseHandle(pi.Thread)
		_ = windows.CloseHandle(pi.Process)
		return StartedPty{}, fmt.Errorf("AssignProcessToJobObject: %w", err)
	}
	if _, err = windows.ResumeThread(pi.Thread); err != nil {
		_ = windows.CloseHandle(pi.Thread)
		_ = windows.CloseHandle(pi.Process)
		return StartedPty{}, fmt.Errorf("ResumeThread: %w", err)
	}
	_ = windows.CloseHandle(pi.Thread)

	// Re-open the process by PID to get an *os.Process to wait on. Do this
	// while pi.Process is still open: Windows won't recycle a PID while any
	// handle to the process remains, so FindProcess can't latch onto a
	// different process that has since reused the PID. Release the original
	// handle once we have our own.
	proc, err := os.FindProcess(int(pi.ProcessId))
	_ = windows.CloseHandle(pi.Process)
	if err != nil {
		return StartedPty{}, err
	}

	wp := &winPty{
		hpc:     hpc,
		job:     job,
		conhost: conhost,
		inWrite: os.NewFile(uintptr(inWrite), "conpty-in"),
		outRead: os.NewFile(uintptr(outRead), "conpty-out"),
	}
	return StartedPty{
		Pty:     wp,
		Process: proc,
		Wait:    startWaiter(proc, wp),
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
