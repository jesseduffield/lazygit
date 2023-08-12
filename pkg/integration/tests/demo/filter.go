package demo

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Filter = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Filter branches",
	ExtraCmdArgs: []string{},
	Skip:         false,
	IsDemo:       true,
	SetupConfig: func(config *config.AppConfig) {
		setDefaultDemoConfig(config)
	},
	SetupRepo: func(shell *Shell) {
		shell.CreateNCommitsWithRandomMessages(30)
		shell.NewBranch("feature/user-authentication")
		shell.NewBranch("feature/payment-processing")
		shell.NewBranch("feature/search-functionality")
		shell.NewBranch("feature/mobile-responsive")
		shell.NewBranch("bugfix/fix-login-issue")
		shell.NewBranch("bugfix/fix-crash-bug")
		shell.NewBranch("bugfix/fix-validation-error")
		shell.NewBranch("refactor/improve-performance")
		shell.NewBranch("refactor/code-cleanup")
		shell.NewBranch("refactor/extract-method")
		shell.NewBranch("docs/update-readme")
		shell.NewBranch("docs/add-user-guide")
		shell.NewBranch("docs/api-documentation")
		shell.NewBranch("experiment/new-feature-idea")
		shell.NewBranch("experiment/try-new-library")
		shell.NewBranch("chore/update-dependencies")
		shell.NewBranch("chore/add-test-cases")
		shell.NewBranch("chore/migrate-database")
		shell.NewBranch("hotfix/critical-bug")
		shell.NewBranch("hotfix/security-patch")
		shell.NewBranch("feature/social-media-integration")
		shell.NewBranch("feature/email-notifications")
		shell.NewBranch("feature/admin-panel")
		shell.NewBranch("feature/analytics-dashboard")
		shell.NewBranch("bugfix/fix-registration-flow")
		shell.NewBranch("bugfix/fix-payment-bug")
		shell.NewBranch("refactor/improve-error-handling")
		shell.NewBranch("refactor/optimize-database-queries")
		shell.NewBranch("docs/improve-tutorials")
		shell.NewBranch("docs/add-faq-section")
		shell.NewBranch("experiment/try-alternative-algorithm")
		shell.NewBranch("experiment/implement-design-concept")
		shell.NewBranch("chore/update-documentation")
		shell.NewBranch("chore/improve-test-coverage")
		shell.NewBranch("chore/cleanup-codebase")
		shell.NewBranch("hotfix/critical-security-vulnerability")
		shell.NewBranch("hotfix/fix-production-issue")
		shell.NewBranch("feature/integrate-third-party-api")
		shell.NewBranch("feature/image-upload-functionality")
		shell.NewBranch("feature/localization-support")
		shell.NewBranch("feature/chat-feature")
		shell.NewBranch("bugfix/fix-broken-link")
		shell.NewBranch("bugfix/fix-css-styling")
		shell.NewBranch("refactor/improve-logging")
		shell.NewBranch("refactor/extract-reusable-component")
		shell.NewBranch("docs/add-changelog")
		shell.NewBranch("docs/update-api-reference")
		shell.NewBranch("experiment/implement-new-design")
		shell.NewBranch("experiment/try-different-architecture")
		shell.NewBranch("chore/clean-up-git-history")
		shell.NewBranch("chore/update-environment-configuration")
		shell.CreateFileAndAdd("env_config.rb", "EnvConfig.call(false)\n")
		shell.Commit("Update env config")
		shell.CreateFileAndAdd("env_config.rb", "# Turns out we need to pass true for this to work\nEnvConfig.call(true)\n")
		shell.Commit("Fix env config issue")
		shell.Checkout("docs/add-faq-section")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.SetCaptionPrefix("Fuzzy filter branches")
		t.Wait(1000)

		t.Views().Branches().
			Focus().
			Wait(500).
			Press(keys.Universal.StartSearch).
			Tap(func() {
				t.Wait(500)

				t.ExpectSearch().Type("environ").Confirm()
			}).
			Wait(500).
			PressEnter()
	},
})
