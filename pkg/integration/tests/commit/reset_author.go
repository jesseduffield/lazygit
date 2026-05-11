package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ResetAuthor = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Reset author on a commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.SetConfig("user.email", "Bill@example.com")
		shell.SetConfig("user.name", "Bill Smith")

		shell.EmptyCommit("one")

		shell.SetConfig("user.email", "John@example.com")
		shell.SetConfig("user.name", "John Smith")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("BS").Contains("one").IsSelected(),
			).
			Press(keys.Commits.ResetCommitAuthor).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Amend commit attribute")).
					Select(Contains("Reset author")).
					Confirm()
			}).
			Lines(
				Contains("JS").Contains("one").IsSelected(),
			)
	},
})
