package git_commands

import (
	"regexp"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

type FlowCommands struct {
	*GitCommon
}

func NewFlowCommands(
	gitCommon *GitCommon,
) *FlowCommands {
	return &FlowCommands{
		GitCommon: gitCommon,
	}
}

func (self *FlowCommands) GitFlowEnabled() bool {
	return self.config.GetGitFlowPrefixes() != ""
}

func (self *FlowCommands) FinishCmdObj(branchName string) (oscommands.ICmdObj, error) {
	prefixes := self.config.GetGitFlowPrefixes()

	// need to find out what kind of branch this is
	prefix := strings.SplitAfterN(branchName, "/", 2)[0]
	suffix := strings.Replace(branchName, prefix, "", 1)

	branchType := ""
	for _, line := range strings.Split(strings.TrimSpace(prefixes), "\n") {
		if strings.HasPrefix(line, "gitflow.prefix.") && strings.HasSuffix(line, prefix) {
			regex := regexp.MustCompile("gitflow.prefix.([^ ]*) .*")
			matches := regex.FindAllStringSubmatch(line, 1)

			if len(matches) > 0 && len(matches[0]) > 1 {
				branchType = matches[0][1]
				break
			}
		}
	}

	if branchType == "" {
		return nil, errors.New(self.Tr.NotAGitFlowBranch)
	}

	return self.cmd.New("git flow " + branchType + " finish " + suffix), nil
}

func (self *FlowCommands) StartCmdObj(branchType string, name string) oscommands.ICmdObj {
	return self.cmd.New("git flow " + branchType + " start " + name)
}
