package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ToggleOverviewMode = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Toggle overview mode to only show merge commits and commits that a ref points to",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		// leave the graph on so that rendering it over the rewritten parents
		// is exercised too
		config.GetUserConfig().Git.Log.ShowGraph = "always"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommitWithDate("base", "2023-04-07 10:00:00").
			NewBranch("feature").
			EmptyCommitWithDate("feature-1", "2023-04-07 11:00:00").
			EmptyCommitWithDate("feature-2", "2023-04-07 12:00:00").
			Checkout("master").
			EmptyCommitWithDate("master-1", "2023-04-07 13:00:00").
			Merge("feature").
			EmptyCommitWithDate("master-2", "2023-04-07 14:00:00").
			CreateLightweightTag("some-tag", "HEAD").
			EmptyCommitWithDate("master-3", "2023-04-07 15:00:00")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("master-3"),
				Contains("master-2"),
				Contains("Merge branch 'feature'"),
				Contains("feature-2"),
				Contains("feature-1"),
				Contains("master-1"),
				Contains("base"),
			).
			Press(keys.Commits.ToggleOverviewMode).
			Lines(
				// the tip of the checked-out branch, the merge commit, the
				// tagged commit, and the tip of the feature branch remain;
				// the graph keeps the shape it has in the full list, with the
				// merge's join line and the feature branch's lane intact
				Contains("master-3"),
				Contains("master-2"),
				Contains("Merge branch 'feature'").Contains("◎─╮"),
				Contains("feature-2").Contains("│ ○"),
			)

		t.Views().Information().Content(Contains("Showing commits overview"))

		t.Views().Commits().
			Title(Equals("Commits (overview)")).
			Press(keys.Commits.ToggleOverviewMode).
			Title(Equals("Commits")).
			Lines(
				Contains("master-3"),
				Contains("master-2"),
				Contains("Merge branch 'feature'"),
				Contains("feature-2"),
				Contains("feature-1"),
				Contains("master-1"),
				Contains("base"),
			)

		t.Views().Information().Content(DoesNotContain("Showing commits overview"))

		// collapsing while on a commit that gets hidden moves the selection to
		// the nearest visible commit above it, and expanding again jumps back
		t.Views().Commits().
			NavigateToLine(Contains("feature-1")).
			Press(keys.Commits.ToggleOverviewMode).
			Lines(
				Contains("master-3"),
				Contains("master-2"),
				Contains("Merge branch 'feature'"),
				Contains("feature-2").IsSelected(),
			).
			Press(keys.Commits.ToggleOverviewMode).
			Lines(
				Contains("master-3"),
				Contains("master-2"),
				Contains("Merge branch 'feature'"),
				Contains("feature-2"),
				Contains("feature-1").IsSelected(),
				Contains("master-1"),
				Contains("base"),
			)

		// ... but not if the selection was moved while collapsed
		t.Views().Commits().
			NavigateToLine(Contains("master-1")).
			Press(keys.Commits.ToggleOverviewMode).
			Lines(
				Contains("master-3"),
				Contains("master-2"),
				Contains("Merge branch 'feature'"),
				Contains("feature-2").IsSelected(),
			).
			NavigateToLine(Contains("master-3")).
			Press(keys.Commits.ToggleOverviewMode).
			Lines(
				Contains("master-3").IsSelected(),
				Contains("master-2"),
				Contains("Merge branch 'feature'"),
				Contains("feature-2"),
				Contains("feature-1"),
				Contains("master-1"),
				Contains("base"),
			)

		// commands that rebase relative to the selection's position in the
		// list require exiting overview mode first
		t.Views().Commits().
			Press(keys.Commits.ToggleOverviewMode).
			Lines(
				Contains("master-3"),
				Contains("master-2"),
				Contains("Merge branch 'feature'"),
				Contains("feature-2"),
			).
			Press(keys.Commits.SquashDown)

		t.ExpectPopup().Confirmation().
			Title(Equals("Command not available")).
			Content(Equals("Command not available in overview mode. Exit overview mode?")).
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("master-3"),
				Contains("master-2"),
				Contains("Merge branch 'feature'"),
				Contains("feature-2"),
				Contains("feature-1"),
				Contains("master-1"),
				Contains("base"),
			)

		t.Views().Information().Content(DoesNotContain("Showing commits overview"))
	},
})
