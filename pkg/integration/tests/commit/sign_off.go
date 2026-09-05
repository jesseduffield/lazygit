package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SignOff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Sign off a commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("initial commit").IsSelected(),
			).
			Press(keys.Commits.ResetCommitAuthor).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Amend commit attribute")).
					Select(Contains("Sign off")).
					Confirm()
			})

		t.Views().Main().ContainsLines(
			Equals("    initial commit"),
			Equals("    "),
			Equals("    Signed-off-by: CI <CI@example.com>"),
		)
	},
})
