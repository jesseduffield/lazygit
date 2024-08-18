package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectedPath = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Use the {{ .SelectedPath }} template variable in different contexts",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.CreateDir("folder1")
		shell.CreateFileAndAdd("folder1/file1", "")
		shell.Commit("commit")
		shell.CreateDir("folder2")
		shell.CreateFile("folder2/file2", "")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:     "X",
				Context: "global",
				Command: "printf '%s' '{{ .SelectedPath }}' > file.txt",
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			NavigateToLine(Contains("file2"))
		t.GlobalPress("X")
		t.FileSystem().FileContent("file.txt", Equals("folder2/file2"))

		t.Views().Commits().
			Focus().
			PressEnter()
		t.Views().CommitFiles().
			IsFocused().
			NavigateToLine(Contains("file1"))
		t.GlobalPress("X")
		t.FileSystem().FileContent("file.txt", Equals("folder1/file1"))
	},
})
