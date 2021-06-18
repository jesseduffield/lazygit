package i18n

import (
	"strings"

	"github.com/cloudfoundry/jibber_jabber"
	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"
)

// Localizer will translate a message into the user's language
type Localizer struct {
	Log *logrus.Entry
	S   TranslationSet
}

// NewTranslationSet creates a new Localizer
func NewTranslationSet(log *logrus.Entry) *TranslationSet {
	userLang := detectLanguage(jibber_jabber.DetectLanguage)

	log.Info("language: " + userLang)

	baseSet := EnglishTranslationSet()

	for languageCode, translationSet := range GetTranslationSets() {
		if strings.HasPrefix(userLang, languageCode) {
			_ = mergo.Merge(baseSet, translationSet, mergo.WithOverride)
		}
	}

	return baseSet
}

// GetTranslationSets gets all the translation sets, keyed by language code
func GetTranslationSets() map[string]*TranslationSet {
	return map[string]*TranslationSet{
		"pl": polishTranslationSet(),
		"nl": dutchTranslationSet(),
		"en": EnglishTranslationSet(),
	}
}

// detectLanguage extracts user language from environment
func detectLanguage(langDetector func() (string, error)) string {
	if userLang, err := langDetector(); err == nil {
		return userLang
	}

	return "C"
}
