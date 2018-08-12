package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
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

// GetStashEntries stash entryies
func (c *GitCommand) GetStashEntries() []StashEntry {
	stashEntries := make([]StashEntry, 0)
	rawString, _ := c.OSCommand.RunDirectCommand("git stash list --pretty='%gs'")
	for i, line := range utils.SplitLines(rawString) {
		stashEntries = append(stashEntries, stashEntryFromLine(line, i))
	}
	return stashEntries
}

func stashEntryFromLine(line string, index int) StashEntry {
	return StashEntry{
		Name:          line,
		Index:         index,
		DisplayString: line,
	}
}

// GetStashEntryDiff stash diff
func (c *GitCommand) GetStashEntryDiff(index int) (string, error) {
	return c.OSCommand.RunCommand("git stash show -p --color stash@{" + fmt.Sprint(index) + "}")
}

func includes(array []string, str string) bool {
	for _, arrayStr := range array {
		if arrayStr == str {
			return true
		}
	}
	return false
}

// GetStatusFiles git status files
func (c *GitCommand) GetStatusFiles() []GitFile {
	statusOutput, _ := c.GitStatus()
	statusStrings := utils.SplitLines(statusOutput)
	gitFiles := make([]GitFile, 0)

	for _, statusString := range statusStrings {
		change := statusString[0:2]
		stagedChange := change[0:1]
		unstagedChange := statusString[1:2]
		filename := statusString[3:]
		tracked := !includes([]string{"??", "A "}, change)
		gitFile := GitFile{
			Name:               filename,
			DisplayString:      statusString,
			HasStagedChanges:   !includes([]string{" ", "U", "?"}, stagedChange),
			HasUnstagedChanges: unstagedChange != " ",
			Tracked:            tracked,
			Deleted:            unstagedChange == "D" || stagedChange == "D",
			HasMergeConflicts:  change == "UU",
		}
		gitFiles = append(gitFiles, gitFile)
	}
	c.Log.Info(gitFiles) // TODO: use a dumper-esque log here
	return gitFiles
}

// StashDo modify stash
func (c *GitCommand) StashDo(index int, method string) (string, error) {
	return c.OSCommand.RunCommand("git stash " + method + " stash@{" + fmt.Sprint(index) + "}")
}

// StashSave save stash
func (c *GitCommand) StashSave(message string) (string, error) {
	output, err := c.OSCommand.RunCommand("git stash save \"" + message + "\"")
	if err != nil {
		return output, err
	}
	// if there are no local changes to save, the exit code is 0, but we want
	// to raise an error
	if output == "No local changes to save\n" {
		return output, errors.New(output)
	}
	return output, nil
}

// MergeStatusFiles merge status files
func (c *GitCommand) MergeStatusFiles(oldGitFiles, newGitFiles []GitFile) []GitFile {
	if len(oldGitFiles) == 0 {
		return newGitFiles
	}

	appendedIndexes := make([]int, 0)

	// retain position of files we already could see
	result := make([]GitFile, 0)
	for _, oldGitFile := range oldGitFiles {
		for newIndex, newGitFile := range newGitFiles {
			if oldGitFile.Name == newGitFile.Name {
				result = append(result, newGitFile)
				appendedIndexes = append(appendedIndexes, newIndex)
				break
			}
		}
	}

	// append any new files to the end
	for index, newGitFile := range newGitFiles {
		if !includesInt(appendedIndexes, index) {
			result = append(result, newGitFile)
		}
	}

	return result
}

func (c *GitCommand) verifyInGitRepo() {
	if output, err := c.OSCommand.RunCommand("git status"); err != nil {
		fmt.Println(output)
		os.Exit(1)
	}
}

// GetBranchName branch name
func (c *GitCommand) GetBranchName() (string, error) {
	return c.OSCommand.RunDirectCommand("git symbolic-ref --short HEAD")
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
	r, err := gogit.PlainOpen(".")
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
	return c.Worktree.Reset(&gogit.ResetOptions{Mode: gogit.HardReset})
}

// UpstreamDifferenceCount checks how many pushables/pullables there are for the
// current branch
func (c *GitCommand) UpstreamDifferenceCount() (string, string) {
	pushableCount, err := c.OSCommand.RunDirectCommand("git rev-list @{u}..head --count")
	if err != nil {
		return "?", "?"
	}
	pullableCount, err := c.OSCommand.RunDirectCommand("git rev-list head..@{u} --count")
	if err != nil {
		return "?", "?"
	}
	return strings.TrimSpace(pushableCount), strings.TrimSpace(pullableCount)
}

