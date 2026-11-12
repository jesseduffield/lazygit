package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ClearCommitMessage = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Clear the commit message in the commit message panel",
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
			Type("first line").
			AddNewline().
			Type("second line").
			SwitchToSummary().
			OpenCommitMenu()

		t.ExpectPopup().Menu().Title(Equals("Commit Menu")).
			Select(Contains("Clear commit message")).
			Confirm()

		t.ExpectPopup().CommitMessagePanel().
			Content(Equals("")).
			SwitchToDescription().
			Content(Equals(""))
	},
})
