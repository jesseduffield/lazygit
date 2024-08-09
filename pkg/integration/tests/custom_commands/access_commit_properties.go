package custom_commands

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AccessCommitProperties = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Run a command that accesses properties of a commit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("my change")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:     "X",
				Context: "commits",
				Command: "printf '%s\n%s\n%s' '{{ .SelectedLocalCommit.Name }}' '{{ .SelectedLocalCommit.Hash }}' '{{ .SelectedLocalCommit.Sha }}' > file.txt",
			},
		}
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("my change").IsSelected(),
			).
			Press("X")

		hash := t.Git().GetCommitHash("HEAD")
		t.FileSystem().FileContent("file.txt", Equals(fmt.Sprintf("my change\n%s\n%s", hash, hash)))
	},
})
