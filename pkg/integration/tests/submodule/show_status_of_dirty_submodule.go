package submodule

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ShowStatusOfDirtySubmodule = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "The main view shows the status of a submodule that has changes that aren't committed, above its diff",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Gui.ShowFileTree = false
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("file", "content")
		shell.Commit("first commit")
		shell.CloneIntoSubmodule("dirty_name", "dirty_path")
		shell.CloneIntoSubmodule("new_commit_name", "new_commit_path")
		shell.Commit("add submodules")

		shell.UpdateFile("dirty_path/file", "changed content")
		shell.CreateFile("dirty_path/untracked_file", "content")

		shell.Chdir("new_commit_path")
		shell.EmptyCommit("new commit")
		shell.Chdir("..")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Equals(" M dirty_path (submodule)").IsSelected(),
				Equals(" M new_commit_path (submodule)"),
			)

		t.Views().Main().
			Content(
				Contains("On branch master").
					Contains("modified:   file").
					Contains("untracked_file").
					Contains("── Diff ─").
					Contains("Submodule dirty_path contains modified content"),
			)

		t.Views().Files().
			NavigateToLine(Contains("new_commit_path"))

		t.Views().Main().
			Content(
				Contains("Submodule new_commit_path").
					Contains("> new commit").
					DoesNotContain("── Diff ─"),
			)

		t.Views().Submodules().
			Focus().
			Lines(
				Contains("dirty_name").IsSelected(),
				Contains("new_commit_name"),
			)

		t.Views().Main().
			Content(
				Contains("Name: dirty_name").
					Contains("On branch master").
					Contains("modified:   file").
					Contains("── Diff ─").
					Contains("Submodule dirty_path contains modified content"),
			)

		t.Views().Submodules().
			NavigateToLine(Contains("new_commit_name"))

		t.Views().Main().
			Content(
				Contains("Name: new_commit_name").
					Contains("> new commit").
					DoesNotContain("── Diff ─"),
			)
	},
})
