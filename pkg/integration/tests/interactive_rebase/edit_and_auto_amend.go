package interactive_rebase

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var EditAndAutoAmend = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Edit a commit, make a change and stage it, then continue the rebase to auto-amend the commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateNCommits(3)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit 03"),
				Contains("commit 02"),
				Contains("commit 01"),
			).
			NavigateToLine(Contains("commit 02")).
			Press(keys.Universal.Edit).
			Lines(
				Contains("commit 03"),
				MatchesRegexp("YOU ARE HERE.*commit 02").IsSelected(),
				Contains("commit 01"),
			)

		t.Shell().CreateFile("fixup-file", "fixup content")
		t.Views().Files().
			Focus().
			Press(keys.Files.RefreshFiles).
			Lines(
				Contains("??").Contains("fixup-file").IsSelected(),
			).
			PressPrimaryAction()

		t.Common().ContinueRebase()

		t.Views().Commits().
			Focus().
			Lines(
				Contains("commit 03"),
				Contains("commit 02").IsSelected(),
				Contains("commit 01"),
			)

		t.Views().Main().
			Content(Contains("fixup content"))
	},
})