// GetCommitsToPush Returns the sha's of the commits that have not yet been pushed
// to the remote branch of the current branch
func (c *GitCommand) GetCommitsToPush() []string {
	pushables, err := c.OSCommand.RunDirectCommand("git rev-list @{u}..head --abbrev-commit")
	if err != nil {
		return make([]string, 0)
	}
	return utils.SplitLines(pushables)
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
	return c.OSCommand.RunDirectCommand("git commit --allow-empty --amend -m \"" + name + "\"")
}

// Fetch fetch git repo
func (c *GitCommand) Fetch() (string, error) {
	return c.OSCommand.RunDirectCommand("git fetch")
}

// ResetToCommit reset to commit
func (c *GitCommand) ResetToCommit(sha string) (string, error) {
	return c.OSCommand.RunDirectCommand("git reset " + sha)
}

// NewBranch create new branch
func (c *GitCommand) NewBranch(name string) (string, error) {
	return c.OSCommand.RunDirectCommand("git checkout -b " + name)
}

// DeleteBranch delete branch
func (c *GitCommand) DeleteBranch(branch string) (string, error) {
	return c.OSCommand.RunCommand("git branch -d " + branch)
}

// ListStash list stash
func (c *GitCommand) ListStash() (string, error) {
	return c.OSCommand.RunDirectCommand("git stash list")
}

// Merge merge
func (c *GitCommand) Merge(branchName string) (string, error) {
	return c.OSCommand.RunDirectCommand("git merge --no-edit " + branchName)
}

// AbortMerge abort merge
func (c *GitCommand) AbortMerge() (string, error) {
	return c.OSCommand.RunDirectCommand("git merge --abort")
}

