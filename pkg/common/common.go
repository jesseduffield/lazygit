package common

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// Commonly used things wrapped into one struct for convenience when passing it around
type Common struct {
	Log        *logrus.Entry
	Tr         *i18n.TranslationSet
	UserConfig *config.UserConfig
	Debug      bool
	// for interacting with the filesystem. We use afero rather than the default
	// `os` package for the sake of mocking the filesystem in tests
	Fs afero.Fs
}
