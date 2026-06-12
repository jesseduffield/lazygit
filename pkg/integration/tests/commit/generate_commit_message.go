package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var GenerateCommitMessage = NewIntegrationTest(NewIntegrationTestArgs{
	Description: "Generate a commit message from a configured command",
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.Commit.MessageGeneratorCommand = `sh -c 'printf "generated subject\n\nroot: %s" "$(basename "$1")"' sh`
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "file content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			OpenCommitMenu()

		t.ExpectPopup().Menu().Title(Equals("Commit Menu")).
			Select(Contains("Generate Commit Message")).
			Confirm()

		t.ExpectPopup().CommitMessagePanel().
			Content(Equals("generated subject")).
			SwitchToDescription().
			Content(Equals("root: repo"))
	},
})
