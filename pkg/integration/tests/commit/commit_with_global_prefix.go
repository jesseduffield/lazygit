package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitWithGlobalPrefix = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Commit with defined config commitPrefix",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(testConfig *config.AppConfig) {
		testConfig.UserConfig.Git.CommitPrefix = &config.CommitPrefixConfig{Pattern: "^\\w+\\/(\\w+-\\w+).*", Replace: "[$1]: "}
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
			InitialText(Equals("[TEST-001]: ")).
			Type("my commit message").
			Cancel()

		t.Views().Files().
			IsFocused().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Title(Equals("Commit summary")).
			InitialText(Equals("[TEST-001]: my commit message")).
			Type(". Added something else").
			Confirm()

		t.Views().Commits().Focus()
		t.Views().Main().Content(Contains("[TEST-001]: my commit message. Added something else"))
	},
})
