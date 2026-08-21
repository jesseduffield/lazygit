package main_view

import (
	"fmt"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectInADiffReadInPart = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Focusing the main view over a diff longer than the part of it that has been read still shows a selection",
	ExtraCmdArgs: []string{},
	Skip:         false,
	// A short terminal, so that the file below is longer than the initial read of its diff.
	Width:       100,
	Height:      20,
	SetupConfig: func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		lines := make([]string, 600)
		for i := range lines {
			lines[i] = fmt.Sprintf("line%03d", i+1)
		}
		shell.CreateFileAndAdd("file1", strings.Join(lines, "\n")+"\n")
		shell.Commit("one big commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			SelectionIsActive()
	},
})
