package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var UnsetUpstream = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Unset upstream of selected branch, both when it exists and when it doesn't",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("one").
			NewBranch("branch_to_remove").
			Checkout("master").
			CloneIntoRemote("origin").
			SetBranchUpstream("master", "origin/master").
			SetBranchUpstream("branch_to_remove", "origin/branch_to_remove").
			// to get the "(upstream gone)" branch status
			RunCommand([]string{"git", "push", "origin", "--delete", "branch_to_remove"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Press(keys.Universal.NextScreenMode). // we need to enlargen the window to see the upstream
			SelectedLines(
				Contains("master").Contains("origin master"),
			).
			Press(keys.Branches.SetUpstream).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Upstream options")).
					Select(Contains("Unset upstream of selected branch")).
					Confirm()
			}).
			SelectedLines(
				Contains("master").DoesNotContain("origin master"),
			)

		t.Views().Branches().
			Focus().
			SelectNextItem().
			SelectedLines(
				Contains("branch_to_remove").Contains("origin branch_to_remove").Contains("upstream gone"),
			).
			Press(keys.Branches.SetUpstream).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Upstream options")).
					Select(Contains("Unset upstream of selected branch")).
					Confirm()
			}).
			SelectedLines(
				Contains("branch_to_remove").DoesNotContain("origin branch_to_remove").DoesNotContain("upstream gone"),
			)
	},
})
