package status

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var LogCmd = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cycle between two different log commands in the Status view",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.Git.AllBranchesLogCmd = `echo "view1"`
		config.UserConfig.Git.AllBranchesLogCmds = []string{`echo "view2"`}
	},
	SetupRepo: func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().
			Focus().
			Press(keys.Status.AllBranchesLogGraph)
		t.Views().Main().Content(Contains("view1"))

		t.Views().Status().
			Focus().
			Press(keys.Status.AllBranchesLogGraph)
		t.Views().Main().Content(Contains("view2").DoesNotContain("view1"))

		t.Views().Status().
			Focus().
			Press(keys.Status.AllBranchesLogGraph)
		t.Views().Main().Content(Contains("view1").DoesNotContain("view2"))
	},
})
