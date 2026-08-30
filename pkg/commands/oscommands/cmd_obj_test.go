package oscommands

import (
	"os/exec"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/stretchr/testify/assert"
)

func TestRemoveEnvVar(t *testing.T) {
	cmd := exec.Command("git", "status")
	cmd.Env = []string{
		"PATH=/usr/bin",
		"GIT_OPTIONAL_LOCKS=0",
		"GIT_OPTIONAL_LOCKS_OTHER=1", // name is a prefix of ours but not the same var
		"GIT_OPTIONAL_LOCKS=0",       // duplicates must all be removed
		"HOME=/home/me",
	}
	cmdObj := &CmdObj{cmd: cmd}

	cmdObj.RemoveEnvVar("GIT_OPTIONAL_LOCKS")

	assert.Equal(t, []string{
		"PATH=/usr/bin",
		"GIT_OPTIONAL_LOCKS_OTHER=1",
		"HOME=/home/me",
	}, cmdObj.GetEnvVars())
}

func TestCmdObjToString(t *testing.T) {
	quote := func(s string) string {
		return "\"" + s + "\""
	}

	scenarios := []struct {
		cmdArgs  []string
		expected string
	}{
		{
			cmdArgs:  []string{"git", "push", "myfile.txt"},
			expected: "git push myfile.txt",
		},
		{
			cmdArgs:  []string{"git", "push", "my file.txt"},
			expected: "git push \"my file.txt\"",
		},
	}

	for _, scenario := range scenarios {
		cmd := exec.Command(scenario.cmdArgs[0], scenario.cmdArgs[1:]...)
		cmdObj := &CmdObj{cmd: cmd}
		actual := cmdObj.ToString()
		if actual != scenario.expected {
			t.Errorf("Expected %s, got %s", quote(scenario.expected), quote(actual))
		}
	}
}

func TestClone(t *testing.T) {
	task := gocui.NewFakeTask()
	cmdObj := &CmdObj{task: task, cmd: &exec.Cmd{}}
	clone := cmdObj.Clone()
	if clone == cmdObj {
		t.Errorf("Clone should not return the same object")
	}

	if clone.GetTask() == nil {
		t.Errorf("Clone task should not be nil")
	}

	if clone.GetTask() != task {
		t.Errorf("Clone should have the same task")
	}
}
