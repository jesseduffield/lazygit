package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AmendFixupCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Amends a staged file to a fixup commit, and checks that other unrelated fixup commits are not auto-squashed.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(1).
			CreateFileAndAdd("first-fixup-file", "").Commit("fixup! commit 01").
			CreateNCommitsStartingAt(2, 2).
			CreateFileAndAdd("unrelated-fixup-file", "fixup 03").Commit("fixup! commit 03").
			CreateFileAndAdd("fixup-file", "fixup 01")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("fixup! commit 03"),
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("fixup! commit 01"),
				Contains("commit 01"),
			).
			NavigateToLine(Contains("fixup! commit 01")).
			Press(keys.Commits.AmendToCommit).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Amend commit")).
					Content(Contains("Are you sure you want to amend this commit with your staged files?")).
					Confirm()
			}).
			Lines(
				Contains("fixup! commit 03"),
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("fixup! commit 01").IsSelected(),
				Contains("commit 01"),
			)

		t.Views().Main().
			Content(Contains("fixup 01"))
	},
})
