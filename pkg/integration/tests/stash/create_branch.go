package stash

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CreateBranch = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create a branch from a stash entry",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.CreateFile("myfile", "content")
		shell.GitAddAll()
		shell.Stash("stash one")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().IsEmpty()

		t.Views().Stash().
			Focus().
			Lines(
				Contains("stash one").IsSelected(),
			).
			Press(keys.Universal.New).
			Tap(func() {
				t.ExpectPopup().Prompt().
					Title(Contains("New Branch Name (Branch is off of 'stash@{0}: On master: stash one'")).
					Type("new_branch").
					Confirm()
			})

		t.Views().Files().IsEmpty()

		t.Views().Branches().
			IsFocused().
			Lines(
				Contains("new_branch").IsSelected(),
				Contains("master"),
			).
			PressEnter()

		t.Views().SubCommits().
			Lines(
				Contains("On master: stash one").IsSelected(),
				MatchesRegexp(`index on master:.*initial commit`),
				Contains("initial commit"),
			)

		t.Views().Main().Content(Contains("myfile | 1 +"))
	},
})
