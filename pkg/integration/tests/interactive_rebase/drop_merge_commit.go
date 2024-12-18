package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var DropMergeCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Drops a merge commit outside of an interactive rebase",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.22.0"), // first version that supports the --rebase-merges option
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.CreateMergeCommit(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI ⏣─╮ Merge branch 'second-change-branch' into first-change-branch").IsSelected(),
				Contains("CI │ ◯ * second-change-branch unrelated change"),
				Contains("CI │ ◯ second change"),
				Contains("CI ◯ │ first change"),
				Contains("CI ◯─╯ * original"),
				Contains("CI ◯ three"),
				Contains("CI ◯ two"),
				Contains("CI ◯ one"),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Drop commit")).
					Content(Equals("Are you sure you want to drop the selected merge commit? Note that it will also drop all the commits that were merged in by it.")).
					Confirm()
			}).
			Lines(
				Contains("CI ◯ first change").IsSelected(),
				Contains("CI ◯ * original"),
				Contains("CI ◯ three"),
				Contains("CI ◯ two"),
				Contains("CI ◯ one"),
			)
	},
})
