package git_commands

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
)

type WorktreeLoader struct {
	*common.Common
	cmd oscommands.ICmdObjBuilder
}

func NewWorktreeLoader(
	common *common.Common,
	cmd oscommands.ICmdObjBuilder,
) *WorktreeLoader {
	return &WorktreeLoader{
		Common: common,
		cmd:    cmd,
	}
}

func (self *WorktreeLoader) GetWorktrees() ([]*models.Worktree, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cmdArgs := NewGitCmd("worktree").Arg("list", "--porcelain", "-z").ToArgv()
	worktreesOutput, err := self.cmd.New(cmdArgs).DontLog().RunWithOutput()
	if err != nil {
		return nil, err
	}

	splitLines := strings.Split(worktreesOutput, "\x00")

	var worktrees []*models.Worktree
	var currentWorktree *models.Worktree
	for _, splitLine := range splitLines {
		if len(splitLine) == 0 && currentWorktree != nil {

			worktrees = append(worktrees, currentWorktree)
			currentWorktree = nil
			continue
		}
		if strings.HasPrefix(splitLine, "worktree ") {
			main := false
			name := "main"
			path := strings.SplitN(splitLine, " ", 2)[1]
			if len(worktrees) == 0 {
				main = true
			} else {
				name = filepath.Base(path)
			}
			currentWorktree = &models.Worktree{
				Name:    name,
				Path:    path,
				Main:    main,
				Current: path == currentDir,
			}
		}
	}

	/*
		worktree /Users/jbaranick/Source/lazygit
		HEAD f6d6b5dec0432ffa953611700ab9b1ff0089f948
		branch refs/heads/worktree_support

		worktree /Users/jbaranick/Source/lazygit/.worktrees/worktree_tests
		HEAD f6d6b5dec0432ffa953611700ab9b1ff0089f948
		branch refs/heads/worktree_tests
	*/

	return worktrees, nil
}
