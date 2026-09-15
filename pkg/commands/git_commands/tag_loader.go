package git_commands

import (
	"regexp"
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
	gitCommon *GitCommon
}

func NewTagLoader(
	common *common.Common,
	cmd oscommands.ICmdObjBuilder,
	gitCommon *GitCommon,
) *TagLoader {
	return &TagLoader{
		Common: common,
		cmd:    cmd,
		gitCommon: gitCommon,
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

	// SVN 仓库：追加扫描 refs/remotes/git-svn/tags/* 下的引用
	if self.gitCommon != nil && self.gitCommon.IsSvnRepo() {
		tagsPaths, err := self.gitCommon.Svn.GetTagsRefsPaths()
		if err == nil && len(tagsPaths) > 0 {
			for _, tagsPath := range tagsPaths {
				svnTags, err := self.getTagsFromPath(tagsPath)
				if err == nil {
					tags = append(tags, svnTags...)
				}
			}
		}
	}
	return tags, nil
}

// getTagsFromPath 从指定 refs 路径下获取所有 tag 引用
func (self *TagLoader) getTagsFromPath(basePath string) ([]*models.Tag, error) {
	var svnTags []*models.Tag
	cmdArgs := NewGitCmd("for-each-ref").
	Arg("--sort=-creatordate").
	Arg("--format=%(refname)").
	Arg(basePath).
	ToArgv()

	err := self.cmd.New(cmdArgs).DontLog().RunAndProcessLines(func(line string) (bool, error){
		line = strings.TrimSpace(line)
		if line == "" {
			return false, nil
		}

		tagName := strings.TrimPrefix(line, basePath+"/")
		if tagName == "" || tagName == line {
			return false, nil
		}

		svnTags = append(svnTags, &models.Tag{
			Name: tagName,
			Message: "(SVN tag)",
			FullRefNameOverride: line,
		})
		return false, nil
	})

	return svnTags, err
}
