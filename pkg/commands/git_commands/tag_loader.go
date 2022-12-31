package git_commands

import (
	"fmt"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

type TagLoader struct {
	*common.Common
	version *GitVersion
	cmd     oscommands.ICmdObjBuilder
}

func NewTagLoader(
	common *common.Common,
	version *GitVersion,
	cmd oscommands.ICmdObjBuilder,
) *TagLoader {
	return &TagLoader{
		Common:  common,
		version: version,
		cmd:     cmd,
	}
}

func (self *TagLoader) GetTags() ([]*models.Tag, error) {
	// get remote branches, sorted  by creation date (descending)
	// see: https://git-scm.com/docs/git-tag#Documentation/git-tag.txt---sortltkeygt
	sortKey := "-creatordate"
	if self.version.IsOlderThan(2, 7, 0) {
		sortKey = "-v:refname"
	}

	tagsOutput, err := self.cmd.New(fmt.Sprintf(`git tag --list --sort=%s`, sortKey)).DontLog().RunWithOutput()
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
