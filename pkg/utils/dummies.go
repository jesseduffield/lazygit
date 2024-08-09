package utils

import (
	"io"

	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// NewDummyLog creates a new dummy Log for testing
func NewDummyLog() *logrus.Entry {
	log := logrus.New()
	log.Out = io.Discard
	return log.WithField("test", "test")
}

func NewDummyCommon() *common.Common {
	tr := i18n.EnglishTranslationSet()
	cmn := &common.Common{
		Log: NewDummyLog(),
		Tr:  tr,
		Fs:  afero.NewOsFs(),
	}
	cmn.SetUserConfig(config.GetDefaultConfig())
	return cmn
}

func NewDummyCommonWithUserConfigAndAppState(userConfig *config.UserConfig, appState *config.AppState) *common.Common {
	tr := i18n.EnglishTranslationSet()
	cmn := &common.Common{
		Log:      NewDummyLog(),
		Tr:       tr,
		AppState: appState,
		// TODO: remove dependency on actual filesystem in tests and switch to using
		// in-memory for everything
		Fs: afero.NewOsFs(),
	}
	cmn.SetUserConfig(userConfig)
	return cmn
}
