package oscommands

import (
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// The requested size can legitimately be zero: the pty inherits the main
// view's dimensions, and that view is zero-sized while hidden, e.g. in
// full-screen mode with a side panel focused.
func TestStartPtyWithZeroSize(t *testing.T) {
	// The command deliberately produces no output: go test runs with
	// redirected std handles, which CreateProcess duplicates into the child
	// in place of handles to the attached pseudoconsole, so command output
	// would bypass the pty and pollute the test log.
	sp, err := StartPty(exec.Command("cmd", "/c", "exit 0"), 0, 0)
	assert.NoError(t, err)

	if err == nil {
		_ = sp.Wait()
		_ = sp.Pty.Close()
	}
}

// StartPty must identify the conhost.exe that CreatePseudoConsole spawned to
// serve the pty: the teardown in Close reaps it on Windows builds whose
// conhost fails to run down on its own, and a failed identification silently
// degrades to not reaping. If this fails, the child-scan in
// openNewConhostChild no longer matches how Windows hosts pseudoconsoles.
func TestStartPtyIdentifiesConhost(t *testing.T) {
	sp, err := StartPty(exec.Command("cmd", "/c", "exit 0"), 80, 24)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	assert.NotZero(t, sp.Pty.(*winPty).conhost)

	_ = sp.Wait()
	_ = sp.Pty.Close()
}

// TerminateLivePtys must reap a still-running pty synchronously: it runs
// when lazygit is about to exit, where the asynchronous teardown would not
// get to finish. Note that it switches the package's pty teardowns into
// quit mode for the remainder of the test binary's lifetime; that's fine
// for the other tests here, which must hold in either mode (quit mode only
// shortens the teardown's conhost rundown wait).
func TestTerminateLivePtysReapsRunningPty(t *testing.T) {
	// The output redirect is there for the reason described in
	// TestStartPtyWithZeroSize.
	sp, err := StartPty(exec.Command("cmd", "/c", "ping -n 30 127.0.0.1 >nul"), 80, 24)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	_ = sp.Pty.Close()
	TerminateLivePtys()

	// The teardown has completed as part of TerminateLivePtys, so the child
	// must be gone already; the timeout is generosity, not a grace period.
	exited := make(chan struct{})
	go func() {
		_ = sp.Wait()
		close(exited)
	}()
	select {
	case <-exited:
	case <-time.After(time.Second):
		t.Fatal("child process was not terminated by TerminateLivePtys")
	}
}

// Closing the pty must terminate the process tree it was running, even when
// it is closed so soon after starting that the child hasn't attached to the
// pseudoconsole yet: such a child misses the CTRL_CLOSE_EVENT that the close
// delivers to attached clients, and only the job-object kill reaps it.
// Without the kill, cmd and its ping child keep running for ~30 seconds and
// the Wait here times out.
func TestClosePtyTerminatesChildProcessTree(t *testing.T) {
	// The output redirect is there for the reason described in
	// TestStartPtyWithZeroSize.
	sp, err := StartPty(exec.Command("cmd", "/c", "ping -n 30 127.0.0.1 >nul"), 80, 24)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	_ = sp.Pty.Close()

	exited := make(chan struct{})
	go func() {
		_ = sp.Wait()
		close(exited)
	}()
	select {
	case <-exited:
	case <-time.After(5 * time.Second):
		t.Fatal("child process was not terminated by closing the pty")
	}
}
