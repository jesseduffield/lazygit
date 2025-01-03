package staging

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var DiffChangePanelSize = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Change the staged changes panel size",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("file", "first line\nsecond line")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			Focus().
			PressEnter()

		t.Views().Staging().
			IsFocused().
			PressPrimaryAction().
			Title(Equals("Unstaged changes")).
			Content(Contains("+second line").DoesNotContain("+first line")).
			PressTab()

		t.Views().StagingSecondary().
			IsFocused().
			Title(Equals("Staged changes")).
			Content(Contains("+first line").DoesNotContain("+second line")).
			Press(keys.Universal.NextPanelSize).
			Tap(func() {
				t.Views().AppStatus().
					IsInvisible()
				t.Views().Staging().
					IsVisible()
			}).
			Press(keys.Universal.NextPanelSize).
			Tap(func() {
				t.Views().AppStatus().
					IsInvisible()
				t.Views().Staging().
					IsInvisible()
			})
	},
})
