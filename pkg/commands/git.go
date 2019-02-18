package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mgutz/str"

	"github.com/fatih/color"
	"github.com/go-errors/errors"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
	gitconfig "github.com/tcnksm/go-gitconfig"
	gogit "gopkg.in/src-d/go-git.v4"
)

func verifyInGitRepo(runCmd func(string) error) error {
	return runCmd("git status")
}

func navigateToRepoRootDirectory(stat func(string) (os.FileInfo, error), chdir func(string) error) error {
	for {
		f, err := stat(".git")

		if err == nil && f.IsDir() {
			return nil
		}

		if !os.IsNotExist(err) {
			return errors.Wrap(err, 0)
		}

		if err = chdir(".."); err != nil {
			return errors.Wrap(err, 0)
		}
	}
}

func setupRepositoryAndWorktree(openGitRepository func(string) (*gogit.Repository, error), sLocalize func(string) string) (repository *gogit.Repository, worktree *gogit.Worktree, err error) {
	repository, err = openGitRepository(".")

	if err != nil {
		if strings.Contains(err.Error(), `unquoted '\' must be followed by new line`) {
			return nil, nil, errors.New(sLocalize("GitconfigParseErr"))
		}

		return
	}

	worktree, err = repository.Worktree()

	if err != nil {
		return
	}

	return
}

// GitCommand is our main git interface
type GitCommand struct {
	Log                *logrus.Entry
	OSCommand          *OSCommand
	Worktree           *gogit.Worktree
	Repo               *gogit.Repository
	Tr                 *i18n.Localizer
	Config             config.AppConfigurer
	getGlobalGitConfig func(string) (string, error)
	getLocalGitConfig  func(string) (string, error)
	removeFile         func(string) error
}

// NewGitCommand it runs git commands
func NewGitCommand(log *logrus.Entry, osCommand *OSCommand, tr *i18n.Localizer, config config.AppConfigurer) (*GitCommand, error) {
	var worktree *gogit.Worktree
	var repo *gogit.Repository

	fs := []func() error{
		func() error {
			return verifyInGitRepo(osCommand.RunCommand)
		},
		func() error {
			return navigateToRepoRootDirectory(os.Stat, os.Chdir)
		},
		func() error {
			var err error
			repo, worktree, err = setupRepositoryAndWorktree(gogit.PlainOpen, tr.SLocalize)
			return err
		},
	}

	for _, f := range fs {
		if err := f(); err != nil {
			return nil, err
		}
	}

	return &GitCommand{
		Log:                log,
		OSCommand:          osCommand,
		Tr:                 tr,
		Worktree:           worktree,
		Repo:               repo,
		Config:             config,
		getGlobalGitConfig: gitconfig.Global,
		getLocalGitConfig:  gitconfig.Local,
		removeFile:         os.RemoveAll,
	}, nil
}

// GetStashEntries stash entryies
func (c *GitCommand) GetStashEntries() []*StashEntry {
	rawString, _ := c.OSCommand.RunCommandWithOutput("git stash list --pretty='%gs'")
	stashEntries := []*StashEntry{}
	for i, line := range utils.SplitLines(rawString) {
		stashEntries = append(stashEntries, stashEntryFromLine(line, i))
	}
	return stashEntries
}

func stashEntryFromLine(line string, index int) *StashEntry {
	return &StashEntry{
		Name:          line,
		Index:         index,
		DisplayString: line,
	}
}

// GetStashEntryDiff stash diff
func (c *GitCommand) GetStashEntryDiff(index int) (string, error) {
	return c.OSCommand.RunCommandWithOutput("git stash show -p --color stash@{" + fmt.Sprint(index) + "}")
}

// GetStatusFiles git status files
func (c *GitCommand) GetStatusFiles() []*File {
	statusOutput, _ := c.GitStatus()
	statusStrings := utils.SplitLines(statusOutput)
	files := []*File{}

	for _, statusString := range statusStrings {
		change := statusString[0:2]
		stagedChange := change[0:1]
		unstagedChange := statusString[1:2]
		filename := c.OSCommand.Unquote(statusString[3:])
		_, untracked := map[string]bool{"??": true, "A ": true, "AM": true}[change]
		_, hasNoStagedChanges := map[string]bool{" ": true, "U": true, "?": true}[stagedChange]

		file := &File{
			Name:               filename,
			DisplayString:      statusString,
			HasStagedChanges:   !hasNoStagedChanges,
			HasUnstagedChanges: unstagedChange != " ",
			Tracked:            !untracked,
			Deleted:            unstagedChange == "D" || stagedChange == "D",
			HasMergeConflicts:  change == "UU" || change == "AA",
			Type:               c.OSCommand.FileType(filename),
		}
		files = append(files, file)
	}
	return files
}

