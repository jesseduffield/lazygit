package submodule

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ResolveConflict = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Resolve a submodule conflict (both sides moved the gitlink) by picking one side's commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowFileTree = false
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.CloneIntoSubmodule("my_submodule_name", "my_submodule_path")
		shell.GitAddAll()
		shell.Commit("add submodule")

		sub := "my_submodule_path"

		// Two diverging commits in the submodule, so the gitlink can't be
		// fast-forwarded and the merge genuinely conflicts.
		shell.RunCommand([]string{"git", "-C", sub, "checkout", "-b", "left"})
		shell.RunCommand([]string{"git", "-C", sub, "commit", "--allow-empty", "-m", "left"})
		shell.RunCommand([]string{"git", "-C", sub, "checkout", "-b", "right", "HEAD~1"})
		shell.RunCommand([]string{"git", "-C", sub, "commit", "--allow-empty", "-m", "right"})

		// "ours" points the submodule at left, "theirs" at right.
		shell.RunCommand([]string{"git", "checkout", "-b", "ours"})
		shell.RunCommand([]string{"git", "-C", sub, "checkout", "left"})
		shell.RunCommand([]string{"git", "add", sub})
		shell.Commit("ours")

		shell.RunCommand([]string{"git", "checkout", "-b", "theirs", "HEAD~1"})
		shell.RunCommand([]string{"git", "-C", sub, "checkout", "right"})
		shell.RunCommand([]string{"git", "add", sub})
		shell.Commit("theirs")

		shell.RunCommand([]string{"git", "checkout", "ours"})
		shell.RunCommand([]string{"git", "-C", sub, "checkout", "left"})
		shell.RunCommandExpectError([]string{"git", "merge", "theirs"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			Lines(
				Contains("UU my_submodule_path (submodule)").IsSelected(),
			).
			Tap(func() {
				// The main view explains the conflict and shows each side's
				// commits as separate "current" and "incoming" logs.
				t.Views().Main().Content(
					Contains("Conflict: the submodule").
						Contains("Current changes:").Contains("left").
						Contains("Incoming changes:").Contains("right"),
				)
			}).
			// Enter opens the resolution menu instead of entering the submodule.
			// The two candidate commits are shown with their summaries.
			Press(keys.Universal.GoInto).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Merge conflicts")).
					Select(Contains("Take current commit").Contains("left")).
					Select(Contains("Take incoming commit").Contains("right")).
					Cancel()
			}).
			// Space opens the same menu; take the incoming commit to resolve.
			PressPrimaryAction().
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Merge conflicts")).
					Select(Contains("Take incoming commit")).
					Confirm()
			}).
			Lines(
				Contains("M  my_submodule_path (submodule)").IsSelected(),
			)
	},
})
