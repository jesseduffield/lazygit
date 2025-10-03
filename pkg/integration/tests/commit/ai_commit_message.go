package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var AICommitMessage = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Generate AI commit message for staged changes",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		// Configure AI settings for testing
		config.UserConfig.AI.Provider = "openai"
		config.UserConfig.AI.Model = "gpt-4o-mini"
		config.UserConfig.AI.Temperature = 0.2
		config.UserConfig.AI.MaxTokens = 300
		config.UserConfig.AI.StagedOnly = true
		config.UserConfig.AI.CommitStyle = "conventional"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("feature.go", `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`)
		shell.CreateFileAndAdd("README.md", "# My Project\n\nThis is a test project.")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		// Note: This test will only work if the user has a valid API key set
		// In a real CI environment, this test might be skipped or use a mock server
		
		t.Views().Commits().
			IsEmpty()

		t.Views().Files().
			IsFocused().
			Lines(
				Contains("feature.go"),
				Contains("README.md"),
			).
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Title(Equals("Commit summary")).
			Press(keys.Universal.OpenCommitMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Commit menu")).
			Select(Contains("Generate AI commit message")).
			Confirm()

		// The AI generation might take some time, so we wait for the status
		// In a real test, we might mock the AI response for faster execution
		t.Views().CommitMessage().
			Content(MatchesRegexp(`^(feat|fix|docs|style|refactor|perf|test|chore|build|ci)(\(.+\))?: .+`))

		// Verify we can still edit the generated message
		t.Views().CommitMessage().
			Type(" - updated by user").
			Confirm()

		t.Views().Commits().
			Lines(
				MatchesRegexp(`^(feat|fix|docs|style|refactor|perf|test|chore|build|ci)(\(.+\))?: .+ - updated by user`),
			)
	},
})

var AICommitMessageWithError = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Handle AI commit message generation errors gracefully",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		// Configure AI with invalid settings to trigger an error
		config.UserConfig.AI.Provider = "openai"
		config.UserConfig.AI.Model = "" // Missing model should cause validation error
		config.UserConfig.AI.Temperature = 0.2
		config.UserConfig.AI.MaxTokens = 300
		config.UserConfig.AI.StagedOnly = true
		config.UserConfig.AI.CommitStyle = "conventional"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("test.go", "package main")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Press(keys.Universal.OpenCommitMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Commit menu")).
			Select(Contains("Generate AI commit message")).
			Confirm()

		// Should show an error alert
		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("ai.model is required")).
			Confirm()

		// Should return to the commit message panel
		t.ExpectPopup().CommitMessagePanel().
			Type("manual commit message").
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("manual commit message"),
			)
	},
})

var AICommitMessageNoStagedFiles = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Handle AI commit message generation when no files are staged",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.AI.Provider = "openai"
		config.UserConfig.AI.Model = "gpt-4o-mini"
		config.UserConfig.AI.Temperature = 0.2
		config.UserConfig.AI.MaxTokens = 300
		config.UserConfig.AI.StagedOnly = true
		config.UserConfig.AI.CommitStyle = "conventional"
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile("unstaged.go", "package main")
		// Don't stage the file
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Lines(
				Contains("unstaged.go"),
			).
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Press(keys.Universal.OpenCommitMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Commit menu")).
			Select(Contains("Generate AI commit message")).
			Confirm()

		// Should show an error about no staged files
		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("No files staged")).
			Confirm()

		// Should return to the commit message panel
		t.ExpectPopup().CommitMessagePanel().
			Cancel()
	},
})

var AICommitMessagePlainStyle = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Generate AI commit message with plain style",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.AI.Provider = "openai"
		config.UserConfig.AI.Model = "gpt-4o-mini"
		config.UserConfig.AI.Temperature = 0.2
		config.UserConfig.AI.MaxTokens = 300
		config.UserConfig.AI.StagedOnly = true
		config.UserConfig.AI.CommitStyle = "plain" // Use plain style instead of conventional
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateFileAndAdd("simple.txt", "Hello World")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Press(keys.Universal.OpenCommitMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Commit menu")).
			Select(Contains("Generate AI commit message")).
			Confirm()

		// For plain style, we expect a simple commit message without conventional format
		t.Views().CommitMessage().
			Content(Not(MatchesRegexp(`^(feat|fix|docs|style|refactor|perf|test|chore|build|ci)(\(.+\))?: .+`))).
			Content(Not(Equals(""))) // Should not be empty

		t.Views().CommitMessage().
			Confirm()

		t.Views().Commits().
			TopLines(
				Not(MatchesRegexp(`^(feat|fix|docs|style|refactor|perf|test|chore|build|ci)(\(.+\))?: .+`)),
			)
	},
})

var AICommitMessageLargeDiff = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Handle AI commit message generation with large diff that exceeds size limit",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.AI.Provider = "openai"
		config.UserConfig.AI.Model = "gpt-4o-mini"
		config.UserConfig.AI.Temperature = 0.2
		config.UserConfig.AI.MaxTokens = 300
		config.UserConfig.AI.StagedOnly = true
		config.UserConfig.AI.CommitStyle = "conventional"
	},
	SetupRepo: func(shell *Shell) {
		// Create a very large file to exceed the diff size limit
		largeContent := ""
		for i := 0; i < 2000; i++ {
			largeContent += "This is a very long line that will make the diff very large when staged.\n"
		}
		shell.CreateFileAndAdd("large_file.txt", largeContent)
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Files().
			IsFocused().
			Press(keys.Files.CommitChanges)

		t.ExpectPopup().CommitMessagePanel().
			Press(keys.Universal.OpenCommitMenu)

		t.ExpectPopup().Menu().
			Title(Equals("Commit menu")).
			Select(Contains("Generate AI commit message")).
			Confirm()

		// Should show an error about diff being too large
		t.ExpectPopup().Alert().
			Title(Equals("Error")).
			Content(Contains("diff too large for AI processing")).
			Confirm()

		// Should return to the commit message panel
		t.ExpectPopup().CommitMessagePanel().
			Type("manual commit for large file").
			Confirm()

		t.Views().Commits().
			Lines(
				Contains("manual commit for large file"),
			)
	},
})
