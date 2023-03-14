package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageRanges = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage and unstage various ranges of a file in the staging panel",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nthree\nfour\nfive\nsix\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("file1").IsSelected(),
			).
			PressEnter()

		t.Views().Staging().
			IsFocused().
			SelectedLines(
				Contains("+three"),
			).
			Press(keys.Main.ToggleDragSelect).
			NavigateToLine(Contains("+five")).
			SelectedLines(
				Contains("+three"),
				Contains("+four"),
				Contains("+five"),
			).
			// stage the three lines we've just selected
			PressPrimaryAction().
			SelectedLines(
				Contains("+six"),
			).
			ContainsLines(
				Contains(" five"),
				Contains("+six"),
			).
			Tap(func() {
				t.Views().StagingSecondary().
					ContainsLines(
						Contains("+three"),
						Contains("+four"),
						Contains("+five"),
					)
			}).
			Press(keys.Universal.TogglePanel)

		t.Views().StagingSecondary().
			IsFocused().
			SelectedLines(
				Contains("+three"),
			).
			Press(keys.Main.ToggleDragSelect).
			NavigateToLine(Contains("+five")).
			SelectedLines(
				Contains("+three"),
				Contains("+four"),
				Contains("+five"),
			).
			// unstage the three selected lines
			PressPrimaryAction().
			// nothing left in our staging secondary panel
			IsEmpty().
			Tap(func() {
				t.Views().Staging().
					ContainsLines(
						Contains("+three"),
						Contains("+four"),
						Contains("+five"),
						Contains("+six"),
					)
			})

		t.Views().Staging().
			IsFocused().
			// coincidentally we land at '+four' here. Maybe we should instead land
			// at '+three'? given it's at the start of the hunk?
			SelectedLines(
				Contains("+four"),
			).
			Press(keys.Main.ToggleDragSelect).
			SelectNextItem().
			SelectedLines(
				Contains("+four"),
				Contains("+five"),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.Common().ConfirmDiscardLines()
			}).
			ContainsLines(
				Contains("+three"),
				Contains("+six"),
			)
	},
})
