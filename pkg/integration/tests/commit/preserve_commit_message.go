package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var PreserveCommitMessage = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Test that the commit message is preserved correctly when canceling the commit message panel",
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
			InitialText(Equals("")).
			Type("my commit message").
			SwitchToDescription().
			Type("first paragraph").
			AddNewline().
			AddNewline().
			Type("second paragraph").
			Cancel()

		t.Views().Files().
			IsFocused().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Content(Equals("my commit message")).
			SwitchToDescription().
			Content(Equals("first paragraph\n\nsecond paragraph"))
	},
})
