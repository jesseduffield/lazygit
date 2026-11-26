package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ExplainNestedRepositories = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The main view explains what a removed submodule, a nested repository and a worktree inside the repo are, and shows their status",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowFileTree = false
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.CloneIntoSubmodule("sub", "modules/sub")
		shell.Commit("add submodule")
		shell.RunCommand([]string{"git", "rm", "--cached", "modules/sub"})
		shell.RunCommand([]string{"git", "rm", ".gitmodules"})
		shell.Commit("remove submodule")
		shell.CreateFile("modules/sub/new-file", "content")

		shell.RunCommand([]string{"git", "init", "nested"})

		shell.AddWorktree("master", "worktree", "newbranch")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals("?? modules/sub (removed submodule)").IsSelected(),
				Equals("?? nested (nested repository)"),
				Equals("?? worktree (worktree)"),
			)

		t.Views().Main().
			Content(
				Contains("'modules/sub' is left over from a submodule").
					Contains("press d to discard the directory").
					Contains("── Status of 'modules/sub' ─").
					Contains("On branch master").
					Contains("new-file"),
			)

		t.Views().Files().
			NavigateToLine(Contains("nested"))

		t.Views().Main().
			Content(
				Contains("'nested' is a git repository of its own").
					Contains("ignore it with i").
					Contains("of its history.").
					Contains("── Status of 'nested' ─").
					Contains("No commits yet"),
			)

		t.Views().Files().
			NavigateToLine(Contains("worktree"))

		t.Views().Main().
			Content(
				Contains("'worktree' is a linked worktree").
					Contains("press d.").
					Contains("── Status of 'worktree' ─").
					Contains("On branch newbranch"),
			)
	},
})
