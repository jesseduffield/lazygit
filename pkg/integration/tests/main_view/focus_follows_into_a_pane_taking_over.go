package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FocusFollowsIntoAPaneTakingOver = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The pane the focus follows into as it takes the section over gets its selection from the top of the diff it is given",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Gui.UseHunkModeInDiffView = false
	},
	SetupRepo: func(shell *Shell) {
		lines := make([]string, 40)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%02d", i+1)
		}
		shell.CreateFileAndAdd("file1", strings.Join(lines, "\n")+"\n")
		shell.Commit("one")

		// Everything staged and nothing unstaged, so the staged side has the section to
		// itself, with more changes in it than fit.
		for i := range lines {
			lines[i] = strings.ToUpper(lines[i])
		}
		shell.UpdateFileAndAdd("file1", strings.Join(lines, "\n")+"\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("M  file1").IsSelected(),
			)

		// Read a way into the diff before focusing it, so that the selection starts out
		// somewhere other than the first change.
		t.Views().Secondary().
			IsVisible().
			ScrollWheelDown().
			ScrollWheelDown().
			ScrollWheelDown().
			ScrollWheelDown().
			OriginY(8)

		t.Views().Files().Press(keys.Universal.FocusMainView)
		t.Views().Secondary().
			IsFocused().
			SelectedLines(
				Contains("-line04"),
			)

		// The index is reset outside lazygit, so the staged side empties and the
		// unstaged side has everything: the pane the focus is in goes away and the other
		// one takes the section over with a diff that is new to it.
		t.Shell().RunCommand([]string{"git", "reset"})
		t.GlobalPress(keys.Universal.Refresh)

		t.Views().Secondary().IsInvisible()
		t.Views().Main().
			IsFocused().
			Title(Equals("Unstaged changes")).
			OriginY(0).
			SelectedLines(
				Contains("-line01"),
			)
	},
})
