package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// Even though this involves submodules, it's a worktree test since
// it's really exercising lazygit's ability to correctly do pathfinding
// in a complex use case.
var DoubleNestedLinkedSubmodule = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Open lazygit in a link to a repo's double nested submodules",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowFileTree = false
	},
	SetupRepo: func(shell *Shell) {
		// we're going to have a directory structure like this:
		// project
		//  - repo/outerSubmodule/innerSubmodule/a/b/c
		//  - link (symlink to repo/outerSubmodule/innerSubmodule/a/b/c)
		//
		shell.CreateFileAndAdd("rootFile", "rootStuff")
		shell.Commit("initial repo commit")

		shell.Chdir("..")
		shell.CreateDir("innerSubmodule")
		shell.Chdir("innerSubmodule")
		shell.Init()
		shell.CreateFileAndAdd("a/b/c/blah", "blah\n")
		shell.Commit("initial inner commit")

		shell.Chdir("..")
		shell.CreateDir("outerSubmodule")
		shell.Chdir("outerSubmodule")
		shell.Init()
		shell.CreateFileAndAdd("foo", "foo")
		shell.Commit("initial outer commit")
		// the git config (-c) parameter below is required
		// to let git create a file-protocol/path submodule
		shell.RunCommand([]string{"git", "-c", "protocol.file.allow=always", "submodule", "add", "../innerSubmodule"})
		shell.Commit("add dependency as innerSubmodule")

		shell.Chdir("../repo")
		shell.RunCommand([]string{"git", "-c", "protocol.file.allow=always", "submodule", "add", "../outerSubmodule"})
		shell.Commit("add dependency as outerSubmodule")
		shell.Chdir("outerSubmodule")
		shell.RunCommand([]string{"git", "-c", "protocol.file.allow=always", "submodule", "update", "--init", "--recursive"})

		shell.Chdir("innerSubmodule")
		shell.UpdateFile("a/b/c/blah", "updated content\n")

		shell.Chdir("../../..")
		shell.RunCommand([]string{"ln", "-s", "repo/outerSubmodule/innerSubmodule/a/b/c", "link"})

		shell.Chdir("link")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Lines(
				Contains("HEAD detached"),
				Contains("master"),
			)

		t.Views().Commits().
			Lines(
				Contains("initial inner commit"),
			)

		t.Views().Files().
			IsFocused().
			Lines(
				Contains(" M a/b/c/blah"), // shows as modified
			).
			PressPrimaryAction().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Title(Equals("Commit summary")).
			Type("Update blah").
			Confirm()

		t.Views().Files().
			IsEmpty()

		t.Views().Commits().
			Lines(
				Contains("Update blah"),
				Contains("initial inner commit"),
			)
	},
})
