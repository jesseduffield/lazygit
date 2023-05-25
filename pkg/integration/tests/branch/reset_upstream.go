package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ResetUpstream = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Reset the upstream of a branch",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("one")
		shell.CloneIntoRemote("origin")
		shell.SetBranchUpstream("master", "origin/master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().
			Focus().
			Press(keys.Universal.NextScreenMode). // we need to enlargen the window to see the upstream
			Lines(
				Contains("master").Contains("origin master").IsSelected(),
			).
			Press(keys.Branches.SetUpstream).
			Tap(func() {
				t.ExpectPopup().Menu().
					Title(Equals("Set/Unset upstream")).
					Select(Contains("Unset upstream of selected branch")).
					Confirm()
			}).
			Lines(
				Contains("master").DoesNotContain("origin master").IsSelected(),
			)
	},
})
