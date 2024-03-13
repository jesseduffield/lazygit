package submodule

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Reset = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Enter a submodule, create a commit and stage some changes, then reset the submodule from back in the parent repo. This test captures functionality around getting a dirty submodule out of your files panel.",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.UserConfig.CustomCommands = []config.CustomCommand{
			{
				Key:     "e",
				Context: "files",
				Command: "git commit --allow-empty -m \"empty commit\" && echo \"my_file content\" > my_file",
			},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.CloneIntoSubmodule("my_submodule_name", "my_submodule_path")
		shell.GitAddAll()
		shell.Commit("add submodule")

		shell.CreateFile("other_file", "")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		assertInParentRepo := func() {
			t.Views().Status().Content(Contains("repo"))
		}
		assertInSubmodule := func() {
			if t.Git().Version().IsAtLeast(2, 22, 0) {
				t.Views().Status().Content(Contains("my_submodule_path(my_submodule_name)"))
			} else {
				t.Views().Status().Content(Contains("my_submodule_path"))
			}
		}

		assertInParentRepo()

		t.Views().Submodules().Focus().
			Lines(
				Contains("my_submodule_name").IsSelected(),
			).
			// enter the submodule
			PressEnter()

		assertInSubmodule()

		t.Views().Files().IsFocused().
			Press("e").
			Tap(func() {
				t.Views().Commits().Content(Contains("empty commit"))
				t.Views().Files().Content(Contains("my_file"))
			}).
			Lines(
				Contains("my_file").IsSelected(),
			).
			// stage my_file
			PressPrimaryAction().
			// return to the parent repo
			PressEscape()

		assertInParentRepo()

		t.Views().Submodules().IsFocused()

		t.Views().Main().Content(Contains("Submodule my_submodule_path contains modified content"))

		t.Views().Files().Focus().
			Lines(
				MatchesRegexp(` M.*my_submodule_path \(submodule\)`),
				Contains("other_file").IsSelected(),
			).
			// Verify we can't use range select on submodules
			Press(keys.Universal.ToggleRangeSelect).
			SelectPreviousItem().
			Lines(
				MatchesRegexp(` M.*my_submodule_path \(submodule\)`).IsSelected(),
				Contains("other_file").IsSelected(),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectToast(Contains("Disabled: Range select not supported for submodules"))
			}).
			Press(keys.Universal.ToggleRangeSelect).
			Lines(
				MatchesRegexp(` M.*my_submodule_path \(submodule\)`).IsSelected(),
				Contains("other_file"),
			).
			Press(keys.Universal.Remove).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("my_submodule_path")).
					Select(Contains("Stash uncommitted submodule changes and update")).
					Confirm()
			}).
			Lines(
				Contains("other_file").IsSelected(),
			)

		t.Views().Submodules().Focus().
			PressEnter()

		assertInSubmodule()

		// submodule has been hard reset to the commit the parent repo specifies
		t.Views().Branches().Lines(
			Contains("HEAD detached").IsSelected(),
			Contains("master"),
		)

		// empty commit is gone
		t.Views().Commits().Lines(
			Contains("first commit").IsSelected(),
		)

		// the staged change has been stashed
		t.Views().Files().IsEmpty()

		t.Views().Stash().Focus().
			Lines(
				Contains("WIP on master").IsSelected(),
			)

		t.Views().Main().Content(Contains("my_file content"))
	},
})
