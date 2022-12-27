package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var (
	postMergeFileContent = "post merge file content"
	postMergeFilename    = "post-merge-file"
)

var AmendMerge = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Amends a staged file to a merge commit.",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			NewBranch("development-branch").
			CreateFileAndAdd("initial-file", "initial file content").
			Commit("initial commit").
			NewBranch("feature-branch"). // it's also checked out automatically
			CreateFileAndAdd("new-feature-file", "new content").
			Commit("new feature commit").
			Checkout("development-branch").
			Merge("feature-branch").
			CreateFileAndAdd(postMergeFilename, postMergeFileContent)
	},
	Run: func(shell *Shell, input *Input, keys config.KeybindingConfig) {
		input.Model().CommitCount(3)

		mergeCommitMessage := "Merge branch 'feature-branch' into development-branch"
		input.Model().HeadCommitMessage(Contains(mergeCommitMessage))

		input.Views().Commits().
			Focus().
			Press(keys.Commits.AmendToCommit)

		input.ExpectConfirmation().
			Title(Equals("Amend Commit")).
			Content(Contains("Are you sure you want to amend this commit with your staged files?")).
			Confirm()

		// assuring we haven't added a brand new commit
		input.Model().CommitCount(3)
		input.Model().HeadCommitMessage(Contains(mergeCommitMessage))

		// assuring the post-merge file shows up in the merge commit.
		input.Views().Main().
			Content(Contains(postMergeFilename)).
			Content(Contains("++" + postMergeFileContent))
	},
})
