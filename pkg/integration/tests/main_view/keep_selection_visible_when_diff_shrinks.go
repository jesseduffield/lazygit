package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeepSelectionVisibleWhenDiffShrinks = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The selection stays on the content when a re-render leaves the diff with fewer lines than the selection was on",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{Name: "plain", Command: `printf '\033]1717;1\007'; cat`},
			// The same diff in fewer lines, as a renderer that collapses or elides
			// parts of it would give us: the addition at the end goes, and the hunk
			// header says so, since a diff that contradicts its own header can't be
			// read as one. (It has to read all of its input: one that exits early
			// leaves the render looking like it is still loading, which holds off the
			// clamping this test is about.)
			{Name: "shrinking", Command: `printf '\033]1717;1\007'; ` +
				`sed -e 's/@@ -1,5 +1,5 @@/@@ -1,5 +1,4 @@/' -e '$d'`},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nTHREE\nfour\nFIVE\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Press(keys.Universal.GotoBottom).
			SelectedLines(
				Contains("+FIVE"),
			).
			Press(keys.Universal.CycleDiffRenderers).
			Tap(func() {
				t.ExpectToast(Equals("Diff renderer: shrinking (2 of 2)"))
			})

		// The selection has nowhere to be but the last line there is. Asserting on the
		// index first waits for that to happen: the re-render and the clamp that
		// follows it are a frame apart, and reading the selected line's text in
		// between would be reading past the content.
		t.Views().Main().
			SelectionIsActive().
			SelectedLineIdx(10).
			SelectedLines(
				Contains("-five"),
			)
	},
})
