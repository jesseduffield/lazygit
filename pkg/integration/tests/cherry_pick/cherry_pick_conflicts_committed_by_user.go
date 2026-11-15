package cherry_pick

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var CherryPickConflictsCommittedByUser = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cherry pick two commits, the first of which conflicts; the user commits the resolution themselves before continuing",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "recency"
	},
	SetupRepo: func(shell *Shell) {
		shared.MergeConflictsSetup(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("first-change-branch"),
				Contains("second-change-branch"),
				Contains("original-branch"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			TopLines(
				Contains("second-change-branch unrelated change"),
				Contains("second change"),
			).
			Press(keys.Commits.CherryPickCopy).
			SelectNextItem().
			Press(keys.Commits.CherryPickCopy)

		t.Views().Information().Content(Contains("2 commits copied"))

		t.Views().Commits().
			Focus().
			TopLines(
				Contains("first change").IsSelected(),
			).
			Press(keys.Commits.PasteCommits)

		t.ExpectPopup().Alert().
			Title(Equals("Cherry-pick")).
			Content(Contains("Are you sure you want to cherry-pick the 2 copied commit(s) onto this branch?")).
			Confirm()

		t.Common().AcknowledgeConflicts()

		t.Shell().
			UpdateFileAndAdd("file", shared.SecondChangeFileContent).
			RunCommand([]string{"git", "commit", "--no-edit"})

		t.GlobalPress(keys.Universal.Refresh)

		t.Views().Status().Content(Equals("(cherry-picking) repo → first-change-branch"))

		t.Views().Commits().
			TopLines(
				Contains("─── Pending cherry-picks"),
				Contains("pick").DoesNotContain("CONFLICT").Contains("second-change-branch unrelated change"),
				Contains("─── Commits"),
				Contains("second change"),
				Contains("first change"),
			)

		t.Common().ContinueOnConflictsResolved("cherry-pick")

		t.Views().Status().Content(Equals("repo → first-change-branch"))

		t.Views().Commits().
			TopLines(
				Contains("second-change-branch unrelated change"),
				Contains("second change"),
				Contains("first change"),
			)
	},
})
