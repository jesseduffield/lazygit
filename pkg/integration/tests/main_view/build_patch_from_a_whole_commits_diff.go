package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var BuildPatchFromAWholeCommitsDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Take lines of two files into a custom patch from the whole diff of a commit, without entering its files first",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.CreateFileAndAdd("file2", "alpha\nbeta\ngamma\n")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file1", "one\nTWO\nthree\n")
		shell.UpdateFileAndAdd("file2", "alpha\nBETA\ngamma\n")
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("second commit").IsSelected(),
				Contains("first commit"),
			).
			Press(keys.Universal.FocusMainView)

		// The commit's whole diff spans both files, and hunk mode offers the first
		// changed block of the first of them.
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-two"),
				Contains("+TWO"),
			).
			PressPrimaryAction().
			// The selection moves past the block just taken in, which is file2's block.
			SelectedLines(
				Contains("-beta"),
				Contains("+BETA"),
			).
			PressPrimaryAction()

		// The patch spans both files.
		t.Views().Secondary().
			Content(Contains("file1")).
			Content(Contains("file2")).
			ContainsLines(
				Contains("-two"),
				Contains("+TWO"),
				Contains(" three"),
			).
			ContainsLines(
				Contains("-beta"),
				Contains("+BETA"),
				Contains(" gamma"),
			)

		// Applying it to the working tree puts both files' changes there, which is the
		// proof that the patch really holds what the preview says it does.
		t.Common().SelectPatchOption(Contains("Apply patch in reverse"))

		t.Views().Files().
			Focus().
			ContainsLines(
				Contains("M").Contains("file1"),
				Contains("M").Contains("file2"),
			)
	},
})
