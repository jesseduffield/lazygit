package patch_building

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ResetPatchBuiltFromMainView = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Reset a custom patch that was built straight from the commits panel's main view, without ever having entered the commit files panel",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.UseHunkModeInStagingView = false
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branch-a")
		shell.CreateFileAndAdd("file1", "one\ntwo\nthree\nfour\nfive\n")
		shell.Commit("first commit")

		shell.NewBranch("branch-b")
		shell.UpdateFileAndAdd("file1", "one\ntwo\nTHREE\nfour\nfive\n")
		shell.Commit("update")

		shell.Checkout("branch-a")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Press(keys.Universal.NextItem).
			PressEnter()

		// Build a patch straight from the whole-commit diff, never entering the commit
		// files panel — so its context is never set up with a ref for this commit.
		t.Views().SubCommits().
			IsFocused().
			Press(keys.Universal.FocusMainView)

		t.Views().Main().
			IsFocused().
			Press(keys.Main.ToggleSelectHunk).
			PressPrimaryAction()

		t.Views().Information().Content(Contains("Building patch"))

		// Resetting the patch from the menu refreshes the (never-initialised) commit files
		// context; this used to crash on its nil ref.
		t.Views().Main().Press(keys.Universal.CreatePatchOptionsMenu)
		t.ExpectPopup().Menu().
			Title(Equals("Patch options")).
			Select(Contains("Reset patch")).
			Confirm()

		t.Views().Information().Content(DoesNotContain("Building patch"))
	},
})
