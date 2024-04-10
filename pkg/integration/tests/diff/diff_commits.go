package diff

import (
	"github.com/lobes/lazytask/pkg/config"
	. "github.com/lobes/lazytask/pkg/integration/components"
)

var DiffCommits = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "View the diff between two commits",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "first line\n")
		shell.Commit("first commit")
		shell.UpdateFileAndAdd("file1", "first line\nsecond line\n")
		shell.Commit("second commit")
		shell.UpdateFileAndAdd("file1", "first line\nsecond line\nthird line\n")
		shell.Commit("third commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("third commit").IsSelected(),
				Contains("second commit"),
				Contains("first commit"),
			).
			Press(keys.Universal.DiffingMenu).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Diffing")).Select(MatchesRegexp(`Diff \w+`)).Confirm()

				t.Views().Information().Content(Contains("Showing output for: git diff"))
			}).
			SelectNextItem().
			SelectNextItem().
			SelectedLine(Contains("first commit")).
			Tap(func() {
				t.Views().Main().Content(Contains("-second line\n-third line"))
			}).
			Press(keys.Universal.DiffingMenu).
			Tap(func() {
				t.ExpectPopup().Menu().Title(Equals("Diffing")).Select(Contains("Reverse diff direction")).Confirm()

				t.Views().Main().Content(Contains("+second line\n+third line"))
			}).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			SelectedLine(Contains("file1"))

		t.Views().Main().Content(Contains("+second line\n+third line"))
	},
})
