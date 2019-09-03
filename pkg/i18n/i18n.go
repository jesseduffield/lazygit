package i18n

import (
	"github.com/cloudfoundry/jibber_jabber"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

// Teml is short for template used to make the required map[string]interface{} shorter when using gui.Tr.SLocalize and gui.Tr.TemplateLocalize
type Teml map[string]interface{}

// Localizer will translate a message into the user's language
type Localizer struct {
	i18nLocalizer *i18n.Localizer
	language      string
	Log           *logrus.Entry
}

// NewLocalizer creates a new Localizer
func NewLocalizer(log *logrus.Entry) *Localizer {
	userLang := detectLanguage(jibber_jabber.DetectLanguage)

	log.Info("language: " + userLang)

	return setupLocalizer(log, userLang)
}

// Localize handels the translations
// expects i18n.LocalizeConfig as input: https://godoc.org/github.com/nicksnyder/go-i18n/v2/i18n#Localizer.MustLocalize
// output: translated string
func (l *Localizer) Localize(config *i18n.LocalizeConfig) string {
	return l.i18nLocalizer.MustLocalize(config)
}

// SLocalize (short localize) is for 1 line localizations
// ID: The id that is used in the .toml translation files
// Other: the default message it needs to return if there is no translation found or the system is english
func (l *Localizer) SLocalize(ID string) string {
	return l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: ID,
		},
	})
}

// TemplateLocalize allows the Other input to be dynamic
func (l *Localizer) TemplateLocalize(ID string, TemplateData map[string]interface{}) string {
	return l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: ID,
		},
		TemplateData: TemplateData,
	})
}

// GetLanguage returns the currently selected language, e.g 'en'
func (l *Localizer) GetLanguage() string {
	return l.language
}

// add translation file(s)
func addBundles(log *logrus.Entry, i18nBundle *i18n.Bundle) {
	fs := []func(*i18n.Bundle) error{
		addPolish,
		addDutch,
		addEnglish,
	}

	for _, f := range fs {
		if err := f(i18nBundle); err != nil {
			log.Fatal(err)

		}
	}
}

// detectLanguage extracts user language from environment
func detectLanguage(langDetector func() (string, error)) string {
	if userLang, err := langDetector(); err == nil {
		return userLang
	}

	return "C"
}

// setupLocalizer creates a new localizer using given userLang
func setupLocalizer(log *logrus.Entry, userLang string) *Localizer {
	// create a i18n bundle that can be used to add translations and other things
	i18nBundle := i18n.NewBundle(language.English)

	addBundles(log, i18nBundle)

	// return the new localizer that can be used to translate text
	i18nLocalizer := i18n.NewLocalizer(i18nBundle, userLang)

	return &Localizer{
		i18nLocalizer: i18nLocalizer,
		language:      userLang,
		Log:           log,
	}
}
