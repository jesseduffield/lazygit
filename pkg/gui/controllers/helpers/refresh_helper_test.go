package helpers

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/hosting_service"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/context/traits"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/stefanhaller/git-todo-parser/todo"
	"github.com/stretchr/testify/assert"
)

func TestCaptureLocalCommitSelectionRange(t *testing.T) {
	testCases := []struct {
		name          string
		commits       []*models.Commit
		selectedIdx   int
		rangeStartIdx int
		expected      *localCommitSelectionRange
	}{
		{
			name:          "captures selected commit and range start",
			commits:       makeCommits("a", "b"),
			selectedIdx:   1,
			rangeStartIdx: 0,
			expected: &localCommitSelectionRange{
				selectedHash:   "b",
				rangeStartHash: "a",
				selectedIdx:    1,
				rangeStartIdx:  0,
				mode:           traits.RangeSelectModeSticky,
			},
		},
		{
			name:          "ignores invalid range start index",
			commits:       makeCommits("a"),
			selectedIdx:   0,
			rangeStartIdx: 1,
			expected:      nil,
		},
		{
			name:          "ignores empty selected hash",
			commits:       append(makeCommits("a"), makeTodoCommit(todo.UpdateRef)),
			selectedIdx:   1,
			rangeStartIdx: 0,
			expected:      nil,
		},
		{
			name:          "ignores empty range start hash",
			commits:       append(makeCommits("a"), makeTodoCommit(todo.Exec)),
			selectedIdx:   0,
			rangeStartIdx: 1,
			expected:      nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			selectionRange := captureLocalCommitSelectionRange(
				testCase.commits,
				testCase.selectedIdx,
				testCase.rangeStartIdx,
				traits.RangeSelectModeSticky,
			)

			assert.Equal(t, testCase.expected, selectionRange)
		})
	}
}

func TestFindLocalCommitSelectionRange(t *testing.T) {
	type expectation struct {
		selectedIdx   int
		rangeStartIdx int
		moved         bool
		found         bool
	}

	selectionRange := localCommitSelectionRange{
		selectedHash:   "b",
		rangeStartHash: "c",
		selectedIdx:    1,
		rangeStartIdx:  2,
		mode:           traits.RangeSelectModeSticky,
	}

	testCases := []struct {
		name     string
		commits  []*models.Commit
		expected expectation
	}{
		{
			name:    "finds selection after commits are inserted above it",
			commits: makeCommits("new", "a", "b", "c"),
			expected: expectation{
				selectedIdx:   2,
				rangeStartIdx: 3,
				moved:         true,
				found:         true,
			},
		},
		{
			name:    "finds selection that did not move",
			commits: makeCommits("a", "b", "c"),
			expected: expectation{
				selectedIdx:   1,
				rangeStartIdx: 2,
				found:         true,
			},
		},
		{
			name:     "reports not found when a hash is missing",
			commits:  makeCommits("a", "b"),
			expected: expectation{},
		},
		{
			name: "skips todo entries with the same hash as a selected commit",
			commits: []*models.Commit{
				makeTodoCommitWithHash("b", todo.Revert),
				makeCommits("a")[0],
				makeCommits("b")[0],
				makeCommits("c")[0],
			},
			expected: expectation{
				selectedIdx:   2,
				rangeStartIdx: 3,
				moved:         true,
				found:         true,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			selectedIdx, rangeStartIdx, moved, found := findLocalCommitSelectionRange(testCase.commits, &selectionRange)
			actual := expectation{
				selectedIdx:   selectedIdx,
				rangeStartIdx: rangeStartIdx,
				moved:         moved,
				found:         found,
			}

			assert.Equal(t, testCase.expected, actual)
		})
	}
}

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
				assert.Equal(t, c.expected, result.remote.Name)
			}
		})
	}
}

func TestGetAuthenticatedGithubRemotes(t *testing.T) {
	githubRemotes := []githubRemoteInfo{
		makeGithubRemoteInfo("origin", "github.com"),
		makeGithubRemoteInfo("fork", "github.com"),
		makeGithubRemoteInfo("enterprise", "ghe.example.com"),
		makeGithubRemoteInfo("missing-auth", "no-token.example.com"),
	}

	callsByHost := map[string]int{}
	result := getAuthenticatedGithubRemotes(githubRemotes, func(host string) string {
		callsByHost[host]++
		switch host {
		case "github.com":
			return "github-token"
		case "ghe.example.com":
			return "ghe-token"
		default:
			return ""
		}
	})

	assert.Equal(t, []githubRemoteInfo{
		makeAuthenticatedGithubRemoteInfo("origin", "github.com", "github-token"),
		makeAuthenticatedGithubRemoteInfo("fork", "github.com", "github-token"),
		makeAuthenticatedGithubRemoteInfo("enterprise", "ghe.example.com", "ghe-token"),
	}, result)
	// Two remotes share github.com; the lookup runs only once.
	assert.Equal(t, map[string]int{
		"github.com":           1,
		"ghe.example.com":      1,
		"no-token.example.com": 1,
	}, callsByHost)
}

func makeGithubRemoteInfoList(names ...string) []githubRemoteInfo {
	return lo.Map(names, func(name string, _ int) githubRemoteInfo {
		return makeGithubRemoteInfo(name, name)
	})
}

func makeGithubRemoteInfo(name string, webDomain string) githubRemoteInfo {
	return githubRemoteInfo{
		remote: &models.Remote{Name: name},
		serviceInfo: hosting_service.ServiceInfo{
			RepoName:  name,
			WebDomain: webDomain,
		},
	}
}

func makeAuthenticatedGithubRemoteInfo(name string, webDomain string, authToken string) githubRemoteInfo {
	info := makeGithubRemoteInfo(name, webDomain)
	info.authToken = authToken
	return info
}

func makeCommits(hashes ...string) []*models.Commit {
	hashPool := &utils.StringPool{}
	return lo.Map(hashes, func(hash string, _ int) *models.Commit {
		return models.NewCommit(hashPool, models.NewCommitOpts{Hash: hash})
	})
}

func makeTodoCommit(action todo.TodoCommand) *models.Commit {
	return models.NewCommit(&utils.StringPool{}, models.NewCommitOpts{Action: action})
}

func makeTodoCommitWithHash(hash string, action todo.TodoCommand) *models.Commit {
	return models.NewCommit(&utils.StringPool{}, models.NewCommitOpts{Hash: hash, Action: action})
}
