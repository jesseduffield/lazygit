package git_commands

import (
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
