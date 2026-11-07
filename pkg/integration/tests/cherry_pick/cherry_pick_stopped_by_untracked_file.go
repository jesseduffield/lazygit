package cherry_pick

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CherryPickStoppedByUntrackedFile = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cherry pick two commits; an untracked file stops the second one before git starts to apply it; continue after removing the file",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.LocalBranchSortOrder = "recency"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFileAndAdd("file", "original\n").
			Commit("base").
			NewBranch("source").
			CreateFileAndAdd("file1", "one\n").
			Commit("one").
			CreateFileAndAdd("file2", "two\n").
			Commit("two").
			Checkout("master").
			CreateFile("file2", "untracked\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("master"),
				Contains("source"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("two").IsSelected(),
				Contains("one"),
				Contains("base"),
			).
			Press(keys.Commits.CherryPickCopy).
			SelectNextItem().
			Press(keys.Commits.CherryPickCopy)

		t.Views().Information().Content(Contains("2 commits copied"))

		t.Views().Commits().
			Focus().
			Lines(
				Contains("base").IsSelected(),
			).
			Press(keys.Commits.PasteCommits).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Cherry-pick")).
					Content(Contains("Are you sure you want to cherry-pick the 2 copied commit(s) onto this branch?")).
					Confirm()
			}).
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Error")).
					Content(Contains("The following untracked working tree files would be overwritten by merge")).
					Confirm()
			}).
			/* EXPECTED:
			Lines(
				Contains("─── Pending cherry-picks"),
				Contains("pick").DoesNotContain("CONFLICT").Contains("two"),
				Contains("─── Commits"),
				Contains("one"),
				Contains("base"),
			)
			ACTUAL: */
			Lines(
				Contains("one"),
				Contains("base"),
			)

		t.Views().Status().
			/* EXPECTED:
			Content(Equals("(cherry-picking) repo → master"))
			ACTUAL: */
			Content(Equals("repo → master"))

		t.Shell().DeleteFile("file2")

		t.GlobalPress(keys.Universal.CreateRebaseOptionsMenu)
		/* EXPECTED:
		t.ExpectPopup().Menu().
			Title(Equals("Cherry-pick options")).
			Select(Contains("continue")).
			Confirm()

		t.Views().Status().Content(Equals("repo → master"))

		t.Views().Commits().
			Lines(
				Contains("two"),
				Contains("one"),
				Contains("base"),
			)
		ACTUAL: */
		t.ExpectToast(Equals("Disabled: You are currently neither rebasing nor merging"))
	},
})
