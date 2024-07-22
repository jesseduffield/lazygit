package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var CommitWithNonMatchingBranchName = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Commit with defined config commitPrefixes",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(testConfig *config.AppConfig) {
		testConfig.UserConfig.Git.CommitPrefix = &config.CommitPrefixConfig{
			Pattern: "^\\w+\\/(\\w+-\\w+).*",
			Replace: "[$1]: ",
		}
	},
	SetupRepo: func(shell *Shell) {
		shell.NewBranch("branchnomatch")
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
			InitialText(Equals(""))
	},
})
