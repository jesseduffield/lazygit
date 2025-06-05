package shell_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ComplexShellCommand = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using a custom command provided at runtime to create a new file, via a shell command. We invoke custom commands through a shell already. This test proves that we can run a shell within a shell, which requires complex escaping.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("blah")
	},
	SetupConfig: func(cfg *config.AppConfig) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsEmpty().
			IsFocused().
			Press(keys.Universal.ExecuteShellCommand)

		t.ExpectPopup().Prompt().
			Title(Equals("Shell command:")).
			Type("sh -c \"touch file.txt\"").
			Confirm()

		t.GlobalPress(keys.Files.RefreshFiles)

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file.txt"),
			)
	},
})
