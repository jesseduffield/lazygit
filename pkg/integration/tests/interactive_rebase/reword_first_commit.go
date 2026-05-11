package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Rewording the first commit is tricky because you can't rebase from its parent commit,
// hence having a specific test for this

var RewordFirstCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Rewords the first commit, just to show that it's possible",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(2)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit 02"),
				Contains("commit 01"),
			).
			NavigateToLine(Contains("commit 01")).
			Press(keys.Commits.RenameCommit).
			Tap(func() {
				t.ExpectPopup().CommitMessagePanel().
					Title(Equals("Reword commit")).
					InitialText(Equals("commit 01")).
					Clear().
					Type("renamed 01").
					Confirm()
			}).
			Lines(
				Contains("commit 02"),
				Contains("renamed 01"),
			)
	},
})
