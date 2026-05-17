package oscommands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewShellWindowsKeepsCommandAsSingleArg(t *testing.T) {
	builder := &CmdObjBuilder{
		platform: &Platform{
			OS:       "windows",
			Shell:    "cmd",
			ShellArg: "/c",
		},
	}

	command := `"C:\Program Files\Notepad++\notepad++.exe" -multiInst -nosession -noPlugin -n42 "C:\path\file.txt"`
	cmdObj := builder.NewShell(command, "")

	assert.Equal(t, []string{"cmd", "/c", command}, cmdObj.Args())
}

func TestNewShellUnixSplitsCommand(t *testing.T) {
	builder := &CmdObjBuilder{
		platform: &Platform{
			OS:       "linux",
			Shell:    "bash",
			ShellArg: "-c",
		},
	}

	cmdObj := builder.NewShell(`vim +42 -- "file.txt"`, "")

	assert.Equal(t, []string{"bash", "-c", `vim +42 -- "file.txt"`}, cmdObj.Args())
}
