package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Checkout = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Checkout a commit as a detached head, or checkout an existing branch at a commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")
		shell.NewBranch("branch1")
		shell.NewBranch("branch2")
		shell.EmptyCommit("three")
		shell.EmptyCommit("four")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("four").IsSelected(),
				Contains("three"),
				Contains("two"),
				Contains("one"),
			).
			PressPrimaryAction()

		t.ExpectPopup().Menu().
			Title(Contains("Checkout branch or commit")).
			Lines(
				Contains("Checkout branch").IsSelected(),
				MatchesRegexp("Checkout commit [a-f0-9]+ as detached head"),
				Contains("Cancel"),
			).
			Tooltip(Contains("Disabled: No branches found at selected commit.")).
			Select(MatchesRegexp("Checkout commit [a-f0-9]+ as detached head")).
			Confirm()
		t.Views().Branches().Lines(
			Contains("* (HEAD detached at"),
			Contains("branch2"),
			Contains("branch1"),
			Contains("master"),
		)

		t.Views().Commits().
			NavigateToLine(Contains("two")).
			PressPrimaryAction()

		t.ExpectPopup().Menu().
			Title(Contains("Checkout branch or commit")).
			Lines(
				Contains("Checkout branch 'branch1'").IsSelected(),
				Contains("Checkout branch 'master'"),
				MatchesRegexp("Checkout commit [a-f0-9]+ as detached head"),
				Contains("Cancel"),
			).
			Select(Contains("Checkout branch 'master'")).
			Confirm()
		t.Views().Branches().Lines(
			Contains("master"),
			Contains("branch2"),
			Contains("branch1"),
		)
	},
})
