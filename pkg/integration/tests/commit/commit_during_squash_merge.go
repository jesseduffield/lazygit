package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitDuringSquashMerge = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Commit the changes of a squash merge that was left uncommitted, using the squash merge message",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.CommitPrefix = []config.CommitPrefixConfig{{Pattern: "^(\\w+)-branch$", Replace: "[$1] "}}
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("original-branch").
			EmptyCommit("one").
			NewBranch("feature-branch").
			CreateFileAndAdd("feature-file", "content").
			Commit("feature").
			Checkout("original-branch")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("original-branch").IsSelected(),
				Contains("feature-branch"),
			).
			SelectNextItem().
			Press(keys.Branches.MergeIntoCurrentBranch)

		t.ExpectPopup().Menu().
			Title(Equals("Merge")).
			Select(Contains("Squash merge and leave uncommitted")).
			Confirm()

		t.Views().Files().
			Focus().
			Lines(
				Equals("A  feature-file"),
			).
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Title(Equals("Commit summary")).
			InitialText(Equals("Squash merge feature-branch into original-branch")).
			Confirm()

		t.Views().Status().Content(Equals("repo → original-branch"))

		t.Views().Commits().
			Lines(
				Contains("Squash merge feature-branch into original-branch"),
				Contains("one"),
			)
	},
})
