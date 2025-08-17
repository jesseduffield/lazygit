package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FailHooksThenCommitNoHooks = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that commit message can be reused in commit without hook after failing commit with hooks",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile(".git/hooks/pre-commit", blockingHook)
		shell.MakeExecutable(".git/hooks/pre-commit")

		shell.CreateFileAndAdd("one", "one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("one"),
			).
			Press(keys.Files.CommitChanges).
			Tap(func() {
				t.ExpectPopup().CommitMessagePanel().Type("my message").Confirm()

				t.ExpectPopup().Alert().Title(Equals("Error")).Content(Contains("Git command failed")).Confirm()
			}).
			Press(keys.Files.CommitChangesWithoutHook).
			Tap(func() {
				t.ExpectPopup().CommitMessagePanel().
					InitialText(Equals("my message")). // it remembered the commit message
					Confirm()

				t.Views().Commits().
					Lines(
						Contains("my message"),
					)
			})
		t.Views().Commits().Focus()
		t.Views().Main().Content(Contains("my message"))
	},
})
