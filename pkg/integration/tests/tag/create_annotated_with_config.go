package tag

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CreateAnnotatedWithConfig = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Create an annotated tag when the alwaysAnnotate config option is enabled",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.GetUserConfig().Git.Tag.AlwaysAnnotate = true
	},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Tags().
			Focus().
			IsEmpty().
			Press(keys.Universal.New).
			Tap(func() {
				t.ExpectPopup().CommitMessagePanel().
					Title(Equals("Tag name")).
					Type("new-tag").
					Confirm()
			}).
			Lines(
				MatchesRegexp(`new-tag.*`).IsSelected(),
			).
			PressEnter().
			Tap(func() {
				t.Git().TagIsAnnotated("new-tag")
			})
	},
})
