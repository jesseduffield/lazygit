package i18n

import (
	"strings"

	"github.com/cloudfoundry/jibber_jabber"
	"github.com/go-errors/errors"
	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"
)

// Localizer will translate a message into the user's language
type Localizer struct {
	Log *logrus.Entry
	S   TranslationSet
}

func NewTranslationSetFromConfig(log *logrus.Entry, configLanguage string) (*TranslationSet, error) {
	if configLanguage == "auto" {
		language := detectLanguage(jibber_jabber.DetectLanguage)
		return NewTranslationSet(log, language), nil
	}

	for key := range GetTranslationSets() {
		if key == configLanguage {
			return NewTranslationSet(log, configLanguage), nil
		}
	}

	return NewTranslationSet(log, "en"), errors.New("Language not found: " + configLanguage)
}

func NewTranslationSet(log *logrus.Entry, language string) *TranslationSet {
	log.Info("language: " + language)

	baseSet := EnglishTranslationSet()

	for languageCode, translationSet := range GetTranslationSets() {
		if strings.HasPrefix(language, languageCode) {
			_ = mergo.Merge(&baseSet, translationSet, mergo.WithOverride)
		}
	}
	return &baseSet
}

// GetTranslationSets gets all the translation sets, keyed by language code
func GetTranslationSets() map[string]TranslationSet {
	return map[string]TranslationSet{
		"pl": polishTranslationSet(),
		"nl": dutchTranslationSet(),
		"en": EnglishTranslationSet(),
		"zh": chineseTranslationSet(),
		"ja": japaneseTranslationSet(),
		"ko": koreanTranslationSet(),
	}
}

// detectLanguage extracts user language from environment
func detectLanguage(langDetector func() (string, error)) string {
	if userLang, err := langDetector(); err == nil {
		return userLang
	}

	return "C"
}
