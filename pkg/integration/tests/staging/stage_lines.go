package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageLines = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Stage and unstage various lines of a file in the staging panel",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file1", "one\ntwo\n")
		shell.Commit("one")

		shell.UpdateFile("file1", "one\ntwo\nthree\nfour\n")
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
			SelectedLines(Contains("+three")).
			// stage 'three'
			PressPrimaryAction().
			// 'three' moves over to the staging secondary panel
			Content(DoesNotContain("+three")).
			Tap(func() {
				t.Views().StagingSecondary().
					ContainsLines(
						Contains("+three"),
					)
			}).
			SelectedLines(Contains("+four")).
			// stage 'four'
			PressPrimaryAction().
			// nothing left in our staging panel
			IsEmpty()

		// because we've staged everything we get moved to the staging secondary panel
		// do the same thing as above, moving the lines back to the staging panel
		t.Views().StagingSecondary().
			IsFocused().
			ContainsLines(
				Contains("+three"),
				Contains("+four"),
			).
			SelectedLines(Contains("+three")).
			PressPrimaryAction().
			Content(DoesNotContain("+three")).
			Tap(func() {
				t.Views().Staging().
					ContainsLines(
						Contains("+three"),
					)
			}).
			SelectedLines(Contains("+four")).
			// pressing 'remove' has the same effect as pressing space when in the staging secondary panel
			Press(keys.Universal.Remove).
			IsEmpty()

		// stage one line and then manually toggle to the staging secondary panel
		t.Views().Staging().
			IsFocused().
			ContainsLines(
				Contains("+three"),
				Contains("+four"),
			).
			SelectedLines(Contains("+three")).
			PressPrimaryAction().
			Content(DoesNotContain("+three")).
			Tap(func() {
				t.Views().StagingSecondary().
					Content(Contains("+three"))
			}).
			Press(keys.Universal.TogglePanel)

		// manually toggle back to the staging panel
		t.Views().StagingSecondary().
			IsFocused().
			Press(keys.Universal.TogglePanel)

		t.Views().Staging().
			SelectedLines(Contains("+four")).
			// discard the line
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Discard change")).
					Content(Contains("Are you sure you want to discard this change")).
					Confirm()
			}).
			IsEmpty()

		t.Views().StagingSecondary().
			IsFocused().
			ContainsLines(
				Contains("+three"),
			).
			// return to file
			PressEscape()

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("M  file1").IsSelected(),
			).
			PressEnter()

		// because we only have a staged change we'll land in the staging secondary panel
		t.Views().StagingSecondary().
			IsFocused()
	},
})
