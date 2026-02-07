package git_commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestResolveHookPath(t *testing.T) {
	scenarios := []struct {
		name         string
		hook         ClientHook
		configValue  string
		configError  error
		expectedPath string
	}{
		{
			name:         "core.hooksPath not set",
			hook:         HookPreCommit,
			configValue:  "",
			configError:  errors.New("not set"),
			expectedPath: "/repo/.git/hooks/pre-commit",
		},
		{
			name:         "core.hooksPath absolute",
			hook:         HookCommitMsg,
			configValue:  "/custom/hooks",
			expectedPath: "/custom/hooks/commit-msg",
		},
		{
			name:         "core.hooksPath relative",
			hook:         HookPrepareCommitMsg,
			configValue:  ".githooks",
			expectedPath: "/repo/.githooks/prepare-commit-msg",
		},
		{
			name:         "core.hooksPath trims whitespace",
			hook:         HookPrePush,
			configValue:  "   hooks   ",
			expectedPath: "/repo/hooks/pre-push",
		},
		{
			name:         "git config returns error",
			hook:         HookPostCommit,
			configValue:  "",
			configError:  errors.New("config error"),
			expectedPath: "/repo/.git/hooks/post-commit",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			// Create fresh runner for each test
			runner := oscommands.NewFakeRunner(t)
			cmdArgs := NewGitCmd("config").Arg("--get", "core.hooksPath").ToArgv()
			runner.ExpectArgs(cmdArgs, s.configValue, s.configError)

			fs := afero.NewMemMapFs()
			repoPath := "/repo"
			repoDir := "/repo/.git"
			cmdBuilder := oscommands.NewDummyCmdObjBuilder(runner)

			gitCommon := &GitCommon{
				Common: &common.Common{
					Fs: fs,
				},
				cmd: cmdBuilder,
				repoPaths: &RepoPaths{
					repoPath:           repoPath,
					worktreeGitDirPath: repoDir,
				},
			}
			hooks := NewClientHookCommands(gitCommon)

			hookPath := hooks.resolveHookPath(s.hook)
			assert.Equal(t, s.expectedPath, hookPath)
		})
	}
}

func TestRunClientHook(t *testing.T) {
	scenarios := []struct {
		name          string
		hook          ClientHook
		configValue   string
		configError   error
		hookExists    bool
		hookContent   string
		hookMode      os.FileMode
		hookRunError  error
		expectedError string
	}{
		{
			name:        "hook exists and succeeds",
			hook:        HookPreCommit,
			configError: errors.New("not set"), // core.hooksPath not configured
			hookExists:  true,
			hookContent: "#!/bin/sh\nexit 0",
			hookMode:    0o755,
		},
		{
			name:          "hook exists but fails",
			hook:          HookPrepareCommitMsg,
			configError:   errors.New("not set"),
			hookExists:    true,
			hookContent:   "#!/bin/sh\nexit 1",
			hookMode:      0o755,
			hookRunError:  errors.New("exit status 1"),
			expectedError: "client hook prepare-commit-msg failed",
		},
		{
			name:          "hook exists but is not executable",
			hook:          HookCommitMsg,
			configError:   errors.New("not set"),
			hookExists:    true,
			hookContent:   "#!/bin/sh",
			hookMode:      0o644,
			hookRunError:  errors.New("permission denied"),
			expectedError: "client hook commit-msg failed",
		},
		{
			name:        "hook does not exist - no error",
			hook:        HookPostCommit,
			configError: errors.New("not set"),
			hookExists:  false,
		},
		{
			name:        "custom hooks path - absolute",
			hook:        HookPrePush,
			configValue: "/custom/hooks",
			hookExists:  true,
			hookContent: "#!/bin/sh",
			hookMode:    0o755,
		},
		{
			name:        "custom hooks path - relative",
			hook:        HookPostCheckout,
			configValue: ".githooks",
			hookExists:  true,
			hookContent: "#!/bin/sh",
			hookMode:    0o755,
		},
		{
			name:        "core.hooksPath with whitespace",
			hook:        HookPreRebase,
			configValue: "  .githooks  \n",
			hookExists:  true,
			hookContent: "#!/bin/sh",
			hookMode:    0o755,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			repoPath := "/repo"
			repoDir := "/repo/.git"
			runner := oscommands.NewFakeRunner(t)

			cmdArgs := NewGitCmd("config").Arg("--get", "core.hooksPath").ToArgv()
			runner.ExpectArgs(cmdArgs, s.configValue, s.configError)

			// Calculate expected hook path based on config
			var hookPath string
			trimmedConfig := strings.TrimSpace(s.configValue)
			if s.configError != nil || trimmedConfig == "" {
				// Default: .git/hooks
				hookPath = filepath.Join(repoDir, "hooks", string(s.hook))
			} else if filepath.IsAbs(trimmedConfig) {
				// Absolute path
				hookPath = filepath.Join(trimmedConfig, string(s.hook))
			} else {
				// Relative path - relative to repo root, not .git
				hookPath = filepath.Join(repoPath, trimmedConfig, string(s.hook))
			}

			if s.hookExists {
				hookDir := filepath.Dir(hookPath)
				_ = fs.MkdirAll(hookDir, 0o755)
				_ = afero.WriteFile(fs, hookPath, []byte(s.hookContent), s.hookMode)

				runner.ExpectArgs([]string{hookPath}, "", s.hookRunError)
			}

			cmdBuilder := oscommands.NewDummyCmdObjBuilder(runner)
			gitCommon := &GitCommon{
				Common: &common.Common{
					Fs: fs,
				},
				cmd: cmdBuilder,
				repoPaths: &RepoPaths{
					repoPath:           repoPath,
					worktreeGitDirPath: repoDir,
				},
			}
			hooks := NewClientHookCommands(gitCommon)

			err := hooks.RunClientHook(s.hook)

			if s.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), s.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
