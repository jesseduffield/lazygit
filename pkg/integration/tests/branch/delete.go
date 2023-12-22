package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Delete = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Try all combination of local and remote branch deletions",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CloneIntoRemote("origin").
			EmptyCommit("blah").
			NewBranch("branch-one").
			PushBranch("origin", "branch-one").
			NewBranch("branch-two").
			PushBranch("origin", "branch-two").
			EmptyCommit("deletion blocker").
			NewBranch("branch-three")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				MatchesRegexp(`\*.*branch-three`).IsSelected(),
				MatchesRegexp(`branch-two`),
				MatchesRegexp(`branch-one`),
				MatchesRegexp(`master`),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Tooltip(Contains("You cannot delete the checked out branch!")).
					Title(Equals("Delete branch 'branch-three'?")).
					Select(Contains("Delete local branch")).
					Confirm().
					Tap(func() {
						t.ExpectToast(Contains("You cannot delete the checked out branch!"))
					}).
					Cancel()
			}).
			SelectNextItem().
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Title(Equals("Delete branch 'branch-two'?")).
					Select(Contains("Delete local branch")).
					Confirm()
			}).
			Tap(func() {
				t.ExpectPopup().
					Confirmation().
					Title(Equals("Force delete branch")).
					Content(Equals("'branch-two' is not fully merged. Are you sure you want to delete it?")).
					Confirm()
			}).
			Lines(
				MatchesRegexp(`\*.*branch-three`),
				MatchesRegexp(`branch-one`).IsSelected(),
				MatchesRegexp(`master`),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Title(Equals("Delete branch 'branch-one'?")).
					Select(Contains("Delete remote branch")).
					Confirm()
			}).
			Tap(func() {
				t.ExpectPopup().
					Confirmation().
					Title(Equals("Delete branch 'branch-one'?")).
					Content(Equals("Are you sure you want to delete the remote branch 'branch-one' from 'origin'?")).
					Confirm()
			}).
			Tap(func() {
				t.Views().Remotes().
					Focus().
					Lines(Contains("origin")).
					PressEnter()

				t.Views().
					RemoteBranches().
					Lines(Equals("branch-two")).
					Press(keys.Universal.Return)

				t.Views().
					Branches().
					Focus()
			}).
			Lines(
				MatchesRegexp(`\*.*branch-three`),
				MatchesRegexp(`branch-one \(upstream gone\)`).IsSelected(),
				MatchesRegexp(`master`),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().
					Menu().
					Title(Equals("Delete branch 'branch-one'?")).
					Select(Contains("Delete local branch")).
					Confirm()
			}).
			Lines(
				MatchesRegexp(`\*.*branch-three`),
				MatchesRegexp(`master`).IsSelected(),
			)
	},
})
