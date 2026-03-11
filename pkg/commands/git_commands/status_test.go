package git_commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusCommandsIsInRebase(t *testing.T) {
	type scenario struct {
		testName string
		prepare  func(t *testing.T, gitDir string)
		expected bool
	}

	scenarios := []scenario{
		{
			testName: "returns true when rebase-merge/head-name exists",
			prepare: func(t *testing.T, gitDir string) {
				assert.NoError(t, writeFile(filepath.Join(gitDir, "rebase-merge", "head-name"), "refs/heads/main\n"))
			},
			expected: true,
		},
		{
			testName: "returns true when rebase-apply/head-name exists",
			prepare: func(t *testing.T, gitDir string) {
				assert.NoError(t, writeFile(filepath.Join(gitDir, "rebase-apply", "head-name"), "refs/heads/main\n"))
			},
			expected: true,
		},
		{
			testName: "returns true when REBASE_HEAD exists",
			prepare: func(t *testing.T, gitDir string) {
				assert.NoError(t, writeFile(filepath.Join(gitDir, "REBASE_HEAD"), "abc123\n"))
			},
			expected: true,
		},
		{
			testName: "returns false when rebase directory is stale without markers",
			prepare: func(t *testing.T, gitDir string) {
				assert.NoError(t, createDir(filepath.Join(gitDir, "rebase-merge")))
			},
			expected: false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			repoPath := t.TempDir()
			gitDir := filepath.Join(repoPath, ".git")
			if !assert.NoError(t, createDir(gitDir)) {
				return
			}
			s.prepare(t, gitDir)

			status := NewStatusCommands(buildGitCommon(commonDeps{repoPaths: MockRepoPaths(repoPath)}))
			actual, err := status.IsInRebase()
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, s.expected, actual)
		})
	}
}

func writeFile(path string, contents string) error {
	if err := createDir(filepath.Dir(path)); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(contents), 0o644)
}

func createDir(path string) error {
	return os.MkdirAll(path, 0o755)
}
