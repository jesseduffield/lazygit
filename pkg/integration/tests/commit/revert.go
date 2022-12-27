package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Revert = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Reverts a commit",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("myfile", "myfile content")
		shell.GitAddAll()
		shell.Commit("first commit")
	},
	Run: func(shell *Shell, input *Input, keys config.KeybindingConfig) {
		input.Model().CommitCount(1)

		input.Views().Commits().
			Focus().
			Lines(
				Contains("first commit"),
			).
			Press(keys.Commits.RevertCommit).
			Tap(func() {
				input.ExpectConfirmation().
					Title(Equals("Revert commit")).
					Content(MatchesRegexp(`Are you sure you want to revert \w+?`)).
					Confirm()
			}).
			Lines(
				Contains("Revert \"first commit\"").IsSelected(),
				Contains("first commit"),
			)

		input.Views().Main().Content(Contains("-myfile content"))
		input.FileSystem().PathNotPresent("myfile")
	},
})
