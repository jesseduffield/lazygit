package file

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var StageAllWithoutChangedFiles = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Pressing the stage-all key when there are no changed files says that there are none",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			IsEmpty().
			Press(keys.Files.ToggleStagedAll).
			Tap(func() {
				t.ExpectToast(Contains("No changed files"))
			})
	},
})
