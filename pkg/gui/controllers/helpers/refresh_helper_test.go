package helpers

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestGetGithubBaseRemote(t *testing.T) {
	cases := []struct {
		name             string
		githubRemotes    []githubRemoteInfo
		configuredRemote string
		expected         string
	}{
		{
			name:             "configured remote wins",
			githubRemotes:    makeGithubRemoteInfoList("origin", "upstream", "fork"),
			configuredRemote: "fork",
			expected:         "fork",
		},
		{
			name:             "configured remote not in github remotes returns nil",
			githubRemotes:    makeGithubRemoteInfoList("origin"),
			configuredRemote: "missing",
			expected:         "",
		},
		{
			name:             "single github remote is auto-picked",
			githubRemotes:    makeGithubRemoteInfoList("myremote"),
			configuredRemote: "",
			expected:         "myremote",
		},
		{
			name:             "upstream is preferred when multiple github remotes exist",
			githubRemotes:    makeGithubRemoteInfoList("origin", "upstream", "fork"),
			configuredRemote: "",
			expected:         "upstream",
		},
		{
			name:             "no upstream and multiple remotes returns nil",
			githubRemotes:    makeGithubRemoteInfoList("origin", "fork"),
			configuredRemote: "",
			expected:         "",
		},
		{
			name:             "empty list returns nil",
			githubRemotes:    nil,
			configuredRemote: "",
			expected:         "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := getGithubBaseRemote(c.githubRemotes, c.configuredRemote)
			if c.expected == "" {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, c.expected, result.Name)
			}
		})
	}
}

func makeGithubRemoteInfoList(names ...string) []githubRemoteInfo {
	return lo.Map(names, func(name string, _ int) githubRemoteInfo {
		return githubRemoteInfo{remote: &models.Remote{Name: name}, repoName: name}
	})
}
