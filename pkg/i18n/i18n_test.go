package i18n

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func getDummyLog() *logrus.Entry {
	log := logrus.New()
	log.Out = ioutil.Discard
	return log.WithField("test", "test")
}

// TestNewLocalizer is a function.
func TestNewLocalizer(t *testing.T) {
	assert.NotNil(t, NewLocalizer(getDummyLog()))
}

// TestDetectLanguage is a function.
func TestDetectLanguage(t *testing.T) {
	type scenario struct {
		langDetector func() (string, error)
		expected     string
	}

	scenarios := []scenario{
		{
			func() (string, error) {
				return "", fmt.Errorf("An error occurred")
			},
			"C",
		},
		{
			func() (string, error) {
				return "en", nil
			},
			"en",
		},
	}

	for _, s := range scenarios {
		assert.EqualValues(t, s.expected, detectLanguage(s.langDetector))
	}
}

// TestLocalizer is a function.
func TestLocalizer(t *testing.T) {
	type scenario struct {
		userLang string
		test     func(*Localizer)
	}

	scenarios := []scenario{
		{
			"C",
			func(l *Localizer) {
				assert.EqualValues(t, "C", l.GetLanguage())
				assert.Equal(t, "Diff", l.Localize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID: "DiffTitle",
					},
				}))
				assert.Equal(t, "Diff", l.SLocalize("DiffTitle"))
				assert.Equal(t, "Are you sure you want to delete the branch test?", l.TemplateLocalize("DeleteBranchMessage", Teml{"selectedBranchName": "test"}))
			},
		},
		{
			"nl",
			func(l *Localizer) {
				assert.EqualValues(t, "nl", l.GetLanguage())
				assert.Equal(t, "Diff", l.Localize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID: "DiffTitle",
					},
				}))
				assert.Equal(t, "Diff", l.SLocalize("DiffTitle"))
				assert.Equal(t, "Weet je zeker dat je branch test wilt verwijderen?", l.TemplateLocalize("DeleteBranchMessage", Teml{"selectedBranchName": "test"}))
			},
		},
	}

	for _, s := range scenarios {
		s.test(setupLocalizer(getDummyLog(), s.userLang))
	}
}
