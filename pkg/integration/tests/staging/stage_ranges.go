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
			SelectedLine(Contains("+three")).
			Press(keys.Main.ToggleDragSelect).
			SelectNextItem().
			SelectedLine(Contains("+four")).
			SelectNextItem().
			SelectedLine(Contains("+five")).
			// stage the three lines we've just selected
			PressPrimaryAction().
			Content(Contains(" five\n+six")).
			Tap(func() {
				t.Views().StagingSecondary().
					Content(Contains("+three\n+four\n+five"))
			}).
			Press(keys.Universal.TogglePanel)

		t.Views().StagingSecondary().
			IsFocused().
			SelectedLine(Contains("+three")).
			Press(keys.Main.ToggleDragSelect).
			SelectNextItem().
			SelectedLine(Contains("+four")).
			SelectNextItem().
			SelectedLine(Contains("+five")).
			// unstage the three selected lines
			PressPrimaryAction().
			// nothing left in our staging secondary panel
			IsEmpty().
			Tap(func() {
				t.Views().Staging().
					Content(Contains("+three\n+four\n+five\n+six"))
			})

		t.Views().Staging().
			IsFocused().
			// coincidentally we land at '+four' here. Maybe we should instead land
			// at '+three'? given it's at the start of the hunk?
			SelectedLine(Contains("+four")).
			Press(keys.Main.ToggleDragSelect).
			SelectNextItem().
			SelectedLine(Contains("+five")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.Actions().ConfirmDiscardLines()
			}).
			Content(Contains("+three\n+six"))
	},
})
