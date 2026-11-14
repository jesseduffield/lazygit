package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SearchKeepsTheSecondMatchOnALineOnRefresh = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "With the second of two search matches on the selected commit as the current one, a refresh keeps it current",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		// So that R refreshes in the commits panel instead of rewording the commit
		cfg.GetUserConfig().Keybinding.Commits.RenameCommitWithEditor = config.Keybinding{"<disabled>"}
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch")
		shell.EmptyCommit("commit 01")
		shell.EmptyCommit("commit 02 needle")
		shell.EmptyCommit("commit 03")
		shell.EmptyCommit("commit 04 needle and needle")
		shell.EmptyCommit("commit 05")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			FilterOrSearch("needle").
			SelectedLines(
				Contains("commit 04 needle and needle"),
			).
			Press(keys.Universal.NextMatch).
			SelectedLines(
				Contains("commit 04 needle and needle"),
			).
			Tap(func() {
				t.Views().Search().Content(Contains("matches for 'needle' (2 of 3)"))

				t.GlobalPress(keys.Universal.Refresh)
				t.Views().Search().Content(
					Contains("matches for 'needle' (2 of 3)"))
			}).
			Press(keys.Universal.NextMatch).
			SelectedLines(
				Contains("commit 02 needle"),
			)
	},
})
