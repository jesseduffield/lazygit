package commands

import (
	"os/exec"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/secureexec"
	"github.com/stretchr/testify/assert"
)

// TestGitCommandGetStatusFiles is a function.
func TestGitCommandGetStatusFiles(t *testing.T) {
	type scenario struct {
		testName string
		command  func(string, ...string) *exec.Cmd
		test     func([]*models.File)
	}

	scenarios := []scenario{
		{
			"No files found",
			func(cmd string, args ...string) *exec.Cmd {
				return secureexec.Command("echo")
			},
			func(files []*models.File) {
				assert.Len(t, files, 0)
			},
		},
		{
			"Several files found",
			func(cmd string, args ...string) *exec.Cmd {
				return secureexec.Command(
					"printf",
					`MM file1.txt\0A  file3.txt\0AM file2.txt\0?? file4.txt\0UU file5.txt`,
				)
			},
			func(files []*models.File) {
				assert.Len(t, files, 5)

				expected := []*models.File{
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
						Type:                    "other",
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
						Type:                    "other",
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
						Type:                    "other",
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
						Type:                    "other",
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
						Type:                    "other",
						ShortStatus:             "UU",
					},
				}

				assert.EqualValues(t, expected, files)
			},
		},
		{
			"File with new line char",
			func(cmd string, args ...string) *exec.Cmd {
				return secureexec.Command(
					"printf",
					`MM a\nb.txt`,
				)
			},
			func(files []*models.File) {
				assert.Len(t, files, 1)

				expected := []*models.File{
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
						Type:                    "other",
						ShortStatus:             "MM",
					},
				}

				assert.EqualValues(t, expected, files)
			},
		},
		{
			"Renamed files",
			func(cmd string, args ...string) *exec.Cmd {
				return secureexec.Command(
					"printf",
					`R  after1.txt\0before1.txt\0RM after2.txt\0before2.txt`,
				)
			},
			func(files []*models.File) {
				assert.Len(t, files, 2)

				expected := []*models.File{
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
						Type:                    "other",
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
						Type:                    "other",
						ShortStatus:             "RM",
					},
				}

				assert.EqualValues(t, expected, files)
			},
		},
		{
			"File with arrow in name",
			func(cmd string, args ...string) *exec.Cmd {
				return secureexec.Command(
					"printf",
					`?? a -> b.txt`,
				)
			},
			func(files []*models.File) {
				assert.Len(t, files, 1)

				expected := []*models.File{
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
						Type:                    "other",
						ShortStatus:             "??",
					},
				}

				assert.EqualValues(t, expected, files)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			gitCmd := NewDummyGitCommand()
			gitCmd.OSCommand.Command = s.command

			s.test(gitCmd.GetStatusFiles(GetStatusFileOptions{}))
		})
	}
}
