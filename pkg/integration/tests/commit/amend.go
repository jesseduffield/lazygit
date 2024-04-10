package commit

import (
	"github.com/lobes/lazytask/pkg/config"
	. "github.com/lobes/lazytask/pkg/integration/components"
)

var Amend = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Amends the last commit from the files panel",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("myfile", "myfile content\n")
		shell.Commit("first commit")
		shell.UpdateFileAndAdd("myfile", "myfile content\nmore content\n")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Lines(
				Contains("first commit"),
			)

		t.Views().Files().
			Focus().
			Press(keys.Commits.AmendToCommit)

		t.ExpectPopup().Confirmation().Title(
			Equals("Amend last commit")).
			Content(Contains("Are you sure you want to amend last commit?")).
			Confirm()

		t.Views().Commits().
			Focus().
			Lines(
				Contains("first commit"),
			)

		t.Views().Main().Content(Contains("+myfile content").Contains("+more content"))
	},
})
