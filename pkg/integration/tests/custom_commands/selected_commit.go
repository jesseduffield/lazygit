package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectedCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Use the {{ .SelectedCommit }} template variable in different contexts",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommits(3)
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:     "X",
				Context: "global",
				Command: "printf '%s' '{{ .SelectedCommit.Name }}' > file.txt",
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Select different commits in each of the commit views
		t.Views().Commits().Focus().
			NavigateToLine(Contains("commit 01"))
		t.Views().ReflogCommits().Focus().
			NavigateToLine(Contains("commit 02"))
		t.Views().Branches().Focus().
			Lines(Contains("master").IsSelected()).
			PressEnter()
		t.Views().SubCommits().IsFocused().
			NavigateToLine(Contains("commit 03"))

		// SubCommits
		t.GlobalPress("X")
		t.FileSystem().FileContent("file.txt", Equals("commit 03"))

		t.Views().SubCommits().PressEnter()
		t.GlobalPress("X")
		t.FileSystem().FileContent("file.txt", Equals("commit 03"))

		// ReflogCommits
		t.Views().ReflogCommits().Focus()
		t.GlobalPress("X")
		t.FileSystem().FileContent("file.txt", Equals("commit: commit 02"))

		t.Views().ReflogCommits().PressEnter()
		t.GlobalPress("X")
		t.FileSystem().FileContent("file.txt", Equals("commit: commit 02"))

		// LocalCommits
		t.Views().Commits().Focus()
		t.GlobalPress("X")
		t.FileSystem().FileContent("file.txt", Equals("commit 01"))

		t.Views().Commits().PressEnter()
		t.GlobalPress("X")
		t.FileSystem().FileContent("file.txt", Equals("commit 01"))

		// None of these
		t.Views().Files().Focus()
		t.GlobalPress("X")
		t.FileSystem().FileContent("file.txt", Equals("commit 01"))
	},
})
