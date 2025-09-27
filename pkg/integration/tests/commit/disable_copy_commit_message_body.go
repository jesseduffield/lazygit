package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DisableCopyCommitMessageBody = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Disables copy commit message body when there is no body",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},

	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("commit")
	},

	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit").IsSelected(),
			).
			Press(keys.Commits.CopyCommitAttributeToClipboard)

		t.ExpectPopup().Menu().
			Title(Equals("Copy to clipboard")).
			Select(Contains("Commit message body")).
			Confirm()

		t.ExpectToast(Equals("Disabled: Commit has no message body"))
	},
})
