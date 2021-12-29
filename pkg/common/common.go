package common

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/sirupsen/logrus"
)

// Commonly used things wrapped into one struct for convenience when passing it around
type Common struct {
	Log        *logrus.Entry
	Tr         *i18n.TranslationSet
	UserConfig *config.UserConfig
	Debug      bool
}
