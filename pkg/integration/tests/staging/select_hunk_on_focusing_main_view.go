package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectHunkOnFocusingMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "When hunk mode is the default, focusing the main view selects the first whole hunk, ready to stage",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = true
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nTHREE\nfour\nfive\nsix\nseven\neight\nNINE\nten\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		// No key press: the first hunk is selected just by focusing, so space
		// stages it straight away.
		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("-three"),
				Contains("+THREE"),
			).
			PressPrimaryAction().
			Tap(func() {
				t.Views().Secondary().
					ContainsLines(
						Contains("-three"),
						Contains("+THREE"),
					)
			})
	},
})
