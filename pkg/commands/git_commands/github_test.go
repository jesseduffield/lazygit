package git_commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/hosting_service"
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

func TestGetRepoInfoFromURL(t *testing.T) {
	cases := []struct {
		name     string
		url      string
		expected hosting_service.RepoInformation
	}{
		{
			name: "SSH URL",
			url:  "git@github.com:jesseduffield/lazygit.git",
			expected: hosting_service.RepoInformation{
				Owner:      "jesseduffield",
				Repository: "lazygit",
			},
		},
		{
			name: "HTTPS URL",
			url:  "https://github.com/jesseduffield/lazygit.git",
			expected: hosting_service.RepoInformation{
				Owner:      "jesseduffield",
				Repository: "lazygit",
			},
		},
		{
			name: "HTTPS URL without .git",
			url:  "https://github.com/jesseduffield/lazygit",
			expected: hosting_service.RepoInformation{
				Owner:      "jesseduffield",
				Repository: "lazygit",
			},
		},
		{
			name: "SSH URL with org nesting",
			url:  "git@github.com:my-org/sub-group/lazygit.git",
			expected: hosting_service.RepoInformation{
				Owner:      "my-org/sub-group",
				Repository: "lazygit",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, err := hosting_service.GetRepoInfoFromURL(c.url)
			assert.NoError(t, err)
			assert.Equal(t, c.expected, result)
		})
	}
}

