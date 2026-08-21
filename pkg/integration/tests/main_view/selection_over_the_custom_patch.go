package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectionOverTheCustomPatch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The pane showing a custom patch keeps its selection across a refresh, its content being a diff like any other",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\n")
		shell.Commit("one")
		shell.UpdateFileAndAdd("file1", "one\nTWO\nthree\n")
		shell.Commit("two")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("two")).
			PressEnter()

		t.Views().CommitFiles().
			IsFocused().
			PressPrimaryAction().
			Press(keys.Universal.FocusMainView)

		// The patch built from the commit is shown in the other pane, and it is a diff,
		// so it has a selection of its own — one that a refresh doesn't take away.
		t.Views().Main().
			IsFocused().
			Press(keys.Universal.TogglePanel)

		t.Views().Secondary().
			IsFocused().
			Title(Equals("Custom patch")).
			SelectionIsActive().
			Tap(func() {
				t.GlobalPress(keys.Universal.Refresh)
			}).
			SelectionIsActive()
	},
})
