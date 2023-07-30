package helpers

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestGetSuggestedRemote(t *testing.T) {
	cases := []struct {
		remotes  []*models.Remote
		expected string
	}{
		{mkRemoteList(), "origin"},
		{mkRemoteList("upstream", "origin", "foo"), "origin"},
		{mkRemoteList("upstream", "foo", "bar"), "upstream"},
	}

	for _, c := range cases {
		result := getSuggestedRemote(c.remotes)
		assert.EqualValues(t, c.expected, result)
	}
}

func mkRemoteList(names ...string) []*models.Remote {
	return lo.Map(names, func(name string, _ int) *models.Remote {
		return &models.Remote{Name: name}
	})
}
