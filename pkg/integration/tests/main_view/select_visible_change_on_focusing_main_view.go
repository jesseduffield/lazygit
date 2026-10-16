package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectVisibleChangeOnFocusingMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Focusing the main view selects a change that is already on screen, leaving the diff where it is",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
		// Enough context around each change to scroll into a stretch of the diff that
		// holds none.
		cfg.GetUserConfig().Git.DiffContextSize = 20
	},
	SetupRepo: func(shell *Shell) {
		lines := make([]string, 60)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%02d", i+1)
		}
		shell.CreateFileAndAdd("file1", strings.Join(lines, "\n")+"\n")
		shell.Commit("one")

		lines[4] = strings.ToUpper(lines[4])
		lines[54] = strings.ToUpper(lines[54])
		shell.UpdateFile("file1", strings.Join(lines, "\n")+"\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Universal.ScrollDownMain).
			Press(keys.Universal.ScrollDownMain)

		t.Views().Main().
			OriginY(4).
			Tap(func() {
				t.Views().Files().Press(keys.Universal.FocusMainView)
			}).
			IsFocused().
			// The first change is on screen, so it is the one to point at — and the view
			// hasn't moved to point at it.
			SelectedLines(
				Contains("-line05"),
			).
			OriginY(4).
			PressEscape()

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.ScrollDownMain).
			Press(keys.Universal.ScrollDownMain).
			Press(keys.Universal.ScrollDownMain).
			Press(keys.Universal.ScrollDownMain)

		t.Views().Main().
			OriginY(12).
			Tap(func() {
				t.Views().Files().Press(keys.Universal.FocusMainView)
			}).
			IsFocused().
			// Now the whole screen is context: the changes are above and below it. The
			// selection goes to the middle of what is on screen rather than to a change
			// the user would have to be scrolled to.
			SelectedLines(
				Contains(" line19"),
			).
			SelectedLineIdx(24).
			OriginY(12)
	},
})
