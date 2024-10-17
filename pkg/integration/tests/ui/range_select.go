package ui

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Here's the state machine we need to verify:
// (no range, press 'v') -> sticky range
// (no range, press arrow) -> no range
// (no range, press shift+arrow) -> nonsticky range
// (sticky range, press 'v') -> no range
// (sticky range, press 'escape') -> no range
// (sticky range, press arrow) -> sticky range
// (sticky range, press `<`/`>` or `,`/`.`) -> sticky range
// (sticky range, press shift+arrow) -> nonsticky range
// (nonsticky range, press 'v') -> no range
// (nonsticky range, press 'escape') -> no range
// (nonsticky range, press arrow) -> no range
// (nonsticky range, press shift+arrow) -> nonsticky range

// Importantly, if you press 'v' when in a nonsticky range, it clears the range,
// so no matter which mode you're in, 'v' will cancel the range.
// And, if you press shift+up/down when in a sticky range, it switches to a non-
// sticky range, meaning if you then press up/down without shift, it clears
// the range.

var RangeSelect = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify range select works as expected in list views and in patch explorer views",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// We're testing the commits view as our representative list context,
		// as well as the staging view, and we're using the exact same code to test
		// both to ensure they have the exact same behaviour (they are currently implemented
		// separately)
		// In both views we're going to have 10 lines starting from 'line 1' going down to
		// 'line 10'.
		fileContent := ""
		total := 10
		for i := 1; i <= total; i++ {
			remaining := total - i + 1
			// Commits are displayed in reverse order so to we need to create them in reverse to have them appear as 'line 1', 'line 2' etc.
			shell.EmptyCommit(fmt.Sprintf("line %d", remaining))
			fileContent = fmt.Sprintf("%sline %d\n", fileContent, i)
		}
		shell.CreateFile("file1", fileContent)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		assertRangeSelectBehaviour := func(v *ViewDriver, otherView *ViewDriver, lineIdxOfFirstItem int) {
			v.
				SelectedLines(
					Contains("line 1"),
				).
				// (no range, press 'v') -> sticky range
				Press(keys.Universal.ToggleRangeSelect).
				SelectedLines(
					Contains("line 1"),
				).
				// (sticky range, press arrow) -> sticky range
				SelectNextItem().
				SelectedLines(
					Contains("line 1"),
					Contains("line 2"),
				).
				// (sticky range, press 'v') -> no range
				Press(keys.Universal.ToggleRangeSelect).
				SelectedLines(
					Contains("line 2"),
				).
				// (no range, press arrow) -> no range
				SelectPreviousItem().
				SelectedLines(
					Contains("line 1"),
				).
				// (no range, press shift+arrow) -> nonsticky range
				Press(keys.Universal.RangeSelectDown).
				SelectedLines(
					Contains("line 1"),
					Contains("line 2"),
				).
				// (nonsticky range, press shift+arrow) -> nonsticky range
				Press(keys.Universal.RangeSelectDown).
				SelectedLines(
					Contains("line 1"),
					Contains("line 2"),
					Contains("line 3"),
				).
				Press(keys.Universal.RangeSelectUp).
				SelectedLines(
					Contains("line 1"),
					Contains("line 2"),
				).
				// (nonsticky range, press arrow) -> no range
				SelectNextItem().
				SelectedLines(
					Contains("line 3"),
				).
				Press(keys.Universal.ToggleRangeSelect).
				SelectedLines(
					Contains("line 3"),
				).
				SelectNextItem().
				SelectedLines(
					Contains("line 3"),
					Contains("line 4"),
				).
				// (sticky range, press shift+arrow) -> nonsticky range
				Press(keys.Universal.RangeSelectDown).
				SelectedLines(
					Contains("line 3"),
					Contains("line 4"),
					Contains("line 5"),
				).
				SelectNextItem().
				SelectedLines(
					Contains("line 6"),
				).
				Press(keys.Universal.RangeSelectDown).
				SelectedLines(
					Contains("line 6"),
					Contains("line 7"),
				).
				// (nonsticky range, press 'v') -> no range
				Press(keys.Universal.ToggleRangeSelect).
				SelectedLines(
					Contains("line 7"),
				).
				Press(keys.Universal.RangeSelectDown).
				SelectedLines(
					Contains("line 7"),
					Contains("line 8"),
				).
				// (nonsticky range, press 'escape') -> no range
				PressEscape().
				SelectedLines(
					Contains("line 8"),
				).
				// (sticky range, press '>') -> sticky range
				Press(keys.Universal.ToggleRangeSelect).
				Press(keys.Universal.GotoBottom).
				SelectedLines(
					Contains("line 8"),
					Contains("line 9"),
					Contains("line 10"),
				).
				// (sticky range, press 'escape') -> no range
				PressEscape().
				SelectedLines(
					Contains("line 10"),
				)

			// Click in view, press shift+arrow -> nonsticky range
			otherView.Focus()
			v.Click(1, lineIdxOfFirstItem).
				SelectedLines(
					Contains("line 1"),
				).
				Press(keys.Universal.RangeSelectDown).
				SelectedLines(
					Contains("line 1"),
					Contains("line 2"),
				)
		}

		assertRangeSelectBehaviour(t.Views().Commits().Focus(), t.Views().Branches(), 0)

		t.Views().Files().
			Focus().
			SelectedLine(
				Contains("file1"),
			).
			PressEnter()

		assertRangeSelectBehaviour(t.Views().Staging().IsFocused(), t.Views().Files(), 6)
	},
})
