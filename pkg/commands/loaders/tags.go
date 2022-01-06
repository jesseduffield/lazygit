package loaders

import (
	"strings"

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
	remoteBranchesStr, err := self.cmd.New(`git tag --list --sort=-creatordate`).DontLog().RunWithOutput()
	if err != nil {
		return nil, err
	}

	content := utils.TrimTrailingNewline(remoteBranchesStr)
	if content == "" {
		return nil, nil
	}

	split := strings.Split(content, "\n")

	// first step is to get our remotes from go-git
	tags := make([]*models.Tag, len(split))
	for i, tagName := range split {
		tags[i] = &models.Tag{
			Name: tagName,
		}
	}

	return tags, nil
}
