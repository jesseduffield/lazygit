package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ScrollHorizontallyStopsAtTheEnd = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Scrolling the branches view sideways with H and L stops when the whole name of the longest branch fits, even though the divergence from the base branch is right-aligned",
	ExtraCmdArgs: []string{},
	Skip:         false,
	Width:        120,
	Height:       30,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowDivergenceFromBaseBranch = "arrowAndNumber"
	},
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("one").
			NewBranch("a-branch-whose-name-is-too-long-for-the-branches-view-end").
			Checkout("master").
			EmptyCommit("two").
			Checkout("a-branch-whose-name-is-too-long-for-the-branches-view-end")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Lines(
				Contains("* a-branch-whose-name").DoesNotContain("view-end").Contains("↓1"),
				Contains("master"),
			).
			OriginX(0).
			// A third of the view's width at a time.
			Press(keys.Universal.ScrollRight).
			OriginX(12).
			Press(keys.Universal.ScrollRight).
			OriginX(24).
			// The whole branch name and the divergence fit after 26 columns.
			Press(keys.Universal.ScrollRight).
			OriginX(26).
			Lines(
				Contains("* a-branch-whose-name-is-too-long-for-the-branches-view-end ↓1"),
				Contains("master"),
			).
			Press(keys.Universal.ScrollRight).
			OriginX(26)
	},
})
