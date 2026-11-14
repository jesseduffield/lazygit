package config

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RebindConfirmSearchAndPrompt = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Confirm a search and a text prompt using a rebound confirm keybinding",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		// The search and prompt confirm keys used to be hard-coded to <enter>,
		// so they couldn't be rebound like the other confirm keybindings.
		cfg.GetUserConfig().Keybinding.Universal.ConfirmSearch = []string{"<c-y>"}
		cfg.GetUserConfig().Keybinding.Universal.ConfirmPrompt = []string{"<c-y>"}
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch")
		shell.EmptyCommit("one")
		shell.EmptyCommit("two")
		shell.EmptyCommit("three")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Confirming a search with the rebound key works.
		t.Views().Commits().
			Focus().
			Lines(
				Contains("three").IsSelected(),
				Contains("two"),
				Contains("one"),
			).
			Press(keys.Universal.StartSearch).
			Tap(func() {
				t.ExpectSearch().
					Type("two")

				t.GlobalPress(keys.Universal.ConfirmSearch)

				t.Views().Search().IsVisible().Content(Contains("matches for 'two' (1 of 1)"))
			}).
			Lines(
				Contains("three"),
				Contains("two").IsSelected(),
				Contains("one"),
			)

		// Confirming a text input prompt with the rebound key works.
		t.Views().Branches().
			Focus().
			Press(keys.Universal.New)

		t.ExpectPopup().Prompt().
			Title(Contains("New branch name")).
			Type("new-branch")

		t.GlobalPress(keys.Universal.ConfirmPrompt)

		t.Git().CurrentBranchName("new-branch")
	},
})
