package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrentPagerName(t *testing.T) {
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
			name:     "no name can be derived",
			pager:    PagingConfig{UseExternalDiffGitConfig: true},
			expected: "",
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			userConfig := &UserConfig{}
			userConfig.Git.Pagers = []PagingConfig{s.pager}
			config := NewPagerConfig(func() *UserConfig { return userConfig })

			assert.Equal(t, s.expected, config.CurrentPagerName())
		})
	}
}

func TestCurrentPagerNameWithoutPagers(t *testing.T) {
	config := NewPagerConfig(func() *UserConfig { return &UserConfig{} })

	assert.Equal(t, "", config.CurrentPagerName())
}
