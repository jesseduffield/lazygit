package git_commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestIsValidGhVersion(t *testing.T) {
	type scenario struct {
		versionStr     string
		expectedResult bool
	}

	scenarios := []scenario{
		{
			"",
			false,
		},
		{
			`gh version 1.0.0 (2020-08-23)
			https://github.com/cli/cli/releases/tag/v1.0.0`,
			false,
		},
		{
			`gh version 2.0.0 (2021-08-23)
			https://github.com/cli/cli/releases/tag/v2.0.0`,
			true,
		},
		{
			`gh version 1.1.0 (2021-10-14)
			https://github.com/cli/cli/releases/tag/v1.1.0

			A new release of gh is available: 1.1.0 → v2.2.0
			To upgrade, run: brew update && brew upgrade gh
			https://github.com/cli/cli/releases/tag/v2.2.0`,
			false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.versionStr, func(t *testing.T) {
			result := isGhVersionValid(s.versionStr)
			assert.Equal(t, result, s.expectedResult)
		})
	}
}

func TestGetSuggestedRemoteName(t *testing.T) {
	cases := []struct {
		remotes  []*models.Remote
		expected string
	}{
		{mkRemoteList(), "origin"},
		{mkRemoteList("upstream", "origin", "foo"), "origin"},
		{mkRemoteList("upstream", "foo", "bar"), "upstream"},
	}

	for _, c := range cases {
		result := GetSuggestedRemoteName(c.remotes)
		assert.EqualValues(t, c.expected, result)
	}
}

func mkRemoteList(names ...string) []*models.Remote {
	return lo.Map(names, func(name string, _ int) *models.Remote {
		return &models.Remote{Name: name}
	})
}
