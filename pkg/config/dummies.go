package config

import (
	yaml "github.com/jesseduffield/yaml"
)

// NewDummyAppConfig creates a new dummy AppConfig for testing
func NewDummyAppConfig() *AppConfig {
	appConfig := &AppConfig{
		Name:       "lazygit",
		Version:    "unversioned",
		Debug:      false,
		UserConfig: GetDefaultConfig(),
		AppState:   &AppState{},
	}
	_ = yaml.Unmarshal([]byte{}, appConfig.AppState)
	return appConfig
}
