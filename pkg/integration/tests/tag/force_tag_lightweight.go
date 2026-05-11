package tag

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ForceTagLightweight = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Overwrite a lightweight tag that already exists",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("first commit")
		shell.CreateLightweightTag("new-tag", "HEAD")
		shell.EmptyCommit("second commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			Focus().
			Lines(
				Contains("second commit").IsSelected(),
				Contains("new-tag").Contains("first commit"),
			).
			Press(keys.Commits.CreateTag).
			Tap(func() {
				t.ExpectPopup().CommitMessagePanel().
					Title(Equals("Tag name")).
					Type("new-tag").
					Confirm()
			}).
			Tap(func() {
				t.ExpectPopup().Confirmation().
					Title(Equals("Force Tag")).
					Content(Contains("The tag 'new-tag' exists already. Press <esc> to cancel, or <enter> to overwrite.")).
					Confirm()
			}).
			Lines(
				Contains("new-tag").Contains("second commit"),
				DoesNotContain("new-tag").Contains("first commit"),
			)
	},
})
