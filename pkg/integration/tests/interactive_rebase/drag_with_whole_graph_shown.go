package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DragWithWholeGraphShown = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Drag a commit past another branch's commit while the whole graph is shown",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.Log.Order = "date-order"
		config.GetUserConfig().Git.Log.ShowWholeGraph = true
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommitWithDate("base", "2026-01-01T00:00:00Z").
			NewBranch("dev-a").
			EmptyCommitWithDate("dev-a-1", "2026-01-01T00:02:00Z").
			NewBranchFrom("dev-b", "master").
			EmptyCommitWithDate("dev-b-1", "2026-01-01T00:01:00Z").
			EmptyCommitWithDate("dev-b-2", "2026-01-01T00:03:00Z").
			EmptyCommitWithDate("dev-b-3", "2026-01-01T00:04:00Z")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("dev-b-2")).
			Lines(
				Contains("dev-b-3"),
				Contains("dev-b-2").IsSelected(),
				Contains("dev-a-1"),
				Contains("dev-b-1"),
				Contains("base"),
			).
			ClickAndHold(1, 1).
			MouseMove(1, 3).
			Lines(
				Contains("dev-b-3"),
				Contains("dev-b-2").IsSelected(),
				Contains("dev-a-1"),
				Contains("dev-b-1"),
				Contains("drop here"),
				Contains("base"),
			).
			MouseRelease()

		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("dev-b")).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			/* EXPECTED:
			Lines(
				Contains("dev-b-3"),
				Contains("dev-b-1"),
				Contains("dev-b-2"),
				Contains("base"),
			)
			ACTUAL: */
			Lines(
				Contains("dev-b-3"),
				Contains("dev-b-1"),
				Contains("base"),
				Contains("dev-b-2"),
			)
	},
})
