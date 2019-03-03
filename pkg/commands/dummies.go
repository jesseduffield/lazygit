package commands

import (
	"io/ioutil"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// This file exports dummy constructors for use by tests in other packages

// NewDummyOSCommand creates a new dummy OSCommand for testing
func NewDummyOSCommand() *OSCommand {
	return NewOSCommand(NewDummyLog(), NewDummyAppConfig())
}

// NewDummyAppConfig creates a new dummy AppConfig for testing
func NewDummyAppConfig() *config.AppConfig {
	appConfig := &config.AppConfig{
		Name:        "lazygit",
		Version:     "unversioned",
		Commit:      "",
		BuildDate:   "",
		Debug:       false,
		BuildSource: "",
		UserConfig:  viper.New(),
	}
	_ = yaml.Unmarshal([]byte{}, appConfig.AppState)
	return appConfig
}

// NewDummyLog creates a new dummy Log for testing
func NewDummyLog() *logrus.Entry {
	log := logrus.New()
	log.Out = ioutil.Discard
	return log.WithField("test", "test")
}

// NewDummyGitCommand creates a new dummy GitCommand for testing
func NewDummyGitCommand() *GitCommand {
	return NewDummyGitCommandWithOSCommand(NewDummyOSCommand())
}

// NewDummyGitCommandWithOSCommand creates a new dummy GitCommand for testing
func NewDummyGitCommandWithOSCommand(osCommand *OSCommand) *GitCommand {
	return &GitCommand{
		Log:                NewDummyLog(),
		OSCommand:          osCommand,
		Tr:                 i18n.NewLocalizer(NewDummyLog()),
		Config:             NewDummyAppConfig(),
		getGlobalGitConfig: func(string) (string, error) { return "", nil },
		getLocalGitConfig:  func(string) (string, error) { return "", nil },
		removeFile:         func(string) error { return nil },
	}
}
