package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEditTemplate(t *testing.T) {
	trueVal := true

	scenarios := []struct {
		name                              string
		osConfig                          *OSConfig
		guessDefaultEditor                func() string
		expectedEditTemplate              string
		expectedEditAtLineTemplate        string
		expectedEditAtLineAndWaitTemplate string
		expectedSuspend                   bool
	}{
		{
			"Default template is vim",
			&OSConfig{},
			func() string { return "" },
			"vim -- {{filename}}",
			"vim +{{line}} -- {{filename}}",
			"vim +{{line}} -- {{filename}}",
			true,
		},
		{
			"Setting a preset",
			&OSConfig{
				EditPreset: "vscode",
			},
			func() string { return "" },
			"code --reuse-window -- {{filename}}",
			"code --reuse-window --goto -- {{filename}}:{{line}}",
			"code --reuse-window --goto --wait -- {{filename}}:{{line}}",
			false,
		},
		{
			"Setting a preset wins over guessed editor",
			&OSConfig{
				EditPreset: "vscode",
			},
			func() string { return "nano" },
			"code --reuse-window -- {{filename}}",
			"code --reuse-window --goto -- {{filename}}:{{line}}",
			"code --reuse-window --goto --wait -- {{filename}}:{{line}}",
			false,
		},
		{
			"Overriding a preset with explicit config (edit)",
			&OSConfig{
				EditPreset:    "vscode",
				Edit:          "myeditor {{filename}}",
				SuspendOnEdit: &trueVal,
			},
			func() string { return "" },
			"myeditor {{filename}}",
			"code --reuse-window --goto -- {{filename}}:{{line}}",
			"code --reuse-window --goto --wait -- {{filename}}:{{line}}",
			true,
		},
		{
			"Overriding a preset with explicit config (edit at line)",
			&OSConfig{
				EditPreset:    "vscode",
				EditAtLine:    "myeditor --line={{line}} {{filename}}",
				SuspendOnEdit: &trueVal,
			},
			func() string { return "" },
			"code --reuse-window -- {{filename}}",
			"myeditor --line={{line}} {{filename}}",
			"code --reuse-window --goto --wait -- {{filename}}:{{line}}",
			true,
		},
		{
			"Overriding a preset with explicit config (edit at line and wait)",
			&OSConfig{
				EditPreset:        "vscode",
				EditAtLineAndWait: "myeditor --line={{line}} -w {{filename}}",
				SuspendOnEdit:     &trueVal,
			},
			func() string { return "" },
			"code --reuse-window -- {{filename}}",
			"code --reuse-window --goto -- {{filename}}:{{line}}",
			"myeditor --line={{line}} -w {{filename}}",
			true,
		},
		{
			"Unknown preset name",
			&OSConfig{
				EditPreset: "thisPresetDoesNotExist",
			},
			func() string { return "" },
			"vim -- {{filename}}",
			"vim +{{line}} -- {{filename}}",
			"vim +{{line}} -- {{filename}}",
			true,
		},
		{
			"Guessing a preset from guessed editor",
			&OSConfig{},
			func() string { return "emacs" },
			"emacs -- {{filename}}",
			"emacs +{{line}} -- {{filename}}",
			"emacs +{{line}} -- {{filename}}",
			true,
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			template, suspend := GetEditTemplate(s.osConfig, s.guessDefaultEditor)
			assert.Equal(t, s.expectedEditTemplate, template)
			assert.Equal(t, s.expectedSuspend, suspend)

			template, suspend = GetEditAtLineTemplate(s.osConfig, s.guessDefaultEditor)
			assert.Equal(t, s.expectedEditAtLineTemplate, template)
			assert.Equal(t, s.expectedSuspend, suspend)

			template = GetEditAtLineAndWaitTemplate(s.osConfig, s.guessDefaultEditor)
			assert.Equal(t, s.expectedEditAtLineAndWaitTemplate, template)
		})
	}
}
