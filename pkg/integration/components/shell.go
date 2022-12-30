package components

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/mgutz/str"
)

// this is for running shell commands, mostly for the sake of setting up the repo
// but you can also run the commands from within lazygit to emulate things happening
// in the background.
type Shell struct {
	// working directory the shell is invoked in
	dir string
	// when running the shell outside the gui we can directly panic on failure,
	// but inside the gui we need to close the gui before panicking
	fail func(string)
}

func NewShell(dir string, fail func(string)) *Shell {
	return &Shell{dir: dir, fail: fail}
}

func (self *Shell) RunCommand(cmdStr string) *Shell {
	args := str.ToArgv(cmdStr)
	cmd := secureexec.Command(args[0], args[1:]...)
	cmd.Env = os.Environ()
	cmd.Dir = self.dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		self.fail(fmt.Sprintf("error running command: %s\n%s", cmdStr, string(output)))
	}

	return self
}

func (self *Shell) runCommandWithOutput(cmdStr string) (string, error) {
	args := str.ToArgv(cmdStr)
	cmd := secureexec.Command(args[0], args[1:]...)
	cmd.Env = os.Environ()
	cmd.Dir = self.dir

	output, err := cmd.CombinedOutput()

	return string(output), err
}

func (self *Shell) RunShellCommand(cmdStr string) *Shell {
	cmd := secureexec.Command("sh", "-c", cmdStr)
	cmd.Env = os.Environ()
	cmd.Dir = self.dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		self.fail(fmt.Sprintf("error running shell command: %s\n%s", cmdStr, string(output)))
	}

	return self
}

func (self *Shell) RunShellCommandExpectError(cmdStr string) *Shell {
	cmd := secureexec.Command("sh", "-c", cmdStr)
	cmd.Env = os.Environ()
	cmd.Dir = self.dir

	output, err := cmd.CombinedOutput()
	if err == nil {
		self.fail(fmt.Sprintf("Expected error running shell command: %s\n%s", cmdStr, string(output)))
	}

	return self
}

func (self *Shell) CreateFile(path string, content string) *Shell {
	fullPath := filepath.Join(self.dir, path)
	err := os.WriteFile(fullPath, []byte(content), 0o644)
	if err != nil {
		self.fail(fmt.Sprintf("error creating file: %s\n%s", fullPath, err))
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
	return self.RunCommand("git checkout -b " + name)
}

func (self *Shell) Checkout(name string) *Shell {
	return self.RunCommand("git checkout " + name)
}

func (self *Shell) Merge(name string) *Shell {
	return self.RunCommand("git merge --commit --no-ff " + name)
}

func (self *Shell) GitAdd(path string) *Shell {
	return self.RunCommand(fmt.Sprintf("git add \"%s\"", path))
}

func (self *Shell) GitAddAll() *Shell {
	return self.RunCommand("git add -A")
}

func (self *Shell) Commit(message string) *Shell {
	return self.RunCommand(fmt.Sprintf("git commit -m \"%s\"", message))
}

func (self *Shell) EmptyCommit(message string) *Shell {
	return self.RunCommand(fmt.Sprintf("git commit --allow-empty -m \"%s\"", message))
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

// creates commits 01, 02, 03, ..., n with a new file in each
// The reason for padding with zeroes is so that it's easier to do string
// matches on the commit messages when there are many of them
func (self *Shell) CreateNCommits(n int) *Shell {
	for i := 1; i <= n; i++ {
		self.CreateFileAndAdd(
			fmt.Sprintf("file%02d.txt", i),
			fmt.Sprintf("file%02d content", i),
		).
			Commit(fmt.Sprintf("commit %02d", i))
	}

	return self
}

func (self *Shell) StashWithMessage(message string) *Shell {
	self.RunCommand(fmt.Sprintf(`git stash -m "%s"`, message))
	return self
}

func (self *Shell) SetConfig(key string, value string) *Shell {
	self.RunCommand(fmt.Sprintf(`git config --local "%s" %s`, key, value))
	return self
}

// creates a clone of the repo in a sibling directory and adds the clone
// as a remote, then fetches it.
func (self *Shell) CloneIntoRemote(name string) *Shell {
	self.RunCommand(fmt.Sprintf("git clone --bare . ../%s", name))
	self.RunCommand(fmt.Sprintf("git remote add %s ../%s", name, name))
	self.RunCommand(fmt.Sprintf("git fetch %s", name))

	return self
}

// e.g. branch: 'master', upstream: 'origin/master'
func (self *Shell) SetBranchUpstream(branch string, upstream string) *Shell {
	self.RunCommand(fmt.Sprintf("git branch --set-upstream-to=%s %s", upstream, branch))

	return self
}

func (self *Shell) RemoveRemoteBranch(remoteName string, branch string) *Shell {
	self.RunCommand(fmt.Sprintf("git -C ../%s branch -d %s", remoteName, branch))

	return self
}

func (self *Shell) HardReset(ref string) *Shell {
	self.RunCommand(fmt.Sprintf("git reset --hard %s", ref))

	return self
}
