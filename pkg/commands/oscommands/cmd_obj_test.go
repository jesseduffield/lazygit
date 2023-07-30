package oscommands

import (
	"os/exec"
	"testing"

	"github.com/jesseduffield/gocui"
)

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
