package git_commands

import (
	"errors"
	"strings"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/stretchr/testify/assert"
)

func TestObtainDeletedBranches(t *testing.T) {
	type scenario struct {
		testName     string
		entries      []*reflogEntry
		existingRefs []string
		expected     []*models.DeletedBranch
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
			existingRefs: []string{"main", "HEAD"},
			expected: []*models.DeletedBranch{
				{Name: "feature", CommitHash: "b", DisplayName: "feature", Recency: "56y", UnixTimestamp: 200},
			},
		},
		{
			testName: "deleted branch left in favor of another branch, both tracked",
			// newest-first reflog
			entries: []*reflogEntry{
				{hash: "a", timestamp: 300, from: "feat/a", to: "main"}, // newest: leave feat/a
				{hash: "b", timestamp: 250},                             // commit on feat/a
				{hash: "a", timestamp: 200, from: "feat/b", to: "feat/a"},
				{hash: "c", timestamp: 150},                             // commit on feat/b
				{hash: "a", timestamp: 100, from: "main", to: "feat/b"}, // oldest: create feat/b
			},
			existingRefs: []string{"main", "HEAD"},
			expected: []*models.DeletedBranch{
				{Name: "feat/a", CommitHash: "b", DisplayName: "feat/a", Recency: "56y", UnixTimestamp: 250},
				{Name: "feat/b", CommitHash: "c", DisplayName: "feat/b", Recency: "56y", UnixTimestamp: 150},
			},
		},
		{
			testName: "existing branches are excluded",
			entries: []*reflogEntry{
				{hash: "a", timestamp: 300, from: "main", to: "other"},
				{hash: "b", timestamp: 200},
				{hash: "c", timestamp: 100, from: "other", to: "main"},
			},
			existingRefs: []string{"main", "other", "HEAD"},
			expected:     nil,
		},
		{
			testName: "no checkout entries means nothing recoverable",
			entries: []*reflogEntry{
				{hash: "a", timestamp: 300},
				{hash: "b", timestamp: 200},
			},
			existingRefs: []string{"main", "HEAD"},
			expected:     nil,
		},
		{
			testName: "HEAD is not treated as a deleted branch",
			entries: []*reflogEntry{
				{hash: "a", timestamp: 300, from: "main", to: "HEAD"},
				{hash: "b", timestamp: 200},
			},
			existingRefs: []string{"main", "HEAD"},
			expected:     nil,
		},
		{
			testName: "branch created via checkout with no commits is recoverable",
			// newest-first reflog: create feature, then immediately leave it
			entries: []*reflogEntry{
				{hash: "a", timestamp: 300, from: "feature", to: "main"}, // leave feature
				{hash: "b", timestamp: 200, from: "main", to: "feature"}, // create feature
			},
			existingRefs: []string{"main", "HEAD"},
			expected: []*models.DeletedBranch{
				{Name: "feature", CommitHash: "b", DisplayName: "feature", Recency: "56y", UnixTimestamp: 200},
			},
		},
		{
			testName: "detached head checkout is not treated as a deleted branch",
			entries: []*reflogEntry{
				{hash: "a", timestamp: 300, from: "abc1234567890abcdefabcdefabcdefabcdefabcdef", to: "main"},
				{hash: "b", timestamp: 200, from: "main", to: "abc1234567890abcdefabcdefabcdefabcdefabcdef"},
			},
			existingRefs: []string{"main", "HEAD"},
			expected:     nil,
		},
		{
			testName: "existing remote-tracking ref is excluded from candidates",
			entries: []*reflogEntry{
				{hash: "a", timestamp: 300, from: "origin/feature", to: "main"},
				{hash: "b", timestamp: 200, from: "main", to: "origin/feature"},
			},
			existingRefs: []string{"main", "origin/feature", "HEAD"},
			expected:     nil,
		},
		{
			testName: "existing tag checkout is not treated as a deleted branch",
			entries: []*reflogEntry{
				{hash: "a", timestamp: 300, from: "v1.0.0", to: "main"},
				{hash: "b", timestamp: 200, from: "main", to: "v1.0.0"},
			},
			existingRefs: []string{"main", "v1.0.0", "HEAD"},
			expected:     nil,
		},
		{
			testName: "checkout to a commit expression is not treated as a deleted branch",
			entries: []*reflogEntry{
				{hash: "a", timestamp: 300, from: "HEAD~1", to: "main"},
				{hash: "b", timestamp: 200, from: "main", to: "HEAD~1"},
			},
			existingRefs: []string{"main", "HEAD"},
			expected:     nil,
		},
		{
			testName: "checkout to a sha-256 hash is not treated as a deleted branch",
			entries: []*reflogEntry{
				{hash: "a", timestamp: 300, from: "8f0f1f2f3f4f5f6f7f8f9fafbfcfdfefff0f1f2f3f4f5f6f7f8f9fafbfcfdfeff", to: "main"},
				{hash: "b", timestamp: 200, from: "main", to: "8f0f1f2f3f4f5f6f7f8f9fafbfcfdfefff0f1f2f3f4f5f6f7f8f9fafbfcfdfeff"},
			},
			existingRefs: []string{"main", "HEAD"},
			expected:     nil,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			result := obtainDeletedBranches(s.entries, s.existingRefs, isValidRefFormatStub)
			assert.Equal(t, s.expected, result)
		})
	}
}

