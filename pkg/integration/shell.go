package integration

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

func (s *ShellImpl) GitAddAll() types.Shell {
	return s.RunCommand("git add -A")
}

func (s *ShellImpl) Commit(message string) types.Shell {
	return s.RunCommand(fmt.Sprintf("git commit -m \"%s\"", message))
}

func (s *ShellImpl) EmptyCommit(message string) types.Shell {
	return s.RunCommand(fmt.Sprintf("git commit --allow-empty -m \"%s\"", message))
}
