package git_commands

import (
	"encoding/binary"
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
			}

			assert.EqualValues(t, s.expectedFiles, loader.GetStatusFiles(GetStatusFileOptions{}))
		})
	}
}

type FakeFileLoaderConfig struct {
	showUntrackedFiles string
}

func (self *FakeFileLoaderConfig) GetShowUntrackedFiles() string {
	return self.showUntrackedFiles
}

// writeFakeIndexHeader writes a minimal git index file (just the 12-byte header:
// "DIRC" + version + big-endian entry count) at <worktreeDir>/.git/index, which
// is where MockRepoPaths expects it.
func writeFakeIndexHeader(t *testing.T, worktreeDir string, entryCount uint32) {
	t.Helper()
	gitDir := filepath.Join(worktreeDir, ".git")
	if err := os.MkdirAll(gitDir, 0o700); err != nil {
		t.Fatalf("failed to create .git dir: %v", err)
	}
	header := make([]byte, 12)
	copy(header, "DIRC")
	binary.BigEndian.PutUint32(header[4:8], 2) // index format version
	binary.BigEndian.PutUint32(header[8:12], entryCount)
	if err := os.WriteFile(filepath.Join(gitDir, "index"), header, 0o600); err != nil {
		t.Fatalf("failed to write index: %v", err)
	}
}

// TestFileLoaderUntrackedFilesArg covers how GetStatusFiles chooses the
// --untracked-files argument: an explicit status.showUntrackedFiles setting is
// always honored, an unset setting defaults to "all" but downgrades to "normal"
// in large repos, and the ForceShowUntracked filter forces "all" regardless.
func TestFileLoaderUntrackedFilesArg(t *testing.T) {
	type scenario struct {
		testName string
		// git config status.showUntrackedFiles ("" means unset)
		showUntrackedFiles string
		// number of entries to record in the fake index header
		indexEntryCount uint32
		// if false, no index file is written (simulates an unreadable/missing index)
		writeIndex         bool
		forceShowUntracked bool
		expectedArg        string
	}

	scenarios := []scenario{
		{
			testName:        "unset config, large repo -> normal",
			indexEntryCount: untrackedFilesAllMaxTrackedFiles + 1,
			writeIndex:      true,
			expectedArg:     "--untracked-files=normal",
		},
		{
			testName:        "unset config, small repo -> all",
			indexEntryCount: 100,
			writeIndex:      true,
			expectedArg:     "--untracked-files=all",
		},
		{
			testName:        "unset config, count exactly at threshold -> all (strict >)",
			indexEntryCount: untrackedFilesAllMaxTrackedFiles,
			writeIndex:      true,
			expectedArg:     "--untracked-files=all",
		},
		{
			testName:    "unset config, missing/unreadable index -> all (fail-safe)",
			writeIndex:  false,
			expectedArg: "--untracked-files=all",
		},
		{
			testName:           "explicit 'all' honored even in large repo",
			showUntrackedFiles: "all",
			indexEntryCount:    untrackedFilesAllMaxTrackedFiles + 1,
			writeIndex:         true,
			expectedArg:        "--untracked-files=all",
		},
		{
			testName:           "explicit 'no' honored even in large repo",
			showUntrackedFiles: "no",
			indexEntryCount:    untrackedFilesAllMaxTrackedFiles + 1,
			writeIndex:         true,
			expectedArg:        "--untracked-files=no",
		},
		{
			testName:           "explicit 'normal' honored in small repo",
			showUntrackedFiles: "normal",
			indexEntryCount:    100,
			writeIndex:         true,
			expectedArg:        "--untracked-files=normal",
		},
		{
			testName:           "ForceShowUntracked forces all in large repo",
			indexEntryCount:    untrackedFilesAllMaxTrackedFiles + 1,
			writeIndex:         true,
			forceShowUntracked: true,
			expectedArg:        "--untracked-files=all",
		},
		{
			testName:           "ForceShowUntracked overrides explicit 'no'",
			showUntrackedFiles: "no",
			indexEntryCount:    100,
			writeIndex:         true,
			forceShowUntracked: true,
			expectedArg:        "--untracked-files=all",
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			worktreeDir := t.TempDir()
			if s.writeIndex {
				writeFakeIndexHeader(t, worktreeDir, s.indexEntryCount)
			}

			runner := oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"status", s.expectedArg, "--porcelain", "-z", "--find-renames=50%"}, "", nil)
			cmd := oscommands.NewDummyCmdObjBuilder(runner)

			userConfig := &config.UserConfig{}
			userConfig.Git.RenameSimilarityThreshold = 50

			loader := &FileLoader{
				GitCommon: buildGitCommon(commonDeps{
					appState:   &config.AppState{},
					userConfig: userConfig,
					repoPaths:  MockRepoPaths(worktreeDir),
				}),
				cmd:         cmd,
				config:      &FakeFileLoaderConfig{showUntrackedFiles: s.showUntrackedFiles},
				getFileType: func(string) string { return "file" },
			}

			loader.GetStatusFiles(GetStatusFileOptions{ForceShowUntracked: s.forceShowUntracked})
			runner.CheckForMissingCalls()
		})
	}
}

func TestTrackedFileCountFromIndex(t *testing.T) {
	validHeader := func(count uint32) []byte {
		b := make([]byte, 12)
		copy(b, "DIRC")
		binary.BigEndian.PutUint32(b[4:8], 2)
		binary.BigEndian.PutUint32(b[8:12], count)
		return b
	}
	writeIndex := func(t *testing.T, content []byte) string {
		t.Helper()
		path := filepath.Join(t.TempDir(), "index")
		if err := os.WriteFile(path, content, 0o600); err != nil {
			t.Fatalf("failed to write index: %v", err)
		}
		return path
	}

	t.Run("valid header returns entry count", func(t *testing.T) {
		count, ok := trackedFileCountFromIndex(writeIndex(t, validHeader(123456)))
		assert.True(t, ok)
		assert.Equal(t, 123456, count)
	})

	t.Run("trailing entry bytes don't affect the count", func(t *testing.T) {
		count, ok := trackedFileCountFromIndex(writeIndex(t, append(validHeader(42), []byte("trailing entry data")...)))
		assert.True(t, ok)
		assert.Equal(t, 42, count)
	})

	t.Run("missing file returns not-ok", func(t *testing.T) {
		count, ok := trackedFileCountFromIndex(filepath.Join(t.TempDir(), "does-not-exist"))
		assert.False(t, ok)
		assert.Equal(t, 0, count)
	})

	t.Run("truncated header returns not-ok", func(t *testing.T) {
		count, ok := trackedFileCountFromIndex(writeIndex(t, []byte("DIRC")))
		assert.False(t, ok)
		assert.Equal(t, 0, count)
	})

	t.Run("wrong signature returns not-ok", func(t *testing.T) {
		bad := validHeader(999)
		copy(bad, "XXXX")
		count, ok := trackedFileCountFromIndex(writeIndex(t, bad))
		assert.False(t, ok)
		assert.Equal(t, 0, count)
	})
}
