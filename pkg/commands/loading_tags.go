package commands

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func (c *GitCommand) GetTags() ([]*models.Tag, error) {
	// get remote branches, sorted  by creation date (descending)
	// see: https://git-scm.com/docs/git-tag#Documentation/git-tag.txt---sortltkeygt
	remoteBranchesStr, err := c.GetOSCommand().RunCommandWithOutput(
		BuildGitCmdObjFromStr("tag --list --sort=-creatordate"),
	)
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
