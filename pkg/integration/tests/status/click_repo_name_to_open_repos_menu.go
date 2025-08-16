package status

import (
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ClickRepoNameToOpenReposMenu = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Click on the repo name in the status side panel to open the recent repositories menu",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		// Create another repo that we can switch to
		shell.CloneIntoRemote("other_repo")
		shell.CreateFileAndAdd("other_repo/README.md", "# other_repo")
		shell.Commit("initial commit")

		// Create a third repo
		shell.CloneIntoRemote("third_repo")
		shell.CreateFileAndAdd("third_repo/README.md", "# third_repo")
		shell.Commit("initial commit")

		// Add these repos to the recent repos list
		otherRepoPath, err := filepath.Abs("other_repo")
		if err != nil {
			panic(err)
		}
		thirdRepoPath, err := filepath.Abs("third_repo")
		if err != nil {
			panic(err)
		}

		// Switch to the other repo and back to add it to the recent repos list
		shell.RunCommand([]string{"cd", otherRepoPath})
		shell.RunCommand([]string{"cd", "-"})
		shell.RunCommand([]string{"cd", thirdRepoPath})
		shell.RunCommand([]string{"cd", "-"})
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().Click(1, 0)
		t.ExpectPopup().Menu().Title(Equals("Recent repositories")).
			Lines(
				Contains("other_repo").IsSelected(),
				Contains("third_repo"),
				Contains("Cancel"),
			)
	},
})
