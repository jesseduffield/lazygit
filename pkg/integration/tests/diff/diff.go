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
	Run: func(shell *Shell, input *Input, keys config.KeybindingConfig) {
		input.Views().Branches().
			Focus().
			TopLines(
				Contains("branch-a"),
				Contains("branch-b"),
			).
			Press(keys.Universal.DiffingMenu)

		input.ExpectMenu().Title(Equals("Diffing")).Select(Contains(`diff branch-a`)).Confirm()

		input.Views().Branches().
			IsFocused().
			Tap(func() {
				input.Views().Information().Content(Contains("showing output for: git diff branch-a branch-a"))
			}).
			SelectNextItem().
			Tap(func() {
				input.Views().Information().Content(Contains("showing output for: git diff branch-a branch-b"))
				input.Views().Main().Content(Contains("+second line"))
			}).
			PressEnter()

		input.Views().SubCommits().
			IsFocused().
			SelectedLine(Contains("update")).
			Tap(func() {
				input.Views().Main().Content(Contains("+second line"))
			}).
			PressEnter()

		input.Views().CommitFiles().
			IsFocused().
			SelectedLine(Contains("file1")).
			Tap(func() {
				input.Views().Main().Content(Contains("+second line"))
			}).
			PressEscape()

		input.Views().SubCommits().PressEscape()

		input.Views().Branches().
			IsFocused().
			Press(keys.Universal.DiffingMenu)

		input.ExpectMenu().Title(Equals("Diffing")).Select(Contains("reverse diff direction")).Confirm()
		input.Views().Information().Content(Contains("showing output for: git diff branch-a branch-b -R"))
		input.Views().Main().Content(Contains("-second line"))
	},
})
