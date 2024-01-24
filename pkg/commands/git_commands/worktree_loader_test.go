package git_commands

import (
	"testing"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGetWorktrees(t *testing.T) {
	type scenario struct {
		testName          string
		repoPaths         *RepoPaths
		before            func(runner *oscommands.FakeCmdObjRunner, fs afero.Fs, getRevParseArgs argFn)
		expectedWorktrees []*models.Worktree
		expectedErr       string
	}

	scenarios := []scenario{
		{
			testName: "Single worktree (main)",
			repoPaths: &RepoPaths{
				repoPath:     "/path/to/repo",
				worktreePath: "/path/to/repo",
			},
			before: func(runner *oscommands.FakeCmdObjRunner, fs afero.Fs, getRevParseArgs argFn) {
				runner.ExpectGitArgs([]string{"worktree", "list", "--porcelain"},
					`worktree /path/to/repo
HEAD d85cc9d281fa6ae1665c68365fc70e75e82a042d
branch refs/heads/mybranch
`,
					nil)

				gitArgsMainWorktree := append(append([]string{"-C", "/path/to/repo"}, getRevParseArgs()...), "--absolute-git-dir")
				runner.ExpectGitArgs(gitArgsMainWorktree, "/path/to/repo/.git", nil)
				_ = fs.MkdirAll("/path/to/repo/.git", 0o755)
			},
			expectedWorktrees: []*models.Worktree{
				{
					IsMain:        true,
					IsCurrent:     true,
					Path:          "/path/to/repo",
					IsPathMissing: false,
					GitDir:        "/path/to/repo/.git",
					Branch:        "mybranch",
					Name:          "repo",
				},
			},
			expectedErr: "",
		},
		{
			testName: "Multiple worktrees (main + linked)",
			repoPaths: &RepoPaths{
				repoPath:     "/path/to/repo",
				worktreePath: "/path/to/repo",
			},
			before: func(runner *oscommands.FakeCmdObjRunner, fs afero.Fs, getRevParseArgs argFn) {
				runner.ExpectGitArgs([]string{"worktree", "list", "--porcelain"},
					`worktree /path/to/repo
HEAD d85cc9d281fa6ae1665c68365fc70e75e82a042d
branch refs/heads/mybranch

worktree /path/to/repo-worktree
HEAD 775955775e79b8f5b4c4b56f82fbf657e2d5e4de
branch refs/heads/mybranch-worktree
`,
					nil)
				gitArgsMainWorktree := append(append([]string{"-C", "/path/to/repo"}, getRevParseArgs()...), "--absolute-git-dir")
				runner.ExpectGitArgs(gitArgsMainWorktree, "/path/to/repo/.git", nil)
				gitArgsLinkedWorktree := append(append([]string{"-C", "/path/to/repo-worktree"}, getRevParseArgs()...), "--absolute-git-dir")
				runner.ExpectGitArgs(gitArgsLinkedWorktree, "/path/to/repo/.git/worktrees/repo-worktree", nil)

				_ = fs.MkdirAll("/path/to/repo/.git", 0o755)
				_ = fs.MkdirAll("/path/to/repo-worktree", 0o755)
				_ = fs.MkdirAll("/path/to/repo/.git/worktrees/repo-worktree", 0o755)
				_ = afero.WriteFile(fs, "/path/to/repo-worktree/.git", []byte("gitdir: /path/to/repo/.git/worktrees/repo-worktree"), 0o755)
			},
			expectedWorktrees: []*models.Worktree{
				{
					IsMain:        true,
					IsCurrent:     true,
					Path:          "/path/to/repo",
					IsPathMissing: false,
					GitDir:        "/path/to/repo/.git",
					Branch:        "mybranch",
					Name:          "repo",
				},
				{
					IsMain:        false,
					IsCurrent:     false,
					Path:          "/path/to/repo-worktree",
					IsPathMissing: false,
					GitDir:        "/path/to/repo/.git/worktrees/repo-worktree",
					Branch:        "mybranch-worktree",
					Name:          "repo-worktree",
				},
			},
			expectedErr: "",
		},
		{
			testName: "Worktree missing path",
			repoPaths: &RepoPaths{
				repoPath:     "/path/to/repo",
				worktreePath: "/path/to/repo",
			},
			before: func(runner *oscommands.FakeCmdObjRunner, fs afero.Fs, getRevParseArgs argFn) {
				runner.ExpectGitArgs([]string{"worktree", "list", "--porcelain"},
					`worktree /path/to/worktree
HEAD 775955775e79b8f5b4c4b56f82fbf657e2d5e4de
branch refs/heads/missingbranch
`,
					nil)

				_ = fs.MkdirAll("/path/to/repo/.git", 0o755)
			},
			expectedWorktrees: []*models.Worktree{
				{
					IsMain:        false,
					IsCurrent:     false,
					Path:          "/path/to/worktree",
					IsPathMissing: true,
					GitDir:        "",
					Branch:        "missingbranch",
					Name:          "worktree",
				},
			},
			expectedErr: "",
		},
		{
			testName: "In linked worktree",
			repoPaths: &RepoPaths{
				repoPath:     "/path/to/repo",
				worktreePath: "/path/to/repo-worktree",
			},
			before: func(runner *oscommands.FakeCmdObjRunner, fs afero.Fs, getRevParseArgs argFn) {
				runner.ExpectGitArgs([]string{"worktree", "list", "--porcelain"},
					`worktree /path/to/repo
HEAD d85cc9d281fa6ae1665c68365fc70e75e82a042d
branch refs/heads/mybranch

worktree /path/to/repo-worktree
HEAD 775955775e79b8f5b4c4b56f82fbf657e2d5e4de
branch refs/heads/mybranch-worktree
`,
					nil)
				gitArgsMainWorktree := append(append([]string{"-C", "/path/to/repo"}, getRevParseArgs()...), "--absolute-git-dir")
				runner.ExpectGitArgs(gitArgsMainWorktree, "/path/to/repo/.git", nil)
				gitArgsLinkedWorktree := append(append([]string{"-C", "/path/to/repo-worktree"}, getRevParseArgs()...), "--absolute-git-dir")
				runner.ExpectGitArgs(gitArgsLinkedWorktree, "/path/to/repo/.git/worktrees/repo-worktree", nil)

				_ = fs.MkdirAll("/path/to/repo/.git", 0o755)
				_ = fs.MkdirAll("/path/to/repo-worktree", 0o755)
				_ = fs.MkdirAll("/path/to/repo/.git/worktrees/repo-worktree", 0o755)
				_ = afero.WriteFile(fs, "/path/to/repo-worktree/.git", []byte("gitdir: /path/to/repo/.git/worktrees/repo-worktree"), 0o755)
			},
			expectedWorktrees: []*models.Worktree{
				{
					IsMain:        false,
					IsCurrent:     true,
					Path:          "/path/to/repo-worktree",
					IsPathMissing: false,
					GitDir:        "/path/to/repo/.git/worktrees/repo-worktree",
					Branch:        "mybranch-worktree",
					Name:          "repo-worktree",
				},
				{
					IsMain:        true,
					IsCurrent:     false,
					Path:          "/path/to/repo",
					IsPathMissing: false,
					GitDir:        "/path/to/repo/.git",
					Branch:        "mybranch",
					Name:          "repo",
				},
			},
			expectedErr: "",
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			runner := oscommands.NewFakeRunner(t)
			fs := afero.NewMemMapFs()
			version, err := GetGitVersion(oscommands.NewDummyOSCommand())
			if err != nil {
				t.Fatal(err)
			}

			getRevParseArgs := func() []string {
				args := []string{"rev-parse"}
				if version.IsAtLeast(2, 31, 0) {
					args = append(args, "--path-format=absolute")
				}
				return args
			}

			s.before(runner, fs, getRevParseArgs)

			loader := &WorktreeLoader{
				GitCommon: buildGitCommon(commonDeps{runner: runner, fs: fs, repoPaths: s.repoPaths, gitVersion: version}),
			}

			worktrees, err := loader.GetWorktrees()
			if s.expectedErr != "" {
				assert.EqualError(t, errors.New(s.expectedErr), err.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, s.expectedWorktrees, worktrees)
			}
		})
	}
}

func TestGetUniqueNamesFromPaths(t *testing.T) {
	for _, scenario := range []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{},
			expected: []string{},
		},
		{
			input: []string{
				"/my/path/feature/one",
			},
			expected: []string{
				"one",
			},
		},
		{
			input: []string{
				"/my/path/feature/one/",
			},
			expected: []string{
				"one",
			},
		},
		{
			input: []string{
				"/a/b/c/d",
				"/a/b/c/e",
				"/a/b/f/d",
				"/a/e/c/d",
			},
			expected: []string{
				"b/c/d",
				"e",
				"f/d",
				"e/c/d",
			},
		},
	} {
		actual := getUniqueNamesFromPaths(scenario.input)
		assert.EqualValues(t, scenario.expected, actual)
	}
}
