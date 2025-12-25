package git_commands

import (
	"fmt"
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/spf13/afero"
)

type ClientHooks struct {
	PreCommit        string
	PrepareCommitMsg string
	CommitMsg        string
	PreRebase        string
	PostRewrite      string
	PostMerge        string
	PostCommit       string
	PrePush          string
	PostCheckout     string
}

func DefaultClientHooks() ClientHooks {
	return ClientHooks{
		PreCommit:        "pre-commit",
		PrepareCommitMsg: "prepare-commit-msg",
		CommitMsg:        "commit-msg",
		PreRebase:        "pre-rebase",
		PostRewrite:      "post-rewrite",
		PostMerge:        "post-merge",
		PostCommit:       "post-commit",
		PrePush:          "pre-push",
		PostCheckout:     "post-checkout",
	}
}

func NewHookCommands(gitCommon *GitCommon) *ClientHookCommands {
	return &ClientHookCommands{
		gitCommon: gitCommon,
		Hooks:     DefaultClientHooks(),
	}
}

type ClientHookCommands struct {
	gitCommon *GitCommon
	Hooks     ClientHooks
}

func (self *ClientHookCommands) RunClientHook(hookName string, args ...string) error {
	hookPath := filepath.Join(
		self.gitCommon.repoPaths.WorktreeGitDirPath(),
		"hooks",
		hookName,
	)
	// Git silently ignores hooks that are missing or not executable.
	// In those cases this is a no-op and the Git operation proceeds normally.
	err := self.isHookValid(self.gitCommon.Common.Fs, hookPath)
	if err != nil {
		return nil
	}

	// Build argv: [ "/path/to/hook", arg... ]
	argv := append([]string{hookPath}, args...)

	cmd := self.gitCommon.cmd.New(argv)

	err = cmd.Run()
	// If the hook exists and is executable, Git *does* surface any non-zero exit code.
	// Hook failures should be returned to the caller.
	if err != nil {
		return fmt.Errorf("client hook %s failed: %w", hookName, err)
	}

	return nil
}

func (self *ClientHookCommands) isHookValid(fs afero.Fs, hookPath string) error {
	info, err := fs.Stat(hookPath)
	if err != nil {
		return err
	}
	mode := info.Mode()
	executable, err := mode&0o111 != 0, nil

	if !executable {
		err = errors.Errorf("File '%s' is not executable", hookPath)
	}

	return err
}
