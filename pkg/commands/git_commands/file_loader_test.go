package git_commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestFileGetStatusFiles(t *testing.T) {
	type scenario struct {
		testName               string
		similarityThreshold    int
		runner                 oscommands.ICmdObjRunner
		showNumstatInFilesView bool
		expectedFiles          []*models.File
	}

	scenarios := []scenario{
		{
			testName:            "No files found",
			similarityThreshold: 50,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"status", "--untracked-files=yes", "--porcelain", "-z", "--find-renames=50%"}, "", nil),
			expectedFiles: []*models.File{},
		},
		{
			testName:            "Several files found",
			similarityThreshold: 50,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"status", "--untracked-files=yes", "--porcelain", "-z", "--find-renames=50%"},
					"MM file1.txt\x00A  file3.txt\x00AM file2.txt\x00?? file4.txt\x00UU file5.txt",
					nil,
				).
				ExpectGitArgs([]string{"diff", "--numstat", "-z", "HEAD"},
					"4\t1\tfile1.txt\x001\t0\tfile2.txt\x002\t2\tfile3.txt\x000\t2\tfile4.txt\x002\t2\tfile5.txt",
					nil,
				).
				ExpectGitArgs([]string{"check-attr", "-z", "--stdin", "conflict-marker-size"},
					"file5.txt\x00conflict-marker-size\x00unspecified\x00",
					nil,
				),
			showNumstatInFilesView: true,
			expectedFiles: []*models.File{
				{
					Path:                    "file1.txt",
					HasStagedChanges:        true,
					HasUnstagedChanges:      true,
					Tracked:                 true,
					Added:                   false,
					Deleted:                 false,
					HasMergeConflicts:       false,
					HasInlineMergeConflicts: false,
					DisplayString:           "MM file1.txt",
					ShortStatus:             "MM",
					LinesAdded:              4,
					LinesDeleted:            1,
				},
				{
					Path:                    "file3.txt",
					HasStagedChanges:        true,
					HasUnstagedChanges:      false,
					Tracked:                 false,
					Added:                   true,
					Deleted:                 false,
					HasMergeConflicts:       false,
					HasInlineMergeConflicts: false,
					DisplayString:           "A  file3.txt",
					ShortStatus:             "A ",
					LinesAdded:              2,
					LinesDeleted:            2,
				},
				{
					Path:                    "file2.txt",
					HasStagedChanges:        true,
					HasUnstagedChanges:      true,
					Tracked:                 false,
					Added:                   true,
					Deleted:                 false,
					HasMergeConflicts:       false,
					HasInlineMergeConflicts: false,
					DisplayString:           "AM file2.txt",
					ShortStatus:             "AM",
					LinesAdded:              1,
					LinesDeleted:            0,
				},
				{
					Path:                    "file4.txt",
					HasStagedChanges:        false,
					HasUnstagedChanges:      true,
					Tracked:                 false,
					Added:                   true,
					Deleted:                 false,
					HasMergeConflicts:       false,
					HasInlineMergeConflicts: false,
					DisplayString:           "?? file4.txt",
					ShortStatus:             "??",
					LinesAdded:              0,
					LinesDeleted:            2,
				},
				{
					Path:                    "file5.txt",
					HasStagedChanges:        false,
					HasUnstagedChanges:      true,
					Tracked:                 true,
					Added:                   false,
					Deleted:                 false,
					HasMergeConflicts:       true,
					HasInlineMergeConflicts: true,
					DisplayString:           "UU file5.txt",
					ShortStatus:             "UU",
					LinesAdded:              2,
					LinesDeleted:            2,
				},
			},
		},
		{
			testName:            "Conflicted files with a conflict-marker-size attribute",
			similarityThreshold: 50,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"status", "--untracked-files=yes", "--porcelain", "-z", "--find-renames=50%"},
					"UU file1.txt\x00UU file2.txt\x00UU file3.txt\x00 M file4.txt",
					nil,
				).
				ExpectGitArgs([]string{"check-attr", "-z", "--stdin", "conflict-marker-size"},
					"file1.txt\x00conflict-marker-size\x0032\x00"+
						"file2.txt\x00conflict-marker-size\x00unspecified\x00"+
						"file3.txt\x00conflict-marker-size\x00nonsense\x00",
					nil,
				),
			expectedFiles: []*models.File{
				{
					Path:                    "file1.txt",
					HasUnstagedChanges:      true,
					Tracked:                 true,
					HasMergeConflicts:       true,
					HasInlineMergeConflicts: true,
					ConflictMarkerSize:      32,
					DisplayString:           "UU file1.txt",
					ShortStatus:             "UU",
				},
				{
					Path:                    "file2.txt",
					HasUnstagedChanges:      true,
					Tracked:                 true,
					HasMergeConflicts:       true,
					HasInlineMergeConflicts: true,
					DisplayString:           "UU file2.txt",
					ShortStatus:             "UU",
				},
				{
					Path:                    "file3.txt",
					HasUnstagedChanges:      true,
					Tracked:                 true,
					HasMergeConflicts:       true,
					HasInlineMergeConflicts: true,
					DisplayString:           "UU file3.txt",
					ShortStatus:             "UU",
				},
				{
					Path:               "file4.txt",
					HasUnstagedChanges: true,
					Tracked:            true,
					DisplayString:      " M file4.txt",
					ShortStatus:        " M",
				},
			},
		},
		{
			testName:            "File with new line char",
			similarityThreshold: 50,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"status", "--untracked-files=yes", "--porcelain", "-z", "--find-renames=50%"}, "MM a\nb.txt", nil),
			expectedFiles: []*models.File{
				{
					Path:                    "a\nb.txt",
					HasStagedChanges:        true,
					HasUnstagedChanges:      true,
					Tracked:                 true,
					Added:                   false,
					Deleted:                 false,
					HasMergeConflicts:       false,
					HasInlineMergeConflicts: false,
					DisplayString:           "MM a\nb.txt",
					ShortStatus:             "MM",
				},
			},
		},
		{
			testName:            "Directory that git doesn't descend into",
			similarityThreshold: 50,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"status", "--untracked-files=yes", "--porcelain", "-z", "--find-renames=50%"}, "?? dir/nested/", nil),
			expectedFiles: []*models.File{
				{
					Path:                    "dir/nested",
					HasStagedChanges:        false,
					HasUnstagedChanges:      true,
					Tracked:                 false,
					Added:                   true,
					Deleted:                 false,
					HasMergeConflicts:       false,
					HasInlineMergeConflicts: false,
					DisplayString:           "?? dir/nested/",
					ShortStatus:             "??",
					DirectoryKind:           models.NestedRepository,
				},
			},
		},
		{
			testName:            "Renamed files",
			similarityThreshold: 50,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"status", "--untracked-files=yes", "--porcelain", "-z", "--find-renames=50%"},
					"R  after1.txt\x00before1.txt\x00RM after2.txt\x00before2.txt",
					nil,
				),
			expectedFiles: []*models.File{
				{
					Path:                    "after1.txt",
					PreviousPath:            "before1.txt",
					HasStagedChanges:        true,
					HasUnstagedChanges:      false,
					Tracked:                 true,
					Added:                   false,
					Deleted:                 false,
					HasMergeConflicts:       false,
					HasInlineMergeConflicts: false,
					DisplayString:           "R  before1.txt -> after1.txt",
					ShortStatus:             "R ",
				},
				{
					Path:                    "after2.txt",
					PreviousPath:            "before2.txt",
					HasStagedChanges:        true,
					HasUnstagedChanges:      true,
					Tracked:                 true,
					Added:                   false,
					Deleted:                 false,
					HasMergeConflicts:       false,
					HasInlineMergeConflicts: false,
					DisplayString:           "RM before2.txt -> after2.txt",
					ShortStatus:             "RM",
				},
			},
		},
		{
			testName:            "File with arrow in name",
			similarityThreshold: 50,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"status", "--untracked-files=yes", "--porcelain", "-z", "--find-renames=50%"},
					`?? a -> b.txt`,
					nil,
				),
			expectedFiles: []*models.File{
				{
					Path:                    "a -> b.txt",
					HasStagedChanges:        false,
					HasUnstagedChanges:      true,
					Tracked:                 false,
					Added:                   true,
					Deleted:                 false,
					HasMergeConflicts:       false,
					HasInlineMergeConflicts: false,
					DisplayString:           "?? a -> b.txt",
					ShortStatus:             "??",
				},
			},
		},
		{
			testName:            "Copied files",
			similarityThreshold: 50,
			runner: oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"status", "--untracked-files=yes", "--porcelain", "-z", "--find-renames=50%"},
					"C  copy1.txt\x00original.txt\x00CM copy2.txt\x00original.txt",
					nil,
				),
			expectedFiles: []*models.File{
				{
					Path:                    "copy1.txt",
					PreviousPath:            "original.txt",
					HasStagedChanges:        true,
					HasUnstagedChanges:      false,
					Tracked:                 true,
					Added:                   false,
					Deleted:                 false,
					HasMergeConflicts:       false,
					HasInlineMergeConflicts: false,
					DisplayString:           "C  original.txt -> copy1.txt",
					ShortStatus:             "C ",
				},
				{
					Path:                    "copy2.txt",
					PreviousPath:            "original.txt",
					HasStagedChanges:        true,
					HasUnstagedChanges:      true,
					Tracked:                 true,
					Added:                   false,
					Deleted:                 false,
					HasMergeConflicts:       false,
					HasInlineMergeConflicts: false,
					DisplayString:           "CM original.txt -> copy2.txt",
					ShortStatus:             "CM",
				},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			cmd := oscommands.NewDummyCmdObjBuilder(s.runner)

			userConfig := &config.UserConfig{}
			userConfig.Gui.ShowNumstatInFilesView = s.showNumstatInFilesView
			userConfig.Git.RenameSimilarityThreshold = s.similarityThreshold

			loader := &FileLoader{
				GitCommon:   buildGitCommon(commonDeps{appState: &config.AppState{}, userConfig: userConfig}),
				cmd:         cmd,
				config:      &FakeFileLoaderConfig{showUntrackedFiles: "yes"},
				getFileType: func(string) string { return "file" },
				getDirectoryKind: func(string, string) models.DirectoryKind {
					return models.NestedRepository
				},
			}

			assert.EqualValues(t, s.expectedFiles, loader.GetStatusFiles(GetStatusFileOptions{}))
		})
	}
}

