package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// We're emulating the clipboard by writing to a file called clipboard

var CopyTagToClipboard = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Copy a commit tag to the clipboard",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		// Include delimiters around the text so that we can assert on the entire content
		config.GetUserConfig().OS.CopyToClipboardCmd = "echo _{{text}}_ > clipboard"
	},

	SetupRepo: func(shell *Shell) {
		shell.SetAuthor("John Doe", "john@doe.com")
		shell.EmptyCommit("commit")
		shell.CreateLightweightTag("tag1", "HEAD")
		shell.CreateLightweightTag("tag2", "HEAD")
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
			Select(Contains("Commit tags")).
			Confirm()

		t.ExpectToast(Equals("Commit tags copied to clipboard"))

		t.Views().Files().
			Focus().
			Press(keys.Files.RefreshFiles).
			Lines(
				Contains("clipboard").IsSelected(),
			)

		t.Views().Main().Content(Contains("+_tag2"))
		t.Views().Main().Content(Contains("+tag1_"))
	},
})
