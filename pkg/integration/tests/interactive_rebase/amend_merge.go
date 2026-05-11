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
	ExtraCmdArgs: []string{},
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
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		mergeCommitMessage := "Merge branch 'feature-branch' into development-branch"

		t.Views().Commits().
			Lines(
				Contains(mergeCommitMessage),
				Contains("new feature commit"),
				Contains("initial commit"),
			)

		t.Views().Commits().
			Focus().
			Press(keys.Commits.AmendToCommit)

		t.ExpectPopup().Confirmation().
			Title(Equals("Amend commit")).
			Content(Contains("Are you sure you want to amend this commit with your staged files?")).
			Confirm()

		// assuring we haven't added a brand new commit
		t.Views().Commits().
			Lines(
				Contains(mergeCommitMessage),
				Contains("new feature commit"),
				Contains("initial commit"),
			)

		// assuring the post-merge file shows up in the merge commit.
		t.Views().Main().
			Content(Contains(postMergeFilename)).
			Content(Contains("++" + postMergeFileContent))
	},
})
