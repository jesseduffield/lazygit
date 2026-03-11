package helpers

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/stretchr/testify/assert"
)

func TestGetPreferredRemoteName(t *testing.T) {
	tests := []struct {
		name                 string
		selectedRemoteBranch *models.RemoteBranch
		selectedRemote       *models.Remote
		expected             string
	}{
		{
			name:                 "uses selected remote branch when available",
			selectedRemoteBranch: &models.RemoteBranch{RemoteName: "upstream"},
			selectedRemote:       &models.Remote{Name: "origin"},
			expected:             "upstream",
		},
		{
			name:                 "falls back to selected remote",
			selectedRemoteBranch: nil,
			selectedRemote:       &models.Remote{Name: "mirror"},
			expected:             "mirror",
		},
		{
			name:                 "defaults to origin",
			selectedRemoteBranch: nil,
			selectedRemote:       nil,
			expected:             "origin",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, getPreferredRemoteName(test.selectedRemoteBranch, test.selectedRemote))
		})
	}
}
