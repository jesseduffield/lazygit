package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DoNotShowBranchMarkerForHeadCommit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Verify that no branch heads are shown for the branch head if there is a tag with the same name as the branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	GitVersion:   AtLeast("2.38.0"),
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.NewBranch("branch1")
		shell.EmptyCommit("two")
		shell.EmptyCommit("three")
		shell.CreateLightweightTag("branch1", "master")

		shell.SetConfig("rebase.updateRefs", "true")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Check that the local commits view does show a branch marker for the head commit
		t.Views().Commits().
			Lines(
				Contains("CI three"),
				Contains("CI two"),
				Contains("CI branch1 one"),
			)
	},
})
