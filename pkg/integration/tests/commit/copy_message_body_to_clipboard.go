package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// We're emulating the clipboard by writing to a file called clipboard

var CopyMessageBodyToClipboard = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Copy a commit message body to the clipboard",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.CopyToClipboardCmd = "printf '%s' {{text}} > clipboard"
	},

	SetupRepo: func(shell *Shell) {
		shell.EmptyCommitWithBody("My Subject", "My awesome commit message body")
	},

	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("My Subject").IsSelected(),
			).
			Press(keys.Commits.CopyCommitAttributeToClipboard)

		t.ExpectPopup().Menu().
			Title(Equals("Copy to clipboard")).
			Select(Contains("Commit message body")).
			Confirm()

		t.ExpectToast(Equals("Commit message body copied to clipboard"))

		t.FileSystem().FileContent("clipboard", Equals("My awesome commit message body"))
	},
})
