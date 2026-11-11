package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PasteMultiLineCommitMessage = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Paste a multi-line commit message into the summary with the terminal's paste command",
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
			Paste("subject\n\nbody line 1\nbody line 2").
			Content(Equals("subject"))

		t.Views().CommitDescription().
			IsFocused().
			Content(Equals(
				/* EXPECTED:
				"body line 1\nbody line 2"))
				ACTUAL: */
				"\nbody line 1\nbody line 2"))
	},
})
