package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AmendHeadDuringRebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Amend the HEAD commit during a rebase.",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(5)
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		input.SwitchToCommitsWindow()
		assert.CurrentViewName("commits")

		input.NavigateToListItemContainingText("commit 02")
		input.PressKeys(keys.Universal.Edit)
		assert.MatchSelectedLine(Contains("YOU ARE HERE"))

		shell.CreateFileAndAdd("password.txt", "hunter2")

		input.SwitchToFilesWindow()
		input.PressKeys(keys.Files.RefreshFiles)
		assert.WorkingTreeFileCount(1)
		input.SwitchToCommitsWindow()
		assert.CurrentViewName("commits")

		input.PressKeys(keys.Commits.AmendToCommit)
		input.ProceedWhenAsked(Contains("Are you sure you want to amend"))

		input.ContinueRebase()

		assert.CommitCount(5)
		assert.WorkingTreeFileCount(0)
		assert.MatchMainViewContent(Contains("password.txt"))
		assert.MatchMainViewContent(Contains("hunter2"))
	},
})
