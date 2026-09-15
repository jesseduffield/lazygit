package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PreserveCommitMessageWhitespace = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Whitespace in the description (e.g. leading blank lines, indented first line) should be preserved when canceling and reopening the commit message panel",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("myfile", "myfile content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Type("my commit message").
			SwitchToDescription().
			AddNewline().
			AddNewline().
			Type("body  ").
			Cancel()

		t.Views().Files().
			IsFocused().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Content(Equals("my commit message")).
			SwitchToDescription().
			Content(Equals("\n\nbody  ")).
			Cancel()
	},
})
