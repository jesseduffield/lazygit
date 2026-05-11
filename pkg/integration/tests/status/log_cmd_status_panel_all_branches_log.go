package status

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var LogCmdStatusPanelAllBranchesLog = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Cycle between two different log commands in the Status view when it has status panel AllBranchesLog",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.AllBranchesLogCmds = []string{
			`echo "view1"`,
			`echo "view2"`,
			`echo "view3"`,
		}
		config.GetUserConfig().Gui.StatusPanelView = "allBranchesLog"
	},
	SetupRepo: func(shell *Shell) {},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Status().
			Focus()
		t.Views().Main().Content(Contains("view1"))

		// We head to the branches view and return
		t.Views().Branches().
			Focus()
		t.Views().Status().
			Focus()

		t.Views().Main().Content(Contains("view1"))

		t.Views().Status().
			Press(keys.Status.AllBranchesLogGraph)
		t.Views().Main().Content(Contains("view2"))

		t.Views().Status().
			Press(keys.Status.AllBranchesLogGraph)
		t.Views().Main().Content(Contains("view3"))

		t.Views().Status().
			Press(keys.Status.AllBranchesLogGraph)
		t.Views().Main().Content(Contains("view1"))

		t.Views().Status().
			Press(keys.Status.AllBranchesLogGraphReverse)
		t.Views().Main().Content(Contains("view3"))

		t.Views().Status().
			Press(keys.Status.AllBranchesLogGraphReverse)
		t.Views().Main().Content(Contains("view2"))

		t.Views().Status().
			Press(keys.Status.AllBranchesLogGraphReverse)
		t.Views().Main().Content(Contains("view1"))
	},
})
