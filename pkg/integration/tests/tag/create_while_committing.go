package tag

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CreateWhileCommitting = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Draft a commit message, escape out, and make a tag. Verify the draft message doesn't appear in the tag create prompt",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CreateFileAndAdd("file.txt", "file contents")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Press(keys.Files.CommitChanges).
			Tap(func() {
				t.ExpectPopup().CommitMessagePanel().
					Title(Equals("Commit summary")).
					Type("draft message").
					Cancel()
			})

		t.Views().Tags().
			Focus().
			IsEmpty().
			Press(keys.Universal.New).
			Tap(func() {
				t.ExpectPopup().CommitMessagePanel().
					Title(Equals("Tag name")).
					InitialText(Equals(""))
			})
	},
})
