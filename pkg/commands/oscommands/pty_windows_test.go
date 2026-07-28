package oscommands

import (
	"os/exec"
	"testing"

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
