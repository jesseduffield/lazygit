package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PasteMultiLineCommitMessageWithHardLineBreaks = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Paste a commit message whose body was wrapped at the auto-wrap width into the summary with the terminal's paste command, and check that the body reflows when editing it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "file content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Paste("subject\n\n" +
				"The body of this commit was wrapped at the auto-wrap width when it was\n" +
				"committed, so it contains a hard line break.").
			Content(Equals("subject"))

		t.ExpectPopup().CommitDescriptionPanel().
			GoToBeginning().
			Type("Now ").
			Content(Equals(
				/* EXPECTED:
				"Now The body of this commit was wrapped at the auto-wrap width when it \n" +
					"was committed, so it contains a hard line break."))
				ACTUAL: */
				"Now The body of this commit was wrapped at the auto-wrap width when it \n" +
					"was\n" +
					"committed, so it contains a hard line break."))
	},
})
