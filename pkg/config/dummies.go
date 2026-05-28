package config

import (
	"gopkg.in/yaml.v3"
)

// NewDummyAppConfig creates a new dummy AppConfig for testing
func NewDummyAppConfig() *AppConfig {
	userConfig := GetDefaultConfig()
	userConfig.Keybinding.MergeLegacyAltKeybindings()
	appConfig := &AppConfig{
		name:       "lazygit",
		version:    "unversioned",
		debug:      false,
		userConfig: userConfig,
		appState:   &AppState{},
	}
	_ = yaml.Unmarshal([]byte{}, appConfig.appState)
	return appConfig
}
