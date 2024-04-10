package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserConfigValidate_enums(t *testing.T) {
	type testCase struct {
		value string
		valid bool
	}

	scenarios := []struct {
		name      string
		setup     func(config *UserConfig, value string)
		testCases []testCase
	}{
		{
			name: "Gui.StatusPanelView",
			setup: func(config *UserConfig, value string) {
				config.Gui.StatusPanelView = value
			},
			testCases: []testCase{
				{value: "dashboard", valid: true},
				{value: "allBranchesLog", valid: true},
				{value: "", valid: false},
				{value: "invalid_value", valid: false},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			for _, testCase := range s.testCases {
				config := GetDefaultConfig()
				s.setup(config, testCase.value)
				err := config.Validate()

				if testCase.valid {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
				}
			}
		})
	}
}
