package main

import (

	// "log"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jesseduffield/gocui"
	gitconfig "github.com/tcnksm/go-gitconfig"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var (
	// ErrNoOpenCommand : When we don't know which command to use to open a file
	ErrNoOpenCommand = errors.New("Unsure what command to use to open this file")
)

// GitFile : A staged/unstaged file
// TODO: decide whether to give all of these the Git prefix
type GitFile struct {
	Name               string
	HasStagedChanges   bool
	HasUnstagedChanges bool
	Tracked            bool
	Deleted            bool
	HasMergeConflicts  bool
	DisplayString      string
}

// Commit : A git commit
type Commit struct {
	Sha           string
	Name          string
	Pushed        bool
	DisplayString string
}

// StashEntry : A git stash entry
type StashEntry struct {
	Index         int
	Name          string
	DisplayString string
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

func mergeGitStatusFiles(oldGitFiles, newGitFiles []GitFile) []GitFile {
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

func runDirectCommand(command string) (string, error) {
	timeStart := time.Now()
	commandLog(command)

	cmdOut, err := exec.
		Command(state.Platform.shell, state.Platform.shellArg, command).
		CombinedOutput()
	devLog("run direct command time for command: ", command, time.Now().Sub(timeStart))
	return sanitisedCommandOutput(cmdOut, err)
}

func branchStringParts(branchString string) (string, string) {
	// expect string to be something like '4w    master`
	splitBranchName := strings.Split(branchString, "\t")
	// if we have no \t then we have no recency, so just output that as blank
	if len(splitBranchName) == 1 {
		return "", branchString
	}
	return splitBranchName[0], splitBranchName[1]
}

// TODO: DRY up this function and getGitBranches
func getGitStashEntries() []StashEntry {
	stashEntries := make([]StashEntry, 0)
	rawString, _ := runDirectCommand("git stash list --pretty='%gs'")
	for i, line := range splitLines(rawString) {
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

func getStashEntryDiff(index int) (string, error) {
	return runCommand("git stash show -p --color stash@{" + fmt.Sprint(index) + "}")
}

func includes(array []string, str string) bool {
	for _, arrayStr := range array {
		if arrayStr == str {
			return true
		}
	}
	return false
}

func getGitStatusFiles() []GitFile {
	statusOutput, _ := getGitStatus()
	statusStrings := splitLines(statusOutput)
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
	devLog(gitFiles)
	return gitFiles
}

func gitStashDo(index int, method string) (string, error) {
	return runCommand("git stash " + method + " stash@{" + fmt.Sprint(index) + "}")
}

func gitStashSave(message string) (string, error) {
	output, err := runCommand("git stash save \"" + message + "\"")
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

func gitCheckout(branch string, force bool) (string, error) {
	forceArg := ""
	if force {
		forceArg = "--force "
	}
	return runCommand("git checkout " + forceArg + branch)
}

func sanitisedCommandOutput(output []byte, err error) (string, error) {
	outputString := string(output)
	if outputString == "" && err != nil {
		return err.Error(), err
	}
	return outputString, err
}

func runCommand(command string) (string, error) {
	commandStartTime := time.Now()
	commandLog(command)
	splitCmd := strings.Split(command, " ")
	devLog(splitCmd)
	cmdOut, err := exec.Command(splitCmd[0], splitCmd[1:]...).CombinedOutput()
	devLog("run command time: ", time.Now().Sub(commandStartTime))
	return sanitisedCommandOutput(cmdOut, err)
}

func vsCodeOpenFile(g *gocui.Gui, filename string) (string, error) {
	return runCommand("code -r " + filename)
}

func sublimeOpenFile(g *gocui.Gui, filename string) (string, error) {
	return runCommand("subl " + filename)
}

func openFile(g *gocui.Gui, filename string) (string, error) {
	cmdName, cmdTrail, err := getOpenCommand()
	if err != nil {
		return "", err
	}
	return runCommand(cmdName + " " + filename + cmdTrail)
}

func getOpenCommand() (string, string, error) {
	//NextStep open equivalents: xdg-open (linux), cygstart (cygwin), open (OSX)
	trailMap := map[string]string{
		"xdg-open": " &>/dev/null &",
		"cygstart": "",
		"open":     "",
	}
	for name, trail := range trailMap {
		if out, _ := runCommand("which " + name); out != "exit status 1" {
			return name, trail, nil
		}
	}
	return "", "", ErrNoOpenCommand
}

func gitAddPatch(g *gocui.Gui, filename string) {
	runSubProcess(g, "git", "add", "--patch", filename)
}

func editFile(g *gocui.Gui, filename string) (string, error) {
	editor, _ := gitconfig.Global("core.editor")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		return "", createErrorPanel(g, "No editor defined in $VISUAL, $EDITOR, or git config.")
	}
	runSubProcess(g, editor, filename)
	return "", nil
}

func runSubProcess(g *gocui.Gui, cmdName string, commandArgs ...string) {
	subprocess = exec.Command(cmdName, commandArgs...)
	subprocess.Stdin = os.Stdin
	subprocess.Stdout = os.Stdout
	subprocess.Stderr = os.Stderr

	g.Update(func(g *gocui.Gui) error {
		return ErrSubprocess
	})
}

func getBranchGraph(branch string) (string, error) {
	return runCommand("git log --graph --color --abbrev-commit --decorate --date=relative --pretty=medium -100 " + branch)
}

func verifyInGitRepo() {
	if output, err := runCommand("git status"); err != nil {
		fmt.Println(output)
		os.Exit(1)
	}
}

func getCommits() []Commit {
	pushables := gitCommitsToPush()
	log := getLog()
	commits := make([]Commit, 0)
	// now we can split it up and turn it into commits
	lines := splitLines(log)
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

func getLog() string {
	// currently limiting to 30 for performance reasons
	// TODO: add lazyloading when you scroll down
	result, err := runDirectCommand("git log --oneline -30")
	if err != nil {
		// assume if there is an error there are no commits yet for this branch
		return ""
	}
	return result
}

func gitIgnore(filename string) {
	if _, err := runDirectCommand("echo '" + filename + "' >> .gitignore"); err != nil {
		panic(err)
	}
}

func gitShow(sha string) string {
	result, err := runDirectCommand("git show --color " + sha)
	if err != nil {
		panic(err)
	}
	return result
}

func getDiff(file GitFile) string {
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
	s, _ := runCommand(command)
	return s
}

func catFile(file string) (string, error) {
	return runDirectCommand("cat " + file)
}

func stageFile(file string) error {
	_, err := runCommand("git add " + file)
	return err
}

func unStageFile(file string, tracked bool) error {
	var command string
	if tracked {
		command = "git reset HEAD "
	} else {
		command = "git rm --cached "
	}
	devLog(command)
	_, err := runCommand(command + file)
	return err
}

func getGitStatus() (string, error) {
	return runCommand("git status --untracked-files=all --short")
}

func isInMergeState() (bool, error) {
	output, err := runCommand("git status --untracked-files=all")
	if err != nil {
		return false, err
	}
	return strings.Contains(output, "conclude merge") || strings.Contains(output, "unmerged paths"), nil
}

func removeFile(file GitFile) error {
	// if the file isn't tracked, we assume you want to delete it
	if !file.Tracked {
		_, err := runCommand("rm -rf ./" + file.Name)
		return err
	}
	// if the file is tracked, we assume you want to just check it out
	_, err := runCommand("git checkout " + file.Name)
	return err
}

func gitCommit(g *gocui.Gui, message string) (string, error) {
	gpgsign, _ := gitconfig.Global("commit.gpgsign")
	if gpgsign != "" {
		runSubProcess(g, "bash", "-c", "git commit -m \""+message+"\"")
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

func gitRenameCommit(message string) (string, error) {
	return runDirectCommand("git commit --allow-empty --amend -m \"" + message + "\"")
}

func gitFetch() (string, error) {
	return runDirectCommand("git fetch")
}

func gitResetToCommit(sha string) (string, error) {
	return runDirectCommand("git reset " + sha)
}

func gitNewBranch(name string) (string, error) {
	return runDirectCommand("git checkout -b " + name)
}

func gitDeleteBranch(branch string) (string, error) {
        return runCommand("git branch -d " + branch)
}

func gitListStash() (string, error) {
	return runDirectCommand("git stash list")
}

func gitMerge(branchName string) (string, error) {
	return runDirectCommand("git merge --no-edit " + branchName)
}

func gitAbortMerge() (string, error) {
	return runDirectCommand("git merge --abort")
}

func gitUpstreamDifferenceCount() (string, string) {
	pushableCount, err := runDirectCommand("git rev-list @{u}..head --count")
	if err != nil {
		return "?", "?"
	}
	pullableCount, err := runDirectCommand("git rev-list head..@{u} --count")
	if err != nil {
		return "?", "?"
	}
	return strings.TrimSpace(pushableCount), strings.TrimSpace(pullableCount)
}

func gitCommitsToPush() []string {
	pushables, err := runDirectCommand("git rev-list @{u}..head --abbrev-commit")
	if err != nil {
		return make([]string, 0)
	}
	return splitLines(pushables)
}

func getGitBranches() []Branch {
	builder := newBranchListBuilder()
	return builder.build()
}

func branchIncluded(branchName string, branches []Branch) bool {
	for _, existingBranch := range branches {
		if strings.ToLower(existingBranch.Name) == strings.ToLower(branchName) {
			return true
		}
	}
	return false
}

func gitResetHard() error {
	return w.Reset(&git.ResetOptions{Mode: git.HardReset})
}