// GitCommit commit to git
func (c *GitCommand) GitCommit(g *gocui.Gui, message string) (string, error) {
	gpgsign, _ := gitconfig.Global("commit.gpgsign")
	if gpgsign != "" {
		sub, err := c.OSCommand.RunSubProcess("git", "commit")
		return "", nil
	}
	userName, err := gitconfig.Username()
	if userName == "" {
		return "", errNoUsername
	}
	userEmail, err := gitconfig.Email()
	_, err = c.Worktree.Commit(message, &gogit.CommitOptions{
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

// GitPull pull from repo
func (c *GitCommand) GitPull() (string, error) {
	return c.OSCommand.RunCommand("git pull --no-edit")
}

// GitPush push to a branch
func (c *GitCommand) GitPush() (string, error) {
	return c.OSCommand.RunDirectCommand("git push -u origin " + state.Branches[0].Name)
}

// SquashPreviousTwoCommits squashes a commit down to the one below it
// retaining the message of the higher commit
func (c *GitCommand) SquashPreviousTwoCommits(message string) (string, error) {
	return c.OSCommand.RunDirectCommand("git reset --soft HEAD^ && git commit --amend -m \"" + message + "\"")
}

// SquashFixupCommit squashes a 'FIXUP' commit into the commit beneath it,
// retaining the commit message of the lower commit
func (c *GitCommand) SquashFixupCommit(branchName string, shaValue string) (string, error) {
	var err error
	commands := []string{
		"git checkout -q " + shaValue,
		"git reset --soft " + shaValue + "^",
		"git commit --amend -C " + shaValue + "^",
		"git rebase --onto HEAD " + shaValue + " " + branchName,
	}
	ret := ""
	for _, command := range commands {
		c.Log.Info(command)
		output, err := c.OSCommand.RunDirectCommand(command)
		ret += output
		if err != nil {
			c.Log.Info(ret)
			break
		}
	}
	if err != nil {
		// We are already in an error state here so we're just going to append
		// the output of these commands
		output, _ := c.OSCommand.RunDirectCommand("git branch -d " + shaValue)
		ret += output
		output, _ = c.OSCommand.RunDirectCommand("git checkout " + branchName)
		ret += output
	}
	return ret, err
}

// CatFile obtain the contents of a file
func (c *GitCommand) CatFile(file string) (string, error) {
	return c.OSCommand.RunDirectCommand("cat " + file)
}

// StageFile stages a file
func (c *GitCommand) StageFile(file string) error {
	_, err := c.OSCommand.RunCommand("git add " + file)
	return err
}

// UnStageFile unstages a file
func (c *GitCommand) UnStageFile(file string, tracked bool) error {
	var command string
	if tracked {
		command = "git reset HEAD "
	} else {
		command = "git rm --cached "
	}
	_, err := c.OSCommand.RunCommand(command + file)
	return err
}

// GitStatus returns the plaintext short status of the repo
func (c *GitCommand) GitStatus() (string, error) {
	return c.OSCommand.RunCommand("git status --untracked-files=all --short")
}

// IsInMergeState states whether we are still mid-merge
func (c *GitCommand) IsInMergeState() (bool, error) {
	output, err := c.OSCommand.RunCommand("git status --untracked-files=all")
	if err != nil {
		return false, err
	}
	return strings.Contains(output, "conclude merge") || strings.Contains(output, "unmerged paths"), nil
}

// RemoveFile directly
func (c *GitCommand) RemoveFile(file GitFile) error {
	// if the file isn't tracked, we assume you want to delete it
	if !file.Tracked {
		_, err := c.OSCommand.RunCommand("rm -rf ./" + file.Name)
		return err
	}
	// if the file is tracked, we assume you want to just check it out
	_, err := c.OSCommand.RunCommand("git checkout " + file.Name)
	return err
}

// Checkout checks out a branch, with --force if you set the force arg to true
func (c *GitCommand) Checkout(branch string, force bool) (string, error) {
	forceArg := ""
	if force {
		forceArg = "--force "
	}
	return c.OSCommand.RunCommand("git checkout " + forceArg + branch)
}

// AddPatch runs a subprocess for adding a patch by patch
// this will eventually be swapped out for a better solution inside the Gui
func (c *GitCommand) AddPatch(g *gocui.Gui, filename string) (*exec.Cmd, error) {
	return c.OSCommand.RunSubProcess("git", "add", "--patch", filename)
}

// GetBranchGraph gets the color-formatted graph of the log for the given branch
// Currently it limits the result to 100 commits, but when we get async stuff
// working we can do lazy loading
func (c *GitCommand) GetBranchGraph(branchName string) (string, error) {
	return c.OSCommand.RunCommand("git log --graph --color --abbrev-commit --decorate --date=relative --pretty=medium -100 " + branchName)
}

// Map (from https://gobyexample.com/collection-functions)
func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func includesString(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// not sure how to genericise this because []interface{} doesn't accept e.g.
// []int arguments
func includesInt(list []int, a int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// GetCommits obtains the commits of the current branch
func (c *GitCommand) GetCommits() []Commit {
	pushables := gogit.GetCommitsToPush()
	log := getLog()
	commits := make([]Commit, 0)
	// now we can split it up and turn it into commits
	lines := utils.RplitLines(log)
	for _, line := range lines {
		splitLine := strings.Split(line, " ")
		sha := splitLine[0]
		pushed := includesString(pushables, sha)
		commits = append(commits, Commit{
			Sha:           sha,
			Name:          strings.Join(splitLine[1:], " "),
			Pushed:        pushed,
			DisplayString: strings.Join(splitLine, " "),
		})
	}
	return commits
}

// GetLog gets the git log (currently limited to 30 commits for performance
// until we work out lazy loading
func (c *GitCommand) GetLog() string {
	// currently limiting to 30 for performance reasons
	// TODO: add lazyloading when you scroll down
	result, err := c.OSCommand.RunDirectCommand("git log --oneline -30")
	if err != nil {
		// assume if there is an error there are no commits yet for this branch
		return ""
	}
	return result
}

// Ignore adds a file to the gitignore for the repo
func (c *GitCommand) Ignore(filename string) {
	if _, err := c.OSCommand.RunDirectCommand("echo '" + filename + "' >> .gitignore"); err != nil {
		panic(err)
	}
}

// Show shows the diff of a commit
func (c *GitCommand) Show(sha string) string {
	result, err := c.OSCommand.RunDirectCommand("git show --color " + sha)
	if err != nil {
		panic(err)
	}
	return result
}

// Diff returns the diff of a file
func (c *GitCommand) Diff(file GitFile) string {
	cachedArg := ""
	if file.HasStagedChanges && !file.HasUnstagedChanges {
		cachedArg = "--cached "
	}
	deletedArg := ""
	if file.Deleted {
		deletedArg = "-- "
	}
	trackedArg := ""
	if !file.Tracked && !file.HasStagedChanges {
		trackedArg = "--no-index /dev/null "
	}
	command := "git diff --color " + cachedArg + deletedArg + trackedArg + file.Name
	// for now we assume an error means the file was deleted
	s, _ := c.OSCommand.RunCommand(command)
	return s
}
