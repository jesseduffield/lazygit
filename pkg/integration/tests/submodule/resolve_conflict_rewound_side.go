package submodule

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ResolveConflictRewoundSide = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "When a side of a submodule conflict added no commits of its own (it was rewound), the main view shows the commit it points at instead of an empty log",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowFileTree = false
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.CloneIntoSubmodule("sub_name", "sub_path")
		shell.GitAddAll()
		shell.Commit("add submodule")

		sub := "sub_path"

		// Mark the submodule's initial commit, then advance it; the merge base
		// will point the submodule here.
		shell.RunCommand([]string{"git", "-C", sub, "branch", "initial"})
		shell.RunCommand([]string{"git", "-C", sub, "commit", "--allow-empty", "-m", "s1"})
		shell.RunCommand([]string{"git", "add", sub})
		shell.Commit("base at s1")

		// "ours" rewinds the submodule to its initial commit (so it has no
		// commits of its own relative to "theirs").
		shell.RunCommand([]string{"git", "checkout", "-b", "ours"})
		shell.RunCommand([]string{"git", "-C", sub, "checkout", "initial"})
		shell.RunCommand([]string{"git", "add", sub})
		shell.Commit("ours rewinds submodule")

		// "theirs" advances the submodule with a further commit.
		shell.RunCommand([]string{"git", "checkout", "-b", "theirs", "HEAD~1"})
		shell.RunCommand([]string{"git", "-C", sub, "checkout", "master"})
		shell.RunCommand([]string{"git", "-C", sub, "commit", "--allow-empty", "-m", "s2"})
		shell.RunCommand([]string{"git", "add", sub})
		shell.Commit("theirs advances submodule")

		shell.RunCommand([]string{"git", "checkout", "ours"})
		shell.RunCommand([]string{"git", "-C", sub, "checkout", "initial"})
		shell.RunCommandExpectError([]string{"git", "merge", "theirs"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			Lines(
				Contains("UU sub_path (submodule)").IsSelected(),
			).
			Tap(func() {
				// "ours" has no commits of its own, so its section falls back to
				// the commit it points at; "theirs" lists the commits it added.
				t.Views().Main().Content(
					Contains("Current changes:").Contains("first commit").
						Contains("Incoming changes:").Contains("s1").Contains("s2"),
				)
			})
	},
})
