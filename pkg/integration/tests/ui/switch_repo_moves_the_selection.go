package ui

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SwitchRepoMovesTheSelection = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The selection follows the focus of the repo being switched to, rather than the one being left",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		otherRepo, _ := filepath.Abs("../other")
		config.GetAppState().RecentRepos = []string{otherRepo}
	},
	SetupRepo: func(shell *Shell) {
		shell.CloneNonBare("other")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		switchToRepo := func(repo string) {
			t.GlobalPress(keys.Universal.OpenRecentRepos)
			t.ExpectPopup().Menu().Title(Equals("Recent repositories")).
				Lines(
					Contains(repo).IsSelected(),
					Contains("Cancel"),
				).Confirm()
			t.Views().Status().Content(Contains(repo + " → master"))
		}

		t.Views().Branches().
			Focus().
			SelectionIsActive()

		// The other repo has its own focus, which is the files panel it starts in
		switchToRepo("other")
		t.Views().Files().IsFocused()
		t.Views().Branches().SelectionIsHidden()

		// And coming back, this repo still has the focus we left it with
		switchToRepo("repo")
		t.Views().Branches().
			IsFocused().
			SelectionIsActive()
		t.Views().Files().SelectionIsHidden()
	},
})
