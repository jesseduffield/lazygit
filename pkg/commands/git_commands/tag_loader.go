package git_commands

import (
	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/utils"
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
	tagsOutput, err := self.cmd.New(`git tag --list --sort=-creatordate`).DontLog().RunWithOutput()
	if err != nil {
		return nil, err
	}

	split := utils.SplitLines(tagsOutput)

	tags := slices.Map(split, func(tagName string) *models.Tag {
		return &models.Tag{
			Name: tagName,
		}
	})

	return tags, nil
}
