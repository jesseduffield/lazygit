package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RemoveLinesFromTheCustomPatch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Take a line back out of the custom patch from the pane previewing it, where the patch's own numbering is not the commit's",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\n")
		shell.Commit("first commit")

		// Additions with a line of the file between them, so that leaving the first of
		// them out of the patch puts the others at line numbers the commit's diff has
		// unchanged lines at: a line of the patch can only be found by counting the
		// patch's own changes.
		shell.UpdateFileAndAdd("file1", "one\nadded a\ntwo\nadded b\nthree\nadded c\nfour\n")
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Universal.FocusMainView)

		// Take the second and third additions into the patch, leaving the first out.
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+added a"),
			).
			NavigateToLine(Contains("+added b")).
			Press(keys.Universal.ToggleRangeSelect).
			NavigateToLine(Contains("+added c")).
			PressPrimaryAction().
			MarkedLines(
				Contains("+added b"),
				Contains("+added c"),
			)

		t.Views().Secondary().ContainsLines(
			Contains("+added b"),
			Contains(" three"),
			Contains("+added c"),
		)

		// Point at the first of the patch's lines and take it back out.
		t.Views().Main().Press(keys.Universal.TogglePanel)

		t.Views().Secondary().
			IsFocused().
			SelectedLines(
				Contains("+added b"),
			).
			PressPrimaryAction().
			// What is left is the line that wasn't pointed at, and the selection has
			// stayed with it.
			SelectedLines(
				Contains("+added c"),
			).
			Content(DoesNotContain("+added b"))

		t.Views().Main().MarkedLines(
			Contains("+added c"),
		)

		// And the patch really is only that line: applying it to the working tree brings
		// back nothing else.
		t.Common().SelectPatchOption(Contains("Apply patch in reverse"))

		t.Views().Files().
			Focus().
			Lines(
				Contains("M").Contains("file1"),
			)

		// The patch went to the index as well as the working tree, so the file's changes
		// are on the staged side of its diff.
		t.Views().Secondary().
			ContainsLines(
				Contains("-added c"),
			).
			Content(DoesNotContain("+added b"))
	},
})
