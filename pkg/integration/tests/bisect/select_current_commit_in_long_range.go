package bisect

import (
	"fmt"

	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var SelectCurrentCommitInLongRange = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Select the current bisect commit when it is further down than the commits that we load at first",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		for i := range 700 {
			shell.EmptyCommit(fmt.Sprintf("commit %03d", i+1))
		}
		shell.StartBisect("HEAD", "HEAD~699")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.Log.ShowGraph = "never"
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			SelectedLine(Contains("<-- current"))
	},
})
