package config

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/stretchr/testify/assert"
)

func TestCurrentPagerName(t *testing.T) {
	tr := i18n.EnglishTranslationSet()

	scenarios := []struct {
		name     string
		pager    PagingConfig
		expected string
	}{
		{
			name:     "explicit name takes precedence over the command",
			pager:    PagingConfig{Name: "delta side-by-side", Pager: "delta --side-by-side"},
			expected: "delta side-by-side",
		},
		{
			name:     "derived from the first word of the pager command",
			pager:    PagingConfig{Pager: "delta --side-by-side"},
			expected: "delta",
		},
		{
			name:     "surrounding whitespace in the command is ignored",
			pager:    PagingConfig{Pager: "  diff-so-fancy  "},
			expected: "diff-so-fancy",
		},
		{
			name:     "falls back to the external diff command when there is no pager",
			pager:    PagingConfig{ExternalDiffCommand: "difft --color=always"},
			expected: "difft",
		},
		{
			name:     "no name can be derived for external diff",
			pager:    PagingConfig{UseExternalDiffGitConfig: true},
			expected: tr.ExternalDiffPagerName,
		},
		{
			name:     "no name can be derived for raw diff",
			pager:    PagingConfig{},
			expected: tr.DefaultPagerName,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			userConfig := &UserConfig{}
			userConfig.Git.Pagers = []PagingConfig{s.pager}
			config := NewPagerConfig(func() *UserConfig { return userConfig })

			assert.Equal(t, s.expected, config.CurrentPagerName(tr))
		})
	}
}

func TestCurrentPagerNameWithoutPagers(t *testing.T) {
	config := NewPagerConfig(func() *UserConfig { return &UserConfig{} })

	tr := i18n.EnglishTranslationSet()
	assert.Equal(t, tr.DefaultPagerName, config.CurrentPagerName(tr))
}

func TestCyclePagers(t *testing.T) {
	userConfig := &UserConfig{}
	userConfig.Git.Pagers = []PagingConfig{{Name: "a"}, {Name: "b"}, {Name: "c"}}
	config := NewPagerConfig(func() *UserConfig { return userConfig })

	currentIndex := func() int {
		index, _ := config.CurrentPagerIndex()
		return index
	}

	assert.Equal(t, 0, currentIndex())

	config.CyclePagers()
	assert.Equal(t, 1, currentIndex())
	config.CyclePagers()
	assert.Equal(t, 2, currentIndex())
	config.CyclePagers()
	assert.Equal(t, 0, currentIndex(), "cycling forward past the last pager wraps to the first")

	config.CyclePagersBackward()
	assert.Equal(t, 2, currentIndex(), "cycling backward past the first pager wraps to the last")
	config.CyclePagersBackward()
	assert.Equal(t, 1, currentIndex())
}
