package ui

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var KeybindingSuggestionsWhenSwitchingRepos = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Show correct keybinding suggestions after switching between repos",
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

		t.Views().Files().Focus()
		t.Views().Options().Content(
			Equals("Commit: c | Stash: s | Reset: D | Keybindings: ? | Cancel: <esc>"))

		switchToRepo("other")
		switchToRepo("repo")

		t.Views().Options().Content(
			Equals("Commit: c | Stash: s | Reset: D | Keybindings: ? | Cancel: <esc>"))
	},
})
