package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DropCommitInCopiedBranchWithUpdateRef = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Drops a commit in a branch that is a copy of another branch, and verify that the other branch is left alone",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.38.0"),
	SetupConfig: func(config *config.AppConfig) {
		config.AppState.GitLogShowGraph = "never"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			NewBranch("branch1").
			CreateNCommits(3).
			NewBranch("branch2")

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("CI * commit 03").IsSelected(),
				Contains("CI commit 02"),
				Contains("CI commit 01"),
			).
			NavigateToLine(Contains("commit 02")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Drop commit")).
					Content(Equals("Are you sure you want to drop the selected commit(s)?")).
					Confirm()
			}).
			Lines(
				Contains("CI commit 03"), // no start on this commit because branch1 is no longer pointing to it
				Contains("CI commit 01"),
			)

		t.Views().Branches().
			Focus().
			NavigateToLine(Contains("branch1")).
			PressPrimaryAction()

		t.Views().Commits().Lines(
			Contains("CI commit 03"),
			Contains("CI commit 02"),
			Contains("CI commit 01"),
		)
	},
})
