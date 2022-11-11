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
}

func NewShell(dir string) *Shell {
	return &Shell{dir: dir}
}

func (s *Shell) RunCommand(cmdStr string) *Shell {
	args := str.ToArgv(cmdStr)
	cmd := secureexec.Command(args[0], args[1:]...)
	cmd.Env = os.Environ()
	cmd.Dir = s.dir

	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("error running command: %s\n%s", cmdStr, string(output)))
	}

	return s
}

func (s *Shell) CreateFile(path string, content string) *Shell {
	fullPath := filepath.Join(s.dir, path)
	err := os.WriteFile(fullPath, []byte(content), 0o644)
	if err != nil {
		panic(fmt.Sprintf("error creating file: %s\n%s", fullPath, err))
	}

	return s
}

func (s *Shell) CreateDir(path string) *Shell {
	fullPath := filepath.Join(s.dir, path)
	if err := os.MkdirAll(fullPath, 0o755); err != nil {
		panic(fmt.Sprintf("error creating directory: %s\n%s", fullPath, err))
	}

	return s
}

func (s *Shell) UpdateFile(path string, content string) *Shell {
	fullPath := filepath.Join(s.dir, path)
	err := os.WriteFile(fullPath, []byte(content), 0o644)
	if err != nil {
		panic(fmt.Sprintf("error updating file: %s\n%s", fullPath, err))
	}

	return s
}

func (s *Shell) NewBranch(name string) *Shell {
	return s.RunCommand("git checkout -b " + name)
}

func (s *Shell) Checkout(name string) *Shell {
	return s.RunCommand("git checkout " + name)
}

func (s *Shell) Merge(name string) *Shell {
	return s.RunCommand("git merge --commit --no-ff " + name)
}

func (s *Shell) GitAdd(path string) *Shell {
	return s.RunCommand(fmt.Sprintf("git add \"%s\"", path))
}

func (s *Shell) GitAddAll() *Shell {
	return s.RunCommand("git add -A")
}

func (s *Shell) Commit(message string) *Shell {
	return s.RunCommand(fmt.Sprintf("git commit -m \"%s\"", message))
}

func (s *Shell) EmptyCommit(message string) *Shell {
	return s.RunCommand(fmt.Sprintf("git commit --allow-empty -m \"%s\"", message))
}

// convenience method for creating a file and adding it
func (s *Shell) CreateFileAndAdd(fileName string, fileContents string) *Shell {
	return s.
		CreateFile(fileName, fileContents).
		GitAdd(fileName)
}

// convenience method for updating a file and adding it
func (s *Shell) UpdateFileAndAdd(fileName string, fileContents string) *Shell {
	return s.
		UpdateFile(fileName, fileContents).
		GitAdd(fileName)
}

// creates commits 01, 02, 03, ..., n with a new file in each
// The reason for padding with zeroes is so that it's easier to do string
// matches on the commit messages when there are many of them
func (s *Shell) CreateNCommits(n int) *Shell {
	for i := 1; i <= n; i++ {
		s.CreateFileAndAdd(
			fmt.Sprintf("file%02d.txt", i),
			fmt.Sprintf("file%02d content", i),
		).
			Commit(fmt.Sprintf("commit %02d", i))
	}

	return s
}

func (s *Shell) StashWithMessage(message string) *Shell {
	s.RunCommand(fmt.Sprintf(`git stash -m "%s"`, message))
	return s
}