func TestDirectoryKind(t *testing.T) {
	repoDir := t.TempDir()
	worktreeGitDir := filepath.Join(repoDir, ".git")

	writeGitFile := func(dir string, gitDir string) {
		assert.NoError(t, os.MkdirAll(dir, 0o755))
		assert.NoError(t, os.WriteFile(filepath.Join(dir, ".git"), []byte("gitdir: "+gitDir+"\n"), 0o644))
	}

	plainDir := filepath.Join(repoDir, "plain")
	assert.NoError(t, os.MkdirAll(plainDir, 0o755))

	nestedRepoDir := filepath.Join(repoDir, "nested")
	assert.NoError(t, os.MkdirAll(filepath.Join(nestedRepoDir, ".git"), 0o755))

	removedSubmoduleDir := filepath.Join(repoDir, "modules", "sub")
	writeGitFile(removedSubmoduleDir, "../../.git/modules/sub")

	otherReposSubmoduleDir := filepath.Join(repoDir, "other-sub")
	writeGitFile(otherReposSubmoduleDir, filepath.Join(t.TempDir(), ".git", "modules", "sub"))

	worktreeDir := filepath.Join(repoDir, "worktree")
	writeGitFile(worktreeDir, filepath.Join(worktreeGitDir, "worktrees", "worktree"))

	scenarios := []struct {
		name     string
		dirPath  string
		expected models.DirectoryKind
	}{
		{"directory without a .git", plainDir, models.PlainDirectory},
		{"repository with a .git directory", nestedRepoDir, models.NestedRepository},
		{"submodule of this repo", removedSubmoduleDir, models.RemovedSubmodule},
		{"submodule of another repo", otherReposSubmoduleDir, models.NestedRepository},
		{"worktree of this repo", worktreeDir, models.NestedRepository},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expected, directoryKind(s.dirPath, worktreeGitDir))
		})
	}
}

type FakeFileLoaderConfig struct {
	showUntrackedFiles string
}

func (self *FakeFileLoaderConfig) GetShowUntrackedFiles() string {
	return self.showUntrackedFiles
}
