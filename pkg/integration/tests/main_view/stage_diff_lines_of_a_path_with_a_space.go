package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageDiffLinesOfAPathWithASpace = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage a line of a file whose path contains a space from the focused main view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		// git ends the path field of the diff's "---" and "+++" lines with a tab
		// when the path contains a space, and the view shows the tab as spaces.
		shell.CreateFileAndAdd("my file", "one\ntwo\nthree\n")
		shell.Commit("one")

		shell.UpdateFile("my file", "one\nADD1\nADD2\ntwo\nthree\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("my file").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(
				Contains("+ADD1"),
			).
			PressPrimaryAction()

		t.Views().Files().Lines(
			Contains("MM my file"),
		)
		t.Views().Secondary().
			ContainsLines(
				Contains("+ADD1"),
			).
			Content(DoesNotContain("+ADD2"))
		t.Views().Main().
			Content(DoesNotContain("+ADD1")).
			ContainsLines(
				Contains("+ADD2"),
			)
	},
})
