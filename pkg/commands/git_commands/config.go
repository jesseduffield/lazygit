package git_commands

import (
	"errors"
	"regexp"
	"strings"

	gogit "github.com/jesseduffield/go-git/v5"
	"github.com/jesseduffield/go-git/v5/config"
	"github.com/jesseduffield/lazygit/pkg/commands/git_config"
	"github.com/jesseduffield/lazygit/pkg/common"
)

type ConfigCommands struct {
	*common.Common

	gitConfig git_config.IGitConfig
	repo      *gogit.Repository
}

func NewConfigCommands(
	common *common.Common,
	gitConfig git_config.IGitConfig,
	repo *gogit.Repository,
) *ConfigCommands {
	return &ConfigCommands{
		Common:    common,
		gitConfig: gitConfig,
		repo:      repo,
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
func (self *ConfigCommands) Branches() (map[string]*config.Branch, error) {
	if self.repo == nil {
		return nil, errors.New("repository is nil")
	}
	conf, err := self.repo.Config()
	if err != nil {
		return nil, err
	}

	return conf.Branches, nil
}

// git-flow config key patterns: legacy uses gitflow.prefix.<type>, git-flow-next uses gitflow.branch.<type>.prefix
const (
	gitFlowLegacyConfigArgs = "--local --get-regexp gitflow.prefix"
	gitFlowNextConfigArgs   = "--local --get-regexp gitflow\\.branch\\..*\\.prefix"
)

func (self *ConfigCommands) GetGitFlowPrefixes() string {
	return self.gitConfig.GetGeneral(gitFlowLegacyConfigArgs)
}

// GetGitFlowNextPrefixes returns raw output of git-flow-next prefix config (gitflow.branch.<type>.prefix).
func (self *ConfigCommands) GetGitFlowNextPrefixes() string {
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

// parseGitFlowPrefixMap parses legacy and git-flow-next config output into a unified prefix -> branchType map.
// Legacy line format: "gitflow.prefix.<type> <prefix>"
// Next line format: "gitflow.branch.<type>.prefix <prefix>"
// Prefixes are normalized with a trailing slash for consistent matching. Legacy entries win on duplicate prefix.
func parseGitFlowPrefixMap(legacyOutput, nextOutput string) map[string]string {
	legacyRegexp := regexp.MustCompile(`gitflow\.prefix\.([^\s]+)\s+(.*)`)
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

// GetGitFlowPrefixMap returns a unified prefix -> branchType map from both legacy and git-flow-next config.
// Git-flow is enabled when this map is non-empty.
func (self *ConfigCommands) GetGitFlowPrefixMap() map[string]string {
	return parseGitFlowPrefixMap(self.GetGitFlowPrefixes(), self.GetGitFlowNextPrefixes())
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
