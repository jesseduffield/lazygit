package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"

	"github.com/cloudfoundry/jibber_jabber"
	"github.com/go-errors/errors"
	"github.com/imdario/mergo"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func NewTranslationSetFromConfig(log *logrus.Entry, configLanguage string) (*TranslationSet, error) {
	languageCodes, err := getSupportedLanguageCodes()
	if err != nil {
		return nil, err
	}

	if configLanguage == "auto" {
		language := detectLanguage(jibber_jabber.DetectIETF)
		for _, languageCode := range languageCodes {
			if strings.HasPrefix(language, languageCode) {
				return newTranslationSet(log, languageCode)
			}
		}

		// Detecting a language that we don't have a translation for is not an
		// error, we'll just use English.
		return EnglishTranslationSet(), nil
	}

	if configLanguage == "en" {
		return EnglishTranslationSet(), nil
	}

	for _, key := range languageCodes {
		if key == configLanguage {
			return newTranslationSet(log, configLanguage)
		}
	}

	// Configuring a language that we don't have a translation for *is* an
	// error, though.
	return nil, errors.New("Language not found: " + configLanguage)
}

func newTranslationSet(log *logrus.Entry, language string) (*TranslationSet, error) {
	log.Info("language: " + language)

	baseSet := EnglishTranslationSet()

	if language != "en" {
		translationSet, err := readLanguageFile(language)
		if err != nil {
			return nil, err
		}
		err = mergo.Merge(baseSet, *translationSet, mergo.WithOverride)
		if err != nil {
			return nil, err
		}
	}

	return baseSet, nil
}

//go:embed translations/*.json
var embedFS embed.FS

// getSupportedLanguageCodes gets all the supported language codes.
// Note: this doesn't include "en"
func getSupportedLanguageCodes() ([]string, error) {
	dirEntries, err := embedFS.ReadDir("translations")
	if err != nil {
		return nil, err
	}
	return lo.Map(dirEntries, func(entry fs.DirEntry, _ int) string {
		return strings.TrimSuffix(entry.Name(), ".json")
	}), nil
}

func readLanguageFile(languageCode string) (*TranslationSet, error) {
	jsonData, err := embedFS.ReadFile(fmt.Sprintf("translations/%s.json", languageCode))
	if err != nil {
		return nil, err
	}
	var translationSet TranslationSet
	err = json.Unmarshal(jsonData, &translationSet)
	if err != nil {
		return nil, err
	}
	return &translationSet, nil
}

// GetTranslationSets gets all the translation sets, keyed by language code
// This includes "en".
func GetTranslationSets() (map[string]*TranslationSet, error) {
	languageCodes, err := getSupportedLanguageCodes()
	if err != nil {
		return nil, err
	}

	result := make(map[string]*TranslationSet)
	result["en"] = EnglishTranslationSet()

	for _, languageCode := range languageCodes {
		translationSet, err := readLanguageFile(languageCode)
		if err != nil {
			return nil, err
		}
		result[languageCode] = translationSet
	}

	return result, nil
}

// detectLanguage extracts user language from environment
func detectLanguage(langDetector func() (string, error)) string {
	if userLang, err := langDetector(); err == nil {
		return userLang
	}

	return "C"
}
