package git_commands

import (
	"testing"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestWorkingTreeStageFile(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"add", "--", "test.txt"}, "", nil)

	instance := buildWorkingTreeCommands(commonDeps{runner: runner})

	assert.NoError(t, instance.StageFile("test.txt"))
	runner.CheckForMissingCalls()
}

func TestWorkingTreeStageFiles(t *testing.T) {
	runner := oscommands.NewFakeRunner(t).
		ExpectGitArgs([]string{"add", "--", "test.txt", "test2.txt"}, "", nil)

	instance := buildWorkingTreeCommands(commonDeps{runner: runner})

	assert.NoError(t, instance.StageFiles([]string{"test.txt", "test2.txt"}, nil))
	runner.CheckForMissingCalls()
}

func TestWorkingTreeUnstageFile(t *testing.T) {
	type scenario struct {
		testName string
		reset    bool
		runner   *oscommands.FakeCmdObjRunner
		test     func(error)
	}

	scenarios := []scenario{
		{
			testName: "Remove an untracked file from staging",
			reset:    false,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"rm", "--cached", "--force", "--", "test.txt"}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testName: "Remove a tracked file from staging",
			reset:    true,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"reset", "HEAD", "--", "test.txt"}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			instance := buildWorkingTreeCommands(commonDeps{runner: s.runner})
			s.test(instance.UnStageFile([]string{"test.txt"}, s.reset))
		})
	}
}

