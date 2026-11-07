package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
	"github.com/jesseduffield/lazygit/pkg/integration/tests/shared"
)

var LockErrorPartwayThroughRebase = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Drops a commit while another git process takes index.lock partway through the rebase",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(4)
		// The rebase makes a commit when it picks commit-03; the lock then
		// stops it when it tries to pick commit-04.
		shared.HoldIndexLockAfterNextCommit(shell)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit-04").IsSelected(),
				Contains("commit-03"),
				Contains("commit-02"),
				Contains("commit-01"),
			).
			NavigateToLine(Contains("commit-02")).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Drop commit")).
					Content(Contains("Are you sure you want to drop the selected commit(s)?")).
					Confirm()
			}).
			/* EXPECTED:
			IsFocused().
			Lines(
				Contains("commit-04"),
				Contains("commit-03"),
				Contains("commit-01"),
			)
			ACTUAL: */
			Tap(func() {
				t.ExpectPopup().Alert().
					Title(Equals("Error")).
					Content(Contains("It seems that there is already a rebase-merge directory")).
					Confirm()
			}).
			Lines(
				Contains("─── Pending rebase todos"),
				Contains("pick").Contains("commit-04"),
				Contains("─── Commits"),
				Contains("commit-03"),
				Contains("commit-01"),
			)
	},
})
