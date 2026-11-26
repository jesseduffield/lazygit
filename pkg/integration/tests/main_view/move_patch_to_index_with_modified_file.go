package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MovePatchToIndexWithModifiedFile = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move a patch from a commit to the index with a conflicting working-tree change",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "1\n2\n3\n4\n")
		shell.Commit("first commit")
		shell.UpdateFileAndAdd("file1", "11\n2\n3\n4\n")
		shell.Commit("second commit")
		shell.UpdateFile("file1", "111\n2\n3\n4\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("second commit").IsSelected(),
				Contains("first commit"),
			).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Equals("M file1"),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("-1")).
			Press(keys.Universal.ToggleRangeSelect).
			Press(keys.Universal.NextItem).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))
		t.Views().Secondary().Content(Contains("-1\n+11"))

		t.Common().SelectPatchOption(Contains("Move patch out into index"))

		t.ExpectPopup().Confirmation().Title(Equals("Must stash")).
			Content(Contains("Pulling a patch out into the index requires stashing and unstashing your changes.")).
			Confirm()

		t.Views().Files().
			Focus().
			Lines(
				Equals("MM file1"),
			)

		t.Views().Main().Content(Contains("-11\n+111\n"))
		t.Views().Secondary().Content(Contains("-1\n+11\n"))
		t.Views().Stash().IsEmpty()
	},
})
