package git_commands

import (
	"os"
	"path/filepath"
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

func TestStatusWorkingTreeStateOfCherryPickOrRevert(t *testing.T) {
	scenarios := []struct {
		testName              string
		files                 map[string]string
		expectedCherryPicking bool
		expectedReverting     bool
	}{
		{
			testName: "nothing in progress",
			files:    map[string]string{},
		},
		{
			testName:              "CHERRY_PICK_HEAD",
			files:                 map[string]string{"CHERRY_PICK_HEAD": "1234567\n"},
			expectedCherryPicking: true,
		},
		{
			testName:          "REVERT_HEAD",
			files:             map[string]string{"REVERT_HEAD": "1234567\n"},
			expectedReverting: true,
		},
		{
			testName: "CHERRY_PICK_HEAD of a rebase that stopped because of conflicts",
			files: map[string]string{
				"CHERRY_PICK_HEAD":         "2345678abcdef\n",
				"rebase-merge/stopped-sha": "2345678\n",
				"rebase-merge/done":        "pick 1234567abcdef one\npick 2345678abcdef two\n",
			},
			expectedCherryPicking: false,
		},
		{
			testName: "CHERRY_PICK_HEAD of a rebase that couldn't commit a pick",
			files: map[string]string{
				"CHERRY_PICK_HEAD":  "2345678abcdef\n",
				"rebase-merge/done": "pick 1234567abcdef one\npick 2345678abcdef two\n",
			},
			expectedCherryPicking: false,
		},
		{
			testName: "cherry-pick while a rebase is stopped",
			files: map[string]string{
				"CHERRY_PICK_HEAD":  "3456789abcdef\n",
				"rebase-merge/done": "pick 1234567abcdef one\nedit 2345678abcdef two\n",
			},
			expectedCherryPicking: true,
		},
		{
			testName:              "cherry-pick of several commits that stopped before applying one",
			files:                 map[string]string{"sequencer/todo": "pick 1234567 one\npick 2345678 two\n"},
			expectedCherryPicking: true,
		},
		{
			testName:          "revert of several commits that stopped before applying one",
			files:             map[string]string{"sequencer/todo": "revert 1234567 one\nrevert 2345678 two\n"},
			expectedReverting: true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			repoDir := t.TempDir()
			for path, content := range s.files {
				fullPath := filepath.Join(repoDir, ".git", path)
				assert.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0o755))
				assert.NoError(t, os.WriteFile(fullPath, []byte(content), 0o644))
			}
			instance := buildStatusCommands(commonDeps{repoPaths: MockRepoPaths(repoDir)})

			state := instance.WorkingTreeState()

			assert.Equal(t, s.expectedCherryPicking, state.CherryPicking)
			assert.Equal(t, s.expectedReverting, state.Reverting)
		})
	}
}

func TestRebaseStoppedBeforeCommittingPick(t *testing.T) {
	scenarios := []struct {
		testName string
		files    map[string]string
		expected bool
	}{
		{
			testName: "the rebase couldn't commit a pick",
			files: map[string]string{
				"CHERRY_PICK_HEAD":  "2345678abcdef\n",
				"rebase-merge/done": "pick 1234567abcdef one\npick 2345678abcdef two\n",
			},
			expected: true,
		},
		{
			testName: "the rebase stopped because a pick conflicted",
			files: map[string]string{
				"CHERRY_PICK_HEAD":         "2345678abcdef\n",
				"rebase-merge/stopped-sha": "2345678\n",
				"rebase-merge/message":     "two\n",
				"rebase-merge/done":        "pick 1234567abcdef one\npick 2345678abcdef two\n",
			},
			expected: false,
		},
		{
			testName: "no CHERRY_PICK_HEAD",
			files: map[string]string{
				"rebase-merge/done": "pick 1234567abcdef one\nedit 2345678abcdef two\n",
			},
			expected: false,
		},
		{
			testName: "CHERRY_PICK_HEAD of a cherry-pick while a rebase is stopped",
			files: map[string]string{
				"CHERRY_PICK_HEAD":  "3456789abcdef\n",
				"rebase-merge/done": "pick 1234567abcdef one\nedit 2345678abcdef two\n",
			},
			expected: false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitDir := t.TempDir()
			for path, content := range s.files {
				fullPath := filepath.Join(gitDir, path)
				assert.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0o755))
				assert.NoError(t, os.WriteFile(fullPath, []byte(content), 0o644))
			}

			assert.Equal(t, s.expected, rebaseStoppedBeforeCommittingPick(gitDir, '#'))
		})
	}
}

func TestStatusWorkingTreeStateOfRebaseApplyDir(t *testing.T) {
	scenarios := []struct {
		testName                string
		files                   map[string]string
		expectedRebasing        bool
		expectedApplyingPatches bool
	}{
		{
			testName:                "git am",
			files:                   map[string]string{"rebase-apply/applying": ""},
			expectedRebasing:        false,
			expectedApplyingPatches: true,
		},
		{
			testName:                "rebase that uses the apply backend",
			files:                   map[string]string{"rebase-apply/rebasing": ""},
			expectedRebasing:        true,
			expectedApplyingPatches: false,
		},
		{
			testName: "git am while an interactive rebase is stopped",
			files: map[string]string{
				"rebase-apply/applying":        "",
				"rebase-merge/git-rebase-todo": "",
			},
			expectedRebasing:        true,
			expectedApplyingPatches: true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			repoDir := t.TempDir()
			for path, content := range s.files {
				fullPath := filepath.Join(repoDir, ".git", path)
				assert.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0o755))
				assert.NoError(t, os.WriteFile(fullPath, []byte(content), 0o644))
			}
			instance := buildStatusCommands(commonDeps{repoPaths: MockRepoPaths(repoDir)})

			state := instance.WorkingTreeState()

			assert.Equal(t, s.expectedRebasing, state.Rebasing)
			assert.Equal(t, s.expectedApplyingPatches, state.ApplyingPatches)
		})
	}
}
