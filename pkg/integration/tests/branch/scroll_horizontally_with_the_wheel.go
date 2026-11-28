package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ScrollHorizontallyWithTheWheel = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The horizontal mouse wheel reveals the end of a branch name that is too long for the branches view, and stops there even though the divergence from the base branch is right-aligned",
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
		t.Views().Files().
			IsFocused()

		t.Views().Branches().
			OriginX(0).
			Lines(
				Contains("* a-branch-whose-name").DoesNotContain("view-end").Contains("↓1"),
				Contains("master"),
			).
			// Twice gui.scrollHeight columns at a time.
			ScrollWheelRight().
			OriginX(4).
			ScrollWheelRight().
			ScrollWheelRight().
			ScrollWheelRight().
			ScrollWheelRight().
			ScrollWheelRight().
			OriginX(24).
			// The whole branch name and the divergence fit now.
			ScrollWheelRight().
			OriginX(26).
			Lines(
				Contains("* a-branch-whose-name-is-too-long-for-the-branches-view-end ↓1"),
				Contains("master"),
			).
			ScrollWheelRight().
			OriginX(26)
	},
})
