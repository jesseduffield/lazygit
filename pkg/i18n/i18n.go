package lang

import (
	"github.com/cloudfoundry/jibber_jabber"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// the function to setup the localizer
func getlocalizer() *i18n.Localizer {

	// detect the user's language
	userLang, _ := jibber_jabber.DetectLanguage()

	// create a i18n bundle that can be used to add translations and other things
	var i18nObject = &i18n.Bundle{DefaultLanguage: language.English}

	// add translation file(s)
	i18nObject = addDutch(i18nObject)

	// return the new localizer that can be used to translate text
	return i18n.NewLocalizer(i18nObject, userLang)
}

// setup the localizer for later use
var localizer = getlocalizer()

// Localize handels the translations
// expects i18n.LocalizeConfig as input: https://godoc.org/github.com/nicksnyder/go-i18n/v2/i18n#Localizer.MustLocalize
// output: translated string
func Localize(config *i18n.LocalizeConfig) string {
	return localizer.MustLocalize(config)
}

// SLocalize (short localize) is for 1 line localizations
// ID: The id that is used in the .toml translation files
// Other: the default message it needs to return if there is no translation found or the system is english
func SLocalize(ID string, Other string) string {
	return Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    ID,
			Other: Other,
		},
	})
}
