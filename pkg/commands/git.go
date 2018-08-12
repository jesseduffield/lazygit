package commands

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
)

// GitCommand is our main git interface
type GitCommand struct {
	Log       *logrus.Logger
	OSCommand *OSCommand
	Worktree  *git.Worktree
	Repo      *git.Repository
}

// NewGitCommand it runs git commands
func NewGitCommand(log *logrus.Logger, osCommand *OSCommand) (*GitCommand, error) {
	gitCommand := &GitCommand{
		Log:       log,
		OSCommand: osCommand,
	}
	return gitCommand, nil
}

// SetupGit sets git repo up
func (c *GitCommand) SetupGit() {
	c.verifyInGitRepo()
	c.navigateToRepoRootDirectory()
	c.setupWorktree()
}

func (c *GitCommand) GitIgnore(filename string) {
	if _, err := c.OSCommand.RunDirectCommand("echo '" + filename + "' >> .gitignore"); err != nil {
		panic(err)
	}
}

func (c *GitCommand) verifyInGitRepo() {
	if output, err := c.OSCommand.RunCommand("git status"); err != nil {
		fmt.Println(output)
		os.Exit(1)
	}
}

func (c *GitCommand) navigateToRepoRootDirectory() {
	_, err := os.Stat(".git")
	for os.IsNotExist(err) {
		c.Log.Debug("going up a directory to find the root")
		os.Chdir("..")
		_, err = os.Stat(".git")
	}
}

func (c *GitCommand) setupWorktree() {
	var err error
	r, err := git.PlainOpen(".")
	if err != nil {
		panic(err)
	}
	c.Repo = r

	w, err := r.Worktree()
	c.Worktree = w
	if err != nil {
		panic(err)
	}
	c.Worktree = w
}
