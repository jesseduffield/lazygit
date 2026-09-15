package commands

import (
	"errors"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stretchr/testify/assert"
)

type runnerResult struct {
	output string
	err    error
}

// scriptedRunner is an ICmdObjRunner stub that returns a preconfigured result
// for each successive call, letting us drive the retry loop deterministically.
// It counts calls so tests can assert whether a command was retried.
type scriptedRunner struct {
	results []runnerResult
	calls   int
}

func (self *scriptedRunner) next() (string, error) {
	result := self.results[self.calls]
	self.calls++
	return result.output, result.err
}

func (self *scriptedRunner) Run(*oscommands.CmdObj) error {
	_, err := self.next()
	return err
}

func (self *scriptedRunner) RunWithOutput(*oscommands.CmdObj) (string, error) {
	return self.next()
}

func (self *scriptedRunner) RunWithOutputs(*oscommands.CmdObj) (string, string, error) {
	output, err := self.next()
	return output, "", err
}

func (self *scriptedRunner) RunAndProcessLines(*oscommands.CmdObj, func(string) (bool, error)) error {
	panic("not implemented")
}

func newTestRunner(inner *scriptedRunner) *gitCmdObjRunner {
	return &gitCmdObjRunner{
		log:         utils.NewDummyLog(),
		innerRunner: inner,
		// don't actually sleep between retries
		initialRetryDelay: 0,
	}
}

// dummyCmdObj returns a throwaway command; only its clonability matters, since
// the scriptedRunner ignores it and returns preconfigured results.
func dummyCmdObj() *oscommands.CmdObj {
	return oscommands.NewDummyCmdObjBuilder(nil).New([]string{"git", "status"})
}

func TestRunWithOutputReturnsSuccessWithoutRetrying(t *testing.T) {
	inner := &scriptedRunner{results: []runnerResult{{output: "done", err: nil}}}

	output, err := newTestRunner(inner).RunWithOutput(dummyCmdObj())

	assert.NoError(t, err)
	assert.Equal(t, "done", output)
	assert.Equal(t, 1, inner.calls)
}

func TestRunWithOutputDoesNotRetryNonLockError(t *testing.T) {
	inner := &scriptedRunner{results: []runnerResult{{output: "boom", err: errors.New("boom")}}}

	_, err := newTestRunner(inner).RunWithOutput(dummyCmdObj())

	assert.Error(t, err)
	assert.Equal(t, 1, inner.calls)
}

func TestRunWithOutputRetriesWhenLockErrorIsInOutput(t *testing.T) {
	inner := &scriptedRunner{results: []runnerResult{
		{output: "fatal: Unable to create '/repo/.git/index.lock': File exists.", err: errors.New("exit status 128")},
		{output: "done", err: nil},
	}}

	output, err := newTestRunner(inner).RunWithOutput(dummyCmdObj())

	assert.NoError(t, err)
	assert.Equal(t, "done", output)
	assert.Equal(t, 2, inner.calls)
}

func TestRunWithOutputRetriesWhenLockErrorIsOnlyInError(t *testing.T) {
	// A streamed command (e.g. an amend run through the gpg helper) doesn't
	// capture its output, so a lock failure surfaces only in the returned error
	// with an empty output string. The retry logic must still recognize it.
	inner := &scriptedRunner{results: []runnerResult{
		{output: "", err: errors.New("fatal: Unable to create '/repo/.git/index.lock': File exists.")},
		{output: "", err: nil},
	}}

	_, err := newTestRunner(inner).RunWithOutput(dummyCmdObj())

	assert.NoError(t, err)
	assert.Equal(t, 2, inner.calls)
}

func TestRunWithOutputGivesUpAfterMaxRetries(t *testing.T) {
	results := make([]runnerResult, maxRetries)
	for i := range results {
		results[i] = runnerResult{err: errors.New("fatal: Unable to create '/repo/.git/index.lock': File exists.")}
	}
	inner := &scriptedRunner{results: results}

	_, err := newTestRunner(inner).RunWithOutput(dummyCmdObj())

	assert.Error(t, err)
	assert.Equal(t, maxRetries, inner.calls)
}

func TestRunWithOutputRetriesLockErrorInLinkedWorktree(t *testing.T) {
	// In a linked worktree the lock lives at .git/worktrees/<name>/index.lock
	// rather than .git/index.lock, so only matching the bare "index.lock"
	// fragment lets the retry fire there too.
	inner := &scriptedRunner{results: []runnerResult{
		{output: "", err: errors.New("fatal: Unable to create '/repo/.git/worktrees/feature/index.lock': File exists.")},
		{output: "", err: nil},
	}}

	_, err := newTestRunner(inner).RunWithOutput(dummyCmdObj())

	assert.NoError(t, err)
	assert.Equal(t, 2, inner.calls)
}
