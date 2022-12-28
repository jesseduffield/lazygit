package diff

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Diff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "View the diff of two branches, then view the reverse diff",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch-a")
		shell.CreateFileAndAdd("file1", "first line")
		shell.Commit("first commit")

		shell.NewBranch("branch-b")
		shell.UpdateFileAndAdd("file1", "first line\nsecond line")
		shell.Commit("update")

		shell.Checkout("branch-a")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			TopLines(
				Contains("branch-a"),
				Contains("branch-b"),
			).
			Press(keys.Universal.DiffingMenu)

		t.ExpectPopup().Menu().Title(Equals("Diffing")).Select(Contains(`diff branch-a`)).Confirm()

		t.Views().Branches().
			IsFocused().
			Tap(func() {
				t.Views().Information().Content(Contains("showing output for: git diff branch-a branch-a"))
			}).
			SelectNextItem().
			Tap(func() {
				t.Views().Information().Content(Contains("showing output for: git diff branch-a branch-b"))
				t.Views().Main().Content(Contains("+second line"))
			}).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			SelectedLine(Contains("update")).
			Tap(func() {
				t.Views().Main().Content(Contains("+second line"))
			}).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			SelectedLine(Contains("file1")).
			Tap(func() {
				t.Views().Main().Content(Contains("+second line"))
			}).
			PressEscape()

		t.Views().SubCommits().PressEscape()

		t.Views().Branches().
			IsFocused().
			Press(keys.Universal.DiffingMenu)

		t.ExpectPopup().Menu().Title(Equals("Diffing")).Select(Contains("reverse diff direction")).Confirm()
		t.Views().Information().Content(Contains("showing output for: git diff branch-a branch-b -R"))
		t.Views().Main().Content(Contains("-second line"))
	},
})
