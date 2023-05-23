package oscommands

import (
	"testing"
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
		cmdObj := &CmdObj{args: scenario.cmdArgs}
		actual := cmdObj.ToString()
		if actual != scenario.expected {
			t.Errorf("Expected %s, got %s", quote(scenario.expected), quote(actual))
		}
	}
}
