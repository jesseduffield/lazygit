package git_commands

import (
	"strings"

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
	// get tags, sorted by creation date (descending)
	// see: https://git-scm.com/docs/git-tag#Documentation/git-tag.txt---sortltkeygt
	cmdArgs := NewGitCmd("for-each-ref").
		Arg("--sort=-creatordate").
		Arg("--format=%(refname)%00%(objecttype)%00%(contents:subject)").
		Arg("refs/tags").
		ToArgv()
	tagsOutput, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err != nil {
		return nil, err
	}

	split := utils.SplitLines(tagsOutput)

	tags := lo.FilterMap(split, func(line string, _ int) (*models.Tag, bool) {
		fields := strings.SplitN(line, "\x00", 3)
		if len(fields) != 3 {
			return nil, false
		}
		tagName := fields[0]
		objectType := fields[1]
		message := fields[2]

		return &models.Tag{
			Name:        strings.TrimPrefix(tagName, "refs/tags/"),
			Message:     message,
			IsAnnotated: objectType == "tag",
		}, true
	})

	return tags, nil
}
