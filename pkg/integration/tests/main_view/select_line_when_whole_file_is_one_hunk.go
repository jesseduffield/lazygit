package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectLineWhenWholeFileIsOneHunk = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Hunk mode falls back to a single line for a file that is one solid block of changes, rather than selecting all of it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = true
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.CreateFileAndAdd("added", "one\ntwo\nthree\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("added").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		// Every line of the file is an addition, so widening to the change block would
		// select the file entire; hunk mode gives way to a single line.
		t.Views().Secondary().
			IsFocused().
			SelectedLines(
				Contains("+one"),
			).
			// Toggling hunk mode on explicitly still selects the whole block: the fallback
			// is about what the default does, not about forbidding the selection.
			Press(keys.Main.ToggleSelectHunk).
			SelectedLines(
				Contains("+one"),
				Contains("+two"),
				Contains("+three"),
			)
	},
})
