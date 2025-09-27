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
		config.GetUserConfig().OS.CopyToClipboardCmd = "printf '%s' {{text}} > clipboard"
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

		t.FileSystem().FileContent("clipboard", Equals("tag2\ntag1"))
	},
})
