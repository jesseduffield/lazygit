package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestSidePanelMaxWidth(t *testing.T) {
	for _, test := range []struct {
		name  string
		yaml  string
		width int
		error string
	}{
		{name: "default", yaml: "{}"},
		{name: "unlimited", yaml: "gui:\n  sidePanelMaxWidth: 0"},
		{name: "positive", yaml: "gui:\n  sidePanelMaxWidth: 60", width: 60},
		{name: "negative", yaml: "gui:\n  sidePanelMaxWidth: -1", width: -1, error: "gui.sidePanelMaxWidth must not be negative."},
	} {
		t.Run(test.name, func(t *testing.T) {
			config := GetDefaultConfig()
			assert.NoError(t, yaml.Unmarshal([]byte(test.yaml), config))
			assert.Equal(t, test.width, config.Gui.SidePanelMaxWidth)
			if test.error == "" {
				assert.NoError(t, config.Validate())
			} else {
				assert.EqualError(t, config.Validate(), test.error)
			}
		})
	}
}
