package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/git"
	gitconfig "github.com/tcnksm/go-gitconfig"
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// GitCommand is our main git interface
type GitCommand struct {
	Log       *logrus.Logger
	OSCommand *OSCommand
	Worktree  *gogit.Worktree
	Repo      *gogit.Repository
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

// GitIgnore adds a file to the .gitignore of the repo
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
	r, err := git.PlainOpen(".")
	if err != nil {
		panic(err)
	}
	c.Repo = r

	w, err := r.Worktree()
	if err != nil {
		panic(err)
	}
	c.Worktree = w
}

// ResetHard does the equivalent of `git reset --hard HEAD`
func (c *GitCommand) ResetHard() error {
	return c.Worktree.Reset(&gogit.ResetOptions{Mode: git.HardReset})
}

// UpstreamDifferenceCount checks how many pushables/pullables there are for the
// current branch
func (c *GitCommand) UpstreamDifferenceCount() (string, string) {
	pushableCount, err := c.OSCommand.runDirectCommand("git rev-list @{u}..head --count")
	if err != nil {
		return "?", "?"
	}
	pullableCount, err := c.OSCommand.runDirectCommand("git rev-list head..@{u} --count")
	if err != nil {
		return "?", "?"
	}
	return strings.TrimSpace(pushableCount), strings.TrimSpace(pullableCount)
}

// GetCommitsToPush Returns the sha's of the commits that have not yet been pushed
// to the remote branch of the current branch
func (c *GitCommand) GetCommitsToPush() []string {
	pushables, err := c.OSCommand.runDirectCommand("git rev-list @{u}..head --abbrev-commit")
	if err != nil {
		return make([]string, 0)
	}
	return splitLines(pushables)
}

// GetGitBranches returns a list of branches for the current repo, with recency
// values stored against those that are in the reflog
func (c *GitCommand) GetGitBranches() []Branch {
	builder := git.newBranchListBuilder()
	return builder.build()
}

// BranchIncluded states whether a branch is included in a list of branches,
// with a case insensitive comparison on name
func (c *GitCommand) BranchIncluded(branchName string, branches []Branch) bool {
	for _, existingBranch := range branches {
		if strings.ToLower(existingBranch.Name) == strings.ToLower(branchName) {
			return true
		}
	}
	return false
}

// RenameCommit renames the topmost commit with the given name
func (c *GitCommand) RenameCommit(name string) (string, error) {
	return c.OSCommand.runDirectCommand("git commit --allow-empty --amend -m \"" + name + "\"")
}

func (c *GitCommand) Fetch() (string, error) {
	return c.OSCommand.runDirectCommand("git fetch")
}

func (c *GitCommand) ResetToCommit(sha string) (string, error) {
	return c.OSCommand.runDirectCommand("git reset " + sha)
}

func (c *GitCommand) NewBranch(name string) (string, error) {
	return c.OSCommand.runDirectCommand("git checkout -b " + name)
}

func (c *GitCommand) DeleteBranch(branch string) (string, error) {
	return runCommand("git branch -d " + branch)
}

func (c *GitCommand) ListStash() (string, error) {
	return c.OSCommand.runDirectCommand("git stash list")
}

func (c *GitCommand) Merge(branchName string) (string, error) {
	return c.OSCommand.runDirectCommand("git merge --no-edit " + branchName)
}

func (c *GitCommand) AbortMerge() (string, error) {
	return c.OSCommand.runDirectCommand("git merge --abort")
}

func gitCommit(g *gocui.Gui, message string) (string, error) {
	gpgsign, _ := gitconfig.Global("commit.gpgsign")
	if gpgsign != "" {
		runSubProcess(g, "git", "commit")
		return "", nil
	}
	userName, err := gitconfig.Username()
	if userName == "" {
		return "", errNoUsername
	}
	userEmail, err := gitconfig.Email()
	_, err = w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  userName,
			Email: userEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err.Error(), err
	}
	return "", nil
}

func gitPull() (string, error) {
	return runDirectCommand("git pull --no-edit")
}

func gitPush() (string, error) {
	return runDirectCommand("git push -u origin " + state.Branches[0].Name)
}

func gitSquashPreviousTwoCommits(message string) (string, error) {
	return runDirectCommand("git reset --soft HEAD^ && git commit --amend -m \"" + message + "\"")
}

func gitSquashFixupCommit(branchName string, shaValue string) (string, error) {
	var err error
	commands := []string{
		"git checkout -q " + shaValue,
		"git reset --soft " + shaValue + "^",
		"git commit --amend -C " + shaValue + "^",
		"git rebase --onto HEAD " + shaValue + " " + branchName,
	}
	ret := ""
	for _, command := range commands {
		devLog(command)
		output, err := runDirectCommand(command)
		ret += output
		if err != nil {
			devLog(ret)
			break
		}
	}
	if err != nil {
		// We are already in an error state here so we're just going to append
		// the output of these commands
		ret += runDirectCommandIgnoringError("git branch -d " + shaValue)
		ret += runDirectCommandIgnoringError("git checkout " + branchName)
	}
	return ret, err
}
