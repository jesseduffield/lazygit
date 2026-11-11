package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PasteCommitMessageWithHardLineBreaks = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Paste a commit message whose body was wrapped at the auto-wrap width, and check that the body reflows when editing it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().OS.CopyToClipboardCmd = "printf '%s' {{text}} > ../clipboard"
		config.GetUserConfig().OS.ReadFromClipboardCmd = "cat ../clipboard"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("subject\n\n" +
			"The body of this commit was wrapped at the auto-wrap width when it was\n" +
			"committed, so it contains a hard line break.")
		shell.CreateFileAndAdd("file", "file content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Commits.CopyCommitAttributeToClipboard)

		t.ExpectPopup().Menu().Title(Equals("Copy to clipboard")).
			Select(Contains("Commit message (subject and body)")).Confirm()

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
			GoToBeginning().
			Type("Now ").
			Content(Equals(
				"Now The body of this commit was wrapped at the auto-wrap width when it \n" +
					"was committed, so it contains a hard line break."))
	},
})
