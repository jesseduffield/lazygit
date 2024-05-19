package components

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// this is for running shell commands, mostly for the sake of setting up the repo
// but you can also run the commands from within lazygit to emulate things happening
// in the background.
type Shell struct {
	// working directory the shell is invoked in
	dir string
	// passed into each command
	env []string

	// when running the shell outside the gui we can directly panic on failure,
	// but inside the gui we need to close the gui before panicking
	fail func(string)

	randomFileContentIndex int
}

func NewShell(dir string, env []string, fail func(string)) *Shell {
	return &Shell{dir: dir, env: env, fail: fail}
}

func (self *Shell) RunCommand(args []string) *Shell {
	return self.RunCommandWithEnv(args, []string{})
}

// Run a command with additional environment variables set
func (self *Shell) RunCommandWithEnv(args []string, env []string) *Shell {
	output, err := self.runCommandWithOutputAndEnv(args, env)
	if err != nil {
		self.fail(fmt.Sprintf("error running command: %v\n%s", args, output))
	}

	return self
}

func (self *Shell) RunCommandExpectError(args []string) *Shell {
	output, err := self.runCommandWithOutput(args)
	if err == nil {
		self.fail(fmt.Sprintf("Expected error running shell command: %v\n%s", args, output))
	}

	return self
}

func (self *Shell) runCommandWithOutput(args []string) (string, error) {
	return self.runCommandWithOutputAndEnv(args, []string{})
}

func (self *Shell) runCommandWithOutputAndEnv(args []string, env []string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = append(self.env, env...)
	cmd.Dir = self.dir

	output, err := cmd.CombinedOutput()

	return string(output), err
}

func (self *Shell) RunShellCommand(cmdStr string) *Shell {
	shell := "sh"
	shellArg := "-c"
	if runtime.GOOS == "windows" {
		shell = "cmd"
		shellArg = "/C"
	}

	cmd := exec.Command(shell, shellArg, cmdStr)
	cmd.Env = os.Environ()
	cmd.Dir = self.dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		self.fail(fmt.Sprintf("error running shell command: %s\n%s", cmdStr, string(output)))
	}

	return self
}

func (self *Shell) CreateFile(path string, content string) *Shell {
	fullPath := filepath.Join(self.dir, path)

	// create any required directories
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		self.fail(fmt.Sprintf("error creating directory: %s\n%s", dir, err))
	}

	err := os.WriteFile(fullPath, []byte(content), 0o644)
	if err != nil {
		self.fail(fmt.Sprintf("error creating file: %s\n%s", fullPath, err))
	}

	return self
}

func (self *Shell) DeleteFile(path string) *Shell {
	fullPath := filepath.Join(self.dir, path)
	err := os.RemoveAll(fullPath)
	if err != nil {
		self.fail(fmt.Sprintf("error deleting file: %s\n%s", fullPath, err))
	}

	return self
}

func (self *Shell) CreateDir(path string) *Shell {
	fullPath := filepath.Join(self.dir, path)
	if err := os.MkdirAll(fullPath, 0o755); err != nil {
		self.fail(fmt.Sprintf("error creating directory: %s\n%s", fullPath, err))
	}

	return self
}

func (self *Shell) UpdateFile(path string, content string) *Shell {
	fullPath := filepath.Join(self.dir, path)
	err := os.WriteFile(fullPath, []byte(content), 0o644)
	if err != nil {
		self.fail(fmt.Sprintf("error updating file: %s\n%s", fullPath, err))
	}

	return self
}

func (self *Shell) NewBranch(name string) *Shell {
	return self.RunCommand([]string{"git", "checkout", "-b", name})
}

func (self *Shell) NewBranchFrom(name string, from string) *Shell {
	return self.RunCommand([]string{"git", "checkout", "-b", name, from})
}

func (self *Shell) Checkout(name string) *Shell {
	return self.RunCommand([]string{"git", "checkout", name})
}

func (self *Shell) Merge(name string) *Shell {
	return self.RunCommand([]string{"git", "merge", "--commit", "--no-ff", name})
}

func (self *Shell) ContinueMerge() *Shell {
	return self.RunCommand([]string{"git", "-c", "core.editor=true", "merge", "--continue"})
}

func (self *Shell) GitAdd(path string) *Shell {
	return self.RunCommand([]string{"git", "add", path})
}

func (self *Shell) GitAddAll() *Shell {
	return self.RunCommand([]string{"git", "add", "-A"})
}

func (self *Shell) Commit(message string) *Shell {
	return self.RunCommand([]string{"git", "commit", "-m", message})
}

func (self *Shell) EmptyCommit(message string) *Shell {
	return self.RunCommand([]string{"git", "commit", "--allow-empty", "-m", message})
}

func (self *Shell) EmptyCommitDaysAgo(message string, daysAgo int) *Shell {
	return self.RunCommand([]string{"git", "commit", "--allow-empty", "--date", fmt.Sprintf("%d days ago", daysAgo), "-m", message})
}

