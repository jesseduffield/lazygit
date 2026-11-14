package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MergeIntoDetachedHead = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Merge and squash merge branches into a detached HEAD",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("one").
			NewBranch("branch1").
			EmptyCommit("branch1").
			Checkout("master").
			NewBranch("branch2").
			CreateFileAndAdd("file", "content").
			Commit("branch2").
			Checkout("master").
			EmptyCommit("two").
			Checkout("master^0")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().
			/* EXPECTED:
			Content(MatchesRegexp(`^repo → \(HEAD detached at [0-9a-f]+\)$`))
			ACTUAL: */
			Content(MatchesRegexp(`^repo → [0-9a-f]{40}$`))

		t.Views().Branches().
			Focus().
			Lines(
				Contains("HEAD detached").IsSelected(),
				Contains("branch1"),
				Contains("branch2"),
				Contains("master"),
			).
			NavigateToLine(Contains("branch1")).
			Press(keys.Branches.MergeIntoCurrentBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Merge")).
			Select(Contains("Regular merge (with merge commit)")).
			Tooltip(Equals("Merge 'branch1' into 'HEAD', creating a merge commit.")).
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("◎─╮ Merge branch 'branch1' into HEAD").IsSelected(),
				Contains("│ ○ * branch1"),
				Contains("○ │ two"),
				Contains("○─╯ one"),
			)

		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("branch2")).
			Press(keys.Branches.MergeIntoCurrentBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Merge")).
			Select(Contains("Squash merge and commit")).
			Tooltip(Equals("Squash merge 'branch2' into 'HEAD' as a single commit.")).
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("Squash merge branch2 into HEAD").IsSelected(),
				Contains("Merge branch 'branch1' into HEAD"),
				Contains("branch1"),
				Contains("two"),
				Contains("one"),
			)

		t.Views().Branches().
			Lines(
				Contains("HEAD detached"),
				Contains("branch1"),
				Contains("branch2"),
				Contains("master"),
			)
	},
})
