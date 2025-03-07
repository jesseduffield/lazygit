package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitSkipHooks = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Commit with skip hook using CommitChangesWithoutHook",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("feature/TEST-002")
		shell.CreateFile("test-wip-commit-prefix", "This is foo bar")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			IsEmpty()

		t.Views().Files().
			IsFocused().
			PressPrimaryAction().
			Press(keys.Files.CommitChangesWithoutHook)

		t.ExpectPopup().CommitMessagePanel().
			Title(Equals("Commit summary")).
			Type("foo bar").
			Confirm()

		t.Views().Commits().Focus()
		t.Views().Main().Content(Contains("foo bar"))
		t.Views().Extras().Content(Contains("--no-verify"))
	},
})
