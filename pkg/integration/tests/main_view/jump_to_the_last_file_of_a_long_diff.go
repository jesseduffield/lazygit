package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var JumpToTheLastFileOfALongDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Jumping to the last file of a long diff from the commits panel doesn't scroll past the end of the diff",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        100,
	Height:       20,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ScrollPastBottom = false
	},
	SetupRepo: func(shell *Shell) {
		// Longer than the 5000 lines that a render reads at most before it waits for
		// the user to scroll further.
		lines := make([]string, 6000)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%04d", i+1)
		}
		shell.CreateFileAndAdd("long.txt", strings.Join(lines, "\n")+"\n")
		shell.CreateFileAndAdd("short.txt", "one\n")
		shell.Commit("first commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			SelectedLine(Contains("first commit")).
			Press(keys.Universal.JumpToFile)

		t.ExpectPopup().Menu().
			Title(Equals("Jump to file in diff")).
			Select(Equals("short.txt")).
			Confirm()

		t.Views().Commits().
			IsFocused()

		// The view shows the last screenful of the diff: the end of long.txt's diff,
		// and all of short.txt's below it.
		t.Views().Main().
			TopVisibleLine(Equals("+line5994"))
	},
})
