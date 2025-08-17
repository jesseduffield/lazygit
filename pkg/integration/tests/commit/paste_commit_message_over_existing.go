package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PasteCommitMessageOverExisting = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Paste a commit message into the commit message panel when there is already text in the panel, causing a confirmation",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.CopyToClipboardCmd = "printf '%s' {{text}} > ../clipboard"
		config.GetUserConfig().OS.ReadFromClipboardCmd = "cat ../clipboard"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("subject\n\nbody 1st line\nbody 2nd line")
		shell.CreateFileAndAdd("file", "file content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			ContainsLines(
				Contains("subject").IsSelected(),
			).
			Press(keys.Commits.CopyCommitAttributeToClipboard)

		t.ExpectPopup().Menu().Title(Equals("Copy to clipboard")).
			Select(Contains("Commit message (subject and body)")).Confirm()

		t.ExpectToast(Equals("Commit message copied to clipboard"))

		t.Views().Files().
			Focus().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Type("existing message").
			OpenCommitMenu()

		t.ExpectPopup().Menu().Title(Equals("Commit Menu")).
			Select(Contains("Paste commit message from clipboard")).
			Confirm()

		t.ExpectPopup().Alert().Title(Equals("Paste commit message from clipboard")).
			Content(Equals("Pasting will overwrite the current commit message, continue?")).
			Confirm()

		t.ExpectPopup().CommitMessagePanel().
			Content(Equals("subject")).
			SwitchToDescription().
			Content(Equals("body 1st line\nbody 2nd line"))
	},
})
