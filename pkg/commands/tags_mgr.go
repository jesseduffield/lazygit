package commands

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	. "github.com/jesseduffield/lazygit/pkg/commands/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

//counterfeiter:generate . ITagsMgr
type ITagsMgr interface {
	Delete(tagName string) error
	LightweightCreate(tagName string, commitSha string) error
	Load() ([]*models.Tag, error)
}

type TagsMgr struct {
	ICommander

	config IGitConfigMgr
}

func NewTagsMgr(
	commander ICommander,
	config IGitConfigMgr,
) *TagsMgr {
	return &TagsMgr{
		ICommander: commander,
		config:     config,
	}
}

func (c *TagsMgr) Load() ([]*models.Tag, error) {
	// get remote branches, sorted  by creation date (descending)
	// see: https://git-scm.com/docs/git-tag#Documentation/git-tag.txt---sortltkeygt
	remoteBranchesStr, err := c.RunWithOutput(
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

func (c *TagsMgr) LightweightCreate(tagName string, commitSha string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("tag %s %s", tagName, commitSha))
}

func (c *TagsMgr) Delete(tagName string) error {
	return c.RunGitCmdFromStr(fmt.Sprintf("tag -d %s", tagName))
}
