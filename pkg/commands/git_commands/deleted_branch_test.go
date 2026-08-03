package git_commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/stretchr/testify/assert"
)

func TestObtainDeletedBranches(t *testing.T) {
	type scenario struct {
		testName           string
		entries            []*reflogEntry
		currentBranchNames []string
		expected           []*models.DeletedBranch
	}

	scenarios := []scenario{
		{
			testName: "recover deleted branch that was committed to",
			// newest-first reflog like `git log -g --format=+%H%x00%ct%x00%gs`
			entries: []*reflogEntry{
				{hash: "a", timestamp: 300, from: "feature", to: "main"},
				{hash: "b", timestamp: 200},
				{hash: "c", timestamp: 100, from: "main", to: "feature"},
				{hash: "d", timestamp: 50},
			},
			currentBranchNames: []string{"main"},
			expected: []*models.DeletedBranch{
				{Name: "feature", CommitHash: "b", DisplayName: "feature", Recency: "56y", UnixTimestamp: 200},
			},
		},
		{
			testName: "deleted branch left in favor of another branch, both tracked",
			// newest-first reflog
			entries: []*reflogEntry{
				{hash: "a", timestamp: 300, from: "feat/a", to: "main"}, // newest: leave feat/a
				{hash: "b", timestamp: 250},                            // commit on feat/a
				{hash: "a", timestamp: 200, from: "feat/b", to: "feat/a"},
				{hash: "c", timestamp: 150}, // commit on feat/b
				{hash: "a", timestamp: 100, from: "main", to: "feat/b"}, // oldest: create feat/b
			},
			currentBranchNames: []string{"main"},
			expected: []*models.DeletedBranch{
				{Name: "feat/a", CommitHash: "b", DisplayName: "feat/a", Recency: "56y", UnixTimestamp: 250},
				{Name: "feat/b", CommitHash: "c", DisplayName: "feat/b", Recency: "56y", UnixTimestamp: 150},
			},
		},
		{
			testName:           "existing branches are excluded",
			entries:            []*reflogEntry{
				{hash: "a", timestamp: 300, from: "main", to: "other"},
				{hash: "b", timestamp: 200},
				{hash: "c", timestamp: 100, from: "other", to: "main"},
			},
			currentBranchNames: []string{"main", "other"},
			expected:           nil,
		},
		{
			testName:           "no checkout entries means nothing recoverable",
			entries:            []*reflogEntry{
				{hash: "a", timestamp: 300},
				{hash: "b", timestamp: 200},
			},
			currentBranchNames: []string{"main"},
			expected:           nil,
		},
		{
			testName: "HEAD is not treated as a deleted branch",
			entries: []*reflogEntry{
				{hash: "a", timestamp: 300, from: "main", to: "HEAD"},
				{hash: "b", timestamp: 200},
			},
			currentBranchNames: []string{"main"},
			expected:           nil,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			result := obtainDeletedBranches(s.entries, s.currentBranchNames)
			assert.Equal(t, s.expected, result)
		})
	}
}

func TestParseReflogCheckoutSubject(t *testing.T) {
	type scenario struct {
		testName    string
		subject     string
		expected    []string
	}

	scenarios := []scenario{
		{
			testName: "normal checkout",
			subject:  "checkout: moving from feature to main",
			expected: []string{"feature", "main"},
		},
		{
			testName: "not a checkout",
			subject:  "commit: message",
			expected: []string{"", ""},
		},
		{
			testName: "checkout from detached head",
			subject:  "checkout: moving from HEAD to main",
			expected: []string{"HEAD", "main"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			from, to := parseReflogCheckoutSubject(s.subject)
			assert.Equal(t, s.expected[0], from)
			assert.Equal(t, s.expected[1], to)
		})
	}
}

func TestBranchRestoreBranch(t *testing.T) {
	type scenario struct {
		testName       string
		runner         *oscommands.FakeCmdObjRunner
		expectedErr    bool
		expectedUpstream string
	}

	scenarios := []scenario{
		{
			testName: "restore branch with no surviving remote branch",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"branch", "feature", "abc123"}, "", nil).
				ExpectGitArgs([]string{"for-each-ref", "--format=%(refname:short)", "refs/remotes"}, "", nil),
			expectedErr:       false,
			expectedUpstream: "",
		},
		{
			testName: "restore branch and reattach upstream when remote branch survives",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"branch", "feature", "abc123"}, "", nil).
				ExpectGitArgs([]string{"for-each-ref", "--format=%(refname:short)", "refs/remotes"}, "origin/feature\n", nil).
				ExpectGitArgs([]string{"branch", "--set-upstream-to=origin/feature", "feature"}, "", nil),
			expectedErr:       false,
			expectedUpstream: "origin/feature",
		},
		{
			testName: "restore branch when multiple remotes share the name does not set upstream",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"branch", "feature", "abc123"}, "", nil).
				ExpectGitArgs([]string{"for-each-ref", "--format=%(refname:short)", "refs/remotes"}, "origin/feature\nfork/feature\n", nil),
			expectedErr:       false,
			expectedUpstream: "",
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			instance := buildBranchCommands(commonDeps{runner: s.runner})

			upstream, err := instance.RestoreBranch("feature", "abc123")
			if s.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, s.expectedUpstream, upstream)
			s.runner.CheckForMissingCalls()
		})
	}
}
