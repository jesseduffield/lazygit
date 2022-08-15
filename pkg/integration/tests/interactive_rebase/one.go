package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var One = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Begins an interactive rebase, then fixups, drops, and squashes some commits",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(5) // these will appears at commit 05, 04, 04, down to 01
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		input.SwitchToCommitsWindow()
		assert.CurrentViewName("commits")

		input.NavigateToListItemContainingText("commit 02")
		input.PressKeys(keys.Universal.Edit)
		assert.MatchSelectedLine(Contains("YOU ARE HERE"))

		input.PreviousItem()
		input.PressKeys(keys.Commits.MarkCommitAsFixup)
		assert.MatchSelectedLine(Contains("fixup"))

		input.PreviousItem()
		input.PressKeys(keys.Universal.Remove)
		assert.MatchSelectedLine(Contains("drop"))

		input.PreviousItem()
		input.PressKeys(keys.Commits.SquashDown)
		assert.MatchSelectedLine(Contains("squash"))

		input.ContinueRebase()

		assert.CommitCount(2)
	},
})
