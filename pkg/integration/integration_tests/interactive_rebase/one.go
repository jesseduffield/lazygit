package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/integration/types"
)

var One = types.NewTest(types.NewTestArgs{
	Description:  "Begins an interactive rebase, then fixups, drops, and squashes some commits",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell types.Shell) {
		shell.
			CreateNCommits(5) // these will appears at commit 05, 04, 04, down to 01
	},
	Run: func(shell types.Shell, input types.Input, assert types.Assert, keys config.KeybindingConfig) {
		input.SwitchToCommitsWindow()
		assert.CurrentViewName("commits")

		input.NavigateToListItemContainingText("commit 02")
		input.PressKeys(keys.Universal.Edit)
		assert.SelectedLineContains("YOU ARE HERE")

		input.PreviousItem()
		input.PressKeys(keys.Commits.MarkCommitAsFixup)
		assert.SelectedLineContains("fixup")

		input.PreviousItem()
		input.PressKeys(keys.Universal.Remove)
		assert.SelectedLineContains("drop")

		input.PreviousItem()
		input.PressKeys(keys.Commits.SquashDown)
		assert.SelectedLineContains("squash")

		input.ContinueRebase()

		assert.CommitCount(2)
	},
})
