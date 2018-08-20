package i18n

import (
	"os"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewLocalizer(t *testing.T) {
	type scenario struct {
		setup    func()
		test     func(*Localizer, error)
		teardown func()
	}

	LCALL := os.Getenv("LC_ALL")
	LANG := os.Getenv("LANG")

	scenarios := []scenario{
		{
			func() {
				os.Setenv("LC_ALL", "")
				os.Setenv("LANG", "")
			},
			func(l *Localizer, err error) {
				assert.EqualValues(t, "C", l.GetLanguage())
			},
			func() {
				os.Setenv("LC_ALL", LCALL)
				os.Setenv("LANG", LANG)
			},
		},
		{
			func() {
				os.Setenv("LC_ALL", "whatever")
				os.Setenv("LANG", "whatever")
			},
			func(l *Localizer, err error) {
				assert.NoError(t, err)

				assert.EqualValues(t, "whatever", l.GetLanguage())
				assert.Equal(t, "Diff", l.Localize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID: "DiffTitle",
					},
				}))
				assert.Equal(t, "Diff", l.SLocalize("DiffTitle"))
				assert.Equal(t, "Are you sure you want delete the branch test ?", l.TemplateLocalize("DeleteBranchMessage", Teml{"selectedBranchName": "test"}))
			},
			func() {
				os.Setenv("LC_ALL", LCALL)
				os.Setenv("LANG", LANG)
			},
		},
		{
			func() {
				os.Setenv("LC_ALL", "nl")
			},
			func(l *Localizer, err error) {
				assert.NoError(t, err)

				assert.EqualValues(t, "nl", l.GetLanguage())
				assert.Equal(t, "Diff", l.Localize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID: "DiffTitle",
					},
				}))
				assert.Equal(t, "Diff", l.SLocalize("DiffTitle"))
				assert.Equal(t, "Weet je zeker dat je test branch wil verwijderen?", l.TemplateLocalize("DeleteBranchMessage", Teml{"selectedBranchName": "test"}))
			},
			func() {
				os.Setenv("LC_ALL", LCALL)
			},
		},
	}

	for _, s := range scenarios {
		s.setup()
		s.test(NewLocalizer(logrus.New()))
		s.teardown()
	}
}
