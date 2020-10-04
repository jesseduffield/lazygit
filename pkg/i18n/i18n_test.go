package i18n

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func getDummyLog() *logrus.Entry {
	log := logrus.New()
	log.Out = ioutil.Discard
	return log.WithField("test", "test")
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
