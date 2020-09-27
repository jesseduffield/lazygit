package commands

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mgutz/str"

	"github.com/go-errors/errors"

	gogit "github.com/go-git/go-git/v5"
	"github.com/jesseduffield/lazygit/pkg/commands/patch"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/env"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sirupsen/logrus"
	gitconfig "github.com/tcnksm/go-gitconfig"
)

// this takes something like:
// * (HEAD detached at 264fc6f5)
//	remotes
// and returns '264fc6f5' as the second match
const CurrentBranchNameRegex = `(?m)^\*.*?([^ ]*?)\)?$`

func verifyInGitRepo(runCmd func(string, ...interface{}) error) error {
	return runCmd("git status")
}

func navigateToRepoRootDirectory(stat func(string) (os.FileInfo, error), chdir func(string) error) error {
	gitDir := env.GetGitDirEnv()
	if gitDir != "" {
		// we've been given the git directory explicitly so no need to navigate to it
		_, err := stat(gitDir)
		if err != nil {
			return WrapError(err)
		}

		return nil
	}

	// we haven't been given the git dir explicitly so we assume it's in the current working directory as `.git/` (or an ancestor directory)

	for {
		_, err := stat(".git")

		if err == nil {
			return nil
		}

		if !os.IsNotExist(err) {
			return WrapError(err)
		}

		if err = chdir(".."); err != nil {
			return WrapError(err)
		}
	}
}

// resolvePath takes a path containing a symlink and returns the true path
func resolvePath(path string) (string, error) {
	l, err := os.Lstat(path)
	if err != nil {
		return "", err
	}

	if l.Mode()&os.ModeSymlink == 0 {
		return path, nil
	}

	return filepath.EvalSymlinks(path)
}

func setupRepository(openGitRepository func(string) (*gogit.Repository, error), sLocalize func(string) string) (*gogit.Repository, error) {
	unresolvedPath := env.GetGitDirEnv()
	if unresolvedPath == "" {
		var err error
		unresolvedPath, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	path, err := resolvePath(unresolvedPath)
	if err != nil {
		return nil, err
	}

	repository, err := openGitRepository(path)

	if err != nil {
		if strings.Contains(err.Error(), `unquoted '\' must be followed by new line`) {
			return nil, errors.New(sLocalize("GitconfigParseErr"))
		}

		return nil, err
	}

	return repository, err
}

// GitCommand is our main git interface
type GitCommand struct {
	Log                  *logrus.Entry
	OSCommand            *OSCommand
	Repo                 *gogit.Repository
	Tr                   *i18n.Localizer
	Config               config.AppConfigurer
	getGlobalGitConfig   func(string) (string, error)
	getLocalGitConfig    func(string) (string, error)
	removeFile           func(string) error
	DotGitDir            string
	onSuccessfulContinue func() error
	PatchManager         *patch.PatchManager

	// Push to current determines whether the user has configured to push to the remote branch of the same name as the current or not
	PushToCurrent bool
}

// NewGitCommand it runs git commands
func NewGitCommand(log *logrus.Entry, osCommand *OSCommand, tr *i18n.Localizer, config config.AppConfigurer) (*GitCommand, error) {
	var repo *gogit.Repository

	// see what our default push behaviour is
	output, err := osCommand.RunCommandWithOutput("git config --get push.default")
	pushToCurrent := false
	if err != nil {
		log.Errorf("error reading git config: %v", err)
	} else {
		pushToCurrent = strings.TrimSpace(output) == "current"
	}

	if err := verifyInGitRepo(osCommand.RunCommand); err != nil {
		return nil, err
	}

	if err := navigateToRepoRootDirectory(os.Stat, os.Chdir); err != nil {
		return nil, err
	}

	if repo, err = setupRepository(gogit.PlainOpen, tr.SLocalize); err != nil {
		return nil, err
	}

	dotGitDir, err := findDotGitDir(os.Stat, ioutil.ReadFile)
	if err != nil {
		return nil, err
	}

	gitCommand := &GitCommand{
		Log:                log,
		OSCommand:          osCommand,
		Tr:                 tr,
		Repo:               repo,
		Config:             config,
		getGlobalGitConfig: gitconfig.Global,
		getLocalGitConfig:  gitconfig.Local,
		removeFile:         os.RemoveAll,
		DotGitDir:          dotGitDir,
		PushToCurrent:      pushToCurrent,
	}

	gitCommand.PatchManager = patch.NewPatchManager(log, gitCommand.ApplyPatch, gitCommand.ShowFileDiff)

	return gitCommand, nil
}

func findDotGitDir(stat func(string) (os.FileInfo, error), readFile func(filename string) ([]byte, error)) (string, error) {
	if env.GetGitDirEnv() != "" {
		return env.GetGitDirEnv(), nil
	}

	f, err := stat(".git")
	if err != nil {
		return "", err
	}

	if f.IsDir() {
		return ".git", nil
	}

	fileBytes, err := readFile(".git")
	if err != nil {
		return "", err
	}
	fileContent := string(fileBytes)
	if !strings.HasPrefix(fileContent, "gitdir: ") {
		return "", errors.New(".git is a file which suggests we are in a submodule but the file's contents do not contain a gitdir pointing to the actual .git directory")
	}
	return strings.TrimSpace(strings.TrimPrefix(fileContent, "gitdir: ")), nil
}

func (c *GitCommand) getUnfilteredStashEntries() []*StashEntry {
	unescaped := "git stash list --pretty='%gs'"
	rawString, _ := c.OSCommand.RunCommandWithOutput(unescaped)
	stashEntries := []*StashEntry{}
	for i, line := range utils.SplitLines(rawString) {
		stashEntries = append(stashEntries, stashEntryFromLine(line, i))
	}
	return stashEntries
}

// GetStashEntries stash entries
func (c *GitCommand) GetStashEntries(filterPath string) []*StashEntry {
	if filterPath == "" {
		return c.getUnfilteredStashEntries()
	}

	unescaped := fmt.Sprintf("git stash list --name-only")
	rawString, err := c.OSCommand.RunCommandWithOutput(unescaped)
	if err != nil {
		return c.getUnfilteredStashEntries()
	}
	stashEntries := []*StashEntry{}
	var currentStashEntry *StashEntry
	lines := utils.SplitLines(rawString)
	isAStash := func(line string) bool { return strings.HasPrefix(line, "stash@{") }
	re := regexp.MustCompile(`stash@\{(\d+)\}`)

outer:
	for i := 0; i < len(lines); i++ {
		if !isAStash(lines[i]) {
			continue
		}
		match := re.FindStringSubmatch(lines[i])
		idx, err := strconv.Atoi(match[1])
		if err != nil {
			return c.getUnfilteredStashEntries()
		}
		currentStashEntry = stashEntryFromLine(lines[i], idx)
		for i+1 < len(lines) && !isAStash(lines[i+1]) {
			i++
			if lines[i] == filterPath {
				stashEntries = append(stashEntries, currentStashEntry)
				continue outer
			}
		}
	}
	return stashEntries
}

func stashEntryFromLine(line string, index int) *StashEntry {
	return &StashEntry{
		Name:  line,
		Index: index,
	}
}

// GetStashEntryDiff stash diff
func (c *GitCommand) ShowStashEntryCmdStr(index int) string {
	return fmt.Sprintf("git stash show -p --stat --color=%s stash@{%d}", c.colorArg(), index)
}

// GetStatusFiles git status files
type GetStatusFileOptions struct {
	NoRenames bool
}

func (c *GitCommand) GetConfigValue(key string) string {
	output, _ := c.OSCommand.RunCommandWithOutput("git config --get %s", key)
	// looks like this returns an error if there is no matching value which we're okay with
	return strings.TrimSpace(output)
}

func (c *GitCommand) GetSubmoduleNames() ([]string, error) {
	file, err := os.Open(".gitmodules")
	if err != nil {
		if err == os.ErrNotExist {
			return nil, nil
		}
		return nil, err
	}

	submoduleNames := []string{}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		re := regexp.MustCompile(`\[submodule "(.*)"\]`)
		matches := re.FindStringSubmatch(line)

		if len(matches) > 0 {
			submoduleNames = append(submoduleNames, matches[1])
		}
	}

	return submoduleNames, nil
}

