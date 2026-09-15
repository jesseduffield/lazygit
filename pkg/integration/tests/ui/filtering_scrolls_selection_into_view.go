package ui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var FilteringScrollsSelectionIntoView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Entering and leaving filtering mode scrolls the selected commit into view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		for i := range 40 {
			file := "otherFile"
			if i%2 == 0 {
				file = "filterFile"
			}
			shell.UpdateFileAndAdd(file, fmt.Sprintf("content %02d", i))
			shell.Commit(fmt.Sprintf("commit %02d", i))
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Press(keys.Universal.GotoBottom).
			SelectedLine(Contains("commit 00")).
			OriginYAtLeast(1).
			Press(keys.Universal.FilteringMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Filtering")).
			Select(Contains("Enter path to filter by")).
			Confirm()
		t.ExpectPopup().Prompt().
			Title(Equals("Enter path:")).
			Type("filterFile").
			Confirm()

		// The filtered list has nothing to do with the one that was showing, so
		// its scroll position doesn't either: we start at the top again
		t.Views().Commits().
			IsFocused().
			SelectedLine(Contains("commit 38")).
			SelectedLineIdx(0).
			OriginY(0).
			Press(keys.Universal.GotoBottom).
			SelectedLine(Contains("commit 00")).
			PressEscape()

		// Leaving filtering mode keeps the commit selected, at its position in
		// the full list, which needs scrolling to again
		t.Views().Commits().
			IsFocused().
			SelectedLine(Contains("commit 00")).
			SelectedLineIsVisible()
	},
})
