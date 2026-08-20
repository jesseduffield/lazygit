package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageUnderConformingDiffRenderer = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A diff renderer that announces the metadata protocol is taken at its word, so its diff is what stays on screen when the main view is focused",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		// A renderer that announces the protocol with a version-only record before
		// anything else, then says who it is and passes the diff through. The diff it
		// passes through keeps its structure, so the rows can be placed by reading it.
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{Command: `printf '\033]1717;1\007RENDERED BY ME\n'; cat`},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nTHREE\nfour\nfive\nsix\nseven\neight\nNINE\nten\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsFocused()

		// The announcement itself leaves nothing on the screen; the line after it shows
		// whose output this is.
		t.Views().Main().Content(Contains("RENDERED BY ME"))

		t.Views().Files().Press(keys.Universal.FocusMainView)

		// Focusing left the renderer's output alone, and the selection went on it.
		t.Views().Main().
			IsFocused().
			Content(Contains("RENDERED BY ME")).
			SelectedLines(
				Contains("-three"),
			).
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(
				Contains("-three"),
				Contains("+THREE"),
			).
			PressPrimaryAction()

		t.Views().Secondary().ContainsLines(
			Contains("-three"),
			Contains("+THREE"),
		)
	},
})
