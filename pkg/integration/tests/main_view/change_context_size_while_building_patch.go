package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChangeContextSizeWhileBuildingPatch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Change diff context size while building a custom patch, then add another line",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInDiffView = false
		config.GetUserConfig().Git.DiffContextSize = 1
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nten\n")
		shell.Commit("first commit")
		shell.UpdateFileAndAdd("file1", "ONE\ntwo\nthree\nfour\nfive\nsix\nseven\neight\nnine\nTEN\n")
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
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectedLines(Contains("-one")).
			PressPrimaryAction().
			Press(keys.Universal.IncreaseContextInDiffView).
			Tap(func() {
				t.ExpectToast(Equals("Changed diff context size to 2"))
			}).
			NavigateToLine(Contains("-ten")).
			PressPrimaryAction()

		t.Views().Secondary().Content(
			Contains("-one").Contains("-ten"),
		)
		t.Views().Main().MarkedLines(
			Contains("-one"),
			Contains("-ten"),
		)
	},
})
