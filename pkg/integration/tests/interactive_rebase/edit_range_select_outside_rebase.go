package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var EditRangeSelectOutsideRebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Select a range of commits to edit outside of a rebase",
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
			TopLines(
				Contains("Merge branch 'second-change-branch' into first-change-branch").IsSelected(),
			).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.RangeSelectDown).
			Lines(
				Contains("CI ⏣─╮ Merge branch 'second-change-branch' into first-change-branch").IsSelected(),
				Contains("CI │ ◯ * second-change-branch unrelated change").IsSelected(),
				Contains("CI │ ◯ second change").IsSelected(),
				Contains("CI ◯ │ first change").IsSelected(),
				Contains("CI ◯─╯ * original").IsSelected(),
				Contains("CI ◯ three").IsSelected(),
				Contains("CI ◯ two"),
				Contains("CI ◯ one"),
			).
			Press(keys.Universal.Edit).
			Lines(
				Contains("merge  CI Merge branch 'second-change-branch' into first-change-branch").IsSelected(),
				Contains("edit   CI first change").IsSelected(),
				Contains("edit   CI * second-change-branch unrelated change").IsSelected(),
				Contains("edit   CI second change").IsSelected(),
				Contains("edit   CI * original").IsSelected(),
				Contains("       CI ◯ <-- YOU ARE HERE --- three").IsSelected(),
				Contains("       CI ◯ two"),
				Contains("       CI ◯ one"),
			)
	},
})
