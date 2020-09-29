package config

import (
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// NewDummyAppConfig creates a new dummy AppConfig for testing
func NewDummyAppConfig() *AppConfig {
	userConfig := viper.New()
	userConfig.SetConfigType("yaml")
	if err := LoadDefaults(userConfig, GetDefaultConfig()); err != nil {
		panic(err)
	}
	appConfig := &AppConfig{
		Name:        "lazygit",
		Version:     "unversioned",
		Commit:      "",
		BuildDate:   "",
		Debug:       false,
		BuildSource: "",
		UserConfig:  userConfig,
	}
	_ = yaml.Unmarshal([]byte{}, appConfig.AppState)
	return appConfig
}
