package git_commands

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
)

// BranchConfig holds the tracking configuration for a branch.
type BranchConfig struct {
	Remote string
	Merge  string // short ref name of upstream branch
}

type ConfigCommands struct {
	*common.Common

	gitConfig git_config.IGitConfig
}

func NewConfigCommands(
	common *common.Common,
	gitConfig git_config.IGitConfig,
) *ConfigCommands {
	return &ConfigCommands{
		Common:    common,
		gitConfig: gitConfig,
	}
}

type GpgConfigKey string

const (
	CommitGpgSign GpgConfigKey = "commit.gpgSign"
	TagGpgSign    GpgConfigKey = "tag.gpgSign"
)

// NeedsGpgSubprocess tells us whether the user has gpg enabled for the specified action type
// and needs a subprocess because they have a process where they manually
// enter their password every time a GPG action is taken
func (self *ConfigCommands) NeedsGpgSubprocess(key GpgConfigKey) bool {
	overrideGpg := self.UserConfig().Git.OverrideGpg
	if overrideGpg {
		return false
	}

	return self.gitConfig.GetBool(string(key))
}

func (self *ConfigCommands) NeedsGpgSubprocessForCommit() bool {
	return self.NeedsGpgSubprocess(CommitGpgSign)
}

func (self *ConfigCommands) GetGpgTagSign() bool {
	return self.gitConfig.GetBool(string(TagGpgSign))
}

func (self *ConfigCommands) GetCoreEditor() string {
	return self.gitConfig.Get("core.editor")
}

// GetRemoteURL returns current repo remote url
func (self *ConfigCommands) GetRemoteURL() string {
	return self.gitConfig.Get("remote.origin.url")
}

func (self *ConfigCommands) GetShowUntrackedFiles() string {
	return self.gitConfig.Get("status.showUntrackedFiles")
}

// this determines whether the user has configured to push to the remote branch of the same name as the current or not
func (self *ConfigCommands) GetPushToCurrent() bool {
	return self.gitConfig.Get("push.default") == "current"
}

// returns the repo's branches as specified in the git config
func (self *ConfigCommands) Branches(cmd oscommands.ICmdObjBuilder) map[string]*BranchConfig {
	cmdArgs := NewGitCmd("config").
		Arg("--local", "--get-regexp", `^branch\.`).ToArgv()
	output, err := cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err != nil {
		// exit code 1 means no matching keys (no branches with config)
		return nil
	}

	result := make(map[string]*BranchConfig)
	for _, line := range strings.Split(output, "\n") {
		key, value, found := strings.Cut(strings.TrimSpace(line), " ")
		if !found {
			continue
		}
		// key is like "branch.<name>.remote" or "branch.<name>.merge"
		lastDot := strings.LastIndex(key, ".")
		// ignore key like branch.autosetuprebase
		if lastDot < len("branch.") {
			continue
		}
		configKey := key[lastDot+1:]
		branchName := key[len("branch."):lastDot]
		if _, ok := result[branchName]; !ok {
			result[branchName] = &BranchConfig{}
		}
		switch configKey {
		case "remote":
			result[branchName].Remote = value
		case "merge":
			result[branchName].Merge = strings.TrimPrefix(value, "refs/heads/")
		}
	}
	return result
}

func (self *ConfigCommands) GetGitFlowPrefixes() string {
	return self.gitConfig.GetGeneral("--local --get-regexp gitflow.prefix")
}

func (self *ConfigCommands) GetCoreCommentChar() byte {
	if commentCharStr := self.gitConfig.Get("core.commentChar"); len(commentCharStr) == 1 {
		return commentCharStr[0]
	}

	return '#'
}

func (self *ConfigCommands) GetRebaseUpdateRefs() bool {
	return self.gitConfig.GetBool("rebase.updateRefs")
}

func (self *ConfigCommands) GetMergeFF() string {
	return self.gitConfig.Get("merge.ff")
}

func (self *ConfigCommands) DropConfigCache() {
	self.gitConfig.DropCache()
}
