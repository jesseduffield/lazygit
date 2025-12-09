package git_commands

import (
	"path/filepath"
	"testing"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestRunClientHook(t *testing.T) {
	fs := afero.NewMemMapFs()
	repoDir := "/repo/.git"
	hooksDir := filepath.Join(repoDir, "hooks")
	_ = fs.MkdirAll(hooksDir, 0o755)

	runner := oscommands.NewFakeRunner(t)
	cmdBuilder := oscommands.NewDummyCmdObjBuilder(runner)

	gitCommon := &GitCommon{
		Common: &common.Common{
			Fs: fs,
		},
		cmd:       cmdBuilder,
		repoPaths: &RepoPaths{worktreeGitDirPath: repoDir},
	}

	hooks := NewHookCommands(gitCommon)

	t.Run("hook is executable", func(t *testing.T) {
		hookPath := filepath.Join(hooksDir, hooks.Hooks.PreCommit)
		_ = afero.WriteFile(fs, hookPath, []byte("#!/bin/sh"), 0o755)
		runner.ExpectArgs([]string{hookPath}, "", nil)
		err := hooks.RunClientHook(hooks.Hooks.PreCommit)
		assert.NoError(t, err)
	})

	t.Run("hook exists but fails", func(t *testing.T) {
		hookPath := filepath.Join(hooksDir, hooks.Hooks.PrepareCommitMsg)
		_ = afero.WriteFile(fs, hookPath, []byte("#!/bin/sh"), 0o755)
		runner.ExpectArgs([]string{hookPath}, "", errors.New("boom"))
		err := hooks.RunClientHook(hooks.Hooks.PrepareCommitMsg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "boom")
	})

	t.Run("hook exists but is not executable", func(t *testing.T) {
		hookPath := filepath.Join(hooksDir, hooks.Hooks.CommitMsg)
		_ = afero.WriteFile(fs, hookPath, []byte("#!/bin/sh"), 0o644)
		err := hooks.RunClientHook(hooks.Hooks.CommitMsg)
		assert.NoError(t, err)
	})

	t.Run("hook does not exist", func(t *testing.T) {
		err := hooks.RunClientHook("missing-hook")
		assert.NoError(t, err)
	})
}
