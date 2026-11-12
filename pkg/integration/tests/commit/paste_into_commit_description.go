package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PasteIntoCommitDescription = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Paste text into the commit description with the terminal's paste command",
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
			Type("subject").
			SwitchToDescription().
			Type("last line").
			GoToBeginning().
			Paste("first line\n\tindented line\n").
			Content(Equals("first line\n    indented line\nlast line"))

		t.ExpectPopup().CommitDescriptionPanel()
	},
})
