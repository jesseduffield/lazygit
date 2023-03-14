package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var RevertMerge = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Reverts a merge commit and chooses to revert to the parent commit",
	ExtraCmdArgs: "",
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

		t.ExpectPopup().Menu().
			Title(Equals("Select parent commit for merge")).
			Lines(
				Contains("first change"),
				Contains("second-change-branch unrelated change"),
				Contains("cancel"),
			).
			Select(Contains("first change")).
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
