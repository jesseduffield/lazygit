package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RevertEmptyCommitResolution = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Handles a revert whose commit becomes empty and offers skip/create options",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("myfile", "")
		shell.Commit("add empty file")
		shell.CreateFileAndAdd("myfile", "first line\n")
		shell.Commit("add first line")
		shell.UpdateFileAndAdd("myfile", "first line\nsecond line\n")
		shell.Commit("add second line")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI ◯ add second line").IsSelected(),
				Contains("CI ◯ add first line"),
				Contains("CI ◯ add empty file"),
			).
			SelectNextItem().
			Press(keys.Commits.RevertCommit).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Revert commit")).
					Content(MatchesRegexp(`Are you sure you want to revert \w+?`)).
					Confirm()
				t.Common().AcknowledgeConflicts()
			})

		t.Shell().
			UpdateFile("myfile", "first line\nsecond line\n").
			RunCommand([]string{"git", "add", "myfile"})

		t.Views().Commits().Focus()
		t.Common().ContinueRebase()

		t.ExpectPopup().Menu().
			Title(Equals("Commit produced no changes")).
			Lines(
				Contains("Skip this revert step"),
				Contains("Create empty commit and continue"),
			).
			Select(Contains("Create empty commit and continue")).
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("CI ◯ add second line").IsSelected(),
				Contains("CI ◯ add first line"),
				Contains("CI ◯ add empty file"),
			)

		t.Views().Commits().Content(DoesNotContain("Pending reverts"))
		t.Views().Options().Content(DoesNotContain("View revert options")).
			Content(DoesNotContain("You are currently neither rebasing nor merging"))
	},
})
