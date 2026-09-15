package daemon

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/env"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/stefanhaller/git-todo-parser/todo"
)

type ChangeTodoAction struct {
	Hash      string
	NewAction todo.TodoCommand
	Flag      string
}

func handleInteractiveRebase(log *logrus.Entry, f func(path string) error) error {
	log.Info("Lazygit invoked as interactive rebase demon")
	log.Info("args: ", os.Args)
	path := os.Args[1]

	if strings.HasSuffix(path, "git-rebase-todo") {
		err := utils.RemoveUpdateRefsForCopiedBranch(path, getCommentChar())
		if err != nil {
			return err
		}
		return f(path)
	} else if strings.HasSuffix(path, filepath.Join(gitDir(), "COMMIT_EDITMSG")) { // TODO: test
		// if we are rebasing and squashing, we'll see a COMMIT_EDITMSG
		// but in this case we don't need to edit it, so we'll just return
	} else {
		log.Info("Lazygit demon did not match on any use cases")
	}

	return nil
}

func gitDir() string {
	dir := env.GetGitDirEnv()
	if dir == "" {
		return ".git"
	}
	return dir
}