func (c *GitCommand) GetStatusFiles(opts GetStatusFileOptions) []*File {
	// check if config wants us ignoring untracked files
	untrackedFilesSetting := c.GetConfigValue("status.showUntrackedFiles")

	if untrackedFilesSetting == "" {
		untrackedFilesSetting = "all"
	}
	untrackedFilesArg := fmt.Sprintf("--untracked-files=%s", untrackedFilesSetting)

	statusOutput, err := c.GitStatus(GitStatusOptions{NoRenames: opts.NoRenames, UntrackedFilesArg: untrackedFilesArg})
	if err != nil {
		c.Log.Error(err)
	}
	statusStrings := utils.SplitLines(statusOutput)
	files := []*File{}

	submoduleNames, err := c.GetSubmoduleNames()
	if err != nil {
		c.Log.Error(err)
	}

	for _, statusString := range statusStrings {
		if strings.HasPrefix(statusString, "warning") {
			c.Log.Warningf("warning when calling git status: %s", statusString)
			continue
		}
		change := statusString[0:2]
		stagedChange := change[0:1]
		unstagedChange := statusString[1:2]
		filename := c.OSCommand.Unquote(statusString[3:])
		untracked := utils.IncludesString([]string{"??", "A ", "AM"}, change)
		hasNoStagedChanges := utils.IncludesString([]string{" ", "U", "?"}, stagedChange)
		hasMergeConflicts := utils.IncludesString([]string{"DD", "AA", "UU", "AU", "UA", "UD", "DU"}, change)
		hasInlineMergeConflicts := utils.IncludesString([]string{"UU", "AA"}, change)
		isSubmodule := utils.IncludesString(submoduleNames, filename)

		file := &File{
			Name:                    filename,
			DisplayString:           statusString,
			HasStagedChanges:        !hasNoStagedChanges,
			HasUnstagedChanges:      unstagedChange != " ",
			Tracked:                 !untracked,
			Deleted:                 unstagedChange == "D" || stagedChange == "D",
			HasMergeConflicts:       hasMergeConflicts,
			HasInlineMergeConflicts: hasInlineMergeConflicts,
			Type:                    c.OSCommand.FileType(filename),
			ShortStatus:             change,
			IsSubmodule:             isSubmodule,
		}
		files = append(files, file)
	}
	return files
}

// StashDo modify stash
func (c *GitCommand) StashDo(index int, method string) error {
	return c.OSCommand.RunCommand("git stash %s stash@{%d}", method, index)
}

// StashSave save stash
// TODO: before calling this, check if there is anything to save
func (c *GitCommand) StashSave(message string) error {
	return c.OSCommand.RunCommand("git stash save %s", c.OSCommand.Quote(message))
}

