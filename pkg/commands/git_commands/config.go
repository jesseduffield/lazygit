package git_commands

import (
	"regexp"
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

// git-flow config key patterns: legacy uses gitflow.prefix.<type>, git-flow-next uses gitflow.branch.<type>.prefix
const (
	gitFlowLegacyConfigArgs = "--local --get-regexp gitflow.prefix"
	gitFlowNextConfigArgs   = "--local --get-regexp gitflow\\.branch\\..*\\.prefix"
)

func (self *ConfigCommands) getGitFlowPrefixes() string {
	return self.gitConfig.GetGeneral(gitFlowLegacyConfigArgs)
}

func (self *ConfigCommands) getGitFlowNextPrefixes() string {
	return self.gitConfig.GetGeneral(gitFlowNextConfigArgs)
}

// parseGitFlowLines parses lines matching re (submatch 1 = branch type, 2 = prefix) into prefixToType.
// When overwrite is false, existing keys are left unchanged so legacy entries win over next.
func parseGitFlowLines(output string, re *regexp.Regexp, prefixToType map[string]string, overwrite bool) {
	for line := range strings.SplitSeq(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if m := re.FindStringSubmatch(line); len(m) == 3 {
			prefix := normalizeGitFlowPrefix(m[2])
			if prefix == "" {
				continue
			}
			if overwrite || prefixToType[prefix] == "" {
				prefixToType[prefix] = m[1]
			}
		}
	}
}

// parseGitFlowPrefixMap parses legacy and git-flow-next config output into a unified prefix → branchType map.
// Legacy line format: "gitflow.prefix.<type> <prefix>"
// Next line format: "gitflow.branch.<type>.prefix <prefix>"
// Prefixes are normalized to end in "/". Legacy entries win on duplicate prefix.
func parseGitFlowPrefixMap(legacyOutput, nextOutput string) map[string]string {
	legacyRegexp := regexp.MustCompile(`gitflow\.prefix\.(\S+)\s+(.*)`)
	nextRegexp := regexp.MustCompile(`gitflow\.branch\.([^.]+)\.prefix\s+(.*)`)
	prefixToType := make(map[string]string)
	parseGitFlowLines(legacyOutput, legacyRegexp, prefixToType, true)
	parseGitFlowLines(nextOutput, nextRegexp, prefixToType, false)
	return prefixToType
}

func normalizeGitFlowPrefix(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return ""
	}
	if !strings.HasSuffix(prefix, "/") {
		return prefix + "/"
	}
	return prefix
}

func (self *ConfigCommands) GetGitFlowPrefixMap() map[string]string {
	return parseGitFlowPrefixMap(self.getGitFlowPrefixes(), self.getGitFlowNextPrefixes())
}

func (self *ConfigCommands) GetGitFlowFinishArgs() []string {
	return self.UserConfig().Git.GitFlowFinishArgs
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
