package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var GenerateCommitMessageError = NewIntegrationTest(NewIntegrationTestArgs{
	Description: "Show stderr when commit message generation fails",
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.Commit.MessageGeneratorCommand = `sh -c 'echo generator failed >&2; exit 1' sh`
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "file content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Type("existing message").
			OpenCommitMenu()

		t.ExpectPopup().Menu().Title(Equals("Commit Menu")).
			Select(Contains("Generate Commit Message")).
			Confirm()

		t.ExpectPopup().Alert().
			Title(Equals("Generate commit message command failed")).
			Content(Equals("generator failed")).
			Confirm()

		t.ExpectPopup().CommitMessagePanel().
			Content(Equals("existing message"))
	},
})
