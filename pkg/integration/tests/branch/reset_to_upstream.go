package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ResetToUpstream = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Hard reset the current branch to the selected branch upstream",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CloneIntoRemote("origin").
			NewBranch("hard-branch").
			EmptyCommit("hard commit").
			PushBranch("origin", "hard-branch").
			NewBranch("soft-branch").
			EmptyCommit("soft commit").
			PushBranch("origin", "soft-branch").
			NewBranch("base").
			EmptyCommit("base-branch commit").
			CreateFile("file-1", "content").
			GitAdd("file-1").
			Commit("commit with file").
			CreateFile("file-2", "content").
			GitAdd("file-2")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// soft reset
		t.Views().Branches().
			Focus().
			Lines(
				Contains("base").IsSelected(),
				Contains("soft-branch"),
				Contains("hard-branch"),
			).
			Press(keys.Branches.SetUpstream).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Upstream options")).
					Select(Contains("Reset checked-out branch onto upstream of selected branch")).
					Tooltip(Contains("Disabled: The selected branch has no upstream (or the upstream is not stored locally)")).
					Confirm().
					Tap(func() {
						t.ExpectToast(Equals("Disabled: The selected branch has no upstream (or the upstream is not stored locally)"))
					}).
					Cancel()
			}).
			SelectNextItem().
			Lines(
				Contains("base"),
				Contains("soft-branch").IsSelected(),
				Contains("hard-branch"),
			).
			Press(keys.Branches.SetUpstream).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Upstream options")).
					Select(Contains("Reset checked-out branch onto origin/soft-branch...")).
					Confirm()

				t.ExpectPopup().Menu().
					Title(Equals("Reset to origin/soft-branch")).
					Select(Contains("Soft reset")).
					Confirm()
			})
		t.Views().Commits().Lines(
			Contains("soft commit"),
			Contains("hard commit"),
		)
		t.Views().Files().Lines(
			Contains("file-1").Contains("A"),
			Contains("file-2").Contains("A"),
		)

		// hard reset
		t.Views().Branches().
			Focus().
			Lines(
				Contains("base"),
				Contains("soft-branch").IsSelected(),
				Contains("hard-branch"),
			).
			NavigateToLine(Contains("hard-branch")).
			Press(keys.Branches.SetUpstream).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Upstream options")).
					Select(Contains("Reset checked-out branch onto origin/hard-branch...")).
					Confirm()

				t.ExpectPopup().Menu().
					Title(Equals("Reset to origin/hard-branch")).
					Select(Contains("Hard reset")).
					Confirm()
			})
		t.Views().Commits().Lines(Contains("hard commit"))
		t.Views().Files().IsEmpty()
	},
})
