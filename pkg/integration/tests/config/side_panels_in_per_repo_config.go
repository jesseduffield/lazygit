package config

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SidePanelsInPerRepoConfig = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "A per-repo config can set the side panel layout, and switching repos re-applies each repo's own layout",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		otherRepo, _ := filepath.Abs("../other")
		cfg.GetAppState().RecentRepos = []string{otherRepo}
	},
	SetupRepo: func(shell *Shell) {
		shell.CloneNonBare("other")
		// The other repo swaps the branches and commits panels.
		shell.CreateFile("../other/.git/lazygit.yml", `
gui:
  sidePanels:
    - status
    - [files, worktrees, submodules]
    - [commits, reflog]
    - [branches, remotes, tags]
    - stash`)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// This repo uses the default layout, so the third panel is branches.
		t.GlobalPress(keys.Universal.JumpToBlock[2])
		t.Views().Branches().IsFocused()

		// Switch to the other repo, whose per-repo config swaps branches and commits.
		t.GlobalPress(keys.Universal.OpenRecentRepos)
		t.ExpectPopup().Menu().Title(Equals("Recent repositories")).
			Lines(
				Contains("other").IsSelected(),
				Contains("Cancel"),
			).Confirm()
		t.Views().Status().Content(Contains("other → master"))

		// Now the third panel is commits.
		t.GlobalPress(keys.Universal.JumpToBlock[2])
		t.Views().Commits().IsFocused()

		// Switch back to the first repo; its default layout is intact even though
		// its contexts were built before we visited the other repo.
		t.GlobalPress(keys.Universal.JumpToBlock[1])
		t.Views().Files().IsFocused()
		t.GlobalPress(keys.Universal.OpenRecentRepos)
		t.ExpectPopup().Menu().Title(Equals("Recent repositories")).Confirm()

		t.GlobalPress(keys.Universal.JumpToBlock[2])
		t.Views().Branches().IsFocused()
	},
})
