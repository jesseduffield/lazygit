package git_commands

import (
	"testing"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/samber/lo"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestStatusRefsSnapshot(t *testing.T) {
	const forEachRefOutput = "aaaa refs/heads/main\nbbbb refs/heads/topic\n"
	forEachRefArgs := []string{"for-each-ref", "--format=%(objectname) %(refname)", "refs/heads"}

	scenarios := []struct {
		testName     string
		headFile     *string // nil means: don't create a .git/HEAD file (simulates it being unreadable).
		runner       *oscommands.FakeCmdObjRunner
		expectedHead string
	}{
		{
			// files backend, on a branch: read straight from .git/HEAD, no
			// child process for HEAD.
			testName: "attached, read from HEAD file",
			headFile: lo.ToPtr("ref: refs/heads/main\n"),
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(forEachRefArgs, forEachRefOutput, nil),
			expectedHead: "ref: refs/heads/main",
		},
		{
			// files backend, detached: .git/HEAD holds the raw hash.
			testName: "detached, read from HEAD file",
			headFile: lo.ToPtr("aaaa\n"),
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(forEachRefArgs, forEachRefOutput, nil),
			expectedHead: "aaaa",
		},
		{
			// reftable backend (HEAD is a fixed stub), attached: fall back to
			// symbolic-ref, which succeeds.
			testName: "reftable stub, attached, fall back to symbolic-ref",
			headFile: lo.ToPtr("ref: refs/heads/.invalid\n"),
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(forEachRefArgs, forEachRefOutput, nil).
				ExpectGitArgs([]string{"symbolic-ref", "HEAD"}, "refs/heads/main\n", nil),
			expectedHead: "refs/heads/main",
		},
		{
			// reftable backend, detached: symbolic-ref fails, fall back to
			// rev-parse.
			testName: "reftable stub, detached, fall back to rev-parse",
			headFile: lo.ToPtr("ref: refs/heads/.invalid\n"),
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(forEachRefArgs, forEachRefOutput, nil).
				ExpectGitArgs([]string{"symbolic-ref", "HEAD"}, "", errors.New("fatal: ref HEAD is not a symbolic ref")).
				ExpectGitArgs([]string{"rev-parse", "HEAD"}, "aaaa\n", nil),
			expectedHead: "aaaa",
		},
		{
			// HEAD file missing/unreadable: same fallback as reftable.
			testName: "no HEAD file, fall back to symbolic-ref",
			headFile: nil,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs(forEachRefArgs, forEachRefOutput, nil).
				ExpectGitArgs([]string{"symbolic-ref", "HEAD"}, "refs/heads/main\n", nil),
			expectedHead: "refs/heads/main",
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			if s.headFile != nil {
				assert.NoError(t, afero.WriteFile(fs, "/repo/.git/HEAD", []byte(*s.headFile), 0o600))
			}

			instance := buildStatusCommands(commonDeps{
				runner:    s.runner,
				fs:        fs,
				repoPaths: MockRepoPaths("/repo"),
			})

			snapshot, err := instance.RefsSnapshot()
			assert.NoError(t, err)
			assert.Equal(t, forEachRefOutput+s.expectedHead, snapshot)
			s.runner.CheckForMissingCalls()
		})
	}
}
