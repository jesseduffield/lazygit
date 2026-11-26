package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ApplyCustomPatch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Apply a custom patch built from a branch commit's focused diff",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch-a")
		shell.CreateFileAndAdd("file1", "first line\n")
		shell.Commit("first commit")

		shell.NewBranch("branch-b")
		shell.UpdateFileAndAdd("file1", "first line\nsecond line\n")
		shell.Commit("update")
		shell.Checkout("branch-a")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("branch-a").IsSelected(),
				Contains("branch-b"),
			).
			Press(keys.Universal.NextItem).
			PressEnter()

		t.Views().SubCommits().
			IsFocused().
			Lines(
				Contains("update").IsSelected(),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("M file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("+second line")).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))
		t.Views().Secondary().Content(Contains("second line"))
		t.Common().SelectPatchOption(MatchesRegexp(`Apply patch$`))

		t.Views().Files().
			Focus().
			Lines(
				Contains("file1").IsSelected(),
			)
		t.Views().Secondary().Content(Contains("second line"))
	},
})
