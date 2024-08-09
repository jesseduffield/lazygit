package custom_commands

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ShowOutputInPanel = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Run a command and show the output in a panel",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("my change")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().CustomCommands = []config.CustomCommand{
			{
				Key:        "X",
				Context:    "commits",
				Command:    "printf '%s' '{{ .SelectedLocalCommit.Name }}'",
				ShowOutput: true,
			},
			{
				Key:         "Y",
				Context:     "commits",
				Command:     "printf '%s' '{{ .SelectedLocalCommit.Name }}'",
				ShowOutput:  true,
				OutputTitle: "Subject of commit {{ .SelectedLocalCommit.Hash }}",
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

		t.ExpectPopup().Alert().
			// Uses cmd string as title if no outputTitle is provided
			Title(Equals("printf '%s' 'my change'")).
			Content(Equals("my change")).
			Confirm()

		t.Views().Commits().
			Press("Y")

		hash := t.Git().GetCommitHash("HEAD")
		t.ExpectPopup().Alert().
			// Uses provided outputTitle with template fields resolved
			Title(Equals(fmt.Sprintf("Subject of commit %s", hash))).
			Content(Equals("my change"))
	},
})
