package helpers

import (
	"testing"

	"github.com/jesseduffield/generics/slices"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
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
	return slices.Map(names, func(name string) *models.Remote {
		return &models.Remote{Name: name}
	})
}
