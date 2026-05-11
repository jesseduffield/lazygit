package misc

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CopyConfirmationMessageToClipboard = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Copy the text of a confirmation popup to the clipboard",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.CopyToClipboardCmd = "printf '%s' {{text}} > clipboard"
	},

	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("commit")
	},

	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit").IsSelected(),
			).
			Press(keys.Universal.Remove)

		t.ExpectPopup().Alert().
			Title(Equals("Drop commit")).
			Content(Equals("Are you sure you want to drop the selected commit(s)?")).
			Tap(func() {
				t.GlobalPress(keys.Universal.CopyToClipboard)
				t.ExpectToast(Equals("Message copied to clipboard"))
			}).
			Confirm()

		t.FileSystem().FileContent("clipboard",
			Equals("Are you sure you want to drop the selected commit(s)?"))
	},
})
