package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DeleteUpdateRefTodo = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Delete an update-ref item from the rebase todo list",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.38.0"),
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			NewBranch("branch1").
			CreateNCommits(3).
			NewBranch("branch2").
			CreateNCommitsStartingAt(3, 4)

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("commit 01")).
			Press(keys.Universal.Edit).
			Lines(
				Contains("pick").Contains("CI commit 06"),
				Contains("pick").Contains("CI commit 05"),
				Contains("pick").Contains("CI commit 04"),
				Contains("update-ref").Contains("branch1"),
				Contains("pick").Contains("CI commit 03"),
				Contains("pick").Contains("CI commit 02"),
				Contains("CI ◯ <-- YOU ARE HERE --- commit 01"),
			).
			NavigateToLine(Contains("update-ref")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Drop commit")).
					Content(Contains("Are you sure you want to delete the selected update-ref todo(s)?")).
					Confirm()
			}).
			Lines(
				Contains("pick").Contains("CI commit 06"),
				Contains("pick").Contains("CI commit 05"),
				Contains("pick").Contains("CI commit 04"),
				Contains("pick").Contains("CI commit 03").IsSelected(),
				Contains("pick").Contains("CI commit 02"),
				Contains("CI ◯ <-- YOU ARE HERE --- commit 01"),
			).
			NavigateToLine(Contains("commit 02")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.Common().ContinueRebase()
			}).
			Lines(
				Contains("CI ◯ commit 06"),
				Contains("CI ◯ commit 05"),
				Contains("CI ◯ commit 04"),
				Contains("CI ◯ commit 03"), // No star on this commit, so there's no branch head here
				Contains("CI ◯ commit 01"),
			)

		t.Views().Branches().
			Lines(
				Contains("branch2"),
				Contains("branch1"),
			)
	},
})