// StashDo modify stash
func (c *GitCommand) StashDo(index int, method string) error {
	return c.OSCommand.RunCommand(fmt.Sprintf("git stash %s stash@{%d}", method, index))
}

// StashSave save stash
// TODO: before calling this, check if there is anything to save
func (c *GitCommand) StashSave(message string) error {
	return c.OSCommand.RunCommand(fmt.Sprintf("git stash save %s", c.OSCommand.Quote(message)))
}

// MergeStatusFiles merge status files
func (c *GitCommand) MergeStatusFiles(oldFiles, newFiles []*File) []*File {
	if len(oldFiles) == 0 {
		return newFiles
	}

	appendedIndexes := []int{}

	// retain position of files we already could see
	result := []*File{}
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

func includesInt(list []int, a int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// ResetAndClean removes all unstaged changes and removes all untracked files
func (c *GitCommand) ResetAndClean() error {
	if err := c.OSCommand.RunCommand("git reset --hard HEAD"); err != nil {
		return err
	}

	return c.OSCommand.RunCommand("git clean -fd")
}

func (c *GitCommand) GetCurrentBranchUpstreamDifferenceCount() (string, string) {
	return c.GetCommitDifferences("HEAD", "@{u}")
}

func (c *GitCommand) GetBranchUpstreamDifferenceCount(branchName string) (string, string) {
	upstream := "origin" // hardcoded for now
	return c.GetCommitDifferences(branchName, fmt.Sprintf("%s/%s", upstream, branchName))
}

// GetCommitDifferences checks how many pushables/pullables there are for the
// current branch
func (c *GitCommand) GetCommitDifferences(from, to string) (string, string) {
	command := "git rev-list %s..%s --count"
	pushableCount, err := c.OSCommand.RunCommandWithOutput(fmt.Sprintf(command, to, from))
	if err != nil {
		return "?", "?"
	}
	pullableCount, err := c.OSCommand.RunCommandWithOutput(fmt.Sprintf(command, from, to))
	if err != nil {
		return "?", "?"
	}
	return strings.TrimSpace(pushableCount), strings.TrimSpace(pullableCount)
}

// GetUnpushedCommits Returns the sha's of the commits that have not yet been pushed
// to the remote branch of the current branch, a map is returned to ease look up
func (c *GitCommand) GetUnpushedCommits() map[string]bool {
	pushables := map[string]bool{}
	o, err := c.OSCommand.RunCommandWithOutput("git rev-list @{u}..HEAD --abbrev-commit")
	if err != nil {
		return pushables
	}
	for _, p := range utils.SplitLines(o) {
		pushables[p] = true
	}

	return pushables
}

// RenameCommit renames the topmost commit with the given name
func (c *GitCommand) RenameCommit(name string) error {
	return c.OSCommand.RunCommand(fmt.Sprintf("git commit --allow-empty --amend -m %s", c.OSCommand.Quote(name)))
}

func (c *GitCommand) RebaseBranch(onto string) error {
	curBranch, err := c.CurrentBranchName()
	if err != nil {
		return err
	}

	return c.OSCommand.RunCommand(fmt.Sprintf("git rebase --autostash %s %s ", onto, curBranch))
}

// Fetch fetch git repo
func (c *GitCommand) Fetch(unamePassQuestion func(string) string, canAskForCredentials bool) error {
	return c.OSCommand.DetectUnamePass("git fetch", func(question string) string {
		if canAskForCredentials {
			return unamePassQuestion(question)
		}
		return "\n"
	})
}

// ResetToCommit reset to commit
func (c *GitCommand) ResetToCommit(sha string) error {
	return c.OSCommand.RunCommand(fmt.Sprintf("git reset %s", sha))
}

// NewBranch create new branch
func (c *GitCommand) NewBranch(name string) error {
	return c.OSCommand.RunCommand(fmt.Sprintf("git checkout -b %s", name))
}

// CurrentBranchName is a function.
func (c *GitCommand) CurrentBranchName() (string, error) {
	branchName, err := c.OSCommand.RunCommandWithOutput("git symbolic-ref --short HEAD")
	if err != nil {
		branchName, err = c.OSCommand.RunCommandWithOutput("git rev-parse --short HEAD")
		if err != nil {
			return "", err
		}
	}
	return utils.TrimTrailingNewline(branchName), nil
}

// DeleteBranch delete branch
func (c *GitCommand) DeleteBranch(branch string, force bool) error {
	command := "git branch -d"

	if force {
		command = "git branch -D"
	}

	return c.OSCommand.RunCommand(fmt.Sprintf("%s %s", command, branch))
}

// ListStash list stash
func (c *GitCommand) ListStash() (string, error) {
	return c.OSCommand.RunCommandWithOutput("git stash list")
}

// Merge merge
func (c *GitCommand) Merge(branchName string) error {
	return c.OSCommand.RunCommand(fmt.Sprintf("git merge --no-edit %s", branchName))
}

// AbortMerge abort merge
func (c *GitCommand) AbortMerge() error {
	return c.OSCommand.RunCommand("git merge --abort")
}

// usingGpg tells us whether the user has gpg enabled so that we can know
// whether we need to run a subprocess to allow them to enter their password
func (c *GitCommand) usingGpg() bool {
	gpgsign, _ := c.getLocalGitConfig("commit.gpgsign")
	if gpgsign == "" {
		gpgsign, _ = c.getGlobalGitConfig("commit.gpgsign")
	}
	value := strings.ToLower(gpgsign)

	return value == "true" || value == "1" || value == "yes" || value == "on"
}

// Commit commits to git
func (c *GitCommand) Commit(message string, amend bool) (*exec.Cmd, error) {
	amendParam := ""
	if amend {
		amendParam = " --amend"
	}
	command := fmt.Sprintf("git commit%s -m %s", amendParam, c.OSCommand.Quote(message))
	if c.usingGpg() {
		return c.OSCommand.PrepareSubProcess(c.OSCommand.Platform.shell, c.OSCommand.Platform.shellArg, command), nil
	}

	return nil, c.OSCommand.RunCommand(command)
}

// Pull pulls from repo
func (c *GitCommand) Pull(ask func(string) string) error {
	return c.OSCommand.DetectUnamePass("git pull --no-edit", ask)
}

// Push pushes to a branch
func (c *GitCommand) Push(branchName string, force bool, ask func(string) string) error {
	forceFlag := ""
	if force {
		forceFlag = "--force-with-lease "
	}

	cmd := fmt.Sprintf("git push %s-u origin %s", forceFlag, branchName)
	return c.OSCommand.DetectUnamePass(cmd, ask)
}

// SquashPreviousTwoCommits squashes a commit down to the one below it
// retaining the message of the higher commit
func (c *GitCommand) SquashPreviousTwoCommits(message string) error {
	// TODO: test this
	if err := c.OSCommand.RunCommand("git reset --soft HEAD^"); err != nil {
		return err
	}
	// TODO: if password is required, we need to return a subprocess
	return c.OSCommand.RunCommand(fmt.Sprintf("git commit --amend -m %s", c.OSCommand.Quote(message)))
}

// SquashFixupCommit squashes a 'FIXUP' commit into the commit beneath it,
// retaining the commit message of the lower commit
func (c *GitCommand) SquashFixupCommit(branchName string, shaValue string) error {
	commands := []string{
		fmt.Sprintf("git checkout -q %s", shaValue),
		fmt.Sprintf("git reset --soft %s^", shaValue),
		fmt.Sprintf("git commit --amend -C %s^", shaValue),
		fmt.Sprintf("git rebase --onto HEAD %s %s", shaValue, branchName),
	}
	for _, command := range commands {
		c.Log.Info(command)

		if output, err := c.OSCommand.RunCommandWithOutput(command); err != nil {
			ret := output
			// We are already in an error state here so we're just going to append
			// the output of these commands
			output, _ := c.OSCommand.RunCommandWithOutput(fmt.Sprintf("git branch -d %s", shaValue))
			ret += output
			output, _ = c.OSCommand.RunCommandWithOutput(fmt.Sprintf("git checkout %s", branchName))
			ret += output

			c.Log.Info(ret)
			return errors.New(ret)
		}
	}

	return nil
}

// CatFile obtains the content of a file
func (c *GitCommand) CatFile(fileName string) (string, error) {
	return c.OSCommand.RunCommandWithOutput(fmt.Sprintf("cat %s", c.OSCommand.Quote(fileName)))
}

// StageFile stages a file
func (c *GitCommand) StageFile(fileName string) error {
	return c.OSCommand.RunCommand(fmt.Sprintf("git add %s", c.OSCommand.Quote(fileName)))
}

// StageAll stages all files
func (c *GitCommand) StageAll() error {
	return c.OSCommand.RunCommand("git add -A")
}

// UnstageAll stages all files
func (c *GitCommand) UnstageAll() error {
	return c.OSCommand.RunCommand("git reset")
}

// UnStageFile unstages a file
func (c *GitCommand) UnStageFile(fileName string, tracked bool) error {
	command := "git rm --cached %s"
	if tracked {
		command = "git reset HEAD %s"
	}
	return c.OSCommand.RunCommand(fmt.Sprintf(command, c.OSCommand.Quote(fileName)))
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

// RebaseMode returns "" for non-rebase mode, "normal" for normal rebase
// and "interactive" for interactive rebase
func (c *GitCommand) RebaseMode() (string, error) {
	exists, err := c.OSCommand.FileExists(".git/rebase-apply")
	if err != nil {
		return "", err
	}
	if exists {
		return "normal", nil
	}
	exists, err = c.OSCommand.FileExists(".git/rebase-merge")
	if exists {
		return "interactive", err
	} else {
		return "", err
	}
}

// RemoveFile directly
func (c *GitCommand) RemoveFile(file *File) error {
	// if the file isn't tracked, we assume you want to delete it
	if file.HasStagedChanges {
		if err := c.OSCommand.RunCommand(fmt.Sprintf("git reset -- %s", file.Name)); err != nil {
			return err
		}
	}
	if !file.Tracked {
		return c.removeFile(file.Name)
	}
	// if the file is tracked, we assume you want to just check it out
	return c.OSCommand.RunCommand(fmt.Sprintf("git checkout -- %s", file.Name))
}

// Checkout checks out a branch, with --force if you set the force arg to true
func (c *GitCommand) Checkout(branch string, force bool) error {
	forceArg := ""
	if force {
		forceArg = "--force "
	}
	return c.OSCommand.RunCommand(fmt.Sprintf("git checkout %s %s", forceArg, branch))
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

// PrepareCommitAmendSubProcess prepares a subprocess for `git commit --amend --allow-empty`
func (c *GitCommand) PrepareCommitAmendSubProcess() *exec.Cmd {
	return c.OSCommand.PrepareSubProcess("git", "commit", "--amend", "--allow-empty")
}

// GetBranchGraph gets the color-formatted graph of the log for the given branch
// Currently it limits the result to 100 commits, but when we get async stuff
// working we can do lazy loading
func (c *GitCommand) GetBranchGraph(branchName string) (string, error) {
	return c.OSCommand.RunCommandWithOutput(fmt.Sprintf("git log --graph --color --abbrev-commit --decorate --date=relative --pretty=medium -100 %s", branchName))
}

func (c *GitCommand) getMergeBase() (string, error) {
	currentBranch, err := c.CurrentBranchName()
	if err != nil {
		return "", err
	}

	baseBranch := "master"
	if strings.HasPrefix(currentBranch, "feature/") {
		baseBranch = "develop"
	}

	output, err := c.OSCommand.RunCommandWithOutput(fmt.Sprintf("git merge-base HEAD %s", baseBranch))
	if err != nil {
		// swallowing error because it's not a big deal; probably because there are no commits yet
	}
	return output, nil
}

// GetRebasingCommits obtains the commits that we're in the process of rebasing
func (c *GitCommand) GetRebasingCommits() ([]*Commit, error) {
	rebaseMode, err := c.RebaseMode()
	if err != nil {
		return nil, err
	}
	switch rebaseMode {
	case "normal":
		return c.GetNormalRebasingCommits()
	case "interactive":
		return c.GetInteractiveRebasingCommits()
	default:
		return nil, nil
	}
}

func (c *GitCommand) GetNormalRebasingCommits() ([]*Commit, error) {
	rewrittenCount := 0
	bytesContent, err := ioutil.ReadFile(".git/rebase-apply/rewritten")
	if err == nil {
		content := string(bytesContent)
		rewrittenCount = len(strings.Split(content, "\n"))
	}

	// we know we're rebasing, so lets get all the files whose names have numbers
	commits := []*Commit{}
	err = filepath.Walk(".git/rebase-apply", func(path string, f os.FileInfo, err error) error {
		if rewrittenCount > 0 {
			rewrittenCount -= 1
			return nil
		}
		if err != nil {
			return err
		}
		re := regexp.MustCompile(`^\d+$`)
		if !re.MatchString(f.Name()) {
			return nil
		}
		bytesContent, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(bytesContent)
		commit, err := c.CommitFromPatch(content)
		if err != nil {
			return err
		}
		commits = append([]*Commit{commit}, commits...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return commits, nil
}

// git-rebase-todo example:
// pick ac446ae94ee560bdb8d1d057278657b251aaef17 ac446ae
// pick afb893148791a2fbd8091aeb81deba4930c73031 afb8931

// git-rebase-todo.backup example:
// pick 49cbba374296938ea86bbd4bf4fee2f6ba5cccf6 third commit on master
// pick ac446ae94ee560bdb8d1d057278657b251aaef17 blah  commit on master
// pick afb893148791a2fbd8091aeb81deba4930c73031 fourth commit on master

// GetInteractiveRebasingCommits takes our git-rebase-todo and our git-rebase-todo.backup files
// and extracts out the sha and names of commits that we still have to go
// in the rebase:
func (c *GitCommand) GetInteractiveRebasingCommits() ([]*Commit, error) {
	bytesContent, err := ioutil.ReadFile(".git/rebase-merge/git-rebase-todo")
	var content []string
	if err == nil {
		content = strings.Split(string(bytesContent), "\n")
		if len(content) > 0 && content[len(content)-1] == "" {
			content = content[0 : len(content)-1]
		}
	}

	// for each of them, grab the matching commit name in the backup
	bytesContent, err = ioutil.ReadFile(".git/rebase-merge/git-rebase-todo.backup")
	var backupContent []string
	if err == nil {
		backupContent = strings.Split(string(bytesContent), "\n")
	}

	commits := []*Commit{}
	for _, todoLine := range content {
		commit := c.extractCommit(todoLine, backupContent)
		if commit != nil {
			commits = append([]*Commit{commit}, commits...)
		}
	}

	return commits, nil
}

func (c *GitCommand) extractCommit(todoLine string, backupContent []string) *Commit {
	for _, backupLine := range backupContent {
		split := strings.Split(todoLine, " ")
		prefix := strings.Join(split[0:2], " ")
		if strings.HasPrefix(backupLine, prefix) {
			return &Commit{
				Sha:    split[2],
				Name:   strings.TrimPrefix(backupLine, prefix+" "),
				Status: "rebasing",
			}
		}
	}
	return nil
}

// assuming the file starts like this:
// From e93d4193e6dd45ca9cf3a5a273d7ba6cd8b8fb20 Mon Sep 17 00:00:00 2001
// From: Lazygit Tester <test@example.com>
// Date: Wed, 5 Dec 2018 21:03:23 +1100
// Subject: second commit on master
func (c *GitCommand) CommitFromPatch(content string) (*Commit, error) {
	lines := strings.Split(content, "\n")
	sha := strings.Split(lines[0], " ")[1][0:7]
	name := strings.TrimPrefix(lines[3], "Subject: ")
	return &Commit{
		Sha:    sha,
		Name:   name,
		Status: "rebasing",
	}, nil
}

// GetCommits obtains the commits of the current branch
func (c *GitCommand) GetCommits() ([]*Commit, error) {
	commits := []*Commit{}
	// here we want to also prepend the commits that we're in the process of rebasing
	rebasingCommits, err := c.GetRebasingCommits()
	if err != nil {
		return nil, err
	}
	if len(rebasingCommits) > 0 {
		commits = append(commits, rebasingCommits...)
	}

	unpushedCommits := c.GetUnpushedCommits()
	log := c.GetLog()

	// now we can split it up and turn it into commits
	for _, line := range utils.SplitLines(log) {
		splitLine := strings.Split(line, " ")
		sha := splitLine[0]
		_, unpushed := unpushedCommits[sha]
		status := map[bool]string{true: "unpushed", false: "pushed"}[unpushed]
		commits = append(commits, &Commit{
			Sha:           sha,
			Name:          strings.Join(splitLine[1:], " "),
			Status:        status,
			DisplayString: strings.Join(splitLine, " "),
		})
	}
	if len(rebasingCommits) > 0 {
		currentCommit := commits[len(rebasingCommits)]
		blue := color.New(color.FgYellow)
		youAreHere := blue.Sprint("<-- YOU ARE HERE ---")
		currentCommit.Name = fmt.Sprintf("%s %s", youAreHere, currentCommit.Name)
	}
	return c.setCommitMergedStatuses(commits)
}

func (c *GitCommand) setCommitMergedStatuses(commits []*Commit) ([]*Commit, error) {
	ancestor, err := c.getMergeBase()
	if err != nil {
		return nil, err
	}
	if ancestor == "" {
		return commits, nil
	}
	passedAncestor := false
	for i, commit := range commits {
		if strings.HasPrefix(ancestor, commit.Sha) {
			passedAncestor = true
		}
		if commit.Status != "pushed" {
			continue
		}
		if passedAncestor {
			commits[i].Status = "merged"
		}
	}
	return commits, nil
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
func (c *GitCommand) Show(sha string) (string, error) {
	show, err := c.OSCommand.RunCommandWithOutput(fmt.Sprintf("git show --color %s", sha))
	if err != nil {
		return "", err
	}

	// if this is a merge commit, we need to go a step further and get the diff between the two branches we merged
	revList, err := c.OSCommand.RunCommandWithOutput(fmt.Sprintf("git rev-list -1 --merges %s^...%s", sha, sha))
	if err != nil {
		// turns out we get an error here when it's the first commit. We'll just return the original show
		return show, nil
	}
	if len(revList) == 0 {
		return show, nil
	}

	// we want to pull out 1a6a69a and 3b51d7c from this:
	// commit ccc771d8b13d5b0d4635db4463556366470fd4f6
	// Merge: 1a6a69a 3b51d7c
	lines := utils.SplitLines(show)
	if len(lines) < 2 {
		return show, nil
	}

	secondLineWords := strings.Split(lines[1], " ")
	if len(secondLineWords) < 3 {
		return show, nil
	}

	mergeDiff, err := c.OSCommand.RunCommandWithOutput(fmt.Sprintf("git diff --color %s...%s", secondLineWords[1], secondLineWords[2]))
	if err != nil {
		return "", err
	}
	return show + mergeDiff, nil
}

// GetRemoteURL returns current repo remote url
func (c *GitCommand) GetRemoteURL() string {
	url, _ := c.OSCommand.RunCommandWithOutput("git config --get remote.origin.url")
	return utils.TrimTrailingNewline(url)
}

// CheckRemoteBranchExists Returns remote branch
func (c *GitCommand) CheckRemoteBranchExists(branch *Branch) bool {
	_, err := c.OSCommand.RunCommandWithOutput(fmt.Sprintf(
		"git show-ref --verify -- refs/remotes/origin/%s",
		branch.Name,
	))

	return err == nil
}

// Diff returns the diff of a file
func (c *GitCommand) Diff(file *File, plain bool) string {
	cachedArg := ""
	trackedArg := "--"
	colorArg := "--color"
	fileName := c.OSCommand.Quote(file.Name)
	if file.HasStagedChanges && !file.HasUnstagedChanges {
		cachedArg = "--cached"
	}
	if !file.Tracked && !file.HasStagedChanges {
		trackedArg = "--no-index /dev/null"
	}
	if plain {
		colorArg = ""
	}

	command := fmt.Sprintf("git diff %s %s %s %s", colorArg, cachedArg, trackedArg, fileName)

	// for now we assume an error means the file was deleted
	s, _ := c.OSCommand.RunCommandWithOutput(command)
	return s
}

func (c *GitCommand) ApplyPatch(patch string) (string, error) {
	filename, err := c.OSCommand.CreateTempFile("patch", patch)
	if err != nil {
		c.Log.Error(err)
		return "", err
	}

	defer func() { _ = c.OSCommand.RemoveFile(filename) }()

	return c.OSCommand.RunCommandWithOutput(fmt.Sprintf("git apply --cached %s", filename))
}

func (c *GitCommand) FastForward(branchName string) error {
	upstream := "origin" // hardcoding for now
	return c.OSCommand.RunCommand(fmt.Sprintf("git fetch %s %s:%s", upstream, branchName, branchName))
}

// GenericMerge takes a commandType of "merging" or "rebasing" and a command of "abort", "skip" or "continue"
// By default we skip the editor in the case where a commit will be made
func (c *GitCommand) GenericMerge(commandType string, command string) error {
	gitCommand := fmt.Sprintf("git %s %s --%s", c.OSCommand.Platform.skipEditorArg, commandType, command)
	return c.OSCommand.RunCommand(gitCommand)
}

func (c *GitCommand) RewordCommit(commits []*Commit, index int) (*exec.Cmd, error) {
	todo, err := c.GenerateGenericRebaseTodo(commits, index, "reword")
	if err != nil {
		return nil, err
	}

	return c.PrepareInteractiveRebaseCommand(commits[index+1].Sha, todo, true)
}

func (c *GitCommand) MoveCommitDown(commits []*Commit, index int) error {
	// we must ensure that we have at least two commits after the selected one
	if len(commits) <= index+2 {
		// assuming they aren't picking the bottom commit
		// TODO: support more than say 30 commits and ensure this logic is correct, and i18n
		return errors.New("Not enough room")
	}

	todo := ""
	orderedCommits := append(commits[0:index], commits[index+1], commits[index])
	for _, commit := range orderedCommits {
		todo = "pick " + commit.Sha + "\n" + todo
	}

	cmd, err := c.PrepareInteractiveRebaseCommand(commits[index+2].Sha, todo, true)
	if err != nil {
		return err
	}

	return c.OSCommand.RunPreparedCommand(cmd)
}

func (c *GitCommand) InteractiveRebase(commits []*Commit, index int, action string) error {
	todo, err := c.GenerateGenericRebaseTodo(commits, index, action)
	if err != nil {
		return err
	}

	autoStash := action != "edit"
	cmd, err := c.PrepareInteractiveRebaseCommand(commits[index+1].Sha, todo, autoStash)
	if err != nil {
		return err
	}

	return c.OSCommand.RunPreparedCommand(cmd)
}

func (c *GitCommand) PrepareInteractiveRebaseCommand(baseSha string, todo string, autoStash bool) (*exec.Cmd, error) {
	ex, err := os.Executable() // get the executable path for git to use
	if err != nil {
		ex = os.Args[0] // fallback to the first call argument if needed
	}

	debug := "FALSE"
	if c.OSCommand.Config.GetDebug() == true {
		debug = "TRUE"
	}

	// we do not want to autostash if we are editing

	splitCmd := str.ToArgv(fmt.Sprintf("git rebase --interactive --autostash %s", baseSha))

	cmd := exec.Command(splitCmd[0], splitCmd[1:]...)

	cmd.Env = os.Environ()
	cmd.Env = append(
		cmd.Env,
		"LAZYGIT_CONTEXT=INTERACTIVE_REBASE",
		"LAZYGIT_REBASE_TODO="+todo,
		"DEBUG="+debug,
		"LANG=en_US.UTF-8",   // Force using EN as language
		"LC_ALL=en_US.UTF-8", // Force using EN as language
		"GIT_SEQUENCE_EDITOR="+ex,
	)

	return cmd, nil
}

func (c *GitCommand) HardReset(baseSha string) error {
	return c.OSCommand.RunCommand("git reset --hard " + baseSha)
}

func (v *GitCommand) GenerateGenericRebaseTodo(commits []*Commit, index int, action string) (string, error) {
	if len(commits) <= index+1 {
		// assuming they aren't picking the bottom commit
		// TODO: support more than say 30 commits and ensure this logic is correct, and i18n
		return "", errors.New("You cannot interactive rebase onto the first commit")
	}

	todo := ""
	for i, commit := range commits[0 : index+1] {
		a := "pick"
		if i == index {
			a = action
		}
		todo = a + " " + commit.Sha + "\n" + todo
	}
	return todo, nil
}
