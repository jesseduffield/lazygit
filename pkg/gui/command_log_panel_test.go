package gui

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/stretchr/testify/assert"
)

func TestRandomTipsResolveLocalizedTemplates(t *testing.T) {
	translationSet := i18n.EnglishTranslationSet()
	translationSet.RandomTipForcePush = "按 '{{.pushKey}}' 强制推送"

	tips := randomTips(translationSet, config.KeybindingConfig{
		Universal: config.KeybindingUniversalConfig{
			Push: config.Keybinding{"P"},
		},
	})

	assert.Equal(t, "按 'P' 强制推送", tips[0])
}
