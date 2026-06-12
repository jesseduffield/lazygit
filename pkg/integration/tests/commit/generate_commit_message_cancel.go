package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var GenerateCommitMessageCancel = NewIntegrationTest(NewIntegrationTestArgs{
	Description: "Cancel a running commit message generator",
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.Commit.MessageGeneratorCommand = `sh -c 'sleep 5; printf "generated subject"' sh`
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

		menu := t.ExpectPopup().Menu().Title(Equals("Commit Menu")).
			Select(Contains("Generate Commit Message"))

		go func() {
			t.Wait(100)
			t.GlobalPressWithoutWaiting(keys.Universal.Return)
		}()

		menu.Confirm()

		t.ExpectPopup().CommitMessagePanel().
			Content(Equals("existing message"))
	},
})
