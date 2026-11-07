package submodule

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var IgnoreSuspiciousPath = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Ignore a submodule whose path in .gitmodules leads outside the repository",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// Two unrelated repos next to ours, each with an uncommitted change
		for _, otherRepo := range []string{"other_repo", "linked_repo"} {
			shell.Chdir("..")
			shell.CreateDir(otherRepo)
			shell.Chdir(otherRepo)
			shell.Init()
			shell.CreateFileAndAdd("file", "original")
			shell.Commit("initial commit")
			shell.UpdateFile("file", "modified")
			shell.Chdir("../repo")
		}

		shell.EmptyCommit("first commit")
		shell.RunCommand([]string{"ln", "-s", "../linked_repo", "link"})
		shell.GitAdd("link")
		shell.CreateFileAndAdd(".gitmodules", `[submodule "other"]
	path = ../other_repo
	url = ../other_repo
[submodule "linked"]
	path = link
	url = ../linked_repo
`)
		shell.Commit("add submodules")
		shell.CreateFile("untracked", "content")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Submodules().
			IsEmpty()

		t.Views().Files().Focus().
			Lines(
				Contains("untracked"),
			).
			Press(keys.Files.ViewResetOptions).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("")).
					Select(Contains("Nuke working tree")).
					Confirm()

				t.ExpectPopup().Confirmation().
					Title(Equals("Nuke working tree")).
					Content(AnyString()).
					Confirm()
			}).
			IsEmpty()

		t.FileSystem().
			FileContent("../other_repo/file", Equals("modified")).
			FileContent("../linked_repo/file", Equals("modified"))
	},
})
