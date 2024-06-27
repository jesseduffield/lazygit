package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PasteCommitMessage = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Paste a commit message into the commit message panel",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.OS.CopyToClipboardCmd = "echo {{text}} > ../clipboard"
		config.UserConfig.OS.ReadFromClipboardCmd = "cat ../clipboard"
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
			Select(Contains("Commit message")).Confirm()

		t.ExpectToast(Equals("Commit message copied to clipboard"))

		t.Views().Files().
			Focus().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			OpenCommitMenu()

		t.ExpectPopup().Menu().Title(Equals("Commit Menu")).
			Select(Contains("Paste commit message from clipboard")).
			Confirm()

		t.ExpectPopup().CommitMessagePanel().
			Content(Equals("subject")).
			SwitchToDescription().
			Content(Equals("body 1st line\nbody 2nd line"))
	},
})
