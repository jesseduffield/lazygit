package helpers

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/stretchr/testify/assert"
)

func TestReadHeadInfo(t *testing.T) {
	scenarios := []struct {
		name string
		// The files to lay out below a temporary root directory, by their path
		// relative to it. "$root" in a file's content is replaced with the
		// root's path, so that a scenario can write an absolute path.
		files map[string]string
		// The repo to read, relative to the root. The directory is created
		// whether or not the scenario puts any files in it.
		repoPath   string
		expected   headInfo
		expectedOk bool
	}{
		{
			name:       "ordinary repo on a branch",
			files:      map[string]string{"repo/.git/HEAD": "ref: refs/heads/mybranch\n"},
			repoPath:   "repo",
			expected:   headInfo{branch: "mybranch"},
			expectedOk: true,
		},
		{
			name:       "ordinary repo at a detached head",
			files:      map[string]string{"repo/.git/HEAD": "d85cc9d2f5d0dc0b8f0e4d8e5b2ba0d1e7c8a3f6\n"},
			repoPath:   "repo",
			expected:   headInfo{hash: "d85cc9d2f5d0dc0b8f0e4d8e5b2ba0d1e7c8a3f6"},
			expectedOk: true,
		},
		{
			name: "worktree whose .git file names the git dir absolutely",
			files: map[string]string{
				"repo/.git/worktrees/wt/HEAD": "ref: refs/heads/mybranch\n",
				"wt/.git":                     "gitdir: $root/repo/.git/worktrees/wt\n",
			},
			repoPath:   "wt",
			expected:   headInfo{branch: "mybranch"},
			expectedOk: true,
		},
		{
			name: "worktree whose .git file names the git dir relatively",
			files: map[string]string{
				"repo/.git/worktrees/wt/HEAD": "ref: refs/heads/mybranch\n",
				"wt/.git":                     "gitdir: ../repo/.git/worktrees/wt\n",
			},
			repoPath:   "wt",
			expected:   headInfo{branch: "mybranch"},
			expectedOk: true,
		},
		{
			name: "submodule, whose .git file always names the git dir relatively",
			files: map[string]string{
				"repo/.git/modules/sub/HEAD": "ref: refs/heads/mybranch\n",
				"repo/sub/.git":              "gitdir: ../.git/modules/sub\n",
			},
			repoPath:   "repo/sub",
			expected:   headInfo{branch: "mybranch"},
			expectedOk: true,
		},
		{
			name: "repo that keeps its refs in a reftable, so HEAD holds a placeholder",
			files: map[string]string{
				"repo/.git/HEAD": "ref: refs/heads/.invalid\n",
			},
			repoPath:   "repo",
			expectedOk: false,
		},
		{
			name:       "directory without a .git entry",
			repoPath:   "notarepo",
			expectedOk: false,
		},
		{
			name:       ".git file that doesn't name a git dir",
			files:      map[string]string{"repo/.git": "not what git writes\n"},
			repoPath:   "repo",
			expectedOk: false,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			root := t.TempDir()
			for path, content := range s.files {
				fullPath := filepath.Join(root, filepath.FromSlash(path))
				assert.NoError(t, os.MkdirAll(filepath.Dir(fullPath), 0o700))
				content = strings.ReplaceAll(content, "$root", root)
				assert.NoError(t, os.WriteFile(fullPath, []byte(content), 0o600))
			}
			repoPath := filepath.Join(root, filepath.FromSlash(s.repoPath))
			assert.NoError(t, os.MkdirAll(repoPath, 0o700))

			head, ok := readHeadInfo(repoPath)

			assert.Equal(t, s.expectedOk, ok)
			assert.Equal(t, s.expected, head)
		})
	}
}

type osGuiCommon struct {
	types.IGuiCommon
	os *oscommands.OSCommand
}

func (self osGuiCommon) OS() *oscommands.OSCommand { return self.os }

func TestWithDirtyState(t *testing.T) {
	root := t.TempDir()
	run := func(dir string, args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		assert.NoError(t, err, string(out))
	}
	clean := filepath.Join(root, "clean")
	dirty := filepath.Join(root, "dirty")
	for _, dir := range []string{clean, dirty} {
		assert.NoError(t, os.Mkdir(dir, 0o755))
		run(dir, "init", "-q")
	}
	assert.NoError(t, os.WriteFile(filepath.Join(dirty, "new.txt"), []byte("x"), 0o644))

	input := []*models.Repo{{Path: clean, Name: "clean"}, {Path: dirty, Name: "dirty"}}

	result := newDirtyTestHelper().WithDirtyState(input, func() bool { return false })

	assert.Len(t, result, 2)
	assert.False(t, result[0].Dirty)
	assert.True(t, result[1].Dirty)
	assert.Equal(t, "dirty", result[1].Name)
	assert.False(t, input[1].Dirty)
	assert.NotSame(t, input[1], result[1])
}

func newDirtyTestHelper() *ReposHelper {
	return &ReposHelper{c: &HelperCommon{
		Common:     common.NewDummyCommon(),
		IGuiCommon: osGuiCommon{os: oscommands.NewDummyOSCommand()},
	}}
}

func TestWithDirtyStateNotARepo(t *testing.T) {
	dir := t.TempDir()
	result := newDirtyTestHelper().WithDirtyState(
		[]*models.Repo{{Path: dir, Name: "plain"}}, func() bool { return false })
	assert.Len(t, result, 1)
	assert.False(t, result[0].Dirty)
}

func TestWithDirtyStateStaleSkipsChecks(t *testing.T) {
	root := t.TempDir()
	cmd := exec.Command("git", "init", "-q")
	cmd.Dir = root
	assert.NoError(t, cmd.Run())
	assert.NoError(t, os.WriteFile(filepath.Join(root, "new.txt"), []byte("x"), 0o644))

	result := newDirtyTestHelper().WithDirtyState(
		[]*models.Repo{{Path: root}}, func() bool { return true })
	assert.False(t, result[0].Dirty)
}

func TestLoadRepo(t *testing.T) {
	root := t.TempDir()
	cmd := exec.Command("git", "init", "-q", "-b", "topic")
	cmd.Dir = root
	assert.NoError(t, cmd.Run())
	assert.NoError(t, os.WriteFile(filepath.Join(root, "new.txt"), []byte("x"), 0o644))

	repo := newDirtyTestHelper().LoadRepo(root)

	assert.Equal(t, &models.Repo{Path: root, Branch: "topic", Dirty: true}, repo)
}

func TestWithRepoRowUpdated(t *testing.T) {
	alpha := &models.Repo{Path: filepath.FromSlash("/w/alpha"), Name: "alpha", Branch: "main"}
	beta := &models.Repo{Path: filepath.FromSlash("/w/beta"), Name: "beta", Branch: "main"}
	repos := []*models.Repo{alpha, beta}

	t.Run("updates the row at the path", func(t *testing.T) {
		result := withRepoRowUpdated(repos, &models.Repo{Path: "/w/beta/", Branch: "topic", Dirty: true})

		assert.Same(t, alpha, result[0])
		assert.Equal(t, &models.Repo{Path: beta.Path, Name: "beta", Branch: "topic", Dirty: true}, result[1])
		assert.Equal(t, "main", beta.Branch, "input must not change")
		assert.False(t, beta.Dirty, "input must not change")
	})

	t.Run("keeps the rows if none is at the path", func(t *testing.T) {
		result := withRepoRowUpdated(repos, &models.Repo{Path: "/w/gamma", Branch: "topic"})

		assert.Equal(t, repos, result)
	})
}
