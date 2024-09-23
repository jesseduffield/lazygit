package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// We're emulating the clipboard by writing to a file called clipboard

var CopyAuthorToClipboard = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Copy a commit author name to the clipboard",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		// Include delimiters around the text so that we can assert on the entire content
		config.GetUserConfig().OS.CopyToClipboardCmd = "echo /{{text}}/ > clipboard"
	},

	SetupRepo: func(shell *Shell) {
		shell.SetAuthor("John Doe", "john@doe.com")
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
			Select(Contains("Commit author")).
			Confirm()

		t.ExpectToast(Equals("Commit author copied to clipboard"))

		t.Views().Files().
			Focus().
			Press(keys.Files.RefreshFiles).
			Lines(
				Contains("clipboard").IsSelected(),
			)

		t.Views().Main().Content(Contains("/John Doe <john@doe.com>/"))
	},
})