func TestGenerateGithubPullRequestMap(t *testing.T) {
	cases := []struct {
		name     string
		prs      []*models.GithubPullRequest
		branches []*models.Branch
		remotes  []*models.Remote
		expected map[string]*models.GithubPullRequest
	}{
		{
			name:     "empty inputs",
			prs:      []*models.GithubPullRequest{},
			branches: []*models.Branch{},
			remotes:  []*models.Remote{},
			expected: map[string]*models.GithubPullRequest{},
		},
		{
			name: "matches PR to branch tracking origin",
			prs: []*models.GithubPullRequest{
				{
					HeadRefName:         "feature-branch",
					Number:              42,
					Title:               "Add feature",
					State:               "OPEN",
					Url:                 "https://github.com/jesseduffield/lazygit/pull/42",
					HeadRepositoryOwner: models.GithubRepositoryOwner{Login: "jesseduffield"},
				},
			},
			branches: []*models.Branch{
				{
					Name:           "feature-branch",
					UpstreamRemote: "origin",
					UpstreamBranch: "feature-branch",
				},
			},
			remotes: []*models.Remote{
				{
					Name: "origin",
					Urls: []string{"git@github.com:jesseduffield/lazygit.git"},
				},
			},
			expected: map[string]*models.GithubPullRequest{
				"feature-branch": {
					HeadRefName:         "feature-branch",
					Number:              42,
					Title:               "Add feature",
					State:               "OPEN",
					Url:                 "https://github.com/jesseduffield/lazygit/pull/42",
					HeadRepositoryOwner: models.GithubRepositoryOwner{Login: "jesseduffield"},
				},
			},
		},
		{
			name: "does not match branch without upstream",
			prs: []*models.GithubPullRequest{
				{
					HeadRefName:         "feature-branch",
					Number:              42,
					Title:               "Add feature",
					State:               "OPEN",
					HeadRepositoryOwner: models.GithubRepositoryOwner{Login: "jesseduffield"},
				},
			},
			branches: []*models.Branch{
				{
					Name: "feature-branch",
					// no upstream set
				},
			},
			remotes: []*models.Remote{
				{
					Name: "origin",
					Urls: []string{"git@github.com:jesseduffield/lazygit.git"},
				},
			},
			expected: map[string]*models.GithubPullRequest{},
		},
		{
			name: "matches fork PR to branch tracking fork remote",
			prs: []*models.GithubPullRequest{
				{
					HeadRefName:         "fix-bug",
					Number:              99,
					Title:               "Fix bug",
					State:               "OPEN",
					HeadRepositoryOwner: models.GithubRepositoryOwner{Login: "contributor"},
				},
			},
			branches: []*models.Branch{
				{
					Name:           "fix-bug",
					UpstreamRemote: "contributor",
					UpstreamBranch: "fix-bug",
				},
			},
			remotes: []*models.Remote{
				{
					Name: "origin",
					Urls: []string{"git@github.com:jesseduffield/lazygit.git"},
				},
				{
					Name: "contributor",
					Urls: []string{"git@github.com:contributor/lazygit.git"},
				},
			},
			expected: map[string]*models.GithubPullRequest{
				"fix-bug": {
					HeadRefName:         "fix-bug",
					Number:              99,
					Title:               "Fix bug",
					State:               "OPEN",
					HeadRepositoryOwner: models.GithubRepositoryOwner{Login: "contributor"},
				},
			},
		},
		{
			name: "does not match when owner differs",
			prs: []*models.GithubPullRequest{
				{
					HeadRefName:         "feature-branch",
					Number:              42,
					Title:               "Add feature",
					State:               "OPEN",
					HeadRepositoryOwner: models.GithubRepositoryOwner{Login: "someone-else"},
				},
			},
			branches: []*models.Branch{
				{
					Name:           "feature-branch",
					UpstreamRemote: "origin",
					UpstreamBranch: "feature-branch",
				},
			},
			remotes: []*models.Remote{
				{
					Name: "origin",
					Urls: []string{"git@github.com:jesseduffield/lazygit.git"},
				},
			},
			expected: map[string]*models.GithubPullRequest{},
		},
		{
			name: "matches when UpstreamRemote is a full URL",
			prs: []*models.GithubPullRequest{
				{
					HeadRefName:         "my-branch",
					Number:              55,
					Title:               "Full URL upstream",
					State:               "OPEN",
					HeadRepositoryOwner: models.GithubRepositoryOwner{Login: "contributor"},
				},
			},
			branches: []*models.Branch{
				{
					Name:           "my-branch",
					UpstreamRemote: "git@github.com:contributor/lazygit.git",
					UpstreamBranch: "my-branch",
				},
			},
			remotes: []*models.Remote{
				{
					Name: "origin",
					Urls: []string{"git@github.com:jesseduffield/lazygit.git"},
				},
			},
			expected: map[string]*models.GithubPullRequest{
				"my-branch": {
					HeadRefName:         "my-branch",
					Number:              55,
					Title:               "Full URL upstream",
					State:               "OPEN",
					HeadRepositoryOwner: models.GithubRepositoryOwner{Login: "contributor"},
				},
			},
		},
		{
			name: "matches with HTTPS remote URL",
			prs: []*models.GithubPullRequest{
				{
					HeadRefName:         "my-pr",
					Number:              10,
					Title:               "My PR",
					State:               "MERGED",
					HeadRepositoryOwner: models.GithubRepositoryOwner{Login: "jesseduffield"},
				},
			},
			branches: []*models.Branch{
				{
					Name:           "my-pr",
					UpstreamRemote: "origin",
					UpstreamBranch: "my-pr",
				},
			},
			remotes: []*models.Remote{
				{
					Name: "origin",
					Urls: []string{"https://github.com/jesseduffield/lazygit.git"},
				},
			},
			expected: map[string]*models.GithubPullRequest{
				"my-pr": {
					HeadRefName:         "my-pr",
					Number:              10,
					Title:               "My PR",
					State:               "MERGED",
					HeadRepositoryOwner: models.GithubRepositoryOwner{Login: "jesseduffield"},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := GenerateGithubPullRequestMap(c.prs, c.branches, c.remotes)
			assert.Equal(t, c.expected, result)
		})
	}
}

func mkRemoteList(names ...string) []*models.Remote {
	return lo.Map(names, func(name string, _ int) *models.Remote {
		return &models.Remote{Name: name}
	})
}
