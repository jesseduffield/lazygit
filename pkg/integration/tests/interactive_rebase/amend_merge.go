package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AmendMerge = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Amends a staged file to a merge commit.",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			NewBranch("development-branch").
			CreateFileAndAdd("initial-file", "content").
			Commit("initial commit").
			NewBranch("feature-branch"). // it's also checked out automatically
			CreateFileAndAdd("new-feature-file", "new content").
			Commit("new feature commit").
			CheckoutBranch("development-branch").
			Merge("feature-branch").
			CreateFileAndAdd("post-merge-file", "content")
	},
	Run: func(shell *Shell, input *Input, assert *Assert, keys config.KeybindingConfig) {
		assert.CommitCount(3)

		input.SwitchToCommitsWindow()
		assert.CurrentViewName("commits")

		input.PressKeys(keys.Commits.AmendToCommit)
		input.PressKeys(keys.Universal.Return)

		assert.MatchHeadCommitMessage(Contains("Merge"))
		assert.CommitCount(3)
	},
})
