package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var MovePartialPatchToIndex = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Move part of a file's changes from a commit to the index",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "first line\nsecond line\nthird line\n")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file1", "first line2\nsecond line\nthird line2\n")
		shell.Commit("second commit")

		shell.CreateFileAndAdd("file2", "file1 content")
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
			NavigateToLine(Contains("second commit")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains(`-first line`)).
			Press(keys.Universal.ToggleRangeSelect).
			Press(keys.Universal.NextItem).
			SelectedLines(
				Contains(`-first line`),
				Contains(`+first line2`),
			).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))
		t.Views().Secondary().ContainsLines(
			Contains(`-first line`),
			Contains(`+first line2`),
			Contains(` second line`),
			Contains(` third line`),
		)

		t.Common().SelectPatchOption(Contains("Move patch out into index"))

		t.Views().Files().Lines(
			Contains("M").Contains("file1"),
		)
		t.Views().Main().
			IsFocused().
			ContainsLines(
				Contains(` first line`),
				Contains(` second line`),
				Contains(`-third line`),
				Contains(`+third line2`),
			)

		t.Views().Files().Focus()
		t.Views().Secondary().ContainsLines(
			Contains(`-first line`),
			Contains(`+first line2`),
			Contains(` second line`),
			Contains(` third line2`),
		)
	},
})
