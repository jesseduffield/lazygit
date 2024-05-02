package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AddCoAuthor = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Add co-author on a commit",
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
					Select(Contains("Add co-author")).
					Confirm()

				t.ExpectPopup().Prompt().
					Title(Contains("Add co-author")).
					Type("John Smith <jsmith@gmail.com>").
					Confirm()
			})

		t.Views().Main().ContainsLines(
			Equals("    initial commit"),
			Equals("    "),
			Equals("    Co-authored-by: John Smith <jsmith@gmail.com>"),
		)
	},
})
