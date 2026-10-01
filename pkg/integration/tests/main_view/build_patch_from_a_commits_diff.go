package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var BuildPatchFromACommitsDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Take lines of a commit's diff into a custom patch, and back out of it, from the focused main view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file1", "ONE\ntwo\nTHREE\nfour\nfive\n")
		shell.Commit("second commit")
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
				Contains("file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-one"),
			).
			// The line goes into the patch, which the pane beside the diff previews, and
			// the selection moves on to the next change rather than staying where a
			// second press would take it straight back out.
			PressPrimaryAction().
			SelectedLines(
				Contains("+ONE"),
			)

		t.Views().Information().Content(Contains("Building patch"))
		t.Views().Secondary().ContainsLines(
			Contains("-one"),
			Contains(" two"),
		)
		// The line that is in the patch is marked as such over the diff itself.
		t.Views().Main().MarkedLines(
			Contains("-one"),
		)

		// The addition of the same modification goes in too, and the patch holds both.
		t.Views().Main().
			IsFocused().
			PressPrimaryAction().
			SelectedLines(
				Contains("-three"),
			)

		t.Views().Secondary().ContainsLines(
			Contains("-one"),
			Contains("+ONE"),
			Contains(" two"),
		)
		t.Views().Main().MarkedLines(
			Contains("-one"),
			Contains("+ONE"),
		)

		// Pointing at a line that is in the patch takes it back out.
		t.Views().Main().
			IsFocused().
			NavigateToLine(Contains("+ONE")).
			PressPrimaryAction()

		t.Views().Secondary().
			ContainsLines(
				Contains("-one"),
				Contains(" two"),
			).
			Content(DoesNotContain("+ONE"))
		t.Views().Main().MarkedLines(
			Contains("-one"),
		)

		// Taking the last line out ends the patch, so the pane previewing it goes away.
		t.Views().Main().
			IsFocused().
			NavigateToLine(Contains("-one")).
			PressPrimaryAction()

		t.Views().Information().Content(DoesNotContain("Building patch"))
		t.Views().Main().NoMarkedLines()
	},
})