func (self *Shell) EmptyCommitWithDate(message string, date string) *Shell {
	env := []string{
		"GIT_AUTHOR_DATE=" + date,
		"GIT_COMMITTER_DATE=" + date,
	}
	return self.RunCommandWithEnv([]string{"git", "commit", "--allow-empty", "-m", message}, env)
}

func (self *Shell) Revert(ref string) *Shell {
	return self.RunCommand([]string{"git", "revert", ref})
}

func (self *Shell) CreateLightweightTag(name string, ref string) *Shell {
	return self.RunCommand([]string{"git", "tag", name, ref})
}

func (self *Shell) CreateAnnotatedTag(name string, message string, ref string) *Shell {
	return self.RunCommand([]string{"git", "tag", "-a", name, "-m", message, ref})
}

func (self *Shell) PushBranch(upstream, branch string) *Shell {
	return self.RunCommand([]string{"git", "push", upstream, branch})
}

func (self *Shell) PushBranchAndSetUpstream(upstream, branch string) *Shell {
	return self.RunCommand([]string{"git", "push", "--set-upstream", upstream, branch})
}

// convenience method for creating a file and adding it
func (self *Shell) CreateFileAndAdd(fileName string, fileContents string) *Shell {
	return self.
		CreateFile(fileName, fileContents).
		GitAdd(fileName)
}

// convenience method for updating a file and adding it
func (self *Shell) UpdateFileAndAdd(fileName string, fileContents string) *Shell {
	return self.
		UpdateFile(fileName, fileContents).
		GitAdd(fileName)
}

// convenience method for deleting a file and adding it
func (self *Shell) DeleteFileAndAdd(fileName string) *Shell {
	return self.
		DeleteFile(fileName).
		GitAdd(fileName)
}

// creates commits 01, 02, 03, ..., n with a new file in each
// The reason for padding with zeroes is so that it's easier to do string
// matches on the commit messages when there are many of them
func (self *Shell) CreateNCommits(n int) *Shell {
	return self.CreateNCommitsStartingAt(n, 1)
}

func (self *Shell) CreateNCommitsStartingAt(n, startIndex int) *Shell {
	for i := startIndex; i < startIndex+n; i++ {
		self.CreateFileAndAdd(
			fmt.Sprintf("file%02d.txt", i),
			fmt.Sprintf("file%02d content", i),
		).
			Commit(fmt.Sprintf("commit %02d", i))
	}

	return self
}

// Only to be used in demos, because the list might change and we don't want
// tests to break when it does.
func (self *Shell) CreateNCommitsWithRandomMessages(n int) *Shell {
	for i := 0; i < n; i++ {
		file := RandomFiles[i]
		self.CreateFileAndAdd(
			file.Name,
			file.Content,
		).
			Commit(RandomCommitMessages[i])
	}

	return self
}

// This creates a repo history of commits
// It uses a branching strategy where each feature branch is directly branched off
// of the master branch
// Only to be used in demos
func (self *Shell) CreateRepoHistory() *Shell {
	authors := []string{"Yang Wen-li", "Siegfried Kircheis", "Paul Oberstein", "Oscar Reuenthal", "Fredrica Greenhill"}

	numAuthors := 5
	numBranches := 10
	numInitialCommits := 20
	maxCommitsPerBranch := 5
	// Each commit will happen on a separate day
	repoStartDaysAgo := 100

	totalCommits := 0

	// Generate commits
	for i := 0; i < numInitialCommits; i++ {
		author := authors[i%numAuthors]
		commitMessage := RandomCommitMessages[totalCommits%len(RandomCommitMessages)]

		self.SetAuthor(author, "")
		self.EmptyCommitDaysAgo(commitMessage, repoStartDaysAgo-totalCommits)
		totalCommits++
	}

	// Generate branches and merges
	for i := 0; i < numBranches; i++ {
		// We'll have one author creating all the commits in the branch
		author := authors[i%numAuthors]
		branchName := RandomBranchNames[i%len(RandomBranchNames)]

		// Choose a random commit within the last 20 commits on the master branch
		lastMasterCommit := totalCommits - 1
		commitOffset := rand.Intn(min(lastMasterCommit, 5)) + 1

		// Create the feature branch and checkout the chosen commit
		self.NewBranchFrom(branchName, fmt.Sprintf("master~%d", commitOffset))

		numCommitsInBranch := rand.Intn(maxCommitsPerBranch) + 1
		for j := 0; j < numCommitsInBranch; j++ {
			commitMessage := RandomCommitMessages[totalCommits%len(RandomCommitMessages)]

			self.SetAuthor(author, "")
			self.EmptyCommitDaysAgo(commitMessage, repoStartDaysAgo-totalCommits)
			totalCommits++
		}

		self.Checkout("master")

		prevCommitterDate := os.Getenv("GIT_COMMITTER_DATE")
		prevAuthorDate := os.Getenv("GIT_AUTHOR_DATE")

		commitDate := time.Now().Add(time.Duration(totalCommits-repoStartDaysAgo) * time.Hour * 24)
		os.Setenv("GIT_COMMITTER_DATE", commitDate.Format(time.RFC3339))
		os.Setenv("GIT_AUTHOR_DATE", commitDate.Format(time.RFC3339))

		// Merge branch into master
		self.RunCommand([]string{"git", "merge", "--no-ff", branchName, "-m", fmt.Sprintf("Merge %s into master", branchName)})

		os.Setenv("GIT_COMMITTER_DATE", prevCommitterDate)
		os.Setenv("GIT_AUTHOR_DATE", prevAuthorDate)
	}

	return self
}

