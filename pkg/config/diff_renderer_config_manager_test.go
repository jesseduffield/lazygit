package config

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/stretchr/testify/assert"
)

func TestCurrentDiffRendererName(t *testing.T) {
	tr := i18n.EnglishTranslationSet()

	scenarios := []struct {
		name               string
		diffRendererConfig DiffRendererConfig
		expected           string
	}{
		{
			name:               "explicit name takes precedence over the command",
			diffRendererConfig: DiffRendererConfig{Name: "delta side-by-side", Command: "delta --side-by-side"},
			expected:           "delta side-by-side",
		},
		{
			name:               "derived from the first word of the stdinFilter command",
			diffRendererConfig: DiffRendererConfig{Command: "delta --side-by-side"},
			expected:           "delta",
		},
		{
			name:               "surrounding whitespace in the command is ignored",
			diffRendererConfig: DiffRendererConfig{Command: "  diff-so-fancy  "},
			expected:           "diff-so-fancy",
		},
		{
			name:               "derived from the first word of the extDiff command",
			diffRendererConfig: DiffRendererConfig{Type: "extDiff", Command: "difft --color=always"},
			expected:           "difft",
		},
		{
			name:               "no name can be derived for external diff",
			diffRendererConfig: DiffRendererConfig{Type: "extDiff"},
			expected:           tr.ExternalDiffDiffRendererName,
		},
		{
			name:               "derived from first argument of rawGit args",
			diffRendererConfig: DiffRendererConfig{Type: "rawGit", Args: []string{"--color-words"}},
			expected:           "--color-words",
		},
		{
			name:               "no name can be derived for raw diff",
			diffRendererConfig: DiffRendererConfig{Type: "rawGit"},
			expected:           tr.DefaultDiffRendererName,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			userConfig := &UserConfig{}
			userConfig.Git.DiffRenderers = []DiffRendererConfig{s.diffRendererConfig}
			config := NewDiffRendererConfigManager(func() *UserConfig { return userConfig })

			assert.Equal(t, s.expected, config.CurrentDiffRendererName(tr))
		})
	}
}

func TestGetStdinFilterCommand(t *testing.T) {
	scenarios := []struct {
		name               string
		diffRendererConfig DiffRendererConfig
		width              int
		expected           string
	}{
		{
			name:               "a command without template variables is passed through",
			diffRendererConfig: DiffRendererConfig{Command: "delta --paging=never"},
			width:              120,
			expected:           "delta --paging=never",
		},
		{
			name:               "the width the diff is rendered at",
			diffRendererConfig: DiffRendererConfig{Command: "delta --width={{width}}"},
			width:              120,
			expected:           "delta --width=120",
		},
		{
			name:               "the width of one side of a side-by-side rendering",
			diffRendererConfig: DiffRendererConfig{Command: "ydiff -p cat -w {{columnWidth}}"},
			width:              120,
			expected:           "ydiff -p cat -w 54",
		},
		{
			name:               "a template variable can also be written with a leading dot",
			diffRendererConfig: DiffRendererConfig{Command: "delta --width={{.width}}"},
			width:              120,
			expected:           "delta --width=120",
		},
		{
			name:               "nothing is returned for a renderer of another type",
			diffRendererConfig: DiffRendererConfig{Type: "extDiff", Command: "difft --width={{width}}"},
			width:              120,
			expected:           "",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			userConfig := &UserConfig{}
			userConfig.Git.DiffRenderers = []DiffRendererConfig{s.diffRendererConfig}
			config := NewDiffRendererConfigManager(func() *UserConfig { return userConfig })

			assert.Equal(t, s.expected, config.GetStdinFilterCommand(s.width))
		})
	}
}

func TestGetExternalDiffCommand(t *testing.T) {
	scenarios := []struct {
		name               string
		diffRendererConfig DiffRendererConfig
		expected           string
	}{
		{
			name:               "a command without template variables is passed through",
			diffRendererConfig: DiffRendererConfig{Type: "extDiff", Command: "difft --color=always"},
			expected:           "difft --color=always",
		},
		{
			name:               "the width the diff is rendered at",
			diffRendererConfig: DiffRendererConfig{Type: "extDiff", Command: "difft --width={{width}}"},
			expected:           "difft --width=120",
		},
		{
			name:               "the width alongside the diff context size",
			diffRendererConfig: DiffRendererConfig{Type: "extDiff", Command: "difft --width={{width}} --context={{diffContext}}"},
			expected:           "difft --width=120 --context=3",
		},
		{
			name:               "nothing is returned for a renderer of another type",
			diffRendererConfig: DiffRendererConfig{Command: "delta --width={{width}}"},
			expected:           "",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			userConfig := &UserConfig{}
			userConfig.Git.DiffRenderers = []DiffRendererConfig{s.diffRendererConfig}
			config := NewDiffRendererConfigManager(func() *UserConfig { return userConfig })

			assert.Equal(t, s.expected, config.GetExternalDiffCommand(3, 120))
		})
	}
}

func TestCurrentDiffRendererNameWithoutDiffRenderers(t *testing.T) {
	config := NewDiffRendererConfigManager(func() *UserConfig { return &UserConfig{} })

	tr := i18n.EnglishTranslationSet()
	assert.Equal(t, tr.DefaultDiffRendererName, config.CurrentDiffRendererName(tr))
}

func TestCycleDiffRenderers(t *testing.T) {
	userConfig := &UserConfig{}
	userConfig.Git.DiffRenderers = []DiffRendererConfig{{Name: "a"}, {Name: "b"}, {Name: "c"}}
	config := NewDiffRendererConfigManager(func() *UserConfig { return userConfig })

	currentIndex := func() int {
		index, _ := config.CurrentDiffRendererIndex()
		return index
	}

	assert.Equal(t, 0, currentIndex())

	config.CycleDiffRenderers()
	assert.Equal(t, 1, currentIndex())
	config.CycleDiffRenderers()
	assert.Equal(t, 2, currentIndex())
	config.CycleDiffRenderers()
	assert.Equal(t, 0, currentIndex(), "cycling forward past the last diff renderer wraps to the first")

	config.CycleDiffRenderersBackward()
	assert.Equal(t, 2, currentIndex(), "cycling backward past the first diff renderer wraps to the last")
	config.CycleDiffRenderersBackward()
	assert.Equal(t, 1, currentIndex())
}
