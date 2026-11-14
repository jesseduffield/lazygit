package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SquashMergeDuringRebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Squash merge a branch with conflicts at the base of an interactive rebase, then continue the merge and the rebase",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "alphabetical"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("file", "original\n").
			Commit("original").
			NewBranch("other").
			UpdateFileAndAdd("file", "other change\n").
			Commit("other change").
			Checkout("master").
			UpdateFileAndAdd("file", "master change\n").
			Commit("master change").
			NewBranch("mybranch").
			CreateFileAndAdd("file2", "content\n").
			Commit("mybranch 1").
			CreateFileAndAdd("file3", "content\n").
			Commit("mybranch 2")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("mybranch 2").IsSelected(),
				Contains("mybranch 1"),
				Contains("master change"),
				Contains("original"),
			).
			Press(keys.Commits.StartInteractiveRebase).
			Lines(
				Contains("─── Pending rebase todos"),
				Contains("mybranch 2").IsSelected(),
				Contains("mybranch 1"),
				Contains("─── Commits"),
				Contains("master change"),
				Contains("original"),
			)

		t.Views().Branches().
			Focus().
			Lines(
				Contains("(no branch, rebasing mybranch)"),
				Contains("master"),
				Contains("mybranch").IsSelected(),
				Contains("other"),
			).
			NavigateToLine(Contains("other")).
			Press(keys.Branches.MergeIntoCurrentBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Merge")).
			Select(Contains("Squash merge and commit")).
			Tooltip(Equals("Squash merge 'other' into 'mybranch' as a single commit.")).
			Confirm()

		t.Common().AcknowledgeConflicts()

		t.Views().Status().
			/* EXPECTED:
			Content(Equals("(merging) repo → mybranch"))
			ACTUAL: */
			Content(MatchesRegexp(`^\(merging\) repo → [0-9a-f]{40}$`))

		t.GlobalPress(keys.Universal.CreateRebaseOptionsMenu)
		t.ExpectPopup().Menu().
			Title(Equals("Merge options")).
			Cancel()

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("UU file").IsSelected(),
			).
			Tap(func() {
				t.Shell().UpdateFile("file", "resolved\n")
			}).
			Press(keys.Universal.Refresh)

		t.Common().ContinueOnConflictsResolved("merge")

		t.Views().Status().
			/* EXPECTED:
			Content(Equals("(rebasing) repo → mybranch"))
			ACTUAL: */
			Content(MatchesRegexp(`^\(rebasing\) repo → [0-9a-f]{40}$`))

		t.Views().Commits().
			Lines(
				Contains("─── Pending rebase todos"),
				Contains("mybranch 2"),
				Contains("mybranch 1"),
				Contains("─── Commits"),
				Contains("Squash merge other into mybranch"),
				Contains("master change"),
				Contains("original"),
			)

		t.Common().ContinueRebase()

		t.Views().Status().Content(Equals("repo → mybranch"))

		t.Views().Commits().
			Lines(
				Contains("mybranch 2"),
				Contains("mybranch 1"),
				Contains("Squash merge other into mybranch"),
				Contains("master change"),
				Contains("original"),
			)
	},
})