// Creates a commit with a random file
// Only to be used in demos
func (self *Shell) RandomChangeCommit(message string) *Shell {
	index := self.randomFileContentIndex
	self.randomFileContentIndex++
	randomFileName := fmt.Sprintf("random-%d.go", index)
	self.CreateFileAndAdd(randomFileName, RandomFileContents[index%len(RandomFileContents)])
	return self.Commit(message)
}

func (self *Shell) SetConfig(key string, value string) *Shell {
	self.RunCommand([]string{"git", "config", "--local", key, value})
	return self
}

func (self *Shell) CloneIntoRemote(name string) *Shell {
	self.Clone(name)
	self.RunCommand([]string{"git", "remote", "add", name, "../" + name})
	self.RunCommand([]string{"git", "fetch", name})

	return self
}

func (self *Shell) CloneIntoSubmodule(submoduleName string, submodulePath string) *Shell {
	self.Clone("other_repo")
	self.RunCommand([]string{"git", "submodule", "add", "--name", submoduleName, "../other_repo", submodulePath})

	return self
}

func (self *Shell) Clone(repoName string) *Shell {
	self.RunCommand([]string{"git", "clone", "--bare", ".", "../" + repoName})

	return self
}

func (self *Shell) SetBranchUpstream(branch string, upstream string) *Shell {
	self.RunCommand([]string{"git", "branch", "--set-upstream-to=" + upstream, branch})

	return self
}

func (self *Shell) RemoveRemoteBranch(remoteName string, branch string) *Shell {
	self.RunCommand([]string{"git", "-C", "../" + remoteName, "branch", "-d", branch})

	return self
}

func (self *Shell) HardReset(ref string) *Shell {
	self.RunCommand([]string{"git", "reset", "--hard", ref})
	return self
}

func (self *Shell) Stash(message string) *Shell {
	self.RunCommand([]string{"git", "stash", "push", "-m", message})
	return self
}

func (self *Shell) StartBisect(good string, bad string) *Shell {
	self.RunCommand([]string{"git", "bisect", "start", good, bad})
	return self
}

func (self *Shell) Init() *Shell {
	self.RunCommand([]string{"git", "-c", "init.defaultBranch=master", "init"})
	return self
}

func (self *Shell) AddWorktree(base string, path string, newBranchName string) *Shell {
	return self.RunCommand([]string{
		"git", "worktree", "add", "-b",
		newBranchName, path, base,
	})
}

// add worktree and have it checkout the base branch
func (self *Shell) AddWorktreeCheckout(base string, path string) *Shell {
	return self.RunCommand([]string{
		"git", "worktree", "add", path, base,
	})
}

func (self *Shell) AddFileInWorktree(worktreePath string) *Shell {
	self.CreateFile(filepath.Join(worktreePath, "content"), "content")

	self.RunCommand([]string{
		"git", "-C", worktreePath, "add", "content",
	})

	return self
}

func (self *Shell) MakeExecutable(path string) *Shell {
	// 0755 sets the executable permission for owner, and read/execute permissions for group and others
	err := os.Chmod(filepath.Join(self.dir, path), 0o755)
	if err != nil {
		panic(err)
	}

	return self
}

// Help files are located at test/files from the root the lazygit repo.
// E.g. You may want to create a pre-commit hook file there, then call this
// function to copy it into your test repo.
func (self *Shell) CopyHelpFile(source string, destination string) *Shell {
	return self.CopyFile(fmt.Sprintf("../../../../../files/%s", source), destination)
}

func (self *Shell) CopyFile(source string, destination string) *Shell {
	absSourcePath := filepath.Join(self.dir, source)
	absDestPath := filepath.Join(self.dir, destination)
	sourceFile, err := os.Open(absSourcePath)
	if err != nil {
		self.fail(err.Error())
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(absDestPath)
	if err != nil {
		self.fail(err.Error())
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		self.fail(err.Error())
	}

	// copy permissions to destination file too
	sourceFileInfo, err := os.Stat(absSourcePath)
	if err != nil {
		self.fail(err.Error())
	}

	err = os.Chmod(absDestPath, sourceFileInfo.Mode())
	if err != nil {
		self.fail(err.Error())
	}

	return self
}

// The final value passed to Chdir() during setup
// will be the directory the test is run from.
func (self *Shell) Chdir(path string) *Shell {
	self.dir = filepath.Join(self.dir, path)

	return self
}

func (self *Shell) SetAuthor(authorName string, authorEmail string) *Shell {
	self.RunCommand([]string{"git", "config", "--local", "user.name", authorName})
	self.RunCommand([]string{"git", "config", "--local", "user.email", authorEmail})

	return self
}