// these tests don't cover everything, in part because we already have an integration
// test which does cover everything. I don't want to unnecessarily assert on the 'how'
// when the 'what' is what matters
func TestWorkingTreeDiscardAllFileChanges(t *testing.T) {
	type scenario struct {
		testName             string
		file                 *models.File
		removedFileErr       error
		runner               *oscommands.FakeCmdObjRunner
		expectedError        string
		expectedRemovedFiles []string
	}

	scenarios := []scenario{
		{
			testName: "An error occurred when resetting",
			file: &models.File{
				Path:             "test",
				HasStagedChanges: true,
			},
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"reset", "--", "test"}, "", errors.New("error")),
			expectedError: "error",
		},
		{
			testName: "An error occurred when removing file",
			file: &models.File{
				Path:    "test",
				Tracked: false,
				Added:   true,
			},
			removedFileErr:       errors.New("an error occurred when removing file"),
			runner:               oscommands.NewFakeRunner(t),
			expectedError:        "an error occurred when removing file",
			expectedRemovedFiles: []string{"test"},
		},
		{
			testName: "An error occurred with checkout",
			file: &models.File{
				Path:             "test",
				Tracked:          true,
				HasStagedChanges: false,
			},
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"checkout", "--", "test"}, "", errors.New("error")),
			expectedError: "error",
		},
		{
			testName: "Checkout only",
			file: &models.File{
				Path:             "test",
				Tracked:          true,
				HasStagedChanges: false,
			},
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"checkout", "--", "test"}, "", nil),
		},
		{
			testName: "Reset and checkout staged changes",
			file: &models.File{
				Path:             "test",
				Tracked:          true,
				HasStagedChanges: true,
			},
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"reset", "--", "test"}, "", nil).
				ExpectGitArgs([]string{"checkout", "--", "test"}, "", nil),
		},
		{
			testName: "Reset and checkout merge conflicts",
			file: &models.File{
				Path:              "test",
				Tracked:           true,
				HasMergeConflicts: true,
			},
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"reset", "--", "test"}, "", nil).
				ExpectGitArgs([]string{"checkout", "--", "test"}, "", nil),
		},
		{
			testName: "Reset and remove",
			file: &models.File{
				Path:             "test",
				Tracked:          false,
				Added:            true,
				HasStagedChanges: true,
			},
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"reset", "--", "test"}, "", nil),
			expectedRemovedFiles: []string{"test"},
		},
		{
			testName: "Remove only",
			file: &models.File{
				Path:             "test",
				Tracked:          false,
				Added:            true,
				HasStagedChanges: false,
			},
			runner:               oscommands.NewFakeRunner(t),
			expectedRemovedFiles: []string{"test"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			var removedFiles []string
			removeFile := func(path string) error {
				removedFiles = append(removedFiles, path)
				return s.removedFileErr
			}
			instance := buildWorkingTreeCommands(commonDeps{runner: s.runner, removeFile: removeFile})
			err := instance.DiscardAllFileChanges(s.file)

			if s.expectedError == "" {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, s.expectedError, err.Error())
			}
			assert.Equal(t, s.expectedRemovedFiles, removedFiles)
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestWorkingTreeDiff(t *testing.T) {
	type scenario struct {
		testName            string
		file                *models.File
		plain               bool
		cached              bool
		ignoreWhitespace    bool
		contextSize         uint64
		similarityThreshold int
		runner              *oscommands.FakeCmdObjRunner
	}

	const expectedResult = "pretend this is an actual git diff"

	scenarios := []scenario{
		{
			testName: "Default case",
			file: &models.File{
				Path:             "test.txt",
				HasStagedChanges: false,
				Tracked:          true,
			},
			plain:               false,
			cached:              false,
			ignoreWhitespace:    false,
			contextSize:         3,
			similarityThreshold: 50,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"-C", "/path/to/worktree", "diff", "--no-ext-diff", "--submodule", "--unified=3", "--color=always", "--find-renames=50%", "--", "test.txt"}, expectedResult, nil),
		},
		{
			testName: "cached",
			file: &models.File{
				Path:             "test.txt",
				HasStagedChanges: false,
				Tracked:          true,
			},
			plain:               false,
			cached:              true,
			ignoreWhitespace:    false,
			contextSize:         3,
			similarityThreshold: 50,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"-C", "/path/to/worktree", "diff", "--no-ext-diff", "--submodule", "--unified=3", "--color=always", "--find-renames=50%", "--cached", "--", "test.txt"}, expectedResult, nil),
		},
		{
			testName: "plain",
			file: &models.File{
				Path:             "test.txt",
				HasStagedChanges: false,
				Tracked:          true,
			},
			plain:               true,
			cached:              false,
			ignoreWhitespace:    false,
			contextSize:         3,
			similarityThreshold: 50,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"-C", "/path/to/worktree", "diff", "--no-ext-diff", "--submodule", "--unified=3", "--color=never", "--find-renames=50%", "--", "test.txt"}, expectedResult, nil),
		},
		{
			testName: "File not tracked and file has no staged changes",
			file: &models.File{
				Path:             "test.txt",
				HasStagedChanges: false,
				Tracked:          false,
			},
			plain:               false,
			cached:              false,
			ignoreWhitespace:    false,
			contextSize:         3,
			similarityThreshold: 50,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"-C", "/path/to/worktree", "diff", "--no-ext-diff", "--submodule", "--unified=3", "--color=always", "--find-renames=50%", "--no-index", "--", "/dev/null", "test.txt"}, expectedResult, nil),
		},
		{
			testName: "Default case (ignore whitespace)",
			file: &models.File{
				Path:             "test.txt",
				HasStagedChanges: false,
				Tracked:          true,
			},
			plain:               false,
			cached:              false,
			ignoreWhitespace:    true,
			contextSize:         3,
			similarityThreshold: 50,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"-C", "/path/to/worktree", "diff", "--no-ext-diff", "--submodule", "--unified=3", "--color=always", "--ignore-all-space", "--find-renames=50%", "--", "test.txt"}, expectedResult, nil),
		},
		{
			testName: "Show diff with custom context size",
			file: &models.File{
				Path:             "test.txt",
				HasStagedChanges: false,
				Tracked:          true,
			},
			plain:               false,
			cached:              false,
			ignoreWhitespace:    false,
			contextSize:         17,
			similarityThreshold: 50,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"-C", "/path/to/worktree", "diff", "--no-ext-diff", "--submodule", "--unified=17", "--color=always", "--find-renames=50%", "--", "test.txt"}, expectedResult, nil),
		},
		{
			testName: "Show diff with custom similarity threshold",
			file: &models.File{
				Path:             "test.txt",
				HasStagedChanges: false,
				Tracked:          true,
			},
			plain:               false,
			cached:              false,
			ignoreWhitespace:    false,
			contextSize:         3,
			similarityThreshold: 33,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"-C", "/path/to/worktree", "diff", "--no-ext-diff", "--submodule", "--unified=3", "--color=always", "--find-renames=33%", "--", "test.txt"}, expectedResult, nil),
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			userConfig := config.GetDefaultConfig()
			userConfig.Git.IgnoreWhitespaceInDiffView = s.ignoreWhitespace
			userConfig.Git.DiffContextSize = s.contextSize
			userConfig.Git.RenameSimilarityThreshold = s.similarityThreshold
			repoPaths := RepoPaths{
				worktreePath: "/path/to/worktree",
			}

			instance := buildWorkingTreeCommands(commonDeps{runner: s.runner, userConfig: userConfig, appState: &config.AppState{}, repoPaths: &repoPaths})
			result := instance.WorktreeFileDiff(s.file, s.plain, s.cached)
			assert.Equal(t, expectedResult, result)
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestWorkingTreeShowFileDiff(t *testing.T) {
	type scenario struct {
		testName         string
		from             string
		to               string
		reverse          bool
		plain            bool
		ignoreWhitespace bool
		contextSize      uint64
		runner           *oscommands.FakeCmdObjRunner
	}

	const expectedResult = "pretend this is an actual git diff"

	scenarios := []scenario{
		{
			testName:         "Default case",
			from:             "1234567890",
			to:               "0987654321",
			reverse:          false,
			plain:            false,
			ignoreWhitespace: false,
			contextSize:      3,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"-C", "/path/to/worktree", "-c", "diff.noprefix=false", "diff", "--no-ext-diff", "--submodule", "--unified=3", "--no-renames", "--color=always", "1234567890", "0987654321", "--", "test.txt"}, expectedResult, nil),
		},
		{
			testName:         "Show diff with custom context size",
			from:             "1234567890",
			to:               "0987654321",
			reverse:          false,
			plain:            false,
			ignoreWhitespace: false,
			contextSize:      123,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"-C", "/path/to/worktree", "-c", "diff.noprefix=false", "diff", "--no-ext-diff", "--submodule", "--unified=123", "--no-renames", "--color=always", "1234567890", "0987654321", "--", "test.txt"}, expectedResult, nil),
		},
		{
			testName:         "Default case (ignore whitespace)",
			from:             "1234567890",
			to:               "0987654321",
			reverse:          false,
			plain:            false,
			ignoreWhitespace: true,
			contextSize:      3,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"-C", "/path/to/worktree", "-c", "diff.noprefix=false", "diff", "--no-ext-diff", "--submodule", "--unified=3", "--no-renames", "--color=always", "1234567890", "0987654321", "--ignore-all-space", "--", "test.txt"}, expectedResult, nil),
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			userConfig := config.GetDefaultConfig()
			userConfig.Git.IgnoreWhitespaceInDiffView = s.ignoreWhitespace
			userConfig.Git.DiffContextSize = s.contextSize
			repoPaths := RepoPaths{
				worktreePath: "/path/to/worktree",
			}

			instance := buildWorkingTreeCommands(commonDeps{runner: s.runner, userConfig: userConfig, appState: &config.AppState{}, repoPaths: &repoPaths})

			result, err := instance.ShowFileDiff(s.from, s.to, s.reverse, "test.txt", s.plain)
			assert.NoError(t, err)
			assert.Equal(t, expectedResult, result)
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestWorkingTreeCheckoutFile(t *testing.T) {
	type scenario struct {
		testName   string
		commitHash string
		fileName   string
		runner     *oscommands.FakeCmdObjRunner
		test       func(error)
	}

	scenarios := []scenario{
		{
			testName:   "typical case",
			commitHash: "11af912",
			fileName:   "test999.txt",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"checkout", "11af912", "--", "test999.txt"}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			testName:   "returns error if there is one",
			commitHash: "11af912",
			fileName:   "test999.txt",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"checkout", "11af912", "--", "test999.txt"}, "", errors.New("error")),
			test: func(err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			instance := buildWorkingTreeCommands(commonDeps{runner: s.runner})

			s.test(instance.CheckoutFile(s.commitHash, s.fileName))
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestWorkingTreeDiscardUnstagedFileChanges(t *testing.T) {
	type scenario struct {
		testName string
		file     *models.File
		runner   *oscommands.FakeCmdObjRunner
		test     func(error)
	}

	scenarios := []scenario{
		{
			testName: "valid case",
			file:     &models.File{Path: "test.txt"},
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"checkout", "--", "test.txt"}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			instance := buildWorkingTreeCommands(commonDeps{runner: s.runner})
			s.test(instance.DiscardUnstagedFileChanges(s.file))
			s.runner.CheckForMissingCalls()
		})
	}
}