// MergeStatusFiles merge status files
func (c *GitCommand) MergeStatusFiles(oldFiles, newFiles []*File, selectedFile *File) []*File {
	if len(oldFiles) == 0 {
		return newFiles
	}

	appendedIndexes := []int{}

	// retain position of files we already could see
	result := []*File{}
	for _, oldFile := range oldFiles {
		for newIndex, newFile := range newFiles {
			if includesInt(appendedIndexes, newIndex) {
				continue
			}
			// if we just staged B and in doing so created 'A -> B' and we are currently have oldFile: A and newFile: 'A -> B', we want to wait until we come across B so the our cursor isn't jumping anywhere
			waitForMatchingFile := selectedFile != nil && newFile.IsRename() && !selectedFile.IsRename() && newFile.Matches(selectedFile) && !oldFile.Matches(selectedFile)

			if oldFile.Matches(newFile) && !waitForMatchingFile {
				result = append(result, newFile)
				appendedIndexes = append(appendedIndexes, newIndex)
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
	if err := c.ResetHard("HEAD"); err != nil {
		return err
	}

	return c.RemoveUntrackedFiles()
}

func (c *GitCommand) GetCurrentBranchUpstreamDifferenceCount() (string, string) {
	return c.GetCommitDifferences("HEAD", "HEAD@{u}")
}

func (c *GitCommand) GetBranchUpstreamDifferenceCount(branchName string) (string, string) {
	return c.GetCommitDifferences(branchName, branchName+"@{u}")
}

// GetCommitDifferences checks how many pushables/pullables there are for the
// current branch
func (c *GitCommand) GetCommitDifferences(from, to string) (string, string) {
	command := "git rev-list %s..%s --count"
	pushableCount, err := c.OSCommand.RunCommandWithOutput(command, to, from)
	if err != nil {
		return "?", "?"
	}
	pullableCount, err := c.OSCommand.RunCommandWithOutput(command, from, to)
	if err != nil {
		return "?", "?"
	}
	return strings.TrimSpace(pushableCount), strings.TrimSpace(pullableCount)
}

// RenameCommit renames the topmost commit with the given name
func (c *GitCommand) RenameCommit(name string) error {
	return c.OSCommand.RunCommand("git commit --allow-empty --amend -m %s", c.OSCommand.Quote(name))
}

// RebaseBranch interactive rebases onto a branch
func (c *GitCommand) RebaseBranch(branchName string) error {
	cmd, err := c.PrepareInteractiveRebaseCommand(branchName, "", false)
	if err != nil {
		return err
	}

	return c.OSCommand.RunPreparedCommand(cmd)
}

type FetchOptions struct {
	PromptUserForCredential func(string) string
	RemoteName              string
	BranchName              string
}

// Fetch fetch git repo
func (c *GitCommand) Fetch(opts FetchOptions) error {
	command := "git fetch"

	if opts.RemoteName != "" {
		command = fmt.Sprintf("%s %s", command, opts.RemoteName)
	}
	if opts.BranchName != "" {
		command = fmt.Sprintf("%s %s", command, opts.BranchName)
	}

	return c.OSCommand.DetectUnamePass(command, func(question string) string {
		if opts.PromptUserForCredential != nil {
			return opts.PromptUserForCredential(question)
		}
		return "\n"
	})
}

// ResetToCommit reset to commit
func (c *GitCommand) ResetToCommit(sha string, strength string, options RunCommandOptions) error {
	return c.OSCommand.RunCommandWithOptions(fmt.Sprintf("git reset --%s %s", strength, sha), options)
}

// NewBranch create new branch
func (c *GitCommand) NewBranch(name string, base string) error {
	return c.OSCommand.RunCommand("git checkout -b %s %s", name, base)
}

// CurrentBranchName get the current branch name and displayname.
// the first returned string is the name and the second is the displayname
// e.g. name is 123asdf and displayname is '(HEAD detached at 123asdf)'
func (c *GitCommand) CurrentBranchName() (string, string, error) {
	branchName, err := c.OSCommand.RunCommandWithOutput("git symbolic-ref --short HEAD")
	if err == nil && branchName != "HEAD\n" {
		trimmedBranchName := strings.TrimSpace(branchName)
		return trimmedBranchName, trimmedBranchName, nil
	}
	output, err := c.OSCommand.RunCommandWithOutput("git branch --contains")
	if err != nil {
		return "", "", err
	}
	for _, line := range utils.SplitLines(output) {
		re := regexp.MustCompile(CurrentBranchNameRegex)
		match := re.FindStringSubmatch(line)
		if len(match) > 0 {
			branchName = match[1]
			displayBranchName := match[0][2:]
			return branchName, displayBranchName, nil
		}
	}
	return "HEAD", "HEAD", nil
}

// DeleteBranch delete branch
func (c *GitCommand) DeleteBranch(branch string, force bool) error {
	command := "git branch -d"

	if force {
		command = "git branch -D"
	}

	return c.OSCommand.RunCommand("%s %s", command, branch)
}

// ListStash list stash
func (c *GitCommand) ListStash() (string, error) {
	return c.OSCommand.RunCommandWithOutput("git stash list")
}

type MergeOpts struct {
	FastForwardOnly bool
}

// Merge merge
func (c *GitCommand) Merge(branchName string, opts MergeOpts) error {
	mergeArgs := c.Config.GetUserConfig().GetString("git.merging.args")

	command := fmt.Sprintf("git merge --no-edit %s %s", mergeArgs, branchName)
	if opts.FastForwardOnly {
		command = fmt.Sprintf("%s --ff-only", command)
	}

	return c.OSCommand.RunCommand(command)
}

// AbortMerge abort merge
func (c *GitCommand) AbortMerge() error {
	return c.OSCommand.RunCommand("git merge --abort")
}

// usingGpg tells us whether the user has gpg enabled so that we can know
// whether we need to run a subprocess to allow them to enter their password
func (c *GitCommand) usingGpg() bool {
	overrideGpg := c.Config.GetUserConfig().GetBool("git.overrideGpg")
	if overrideGpg {
		return false
	}

	gpgsign, _ := c.getLocalGitConfig("commit.gpgsign")
	if gpgsign == "" {
		gpgsign, _ = c.getGlobalGitConfig("commit.gpgsign")
	}
	value := strings.ToLower(gpgsign)

	return value == "true" || value == "1" || value == "yes" || value == "on"
}

// Commit commits to git
func (c *GitCommand) Commit(message string, flags string) (*exec.Cmd, error) {
	command := fmt.Sprintf("git commit %s -m %s", flags, strconv.Quote(message))
	if c.usingGpg() {
		return c.OSCommand.ShellCommandFromString(command), nil
	}

	return nil, c.OSCommand.RunCommand(command)
}

// Get the subject of the HEAD commit
func (c *GitCommand) GetHeadCommitMessage() (string, error) {
	cmdStr := "git log -1 --pretty=%s"
	message, err := c.OSCommand.RunCommandWithOutput(cmdStr)
	return strings.TrimSpace(message), err
}

func (c *GitCommand) GetCommitMessage(commitSha string) (string, error) {
	cmdStr := "git rev-list --format=%B --max-count=1 " + commitSha
	messageWithHeader, err := c.OSCommand.RunCommandWithOutput(cmdStr)
	message := strings.Join(strings.SplitAfter(messageWithHeader, "\n")[1:], "\n")
	return strings.TrimSpace(message), err
}

// AmendHead amends HEAD with whatever is staged in your working tree
func (c *GitCommand) AmendHead() (*exec.Cmd, error) {
	command := "git commit --amend --no-edit --allow-empty"
	if c.usingGpg() {
		return c.OSCommand.ShellCommandFromString(command), nil
	}

	return nil, c.OSCommand.RunCommand(command)
}

// Push pushes to a branch
func (c *GitCommand) Push(branchName string, force bool, upstream string, args string, promptUserForCredential func(string) string) error {
	forceFlag := ""
	if force {
		forceFlag = "--force-with-lease"
	}

	setUpstreamArg := ""
	if upstream != "" {
		setUpstreamArg = "--set-upstream " + upstream
	}

	cmd := fmt.Sprintf("git push --follow-tags %s %s %s", forceFlag, setUpstreamArg, args)
	return c.OSCommand.DetectUnamePass(cmd, promptUserForCredential)
}

// CatFile obtains the content of a file
func (c *GitCommand) CatFile(fileName string) (string, error) {
	return c.OSCommand.RunCommandWithOutput("%s %s", c.OSCommand.Platform.catCmd, c.OSCommand.Quote(fileName))
}

// StageFile stages a file
func (c *GitCommand) StageFile(fileName string) error {
	// renamed files look like "file1 -> file2"
	fileNames := strings.Split(fileName, " -> ")
	return c.OSCommand.RunCommand("git add %s", c.OSCommand.Quote(fileNames[len(fileNames)-1]))
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

	// renamed files look like "file1 -> file2"
	fileNames := strings.Split(fileName, " -> ")
	for _, name := range fileNames {
		if err := c.OSCommand.RunCommand(command, c.OSCommand.Quote(name)); err != nil {
			return err
		}
	}
	return nil
}

// GitStatus returns the plaintext short status of the repo
type GitStatusOptions struct {
	NoRenames         bool
	UntrackedFilesArg string
}

func (c *GitCommand) GitStatus(opts GitStatusOptions) (string, error) {
	noRenamesFlag := ""
	if opts.NoRenames {
		noRenamesFlag = "--no-renames"
	}

	return c.OSCommand.RunCommandWithOutput("git status %s --porcelain %s", opts.UntrackedFilesArg, noRenamesFlag)
}

// IsInMergeState states whether we are still mid-merge
func (c *GitCommand) IsInMergeState() (bool, error) {
	return c.OSCommand.FileExists(filepath.Join(c.DotGitDir, "MERGE_HEAD"))
}

// RebaseMode returns "" for non-rebase mode, "normal" for normal rebase
// and "interactive" for interactive rebase
func (c *GitCommand) RebaseMode() (string, error) {
	exists, err := c.OSCommand.FileExists(filepath.Join(c.DotGitDir, "rebase-apply"))
	if err != nil {
		return "", err
	}
	if exists {
		return "normal", nil
	}
	exists, err = c.OSCommand.FileExists(filepath.Join(c.DotGitDir, "rebase-merge"))
	if exists {
		return "interactive", err
	} else {
		return "", err
	}
}

func (c *GitCommand) BeforeAndAfterFileForRename(file *File) (*File, *File, error) {

	if !file.IsRename() {
		return nil, nil, errors.New("Expected renamed file")
	}

	// we've got a file that represents a rename from one file to another. Unfortunately
	// our File abstraction fails to consider this case, so here we will refetch
	// all files, passing the --no-renames flag and then recursively call the function
	// again for the before file and after file. At some point we should fix the abstraction itself

	split := strings.Split(file.Name, " -> ")
	filesWithoutRenames := c.GetStatusFiles(GetStatusFileOptions{NoRenames: true})
	var beforeFile *File
	var afterFile *File
	for _, f := range filesWithoutRenames {
		if f.Name == split[0] {
			beforeFile = f
		}
		if f.Name == split[1] {
			afterFile = f
		}
	}

	if beforeFile == nil || afterFile == nil {
		return nil, nil, errors.New("Could not find deleted file or new file for file rename")
	}

	if beforeFile.IsRename() || afterFile.IsRename() {
		// probably won't happen but we want to ensure we don't get an infinite loop
		return nil, nil, errors.New("Nested rename found")
	}

	return beforeFile, afterFile, nil
}

// DiscardAllFileChanges directly
func (c *GitCommand) DiscardAllFileChanges(file *File) error {
	if file.IsRename() {
		beforeFile, afterFile, err := c.BeforeAndAfterFileForRename(file)
		if err != nil {
			return err
		}

		if err := c.DiscardAllFileChanges(beforeFile); err != nil {
			return err
		}

		if err := c.DiscardAllFileChanges(afterFile); err != nil {
			return err
		}

		return nil
	}

	// if the file isn't tracked, we assume you want to delete it
	quotedFileName := c.OSCommand.Quote(file.Name)
	if file.IsSubmodule {
		if err := c.OSCommand.RunCommand(fmt.Sprintf("git submodule update --checkout --force --init %s", quotedFileName)); err != nil {
			return err
		}
	} else if file.HasStagedChanges || file.HasMergeConflicts {
		if err := c.OSCommand.RunCommand("git reset -- %s", quotedFileName); err != nil {
			return err
		}
	}

	if !file.Tracked {
		return c.removeFile(file.Name)
	}
	return c.DiscardUnstagedFileChanges(file)
}

// DiscardUnstagedFileChanges directly
func (c *GitCommand) DiscardUnstagedFileChanges(file *File) error {
	quotedFileName := c.OSCommand.Quote(file.Name)
	return c.OSCommand.RunCommand("git checkout -- %s", quotedFileName)
}

// Checkout checks out a branch (or commit), with --force if you set the force arg to true
type CheckoutOptions struct {
	Force   bool
	EnvVars []string
}

func (c *GitCommand) Checkout(branch string, options CheckoutOptions) error {
	forceArg := ""
	if options.Force {
		forceArg = "--force "
	}
	return c.OSCommand.RunCommandWithOptions(fmt.Sprintf("git checkout %s %s", forceArg, branch), RunCommandOptions{EnvVars: options.EnvVars})
}

// PrepareCommitAmendSubProcess prepares a subprocess for `git commit --amend --allow-empty`
func (c *GitCommand) PrepareCommitAmendSubProcess() *exec.Cmd {
	return c.OSCommand.PrepareSubProcess("git", "commit", "--amend", "--allow-empty")
}

// GetBranchGraph gets the color-formatted graph of the log for the given branch
// Currently it limits the result to 100 commits, but when we get async stuff
// working we can do lazy loading
func (c *GitCommand) GetBranchGraph(branchName string) (string, error) {
	cmdStr := c.GetBranchGraphCmdStr(branchName)
	return c.OSCommand.RunCommandWithOutput(cmdStr)
}

func (c *GitCommand) GetUpstreamForBranch(branchName string) (string, error) {
	output, err := c.OSCommand.RunCommandWithOutput("git rev-parse --abbrev-ref --symbolic-full-name %s@{u}", branchName)
	return strings.TrimSpace(output), err
}

// Ignore adds a file to the gitignore for the repo
func (c *GitCommand) Ignore(filename string) error {
	return c.OSCommand.AppendLineToFile(".gitignore", filename)
}

func (c *GitCommand) ShowCmdStr(sha string, filterPath string) string {
	filterPathArg := ""
	if filterPath != "" {
		filterPathArg = fmt.Sprintf(" -- %s", c.OSCommand.Quote(filterPath))
	}
	return fmt.Sprintf("git show --color=%s --no-renames --stat -p %s %s", c.colorArg(), sha, filterPathArg)
}

func (c *GitCommand) GetBranchGraphCmdStr(branchName string) string {
	branchLogCmdTemplate := c.Config.GetUserConfig().GetString("git.branchLogCmd")
	templateValues := map[string]string{
		"branchName": branchName,
	}
	return utils.ResolvePlaceholderString(branchLogCmdTemplate, templateValues)
}

// GetRemoteURL returns current repo remote url
func (c *GitCommand) GetRemoteURL() string {
	url, _ := c.OSCommand.RunCommandWithOutput("git config --get remote.origin.url")
	return utils.TrimTrailingNewline(url)
}

// CheckRemoteBranchExists Returns remote branch
func (c *GitCommand) CheckRemoteBranchExists(branch *Branch) bool {
	_, err := c.OSCommand.RunCommandWithOutput(
		"git show-ref --verify -- refs/remotes/origin/%s",
		branch.Name,
	)

	return err == nil
}

// WorktreeFileDiff returns the diff of a file
func (c *GitCommand) WorktreeFileDiff(file *File, plain bool, cached bool) string {
	// for now we assume an error means the file was deleted
	s, _ := c.OSCommand.RunCommandWithOutput(c.WorktreeFileDiffCmdStr(file, plain, cached))
	return s
}

func (c *GitCommand) WorktreeFileDiffCmdStr(file *File, plain bool, cached bool) string {
	cachedArg := ""
	trackedArg := "--"
	colorArg := c.colorArg()
	split := strings.Split(file.Name, " -> ") // in case of a renamed file we get the new filename
	fileName := c.OSCommand.Quote(split[len(split)-1])
	if cached {
		cachedArg = "--cached"
	}
	if !file.Tracked && !file.HasStagedChanges && !cached {
		trackedArg = "--no-index /dev/null"
	}
	if plain {
		colorArg = "never"
	}

	return fmt.Sprintf("git diff --submodule --no-ext-diff --color=%s %s %s %s", colorArg, cachedArg, trackedArg, fileName)
}

func (c *GitCommand) ApplyPatch(patch string, flags ...string) error {
	filepath := filepath.Join(c.Config.GetUserConfigDir(), utils.GetCurrentRepoName(), time.Now().Format("Jan _2 15.04.05.000000000")+".patch")
	c.Log.Infof("saving temporary patch to %s", filepath)
	if err := c.OSCommand.CreateFileWithContent(filepath, patch); err != nil {
		return err
	}

	flagStr := ""
	for _, flag := range flags {
		flagStr += " --" + flag
	}

	return c.OSCommand.RunCommand("git apply %s %s", flagStr, c.OSCommand.Quote(filepath))
}

func (c *GitCommand) FastForward(branchName string, remoteName string, remoteBranchName string, promptUserForCredential func(string) string) error {
	command := fmt.Sprintf("git fetch %s %s:%s", remoteName, remoteBranchName, branchName)
	return c.OSCommand.DetectUnamePass(command, promptUserForCredential)
}

func (c *GitCommand) RunSkipEditorCommand(command string) error {
	cmd := c.OSCommand.ExecutableFromString(command)
	lazyGitPath := c.OSCommand.GetLazygitPath()
	cmd.Env = append(
		cmd.Env,
		"LAZYGIT_CLIENT_COMMAND=EXIT_IMMEDIATELY",
		"GIT_EDITOR="+lazyGitPath,
		"EDITOR="+lazyGitPath,
		"VISUAL="+lazyGitPath,
	)
	return c.OSCommand.RunExecutable(cmd)
}

// GenericMerge takes a commandType of "merge" or "rebase" and a command of "abort", "skip" or "continue"
// By default we skip the editor in the case where a commit will be made
func (c *GitCommand) GenericMerge(commandType string, command string) error {
	err := c.RunSkipEditorCommand(
		fmt.Sprintf(
			"git %s --%s",
			commandType,
			command,
		),
	)
	if err != nil {
		if !strings.Contains(err.Error(), "no rebase in progress") {
			return err
		}
		c.Log.Warn(err)
	}

	// sometimes we need to do a sequence of things in a rebase but the user needs to
	// fix merge conflicts along the way. When this happens we queue up the next step
	// so that after the next successful rebase continue we can continue from where we left off
	if commandType == "rebase" && command == "continue" && c.onSuccessfulContinue != nil {
		f := c.onSuccessfulContinue
		c.onSuccessfulContinue = nil
		return f()
	}
	if command == "abort" {
		c.onSuccessfulContinue = nil
	}
	return nil
}

func (c *GitCommand) RewordCommit(commits []*Commit, index int) (*exec.Cmd, error) {
	todo, sha, err := c.GenerateGenericRebaseTodo(commits, index, "reword")
	if err != nil {
		return nil, err
	}

	return c.PrepareInteractiveRebaseCommand(sha, todo, false)
}

func (c *GitCommand) MoveCommitDown(commits []*Commit, index int) error {
	// we must ensure that we have at least two commits after the selected one
	if len(commits) <= index+2 {
		// assuming they aren't picking the bottom commit
		return errors.New(c.Tr.SLocalize("NoRoom"))
	}

	todo := ""
	orderedCommits := append(commits[0:index], commits[index+1], commits[index])
	for _, commit := range orderedCommits {
		todo = "pick " + commit.Sha + " " + commit.Name + "\n" + todo
	}

	cmd, err := c.PrepareInteractiveRebaseCommand(commits[index+2].Sha, todo, true)
	if err != nil {
		return err
	}

	return c.OSCommand.RunPreparedCommand(cmd)
}

func (c *GitCommand) InteractiveRebase(commits []*Commit, index int, action string) error {
	todo, sha, err := c.GenerateGenericRebaseTodo(commits, index, action)
	if err != nil {
		return err
	}

	cmd, err := c.PrepareInteractiveRebaseCommand(sha, todo, true)
	if err != nil {
		return err
	}

	return c.OSCommand.RunPreparedCommand(cmd)
}

// PrepareInteractiveRebaseCommand returns the cmd for an interactive rebase
// we tell git to run lazygit to edit the todo list, and we pass the client
// lazygit a todo string to write to the todo file
func (c *GitCommand) PrepareInteractiveRebaseCommand(baseSha string, todo string, overrideEditor bool) (*exec.Cmd, error) {
	ex := c.OSCommand.GetLazygitPath()

	debug := "FALSE"
	if c.OSCommand.Config.GetDebug() {
		debug = "TRUE"
	}

	cmdStr := fmt.Sprintf("git rebase --interactive --autostash --keep-empty %s", baseSha)
	c.Log.WithField("command", cmdStr).Info("RunCommand")
	splitCmd := str.ToArgv(cmdStr)

	cmd := c.OSCommand.command(splitCmd[0], splitCmd[1:]...)

	gitSequenceEditor := ex
	if todo == "" {
		gitSequenceEditor = "true"
	}

	cmd.Env = os.Environ()
	cmd.Env = append(
		cmd.Env,
		"LAZYGIT_CLIENT_COMMAND=INTERACTIVE_REBASE",
		"LAZYGIT_REBASE_TODO="+todo,
		"DEBUG="+debug,
		"LANG=en_US.UTF-8",   // Force using EN as language
		"LC_ALL=en_US.UTF-8", // Force using EN as language
		"GIT_SEQUENCE_EDITOR="+gitSequenceEditor,
	)

	if overrideEditor {
		cmd.Env = append(cmd.Env, "GIT_EDITOR="+ex)
	}

	return cmd, nil
}

func (c *GitCommand) HardReset(baseSha string) error {
	return c.OSCommand.RunCommand("git reset --hard " + baseSha)
}

func (c *GitCommand) SoftReset(baseSha string) error {
	return c.OSCommand.RunCommand("git reset --soft " + baseSha)
}

func (c *GitCommand) GenerateGenericRebaseTodo(commits []*Commit, actionIndex int, action string) (string, string, error) {
	baseIndex := actionIndex + 1

	if len(commits) <= baseIndex {
		return "", "", errors.New(c.Tr.SLocalize("CannotRebaseOntoFirstCommit"))
	}

	if action == "squash" || action == "fixup" {
		baseIndex++

		if len(commits) <= baseIndex {
			return "", "", errors.New(c.Tr.SLocalize("CannotSquashOntoSecondCommit"))
		}
	}

	todo := ""
	for i, commit := range commits[0:baseIndex] {
		var commitAction string
		if i == actionIndex {
			commitAction = action
		} else if commit.IsMerge {
			// your typical interactive rebase will actually drop merge commits by default. Damn git CLI, you scary!
			// doing this means we don't need to worry about rebasing over merges which always causes problems.
			// you typically shouldn't be doing rebases that pass over merge commits anyway.
			commitAction = "drop"
		} else {
			commitAction = "pick"
		}
		todo = commitAction + " " + commit.Sha + " " + commit.Name + "\n" + todo
	}

	return todo, commits[baseIndex].Sha, nil
}

// AmendTo amends the given commit with whatever files are staged
func (c *GitCommand) AmendTo(sha string) error {
	if err := c.CreateFixupCommit(sha); err != nil {
		return err
	}

	return c.SquashAllAboveFixupCommits(sha)
}

// EditRebaseTodo sets the action at a given index in the git-rebase-todo file
func (c *GitCommand) EditRebaseTodo(index int, action string) error {
	fileName := filepath.Join(c.DotGitDir, "rebase-merge/git-rebase-todo")
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	content := strings.Split(string(bytes), "\n")
	commitCount := c.getTodoCommitCount(content)

	// we have the most recent commit at the bottom whereas the todo file has
	// it at the bottom, so we need to subtract our index from the commit count
	contentIndex := commitCount - 1 - index
	splitLine := strings.Split(content[contentIndex], " ")
	content[contentIndex] = action + " " + strings.Join(splitLine[1:], " ")
	result := strings.Join(content, "\n")

	return ioutil.WriteFile(fileName, []byte(result), 0644)
}

func (c *GitCommand) getTodoCommitCount(content []string) int {
	// count lines that are not blank and are not comments
	commitCount := 0
	for _, line := range content {
		if line != "" && !strings.HasPrefix(line, "#") {
			commitCount++
		}
	}
	return commitCount
}

// MoveTodoDown moves a rebase todo item down by one position
func (c *GitCommand) MoveTodoDown(index int) error {
	fileName := filepath.Join(c.DotGitDir, "rebase-merge/git-rebase-todo")
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	content := strings.Split(string(bytes), "\n")
	commitCount := c.getTodoCommitCount(content)
	contentIndex := commitCount - 1 - index

	rearrangedContent := append(content[0:contentIndex-1], content[contentIndex], content[contentIndex-1])
	rearrangedContent = append(rearrangedContent, content[contentIndex+1:]...)
	result := strings.Join(rearrangedContent, "\n")

	return ioutil.WriteFile(fileName, []byte(result), 0644)
}

// Revert reverts the selected commit by sha
func (c *GitCommand) Revert(sha string) error {
	return c.OSCommand.RunCommand("git revert %s", sha)
}

// CherryPickCommits begins an interactive rebase with the given shas being cherry picked onto HEAD
func (c *GitCommand) CherryPickCommits(commits []*Commit) error {
	todo := ""
	for _, commit := range commits {
		todo = "pick " + commit.Sha + " " + commit.Name + "\n" + todo
	}

	cmd, err := c.PrepareInteractiveRebaseCommand("HEAD", todo, false)
	if err != nil {
		return err
	}

	return c.OSCommand.RunPreparedCommand(cmd)
}

// GetFilesInDiff get the specified commit files
func (c *GitCommand) GetFilesInDiff(from string, to string, reverse bool, patchManager *patch.PatchManager) ([]*CommitFile, error) {
	reverseFlag := ""
	if reverse {
		reverseFlag = " -R "
	}

	filenames, err := c.OSCommand.RunCommandWithOutput("git diff --submodule --no-ext-diff --name-status %s %s %s", reverseFlag, from, to)
	if err != nil {
		return nil, err
	}

	return c.GetCommitFilesFromFilenames(filenames, to, patchManager), nil
}

// filenames string is something like "file1\nfile2\nfile3"
func (c *GitCommand) GetCommitFilesFromFilenames(filenames string, parent string, patchManager *patch.PatchManager) []*CommitFile {
	commitFiles := make([]*CommitFile, 0)

	for _, line := range strings.Split(strings.TrimRight(filenames, "\n"), "\n") {
		// typical result looks like 'A my_file' meaning my_file was added
		if line == "" {
			continue
		}
		changeStatus := line[0:1]
		name := line[2:]
		status := patch.UNSELECTED
		if patchManager != nil && patchManager.To == parent {
			status = patchManager.GetFileStatus(name)
		}

		commitFiles = append(commitFiles, &CommitFile{
			Parent:       parent,
			Name:         name,
			ChangeStatus: changeStatus,
			PatchStatus:  status,
		})
	}

	return commitFiles
}

// ShowFileDiff get the diff of specified from and to. Typically this will be used for a single commit so it'll be 123abc^..123abc
// but when we're in diff mode it could be any 'from' to any 'to'. The reverse flag is also here thanks to diff mode.
func (c *GitCommand) ShowFileDiff(from string, to string, reverse bool, fileName string, plain bool) (string, error) {
	cmdStr := c.ShowFileDiffCmdStr(from, to, reverse, fileName, plain)
	return c.OSCommand.RunCommandWithOutput(cmdStr)
}

func (c *GitCommand) ShowFileDiffCmdStr(from string, to string, reverse bool, fileName string, plain bool) string {
	colorArg := c.colorArg()
	if plain {
		colorArg = "never"
	}

	reverseFlag := ""
	if reverse {
		reverseFlag = " -R "
	}

	return fmt.Sprintf("git diff --submodule --no-ext-diff --no-renames --color=%s %s %s %s -- %s", colorArg, from, to, reverseFlag, fileName)
}

// CheckoutFile checks out the file for the given commit
func (c *GitCommand) CheckoutFile(commitSha, fileName string) error {
	return c.OSCommand.RunCommand("git checkout %s %s", commitSha, fileName)
}

// DiscardOldFileChanges discards changes to a file from an old commit
func (c *GitCommand) DiscardOldFileChanges(commits []*Commit, commitIndex int, fileName string) error {
	if err := c.BeginInteractiveRebaseForCommit(commits, commitIndex); err != nil {
		return err
	}

	// check if file exists in previous commit (this command returns an error if the file doesn't exist)
	if err := c.OSCommand.RunCommand("git cat-file -e HEAD^:%s", fileName); err != nil {
		if err := c.OSCommand.Remove(fileName); err != nil {
			return err
		}
		if err := c.StageFile(fileName); err != nil {
			return err
		}
	} else if err := c.CheckoutFile("HEAD^", fileName); err != nil {
		return err
	}

	// amend the commit
	cmd, err := c.AmendHead()
	if cmd != nil {
		return errors.New("received unexpected pointer to cmd")
	}
	if err != nil {
		return err
	}

	// continue
	return c.GenericMerge("rebase", "continue")
}

// DiscardAnyUnstagedFileChanges discards any unstages file changes via `git checkout -- .`
func (c *GitCommand) DiscardAnyUnstagedFileChanges() error {
	return c.OSCommand.RunCommand("git checkout -- .")
}

// RemoveTrackedFiles will delete the given file(s) even if they are currently tracked
func (c *GitCommand) RemoveTrackedFiles(name string) error {
	return c.OSCommand.RunCommand("git rm -r --cached %s", name)
}

// RemoveUntrackedFiles runs `git clean -fd`
func (c *GitCommand) RemoveUntrackedFiles() error {
	return c.OSCommand.RunCommand("git clean -fd")
}

// ResetHardHead runs `git reset --hard`
func (c *GitCommand) ResetHard(ref string) error {
	return c.OSCommand.RunCommand("git reset --hard " + ref)
}

// ResetSoft runs `git reset --soft HEAD`
func (c *GitCommand) ResetSoft(ref string) error {
	return c.OSCommand.RunCommand("git reset --soft " + ref)
}

// CreateFixupCommit creates a commit that fixes up a previous commit
func (c *GitCommand) CreateFixupCommit(sha string) error {
	return c.OSCommand.RunCommand("git commit --fixup=%s", sha)
}

// SquashAllAboveFixupCommits squashes all fixup! commits above the given one
func (c *GitCommand) SquashAllAboveFixupCommits(sha string) error {
	return c.RunSkipEditorCommand(
		fmt.Sprintf(
			"git rebase --interactive --autostash --autosquash %s^",
			sha,
		),
	)
}

// StashSaveStagedChanges stashes only the currently staged changes. This takes a few steps
// shoutouts to Joe on https://stackoverflow.com/questions/14759748/stashing-only-staged-changes-in-git-is-it-possible
func (c *GitCommand) StashSaveStagedChanges(message string) error {

	if err := c.OSCommand.RunCommand("git stash --keep-index"); err != nil {
		return err
	}

	if err := c.StashSave(message); err != nil {
		return err
	}

	if err := c.OSCommand.RunCommand("git stash apply stash@{1}"); err != nil {
		return err
	}

	if err := c.OSCommand.PipeCommands("git stash show -p", "git apply -R"); err != nil {
		return err
	}

	if err := c.OSCommand.RunCommand("git stash drop stash@{1}"); err != nil {
		return err
	}

	// if you had staged an untracked file, that will now appear as 'AD' in git status
	// meaning it's deleted in your working tree but added in your index. Given that it's
	// now safely stashed, we need to remove it.
	files := c.GetStatusFiles(GetStatusFileOptions{})
	for _, file := range files {
		if file.ShortStatus == "AD" {
			if err := c.UnStageFile(file.Name, false); err != nil {
				return err
			}
		}
	}

	return nil
}

// BeginInteractiveRebaseForCommit starts an interactive rebase to edit the current
// commit and pick all others. After this you'll want to call `c.GenericMerge("rebase", "continue")`
func (c *GitCommand) BeginInteractiveRebaseForCommit(commits []*Commit, commitIndex int) error {
	if len(commits)-1 < commitIndex {
		return errors.New("index outside of range of commits")
	}

	// we can make this GPG thing possible it just means we need to do this in two parts:
	// one where we handle the possibility of a credential request, and the other
	// where we continue the rebase
	if c.usingGpg() {
		return errors.New(c.Tr.SLocalize("DisabledForGPG"))
	}

	todo, sha, err := c.GenerateGenericRebaseTodo(commits, commitIndex, "edit")
	if err != nil {
		return err
	}

	cmd, err := c.PrepareInteractiveRebaseCommand(sha, todo, true)
	if err != nil {
		return err
	}

	if err := c.OSCommand.RunPreparedCommand(cmd); err != nil {
		return err
	}

	return nil
}

func (c *GitCommand) SetUpstreamBranch(upstream string) error {
	return c.OSCommand.RunCommand("git branch -u %s", upstream)
}

func (c *GitCommand) AddRemote(name string, url string) error {
	return c.OSCommand.RunCommand("git remote add %s %s", name, url)
}

func (c *GitCommand) RemoveRemote(name string) error {
	return c.OSCommand.RunCommand("git remote remove %s", name)
}

func (c *GitCommand) IsHeadDetached() bool {
	err := c.OSCommand.RunCommand("git symbolic-ref -q HEAD")
	return err != nil
}

func (c *GitCommand) DeleteRemoteBranch(remoteName string, branchName string) error {
	return c.OSCommand.RunCommand("git push %s --delete %s", remoteName, branchName)
}

func (c *GitCommand) SetBranchUpstream(remoteName string, remoteBranchName string, branchName string) error {
	return c.OSCommand.RunCommand("git branch --set-upstream-to=%s/%s %s", remoteName, remoteBranchName, branchName)
}

func (c *GitCommand) RenameRemote(oldRemoteName string, newRemoteName string) error {
	return c.OSCommand.RunCommand("git remote rename %s %s", oldRemoteName, newRemoteName)
}

func (c *GitCommand) UpdateRemoteUrl(remoteName string, updatedUrl string) error {
	return c.OSCommand.RunCommand("git remote set-url %s %s", remoteName, updatedUrl)
}

func (c *GitCommand) CreateLightweightTag(tagName string, commitSha string) error {
	return c.OSCommand.RunCommand("git tag %s %s", tagName, commitSha)
}

func (c *GitCommand) DeleteTag(tagName string) error {
	return c.OSCommand.RunCommand("git tag -d %s", tagName)
}

func (c *GitCommand) PushTag(remoteName string, tagName string) error {
	return c.OSCommand.RunCommand("git push %s %s", remoteName, tagName)
}

func (c *GitCommand) FetchRemote(remoteName string) error {
	return c.OSCommand.RunCommand("git fetch %s", remoteName)
}

// GetReflogCommits only returns the new reflog commits since the given lastReflogCommit
// if none is passed (i.e. it's value is nil) then we get all the reflog commits

func (c *GitCommand) GetReflogCommits(lastReflogCommit *Commit, filterPath string) ([]*Commit, bool, error) {
	commits := make([]*Commit, 0)
	re := regexp.MustCompile(`(\w+).*HEAD@\{([^\}]+)\}: (.*)`)

	filterPathArg := ""
	if filterPath != "" {
		filterPathArg = fmt.Sprintf(" --follow -- %s", c.OSCommand.Quote(filterPath))
	}

	cmd := c.OSCommand.ExecutableFromString(fmt.Sprintf("git reflog --abbrev=20 --date=unix %s", filterPathArg))
	onlyObtainedNewReflogCommits := false
	err := RunLineOutputCmd(cmd, func(line string) (bool, error) {
		match := re.FindStringSubmatch(line)
		if len(match) <= 1 {
			return false, nil
		}

		unixTimestamp, _ := strconv.Atoi(match[2])

		commit := &Commit{
			Sha:           match[1],
			Name:          match[3],
			UnixTimestamp: int64(unixTimestamp),
			Status:        "reflog",
		}

		if lastReflogCommit != nil && commit.Sha == lastReflogCommit.Sha && commit.UnixTimestamp == lastReflogCommit.UnixTimestamp {
			onlyObtainedNewReflogCommits = true
			// after this point we already have these reflogs loaded so we'll simply return the new ones
			return true, nil
		}

		commits = append(commits, commit)
		return false, nil
	})
	if err != nil {
		return nil, false, err
	}

	return commits, onlyObtainedNewReflogCommits, nil
}

func (c *GitCommand) ConfiguredPager() string {
	if os.Getenv("GIT_PAGER") != "" {
		return os.Getenv("GIT_PAGER")
	}
	if os.Getenv("PAGER") != "" {
		return os.Getenv("PAGER")
	}
	output, err := c.OSCommand.RunCommandWithOutput("git config --get-all core.pager")
	if err != nil {
		return ""
	}
	trimmedOutput := strings.TrimSpace(output)
	return strings.Split(trimmedOutput, "\n")[0]
}

func (c *GitCommand) GetPager(width int) string {
	useConfig := c.Config.GetUserConfig().GetBool("git.paging.useConfig")
	if useConfig {
		pager := c.ConfiguredPager()
		return strings.Split(pager, "| less")[0]
	}

	templateValues := map[string]string{
		"columnWidth": strconv.Itoa(width/2 - 6),
	}

	pagerTemplate := c.Config.GetUserConfig().GetString("git.paging.pager")
	return utils.ResolvePlaceholderString(pagerTemplate, templateValues)
}

func (c *GitCommand) colorArg() string {
	return c.Config.GetUserConfig().GetString("git.paging.colorArg")
}

func (c *GitCommand) RenameBranch(oldName string, newName string) error {
	return c.OSCommand.RunCommand("git branch --move %s %s", oldName, newName)
}

func (c *GitCommand) WorkingTreeState() string {
	rebaseMode, _ := c.RebaseMode()
	if rebaseMode != "" {
		return "rebasing"
	}
	merging, _ := c.IsInMergeState()
	if merging {
		return "merging"
	}
	return "normal"
}

func (c *GitCommand) IsBareRepo() bool {
	// note: could use `git rev-parse --is-bare-repository` if we wanna drop go-git
	_, err := c.Repo.Worktree()
	return err == gogit.ErrIsBareRepository
}
