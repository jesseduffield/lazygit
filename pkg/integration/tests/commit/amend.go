package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AmendHead = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Amending head.",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(5). // these will appear as commit 05, 04, 04, down to 01
			CreateFileAndAdd("password.txt", "hunter2")
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		input.SwitchToCommitsWindow()
		assert.CurrentViewName("commits")

		input.PressKeys(keys.Commits.AmendToCommit)
		assert.InConfirm()
		assert.MatchCurrentViewContent(Contains("Are you sure you want to amend this commit"))
		input.Confirm()

		assert.CommitCount(5)
		assert.WorkingTreeFileCount(0)
		assert.MatchMainViewContent(Contains("password.txt"))
		assert.MatchMainViewContent(Contains("hunter2"))
	},
})
