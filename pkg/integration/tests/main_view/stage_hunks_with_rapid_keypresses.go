package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// The second space is pressed before the diff the first one changed has re-rendered.
// The selection only moves to the next hunk once that render arrives, so until then it
// still covers lines that are no longer in the diff. The second press has to wait for
// the render rather than act on those lines.
var StageHunksWithRapidKeypresses = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage two hunks from the focused main view with two space presses in rapid succession",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInStagingView = true
	},
	SetupRepo: func(shell *Shell) {
		// Seven context lines between the two change blocks, so that git makes them two
		// hunks.
		shell.CreateFileAndAdd("file1", "1\n2\na\nb\nc\nd\ne\nf\ng\n3\n4\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "1b\n2b\na\nb\nc\nd\ne\nf\ng\n3b\n4b\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			PressRapidly(keys.Universal.Select, keys.Universal.Select)

		// Both presses were acted on, so the whole file is staged.
		t.Views().Files().Lines(
			Contains("M  file1"),
		)
		t.Views().Secondary().
			Title(Equals("Staged changes")).
			ContainsLines(
				Contains("+1b"),
				Contains("+2b"),
			).
			ContainsLines(
				Contains("+3b"),
				Contains("+4b"),
			)
	},
})
