package custom_commands

import (
	"runtime"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var BasicCmdAtRuntime = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using a custom command provided at runtime to create a new file",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("blah")
	},
	SetupConfig: func(cfg *config.AppConfig) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		cmd := "touch file.txt"
		if runtime.GOOS == "windows" {
			cmd = "copy NUL file.txt"
		}
		t.Views().Files().
			IsEmpty().
			IsFocused().
			Press(keys.Universal.ExecuteCustomCommand)

		t.ExpectPopup().Prompt().
			Title(Equals("Custom command:")).
			Type(cmd).
			Confirm()

		t.GlobalPress(keys.Files.RefreshFiles)

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file.txt"),
			)
	},
})
