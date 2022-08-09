package helpers

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jesseduffield/lazygit/pkg/integration/types"
	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/mgutz/str"
)

type ShellImpl struct{}

var _ types.Shell = &ShellImpl{}

func (s *ShellImpl) RunCommand(cmdStr string) types.Shell {
	args := str.ToArgv(cmdStr)
	cmd := secureexec.Command(args[0], args[1:]...)
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("error running command: %s\n%s", cmdStr, string(output)))
	}

	return s
}

func (s *ShellImpl) CreateFile(path string, content string) types.Shell {
	err := ioutil.WriteFile(path, []byte(content), 0o644)
	if err != nil {
		panic(fmt.Sprintf("error creating file: %s\n%s", path, err))
	}

	return s
}

func (s *ShellImpl) NewBranch(name string) types.Shell {
	return s.RunCommand("git checkout -b " + name)
}

func (s *ShellImpl) GitAdd(path string) types.Shell {
	return s.RunCommand(fmt.Sprintf("git add \"%s\"", path))
}

func (s *ShellImpl) GitAddAll() types.Shell {
	return s.RunCommand("git add -A")
}

func (s *ShellImpl) Commit(message string) types.Shell {
	return s.RunCommand(fmt.Sprintf("git commit -m \"%s\"", message))
}

func (s *ShellImpl) EmptyCommit(message string) types.Shell {
	return s.RunCommand(fmt.Sprintf("git commit --allow-empty -m \"%s\"", message))
}

func (s *ShellImpl) CreateFileAndAdd(fileName string, fileContents string) types.Shell {
	return s.
		CreateFile(fileName, fileContents).
		GitAdd(fileName)
}

func (s *ShellImpl) CreateNCommits(n int) types.Shell {
	for i := 1; i <= n; i++ {
		s.CreateFileAndAdd(
			fmt.Sprintf("file%02d.txt", i),
			fmt.Sprintf("file%02d content", i),
		).
			Commit(fmt.Sprintf("commit %02d", i))
	}

	return s
}
