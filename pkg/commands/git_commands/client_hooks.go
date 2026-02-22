package git_commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
	"github.com/spf13/afero"
)

const TmpCommitEditMsg = "TMP_COMMIT_EDITMSG"

type ClientHook string

const (
	HookPreCommit        ClientHook = "pre-commit"
	HookPrepareCommitMsg ClientHook = "prepare-commit-msg"
	HookCommitMsg        ClientHook = "commit-msg"
	HookPreRebase        ClientHook = "pre-rebase"
	HookPostRewrite      ClientHook = "post-rewrite"
	HookPostMerge        ClientHook = "post-merge"
	HookPostCommit       ClientHook = "post-commit"
	HookPrePush          ClientHook = "pre-push"
	HookPostCheckout     ClientHook = "post-checkout"
)

type ClientHookCommands struct {
	gitCommon *GitCommon
}

func NewClientHookCommands(gitCommon *GitCommon) *ClientHookCommands {
	return &ClientHookCommands{
		gitCommon: gitCommon,
	}
}

func (self *ClientHookCommands) RunClientHook(hook ClientHook, args ...string) error {
	hookPath := self.resolveHookPath(hook)

	if _, err := self.gitCommon.Common.Fs.Stat(hookPath); err != nil {
		if errors.Is(err, afero.ErrFileNotFound) {
			return nil // Git silently ignores missing hooks
		}
		return err
	}

	argv := append([]string{hookPath}, args...)
	err := self.gitCommon.cmd.New(argv).Run()
	if err != nil {
		// Non-zero exit code from hook
		// Git surfaces error
		return fmt.Errorf("client hook %s failed: %w", hook, err)
	}

	return nil
}

func (self *ClientHookCommands) resolveHookPath(hook ClientHook) string {
	cmdArgs := NewGitCmd("config").
		Arg("--get", "core.hooksPath").
		ToArgv()
	out, err := self.gitCommon.cmd.New(cmdArgs).RunWithOutput()

	hooksPath := strings.TrimSpace(string(out))

	if err != nil || hooksPath == "" {
		// core.hooksPath not set or error, fallback to '.git/hooks'
		return filepath.Join(
			self.gitCommon.repoPaths.WorktreeGitDirPath(),
			"hooks",
			string(hook),
		)
	}

	if !filepath.IsAbs(hooksPath) {
		hooksPath = filepath.Join(self.gitCommon.repoPaths.repoPath, hooksPath)
	}

	return filepath.Join(hooksPath, string(hook))
}