// testNode implements IFileNode for unit tests.
type testNode struct {
	children []*testNode
	path     string
	file     *models.File // non-nil only for file nodes
}

func (n *testNode) ForEachFile(cb func(*models.File) error) error {
	if n.file != nil {
		return cb(n.file)
	}
	for _, child := range n.children {
		if err := child.ForEachFile(cb); err != nil {
			return err
		}
	}
	return nil
}

func (n *testNode) GetFilePathsMatching(test func(*models.File) bool) []string {
	if n.file != nil {
		if test(n.file) {
			return []string{n.path}
		}
		return nil
	}
	return lo.FlatMap(n.children, func(child *testNode, _ int) []string {
		return child.GetFilePathsMatching(test)
	})
}

func (n *testNode) GetPath() string       { return n.path }
func (n *testNode) GetFile() *models.File { return n.file }

func TestWorkingTreeDiscardAllDirChanges(t *testing.T) {
	type scenario struct {
		testName               string
		nodes                  []IFileNode
		runner                 *oscommands.FakeCmdObjRunner
		dirsWithRemainingFiles []string // dirs where isDirEmpty returns false
		expectedRemovedFiles   []string
		expectedRemovedDirs    []string
	}

	scenarios := []scenario{
		{
			testName: "multiple regular tracked files batched into a single checkout call",
			nodes: []IFileNode{&testNode{
				children: []*testNode{
					{path: "a.txt", file: &models.File{Path: "a.txt", Tracked: true}},
					{path: "b.txt", file: &models.File{Path: "b.txt", Tracked: true}},
					{path: "c.txt", file: &models.File{Path: "c.txt", Tracked: true}},
				},
			}},
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"checkout", "--", "a.txt", "b.txt", "c.txt"}, "", nil),
		},
		{
			testName: "staged files batched into a single reset then a single checkout",
			nodes: []IFileNode{&testNode{
				children: []*testNode{
					{path: "a.txt", file: &models.File{Path: "a.txt", Tracked: true, HasStagedChanges: true}},
					{path: "b.txt", file: &models.File{Path: "b.txt", Tracked: true, HasStagedChanges: true}},
				},
			}},
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"reset", "--", "a.txt", "b.txt"}, "", nil).
				ExpectGitArgs([]string{"checkout", "--", "a.txt", "b.txt"}, "", nil),
		},
		{
			testName: "added files with no staged changes are removed from disk without any git call",
			nodes: []IFileNode{&testNode{
				children: []*testNode{
					{path: "new1.txt", file: &models.File{Path: "new1.txt", Added: true}},
					{path: "new2.txt", file: &models.File{Path: "new2.txt", Added: true}},
				},
			}},
			runner:               oscommands.NewFakeRunner(t),
			expectedRemovedFiles: []string{"new1.txt", "new2.txt"},
		},
		{
			testName: "files from multiple nodes are batched into a single git call",
			nodes: []IFileNode{
				&testNode{
					path: "dir1",
					children: []*testNode{
						{path: "dir1/a.txt", file: &models.File{Path: "dir1/a.txt", Tracked: true}},
						{path: "dir1/b.txt", file: &models.File{Path: "dir1/b.txt", Added: true}},
					},
				},
				&testNode{
					path: "dir2",
					children: []*testNode{
						{path: "dir2/c.txt", file: &models.File{Path: "dir2/c.txt", Tracked: true}},
						{path: "dir2/d.txt", file: &models.File{Path: "dir2/d.txt", Added: true}},
					},
				},
			},
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"checkout", "--", "dir1/a.txt", "dir2/c.txt"}, "", nil),
			dirsWithRemainingFiles: []string{"dir1", "dir2"}, // tracked files a.txt / c.txt remain
			expectedRemovedFiles:   []string{"dir1/b.txt", "dir2/d.txt"},
		},
		{
			testName: "empty parent directory is removed after all its added files are deleted",
			nodes: []IFileNode{&testNode{
				path: "dir",
				children: []*testNode{
					{
						path: "dir/newdir",
						children: []*testNode{
							{path: "dir/newdir/a.txt", file: &models.File{Path: "dir/newdir/a.txt", Added: true}},
							{path: "dir/newdir/b.txt", file: &models.File{Path: "dir/newdir/b.txt", Added: true}},
						},
					},
				},
			}},
			runner:                 oscommands.NewFakeRunner(t),
			dirsWithRemainingFiles: []string{"dir"}, // assume there are other tracked files in dir
			expectedRemovedFiles:   []string{"dir/newdir/a.txt", "dir/newdir/b.txt"},
			expectedRemovedDirs:    []string{"dir/newdir"},
		},
		{
			testName: "nested empty directories are removed bottom-up",
			nodes: []IFileNode{&testNode{
				path: "newdir",
				children: []*testNode{
					{
						path: "newdir/sub",
						children: []*testNode{
							{path: "newdir/sub/file.txt", file: &models.File{Path: "newdir/sub/file.txt", Added: true}},
						},
					},
				},
			}},
			runner:               oscommands.NewFakeRunner(t),
			expectedRemovedFiles: []string{"newdir/sub/file.txt"},
			expectedRemovedDirs:  []string{"newdir/sub", "newdir"},
		},
		{
			testName: "empty directory is NOT removed when individual file nodes are selected",
			nodes: []IFileNode{
				&testNode{path: "newdir/a.txt", file: &models.File{Path: "newdir/a.txt", Added: true}},
				&testNode{path: "newdir/b.txt", file: &models.File{Path: "newdir/b.txt", Added: true}},
			},
			runner:               oscommands.NewFakeRunner(t),
			expectedRemovedFiles: []string{"newdir/a.txt", "newdir/b.txt"},
			// newdir becomes empty but was not selected as a directory node, so it is not removed
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			var removedFiles []string
			removeFile := func(path string) error {
				removedFiles = append(removedFiles, path)
				return nil
			}
			isDirEmpty := func(path string) (bool, error) { return !lo.Contains(s.dirsWithRemainingFiles, path), nil }
			var removedDirs []string
			removeDir := func(path string) error {
				removedDirs = append(removedDirs, path)
				return nil
			}
			instance := buildWorkingTreeCommands(commonDeps{
				runner:     s.runner,
				removeFile: removeFile,
				isDirEmpty: isDirEmpty,
				removeDir:  removeDir,
			})
			err := instance.DiscardAllDirChanges(s.nodes)
			assert.NoError(t, err)
			assert.Equal(t, s.expectedRemovedFiles, removedFiles)
			assert.Equal(t, s.expectedRemovedDirs, removedDirs)
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestWorkingTreeDiscardUnstagedDirChanges(t *testing.T) {
	type scenario struct {
		testName               string
		nodes                  []IFileNode
		runner                 *oscommands.FakeCmdObjRunner
		dirsWithRemainingFiles []string // dirs where isDirEmpty returns false
		expectedRemovedFiles   []string
		expectedRemovedDirs    []string
	}

	scenarios := []scenario{
		{
			testName: "directory node: removes untracked files and checks out tracked files by path, not by directory",
			nodes: []IFileNode{&testNode{
				path: "dir",
				children: []*testNode{
					{path: "dir/tracked1.txt", file: &models.File{Path: "dir/tracked1.txt", Tracked: true}},
					{path: "dir/tracked2.txt", file: &models.File{Path: "dir/tracked2.txt", Tracked: true}},
					{path: "dir/new.txt", file: &models.File{Path: "dir/new.txt", Tracked: false}},
				},
			}},
			// Must checkout the individual files, not "dir" — otherwise a filter would be ignored.
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"checkout", "--", "dir/tracked1.txt", "dir/tracked2.txt"}, "", nil),
			dirsWithRemainingFiles: []string{"dir"}, // tracked files remain in dir
			expectedRemovedFiles:   []string{"dir/new.txt"},
		},
		{
			testName: "directory node: staged-but-not-committed file (Tracked=false, HasStagedChanges=true) is left alone; purely untracked file is removed",
			nodes: []IFileNode{&testNode{
				path: "dir",
				children: []*testNode{
					// Staged new files: not removed from disk, but checked out in
					// case they also have unstaged changes on top (AM status).
					{path: "dir/staged-new1.txt", file: &models.File{Path: "dir/staged-new1.txt", Tracked: false, Added: true, HasStagedChanges: true}},
					{path: "dir/staged-new2.txt", file: &models.File{Path: "dir/staged-new2.txt", Tracked: false, Added: true, HasStagedChanges: true}},
					// Purely untracked file: removed from disk, not checked out.
					{path: "dir/untracked.txt", file: &models.File{Path: "dir/untracked.txt", Tracked: false, Added: true, HasStagedChanges: false}},
				},
			}},
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"checkout", "--", "dir/staged-new1.txt", "dir/staged-new2.txt"}, "", nil),
			dirsWithRemainingFiles: []string{"dir"}, // staged files remain in dir
			expectedRemovedFiles:   []string{"dir/untracked.txt"},
		},
		{
			testName: "file node: added and unstaged file is removed from disk",
			nodes: []IFileNode{&testNode{
				path: "new.txt",
				file: &models.File{Path: "new.txt", Added: true, HasStagedChanges: false},
			}},
			runner:               oscommands.NewFakeRunner(t),
			expectedRemovedFiles: []string{"new.txt"},
		},
		{
			testName: "files from multiple nodes are batched into a single checkout call",
			nodes: []IFileNode{
				&testNode{
					path: "dir1",
					children: []*testNode{
						{path: "dir1/tracked.txt", file: &models.File{Path: "dir1/tracked.txt", Tracked: true}},
						{path: "dir1/untracked.txt", file: &models.File{Path: "dir1/untracked.txt", Tracked: false}},
					},
				},
				&testNode{
					path: "dir2",
					children: []*testNode{
						{path: "dir2/tracked.txt", file: &models.File{Path: "dir2/tracked.txt", Tracked: true}},
						{path: "dir2/untracked.txt", file: &models.File{Path: "dir2/untracked.txt", Tracked: false}},
					},
				},
			},
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"checkout", "--", "dir1/tracked.txt", "dir2/tracked.txt"}, "", nil),
			dirsWithRemainingFiles: []string{"dir1", "dir2"}, // tracked files remain
			expectedRemovedFiles:   []string{"dir1/untracked.txt", "dir2/untracked.txt"},
		},
		{
			testName: "empty untracked directory is removed after its files are deleted",
			nodes: []IFileNode{&testNode{
				path: "newdir",
				children: []*testNode{
					{path: "newdir/a.txt", file: &models.File{Path: "newdir/a.txt", Tracked: false}},
					{path: "newdir/b.txt", file: &models.File{Path: "newdir/b.txt", Tracked: false}},
				},
			}},
			runner:               oscommands.NewFakeRunner(t),
			expectedRemovedFiles: []string{"newdir/a.txt", "newdir/b.txt"},
			expectedRemovedDirs:  []string{"newdir"},
		},
		{
			testName: "empty directory is NOT removed when individual file nodes are selected",
			nodes: []IFileNode{
				&testNode{path: "newdir/a.txt", file: &models.File{Path: "newdir/a.txt", Tracked: false}},
				&testNode{path: "newdir/b.txt", file: &models.File{Path: "newdir/b.txt", Tracked: false}},
			},
			runner:               oscommands.NewFakeRunner(t),
			expectedRemovedFiles: []string{"newdir/a.txt", "newdir/b.txt"},
			// newdir becomes empty but was not selected as a directory node, so it is not removed
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			var removedFiles []string
			removeFile := func(path string) error {
				removedFiles = append(removedFiles, path)
				return nil
			}
			isDirEmpty := func(path string) (bool, error) { return !lo.Contains(s.dirsWithRemainingFiles, path), nil }
			var removedDirs []string
			removeDir := func(path string) error {
				removedDirs = append(removedDirs, path)
				return nil
			}
			instance := buildWorkingTreeCommands(commonDeps{
				runner:     s.runner,
				removeFile: removeFile,
				isDirEmpty: isDirEmpty,
				removeDir:  removeDir,
			})
			assert.NoError(t, instance.DiscardUnstagedDirChanges(s.nodes))
			s.runner.CheckForMissingCalls()
			assert.Equal(t, s.expectedRemovedFiles, removedFiles)
			assert.Equal(t, s.expectedRemovedDirs, removedDirs)
		})
	}
}

