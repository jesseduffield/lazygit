package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
	gitconfig "github.com/tcnksm/go-gitconfig"
	gogit "gopkg.in/src-d/go-git.v4"
)

// GitCommand is our main git interface
type GitCommand struct {
	Log       *logrus.Entry
	OSCommand *OSCommand
	Worktree  *gogit.Worktree
	Repo      *gogit.Repository
}

// NewGitCommand it runs git commands
func NewGitCommand(log *logrus.Entry, osCommand *OSCommand) (*GitCommand, error) {
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

// GetStashEntries stash entryies
func (c *GitCommand) GetStashEntries() []StashEntry {
	stashEntries := make([]StashEntry, 0)
	rawString, _ := c.OSCommand.RunCommandWithOutput("git stash list --pretty='%gs'")
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
	return c.OSCommand.RunCommandWithOutput("git stash show -p --color stash@{" + fmt.Sprint(index) + "}")
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
func (c *GitCommand) GetStatusFiles() []File {
	statusOutput, _ := c.GitStatus()
	statusStrings := utils.SplitLines(statusOutput)
	files := make([]File, 0)

	for _, statusString := range statusStrings {
		change := statusString[0:2]
		stagedChange := change[0:1]
		unstagedChange := statusString[1:2]
		filename := statusString[3:]
		tracked := !includes([]string{"??", "A ", "AM"}, change)
		file := File{
			Name:               c.OSCommand.Unquote(filename),
			DisplayString:      statusString,
			HasStagedChanges:   !includes([]string{" ", "U", "?"}, stagedChange),
			HasUnstagedChanges: unstagedChange != " ",
			Tracked:            tracked,
			Deleted:            unstagedChange == "D" || stagedChange == "D",
			HasMergeConflicts:  change == "UU",
		}
		files = append(files, file)
	}
	c.Log.Info(files) // TODO: use a dumper-esque log here
	return files
}

// StashDo modify stash
func (c *GitCommand) StashDo(index int, method string) error {
	return c.OSCommand.RunCommand("git stash " + method + " stash@{" + fmt.Sprint(index) + "}")
}

// StashSave save stash
// TODO: before calling this, check if there is anything to save
func (c *GitCommand) StashSave(message string) error {
	return c.OSCommand.RunCommand("git stash save " + c.OSCommand.Quote(message))
}

// MergeStatusFiles merge status files
func (c *GitCommand) MergeStatusFiles(oldFiles, newFiles []File) []File {
	if len(oldFiles) == 0 {
		return newFiles
	}

	appendedIndexes := make([]int, 0)

	// retain position of files we already could see
	result := make([]File, 0)
	for _, oldFile := range oldFiles {
		for newIndex, newFile := range newFiles {
			if oldFile.Name == newFile.Name {
				result = append(result, newFile)
				appendedIndexes = append(appendedIndexes, newIndex)
				break
			}
		}
	}

	// append any new files to the end
	for index, newFile := range newFiles {
		if !includesInt(appendedIndexes, index) {
			result = append(result, newFile)
		}
	}

	return result
}

func (c *GitCommand) verifyInGitRepo() {
	if output, err := c.OSCommand.RunCommandWithOutput("git status"); err != nil {
		fmt.Println(output)
		os.Exit(1)
	}
}

// GetBranchName branch name
func (c *GitCommand) GetBranchName() (string, error) {
	return c.OSCommand.RunCommandWithOutput("git symbolic-ref --short HEAD")
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
	pushableCount, err := c.OSCommand.RunCommandWithOutput("git rev-list @{u}..head --count")
	if err != nil {
		return "?", "?"
	}
	pullableCount, err := c.OSCommand.RunCommandWithOutput("git rev-list head..@{u} --count")
	if err != nil {
		return "?", "?"
	}
	return strings.TrimSpace(pushableCount), strings.TrimSpace(pullableCount)
}

// GetCommitsToPush Returns the sha's of the commits that have not yet been pushed
// to the remote branch of the current branch
func (c *GitCommand) GetCommitsToPush() []string {
	pushables, err := c.OSCommand.RunCommandWithOutput("git rev-list @{u}..head --abbrev-commit")
	if err != nil {
		return make([]string, 0)
	}
	return utils.SplitLines(pushables)
}

// RenameCommit renames the topmost commit with the given name
func (c *GitCommand) RenameCommit(name string) error {
	return c.OSCommand.RunCommand("git commit --allow-empty --amend -m " + c.OSCommand.Quote(name))
}

// Fetch fetch git repo
func (c *GitCommand) Fetch() error {
	return c.OSCommand.RunCommand("git fetch")
}

// ResetToCommit reset to commit
func (c *GitCommand) ResetToCommit(sha string) error {
	return c.OSCommand.RunCommand("git reset " + sha)
}

// NewBranch create new branch
func (c *GitCommand) NewBranch(name string) error {
	return c.OSCommand.RunCommand("git checkout -b " + name)
}

// DeleteBranch delete branch
func (c *GitCommand) DeleteBranch(branch string, force bool) error {
	var command string
	if force {
		command = "git branch -D "
	} else {
		command = "git branch -d "
	}
	return c.OSCommand.RunCommand(command + branch)
}

// ListStash list stash
func (c *GitCommand) ListStash() (string, error) {
	return c.OSCommand.RunCommandWithOutput("git stash list")
}

// Merge merge
func (c *GitCommand) Merge(branchName string) error {
	return c.OSCommand.RunCommand("git merge --no-edit " + branchName)
}

// AbortMerge abort merge
func (c *GitCommand) AbortMerge() error {
	return c.OSCommand.RunCommand("git merge --abort")
}

// UsingGpg tells us whether the user has gpg enabled so that we can know
// whether we need to run a subprocess to allow them to enter their password
func (c *GitCommand) UsingGpg() bool {
	gpgsign, _ := gitconfig.Global("commit.gpgsign")
	if gpgsign == "" {
		gpgsign, _ = gitconfig.Local("commit.gpgsign")
	}
	if gpgsign == "" {
		return false
	}
	return true
}

// Commit commit to git
func (c *GitCommand) Commit(g *gocui.Gui, message string) (*exec.Cmd, error) {
	command := "git commit -m " + c.OSCommand.Quote(message)
	if c.UsingGpg() {
		return c.OSCommand.PrepareSubProcess(c.OSCommand.Platform.shell, c.OSCommand.Platform.shellArg, command), nil
	}
	return nil, c.OSCommand.RunCommand(command)
}

// Pull pull from repo
func (c *GitCommand) Pull() error {
	return c.OSCommand.RunCommand("git pull --no-edit")
}

// Push push to a branch
func (c *GitCommand) Push(branchName string, force bool) error {
	forceFlag := ""
	if force {
		forceFlag = "--force-with-lease "
	}
	return c.OSCommand.RunCommand("git push " + forceFlag + "-u origin " + branchName)
}

// SquashPreviousTwoCommits squashes a commit down to the one below it
// retaining the message of the higher commit
func (c *GitCommand) SquashPreviousTwoCommits(message string) error {
	// TODO: test this
	err := c.OSCommand.RunCommand("git reset --soft HEAD^")
	if err != nil {
		return err
	}
	// TODO: if password is required, we need to return a subprocess
	return c.OSCommand.RunCommand("git commit --amend -m " + c.OSCommand.Quote(message))
}

// SquashFixupCommit squashes a 'FIXUP' commit into the commit beneath it,
// retaining the commit message of the lower commit
func (c *GitCommand) SquashFixupCommit(branchName string, shaValue string) error {
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
		output, err := c.OSCommand.RunCommandWithOutput(command)
		ret += output
		if err != nil {
			c.Log.Info(ret)
			break
		}
	}
	if err != nil {
		// We are already in an error state here so we're just going to append
		// the output of these commands
		output, _ := c.OSCommand.RunCommandWithOutput("git branch -d " + shaValue)
		ret += output
		output, _ = c.OSCommand.RunCommandWithOutput("git checkout " + branchName)
		ret += output
	}
	if err != nil {
		return errors.New(ret)
	}
	return nil
}

// CatFile obtain the contents of a file
func (c *GitCommand) CatFile(fileName string) (string, error) {
	return c.OSCommand.RunCommandWithOutput("cat " + c.OSCommand.Quote(fileName))
}

// StageFile stages a file
func (c *GitCommand) StageFile(fileName string) error {
	return c.OSCommand.RunCommand("git add " + c.OSCommand.Quote(fileName))
}

// UnStageFile unstages a file
func (c *GitCommand) UnStageFile(fileName string, tracked bool) error {
	var command string
	if tracked {
		command = "git reset HEAD "
	} else {
		command = "git rm --cached "
	}
	return c.OSCommand.RunCommand(command + c.OSCommand.Quote(fileName))
}

// GitStatus returns the plaintext short status of the repo
func (c *GitCommand) GitStatus() (string, error) {
	return c.OSCommand.RunCommandWithOutput("git status --untracked-files=all --short")
}

// IsInMergeState states whether we are still mid-merge
func (c *GitCommand) IsInMergeState() (bool, error) {
	output, err := c.OSCommand.RunCommandWithOutput("git status --untracked-files=all")
	if err != nil {
		return false, err
	}
	return strings.Contains(output, "conclude merge") || strings.Contains(output, "unmerged paths"), nil
}

// RemoveFile directly
func (c *GitCommand) RemoveFile(file File) error {
	// if the file isn't tracked, we assume you want to delete it
	if file.HasStagedChanges {
		if err := c.OSCommand.RunCommand("git reset -- " + file.Name); err != nil {
			return err
		}
	}
	if !file.Tracked {
		return os.RemoveAll(file.Name)
	}
	// if the file is tracked, we assume you want to just check it out
	return c.OSCommand.RunCommand("git checkout -- " + file.Name)
}

// Checkout checks out a branch, with --force if you set the force arg to true
func (c *GitCommand) Checkout(branch string, force bool) error {
	forceArg := ""
	if force {
		forceArg = "--force "
	}
	return c.OSCommand.RunCommand("git checkout " + forceArg + branch)
}

// AddPatch prepares a subprocess for adding a patch by patch
// this will eventually be swapped out for a better solution inside the Gui
func (c *GitCommand) AddPatch(filename string) *exec.Cmd {
	return c.OSCommand.PrepareSubProcess("git", "add", "--patch", filename)
}

// PrepareCommitSubProcess prepares a subprocess for `git commit`
func (c *GitCommand) PrepareCommitSubProcess() *exec.Cmd {
	return c.OSCommand.PrepareSubProcess("git", "commit")
}

// GetBranchGraph gets the color-formatted graph of the log for the given branch
// Currently it limits the result to 100 commits, but when we get async stuff
// working we can do lazy loading
func (c *GitCommand) GetBranchGraph(branchName string) (string, error) {
	return c.OSCommand.RunCommandWithOutput("git log --graph --color --abbrev-commit --decorate --date=relative --pretty=medium -100 " + branchName)
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
	pushables := c.GetCommitsToPush()
	log := c.GetLog()
	commits := make([]Commit, 0)
	// now we can split it up and turn it into commits
	lines := utils.SplitLines(log)
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
	result, err := c.OSCommand.RunCommandWithOutput("git log --oneline -30")
	if err != nil {
		// assume if there is an error there are no commits yet for this branch
		return ""
	}
	return result
}

// Ignore adds a file to the gitignore for the repo
func (c *GitCommand) Ignore(filename string) error {
	return c.OSCommand.AppendLineToFile(".gitignore", filename)
}

// Show shows the diff of a commit
func (c *GitCommand) Show(sha string) string {
	result, err := c.OSCommand.RunCommandWithOutput("git show --color " + sha)
	if err != nil {
		panic(err)
	}
	return result
}

// Diff returns the diff of a file
func (c *GitCommand) Diff(file File) string {
	cachedArg := ""
	fileName := c.OSCommand.Quote(file.Name)
	if file.HasStagedChanges && !file.HasUnstagedChanges {
		cachedArg = "--cached"
	}
	trackedArg := "--"
	if !file.Tracked && !file.HasStagedChanges {
		trackedArg = "--no-index /dev/null"
	}
	command := fmt.Sprintf("%s %s %s %s", "git diff --color ", cachedArg, trackedArg, fileName)

	// for now we assume an error means the file was deleted
	s, _ := c.OSCommand.RunCommandWithOutput(command)
	return s
}
