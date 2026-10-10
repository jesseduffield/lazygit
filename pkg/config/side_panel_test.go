package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSidePanelsForMultiRepo(t *testing.T) {
	scenarios := []struct {
		name     string
		panels   []SidePanel
		expected []SidePanel
	}{
		{
			name:     "replaces status",
			panels:   []SidePanel{{"status"}, {"files", "worktrees"}},
			expected: []SidePanel{{"repos"}, {"files", "worktrees"}},
		},
		{
			name:     "replaces status inside a group",
			panels:   []SidePanel{{"files", "status"}},
			expected: []SidePanel{{"files", "repos"}},
		},
		{
			name:     "keeps the config if it lists repos already",
			panels:   []SidePanel{{"status"}, {"repos", "files"}},
			expected: []SidePanel{{"status"}, {"repos", "files"}},
		},
		{
			name:     "adds repos as the first panel if the config has no status",
			panels:   []SidePanel{{"files"}},
			expected: []SidePanel{{"repos"}, {"files"}},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			original := cloneSidePanels(s.panels)
			assert.Equal(t, s.expected, SidePanelsForMultiRepo(s.panels))
			assert.Equal(t, original, s.panels, "input must not change")
		})
	}
}

func cloneSidePanels(panels []SidePanel) []SidePanel {
	result := make([]SidePanel, len(panels))
	for i, panel := range panels {
		result[i] = append(SidePanel{}, panel...)
	}
	return result
}
