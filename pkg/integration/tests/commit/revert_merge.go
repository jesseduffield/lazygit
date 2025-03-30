package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var RevertMerge = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Reverts a merge commit and chooses to revert to the parent commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.CreateMergeCommit(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().Focus().
			TopLines(
				Contains("Merge branch 'second-change-branch' into first-change-branch").IsSelected(),
			).
			Press(keys.Commits.RevertCommit)

		t.ExpectPopup().Confirmation().
			Title(Equals("Revert commit")).
			Content(MatchesRegexp(`Are you sure you want to revert \w+?`)).
			Confirm()

		t.Views().Commits().IsFocused().
			TopLines(
				Contains("Revert \"Merge branch 'second-change-branch' into first-change-branch\""),
				Contains("Merge branch 'second-change-branch' into first-change-branch").IsSelected(),
			).
			SelectPreviousItem()

		t.Views().Main().Content(Contains("-Second Change").Contains("+First Change"))
	},
})