func TestWorkingTreeDiscardAnyUnstagedFileChanges(t *testing.T) {
	type scenario struct {
		testName string
		runner   *oscommands.FakeCmdObjRunner
		test     func(error)
	}

	scenarios := []scenario{
		{
			testName: "valid case",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"checkout", "--", "."}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			instance := buildWorkingTreeCommands(commonDeps{runner: s.runner})
			s.test(instance.DiscardAnyUnstagedFileChanges())
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestWorkingTreeRemoveUntrackedFiles(t *testing.T) {
	type scenario struct {
		testName string
		runner   *oscommands.FakeCmdObjRunner
		test     func(error)
	}

	scenarios := []scenario{
		{
			testName: "valid case",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"clean", "-fd"}, "", nil),
			test: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			instance := buildWorkingTreeCommands(commonDeps{runner: s.runner})
			s.test(instance.RemoveUntrackedFiles())
			s.runner.CheckForMissingCalls()
		})
	}
}

func TestWorkingTreeResetHard(t *testing.T) {
	type scenario struct {
		testName string
		ref      string
		runner   *oscommands.FakeCmdObjRunner
		test     func(error)
	}

	scenarios := []scenario{
		{
			"valid case",
			"HEAD",
			oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"reset", "--hard", "HEAD"}, "", nil),
			func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			instance := buildWorkingTreeCommands(commonDeps{runner: s.runner})
			s.test(instance.ResetHard(s.ref))
		})
	}
}

func TestWorkingTreeCommands_AllRepoFiles(t *testing.T) {
	scenarios := []struct {
		name     string
		runner   *oscommands.FakeCmdObjRunner
		expected []string
	}{
		{
			name: "no files",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"ls-files", "-z"}, "", nil),
			expected: []string{},
		},
		{
			name: "two files",
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"ls-files", "-z"}, "dir/file1.txt\x00dir2/file2.go\x00", nil),
			expected: []string{"dir/file1.txt", "dir2/file2.go"},
		},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			instance := buildWorkingTreeCommands(commonDeps{runner: s.runner})
			result, err := instance.AllRepoFiles()
			assert.NoError(t, err)
			assert.Equal(t, s.expected, result)
		})
	}
}
