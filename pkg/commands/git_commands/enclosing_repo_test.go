package git_commands

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/stretchr/testify/assert"
)

func TestIsTrackedByEnclosingRepo(t *testing.T) {
	scenarios := []struct {
		testName string
		// Sets up a repo (or not) in the fixture's dir, and returns the path to
		// ask about
		setup       func(f repoFixture) string
		expected    bool
		expectError bool
	}{
		{
			testName: "Not in a repo",
			setup: func(f repoFixture) string {
				return f.writeFile(".lazygit.yml")
			},
			expected: false,
		},
		{
			testName: "Untracked file",
			setup: func(f repoFixture) string {
				f.git("init")
				return f.writeFile(".lazygit.yml")
			},
			expected: false,
		},
		{
			testName: "Tracked file",
			setup: func(f repoFixture) string {
				f.git("init")
				path := f.writeFile(".lazygit.yml")
				f.git("add", ".lazygit.yml")
				return path
			},
			expected: true,
		},
		{
			testName: "Tracked file in a subdirectory",
			setup: func(f repoFixture) string {
				f.git("init")
				path := f.writeFile("sub/dir/.lazygit.yml")
				f.git("add", "sub/dir/.lazygit.yml")
				return path
			},
			expected: true,
		},
		{
			testName: "Tracked file whose name differs in case",
			setup: func(f repoFixture) string {
				f.git("init")
				f.writeFile(".LAZYGIT.YML")
				f.git("add", ".LAZYGIT.YML")
				path := filepath.Join(f.dir, ".lazygit.yml")
				if _, err := os.Lstat(path); err != nil {
					f.t.Skip("the file system is case-sensitive")
				}
				return path
			},
			expected: true,
		},
		{
			testName: "Tracked file in a directory whose name the index has in a different Unicode normalization",
			setup: func(f repoFixture) string {
				f.git("init")
				// "Ü" as "U" plus a combining diaeresis (NFD); with precomposeunicode
				// turned off, git add keeps this spelling in the index
				f.writeFile("Ü/.lazygit.yml")
				f.git("-c", "core.precomposeunicode=false", "add", ".")
				// "Ü" as a single code point (NFC)
				path := filepath.Join(f.dir, "Ü", ".lazygit.yml")
				if _, err := os.Lstat(path); err != nil {
					f.t.Skip("the file system treats the two spellings as different names")
				}
				return path
			},
			expected: true,
		},
		{
			testName: "Tracked file with glob characters in its name",
			setup: func(f repoFixture) string {
				f.git("init")
				path := f.writeFile("[x].yml")
				f.git("add", "[x].yml")
				return path
			},
			expected: true,
		},
		{
			testName: "Name that isn't ASCII",
			setup: func(f repoFixture) string {
				f.git("init")
				return f.writeFile("ü.yml")
			},
			expectError: true,
		},
		{
			testName: "Tracked file with GIT_LITERAL_PATHSPECS set",
			setup: func(f repoFixture) string {
				f.git("init")
				path := f.writeFile(".lazygit.yml")
				f.git("add", ".lazygit.yml")
				f.t.Setenv("GIT_LITERAL_PATHSPECS", "1")
				return path
			},
			expected: true,
		},
		{
			testName: "Tracked file with GIT_DIR pointing at another repo",
			setup: func(f repoFixture) string {
				f.git("init", "other")
				f.git("init", "repo")
				path := f.writeFile("repo/.lazygit.yml")
				f.git("-C", "repo", "add", ".lazygit.yml")
				f.t.Setenv("GIT_DIR", filepath.Join(f.dir, "other", ".git"))
				return path
			},
			expected: true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			// Keep the user's git config (e.g. a global gitignore) out of the test
			globalConfig := repoFixture{t: t, dir: t.TempDir()}.writeFile("gitconfig")
			t.Setenv("GIT_CONFIG_GLOBAL", globalConfig)
			t.Setenv("GIT_CONFIG_NOSYSTEM", "1")

			path := s.setup(repoFixture{t: t, dir: t.TempDir()})

			isTracked, err := IsTrackedByEnclosingRepo(oscommands.NewDummyOSCommand().Cmd, path)
			if s.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, s.expected, isTracked)
		})
	}
}

type repoFixture struct {
	t   *testing.T
	dir string
}

func (f repoFixture) git(args ...string) {
	f.t.Helper()
	output, err := exec.Command("git", append([]string{"-C", f.dir}, args...)...).CombinedOutput()
	if !assert.NoError(f.t, err, string(output)) {
		f.t.FailNow()
	}
}

// Creates an empty file at relPath below the fixture's dir, and returns its
// absolute path
func (f repoFixture) writeFile(relPath string) string {
	f.t.Helper()
	path := filepath.Join(f.dir, relPath)
	if !assert.NoError(f.t, os.MkdirAll(filepath.Dir(path), 0o755)) ||
		!assert.NoError(f.t, os.WriteFile(path, []byte{}, 0o644)) {
		f.t.FailNow()
	}
	return path
}
