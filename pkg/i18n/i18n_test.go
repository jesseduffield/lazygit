package i18n

import (
	"fmt"
	"io"
	"runtime"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

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

// Can't use utils.NewDummyLog() because of a cyclic dependency
func newDummyLog() *logrus.Entry {
	log := logrus.New()
	log.Out = io.Discard
	return log.WithField("test", "test")
}

func TestNewTranslationSetFromConfig(t *testing.T) {
	if runtime.GOOS == "windows" {
		// These tests are based on setting the LANG environment variable, which
		// isn't respected on Windows.
		t.Skip("Skipping test on Windows")
	}

	scenarios := []struct {
		name           string
		configLanguage string
		envLanguage    string
		expected       string
		expectedErr    bool
	}{
		{
			name:           "configLanguage is nl",
			configLanguage: "nl",
			envLanguage:    "en_US",
			expected:       "nl",
			expectedErr:    false,
		},
		{
			name:           "configLanguage is an unsupported language",
			configLanguage: "xy",
			envLanguage:    "en_US",
			expectedErr:    true,
		},
		{
			name:           "auto-detection without LANG set",
			configLanguage: "auto",
			envLanguage:    "",
			expected:       "en",
			expectedErr:    false,
		},
		{
			name:           "auto-detection with LANG set to nl_NL",
			configLanguage: "auto",
			envLanguage:    "nl_NL",
			expected:       "nl",
			expectedErr:    false,
		},
		{
			name:           "auto-detection with LANG set to zh-CN",
			configLanguage: "auto",
			envLanguage:    "zh-CN",
			expected:       "zh-CN",
			expectedErr:    false,
		},
		{
			name:           "auto-detection with LANG set to an unsupported language",
			configLanguage: "auto",
			envLanguage:    "xy_XY",
			expected:       "en",
			expectedErr:    false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			log := newDummyLog()
			t.Setenv("LANG", s.envLanguage)
			actualTranslationSet, err := NewTranslationSetFromConfig(log, s.configLanguage)
			if s.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				expectedTranslationSet, _ := newTranslationSet(log, s.expected)
				assert.Equal(t, expectedTranslationSet, actualTranslationSet)
			}
		})
	}
}
