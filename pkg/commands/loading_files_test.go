package commands

// import (
// 	"os/exec"
// 	"testing"

// 	"github.com/jesseduffield/lazygit/pkg/commands/models"
// 	"github.com/jesseduffield/lazygit/pkg/secureexec"
// 	"github.com/stretchr/testify/assert"
// )

// // TestGitCommandGetStatusFiles is a function.
// func TestGitCommandGetStatusFiles(t *testing.T) {
// 	type scenario struct {
// 		testName string
// 		command  func(string, ...string) *exec.Cmd
// 		test     func([]*models.File)
// 	}

// 	scenarios := []scenario{
// 		{
// 			"No files found",
// 			func(cmd string, args ...string) *exec.Cmd {
// 				return secureexec.Command("echo")
// 			},
// 			func(files []*models.File) {
// 				assert.Len(t, files, 0)
// 			},
// 		},
// 		{
// 			"Several files found",
// 			func(cmd string, args ...string) *exec.Cmd {
// 				return secureexec.Command(
// 					"echo",
// 					"MM file1.txt\nA  file3.txt\nAM file2.txt\n?? file4.txt\nUU file5.txt",
// 				)
// 			},
// 			func(files []*models.File) {
// 				assert.Len(t, files, 5)

// 				expected := []*models.File{
// 					{
// 						Name:                    "file1.txt",
// 						HasStagedChanges:        true,
// 						HasUnstagedChanges:      true,
// 						Tracked:                 true,
// 						Added:                   false,
// 						Deleted:                 false,
// 						HasMergeConflicts:       false,
// 						HasInlineMergeConflicts: false,
// 						DisplayString:           "MM file1.txt",
// 						Type:                    "other",
// 						ShortStatus:             "MM",
// 					},
// 					{
// 						Name:                    "file3.txt",
// 						HasStagedChanges:        true,
// 						HasUnstagedChanges:      false,
// 						Tracked:                 false,
// 						Added:                   true,
// 						Deleted:                 false,
// 						HasMergeConflicts:       false,
// 						HasInlineMergeConflicts: false,
// 						DisplayString:           "A  file3.txt",
// 						Type:                    "other",
// 						ShortStatus:             "A ",
// 					},
// 					{
// 						Name:                    "file2.txt",
// 						HasStagedChanges:        true,
// 						HasUnstagedChanges:      true,
// 						Tracked:                 false,
// 						Added:                   true,
// 						Deleted:                 false,
// 						HasMergeConflicts:       false,
// 						HasInlineMergeConflicts: false,
// 						DisplayString:           "AM file2.txt",
// 						Type:                    "other",
// 						ShortStatus:             "AM",
// 					},
// 					{
// 						Name:                    "file4.txt",
// 						HasStagedChanges:        false,
// 						HasUnstagedChanges:      true,
// 						Tracked:                 false,
// 						Added:                   true,
// 						Deleted:                 false,
// 						HasMergeConflicts:       false,
// 						HasInlineMergeConflicts: false,
// 						DisplayString:           "?? file4.txt",
// 						Type:                    "other",
// 						ShortStatus:             "??",
// 					},
// 					{
// 						Name:                    "file5.txt",
// 						HasStagedChanges:        false,
// 						HasUnstagedChanges:      true,
// 						Tracked:                 true,
// 						Added:                   false,
// 						Deleted:                 false,
// 						HasMergeConflicts:       true,
// 						HasInlineMergeConflicts: true,
// 						DisplayString:           "UU file5.txt",
// 						Type:                    "other",
// 						ShortStatus:             "UU",
// 					},
// 				}

// 				assert.EqualValues(t, expected, files)
// 			},
// 		},
// 	}

// 	for _, s := range scenarios {
// 		t.Run(s.testName, func(t *testing.T) {
// 			gitCmd := NewDummyGit()
// 			gitCmd.GetOSCommand().Command = s.command

// 			s.test(gitCmd.GetStatusFiles(GetStatusFileOptions{}))
// 		})
// 	}
// }
