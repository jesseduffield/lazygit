package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CustomPatchGoesThroughTheDiffRenderer = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The custom patch is shown by the configured diff renderer, being rendered as a diff of real files rather than assembled by us",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		// A renderer that announces the metadata protocol — so that focusing the main
		// view keeps its output — and says who it is above the diff it passes through.
		cfg.GetUserConfig().Git.DiffRenderers = []config.DiffRendererConfig{
			{Command: `printf '\033]1717;1\007RENDERED BY ME\n'; cat`},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("first commit")

		shell.UpdateFileAndAdd("file1", "one\nTWO\nthree\n")
		shell.Commit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Content(Contains("RENDERED BY ME")).
			SelectedLines(
				Contains("-two"),
			).
			PressPrimaryAction()

		// The patch is a diff like any other, so the renderer has had it too — and git
		// worked out its context, which is what gives it the unchanged line either side.
		t.Views().Secondary().
			Content(Contains("RENDERED BY ME")).
			ContainsLines(
				Contains(" one"),
				Contains("-two"),
				Contains(" three"),
			)
	},
})
