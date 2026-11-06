package config

// NewDummyAppConfig creates a new dummy AppConfig for testing
func NewDummyAppConfig() *AppConfig {
	userConfig := GetDefaultConfig()
	userConfig.Keybinding.MergeLegacyAltKeybindings()
	appConfig := &AppConfig{
		name:                   "lazygit",
		version:                "unversioned",
		debug:                  false,
		userConfig:             userConfig,
		githubPullRequestCache: newGithubPullRequestCache(""),
	}
	return appConfig
}
