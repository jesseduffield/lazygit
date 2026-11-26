package main_view

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var JumpToAFileOnlyOverADiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The menu of the diff's files is offered only while the main view is showing a diff",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shared.CreateMergeConflictFiles(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// A commit shows its diff in the main view, so the menu is offered over it.
		t.Views().Commits().
			Focus().
			NavigateToLine(Contains("original")).
			Press(keys.Universal.OptionMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			Tap(func() {
				t.Views().Menu().Content(Contains("Jump to file in diff"))
			}).
			Cancel()

		// The pane binds the same key itself, so the global one isn't offered on top of
		// the pane's while the pane has the focus.
		t.Views().Commits().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Press(keys.Universal.OptionMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			Tap(func() {
				t.Views().Menu().
					Content(Contains("Jump to file")).
					Content(DoesNotContain("Jump to file in diff"))
			}).
			Cancel()

		t.Views().Main().
			IsFocused().
			Press(keys.Universal.Return)

		// Working through a conflicted file gives the main section over to the merge
		// conflicts view, which is no diff to jump around in. The pane behind it goes
		// on holding the diff it last rendered, so it is the view on screen that
		// decides.
		t.Views().Files().
			Focus().
			NavigateToLine(Contains("UU file1")).
			PressEnter()

		t.Views().MergeConflicts().
			IsFocused().
			Press(keys.Universal.OptionMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Keybindings")).
			Tap(func() {
				t.Views().Menu().Content(DoesNotContain("Jump to file"))
			}).
			Cancel()
	},
})