// isValidRefFormatStub is a stand-in for `git check-ref-format` used to keep
// the pure tests independent of the subprocess. It applies the same ref-name
// grammar rules git enforces.
func isValidRefFormatStub(name string) bool {
	if name == "" {
		return false
	}
	if strings.HasPrefix(name, "-") {
		return false
	}
	if strings.HasSuffix(name, ".") || strings.HasSuffix(name, "/") {
		return false
	}
	if strings.Contains(name, "..") || strings.Contains(name, "@{") {
		return false
	}
	for _, c := range name {
		if c <= ' ' || strings.ContainsRune("~^:?*[\\", c) {
			return false
		}
	}
	return true
}

func TestIsValidRefFormat(t *testing.T) {
	type scenario struct {
		testName string
		name     string
		isValid  bool
	}

	scenarios := []scenario{
		{
			testName: "valid branch name",
			name:     "feature/foo",
			isValid:  true,
		},
		{
			// a sha-looking string is still a valid ref to git; filtering it is
			// done separately by looksLikeSha
			testName: "sha-like is a valid ref to git",
			name:     "8f0f1f2f3f4f",
			isValid:  true,
		},
		{
			testName: "rejects commit expression",
			name:     "HEAD~1",
		},
		{
			testName: "rejects trailing slash",
			name:     "feature/",
		},
		{
			testName: "rejects double dot",
			name:     "feature..other",
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			err := errors.New("invalid ref")
			if s.isValid {
				err = nil
			}
			runner := oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"check-ref-format", "--allow-onelevel", s.name}, "", err)

			gitCommon := buildGitCommon(commonDeps{runner: runner})
			loader := &BranchLoader{
				Common:    gitCommon.Common,
				GitCommon: gitCommon,
				cmd:       gitCommon.cmd,
			}

			assert.Equal(t, s.isValid, loader.isValidRefFormat(s.name))
			runner.CheckForMissingCalls()
		})
	}
}

func TestParseReflogCheckoutSubject(t *testing.T) {
	type scenario struct {
		testName string
		subject  string
		expected []string
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
		testName         string
		runner           *oscommands.FakeCmdObjRunner
		expectedErr      bool
		expectedUpstream string
	}

	scenarios := []scenario{
		{
			testName: "restore branch with no surviving remote branch",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"branch", "feature", "abc123"}, "", nil).
				ExpectGitArgs([]string{"for-each-ref", "--format=%(refname:short)", "refs/remotes"}, "", nil),
			expectedErr:      false,
			expectedUpstream: "",
		},
		{
			testName: "restore branch and reattach upstream when remote branch survives",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"branch", "feature", "abc123"}, "", nil).
				ExpectGitArgs([]string{"for-each-ref", "--format=%(refname:short)", "refs/remotes"}, "origin/feature\n", nil).
				ExpectGitArgs([]string{"branch", "--set-upstream-to=origin/feature", "feature"}, "", nil),
			expectedErr:      false,
			expectedUpstream: "origin/feature",
		},
		{
			testName: "restore branch when multiple remotes share the name does not set upstream",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"branch", "feature", "abc123"}, "", nil).
				ExpectGitArgs([]string{"for-each-ref", "--format=%(refname:short)", "refs/remotes"}, "origin/feature\nfork/feature\n", nil),
			expectedErr:      false,
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
