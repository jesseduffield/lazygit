package utils

import (
	"io/ioutil"

	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/sirupsen/logrus"
)

// NewDummyLog creates a new dummy Log for testing
func NewDummyLog() *logrus.Entry {
	log := logrus.New()
	log.Out = ioutil.Discard
	return log.WithField("test", "test")
}

func NewDummyCommon() *common.Common {
	tr := i18n.EnglishTranslationSet()
	return &common.Common{
		Log:        NewDummyLog(),
		Tr:         &tr,
		UserConfig: config.GetDefaultConfig(),
	}
}
