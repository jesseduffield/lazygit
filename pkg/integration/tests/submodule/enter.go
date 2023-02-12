package submodule

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Enter = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Enter a submodule, add a commit, and then stage the change in the parent repo",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.UserConfig.CustomCommands = []config.CustomCommand{
			{
				Key:     "e",
				Context: "files",
				Command: "git commit --allow-empty -m \"empty commit\"",
			},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.RunCommand("git clone --bare . ../other_repo")
		shell.RunCommand("git submodule add ../other_repo my_submodule")
		shell.GitAddAll()
		shell.Commit("add submodule")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Submodules().Focus().
			Lines(
				Contains("my_submodule").IsSelected(),
			).
			// enter the submodule
			PressEnter()

		t.Views().Files().IsFocused().
			Press("e").
			Tap(func() {
				t.Views().Commits().Content(Contains("empty commit"))
			}).
			// return to the parent repo
			PressEscape()

		t.Views().Submodules().IsFocused()

		// we see the new commit in the submodule is ready to be staged in the parent repo
		t.Views().Main().Content(Contains("> empty commit"))

		t.Views().Files().Focus().
			Lines(
				MatchesRegexp(` M.*my_submodule \(submodule\)`).IsSelected(),
			).
			Tap(func() {
				// main view also shows the new commit when we're looking at the submodule within the files view
				t.Views().Main().Content(Contains("> empty commit"))
			}).
			PressPrimaryAction().
			Press(keys.Files.CommitChanges).
			Tap(func() {
				t.ExpectPopup().CommitMessagePanel().Type("submodule change").Confirm()
			}).
			IsEmpty()

		t.Views().Submodules().Focus()

		// we no longer report a new commit because we've committed it in the parent repo
		t.Views().Main().Content(DoesNotContain("> empty commit"))
	},
})
