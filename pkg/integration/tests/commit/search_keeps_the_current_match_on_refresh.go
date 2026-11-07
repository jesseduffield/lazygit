package commit

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SearchKeepsTheCurrentMatchOnRefresh = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "With the selection between two search matches, a refresh keeps the match that the selection made the current one",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		// So that R refreshes in the commits panel instead of rewording the commit
		cfg.GetUserConfig().Keybinding.Commits.RenameCommitWithEditor = config.Keybinding{"<disabled>"}
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch")
		for i := 1; i <= 8; i++ {
			message := fmt.Sprintf("commit %02d", i)
			if i == 7 || i == 3 {
				message += " needle"
			}
			shell.EmptyCommit(message)
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			FilterOrSearch("needle").
			SelectedLines(
				Contains("commit 07 needle"),
			).
			NavigateToLine(Contains("commit 05")).
			Tap(func() {
				t.Views().Search().Content(Contains("matches for 'needle' (1 of 2)"))

				t.GlobalPress(keys.Universal.Refresh)
				t.Views().Search().Content(
					Contains("matches for 'needle' (1 of 2)"))
			}).
			NavigateToLine(Contains("commit 01")).
			Tap(func() {
				t.Views().Search().Content(Contains("matches for 'needle' (2 of 2)"))

				t.GlobalPress(keys.Universal.Refresh)
				t.Views().Search().Content(
					Contains("matches for 'needle' (2 of 2)"))
			}).
			// The next match wraps around to the first one
			Press(keys.Universal.NextMatch).
			SelectedLines(
				Contains("commit 07 needle"),
			)
	},
})
