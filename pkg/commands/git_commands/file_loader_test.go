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
		testName            string
		similarityThreshold int
		runner              oscommands.ICmdObjRunner
		expectedFiles       []*models.File
	}

	scenarios := []scenario{
		{
			"No files found",
			50,
			oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"status", "--untracked-files=yes", "--porcelain", "-z", "--find-renames=50%"}, "", nil),
			[]*models.File{},
		},
		{
			"Several files found",
			50,
			oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"status", "--untracked-files=yes", "--porcelain", "-z", "--find-renames=50%"},
					"MM file1.txt\x00A  file3.txt\x00AM file2.txt\x00?? file4.txt\x00UU file5.txt",
					nil,
				),
			[]*models.File{
				{
					Name:                    "file1.txt",
					HasStagedChanges:        true,
					HasUnstagedChanges:      true,
					Tracked:                 true,
					Added:                   false,
					Deleted:                 false,
					HasMergeConflicts:       false,
					HasInlineMergeConflicts: false,
					DisplayString:           "MM file1.txt",
					ShortStatus:             "MM",
				},
				{
					Name:                    "file3.txt",
					HasStagedChanges:        true,
					HasUnstagedChanges:      false,
					Tracked:                 false,
					Added:                   true,
					Deleted:                 false,
					HasMergeConflicts:       false,
					HasInlineMergeConflicts: false,
					DisplayString:           "A  file3.txt",
					ShortStatus:             "A ",
				},
				{
					Name:                    "file2.txt",
					HasStagedChanges:        true,
					HasUnstagedChanges:      true,
					Tracked:                 false,
					Added:                   true,
					Deleted:                 false,
					HasMergeConflicts:       false,
					HasInlineMergeConflicts: false,
					DisplayString:           "AM file2.txt",
					ShortStatus:             "AM",
				},
				{
					Name:                    "file4.txt",
					HasStagedChanges:        false,
					HasUnstagedChanges:      true,
					Tracked:                 false,
					Added:                   true,
					Deleted:                 false,
					HasMergeConflicts:       false,
					HasInlineMergeConflicts: false,
					DisplayString:           "?? file4.txt",
					ShortStatus:             "??",
				},
				{
					Name:                    "file5.txt",
					HasStagedChanges:        false,
					HasUnstagedChanges:      true,
					Tracked:                 true,
					Added:                   false,
					Deleted:                 false,
					HasMergeConflicts:       true,
					HasInlineMergeConflicts: true,
					DisplayString:           "UU file5.txt",
					ShortStatus:             "UU",
				},
			},
		},
		{
			"File with new line char",
			50,
			oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"status", "--untracked-files=yes", "--porcelain", "-z", "--find-renames=50%"}, "MM a\nb.txt", nil),
			[]*models.File{
				{
					Name:                    "a\nb.txt",
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
			"Renamed files",
			50,
			oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"status", "--untracked-files=yes", "--porcelain", "-z", "--find-renames=50%"},
					"R  after1.txt\x00before1.txt\x00RM after2.txt\x00before2.txt",
					nil,
				),
			[]*models.File{
				{
					Name:                    "after1.txt",
					PreviousName:            "before1.txt",
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
					Name:                    "after2.txt",
					PreviousName:            "before2.txt",
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
			"File with arrow in name",
			50,
			oscommands.NewFakeRunner(t).
				ExpectGitArgs([]string{"status", "--untracked-files=yes", "--porcelain", "-z", "--find-renames=50%"},
					`?? a -> b.txt`,
					nil,
				),
			[]*models.File{
				{
					Name:                    "a -> b.txt",
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
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			cmd := oscommands.NewDummyCmdObjBuilder(s.runner)

			appState := &config.AppState{}
			appState.RenameSimilarityThreshold = s.similarityThreshold

			loader := &FileLoader{
				GitCommon:   buildGitCommon(commonDeps{appState: appState}),
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
