package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitWithFallthroughPrefix = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Commit with multiple CommitPrefixConfig",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.GetUserConfig().Git.CommitPrefix = []config.CommitPrefixConfig{
			{Pattern: "^doesntmatch-(\\w+).*", Replace: "[BAD $1]: "},
			{Pattern: "^\\w+\\/(\\w+-\\w+).*", Replace: "[GOOD $1]: "},
		}
		cfg.GetUserConfig().Git.CommitPrefixes = map[string][]config.CommitPrefixConfig{
			"DifferentProject": {{Pattern: "^otherthatdoesn'tmatch-(\\w+).*", Replace: "[BAD $1]: "}},
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("feature/TEST-001")
		shell.CreateFile("test-commit-prefix", "This is foo bar")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			IsEmpty()

		t.Views().Files().
			IsFocused().
			PressPrimaryAction().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Title(Equals("Commit summary")).
			InitialText(Equals("[GOOD TEST-001]: ")).
			Type("my commit message").
			Cancel()

		t.Views().Files().
			IsFocused().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Title(Equals("Commit summary")).
			InitialText(Equals("[GOOD TEST-001]: my commit message")).
			Type(". Added something else").
			Confirm()

		t.Views().Commits().Focus()
		t.Views().Main().Content(Contains("[GOOD TEST-001]: my commit message. Added something else"))
	},
})
