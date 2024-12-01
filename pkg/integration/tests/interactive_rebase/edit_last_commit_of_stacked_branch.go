package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var EditLastCommitOfStackedBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Edit and amend the last commit of a branch in a stack of branches, and ensure that it doesn't break the stack",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.38.0"),
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.MainBranches = []string{"master"}
		config.GetAppState().GitLogShowGraph = "never"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(1).
			NewBranch("branch1").
			CreateNCommitsStartingAt(2, 2).
			NewBranch("branch2").
			CreateNCommitsStartingAt(2, 4)

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI commit 05").IsSelected(),
				Contains("CI commit 04"),
				Contains("CI * commit 03"),
				Contains("CI commit 02"),
				Contains("CI commit 01"),
			).
			NavigateToLine(Contains("commit 03")).
			Press(keys.Universal.Edit).
			Lines(
				Contains("pick").Contains("CI commit 05"),
				Contains("pick").Contains("CI commit 04"),
				Contains("update-ref").Contains("branch1"),
				Contains("<-- YOU ARE HERE --- * commit 03").IsSelected(),
				Contains("CI commit 02"),
				Contains("CI commit 01"),
			)

		t.Shell().CreateFile("fixup-file", "fixup content")
		t.Views().Files().
			Focus().
			Press(keys.Files.RefreshFiles).
			Lines(
				Contains("??").Contains("fixup-file").IsSelected(),
			).
			PressPrimaryAction().
			Press(keys.Files.AmendLastCommit)
		t.ExpectPopup().Confirmation().
			Title(Equals("Amend last commit")).
			Content(Contains("Are you sure you want to amend last commit?")).
			Confirm()

		t.Common().ContinueRebase()

		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI commit 05"),
				Contains("CI commit 04"),
				Contains("CI * commit 03"),
				Contains("CI commit 02"),
				Contains("CI commit 01"),
			)
	},
})
