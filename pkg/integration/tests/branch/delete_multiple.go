package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DeleteMultiple = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Try some combinations of local and remote branch deletions with a range selection of branches",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetAppState().LocalBranchSortOrder = "alphabetic"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CloneIntoRemote("origin").
			CloneIntoRemote("other-remote").
			EmptyCommit("blah").
			NewBranch("branch-01").
			EmptyCommit("on branch-01 01").
			PushBranchAndSetUpstream("origin", "branch-01").
			EmptyCommit("on branch-01 02").
			NewBranch("branch-02").
			EmptyCommit("on branch-02 01").
			PushBranchAndSetUpstream("origin", "branch-02").
			NewBranchFrom("branch-03", "master").
			EmptyCommit("on branch-03 01").
			NewBranch("current-head").
			EmptyCommit("on current-head").
			NewBranchFrom("branch-04", "master").
			EmptyCommit("on branch-04 01").
			PushBranchAndSetUpstream("other-remote", "branch-04").
			EmptyCommit("on branch-04 02").
			NewBranchFrom("branch-05", "master").
			EmptyCommit("on branch-05 01").
			PushBranchAndSetUpstream("origin", "branch-05").
			NewBranchFrom("branch-06", "master").
			EmptyCommit("on branch-06 01").
			PushBranch("origin", "branch-06").
			PushBranchAndSetUpstream("other-remote", "branch-06").
			EmptyCommit("on branch-06 02").
			Checkout("current-head")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("current-head").IsSelected(),
				Contains("branch-01 ↑1"),
				Contains("branch-02 ✓"),
				Contains("branch-03"),
				Contains("branch-04 ↑1"),
				Contains("branch-05 ✓"),
				Contains("branch-06 ↑1"),
				Contains("master"),
			).
			Press(keys.Universal.RangeSelectDown).

			// Deleting a range that includes the current branch is not possible
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Tooltip(Contains("You cannot delete the checked out branch!")).
					Title(Equals("Delete selected branches?")).
					Select(Contains("Delete local branches")).
					Confirm().
					Tap(func() {
						t.ExpectToast(Contains("You cannot delete the checked out branch!"))
					}).
					Cancel()
			}).

			// Delete branch-03 and branch-04. 04 is not fully merged, so we get
			// a confirmation popup.
			NavigateToLine(Contains("branch-03")).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Title(Equals("Delete selected branches?")).
					Select(Contains("Delete local branches")).
					Confirm()
				t.ExpectPopup().
					Confirmation().
					Title(Equals("Force delete branch")).
					Content(Equals("Some of the selected branches are not fully merged. Are you sure you want to delete them?")).
					Confirm()
			}).
			Lines(
				Contains("current-head"),
				Contains("branch-01 ↑1"),
				Contains("branch-02 ✓"),
				Contains("branch-05 ✓").IsSelected(),
				Contains("branch-06 ↑1"),
				Contains("master"),
			).

			// Delete remote branches of branch-05 and branch-06. They are on different remotes.
			NavigateToLine(Contains("branch-05")).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Title(Equals("Delete selected branches?")).
					Select(Contains("Delete remote branches")).
					Confirm()
			}).
			Tap(func() {
				t.ExpectPopup().
					Confirmation().
					Title(Equals("Delete selected branches?")).
					Content(Equals("Are you sure you want to delete the remote branches of the selected branches from their respective remotes?")).
					Confirm()
			}).
			Tap(func() {
				checkRemoteBranches(t, keys, "origin", []string{
					"branch-01",
					"branch-02",
					"branch-06",
				})
				checkRemoteBranches(t, keys, "other-remote", []string{
					"branch-04",
				})
			}).
			Lines(
				Contains("current-head"),
				Contains("branch-01 ↑1"),
				Contains("branch-02 ✓"),
				Contains("branch-05 (upstream gone)").IsSelected(),
				Contains("branch-06 (upstream gone)").IsSelected(),
				Contains("master"),
			).

			// Try to delete both local and remote branches of branch-02 and
			// branch-05; not possible because branch-05's upstream is gone
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Title(Equals("Delete selected branches?")).
					Select(Contains("Delete local and remote branches")).
					Confirm().
					Tap(func() {
						t.ExpectToast(Contains("Some of the selected branches have no upstream (or the upstream is not stored locally)"))
					}).
					Cancel()
			}).

			// Delete both local and remote branches of branch-01 and branch-02. We get
			// the force-delete warning because branch-01 it is not fully merged.
			NavigateToLine(Contains("branch-01")).
			Press(keys.Universal.RangeSelectDown).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Title(Equals("Delete selected branches?")).
					Select(Contains("Delete local and remote branches")).
					Confirm()
				t.ExpectPopup().
					Confirmation().
					Title(Equals("Delete local and remote branch")).
					Content(Contains("Are you sure you want to delete both the selected branches from your machine, and their remote branches from their respective remotes?").
						Contains("Some of the selected branches are not fully merged. Are you sure you want to delete them?")).
					Confirm()
			}).
			Lines(
				Contains("current-head"),
				Contains("branch-05 (upstream gone)").IsSelected(),
				Contains("branch-06 (upstream gone)"),
				Contains("master"),
			).
			Tap(func() {
				checkRemoteBranches(t, keys, "origin", []string{
					"branch-06",
				})
				checkRemoteBranches(t, keys, "other-remote", []string{
					"branch-04",
				})
			})
	},
})
