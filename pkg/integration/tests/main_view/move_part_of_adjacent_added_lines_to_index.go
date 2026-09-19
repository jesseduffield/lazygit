package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MovePartOfAdjacentAddedLinesToIndex = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move only one of two adjacent added lines from a commit to the index",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file1", "1st line\n2nd line\n")
		shell.Commit("commit to move from")

		shell.UpdateFileAndAdd("unrelated-file", "")
		shell.Commit("third commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("third commit").IsSelected(),
				Contains("commit to move from"),
				Contains("first commit"),
			).
			SelectNextItem().
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("+1st line")).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))
		t.Common().SelectPatchOption(Contains("Move patch out into index"))

		t.Views().Main().
			IsFocused().
			Content(Contains("+2nd line").DoesNotContain("1st line"))
		t.Views().Files().
			Focus().
			ContainsLines(
				Contains("M").Contains("file1"),
			)
		t.Views().Secondary().Content(Contains("+1st line\n 2nd line"))
	},
})
