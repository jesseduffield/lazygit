package git_commands

import (
	"regexp"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type TagLoader struct {
	*common.Common
	cmd oscommands.ICmdObjBuilder
}

func NewTagLoader(
	common *common.Common,
	cmd oscommands.ICmdObjBuilder,
) *TagLoader {
	return &TagLoader{
		Common: common,
		cmd:    cmd,
	}
}

func (self *TagLoader) GetTags() ([]*models.Tag, error) {
	// get remote branches, sorted  by creation date (descending)
	// see: https://git-scm.com/docs/git-tag#Documentation/git-tag.txt---sortltkeygt
	cmdArgs := NewGitCmd("tag").Arg("--list", "-n", "--sort=-creatordate").ToArgv()
	tagsOutput, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err != nil {
		return nil, err
	}

	split := utils.SplitLines(tagsOutput)

	lineRegex := regexp.MustCompile(`^([^\s]+)(\s+)?(.*)$`)

	tags := lo.Map(split, func(line string, _ int) *models.Tag {
		matches := lineRegex.FindStringSubmatch(line)
		tagName := matches[1]
		message := ""
		if len(matches) > 3 {
			message = matches[3]
		}

		return &models.Tag{
			Name:    tagName,
			Message: message,
		}
	})

	return tags, nil
}
