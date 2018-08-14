package i18n

import (
	"github.com/Sirupsen/logrus"
	"github.com/cloudfoundry/jibber_jabber"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Localizer will translate a message into the user's language
type Localizer struct {
	i18nLocalizer *i18n.Localizer
	language      string
	Log           *logrus.Logger
}

// NewLocalizer creates a new Localizer
func NewLocalizer(log *logrus.Logger) (*Localizer, error) {

	// detect the user's language
	userLang, err := jibber_jabber.DetectLanguage()
	if err != nil {
		return nil, err
	}
	log.Info("language: " + userLang)

	// create a i18n bundle that can be used to add translations and other things
	i18nBundle := &i18n.Bundle{DefaultLanguage: language.English}

	addBundles(i18nBundle)

	// return the new localizer that can be used to translate text
	i18nLocalizer := i18n.NewLocalizer(i18nBundle, userLang)

	localizer := &Localizer{
		i18nLocalizer: i18nLocalizer,
		language:      userLang,
		Log:           log,
	}

	return localizer, nil
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
func (l *Localizer) SLocalize(ID string, Other string) string {
	return l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    ID,
			Other: Other,
		},
	})
}

// TemplateLocalize allows the Other input to be dynamic
func (l *Localizer) TemplateLocalize(ID string, Other string, TemplateData map[string]interface{}) string {
	return l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    ID,
			Other: Other,
		},
		TemplateData: TemplateData,
	})
}

// GetLanguage returns the currently selected language, e.g 'en'
func (l *Localizer) GetLanguage() string {
	return l.language
}

// add translation file(s)
func addBundles(i18nBundle *i18n.Bundle) {
	addDutch(i18nBundle)
}
