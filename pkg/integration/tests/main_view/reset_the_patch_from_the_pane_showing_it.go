package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ResetThePatchFromThePaneShowingIt = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Reset a custom patch while the focus is in the pane previewing it",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
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
			SelectedLines(
				Contains("-two"),
				Contains("+TWO"),
			).
			PressPrimaryAction().
			Press(keys.Universal.TogglePanel)

		t.Views().Information().Content(Contains("Building patch"))
		t.Views().Secondary().
			IsFocused().
			Content(Contains("-two"))

		t.Common().SelectPatchOption(Contains("Reset patch"))

		// The pane goes with the patch it was previewing, and the focus follows into the
		// one showing the diff the patch was built from.
		t.Views().Information().Content(DoesNotContain("Building patch"))
		t.Views().Secondary().IsInvisible()
		t.Views().Main().
			IsFocused().
			Content(Contains("-two"))
	},
})
