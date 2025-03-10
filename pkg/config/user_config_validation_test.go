package config

import (
	"strings"
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
		{
			name: "Keybindings",
			setup: func(config *UserConfig, value string) {
				config.Keybinding.Universal.Quit = value
			},
			testCases: []testCase{
				{value: "", valid: true},
				{value: "<disabled>", valid: true},
				{value: "q", valid: true},
				{value: "<c-c>", valid: true},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "JumpToBlock keybinding",
			setup: func(config *UserConfig, value string) {
				config.Keybinding.Universal.JumpToBlock = strings.Split(value, ",")
			},
			testCases: []testCase{
				{value: "", valid: false},
				{value: "1,2,3", valid: false},
				{value: "1,2,3,4,5", valid: true},
				{value: "1,2,3,4,invalid", valid: false},
				{value: "1,2,3,4,5,6", valid: false},
			},
		},
		{
			name: "Custom command keybinding",
			setup: func(config *UserConfig, value string) {
				config.CustomCommands = []CustomCommand{
					{
						Key:     value,
						Command: "echo 'hello'",
					},
				}
			},
			testCases: []testCase{
				{value: "", valid: true},
				{value: "<disabled>", valid: true},
				{value: "q", valid: true},
				{value: "<c-c>", valid: true},
				{value: "invalid_value", valid: false},
			},
		},
		{
			name: "Custom command sub menu",
			setup: func(config *UserConfig, _ string) {
				config.CustomCommands = []CustomCommand{
					{
						Key:         "X",
						Description: "My Custom Commands",
						CommandMenu: []CustomCommand{
							{Key: "1", Command: "echo 'hello'", Context: "global"},
						},
					},
				}
			},
			testCases: []testCase{
				{value: "", valid: true},
			},
		},
		{
			name: "Custom command sub menu",
			setup: func(config *UserConfig, _ string) {
				config.CustomCommands = []CustomCommand{
					{
						Key:     "X",
						Context: "global",
						CommandMenu: []CustomCommand{
							{Key: "1", Command: "echo 'hello'", Context: "global"},
						},
					},
				}
			},
			testCases: []testCase{
				{value: "", valid: false},
			},
		},
		{
			name: "Custom command sub menu",
			setup: func(config *UserConfig, _ string) {
				falseVal := false
				config.CustomCommands = []CustomCommand{
					{
						Key:        "X",
						Subprocess: &falseVal,
						CommandMenu: []CustomCommand{
							{Key: "1", Command: "echo 'hello'", Context: "global"},
						},
					},
				}
			},
			testCases: []testCase{
				{value: "", valid: false},
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
