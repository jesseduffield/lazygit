package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ChangeScreenModeInFocusedDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Enlarge the focused diff, which leaves the pane showing the other side of the file behind",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("file", "first line\nsecond line")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			PressPrimaryAction().
			Title(Equals("Unstaged changes")).
			Content(Contains("+second line").DoesNotContain("+first line")).
			PressTab()

		t.Views().Secondary().
			IsFocused().
			Title(Equals("Staged changes")).
			Content(Contains("+first line").DoesNotContain("+second line")).
			Press(keys.Universal.NextScreenMode).
			Tap(func() {
				// Half screen: the side panels are gone, the two diff panes are not.
				t.Views().AppStatus().IsInvisible()
				t.Views().Main().IsVisible()
			}).
			Press(keys.Universal.NextScreenMode).
			Tap(func() {
				// Full screen: the focused pane has the window to itself.
				t.Views().AppStatus().IsInvisible()
				t.Views().Main().IsInvisible()
			})
	},
})
