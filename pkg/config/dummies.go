package config

import (
	"gopkg.in/yaml.v3"
)

// NewDummyAppConfig creates a new dummy AppConfig for testing
func NewDummyAppConfig() *AppConfig {
	appConfig := &AppConfig{
		name:       "lazygit",
		version:    "unversioned",
		debug:      false,
		userConfig: GetDefaultConfig(),
		appState:   &AppState{},
	}
	_ = yaml.Unmarshal([]byte{}, appConfig.appState)
	return appConfig
}
