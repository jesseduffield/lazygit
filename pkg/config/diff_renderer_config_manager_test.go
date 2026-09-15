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
