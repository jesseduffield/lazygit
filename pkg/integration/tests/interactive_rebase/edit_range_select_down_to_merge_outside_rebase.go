package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var EditRangeSelectDownToMergeOutsideRebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Select a range of commits (the last one being a merge commit) to edit outside of a rebase",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.22.0"), // first version that supports the --rebase-merges option
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.CreateMergeCommit(shell)
		shell.CreateNCommits(2)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			TopLines(
				Contains("CI ◯ commit 02").IsSelected(),
				Contains("CI ◯ commit 01"),
				Contains("Merge branch 'second-change-branch' into first-change-branch"),
			).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.Edit).
			Lines(
				Contains("edit  CI commit 02").IsSelected(),
				Contains("edit  CI commit 01").IsSelected(),
				Contains("      CI ⏣─╮ <-- YOU ARE HERE --- Merge branch 'second-change-branch' into first-change-branch").IsSelected(),
				Contains("      CI │ ◯ * second-change-branch unrelated change"),
				Contains("      CI │ ◯ second change"),
				Contains("      CI ◯ │ first change"),
				Contains("      CI ◯─╯ * original"),
				Contains("      CI ◯ three"),
				Contains("      CI ◯ two"),
				Contains("      CI ◯ one"),
			)
	},
})
